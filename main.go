package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/sq325/svcPodKmsExporter/collector"
)

// flags
var (
	port    *string = pflag.StringP("port", "p", "0", "bind port, default port: 8181")
	version *bool   = pflag.BoolP("version", "v", false, "Version info")
)

func init() {
	// Metrics have to be registered to be exposed
	prometheus.Unregister(promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}))
	prometheus.Unregister(promcollectors.NewGoCollector())
	prometheus.Register(collector.NewPodCollector(prometheus.NewDesc(collector.PodMetricName, collector.PodMetricHelp, collector.PodMetricLabelKeys, nil)))
	prometheus.Register(collector.NewSvcCollector(prometheus.NewDesc(collector.SvcMetricName, collector.SvcMetricHelp, collector.SvcMetricLabelKeys, nil)))
}

func main() {
	pflag.Parse()
	if *version {
		fmt.Println("svcPod_exporter v1.0")
		fmt.Println("Update: 2023-1-8")
		fmt.Println("Autor: Quan Sun")
		os.Exit(0)
	}

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
	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Address:", listener.Addr().String())
	_port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	log.Println("Listening port:", _port)
	log.Println("URL: http://<ip>:" + _port + "/metrics")
	log.Fatal(http.Serve(listener, nil))
}
