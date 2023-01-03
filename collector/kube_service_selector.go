package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	SvcMetricVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kube_service_selector",
			Help: "kube_service_selector contain all selector of services",
		},
		[]string{"namespace", "name", "selector"},
	)
)
