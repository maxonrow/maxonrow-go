package types

import (
	"fmt"
	"strings"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	CINperMXW = 1000000000000000000
	CIN       = "cin"
	MXW       = "mxw"
)

func init() {
	sdkTypes.RegisterDenom(MXW, sdkTypes.NewDec(CINperMXW))
	sdkTypes.RegisterDenom(CIN, sdkTypes.NewDec(1))
}

func MXWtoCIN(coin sdkTypes.Coin) sdkTypes.Coin {
	if coin.Denom != MXW {
		panic(fmt.Sprintf("Invalid token: %s", coin.Denom))
	}

	return sdkTypes.NewCoin(CIN, coin.Amount.Mul(sdkTypes.NewInt(CINperMXW)))
}

func ParseCoins(coinsStr string) (sdkTypes.Coins, error) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return nil, nil
	}

	coinStrs := strings.Split(coinsStr, ",")
	coins := make(sdkTypes.Coins, len(coinStrs))
	for i, coinStr := range coinStrs {
		coin, err := sdkTypes.ParseCoin(coinStr)
		if err != nil {
			return nil, err
		}

		coins[i] = coin
	}

	// sort coins for determinism
	coins.Sort()

	return coins, nil
}
