package deriv

import (
	"fmt"

	"github.com/ksysoev/deriv-api"
)

const (
	defaultLanguage = "en"
)

type Config struct {
	AppID    int    `mapstructure:"app_id"`
	Endpoint string `mapstructure:"endpoint"`
	Origin   string `mapstructure:"origin"`
}

type API struct {
	client *deriv.Client
}

// New creates a new API instance using the provided configuration.
// It validates the configuration and initializes a Deriv API client.
// Returns the initialized API instance and an error if client creation fails or the configuration is invalid.
func New(cfg Config) (*API, error) {
	client, err := deriv.NewDerivAPI(cfg.Endpoint, cfg.AppID, defaultLanguage, cfg.Origin)
	if err != nil {
		return nil, fmt.Errorf("failed to create Deriv API client: %w", err)
	}

	return &API{client: client}, nil
}

// Close releases all resources associated with the API instance.
// It disconnects the underlying client and should be called to clean up properly.
// Returns an error if disconnection fails.
func (a *API) Close() {
	a.client.Disconnect()
}
