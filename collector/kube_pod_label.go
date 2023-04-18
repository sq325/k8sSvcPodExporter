package collector

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sq325/svcPodKmsExporter/resource"
	"github.com/sq325/svcPodKmsExporter/utils"
)

// kube_pod_label metric meta info
var (
	PodMetricName      = "kube_pod_label"
	PodMetricHelp      = "kube_pod_label contains all pods labels"
	PodMetricLabelKeys = []string{"namespace", "pod", "labels"}
)

var (
	podCounterVec *prometheus.CounterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: PodMetricName,
		Help: PodMetricHelp,
	}, PodMetricLabelKeys)
)

type PodCollector struct {
	cv     *prometheus.CounterVec
	factor resource.Factor
}

func (c *PodCollector) Describe(ch chan<- *prometheus.Desc) {
	c.cv.MetricVec.Describe(ch)
}

func (c *PodCollector) Collect(ch chan<- prometheus.Metric) {
	pods, err := c.factor.GetResources()
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, p := range pods {
		labelsStr, err := utils.MapToStr(p.Labels())
		if err != nil {
			log.Fatal(err)
			return
		}
		podName := p.Name()
		namespace := p.Namespace()
		values := []string{namespace, podName, labelsStr}
		c.cv.WithLabelValues(values...).Inc()
	}
	c.cv.Collect(ch)
}

func NewPodCollector(factor resource.Factor) prometheus.Collector {
	return &PodCollector{
		cv:     podCounterVec,
		factor: factor,
	}
}
