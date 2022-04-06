package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	// "github.com/prometheus/client_golang/prometheus"
	"github.com/awcodify/ssl-chain-exporter/exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type sslChainOpts struct {
	domainName []string
}

var (
	listenAddress = flag.String("web.listen-address", ":9102", "Address to listen on for web interface.")
	metricPath    = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
	domains       = flag.String("domains", "", "Which domain will be collected. Comma separated.")
)

func main() {
	flag.Parse()
	domainList := strings.Split(*domains, ",")

	sslOpts := exporter.SSLOptions{}
	for _, domain := range domainList {
		opt := exporter.SSLOption{
			Domain: domain,
		}

		sslOpts.Options = append(sslOpts.Options, opt)
	}

	exporter.Register(&sslOpts)

	log.Fatal(serverMetrics(*listenAddress, *metricPath))
}

func serverMetrics(listenAddress, metricsPath string) error {
	http.Handle(metricsPath, promhttp.Handler())

	return http.ListenAndServe(listenAddress, nil)
}
