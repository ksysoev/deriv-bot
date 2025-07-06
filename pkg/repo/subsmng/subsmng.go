package subsmng

import (
	"sync"

	"github.com/ksysoev/deriv-bot/pkg/core/signal"
)

type SubscriptionManager struct {
	subs map[string]<-chan signal.Tick
	mu   sync.Mutex
}

// New creates and initializes a new SubscriptionManager instance.
// It allocates the necessary internal map to manage subscriptions and ensures readiness for use.
// Returns a pointer to a SubscriptionManager configured with an empty subscription map.
func New() *SubscriptionManager {
	return &SubscriptionManager{
		subs: make(map[string]<-chan signal.Tick),
	}
}

// GetMarketSubscription retrieves the subscription channel for a specific market symbol if it exists.
// It locks the subscription manager during execution to ensure thread safety.
// Takes symbol, the market symbol to search for in the subscription map.
// Returns a channel of signal.Tick for the specified market symbol if a subscription exists, otherwise returns false.
func (s *SubscriptionManager) GetMarketSubscription(symbol string) (<-chan signal.Tick, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub, ok := s.subs[symbol]
	return sub, ok
}

// SetMarketSubscription registers a subscription channel for a given market symbol, overriding any existing subscription.
// It safely updates the internal subscription map while ensuring thread safety using a mutex.
// Takes symbol, the market symbol used as a key, and sub, the channel of signal.Tick used for subscription updates.
func (s *SubscriptionManager) SetMarketSubscription(symbol string, sub <-chan signal.Tick) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.subs[symbol] = sub
}
