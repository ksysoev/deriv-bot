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
func (a *API) Buy(ctx context.Context, symbol string, amount, price float64, leverage int) (int, error) {
	lev := float64(leverage)
	basis := schema.BuyParametersBasisStake

	res, err := a.client.Buy(ctx, schema.Buy{
		Buy:   "1",
		Price: price,
		Parameters: &schema.BuyParameters{
			ContractType: schema.BuyParametersContractTypeMULTUP,
			Basis:        &basis,
			Symbol:       symbol,
			Amount:       &amount,
			ProductType:  schema.BuyParametersProductTypeBasic,
			Multiplier:   &lev,
			Currency:     "USD", // Hardcoded to USD, for simplicity, in future need to accept as a parameter
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
func (a *API) Sell(ctx context.Context, symbol string, amount, price float64, leverage int) (int, error) {
	lev := float64(leverage)
	basis := schema.BuyParametersBasisStake

	res, err := a.client.Buy(ctx, schema.Buy{
		Price: price,
		Parameters: &schema.BuyParameters{
			ContractType: schema.BuyParametersContractTypeMULTDOWN,
			Basis:        &basis,
			Symbol:       symbol,
			Amount:       &amount,
			ProductType:  schema.BuyParametersProductTypeBasic,
			Multiplier:   &lev,
			Currency:     "USD", // Hardcoded to USD, for simplicity, in future need to accept as a parameter
		},
	})

	if err != nil {
		return 0, fmt.Errorf("failed to place sell order for symbol %s: %w", symbol, err)
	}

	return res.Buy.ContractId, nil
}

// ClosePosition closes an open trading position for a given contract ID.
// It performs a sell operation at the market price to close the position.
// Accepts ctx to manage request lifecycle and contractID identifying the position to close.
// Returns an error if the API request to close the position fails.
func (a *API) ClosePosition(ctx context.Context, contractID int) error {
	_, err := a.client.Sell(ctx, schema.Sell{
		Sell:  contractID,
		Price: 0, // Sell at market price, we may want to allow specifying a price in the future
	})

	if err != nil {
		return fmt.Errorf("failed to close position %d: %w", contractID, err)
	}

	return nil
}
