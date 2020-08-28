package fungible

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	ApproveToken         = "APPROVE"
	RejectToken          = "REJECT"
	FreezeToken          = "FREEZE"
	UnfreezeToken        = "UNFREEZE"
	FreezeTokenAccount   = "FREEZE_ACCOUNT"
	UnfreezeTokenAccount = "UNFREEZE_ACCOUNT"

	ApproveTransferTokenOwnership = "APPROVE_TRANSFER_TOKEN_OWNERSHIP"
	RejectTransferTokenOwnership  = "REJECT_TRANSFER_TOKEN_OWNERSHIP"
)

func NewHandler(keeper *Keeper) sdkTypes.Handler {
	return func(ctx sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		switch msg := msg.(type) {
		case MsgCreateFungibleToken:
			return handleMsgCreateFungibleToken(ctx, keeper, msg)
		case MsgSetFungibleTokenStatus:
			return handleMsgSetFungibleTokenStatus(ctx, keeper, msg)
		case MsgMintFungibleToken:
			return handleMsgMintFungibleToken(ctx, keeper, msg)
		case MsgTransferFungibleToken:
			return handleMsgTransferFungibleToken(ctx, keeper, msg)
		case MsgBurnFungibleToken:
			return handleMsgBurnFungibleToken(ctx, keeper, msg)
		case MsgSetFungibleTokenAccountStatus:
			return handleMsgSetFungibleTokenAccountStatus(ctx, keeper, msg)
		case MsgTransferFungibleTokenOwnership:
			return handleMsgTransferTokenOwnership(ctx, keeper, msg)
		case MsgAcceptFungibleTokenOwnership:
			return handleMsgAcceptTokenOwnership(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized fungible Msg type: %v", msg.Type())
			return sdkTypes.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateFungibleToken(ctx sdkTypes.Context, keeper *Keeper, msg MsgCreateFungibleToken) sdkTypes.Result {

	return keeper.CreateFungibleToken(ctx, msg.Name, msg.Symbol, msg.Decimals, msg.Owner, msg.FixedSupply, msg.MaxSupply, msg.Metadata, msg.Fee)
}

func handleMsgSetFungibleTokenStatus(ctx sdkTypes.Context, keeper *Keeper, msg MsgSetFungibleTokenStatus) sdkTypes.Result {

	signaturesErr := keeper.ValidateSignatures(ctx, msg)
	if signaturesErr != nil {
		return signaturesErr.Result()
	}

	//* token.metadata temporaily not in use.
	switch msg.Payload.Token.Status {
	case ApproveToken:
		return keeper.ApproveToken(ctx, msg.Payload.Token.Symbol, msg.Payload.Token.TokenFees, msg.Payload.Token.Burnable, msg.Owner, "")
	case RejectToken:
		return keeper.RejectToken(ctx, msg.Payload.Token.Symbol, msg.Owner)
	case FreezeToken:
		return keeper.FreezeToken(ctx, msg.Payload.Token.Symbol, msg.Owner, "")
	case UnfreezeToken:
		return keeper.UnfreezeToken(ctx, msg.Payload.Token.Symbol, msg.Owner, "")
	case ApproveTransferTokenOwnership:
		return keeper.ApproveTransferTokenOwnership(ctx, msg.Payload.Token.Symbol, msg.Owner)
	case RejectTransferTokenOwnership:
		return keeper.RejectTransferTokenOwnership(ctx, msg.Payload.Token.Symbol, msg.Owner)
	default:
		errMsg := fmt.Sprintf("Unrecognized status: %v", msg.Payload.Token.Status)
		return sdkTypes.ErrUnknownRequest(errMsg).Result()
	}

}

func handleMsgMintFungibleToken(ctx sdkTypes.Context, keeper *Keeper, msg MsgMintFungibleToken) sdkTypes.Result {
	return keeper.MintFungibleToken(ctx, msg.Symbol, msg.Owner, msg.To, msg.Value)
}

func handleMsgTransferFungibleToken(ctx sdkTypes.Context, keeper *Keeper, msg MsgTransferFungibleToken) sdkTypes.Result {
	return keeper.TransferFungibleToken(ctx, msg.Symbol, msg.From, msg.To, msg.Value)
}

func handleMsgBurnFungibleToken(ctx sdkTypes.Context, keeper *Keeper, msg MsgBurnFungibleToken) sdkTypes.Result {
	return keeper.BurnFungibleToken(ctx, msg.Symbol, msg.From, msg.Value)
}

func handleMsgTransferTokenOwnership(ctx sdkTypes.Context, keeper *Keeper, msg MsgTransferFungibleTokenOwnership) sdkTypes.Result {
	return keeper.TransferTokenOwnership(ctx, msg.Symbol, msg.From, msg.To, "")
}

func handleMsgAcceptTokenOwnership(ctx sdkTypes.Context, keeper *Keeper, msg MsgAcceptFungibleTokenOwnership) sdkTypes.Result {
	return keeper.AcceptTokenOwnership(ctx, msg.Symbol, msg.From, "")
}

func handleMsgSetFungibleTokenAccountStatus(ctx sdkTypes.Context, keeper *Keeper, msg MsgSetFungibleTokenAccountStatus) sdkTypes.Result {

	signaturesErr := keeper.ValidateSignatures(ctx, msg)
	if signaturesErr != nil {
		return signaturesErr.Result()
	}

	//* token.metadata temporaily not in use.
	switch msg.TokenAccountPayload.TokenAccount.Status {
	case FreezeTokenAccount:
		return keeper.FreezeFungibleTokenAccount(ctx, msg.TokenAccountPayload.TokenAccount.Symbol, msg.Owner, msg.TokenAccountPayload.TokenAccount.Account, "")
	case UnfreezeTokenAccount:
		return keeper.UnfreezeFungibleTokenAccount(ctx, msg.TokenAccountPayload.TokenAccount.Symbol, msg.Owner, msg.TokenAccountPayload.TokenAccount.Account, "")
	default:
		errMsg := fmt.Sprintf("Unrecognized status: %v", msg.TokenAccountPayload.TokenAccount.Status)
		return sdkTypes.ErrUnknownRequest(errMsg).Result()
	}

}
