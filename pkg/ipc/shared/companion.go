package shared

const (
	OpenKey         = "ipcOpen"
	GetInstancesKey = "ipcGetInstances"
)

type Instance struct {
	ID   string
	Name string
}
