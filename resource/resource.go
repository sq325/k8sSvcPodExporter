package resource

type Resource interface {
	Name() string
	Namespace() string
}

type Factor interface {
	CmdStr() string
	IsEmpty() bool
	GetResources() (Resource, error)
}
