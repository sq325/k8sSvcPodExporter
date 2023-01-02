package resource

import "testing"

// default;kubernetes;
// kube-system;kube-dns;{"k8s-app":"kube-dns"}
// kube-system;metrics-server;{"k8s-app":"metrics-server"}
// monitoring;prometheus;{"app":"prometheus"}

func TestSvcGetResources(t *testing.T) {
	f := NewSvcFactor(kubectlSvcCmd)
	svcs, err := f.GetResources()
	if err != nil {
		t.Log(err)
	}
	for _, s := range svcs {
		t.Log(s)
	}
}
