package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/dirsigler/configcat-exporter/internal/config"
	"github.com/dirsigler/configcat-exporter/internal/exporter"
	"github.com/dirsigler/configcat-exporter/internal/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "configcat-exporter",
		Usage: "Export ConfigCat metrics to Prometheus",
		Flags: config.CLIFlags(),
		Action: func(c *cli.Context) error {
			// Load configuration
			cfg, err := config.Load(c)
			if err != nil {
				return err
			}

			// Setup logger
			logger := cfg.SetupLogger()

			logger.Info("Starting ConfigCat Prometheus Exporter",
				slog.Int("port", cfg.Port),
				slog.Int("scrape_interval", cfg.ScrapeInterval),
				slog.String("log_level", cfg.LogLevel),
				slog.String("organization_id", cfg.OrganizationID),
				slog.String("product_id", cfg.ProductID))

			// Create exporter
			exp := exporter.New(cfg, logger)

			// Setup HTTP server
			mux := http.NewServeMux()
			mux.HandleFunc("/", handler.IndexHandler)
			mux.HandleFunc("/health", handler.HealthHandler)
			mux.Handle("/metrics", promhttp.Handler())

			server := &http.Server{
				Addr:    ":" + strconv.Itoa(cfg.Port),
				Handler: mux,
			}

			// Setup context for graceful shutdown
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Start metrics scraping
			go exp.StartScraping(ctx)

			// Handle shutdown signals
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			// Start HTTP server
			go func() {
				logger.Info("HTTP server starting", slog.Int("port", cfg.Port))
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP server error", slog.String("error", err.Error()))
				}
			}()

			// Wait for shutdown signal
			<-sigChan
			logger.Info("Shutdown signal received")

			// Cancel scraping context
			cancel()

			// Shutdown HTTP server
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()

			if err := server.Shutdown(shutdownCtx); err != nil {
				logger.Error("Server shutdown error", slog.String("error", err.Error()))
			}

			logger.Info("Exporter stopped")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("Application error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
