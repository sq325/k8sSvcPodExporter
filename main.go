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
	"github.com/sq325/svcPodKmsExporter/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	_ "net/http/pprof"
)

// flags
var (
	port       *string = pflag.StringP("port", "p", "0", "listening port")
	kubeconfig *string = pflag.String("kubeconfig", "", "Path to the kubeconfig file to use for CLI requests, default is one of [ $KUBECONFIG, $HOME/.kube/config ]")
	version    *bool   = pflag.BoolP("version", "v", false, "Version info")
)

func main() {
	pflag.Parse()
	if *version {
		fmt.Println("svcPod_exporter v2.1, service and pod label")
		fmt.Println("Update: 2023-4-17")
		os.Exit(0)
	}

	// get kubeconfig
	if *kubeconfig == "" {
		if *kubeconfig = os.Getenv("KUBECONFIG"); *kubeconfig == "" {
			*kubeconfig = os.Getenv("HOME") + "/.kube/config"
		}
	}
	if _, err := os.Stat(*kubeconfig); os.IsNotExist(err) {
		log.Fatal("kubeconfig file not found")
	}

	// Create a Kubernetes clientset
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// register collector
	podfactor := resource.NewPodFactor(clientset)
	svcfactor := resource.NewSvcFactor(clientset)
	prometheus.Unregister(promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}))
	prometheus.Unregister(promcollectors.NewGoCollector())
	prometheus.Register(collector.NewPodCollector(podfactor))
	prometheus.Register(collector.NewSvcCollector(svcfactor))

	// start http server
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
