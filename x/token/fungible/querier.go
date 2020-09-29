package fungible

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/x/fee"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryListTokenSymbol                     = "list_token_symbol"
	QueryTokenData                           = "token_data"
	QueryAccount                             = "account"
	QueryGetFee                              = "get_fee"
	QueryGetTokenTransferFee                 = "get_token_transfer_fee"
	QueryGetFungibleTokenMaintainerAddresses = "get_token_maintainer_addresses"
)

func NewQuerier(cdc *codec.Codec, keeper *Keeper, feeKeeper *fee.Keeper) sdkTypes.Querier {
	return func(ctx sdkTypes.Context, path []string, req abci.RequestQuery) ([]byte, sdkTypes.Error) {
		switch path[0] {
		case QueryListTokenSymbol:
			return queryListTokenSymbol(cdc, ctx, path[1:], req, keeper)
		case QueryTokenData:
			return queryTokenData(cdc, ctx, path[1:], req, keeper)
		case QueryAccount:
			return queryAccount(cdc, ctx, path[1:], req, keeper)
		case QueryGetFungibleTokenMaintainerAddresses:
			return queryGetFungibleTokenMaintainerAddresses(ctx, path[1:], req, keeper)
		default:
			return nil, sdkTypes.ErrUnknownRequest("unknown token query endpoint")
		}
	}
}

func queryListTokenSymbol(cdc *codec.Codec, ctx sdkTypes.Context, _ []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {
	tokens := keeper.ListTokens(ctx)

	var symbols []string
	for _, t := range tokens {
		symbols = append(symbols, t.Symbol)
	}

	respData := cdc.MustMarshalJSON(symbols)

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

func queryAccount(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {
	if len(path) != 2 {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid path %s", strings.Join(path, "/")))
	}

	symbol := path[0]
	accountBech := path[1]

	account, err := sdkTypes.AccAddressFromBech32(accountBech)
	if err != nil {
		return nil, sdkTypes.ErrInvalidAddress("Invalid account address")
	}

	acc, err := keeper.GetAccount(ctx, symbol, account)
	if err != nil {
		return nil, sdkTypes.ErrInternal(err.Error())
	}

	if acc == nil {
		return nil, nil
	}

	accountData := cdc.MustMarshalJSON(acc)

	return accountData, nil
}

type listTokenSymbolResponse struct {
	Fungible    []string `json:"fungible"`
	Nonfungible []string `json:"nonfungible"`
}

type FungibleTokenMaintainerSetting struct {
	Module      string              `json:"module"`
	Maintainers []MaintainerSetting `json:"maintainers"`
}

type MaintainerSetting struct {
	Type    string                `json:"type"`
	Address []sdkTypes.AccAddress `json:"address"`
}

func queryGetFungibleTokenMaintainerAddresses(ctx sdkTypes.Context, path []string, req abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {

	var fungibleTokenMaintainerSettings []MaintainerSetting

	fungibleTokenIssuerAddresses := keeper.GetIssuerAddresses(ctx)
	maintainerData := MaintainerSetting{"Issuer", fungibleTokenIssuerAddresses}
	fungibleTokenMaintainerSettings = append(fungibleTokenMaintainerSettings, maintainerData)

	fungibleTokenProviderAddresses := keeper.GetProviderAddresses(ctx)
	maintainerData = MaintainerSetting{"Provider", fungibleTokenProviderAddresses}
	fungibleTokenMaintainerSettings = append(fungibleTokenMaintainerSettings, maintainerData)

	fungibleTokenAuthorisedAddresses := keeper.GetAuthorisedAddresses(ctx)
	maintainerData = MaintainerSetting{"Middleware", fungibleTokenAuthorisedAddresses}
	fungibleTokenMaintainerSettings = append(fungibleTokenMaintainerSettings, maintainerData)

	fungibleTokenMaintainerSetting := FungibleTokenMaintainerSetting{"token", fungibleTokenMaintainerSettings}
	respData := codec.Cdc.MustMarshalJSON(fungibleTokenMaintainerSetting)
	return respData, nil
}
