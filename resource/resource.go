package resource

import "github.com/prometheus/client_golang/prometheus"

// Resource is the interface a k8s resource has to implement
type Resource interface {
	Name() string
	Namespace() string
	Kind() string
	Labels() prometheus.Labels
	Selector() map[string]string
}

type Resources []Resource

// Factor is the interface a resource factor has to implement
type Factor interface {
	CmdStr() string
	IsEmpty() bool
	GetResources() (Resources, error)
}
