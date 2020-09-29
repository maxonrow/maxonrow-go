package kyc

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/x/fee"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryIsWhitelisted             = "is_whitelisted"
	QueryIsAuthorised              = "is_authorised"
	QueryGetKycAddress             = "get_kyc_address"
	QueryGetFee                    = "get_fee"
	QueryGetKycMaintainerAddresses = "get_kyc_maintainer_addresses"
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
		case QueryGetKycMaintainerAddresses:
			return queryGetKycMaintainerAddresses(ctx, path[1:], req, keeper)
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

type KycMaintainerSetting struct {
	Module      string              `json:"module"`
	Maintainers []MaintainerSetting `json:"maintainers"`
}

type MaintainerSetting struct {
	Type    string                `json:"type"`
	Address []sdkTypes.AccAddress `json:"address"`
}

func queryGetKycMaintainerAddresses(ctx sdkTypes.Context, path []string, req abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {

	var kycMaintainerSettings []MaintainerSetting

	kycIssuerAddresses := keeper.GetIssuerAddresses(ctx)
	maintainerData := MaintainerSetting{"issuer_addresses", kycIssuerAddresses}
	kycMaintainerSettings = append(kycMaintainerSettings, maintainerData)

	kycProviderAddresses := keeper.GetProviderAddresses(ctx)
	maintainerData = MaintainerSetting{"provider_addresses", kycProviderAddresses}
	kycMaintainerSettings = append(kycMaintainerSettings, maintainerData)

	kycAuthorisedAddresses := keeper.GetAuthorisedAddresses(ctx)
	maintainerData = MaintainerSetting{"authorised_addresses", kycAuthorisedAddresses}
	kycMaintainerSettings = append(kycMaintainerSettings, maintainerData)

	kycMaintainerSetting := KycMaintainerSetting{"kyc", kycMaintainerSettings}
	respData := codec.Cdc.MustMarshalJSON(kycMaintainerSetting)
	return respData, nil
}
