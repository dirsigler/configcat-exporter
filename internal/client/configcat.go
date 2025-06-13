package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/dirsigler/configcat-exporter/pkg/types"
)

// ConfigCatClient handles API interactions with ConfigCat
type ConfigCatClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
	logger  *slog.Logger
}

// New creates a new ConfigCat API client
func New(apiKey, baseURL string, logger *slog.Logger) *ConfigCatClient {
	return &ConfigCatClient{
		apiKey:  apiKey,
		baseURL: strings.TrimSuffix(baseURL, "/"),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// makeRequest performs authenticated HTTP requests to ConfigCat API
func (c *ConfigCatClient) makeRequest(ctx context.Context, endpoint string, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ConfigCat-Exporter - https://github.com/dirsigler/configcat-exporter/1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d for endpoint %s", resp.StatusCode, endpoint)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	return nil
}

// GetConfigs retrieves all configs for a product
func (c *ConfigCatClient) GetConfigs(ctx context.Context, productID string) ([]types.Config, error) {
	var configs []types.Config
	endpoint := fmt.Sprintf("/v1/products/%s/configs", productID)

	if err := c.makeRequest(ctx, endpoint, &configs); err != nil {
		return nil, fmt.Errorf("getting configs for product %s: %w", productID, err)
	}

	c.logger.Debug("Retrieved configs",
		slog.String("product_id", productID),
		slog.Int("count", len(configs)))

	return configs, nil
}

// GetEnvironments retrieves all environments for a product
func (c *ConfigCatClient) GetEnvironments(ctx context.Context, productID string) ([]types.Environment, error) {
	var environments []types.Environment
	endpoint := fmt.Sprintf("/v1/products/%s/environments", productID)

	if err := c.makeRequest(ctx, endpoint, &environments); err != nil {
		return nil, fmt.Errorf("getting environments for product %s: %w", productID, err)
	}

	c.logger.Debug("Retrieved environments",
		slog.String("product_id", productID),
		slog.Int("count", len(environments)))

	return environments, nil
}

// GetFeatureFlags retrieves all feature flags for a config
func (c *ConfigCatClient) GetFeatureFlags(ctx context.Context, configID string) ([]types.FeatureFlag, error) {
	var featureFlags []types.FeatureFlag
	endpoint := fmt.Sprintf("/v1/configs/%s/settings", configID)

	if err := c.makeRequest(ctx, endpoint, &featureFlags); err != nil {
		return nil, fmt.Errorf("getting feature flags for config %s: %w", configID, err)
	}

	c.logger.Debug("Retrieved feature flags",
		slog.String("config_id", configID),
		slog.Int("count", len(featureFlags)))

	return featureFlags, nil
}

// GetZombieFlags retrieves all zombieflags for a product
func (c *ConfigCatClient) GetZombieFlags(ctx context.Context, productID string) (types.ZombieFlag, error) {
	var zombieFlags types.ZombieFlag
	endpoint := fmt.Sprintf("/v1/products/%s/staleflags", productID)

	if err := c.makeRequest(ctx, endpoint, &zombieFlags); err != nil {
		return types.ZombieFlag{}, fmt.Errorf("getting zombie flags for product %s: %w", productID, err)
	}

	c.logger.Debug("Retrieved zombie flags",
		slog.String("product_id", productID))

	return zombieFlags, nil
}
