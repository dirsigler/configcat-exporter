package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "configcat"
)

// Metrics holds all Prometheus metrics for ConfigCat
type Metrics struct {
	ConfigsTotal      *prometheus.GaugeVec
	EnvironmentsTotal *prometheus.GaugeVec
	FeatureFlagsTotal *prometheus.GaugeVec
	LastScrapeTime    prometheus.Gauge
	ProductsTotal     prometheus.Gauge
	ScrapeDuration    prometheus.Histogram
	ScrapeErrors      prometheus.Counter
	ZombieFlagsTotal  *prometheus.GaugeVec
}

var (
	productLabels = []string{"product_id", "product_name"}
	configLabels  = []string{"config_id", "config_name"}
)

// New creates and registers Prometheus metrics
func New() *Metrics {
	m := &Metrics{
		ProductsTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:      "products_total",
			Help:      "Total number of ConfigCat products",
			Namespace: namespace,
		}),
		ConfigsTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "configs_total",
			Help:      "Total number of ConfigCat configs per product",
			Namespace: namespace,
		},
			productLabels,
		),
		EnvironmentsTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "environments_total",
			Help:      "Total number of ConfigCat environments per product",
			Namespace: namespace,
		},
			productLabels,
		),
		FeatureFlagsTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "feature_flags_total",
			Help:      "Total number of feature flags per config",
			Namespace: namespace,
		}, append(productLabels, configLabels...)),
		ScrapeErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name:      "scrape_errors_total",
			Help:      "Total number of ConfigCat API scrape errors",
			Namespace: namespace,
		}),
		LastScrapeTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:      "last_scrape_timestamp",
			Help:      "Unix timestamp of the last successful scrape",
			Namespace: namespace,
		}),
		ScrapeDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:      "scrape_duration_seconds",
			Help:      "Duration of ConfigCat API scrapes in seconds",
			Namespace: namespace,
			Buckets:   prometheus.DefBuckets,
		}),
		ZombieFlagsTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "zombie_flags_total",
			Help:      "Total number of zombie flags per product",
			Namespace: namespace,
		}, append(productLabels, configLabels...)),
	}

	// Register metrics with Prometheus
	prometheus.MustRegister(
		m.ProductsTotal,
		m.ConfigsTotal,
		m.EnvironmentsTotal,
		m.FeatureFlagsTotal,
		m.ScrapeErrors,
		m.LastScrapeTime,
		m.ScrapeDuration,
		m.ZombieFlagsTotal,
	)

	return m
}

// Reset clears all metric values (useful for testing)
func (m *Metrics) Reset() {
	m.ProductsTotal.Set(0)
	m.ConfigsTotal.Reset()
	m.EnvironmentsTotal.Reset()
	m.FeatureFlagsTotal.Reset()
	m.ZombieFlagsTotal.Reset()
	// Note: Counters and histograms cannot be reset in Prometheus
}
