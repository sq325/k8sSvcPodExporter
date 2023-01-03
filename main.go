package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sq325/svcPodKmsExporter/collector"
	"github.com/sq325/svcPodKmsExporter/resource"
	"github.com/sq325/svcPodKmsExporter/utils"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	svcMetricVec  = collector.SvcMetricVec
	podMetricVec  = collector.PodMetricVec
	kubectlPodCmd = resource.KubectlPodCmd
	kubectlSvcCmd = resource.KubectlSvcCmd
)

func init() {
	// Metrics have to be registered to be exposed
	prometheus.MustRegister(svcMetricVec)
	prometheus.MustRegister(podMetricVec)
}

func main() {
	pF := resource.NewPodFactor(kubectlPodCmd)
	sF := resource.NewSvcFactor(kubectlSvcCmd)
	pods, err := pF.GetResources()
	if err != nil {
		fmt.Println(err)
		return
	}
	svcs, err := sF.GetResources()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, p := range pods {
		labelsStr, err := utils.MapToStr(p.Labels())
		if err != nil {
			log.Fatal(err)
			return
		}
		podMetricVec.With(prometheus.Labels{
			"name":      p.Name(),
			"namespace": p.Namespace(),
			"labels":    labelsStr,
		}).Set(1)
	}
	for _, s := range svcs {
		selectorStr, err := utils.MapToStr(s.Selector())
		if err != nil {
			log.Fatal(err)
			return
		}
		svcMetricVec.With(prometheus.Labels{
			"name":      s.Name(),
			"namespace": s.Namespace(),
			"selector":  selectorStr,
		}).Set(1)
	}

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8181", nil))
}
