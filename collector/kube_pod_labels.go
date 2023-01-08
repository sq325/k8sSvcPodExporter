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
	PodMetricLabelKeys = []string{"namespace", "name", "labels"}
)

// PodFactor
var (
	kubectlPodCmd                 = resource.KubectlPodCmd
	podFactor     resource.Factor = resource.NewPodFactor(kubectlPodCmd)
)

func NewPodMetric(metricName, help string, labelKeys, labelValues []string, value float64) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, help, labelKeys, nil),
		prometheus.GaugeValue,
		value,
		labelValues...,
	)
}

type PodCollector struct {
	desc   *prometheus.Desc
	factor resource.Factor
}

func (c *PodCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
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
		name := p.Name()
		namespace := p.Namespace()
		keys := PodMetricLabelKeys
		values := []string{namespace, name, labelsStr}
		ch <- NewPodMetric(
			PodMetricName,
			PodMetricHelp,
			keys,
			values,
			float64(1),
		)
	}
}

func NewPodCollector(desc *prometheus.Desc) prometheus.Collector {
	return &PodCollector{
		desc:   desc,
		factor: podFactor,
	}
}
