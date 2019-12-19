package kyc

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/maxonrow/maxonrow-go/x/fee"
)

const (
	QueryIsWhitelisted = "is_whitelisted"
	QueryIsAuthorised  = "is_authorised"
	QueryGetKycAddress = "get_kyc_address"
	QueryGetFee        = "get_fee"
)

func NewQuerier(keeper *Keeper, feeKeeper *fee.Keeper) sdkTypes.Querier {
	return func(ctx sdkTypes.Context, path []string, req abci.RequestQuery) ([]byte, sdkTypes.Error) {
		switch path[0] {
		case QueryIsWhitelisted:
			return queryIsWhitelisted(ctx, path[1:], req, keeper)
		case QueryIsAuthorised:
			return queryIsAuthorised(ctx, path[1:], req, keeper)
		case QueryGetKycAddress:
			return queryGetKycAddress(ctx, path[1:], req, keeper)
		case QueryGetFee:
			return queryGetFee(ctx, path[1:], req, keeper, feeKeeper)
		default:
			return nil, sdkTypes.ErrUnknownRequest("unknown kyc query endpoint")
		}
	}
}

func queryIsWhitelisted(ctx sdkTypes.Context, path []string, req abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {

	addressString := path[0]

	address, err := sdkTypes.AccAddressFromBech32(addressString)
	if err != nil {
		return nil, sdkTypes.ErrInvalidAddress(addressString)
	}

	if keeper.IsWhitelisted(ctx, address) {
		return []byte("True"), nil
	} else {
		return []byte("False"), nil
	}
}

func queryIsAuthorised(ctx sdkTypes.Context, path []string, req abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {

	addressString := path[0]

	address, err := sdkTypes.AccAddressFromBech32(addressString)
	if err != nil {
		return nil, sdkTypes.ErrInvalidAddress(addressString)
	}

	if keeper.IsAuthorised(ctx, address) {
		return []byte("True"), nil
	} else {
		return []byte("False"), nil
	}
}

func queryGetKycAddress(ctx sdkTypes.Context, path []string, req abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {

	addressString := path[0]

	address, err := sdkTypes.AccAddressFromBech32(addressString)
	if err != nil {
		return nil, sdkTypes.ErrInvalidAddress(addressString)
	}

	kycAddress := keeper.GetKycAddress(ctx, address)

	return kycAddress, nil

}

func queryGetFee(ctx sdkTypes.Context, path []string, req abci.RequestQuery, keeper *Keeper, feeKeeper *fee.Keeper) ([]byte, sdkTypes.Error) {

	msgType := path[0]
	totalAmount, parseCoinsErr := sdkTypes.ParseCoins("")
	if parseCoinsErr != nil {
		return nil, sdkTypes.ErrInvalidCoins("Invalid amount")
	}

	feeSetting, feeSettingErr := feeKeeper.GetMsgFeeSetting(ctx, msgType)
	if feeSettingErr != nil {
		return nil, feeSettingErr
	}

	calculatedFee, calculatedFeeErr := fee.DefaultCalculateFee(ctx, feeSetting, totalAmount)
	if calculatedFeeErr != nil {
		return nil, calculatedFeeErr
	}

	respData := sdkTypes.MustSortJSON(codec.Cdc.MustMarshalJSON(calculatedFee))

	return respData, nil

}
