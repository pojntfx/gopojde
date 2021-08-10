package shared

const (
	OpenKey             = "ipcOpen"
	GetInstancesKey     = "ipcGetInstances"
	CreateSSHConnection = "ipcCreateSSHConnection"
)

type Instance struct {
	ID   string
	Name string
}
