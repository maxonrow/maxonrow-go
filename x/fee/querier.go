package fee

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QuerySysFeeSetting      = "get_sys_fee_setting"
	QueryMsgFeeSetting      = "get_msg_fee_setting"
	QueryAccFeeSetting      = "get_acc_fee_setting"
	QueryTokenFeeSetting    = "get_token_fee_setting"
	QueryFeeMultiplier      = "get_fee_multiplier"
	QueryTokenFeeMultiplier = "get_token_fee_multiplier"
	QueryListFeeSettings    = "list_all_fee_settings"
	QueryIsFeeSettingExist  = "is_fee_setting_exist"
	QueryIsFeeSettingInUsed = "is_fee_setting_in_used"
)

func NewQuerier(cdc *codec.Codec, keeper *Keeper) sdkTypes.Querier {
	return func(ctx sdkTypes.Context, path []string, req abci.RequestQuery) ([]byte, sdkTypes.Error) {
		fmt.Println(path)
		switch path[0] {
		case QuerySysFeeSetting:
			return querySysFeeSetting(cdc, ctx, path[1:], req, keeper)
		case QueryMsgFeeSetting:
			return queryMsgFeeSetting(cdc, ctx, path[1:], req, keeper)
		case QueryAccFeeSetting:
			return queryAccFeeSetting(cdc, ctx, path[1:], req, keeper)
		case QueryTokenFeeSetting:
			return queryTokenFeeSetting(cdc, ctx, path[1:], req, keeper)
		case QueryFeeMultiplier:
			return queryFeeMultiplier(cdc, ctx, path[1:], req, keeper)
		case QueryTokenFeeMultiplier:
			return queryTokenFeeMultiplier(cdc, ctx, path[1:], req, keeper)
		case QueryListFeeSettings:
			return queryListFeeSettings(cdc, ctx, path[1:], req, keeper)
		case QueryIsFeeSettingExist:
			return queryIsFeeSettingExist(cdc, ctx, path[1:], req, keeper)
		case QueryIsFeeSettingInUsed:
			return queryIsFeeSettingInUsed(cdc, ctx, path[1:], req, keeper)
		default:
			return nil, sdkTypes.ErrUnknownRequest("unknown fee query endpoint")
		}
	}
}

func queryMsgFeeSetting(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {

	if len(path) != 1 {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid path %s", strings.Join(path, "/")))
	}

	msgType := path[0]

	feeSetting, err := keeper.GetMsgFeeSetting(ctx, msgType)
	if err != nil {
		return nil, err
	}

	respData := cdc.MustMarshalJSON(feeSetting)

	return respData, nil
}

func querySysFeeSetting(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {
	if len(path) != 1 {
		return nil, sdkTypes.ErrUnknownRequest("fee setting type invalid.")
	}

	feeSettingType := path[0]
	feeSetting, err := keeper.GetFeeSettingByName(ctx, feeSettingType)
	if err != nil {
		return nil, err
	}

	feeSettingData := cdc.MustMarshalJSON(feeSetting)

	return feeSettingData, nil
}

func queryIsFeeSettingExist(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {
	if len(path) != 1 {
		return nil, sdkTypes.ErrUnknownRequest("fee setting invalid.")
	}

	feeSetting := path[0]
	isExists := keeper.FeeSettingExists(ctx, feeSetting)

	str := strconv.FormatBool(isExists)
	fmt.Println(str)
	return []byte(str), nil
}

func queryIsFeeSettingInUsed(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {
	if len(path) != 1 {
		return nil, sdkTypes.ErrUnknownRequest("fee setting invalid.")
	}

	feeSetting := path[0]
	isExists := keeper.IsFeeSettingUsed(ctx, feeSetting)

	str := strconv.FormatBool(isExists)
	fmt.Println(str)
	return []byte(str), nil
}

func queryAccFeeSetting(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {
	if len(path) != 1 {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid path %s", strings.Join(path, "/")))
	}

	accStr := path[0]
	acc, accErr := sdkTypes.AccAddressFromBech32(accStr)
	if accErr != nil {
		return nil, sdkTypes.ErrUnknownRequest("Invalid account address.")
	}
	feeSetting, err := keeper.GetAccFeeSetting(ctx, acc)
	if err != nil {
		return nil, err
	}

	respData := cdc.MustMarshalJSON(feeSetting)

	return respData, nil
}

func queryTokenFeeSetting(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {
	if len(path) != 2 {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid path %s", strings.Join(path, "/")))
	}

	tokenSymbol := path[0]
	tokenAction := path[1]

	feeSetting, err := keeper.GetTokenFeeSetting(ctx, tokenSymbol, tokenAction)
	if err != nil {
		return nil, err
	}

	respData := cdc.MustMarshalJSON(feeSetting)

	return respData, nil
}

func queryFeeMultiplier(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {
	if len(path) != 0 {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid path %s", strings.Join(path, "/")))
	}

	multiplier, err := keeper.GetFeeMultiplier(ctx)
	if err != nil {
		return nil, err
	}

	return []byte(multiplier), nil
}

func queryTokenFeeMultiplier(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {
	if len(path) != 0 {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid path %s", strings.Join(path, "/")))
	}

	tokenFeemultiplier, err := keeper.GetTokenFeeMultiplier(ctx)
	if err != nil {
		return nil, err
	}

	return []byte(tokenFeemultiplier), nil
}

func queryListFeeSettings(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {
	if len(path) != 0 {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid path %s", strings.Join(path, "/")))
	}

	feeSettings := keeper.ListAllSysFeeSetting(ctx)

	respData := cdc.MustMarshalJSON(feeSettings)

	return respData, nil
}
