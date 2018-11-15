package corehttp

import (
	"net"
	"net/http"

	core "github.com/udfs/go-udfs/core"

	prometheus "gx/udfs/QmYYv3QFnfQbiwmi1tpkgKF8o4xFnZoBrvpupTiGJwL9nH/client_golang/prometheus"
)

// This adds the scraping endpoint which Prometheus uses to fetch metrics.
func MetricsScrapingOption(path string) ServeOption {
	return func(n *core.UdfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		mux.Handle(path, prometheus.UninstrumentedHandler())
		return mux, nil
	}
}

// This adds collection of net/http-related metrics
func MetricsCollectionOption(handlerName string) ServeOption {
	return func(_ *core.UdfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		childMux := http.NewServeMux()
		mux.HandleFunc("/", prometheus.InstrumentHandler(handlerName, childMux))
		return childMux, nil
	}
}

var (
	peersTotalMetric = prometheus.NewDesc(
		prometheus.BuildFQName("udfs", "p2p", "peers_total"),
		"Number of connected peers", []string{"transport"}, nil)
)

type UdfsNodeCollector struct {
	Node *core.UdfsNode
}

func (_ UdfsNodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- peersTotalMetric
}

func (c UdfsNodeCollector) Collect(ch chan<- prometheus.Metric) {
	for tr, val := range c.PeersTotalValues() {
		ch <- prometheus.MustNewConstMetric(
			peersTotalMetric,
			prometheus.GaugeValue,
			val,
			tr,
		)
	}
}

func (c UdfsNodeCollector) PeersTotalValues() map[string]float64 {
	vals := make(map[string]float64)
	if c.Node.PeerHost == nil {
		return vals
	}
	for _, conn := range c.Node.PeerHost.Network().Conns() {
		tr := ""
		for _, proto := range conn.RemoteMultiaddr().Protocols() {
			tr = tr + "/" + proto.Name
		}
		vals[tr] = vals[tr] + 1
	}
	return vals
}
