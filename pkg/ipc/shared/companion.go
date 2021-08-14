package shared

const (
	OpenKey             = "ipcOpen"
	GetInstancesKey     = "ipcGetInstances"
	PasswordGetterKey   = "ipcPasswordGetter"
	HostKeyValidatorKey = "ipcHostKeyValidator"
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
