package deriv

import (
	"context"
	"fmt"

	"github.com/ksysoev/deriv-api/schema"
)

// Buy places a buy order for a specified symbol with given parameters.
// It uses the provided price, amount, and leverage to configure the order.
// Accepts ctx for request lifecycle management, symbol, the asset's market symbol, amount as the quantity to buy, price for the transaction, and leverage specifying the multiplier.
// Returns the contract ID of the placed buy order and an error if the order fails due to API issues or invalid parameters.
func (a *API) Buy(ctx context.Context, symbol string, amount float64, price float64, leverage int) (int, error) {
	lev := float64(leverage)

	res, err := a.client.Buy(ctx, schema.Buy{
		Price: price,
		Parameters: &schema.BuyParameters{
			ContractType: schema.BuyParametersContractTypeMULTUP,
			Symbol:       symbol,
			Amount:       &amount,
			Multiplier:   &lev,
		},
	})

	if err != nil {
		return 0, fmt.Errorf("failed to place buy order for symbol %s: %w", symbol, err)
	}

	return res.Buy.ContractId, nil
}

// Sell places a sell order for the specified symbol with the provided parameters.
// It uses the given price, amount, and leverage to configure the order.
// Accepts ctx for request lifecycle management, symbol for the market asset, amount as the quantity to sell, price per unit, and leverage for multiplier configuration.
// Returns the contract ID of the placed sell order and an error if the order fails due to API issues or invalid parameters.
func (a *API) Sell(ctx context.Context, symbol string, amount float64, price float64, leverage int) (int, error) {
	lev := float64(leverage)

	res, err := a.client.Buy(ctx, schema.Buy{
		Price: price,
		Parameters: &schema.BuyParameters{
			ContractType: schema.BuyParametersContractTypeMULTDOWN,
			Symbol:       symbol,
			Amount:       &amount,
			Multiplier:   &lev,
		},
	})

	if err != nil {
		return 0, fmt.Errorf("failed to place sell order for symbol %s: %w", symbol, err)
	}

	return res.Buy.ContractId, nil
}
