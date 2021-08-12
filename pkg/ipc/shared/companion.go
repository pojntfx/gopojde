package shared

const (
	OpenKey                = "ipcOpen"
	GetInstancesKey        = "ipcGetInstances"
	CreateSSHConnectionKey = "ipcCreateSSHConnection"
	PasswordGetterKey      = "ipcPasswordGetter"
	HostKeyValidatorKey    = "ipcHostKeyValidator"
)

type Instance struct {
	ID   string
	Name string
}
