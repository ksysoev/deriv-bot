package signal

import (
	"context"
	"fmt"
)

type MarketProvider interface {
	SubscribeToTicks(ctx context.Context, symbol string) (<-chan Tick, error)
}

type SubscribtionManager interface {
	GetMarketSubscription(symbol string) (<-chan Tick, bool)
	SetMarketSubscription(symbol string, sub <-chan Tick)
}

type Service struct {
	markerProv MarketProvider
	subMgr     SubscribtionManager
}

// New creates and initializes a new Service instance with the provided MarketProvider.
// It requires a valid prov implementing the MarketProvider interface.
// Returns a pointer to the newly created Service.
func New(prov MarketProvider, subMgr SubscribtionManager) *Service {
	return &Service{
		markerProv: prov,
		subMgr:     subMgr,
	}
}

// SubscribeOnMarket subscribes to real-time market updates for the specified symbol and provides ticks via a channel.
// It connects to the underlying market provider and streams tick data to subscribers.
// ctx is the context to control cancellation or timeout, and symbol specifies the market symbol of interest.
// Returns a read-only channel streaming Tick updates and an error if the subscription fails.
func (s *Service) SubscribeOnMarket(ctx context.Context, symbol string) (<-chan Tick, error) {
	tickChan, err := s.markerProv.SubscribeToTicks(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to ticks for symbol %s: %w", symbol, err)
	}

	// TODO: here we'll implement logic for sending ticks to multiple subscribers

	return tickChan, nil
}
