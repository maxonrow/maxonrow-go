package nonfungible

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/x/fee"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryListTokenSymbol                        = "list_token_symbol"
	QueryTokenData                              = "token_data"
	QueryItemData                               = "item_data"
	QueryAccount                                = "account"
	QueryGetFee                                 = "get_fee"
	QueryGetTokenTransferFee                    = "get_token_transfer_fee"
	QueryEndorserList                           = "get_endorser_list"
	QueryGetNonfungibleTokenMaintainerAddresses = "get_nonfungible_maintainer_addresses"
)

type ItemInfo struct {
	Owner         sdkTypes.AccAddress
	ID            string
	Properties    string
	Metadata      string
	TransferLimit sdkTypes.Uint
	Frozen        bool
}

func NewQuerier(cdc *codec.Codec, keeper *Keeper, feeKeeper *fee.Keeper) sdkTypes.Querier {
	return func(ctx sdkTypes.Context, path []string, req abci.RequestQuery) ([]byte, sdkTypes.Error) {
		switch path[0] {
		case QueryListTokenSymbol:
			return queryListTokenSymbol(cdc, ctx, path[1:], req, keeper)
		case QueryTokenData:
			return queryTokenData(cdc, ctx, path[1:], req, keeper)
		case QueryItemData:
			return queryItemData(cdc, ctx, path[1:], req, keeper)
		case QueryEndorserList:
			return queryEndorserList(cdc, ctx, path[1:], req, keeper)
		case QueryGetNonfungibleTokenMaintainerAddresses:
			return queryGetNonfungibleTokenMaintainerAddresses(cdc, ctx, path[1:], req, keeper)
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

func queryItemData(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {
	if len(path) != 2 {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid path %s", strings.Join(path, "/")))
	}

	symbol := path[0]
	itemID := path[1]

	item := keeper.GetNonFungibleItem(ctx, symbol, itemID)
	owner := keeper.GetNonFungibleItemOwnerInfo(ctx, symbol, itemID)

	var itemInfo = ItemInfo{
		Owner:         owner,
		ID:            item.ID,
		Properties:    item.Properties,
		Metadata:      item.Metadata,
		TransferLimit: item.TransferLimit,
		Frozen:        item.Frozen,
	}

	js := cdc.MustMarshalJSON(itemInfo)

	return js, nil
}

func queryEndorserList(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {
	if len(path) != 1 {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid path %s", strings.Join(path, "/")))
	}

	symbol := path[0]

	endorserList := keeper.GetEndorserList(ctx, symbol)
	if endorserList != nil {
		return cdc.MustMarshalJSON(endorserList), nil
	}

	return nil, nil
}

type listTokenSymbolResponse struct {
	Fungible    []string `json:"fungible"`
	Nonfungible []string `json:"nonfungible"`
}

type NonfungibleTokenMaintainerSetting struct {
	Module      string              `json:"module"`
	Maintainers []MaintainerSetting `json:"maintainers"`
}

type MaintainerSetting struct {
	Type    string                `json:"type"`
	Address []sdkTypes.AccAddress `json:"address"`
}

func queryGetNonfungibleTokenMaintainerAddresses(cdc *codec.Codec, ctx sdkTypes.Context, path []string, req abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {

	var nonfungibleTokenMaintainerSettings []MaintainerSetting

	nonfungibleTokenIssuerAddresses := keeper.GetIssuerAddresses(ctx)
	maintainerData := MaintainerSetting{"issuer_addresses", nonfungibleTokenIssuerAddresses}
	nonfungibleTokenMaintainerSettings = append(nonfungibleTokenMaintainerSettings, maintainerData)

	nonfungibleTokenProviderAddresses := keeper.GetProviderAddresses(ctx)
	maintainerData = MaintainerSetting{"provider_addresses", nonfungibleTokenProviderAddresses}
	nonfungibleTokenMaintainerSettings = append(nonfungibleTokenMaintainerSettings, maintainerData)

	nonfungibleTokenAuthorisedAddresses := keeper.GetAuthorisedAddresses(ctx)
	maintainerData = MaintainerSetting{"authorised_addresses", nonfungibleTokenAuthorisedAddresses}
	nonfungibleTokenMaintainerSettings = append(nonfungibleTokenMaintainerSettings, maintainerData)

	nonfungibleTokenMaintainerSetting := NonfungibleTokenMaintainerSetting{"nonfungible-token", nonfungibleTokenMaintainerSettings}
	respData := cdc.MustMarshalJSON(nonfungibleTokenMaintainerSetting)
	return respData, nil
}
