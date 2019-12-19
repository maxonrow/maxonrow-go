package bank

import (
	"fmt"

	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/fee"
)

const (
	QueryGetTransferFee = "get_fee"
)

func NewQuerier(feeKeeper fee.Keeper) sdkTypes.Querier {
	return func(ctx sdkTypes.Context, path []string, req abci.RequestQuery) ([]byte, sdkTypes.Error) {
		switch path[0] {
		case QueryGetTransferFee:
			return queryGetTransferFee(ctx, path[1:], req, feeKeeper)

		default:
			return nil, sdkTypes.ErrUnknownRequest("unknown bank query endpoint")
		}
	}
}

func queryGetTransferFee(ctx sdkTypes.Context, path []string, req abci.RequestQuery, feeKeeper fee.Keeper) ([]byte, sdkTypes.Error) {
	if len(path) != 2 {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid path %s", strings.Join(path, "/")))
	}

	msgType := path[0]

	totalAmount, parseCoinsErr := sdkTypes.ParseCoins(path[1] + types.CIN)
	if parseCoinsErr != nil {
		return nil, sdkTypes.ErrInvalidCoins("Invalid amount")
	}

	feeSetting, feeSettingErr := feeKeeper.GetMsgFeeSetting(ctx, msgType)
	if feeSettingErr != nil {
		return nil, feeSettingErr
	}

	calculatedFee, err := fee.DefaultCalculateFee(ctx, feeSetting, totalAmount)
	if err != nil {
		return nil, err
	}

	respData := sdkTypes.MustSortJSON(codec.Cdc.MustMarshalJSON(calculatedFee))

	return respData, nil

}
