package nonfungible

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/maxonrow/maxonrow-go/x/fee"
)

const (
	QueryListTokenSymbol     = "list_token_symbol"
	QueryTokenData           = "token_data"
	QueryAccount             = "account"
	QueryGetFee              = "get_fee"
	QueryGetTokenTransferFee = "get_token_transfer_fee"
)

func NewQuerier(cdc *codec.Codec, keeper *Keeper, feeKeeper *fee.Keeper) sdkTypes.Querier {
	return func(ctx sdkTypes.Context, path []string, req abci.RequestQuery) ([]byte, sdkTypes.Error) {
		switch path[0] {
		case QueryListTokenSymbol:
			return queryListTokenSymbol(cdc, ctx, path[1:], req, keeper)
		case QueryTokenData:
			return queryTokenData(cdc, ctx, path[1:], req, keeper)
		default:
			return nil, sdkTypes.ErrUnknownRequest("unknown token query endpoint")
		}
	}
}

func queryListTokenSymbol(cdc *codec.Codec, ctx sdkTypes.Context, _ []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {
	fungibleToken, nonfungibleToken := keeper.ListTokenData(ctx)

	resp := &listTokenSymbolResponse{
		Fungible:    fungibleToken,
		Nonfungible: nonfungibleToken,
	}
	respData := cdc.MustMarshalJSON(resp)

	return respData, nil
}

func queryTokenData(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {
	if len(path) != 1 {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid path %s", strings.Join(path, "/")))
	}

	symbol := path[0]

	tokenData, err := keeper.GetTokenData(ctx, symbol)
	if err != nil {
		return nil, err
	}

	tokenInfo := cdc.MustMarshalJSON(tokenData)

	return tokenInfo, nil
}

type listTokenSymbolResponse struct {
	Fungible    []string `json:"fungible"`
	Nonfungible []string `json:"nonfungible"`
}
