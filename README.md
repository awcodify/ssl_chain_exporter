
# SSL Chain Exporter for Prometheus

[![GitHub release](https://img.shields.io/github/release/awcodify/ssl_chain_exporter.svg)][release]
![GitHub Downloads](https://img.shields.io/github/downloads/awcodify/ssl_chain_exporter/total.svg)

Help on flags:

<pre>
  -domains string
        Which domain will be collected. Comma separated.
  -web.listen-address string
        Address to listen on for web interface. (default ":9102")
  -web.metrics-path string
        Path under which to expose metrics. (default "/metrics")
</pre>

Check your domains by runnning
```
./ssl_chain_exporter --domains=github.com/gist.github.com
```
## Collectors

SSL Chain metrics:

| Metric                                | Description                                                                       | Type    | Label                 |
|---------------------------------------|-----------------------------------------------------------------------------------|---------|-----------------------|
| ssl_chain_up                          | Is the provided domain can be reached or not.                                     | Gauge   | domain                |
| ssl_chain_expiry                      | The date after which a peer certificate expires. Expressed as a Unix Epoch Time.  | Gauge   | domain, chain, issuer |
| ssl_chain_exporter_scrape_error_total | Number of errors while scraping SSL chain.                                        | Counter | domain                |

Sample:
```
# HELP ssl_chain_collector_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which ssl_chain_collector was built.
# TYPE ssl_chain_collector_build_info gauge
ssl_chain_collector_build_info{branch="",goversion="go1.17.3",revision="",version=""} 1
# HELP ssl_chain_expiry expiration of certification
# TYPE ssl_chain_expiry gauge
ssl_chain_expiry{chain="intermediate",domain="google.co.id",issuer="CN=GTS Root R1,O=Google Trust Services LLC,C=US"} 1.822262442e+09
ssl_chain_expiry{chain="root",domain="google.co.id",issuer="CN=GlobalSign Root CA,OU=Root CA,O=GlobalSign nv-sa,C=BE"} 1.832630442e+09
ssl_chain_expiry{chain="server",domain="google.co.id",issuer="CN=GTS CA 1C3,O=Google Trust Services LLC,C=US"} 1.654776441e+09
# HELP ssl_chain_exporter_scrape_error_total Number of errors while scraping SSL chain.
# TYPE ssl_chain_exporter_scrape_error_total counter
ssl_chain_exporter_scrape_error_total{domain="asdasd.as"} 3
ssl_chain_exporter_scrape_error_total{domain="hasldklaskdlaslda.asa"} 3
# HELP ssl_chain_up Could the server be reached
# TYPE ssl_chain_up gauge
ssl_chain_up{domain="asdasd.as"} 0
ssl_chain_up{domain="google.co.id"} 1
ssl_chain_up{domain="hasldklaskdlaslda.asa"} 0
```

[release]: https://github.com/awcodify/ssl_chain_exporter/releases/latest
