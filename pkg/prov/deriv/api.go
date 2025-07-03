package deriv

import (
	"context"
	"fmt"
	"sync"

	"github.com/ksysoev/deriv-api"
	"github.com/ksysoev/deriv-api/schema"
	"github.com/ksysoev/deriv-bot/pkg/core/executor"
)

const (
	defaultLanguage = "en"
)

type Config struct {
	Endpoint string `mapstructure:"endpoint"`
	Origin   string `mapstructure:"origin"`
	AppID    int    `mapstructure:"app_id"`
}

type API struct {
	client *deriv.Client
	wg     sync.WaitGroup
}

// New creates a new API instance using the provided configuration.
// It validates the configuration and initializes a Deriv API client.
// Returns the initialized API instance and an error if client creation fails or the configuration is invalid.
func New(cfg Config) (*API, error) {
	client, err := deriv.NewDerivAPI(cfg.Endpoint, cfg.AppID, defaultLanguage, cfg.Origin, deriv.Debug)
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
	a.wg.Wait()
}

func (a *API) Authorize(ctx context.Context, token string) (*executor.Account, error) {
	res, err := a.client.Authorize(ctx, schema.Authorize{Authorize: token})
	if err != nil {
		return nil, fmt.Errorf("failed to authorize with Deriv API: %w", err)
	}

	return &executor.Account{
		ID:       *res.Authorize.Loginid,
		Currency: *res.Authorize.Currency,
	}, nil
}
