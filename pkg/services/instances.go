package services

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/pojntfx/gopojde/pkg/api/proto/v1"
	"github.com/pojntfx/gopojde/pkg/orchestration"
	"github.com/pojntfx/gopojde/pkg/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

//go:generate sh -c "mkdir -p ../api/proto/v1 && protoc --go_out=paths=source_relative,plugins=grpc:../api/proto/v1 -I=../../api/proto/v1 ../../api/proto/v1/*.proto"

type InstancesService struct {
	api.UnimplementedInstancesServiceServer

	instancesManager *orchestration.InstancesManager
}

func NewInstancesService(instancesManager *orchestration.InstancesManager) *InstancesService {
	return &InstancesService{
		instancesManager: instancesManager,
	}
}

func (s *InstancesService) GetInstances(ctx context.Context, _ *empty.Empty) (*api.InstanceSummariesMessage, error) {
	instances, err := s.instancesManager.GetInstances(ctx)
	if err != nil {
		return &api.InstanceSummariesMessage{}, err
	}

	out := []*api.InstanceSummaryMessage{}
	for _, instance := range instances {
		ports := []*api.PortMessage{}
		for _, port := range instance.Ports {
			ports = append(ports, &api.PortMessage{
				Service: port.Service,
				Port:    port.Port,
			})
		}

		out = append(out, &api.InstanceSummaryMessage{
			InstanceID: &api.InstanceIDMessage{
				Name: instance.Name,
			},
			Ports:  ports,
			Status: instance.Status,
		})
	}

	return &api.InstanceSummariesMessage{
		Instances: out,
	}, nil
}

func (s *InstancesService) GetLogs(req *api.InstanceIDMessage, stream api.InstancesService_GetLogsServer) error {
	var fatalError error
	ctx, _cancel := context.WithCancel(stream.Context())
	cancel := func(err error) {
		fatalError = err

		_cancel()
	}

	stdoutChan, stderrChan := make(chan []byte), make(chan []byte)
	defer close(stdoutChan)
	defer close(stderrChan)

	go func() {
		for chunk := range stdoutChan {
			if err := stream.Send(&api.ShellOutputMessage{
				Stdout: chunk,
			}); err != nil {
				cancel(err)

				return
			}
		}
	}()

	go func() {
		for chunk := range stderrChan {
			if err := stream.Send(&api.ShellOutputMessage{
				Stderr: chunk,
			}); err != nil {
				cancel(err)

				return
			}
		}
	}()

	go s.instancesManager.GetLogs(ctx, cancel, req.GetName(), stdoutChan, stderrChan)

	<-ctx.Done()

	return fatalError
}

func (s *InstancesService) StartInstance(ctx context.Context, req *api.InstanceIDMessage) (*empty.Empty, error) {
	return &emptypb.Empty{}, s.instancesManager.StartInstance(ctx, req.GetName(), false)
}

func (s *InstancesService) StopInstance(ctx context.Context, req *api.InstanceIDMessage) (*empty.Empty, error) {
	return &emptypb.Empty{}, s.instancesManager.StopInstance(ctx, req.GetName())
}

func (s *InstancesService) RestartInstance(ctx context.Context, req *api.InstanceIDMessage) (*empty.Empty, error) {
	return &emptypb.Empty{}, s.instancesManager.RestartInstance(ctx, req.GetName())
}

func (s *InstancesService) RemoveInstance(ctx context.Context, req *api.InstanceRemovalOptionsMessage) (*empty.Empty, error) {
	return &emptypb.Empty{}, s.instancesManager.RemoveInstance(ctx, req.GetInstanceID().GetName(), orchestration.InstanceRemovalOptions{
		Customizations: req.GetCustomizations(),
		DEBCache:       req.GetDEBCache(),
		Preferences:    req.GetPreferences(),
		UserData:       req.GetUserData(),
		Transfer:       req.GetTransfer(),
	}, false)
}

func (s *InstancesService) GetCACert(ctx context.Context, _ *empty.Empty) (*api.CACertMessage, error) {
	cert, err := s.instancesManager.GetCACert(ctx)
	if err != nil {
		return &api.CACertMessage{}, err
	}

	return &api.CACertMessage{
		Content: cert,
	}, nil
}

func (s *InstancesService) ResetCA(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &emptypb.Empty{}, s.instancesManager.ResetCA(ctx)
}

func (s *InstancesService) GetShell(stream api.InstancesService_GetShellServer) error {
	var fatalError error
	ctx, _cancel := context.WithCancel(stream.Context())
	cancel := func(err error) {
		fatalError = err

		_cancel()
	}

	go func() {
		stdinChan, stdoutChan, stderrChan := make(chan []byte), make(chan []byte), make(chan []byte)
		defer close(stdinChan)
		defer close(stdoutChan)
		defer close(stderrChan)

		open := false
		for {
			msg, err := stream.Recv()
			if err != nil {
				cancel(err)

				return
			}

			if !open {
				go func() {
					for chunk := range stdoutChan {
						if err := stream.Send(&api.ShellOutputMessage{
							Stdout: chunk,
						}); err != nil {
							cancel(err)

							return
						}
					}
				}()

				go func() {
					for chunk := range stderrChan {
						if err := stream.Send(&api.ShellOutputMessage{
							Stderr: chunk,
						}); err != nil {
							cancel(err)

							return
						}
					}
				}()

				go s.instancesManager.GetShell(ctx, cancel, msg.GetInstanceID().GetName(), stdinChan, stdoutChan, stderrChan)

				open = true
			}

			select {
			case <-ctx.Done():
				return
			default:
				stdinChan <- msg.GetStdin()
			}
		}
	}()

	<-ctx.Done()

	return fatalError
}

func (s *InstancesService) ApplyInstance(req *api.InstanceConfigurationMessage, stream api.InstancesService_ApplyInstanceServer) error {
	var fatalError error
	ctx, _cancel := context.WithCancel(stream.Context())
	cancel := func(err error) {
		fatalError = err

		_cancel()
	}

	stdoutChan, stderrChan, statusChan := make(chan []byte), make(chan []byte), make(chan string)
	defer close(stdoutChan)
	defer close(stderrChan)
	defer close(statusChan)

	go func() {
		for chunk := range stdoutChan {
			if err := stream.Send(&api.ShellOutputMessage{
				Stdout: chunk,
			}); err != nil {
				cancel(err)

				return
			}
		}
	}()

	go func() {
		for chunk := range stderrChan {
			if err := stream.Send(&api.ShellOutputMessage{
				Stderr: chunk,
			}); err != nil {
				cancel(err)

				return
			}
		}
	}()

	go func() {
		for status := range statusChan {
			if err := stream.Send(&api.ShellOutputMessage{
				Status: status,
			}); err != nil {
				cancel(err)

				return
			}
		}
	}()

	go s.instancesManager.ApplyInstance(
		ctx,
		cancel,
		req.GetInstanceContainerOptions().GetName(),
		stdoutChan,
		stderrChan,
		statusChan,
		orchestration.InstanceContainerOptions{
			StartPort:       req.GetInstanceContainerOptions().GetStartPort(),
			Isolate:         req.GetInstanceContainerOptions().GetIsolate(),
			Privileged:      req.GetInstanceContainerOptions().GetPrivileged(),
			Recreate:        req.GetInstanceContainerOptions().GetRecreate(),
			PullLatestImage: req.GetInstanceContainerOptions().GetPullLatestImage(),
		},
		orchestration.InstanceAuthentication{
			RootPassword: req.GetInstanceAuthentication().GetRootPassword(),
			UserPassword: req.GetInstanceAuthentication().GetUserPassword(),
		},
		orchestration.InstanceOptions{
			UserName:     req.GetInstanceOptions().GetUserName(),
			UserEmail:    req.GetInstanceOptions().GetUserEmail(),
			UserFullName: req.GetInstanceOptions().GetUserFullName(),
			SSHKeyURL:    req.GetInstanceOptions().GetSSHKeyURL(),

			AdditionalIPs:     req.GetInstanceOptions().GetAdditionalIPs(),
			AdditionalDomains: req.GetInstanceOptions().GetAdditionalDomains(),

			EnabledModules:  req.GetInstanceOptions().GetEnabledModules(),
			EnabledServices: req.GetInstanceOptions().GetEnabledServices(),
		})

	<-ctx.Done()

	return fatalError
}

func (s *InstancesService) GetInstanceConfig(ctx context.Context, req *api.InstanceIDMessage) (*api.InstanceOptionsMessage, error) {
	cfg, err := s.instancesManager.GetInstanceConfig(ctx, req.GetName())
	if err != nil {
		return &api.InstanceOptionsMessage{}, err
	}

	return &api.InstanceOptionsMessage{
		UserName:     cfg.UserEmail,
		UserEmail:    cfg.UserEmail,
		UserFullName: cfg.UserFullName,
		SSHKeyURL:    cfg.SSHKeyURL,

		AdditionalIPs:     cfg.AdditionalIPs,
		AdditionalDomains: cfg.AdditionalDomains,

		EnabledModules:  cfg.EnabledModules,
		EnabledServices: cfg.EnabledServices,
	}, nil
}

func (s *InstancesService) GetSSHKeys(ctx context.Context, req *api.InstanceIDMessage) (*api.SSHKeysMessage, error) {
	sshKeyContents, err := s.instancesManager.GetSSHKeys(ctx, req.GetName())
	if err != nil {
		return &api.SSHKeysMessage{}, err
	}

	sshKeys := []*api.SSHKeyMessage{}
	for _, sshKeyContents := range sshKeyContents {
		sshKeys = append(sshKeys, &api.SSHKeyMessage{
			SSHKeyID: &api.SSHKeyIDMessage{
				InstanceID: req,
				Hash:       util.GetSHA512Hash(sshKeyContents),
			},
			Content: sshKeyContents,
		})
	}

	return &api.SSHKeysMessage{
		SSHKeys: sshKeys,
	}, nil
}

func (s *InstancesService) AddSSHKey(ctx context.Context, req *api.SSHKeyAdditionMessage) (*api.SSHKeyMessage, error) {
	return &api.SSHKeyMessage{
		SSHKeyID: &api.SSHKeyIDMessage{
			InstanceID: req.GetInstanceID(),
			Hash:       util.GetSHA512Hash(req.GetContent()),
		},
		Content: req.GetContent(),
	}, s.instancesManager.AddSSHKey(ctx, req.GetInstanceID().GetName(), req.GetContent())
}

func (s *InstancesService) RemoveSSHKey(ctx context.Context, req *api.SSHKeyIDMessage) (*api.SSHKeyMessage, error) {
	sshKey, err := s.instancesManager.RemoveSSHKey(ctx, req.GetInstanceID().GetName(), req.GetHash())
	if err != nil {
		return &api.SSHKeyMessage{}, err
	}

	return &api.SSHKeyMessage{
		SSHKeyID: req,
		Content:  sshKey,
	}, nil
}
