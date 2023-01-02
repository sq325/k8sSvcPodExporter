package resource

import "testing"

// kubectlPodCmd want
// default;busybox;{"run":"busybox"}
// kube-system;coredns-6d4b75cb6d-bz268;{"k8s-app":"kube-dns","pod-template-hash":"6d4b75cb6d"}
// kube-system;etcd-minikube;{"component":"etcd","tier":"control-plane"}
// kube-system;kindnet-lqzpd;{"app":"kindnet","controller-revision-hash":"78f985b4","k8s-app":"kindnet","pod-template-generation":"1","tier":"node"}
// kube-system;kube-apiserver-minikube;{"component":"kube-apiserver","tier":"control-plane"}
// kube-system;kube-controller-manager-minikube;{"component":"kube-controller-manager","tier":"control-plane"}
// kube-system;kube-proxy-49wpd;{"controller-revision-hash":"58bf5dfbd7","k8s-app":"kube-proxy","pod-template-generation":"1"}
// kube-system;kube-scheduler-minikube;{"component":"kube-scheduler","tier":"control-plane"}
// kube-system;metrics-server-8595bd7d4c-66v84;{"k8s-app":"metrics-server","pod-template-hash":"8595bd7d4c"}
// kube-system;storage-provisioner;{"addonmanager.kubernetes.io/mode":"Reconcile","integration-test":"storage-provisioner"}
// monitoring;prometheus-6856764f47-mw6gc;{"app":"prometheus","pod-template-hash":"6856764f47"}

func TestRuncmd(t *testing.T) {
	f := NewPodFactor(kubectlPodCmd)
	scan, _ := f.runcmd()
	for scan.Scan() {
		t.Log(scan.Text())
	}
}

func TestPodGetResources(t *testing.T) {
	f := NewPodFactor(kubectlPodCmd)
	pods, err := f.GetResources()
	if err != nil {
		t.Log(err)
	}
	for _, p := range pods {
		t.Log(p)
	}
}
