package collectors

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/thelastguardian/tplinkexporter/clients"
)

type TrafficCollector struct {
	namespace     string
	client        clients.TPLINKSwitchClient
	pktMetrics    map[string](*prometheus.GaugeVec)
	statusMetrics map[string](*prometheus.GaugeVec)

	// trafficScrapesTotalMetric              prometheus.Gauge
	// trafficScrapeErrorsTotalMetric         prometheus.Gauge
	// lastTrafficScrapeErrorMetric           prometheus.Gauge
	// lastTrafficScrapeTimestampMetric       prometheus.Gauge
	// lastTrafficScrapeDurationSecondsMetric prometheus.Gauge
}

var statusMetricsFields = []string{
	"State",
	"LinkStatus",
}

var pktMetricsFields = []string{
	"TxGoodPkt",
	"TxBadPkt",
	"RxGoodPkt",
	"RxBadPkt",
}

func NewTrafficCollector(namespace string, client clients.TPLINKSwitchClient) *TrafficCollector {
	pktMetrics := make(map[string]*prometheus.GaugeVec)
	for _, name := range pktMetricsFields {
		pktMetrics[name] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "portstats",
				Name:      strings.ToLower(name),
				Help:      fmt.Sprintf("Value of the '%s' traffic metric from the router", name),
			},
			[]string{"portnum", "host"},
		)
	}
	statusMetrics := make(map[string]*prometheus.GaugeVec)
	for _, name := range statusMetricsFields {
		statusMetrics[name] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "portstats",
				Name:      strings.ToLower(name),
				Help:      fmt.Sprintf("Value of the '%s' traffic metric from the router", name),
			},
			[]string{"portnum", "host"},
		)
	}
	return &TrafficCollector{
		namespace:     namespace,
		client:        client,
		pktMetrics:    pktMetrics,
		statusMetrics: statusMetrics,

		// trafficScrapesTotalMetric:              trafficScrapesTotalMetric,
		// trafficScrapeErrorsTotalMetric:         trafficScrapeErrorsTotalMetric,
		// lastTrafficScrapeErrorMetric:           lastTrafficScrapeErrorMetric,
		// lastTrafficScrapeTimestampMetric:       lastTrafficScrapeTimestampMetric,
		// lastTrafficScrapeDurationSecondsMetric: lastTrafficScrapeDurationSecondsMetric,
	}
}

func (c *TrafficCollector) Collect(ch chan<- prometheus.Metric) {
	// var begun = time.Now()

	stats, err := c.client.GetPortStats()
	if err != nil {
		log.Errorf("Error while collecting traffic statistics: %v", err)
		// c.trafficScrapeErrorsTotalMetric.Inc()
	} else {
		for portnum := 0; portnum < len(stats); portnum++ {
			for name, value := range stats[portnum].PktCount {
				// log.Infof("portnum '%d', metricname '%s', metricvalue '%d'", portnum, name, value)
				c.pktMetrics[name].With(prometheus.Labels{"portnum": strconv.FormatInt(int64(portnum+1), 10), "host": c.client.GetHost()}).Set(float64(value))
			}
			// log.Infof("portnum '%d', state '%d', linkstatus '%d'", portnum, stats[portnum].State, stats[portnum].LinkStatus)
			c.statusMetrics["State"].With(prometheus.Labels{"portnum": strconv.FormatInt(int64(portnum+1), 10), "host": c.client.GetHost()}).Set(float64(stats[portnum].State))
			c.statusMetrics["LinkStatus"].With(prometheus.Labels{"portnum": strconv.FormatInt(int64(portnum+1), 10), "host": c.client.GetHost()}).Set(float64(stats[portnum].LinkStatus))
		}
		for name := range c.pktMetrics {
			c.pktMetrics[name].Collect(ch)
		}
		for name := range c.statusMetrics {
			c.statusMetrics[name].Collect(ch)
		}
	}

	// c.trafficScrapeErrorsTotalMetric.Collect(ch)

	// c.trafficScrapesTotalMetric.Inc()
	// c.trafficScrapesTotalMetric.Collect(ch)

	// c.lastTrafficScrapeErrorMetric.Set(errorMetric)
	// c.lastTrafficScrapeErrorMetric.Collect(ch)

	// c.lastTrafficScrapeTimestampMetric.Set(float64(time.Now().Unix()))
	// c.lastTrafficScrapeTimestampMetric.Collect(ch)

	// c.lastTrafficScrapeDurationSecondsMetric.Set(time.Since(begun).Seconds())
	// c.lastTrafficScrapeDurationSecondsMetric.Collect(ch)
}

func (c *TrafficCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, name := range pktMetricsFields {
		c.pktMetrics[name].Describe(ch)
	}
	for _, name := range statusMetricsFields {
		c.statusMetrics[name].Describe(ch)
	}

	// c.trafficScrapesTotalMetric.Describe(ch)
	// c.trafficScrapeErrorsTotalMetric.Describe(ch)
	// c.lastTrafficScrapeErrorMetric.Describe(ch)
	// c.lastTrafficScrapeTimestampMetric.Describe(ch)
	// c.lastTrafficScrapeDurationSecondsMetric.Describe(ch)
}
