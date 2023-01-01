package resource

const (
	GetAllSvcs string = `kubectl get svc -A -o=jsonpath='{range .items[*]}{.metadata.namespace};{.metadata.name};{.metadata.labels};{.spec.selector}{"\n"}{end}'`
	// default,kubernetes,{"component":"apiserver","provider":"kubernetes"},
	// kube-system,kube-dns,{"k8s-app":"kube-dns","kubernetes.io/cluster-service":"true","kubernetes.io/name":"CoreDNS"},{"k8s-app":"kube-dns"}
	// kube-system,metrics-server,{"addonmanager.kubernetes.io/mode":"Reconcile","k8s-app":"metrics-server","kubernetes.io/minikube-addons":"metrics-server","kubernetes.io/minikube-addons-endpoint":"metrics-server","kubernetes.io/name":"Metrics-server"},{"k8s-app":"metrics-server"}
	// monitoring,prometheus,{"app":"prometheus"},{"app":"prometheus"}
)

// Service define a service resource
type Service struct {
	Name, Namespace string
	Selector        map[string]string
}

type Services []Service
