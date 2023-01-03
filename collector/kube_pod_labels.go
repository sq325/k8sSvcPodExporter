package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	PodMetricVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kube_pod_labels",
			Help: "Contain all labels of pods.",
		},
		[]string{"namespace", "name", "labels"},
	)
)
