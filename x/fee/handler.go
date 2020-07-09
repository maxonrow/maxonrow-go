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
		case MsgAssignFeeToFungibleToken:
			return handleMsgAssignFeeToFungibleToken(ctx, keeper, msg)
		case MsgAssignFeeToNonFungibleToken:
			return handleMsgAssignFeeToNonFungibleToken(ctx, keeper, msg)
		case MsgMultiplier:
			return handleMsgMultiplier(ctx, keeper, msg)
		case MsgFungibleTokenMultiplier:
			return handleMsgFungibleTokenMultiplier(ctx, keeper, msg)
		case MsgNonFungibleTokenMultiplier:
			return handleMsgNonFungibleTokenMultiplier(ctx, keeper, msg)
		case MsgDeleteSysFeeSetting:
			return handleMsgDeleteSysFeeSetting(ctx, keeper, msg)
		case MsgDeleteAccFeeSetting:
			return handleMsgDeleteAccFeeSetting(ctx, keeper, msg)
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

func handleMsgAssignFeeToFungibleToken(ctx sdkTypes.Context, keeper *Keeper, msg MsgAssignFeeToFungibleToken) sdkTypes.Result {
	return keeper.AssignFeeToFungibleToken(ctx, msg)
}

func handleMsgAssignFeeToNonFungibleToken(ctx sdkTypes.Context, keeper *Keeper, msg MsgAssignFeeToNonFungibleToken) sdkTypes.Result {
	return keeper.AssignFeeToNonFungibleToken(ctx, msg)
}

func handleMsgMultiplier(ctx sdkTypes.Context, keeper *Keeper, msg MsgMultiplier) sdkTypes.Result {
	return keeper.CreateMultiplier(ctx, msg)
}

func handleMsgFungibleTokenMultiplier(ctx sdkTypes.Context, keeper *Keeper, msg MsgFungibleTokenMultiplier) sdkTypes.Result {
	return keeper.CreateFungibleTokenMultiplier(ctx, msg)
}

func handleMsgNonFungibleTokenMultiplier(ctx sdkTypes.Context, keeper *Keeper, msg MsgNonFungibleTokenMultiplier) sdkTypes.Result {
	return keeper.CreateNonFungibleTokenMultiplier(ctx, msg)
}

func handleMsgDeleteSysFeeSetting(ctx sdkTypes.Context, keeper *Keeper, msg MsgDeleteSysFeeSetting) sdkTypes.Result {
	return keeper.DeleteFeeSetting(ctx, msg)
}

func handleMsgDeleteAccFeeSetting(ctx sdkTypes.Context, keeper *Keeper, msg MsgDeleteAccFeeSetting) sdkTypes.Result {
	return keeper.DeleteAccFeeSetting(ctx, msg)
}
