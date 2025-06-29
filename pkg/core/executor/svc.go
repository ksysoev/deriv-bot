package executor

import (
	"context"
	"log/slog"

	"github.com/ksysoev/deriv-bot/pkg/core/signal"
)

type MarketSignals interface {
	SubscribeOnMarket(ctx context.Context, symbol string) (<-chan signal.Tick, error)
}

type TradingProvider interface {
	Buy(ctx context.Context, symbol string, amount float64, price float64, leverage int) (int, error)
	Sell(ctx context.Context, symbol string, amount float64, price float64, leverage int) (int, error)
}

type Service struct {
	marketSignals MarketSignals
	tradingProv   TradingProvider
}

// New creates and returns a new Service instance with the provided marketSignals and tradingProv dependencies.
// marketSignals provides market data subscription capabilities.
// tradingProv handles trading operations like buy and sell.
func New(marketSignals MarketSignals, tradingProv TradingProvider) *Service {
	return &Service{
		marketSignals: marketSignals,
		tradingProv:   tradingProv,
	}
}

// ExecuteStrategy monitors market signals for a given symbol and executes a buy operation when the evaluation condition is met.
// It subscribes to market signals and iterates through incoming ticks. If `eval` returns true for a tick, the service buys the symbol.
// ctx is the context for managing the subscription and operation lifecycle.
// symbol specifies the market symbol to trade, and amount is the quantity to buy.
// eval is a callback function that evaluates tick data to decide when to trigger the buy action.
// Returns the transaction ID of the buy operation and an error if subscribing to market signals or executing the buy fails.
func (s *Service) ExecuteStrategy(ctx context.Context, symbol string, amount float64, eval func(tick signal.Tick) bool) error {
	tickChan, err := s.marketSignals.SubscribeOnMarket(ctx, symbol)
	if err != nil {
		return err
	}

	for tick := range tickChan {
		slog.Info("Received tick", slog.Any("tick", tick))
		if eval(tick) {
			_, err := s.tradingProv.Buy(ctx, symbol, amount, tick.Quote, 10)

			return err
		}
	}

	return nil
}
