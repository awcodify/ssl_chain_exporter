package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	// "github.com/prometheus/client_golang/prometheus"
	"github.com/awcodify/ssl_chain_exporter/exporter"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
)

type sslChainOpts struct {
	domainName []string
}

var (
	listenAddress  = flag.String("web.listen-address", ":9102", "Address to listen on for web interface.")
	metricPath     = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
	domains        = flag.String("domains", "", "Which domain will be collected. Comma separated.")
	domainsFromCSV = flag.String("domains-from-csv", "", "Which domain will be collected from csv file. Header required.")
)

func readDomainListFromFile(filePath string) []string {
	// parse domains param from file to array
	// open file
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	var domainList []string
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	// skip the header row
	for _, r := range records[1:] {
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		domainList = append(domainList, r[0])
	}
	return domainList
}
func main() {
	flag.Parse()
	var domainList []string
	if *domainsFromCSV != "" {
		fmt.Println("Reading domains from file")
		domainList = readDomainListFromFile(*domainsFromCSV)
	} else {
		fmt.Println("Reading domains from command line")
		// parse domains param to array
		domainList = strings.Split(*domains, ",")
	}
	sslOpts := exporter.SSLOptions{}
	for _, domain := range domainList {
		opt := exporter.SSLOption{
			Domain: domain,
		}

		sslOpts.Options = append(sslOpts.Options, opt)
	}

	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)
	exporter.Register(&sslOpts, logger)

	level.Info(logger).Log("msg", "Starting ssl_chain_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build", version.BuildContext())
	level.Info(logger).Log("msg", "Starting Server: ", "listen_address", *listenAddress)
	level.Info(logger).Log("msg", "Collect from: ", "scrape_uri", *metricPath)

	log.Fatal(serverMetrics(*listenAddress, *metricPath))
}

func serverMetrics(listenAddress, metricsPath string) error {
	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			 <head><title>Apache Exporter</title></head>
			 <body>
			 <h1>Apache Exporter</h1>
			 <p><a href='` + *&metricsPath + `'>Metrics</a></p>
			 </body>
			 </html>`))
	})

	return http.ListenAndServe(listenAddress, nil)
}
