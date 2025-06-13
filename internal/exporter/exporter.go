package exporter

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/dirsigler/configcat-exporter/internal/client"
	"github.com/dirsigler/configcat-exporter/internal/config"
	"github.com/dirsigler/configcat-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Exporter manages the ConfigCat metrics collection
type Exporter struct {
	client  *client.ConfigCatClient
	metrics *metrics.Metrics
	logger  *slog.Logger
	config  *config.Config
	mu      sync.RWMutex
}

// New creates a new ConfigCat exporter
func New(cfg *config.Config, logger *slog.Logger) *Exporter {
	configCatClient := client.New(cfg.ConfigCatAPIKey, cfg.ConfigCatAPIURL, logger)
	prometheusMetrics := metrics.New()

	return &Exporter{
		client:  configCatClient,
		metrics: prometheusMetrics,
		logger:  logger,
		config:  cfg,
	}
}

// scrapeMetrics collects metrics from ConfigCat API
func (e *Exporter) scrapeMetrics(ctx context.Context) error {
	timer := prometheus.NewTimer(e.metrics.ScrapeDuration)
	defer timer.ObserveDuration()

	e.mu.Lock()
	defer e.mu.Unlock()

	e.logger.Info("Starting ConfigCat metrics scrape")

	// For single product setup, we set products to 1
	e.metrics.ProductsTotal.Set(1)

	// Get configs for the product
	configs, err := e.client.GetConfigs(ctx, e.config.ProductID)
	if err != nil {
		e.metrics.ScrapeErrors.Inc()
		return fmt.Errorf("getting configs: %w", err)
	}

	// Get environments for the product
	environments, err := e.client.GetEnvironments(ctx, e.config.ProductID)
	if err != nil {
		e.metrics.ScrapeErrors.Inc()
		return fmt.Errorf("getting environments: %w", err)
	}

	zombieflags, err := e.client.GetZombieFlags(ctx, e.config.ProductID)
	if err != nil {
		e.metrics.ScrapeErrors.Inc()
		return fmt.Errorf("getting zombieflags: %w", err)
	}

	// Set configs and environments metrics
	productName := fmt.Sprintf("product-%s", e.config.ProductID) // We could enhance this by fetching actual product name
	e.metrics.ConfigsTotal.WithLabelValues(e.config.ProductID, productName).Set(float64(len(configs)))
	e.metrics.EnvironmentsTotal.WithLabelValues(e.config.ProductID, productName).Set(float64(len(environments)))

	var zombieflagCount int
	for _, zombieConfig := range zombieflags.Configs {
		zombieflagCount += len(zombieConfig.Settings)

		e.metrics.ZombieFlagsTotal.WithLabelValues(
			e.config.ProductID,
			productName,
			zombieConfig.ConfigID,
			zombieConfig.Name).Set(float64(zombieflagCount))
	}

	// Get feature flags for each config
	for _, config := range configs {
		featureFlags, err := e.client.GetFeatureFlags(ctx, config.ConfigID)
		if err != nil {
			e.logger.Error("Failed to get feature flags for config",
				slog.String("config_id", config.ConfigID),
				slog.String("config_name", config.Name),
				slog.String("error", err.Error()))
			e.metrics.ScrapeErrors.Inc()
			continue
		}

		e.metrics.FeatureFlagsTotal.WithLabelValues(
			e.config.ProductID,
			productName,
			config.ConfigID,
			config.Name,
		).Set(float64(len(featureFlags)))

		e.logger.Debug("Collected feature flags",
			slog.String("config_name", config.Name),
			slog.String("config_id", config.ConfigID),
			slog.Int("feature_flags_count", len(featureFlags)))
	}

	e.metrics.LastScrapeTime.SetToCurrentTime()

	e.logger.Info("ConfigCat metrics scrape completed",
		slog.Int("configs_count", len(configs)),
		slog.Int("environments_count", len(environments)),
		slog.Int("zombieflags_count", zombieflagCount))

	return nil
}

// StartScraping starts the periodic metrics collection
func (e *Exporter) StartScraping(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(e.config.ScrapeInterval) * time.Second)
	defer ticker.Stop()

	e.logger.Info("Starting metrics scraper",
		slog.Int("interval_seconds", e.config.ScrapeInterval))

	// Initial scrape
	if err := e.scrapeMetrics(ctx); err != nil {
		e.logger.Error("Initial scrape failed", slog.String("error", err.Error()))
	}

	for {
		select {
		case <-ctx.Done():
			e.logger.Info("Scraping stopped")
			return
		case <-ticker.C:
			if err := e.scrapeMetrics(ctx); err != nil {
				e.logger.Error("Scrape failed", slog.String("error", err.Error()))
			}
		}
	}
}
