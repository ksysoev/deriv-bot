package deriv

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ksysoev/deriv-api"
	"github.com/ksysoev/deriv-api/schema"
	"github.com/ksysoev/deriv-bot/pkg/core/signal"
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
	a.wg.Wait()
}

// SubscribeToTicks subscribes to real-time tick data for the specified symbol using the provided context.
// It listens for tick data updates and streams them through a channel of signal.Tick.
// Accepts ctx for managing subscription lifecycle and symbol, the market symbol to subscribe to.
// Returns a read-only channel of signal.Tick containing streaming tick updates and an error if the subscription fails.
func (a *API) SubscribeToTicks(ctx context.Context, symbol string) (<-chan signal.Tick, error) {
	_, sub, err := a.client.SubscribeTicks(ctx, schema.Ticks{Ticks: symbol})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to ticks for symbol %s: %w", symbol, err)
	}

	subChan := sub.GetStream()
	resChan := make(chan signal.Tick)

	a.wg.Add(1)

	go func() {
		defer a.wg.Done()
		defer func() { _ = sub.Forget() }()
		defer close(resChan)

		for {
			select {
			case <-ctx.Done():
				return
			case tick, ok := <-subChan:
				if !ok {
					return
				}

				epoch := time.Unix(int64(*tick.Tick.Epoch), 0)

				resChan <- signal.Tick{
					Time:  epoch,
					Quote: *tick.Tick.Quote,
					Ask:   *tick.Tick.Ask,
					Bid:   *tick.Tick.Bid,
				}
			}
		}
	}()

	return resChan, nil
}
