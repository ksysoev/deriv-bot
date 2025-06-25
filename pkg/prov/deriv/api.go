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

// New initializes and returns a new API instance using the provided configuration or returns an error if creation fails.
func New(cfg Config) (*API, error) {
	client, err := deriv.NewDerivAPI(cfg.Endpoint, cfg.AppID, defaultLanguage, cfg.Origin)
	if err != nil {
		return nil, fmt.Errorf("failed to create Deriv API client: %w", err)
	}

	return &API{client: client}, nil
}
