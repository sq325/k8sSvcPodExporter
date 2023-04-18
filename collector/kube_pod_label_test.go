package collector

import (
	"os"
	"testing"

	"github.com/sq325/svcPodKmsExporter/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func setup() {

}

func BenchmarkPodCollector(b *testing.B) {
	kubeconfig := os.Getenv("HOME") + "/.kube/config"
	// Create a Kubernetes clientset
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, _ := kubernetes.NewForConfig(config)
	factor := resource.NewPodFactor(clientset)
	collector := NewPodCollector(factor)
	for i := 0; i < b.N; i++ {
		collector.Collect(nil)
	}
}
