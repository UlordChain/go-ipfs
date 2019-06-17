package corehttp

import (
	"github.com/ipfs/go-ipfs/core"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"net"
	"net/http"
	"os"
	"path"

	"gx/ipfs/QmTQuFQWHAWy4wMH6ZyPfGiawA5u9T8rs79FENoV8yXaoS/client_golang/prometheus"
	"gx/ipfs/QmTQuFQWHAWy4wMH6ZyPfGiawA5u9T8rs79FENoV8yXaoS/client_golang/prometheus/promhttp"
)

// This adds the scraping endpoint which Prometheus uses to fetch metrics.
func MetricsScrapingOption(path string) ServeOption {
	return func(n *core.IpfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		mux.Handle(path, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
		return mux, nil
	}
}

// This adds collection of net/http-related metrics
func MetricsCollectionOption(handlerName string) ServeOption {
	return func(_ *core.IpfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		// Adapted from github.com/prometheus/client_golang/prometheus/http.go
		// Work around https://github.com/prometheus/client_golang/pull/311
		opts := prometheus.SummaryOpts{
			Subsystem:   "http",
			ConstLabels: prometheus.Labels{"handler": handlerName},
			Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}

		reqCnt := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace:   opts.Namespace,
				Subsystem:   opts.Subsystem,
				Name:        "requests_total",
				Help:        "Total number of HTTP requests made.",
				ConstLabels: opts.ConstLabels,
			},
			[]string{"method", "code"},
		)
		if err := prometheus.Register(reqCnt); err != nil {
			if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
				reqCnt = are.ExistingCollector.(*prometheus.CounterVec)
			} else {
				return nil, err
			}
		}

		opts.Name = "request_duration_seconds"
		opts.Help = "The HTTP request latencies in seconds."
		reqDur := prometheus.NewSummaryVec(opts, nil)
		if err := prometheus.Register(reqDur); err != nil {
			if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
				reqDur = are.ExistingCollector.(*prometheus.SummaryVec)
			} else {
				return nil, err
			}
		}

		opts.Name = "request_size_bytes"
		opts.Help = "The HTTP request sizes in bytes."
		reqSz := prometheus.NewSummaryVec(opts, nil)
		if err := prometheus.Register(reqSz); err != nil {
			if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
				reqSz = are.ExistingCollector.(*prometheus.SummaryVec)
			} else {
				return nil, err
			}
		}

		opts.Name = "response_size_bytes"
		opts.Help = "The HTTP response sizes in bytes."
		resSz := prometheus.NewSummaryVec(opts, nil)
		if err := prometheus.Register(resSz); err != nil {
			if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
				resSz = are.ExistingCollector.(*prometheus.SummaryVec)
			} else {
				return nil, err
			}
		}

		// Construct the mux
		childMux := http.NewServeMux()
		var promMux http.Handler = childMux
		promMux = promhttp.InstrumentHandlerResponseSize(resSz, promMux)
		promMux = promhttp.InstrumentHandlerRequestSize(reqSz, promMux)
		promMux = promhttp.InstrumentHandlerDuration(reqDur, promMux)
		promMux = promhttp.InstrumentHandlerCounter(reqCnt, promMux)
		mux.Handle("/", promMux)

		return childMux, nil
	}
}

var (
	peersTotalMetric = prometheus.NewDesc(
		prometheus.BuildFQName("ipfs", "p2p", "peers_total"),
		"Number of connected peers", []string{"transport"}, nil)
	cpuTotalMetric = prometheus.NewDesc(
		prometheus.BuildFQName("ipfs", "machine", "cpu_total"),
		"Number of cpu cores", []string{}, nil)
	memoryTotalMetric = prometheus.NewDesc(
		prometheus.BuildFQName("ipfs", "machine", "memory_total"),
		"Size of memory total", []string{}, nil)
	memoryFreeMetric = prometheus.NewDesc(
		prometheus.BuildFQName("ipfs", "machine", "memory_free"),
		"Size of memory free", []string{}, nil)

	diskTotalMetric = prometheus.NewDesc(
		prometheus.BuildFQName("ipfs", "machine", "disk_total"),
		"Size of disk total", []string{"path","fstype"}, nil)
	diskFreeMetric = prometheus.NewDesc(
		prometheus.BuildFQName("ipfs", "machine", "disk_free"),
		"Size of disk free", []string{"path","fstype"}, nil)
)

type IpfsNodeCollector struct {
	Node *core.IpfsNode
}

func (_ IpfsNodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- peersTotalMetric
	ch <- cpuTotalMetric
}

func ipfsPath() string {
	p := os.Getenv("IPFS_PATH")
	if p == "" {
		p = path.Join(os.Getenv("HOME"), ".ipfs")
	}
	return p
}

func (c IpfsNodeCollector) Collect(ch chan<- prometheus.Metric) {
	for tr, val := range c.PeersTotalValues() {
		ch <- prometheus.MustNewConstMetric(
			peersTotalMetric,
			prometheus.GaugeValue,
			val,
			tr,
		)
	}

	infos, _ := cpu.Info()
	ch <- prometheus.MustNewConstMetric(
		cpuTotalMetric,
		prometheus.GaugeValue,
		float64(len(infos)),
	)

	vms, _ := mem.VirtualMemory()
	ch <- prometheus.MustNewConstMetric(
		memoryTotalMetric,
		prometheus.GaugeValue,
		float64(vms.Total),
		)

	ch <- prometheus.MustNewConstMetric(
		memoryFreeMetric,
		prometheus.GaugeValue,
		float64(vms.Free),
	)

	us, _ := disk.Usage(ipfsPath())
	ch <- prometheus.MustNewConstMetric(
		diskTotalMetric,
		prometheus.GaugeValue,
		float64(us.Total),
		us.Path,
		us.Fstype,
	)

	ch <- prometheus.MustNewConstMetric(
		diskFreeMetric,
		prometheus.GaugeValue,
		float64(us.Free),
		us.Path,
		us.Fstype,
	)

}

func (c IpfsNodeCollector) PeersTotalValues() map[string]float64 {
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
