package resource

// Resource is the interface a k8s resource has to implement
type Resource interface {
	Name() string
	Namespace() string
	Kind() string
}

type Resources []*Resource

// Factor is the interface a resource factor has to implement
type Factor interface {
	CmdStr() string
	IsEmpty() bool
	GetResources() (Resources, error)
}
