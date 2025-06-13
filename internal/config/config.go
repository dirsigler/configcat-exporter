package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/urfave/cli/v2"
)

// Config holds the application configuration
type Config struct {
	ConfigCatAPIKey string `env:"CONFIGCAT_API_KEY,required"`
	Port            int    `env:"PORT" envDefault:"8080"`
	ScrapeInterval  int    `env:"SCRAPE_INTERVAL" envDefault:"60"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info"`
	ConfigCatAPIURL string `env:"CONFIGCAT_API_URL" envDefault:"https://api.configcat.com"`
	OrganizationID  string `env:"CONFIGCAT_ORGANIZATION_ID,required"`
	ProductID       string `env:"CONFIGCAT_PRODUCT_ID,required"`
}

// Load parses configuration from environment variables and CLI context
func Load(c *cli.Context) (*Config, error) {
	// Parse configuration from environment variables
	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("parsing config from environment: %w", err)
	}

	// Override with CLI flags if provided
	if c.String("configcat-api-key") != "" {
		config.ConfigCatAPIKey = c.String("configcat-api-key")
	}
	if c.Int("port") != 8080 {
		config.Port = c.Int("port")
	}
	if c.Int("scrape-interval") != 60 {
		config.ScrapeInterval = c.Int("scrape-interval")
	}
	if c.String("log-level") != "info" {
		config.LogLevel = c.String("log-level")
	}
	if c.String("organization-id") != "" {
		config.OrganizationID = c.String("organization-id")
	}
	if c.String("product-id") != "" {
		config.ProductID = c.String("product-id")
	}

	// Validate required configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks if all required configuration is present
func (c *Config) Validate() error {
	if c.ConfigCatAPIKey == "" {
		return fmt.Errorf("ConfigCat API key is required")
	}
	if c.OrganizationID == "" {
		return fmt.Errorf("ConfigCat Organization ID is required")
	}
	if c.ProductID == "" {
		return fmt.Errorf("ConfigCat Product ID is required")
	}
	return nil
}

// SetupLogger configures the structured logger based on config
func (c *Config) SetupLogger() *slog.Logger {
	var logLevel slog.Level
	switch strings.ToLower(c.LogLevel) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

// CLIFlags returns the CLI flags for the application
func CLIFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "configcat-api-key",
			Usage:   "ConfigCat API key (can also be set via CONFIGCAT_API_KEY env var)",
			EnvVars: []string{"CONFIGCAT_API_KEY"},
		},
		&cli.IntFlag{
			Name:    "port",
			Usage:   "Port to listen on",
			EnvVars: []string{"PORT"},
			Value:   8080,
		},
		&cli.IntFlag{
			Name:    "scrape-interval",
			Usage:   "Scrape interval in seconds",
			EnvVars: []string{"SCRAPE_INTERVAL"},
			Value:   60,
		},
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "Log level (debug, info, warn, error)",
			EnvVars: []string{"LOG_LEVEL"},
			Value:   "info",
		},
		&cli.StringFlag{
			Name:    "organization-id",
			Usage:   "ConfigCat Organization ID",
			EnvVars: []string{"CONFIGCAT_ORGANIZATION_ID"},
		},
		&cli.StringFlag{
			Name:    "product-id",
			Usage:   "ConfigCat Product ID",
			EnvVars: []string{"CONFIGCAT_PRODUCT_ID"},
		},
	}
}
