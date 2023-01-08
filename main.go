// 问题：
// 1. 每次scrape 都会更新metrics

package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/sq325/svcPodKmsExporter/collector"
)

var (
	port *string = pflag.StringP("port", "d", "8181", "bind port, default port: 8181")
)

var ()

func init() {
	// Metrics have to be registered to be exposed
	prometheus.Unregister(promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}))
	prometheus.Unregister(promcollectors.NewGoCollector())
	prometheus.Register(collector.NewPodCollector(prometheus.NewDesc(collector.PodMetricName, collector.PodMetricHelp, collector.PodMetricLabelKeys, nil)))
	prometheus.Register(collector.NewSvcCollector(prometheus.NewDesc(collector.SvcMetricName, collector.SvcMetricHelp, collector.SvcMetricLabelKeys, nil)))
}

func main() {
	pflag.Parse()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
				<head><title>svcPodKmsExporter</title></head>
				<body>
				<h1>Relationship between services and pods Exporter</h1>
				<p>please click <a href="` + "metrics" + `">Metrics</a></p>
				</body>
				</html>`))
	})
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Listening port:", *port)
	log.Println("URL: http://<ip>:" + *port + "/metrics")
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
