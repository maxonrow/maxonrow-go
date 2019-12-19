package nameservice

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/types"
)

const (
	ApproveAlias = "APPROVE"
	RejectAlias  = "REJECT"
	RevokeAlias  = "REVOKE"
)

func NewHandler(keeper Keeper) sdkTypes.Handler {
	return func(ctx sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		switch msg := msg.(type) {
		case MsgCreateAlias:
			return handleMsgCreateAlias(ctx, keeper, msg)
		case MsgSetAliasStatus:
			return handleMsgSetAliasStatus(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdkTypes.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateAlias(ctx sdkTypes.Context, keeper Keeper, msg MsgCreateAlias) sdkTypes.Result {

	if keeper.IsAliasExists(ctx, msg.Name) {
		return types.ErrAliasIsInUsed().Result()
	}

	return keeper.CreateAlias(ctx, msg.Owner, msg.Name, msg.Fee)

}

func handleMsgSetAliasStatus(ctx sdkTypes.Context, keeper Keeper, msg MsgSetAliasStatus) sdkTypes.Result {

	signaturesErr := keeper.ValidateSignatures(ctx, msg)
	if signaturesErr != nil {
		return signaturesErr.Result()
	}

	//* token.metadata temporaily not in use.
	switch msg.Payload.Alias.Status {
	case ApproveAlias:
		return keeper.ApproveAlias(ctx, msg.Payload.Alias.Name, msg.Owner, "")
	case RejectAlias:
		return keeper.RejectAlias(ctx, msg.Payload.Alias.Name, msg.Owner)
	case RevokeAlias:
		return keeper.RevokeAlias(ctx, msg.Payload.Alias.Name, msg.Owner)
	default:
		errMsg := fmt.Sprintf("Unrecognized status: %v", msg.Payload.Alias.Status)
		return sdkTypes.ErrUnknownRequest(errMsg).Result()
	}
}
