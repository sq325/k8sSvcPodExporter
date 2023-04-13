// PodFactor is only a single instance

package resource

import (
	"context"
	"log"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Service implement Resource interface
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

func (p *Pod) Name() string {
	return p.name
}

func (p *Pod) Namespace() string {
	return p.namespace
}

func (p *Pod) Kind() string {
	return p.kind
}

func (p *Pod) Labels() map[string]string {
	return p.labels
}

func (p *Pod) Selector() map[string]string {
	return nil
}

type Pods []*Pod

// PodFactor implements Factor interface
// PodFactor parse output and produce Pods
type PodFactor struct {
	ClientSet *kubernetes.Clientset
}

func NewPodFactor(clientSet *kubernetes.Clientset) Factor {
	return &PodFactor{ClientSet: clientSet}
}

func (p *PodFactor) GetResources() (Resources, error) {
	podList, err := p.ClientSet.CoreV1().Pods("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	var pods Resources
	for _, npod := range podList.Items {
		pod := NewPod(npod.Name, npod.Namespace, npod.Labels)
		pods = append(pods, pod)
	}

	if len(pods) == 0 {
		log.Println("No pods found in cluster")
		return nil, nil
	}
	return pods, nil
}
