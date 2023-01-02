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
)

// Pod define a pod resource
type Pod struct {
	name      string
	namespace string
	kind      string // title style
	labels    map[string]string
}

func NewPod(name, namespace string, labels map[string]string) *Pod {
	return &Pod{
		name:      name,
		namespace: namespace,
		kind:      "Pod",
		labels:    labels,
	}
}

type Pods []*Pod

// PodFactor implements Factor interface
// PodFactor parse output and produce Pods
type PodFactor struct {
	cmdstr string // kubectl command
}

func NewPodFactor(cmdstr string) *PodFactor {
	return &PodFactor{
		cmdstr: cmdstr,
	}
}

func (p *PodFactor) CmdStr() string {
	return p.cmdstr
}

func (p *PodFactor) runcmd() (*bufio.Scanner, bool) {
	cmd := exec.Command("sh", "-c", p.cmdstr)
	scanner, isempty := utils.RunCmd(cmd)
	return scanner, isempty
}

func (p *PodFactor) IsEmpty() bool {
	_, b := p.runcmd()
	return b
}

func (p *PodFactor) parseLineS(lineS []string) (name, namespace string, m map[string]string, err error) {
	if len(lineS) != 3 {
		return "", "", nil, fmt.Errorf("pod lineS colume num != 3\n%s", p.CmdStr())
	}
	jsonstr := lineS[2]
	if jsonstr != "" {
		m, err = utils.JsonStrToMap(jsonstr)
		if err != nil {
			return "", "", nil, err
		}
	}
	return lineS[1], lineS[0], m, nil
}

func (p *PodFactor) GetResources() (Pods, error) {
	scanner, isempty := p.runcmd()
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
