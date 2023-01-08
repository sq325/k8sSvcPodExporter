package resource

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sq325/svcPodKmsExporter/utils"
)

const (
	KubectlSvcCmd string = `kubectl get svc -A -o=jsonpath='{range .items[*]}{.metadata.namespace};{.metadata.name};{.spec.selector}{"\n"}{end}'`
)

// Service implement Resource interface
// Service define a service resource
type Svc struct {
	name, namespace string
	kind            string // title style
	selector        map[string]string
}
type Svcs []*Svc

func NewSvc(name, namespaces string, selector map[string]string) *Svc {
	return &Svc{
		name:      name,
		namespace: namespaces,
		kind:      "Service",
		selector:  selector,
	}
}

func (s *Svc) Name() string {
	return s.name
}
func (s *Svc) Namespace() string {
	return s.namespace
}

func (s *Svc) Kind() string {
	return s.kind
}

func (s *Svc) Selector() map[string]string {
	return s.selector
}

func (s *Svc) Labels() prometheus.Labels {
	return nil
}

// SvcFactor implement Factor interface
type SvcFactor struct {
	cmdstr string // kubectl command
}

func NewSvcFactor(cmdstr string) Factor {
	return &SvcFactor{
		cmdstr: cmdstr,
	}
}

func (s *SvcFactor) CmdStr() string {
	return s.cmdstr
}

func (s *SvcFactor) runcmd() (*bufio.Scanner, bool) {
	cmd := exec.Command("sh", "-c", s.cmdstr)
	scanner, isempty := utils.RunCmd(cmd)
	return scanner, isempty
}

func (s *SvcFactor) IsEmpty() bool {
	_, b := s.runcmd()
	return b
}

func (s *SvcFactor) parseLineS(lineS []string) (name, namespace string, m map[string]string, err error) {
	if len(lineS) != 3 {
		return "", "", nil, fmt.Errorf("svc lineS colume num != 3\n%s", s.CmdStr())
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

func (s *SvcFactor) GetResources() (Resources, error) {
	scanner, isempty := s.runcmd()
	if isempty {
		return nil, errors.New("no resources found")
	}
	var svcs Resources
	for scanner.Scan() {
		line := scanner.Text()
		lineS := strings.Split(line, ";")
		name, namespace, m, err := s.parseLineS(lineS)
		if err != nil {
			return nil, err
		}
		svc := NewSvc(name, namespace, m)
		svcs = append(svcs, svc)
	}
	if len(svcs) == 0 {
		return nil, fmt.Errorf("svcs is empty, cmd: %s", s.CmdStr())
	}
	return svcs, nil
}
