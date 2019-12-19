package fee

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/types"
)

func DefaultCalculateFee(ctx sdkTypes.Context, feeSetting *FeeSetting, amt sdkTypes.Coins) (sdkTypes.Coins, sdkTypes.Error) {

	amount := amt.AmountOf(types.CIN)
	minFee := feeSetting.Min.AmountOf(types.CIN)
	maxFee := feeSetting.Max.AmountOf(types.CIN)
	percentageStr := feeSetting.Percentage

	bigIntPercentage, ok := sdkTypes.NewIntFromString(percentageStr)
	if !ok {
		return nil, sdkTypes.ErrInternal(fmt.Sprintf("Invalid percentage: %s", percentageStr))
	}

	// percentage in the genesis file always power of -4
	multiplier, _ := sdkTypes.NewIntFromString("10000")
	fee := amount.Mul(bigIntPercentage)
	fee = fee.Quo(multiplier)

	if fee.LT(minFee) {
		fee = minFee
	}
	if fee.GT(maxFee) {
		fee = maxFee
	}

	if fee.IsZero() {
		zero, _ := sdkTypes.NewIntFromString("0")
		var fees sdkTypes.Coins
		fees = sdkTypes.Coins{
			{
				Denom:  types.CIN,
				Amount: zero,
			},
		}

		return fees, nil
	}

	return sdkTypes.Coins{sdkTypes.NewCoin(types.CIN, fee)}, nil
}
