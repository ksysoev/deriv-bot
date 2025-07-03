package executor

import (
	"context"
	"fmt"

	"github.com/ksysoev/deriv-bot/pkg/core/signal"
)

type MarketSignals interface {
	SubscribeOnMarket(ctx context.Context, symbol string) (<-chan signal.Tick, error)
}

type TradingProvider interface {
	Authorize(ctx context.Context, token string) error
	Buy(ctx context.Context, symbol string, amount float64, price float64, leverage int) (int, error)
	Sell(ctx context.Context, symbol string, amount float64, price float64, leverage int) (int, error)
	ClosePosition(ctx context.Context, contractID int) error
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
func (s *Service) ExecuteStrategy(ctx context.Context, stategy Strategy) error {
	if err := s.tradingProv.Authorize(ctx, stategy.Token); err != nil {
		return fmt.Errorf("failed to authorize trading provider: %w", err)
	}

	tickChan, err := s.marketSignals.SubscribeOnMarket(ctx, stategy.Symbol)
	if err != nil {
		return err
	}

	cid := 0

	for tick := range tickChan {
		if cid == 0 && stategy.CheckToOpen(tick) {
			var err error

			switch stategy.Type {
			case StrategyTypeBuy:
				cid, err = s.tradingProv.Buy(ctx, stategy.Symbol, stategy.Amount, tick.Quote, stategy.Leverage)
			case StrategyTypeSell:
				cid, err = s.tradingProv.Sell(ctx, stategy.Symbol, stategy.Amount, tick.Quote, stategy.Leverage)
			case StrategyTypeNotSet:
				return fmt.Errorf("strategy type not set")
			default:
				return fmt.Errorf("unknown strategy type: %d", stategy.Type)
			}

			if err != nil {
				return fmt.Errorf("failed to execute buy for symbol %s: %w", stategy.Amount, err)
			}

			continue
		}

		if cid != 0 && stategy.CheckToClose(tick) {
			if err := s.tradingProv.ClosePosition(ctx, cid); err != nil {
				return fmt.Errorf("failed to close position for contract ID %d: %w", cid, err)
			}

			cid = 0 // Reset contract ID after closing position

			continue
		}
	}

	return nil
}
