package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/fee"
)

const (
	QueryResolve       = "resolve"
	QueryWhois         = "whois"
	QueryGetFee        = "get_fee"
	QueryListUsedAlias = "list_used_alias"
	QueryPendingAlias  = "pending"
)

type Resolve struct {
	Alias   string         `json:"alias"`
	Address sdk.AccAddress `json:"address"`
}

func NewQuerier(cdc *codec.Codec, keeper Keeper, feeKeeper fee.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryResolve:
			return queryResolve(cdc, ctx, path[1:], req, keeper)
		case QueryWhois:
			return queryWhois(ctx, path[1:], req, keeper)
		case QueryGetFee:
			return queryGetFee(cdc, ctx, path[1:], req, keeper, feeKeeper)
		case QueryListUsedAlias:
			return queryListUsedAlias(cdc, ctx, path[1:], req, keeper)
		case QueryPendingAlias:
			return queryPendingAlias(cdc, ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

func queryListUsedAlias(cdc *codec.Codec, ctx sdkTypes.Context, _ []string, _ abci.RequestQuery, keeper Keeper) ([]byte, sdkTypes.Error) {
	usedAlias := keeper.ListUsedAlias(ctx)

	resp := &listAliasResponse{
		UsedAlias: usedAlias,
	}
	respData := cdc.MustMarshalJSON(resp)

	return respData, nil
}

func queryResolve(cdc *codec.Codec, ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	alias := path[0]

	value, err := keeper.ResolveAlias(ctx, alias)

	if err != nil {
		return []byte{}, err
	}

	return []byte(value), nil
}

func queryWhois(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	addressString := path[0]

	address, err := sdk.AccAddressFromBech32(addressString)

	if err != nil {
		return nil, sdk.ErrInvalidAddress(addressString)
	}

	value := keeper.Whois(ctx, address)

	if value == "" {
		return []byte{}, types.ErrAliasCouldNotResolveAddress()
	}

	return []byte(value), nil
}

func queryGetFee(cdc *codec.Codec, ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, feeKeeper fee.Keeper) ([]byte, sdk.Error) {

	msgType := path[0]
	totalAmount, parseCoinsErr := sdk.ParseCoins("")
	if parseCoinsErr != nil {
		return nil, sdk.ErrInvalidCoins("Invalid amount")
	}

	feeSetting, feeSettingErr := feeKeeper.GetMsgFeeSetting(ctx, msgType)
	if feeSettingErr != nil {
		return nil, feeSettingErr
	}

	calculatedFee, calculatedFeeErr := fee.DefaultCalculateFee(ctx, feeSetting, totalAmount)
	if calculatedFeeErr != nil {
		return nil, calculatedFeeErr
	}

	respData := sdk.MustSortJSON(codec.Cdc.MustMarshalJSON(calculatedFee))

	return respData, nil

}

func queryPendingAlias(cdc *codec.Codec, ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {

	addressString := path[0]

	address, err := sdk.AccAddressFromBech32(addressString)

	if err != nil {
		return nil, sdk.ErrInvalidAddress(addressString)
	}

	aliasOwnerData, aliasOwnerDataErr := keeper.getPendingAliasOwnerData(ctx, address)
	if aliasOwnerDataErr != nil {
		return nil, aliasOwnerDataErr
	}

	aliasData, aliasDataErr := keeper.getPendingAlias(ctx, aliasOwnerData.Name)
	if aliasDataErr != nil {
		return nil, aliasDataErr
	}

	respData := sdk.MustSortJSON(codec.Cdc.MustMarshalJSON(aliasData))

	return respData, nil

}

type listAliasResponse struct {
	UsedAlias []string `json:"alias"`
}
