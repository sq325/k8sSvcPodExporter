package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/sq325/svcPodKmsExporter/utils"

	"github.com/prometheus/client_golang/prometheus"
)

type KubeNamespace struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			CreationTimestamp time.Time `json:"creationTimestamp"`
			Labels            struct {
				KubernetesIoMetadataName string `json:"kubernetes.io/metadata.name"`
			} `json:"labels"`
			Name            string `json:"name"`
			ResourceVersion string `json:"resourceVersion"`
			UID             string `json:"uid"`
		} `json:"metadata"`
		Spec struct {
			Finalizers []string `json:"finalizers"`
		} `json:"spec"`
		Status struct {
			Phase string `json:"phase"`
		} `json:"status"`
	} `json:"items"`
	Kind     string `json:"kind"`
	Metadata struct {
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
}

var (
	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})
	hdFailures = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kube_service_selector",
			Help: "Kubernetes service selector converted to Prometheus labels",
		},
		[]string{"namespace", "service"},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(hdFailures)
}

func main() {
	cmd := exec.Command("sh", "-c", `kubectl get pods -A -o=jsonpath='{range .items[*]}{.metadata.namespace},{.metadata.name},{.metadata.labels}{"\n"}{end}'`)
	_, ifEmpty := utils.RunCmd(cmd)
	fmt.Println(ifEmpty)
	// cpuTemp.Set(65.3)
	// hdFailures.With(prometheus.Labels{"device": "/dev/sda", "label2": "laebel2Val"}).Inc()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	// http.Handle("/metrics", promhttp.Handler())
	// log.Fatal(http.ListenAndServe(":8181", nil))
}
