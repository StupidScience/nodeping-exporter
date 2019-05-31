package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

func healthz(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(response, "ok")
}

func main() {
	np, err := NewNodePing("https://api.nodeping.com/api/1", os.Getenv("NODEPING_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	c, err := NewCollector(np)
	if err != nil {
		log.Fatalf("Can't create collector: %v", err)
	}
	prometheus.MustRegister(c)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>NodePing Exporter</title></head>
			<body>
			<h1>NodePing Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})
	log.Infoln("Starting nodeping-exporter")
	log.Fatal(http.ListenAndServe(":9503", nil))
}
