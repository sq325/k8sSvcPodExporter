// PodFactor is only a single instance

package resource

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sq325/svcPodKmsExporter/utils"
)

const (
	kubectlPodCmd string = `kubectl get pods -A -o=jsonpath='{range .items[*]}{.metadata.namespace};{.metadata.name};{.metadata.labels}{"\n"}{end}'`
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
)

// Pod define a pod resource
type Pod struct {
	name      string
	namespace string
	kind      string
	labels    map[string]string
}

func NewPod(name, namespace string, labels map[string]string) *Pod {
	return &Pod{
		name:      name,
		namespace: namespace,
		kind:      "pod",
		labels:    labels,
	}
}

type Pods []*Pod

// PodFactor implements Factor interface
// PodFactor parse output and produce Pods
type PodFactor struct {
	cmdstr string // kubectl command
}

func NewPodFactor(cmd string) *PodFactor {
	return &PodFactor{
		cmdstr: cmd,
	}
}

func (p *PodFactor) CmdStr() string {
	return p.cmdstr
}

func (p *PodFactor) runCmd() (*bufio.Scanner, bool) {
	cmd := exec.Command("sh", "-c", p.cmdstr)
	scanner, isempty := utils.RunCmd(cmd)
	return scanner, isempty
}

func (p *PodFactor) IsEmpty() bool {
	_, b := p.runCmd()
	return b
}

func (p *PodFactor) parseLineS(lineS []string) (name, namespace string, m map[string]string, err error) {
	if len(lineS) != 3 {
		return "", "", nil, fmt.Errorf("pod lineS colume num != 3\n%s", p.CmdStr())
	}
	m, err = utils.JsonStrToMap(lineS[2])
	if err != nil {
		return "", "", nil, err
	}
	return lineS[1], lineS[0], m, nil
}

func (p *PodFactor) GetResources() (Pods, error) {
	scanner, isempty := p.runCmd()
	if isempty {
		return nil, errors.New("no resources found")
	}
	var pods Pods
	for scanner.Scan() {
		line := scanner.Text()
		lineS := strings.Split(line, ";")
		name, namespace, m, err := p.parseLineS(lineS)
		if err != nil {
			return nil, err
		}
		pod := NewPod(name, namespace, m)
		pods = append(pods, pod)
	}
	if len(pods) == 0 {
		return nil, fmt.Errorf("pods is empty, cmd: %s", p.CmdStr())
	}
	return pods, nil
}
