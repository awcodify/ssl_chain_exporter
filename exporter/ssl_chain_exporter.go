package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"

	"crypto/tls"
	"fmt"
)

const prefix = "ssl_chain"

var peerCertificatesIndex = []string{"server", "intermediate", "root"}

type SSLOptions struct {
	Options []SSLOption
}

type SSLOption struct {
	Domain string
}

type sslChainCollector struct {
	expiry *prometheus.Desc

	sslOptions SSLOptions
}

func newSSLChainCollector(opts *SSLOptions) *sslChainCollector {
	return &sslChainCollector{
		expiry: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "expiry"),
			"expiration of certification",
			[]string{"domain", "chain", "issuer"},
			nil,
		),

		sslOptions: *opts,
	}
}

func (collector *sslChainCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.expiry
}

func (collector *sslChainCollector) Collect(ch chan<- prometheus.Metric) {

	//for each descriptor or call other functions that do so.
	for _, opt := range collector.sslOptions.Options {

		conn, err := tls.Dial("tcp", opt.Domain+":443", nil)
		if err != nil {
			panic(opt.Domain + " doesn't support SSL certificate err: " + err.Error())
		}

		err = conn.VerifyHostname(opt.Domain)
		if err != nil {
			panic("Hostname doesn't match with certificate: " + err.Error())
		}

		for id, chain := range peerCertificatesIndex {
			expiry := conn.ConnectionState().PeerCertificates[id].NotAfter
			issuer := fmt.Sprintf("%s", conn.ConnectionState().PeerCertificates[id].Issuer)

			ch <- prometheus.MustNewConstMetric(collector.expiry, prometheus.GaugeValue, float64(expiry.Unix()), opt.Domain, chain, issuer)
		}

	}
}

func Register(options *SSLOptions) {
	collector := newSSLChainCollector(options)
	prometheus.MustRegister(version.NewCollector("volume_exporter"))
	prometheus.MustRegister(collector)
	prometheus.Unregister(prometheus.NewGoCollector())
}
