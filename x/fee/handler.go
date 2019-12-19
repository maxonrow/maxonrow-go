package fee

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper *Keeper) sdkTypes.Handler {
	return func(ctx sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		switch msg := msg.(type) {

		case MsgSysFeeSetting:
			return handleMsgFeeSetting(ctx, keeper, msg)
		case MsgAssignFeeToMsg:
			return handleMsgAssignFeeToMsg(ctx, keeper, msg)
		case MsgAssignFeeToAcc:
			return handleMsgAssignFeeToAcc(ctx, keeper, msg)
		case MsgMultiplier:
			return handleMsgMultiplier(ctx, keeper, msg)
		case MsgTokenMultiplier:
			return handleMsgTokenMultiplier(ctx, keeper, msg)
		case MsgDeleteSysFeeSetting:
			return handleMsgDeleteSysFeeSetting(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized fee Msg type: %v", msg.Type())
			return sdkTypes.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgFeeSetting(ctx sdkTypes.Context, keeper *Keeper, msg MsgSysFeeSetting) sdkTypes.Result {
	return keeper.CreateFeeSetting(ctx, msg)
}

func handleMsgAssignFeeToMsg(ctx sdkTypes.Context, keeper *Keeper, msg MsgAssignFeeToMsg) sdkTypes.Result {
	return keeper.AssignFeeToMsg(ctx, msg)
}

func handleMsgAssignFeeToAcc(ctx sdkTypes.Context, keeper *Keeper, msg MsgAssignFeeToAcc) sdkTypes.Result {
	return keeper.AssignFeeToAcc(ctx, msg)
}

func handleMsgMultiplier(ctx sdkTypes.Context, keeper *Keeper, msg MsgMultiplier) sdkTypes.Result {
	return keeper.CreateMultiplier(ctx, msg)
}

func handleMsgTokenMultiplier(ctx sdkTypes.Context, keeper *Keeper, msg MsgTokenMultiplier) sdkTypes.Result {
	return keeper.CreateTokenMultiplier(ctx, msg)
}

func handleMsgDeleteSysFeeSetting(ctx sdkTypes.Context, keeper *Keeper, msg MsgDeleteSysFeeSetting) sdkTypes.Result {
	return keeper.DeleteFeeSetting(ctx, msg)
}
