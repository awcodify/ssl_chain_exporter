package exporter

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"

	"crypto/tls"
	"fmt"
	"net"
	"time"
)

const namespace = "ssl_chain"

var peerCertificatesIndex = []string{"server", "intermediate", "root"}

type SSLOptions struct {
	Options []SSLOption
}

type SSLOption struct {
	Domain string
}

type sslChainCollector struct {
	up               *prometheus.Desc
	expiry           *prometheus.Desc
	scrapeErrorTotal *prometheus.CounterVec

	sslOptions SSLOptions
	logger     log.Logger
}

func newSSLChainCollector(opts *SSLOptions, logger log.Logger) *sslChainCollector {
	return &sslChainCollector{
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Could the server be reached",
			[]string{"domain"},
			nil,
		),
		expiry: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "expiry"),
			"expiration of certification",
			[]string{"domain", "chain", "issuer"},
			nil,
		),
		scrapeErrorTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "exporter_scrape_error_total",
				Help:      "Number of errors while scraping SSL chain.",
			},
			[]string{"domain"},
		),
		sslOptions: *opts,
		logger:     logger,
	}
}

func (collector *sslChainCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.up
	ch <- collector.expiry
	collector.scrapeErrorTotal.Describe(ch)
}

func (collector *sslChainCollector) Collect(ch chan<- prometheus.Metric) {
	for _, opt := range collector.sslOptions.Options {
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 2 * time.Second}, "tcp", opt.Domain+":443", nil)
		if err != nil {
			// make status down
			ch <- prometheus.MustNewConstMetric(collector.up, prometheus.GaugeValue, 0, opt.Domain)
			// increase error counter
			collector.scrapeErrorTotal.WithLabelValues(opt.Domain).Inc()
			collector.scrapeErrorTotal.WithLabelValues(opt.Domain).Collect(ch)

			level.Info(collector.logger).Log(opt.Domain + " doesn't support SSL certificate err: " + err.Error())
			continue
		}

		err = conn.VerifyHostname(opt.Domain)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(collector.up, prometheus.GaugeValue, 0, opt.Domain)

			level.Info(collector.logger).Log("Hostname doesn't match with certificate: " + err.Error())
			continue
		}

		// success connect with TLS, then make status become UP (1)
		ch <- prometheus.MustNewConstMetric(collector.up, prometheus.GaugeValue, 1, opt.Domain)

		collector.collect(conn, opt, ch)
	}
}

func (collector *sslChainCollector) collect(conn *tls.Conn, opt SSLOption, ch chan<- prometheus.Metric) {
	for id, chain := range peerCertificatesIndex {
		expiry := conn.ConnectionState().PeerCertificates[id].NotAfter
		issuer := fmt.Sprintf("%s", conn.ConnectionState().PeerCertificates[id].Issuer)

		ch <- prometheus.MustNewConstMetric(collector.expiry, prometheus.GaugeValue, float64(expiry.Unix()), opt.Domain, chain, issuer)
	}
}

func Register(options *SSLOptions, logger log.Logger) {
	collector := newSSLChainCollector(options, logger)
	prometheus.MustRegister(version.NewCollector("ssl_chain_collector"))
	prometheus.MustRegister(collector)
	prometheus.Unregister(prometheus.NewGoCollector())
}
