package shared

const (
	OpenKey                     = "ipcOpen"
	GetInstancesKey             = "ipcGetInstances"
	PasswordGetterKey           = "ipcPasswordGetter"
	HostKeyValidatorKey         = "ipcHostKeyValidator"
	ForwardFromLocalToRemoteKey = "ipcForwardFromLocalToRemote"
)

type Instance struct {
	ID      string
	Name    string
	Tunnels []Tunnel
}

type Tunnel struct {
	ID            string
	LocalAddress  string
	RemoteAddress string
}
