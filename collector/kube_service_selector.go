package collector

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sq325/svcPodKmsExporter/resource"
	"github.com/sq325/svcPodKmsExporter/utils"
)

var (
	SvcMetricName      = "kube_service_selector"
	SvcMetricHelp      = "kube_service_selector contains all selectors services"
	SvcMetricLabelKeys = []string{"namespace", "name", "selector"}
)

var (
	kubectlSvcCmd                 = resource.KubectlSvcCmd
	svcPod        resource.Factor = resource.NewSvcFactor(kubectlSvcCmd)
)

func NewSvcMetric(metricName, help string, labelKeys, labelValues []string, value float64) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, help, labelKeys, nil),
		prometheus.GaugeValue,
		value,
		labelValues...,
	)
}

type SvcCollector struct {
	desc   *prometheus.Desc
	factor resource.Factor
}

func (c *SvcCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *SvcCollector) Collect(ch chan<- prometheus.Metric) {
	svcs, err := c.factor.GetResources()
	if err != nil {
		log.Println(err)
		return
	}
	for _, s := range svcs {
		selectorStr, err := utils.MapToStr(s.Selector())
		if err != nil {
			log.Fatal(err)
			return
		}
		name := s.Name()
		namespace := s.Namespace()
		keys := SvcMetricLabelKeys
		values := []string{namespace, name, selectorStr}
		ch <- NewPodMetric(
			SvcMetricName,
			SvcMetricHelp,
			keys,
			values,
			float64(1),
		)
	}
}

func NewSvcCollector(desc *prometheus.Desc) prometheus.Collector {
	return &SvcCollector{
		desc:   desc,
		factor: svcPod,
	}
}
