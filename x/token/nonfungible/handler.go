package nonfungible

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	ApproveToken  = "APPROVE"
	RejectToken   = "REJECT"
	FreezeToken   = "FREEZE"
	UnfreezeToken = "UNFREEZE"
	FreezeItem    = "FREEZE_ITEM"
	UnfreezeItem  = "UNFREEZE_ITEM"

	ApproveTransferTokenOwnership = "APPROVE_TRANFER_TOKEN_OWNERSHIP"
	RejectTransferTokenOwnership  = "REJECT_TRANFER_TOKEN_OWNERSHIP"
)

func NewHandler(keeper *Keeper) sdkTypes.Handler {
	return func(ctx sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		switch msg := msg.(type) {
		case MsgCreateNonFungibleToken:
			return handleMsgCreateNonFungibleToken(ctx, keeper, msg)
		case MsgSetNonFungibleTokenStatus:
			return handleMsgSetNonFungibleTokenStatus(ctx, keeper, msg)
		case MsgMintNonFungibleToken:
			return handleMsgMintNonFungibleToken(ctx, keeper, msg)
		case MsgTransferNonFungibleToken:
			return handleMsgTransferNonFungibleToken(ctx, keeper, msg)
		case MsgBurnNonFungibleToken:
			return handleMsgBurnNonFungibleToken(ctx, keeper, msg)
		case MsgSetNonFungibleItemStatus:
			return handleMsgSetNonFungibleItemStatus(ctx, keeper, msg)
		case MsgTransferNonFungibleTokenOwnership:
			return handleMsgTransferNonFungibleTokenOwnership(ctx, keeper, msg)
		case MsgAcceptNonFungibleTokenOwnership:
			return handleMsgAcceptTokenOwnership(ctx, keeper, msg)
		case MsgEndorsement:
			return handleMsgEndorsement(ctx, keeper, msg)
		case MsgUpdateItemMetadata:
			return handleMsgUpdateItemMetadata(ctx, keeper, msg)
		case MsgUpdateNFTMetadata:
			return handleMsgUpdateNFTMetadata(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized fungible token Msg type: %v", msg.Type())
			return sdkTypes.ErrUnknownRequest(errMsg).Result()
		}
	}

}

func handleMsgCreateNonFungibleToken(ctx sdkTypes.Context, keeper *Keeper, msg MsgCreateNonFungibleToken) sdkTypes.Result {

	return keeper.CreateNonFungibleToken(ctx, msg.Name, msg.Symbol, msg.Owner, msg.Properties, msg.Metadata, msg.Fee)
}

func handleMsgSetNonFungibleTokenStatus(ctx sdkTypes.Context, keeper *Keeper, msg MsgSetNonFungibleTokenStatus) sdkTypes.Result {

	signaturesErr := keeper.ValidateSignatures(ctx, msg)
	if signaturesErr != nil {
		return signaturesErr.Result()
	}

	//* token.metadata temporaily not in use.
	switch msg.Payload.Token.Status {
	case ApproveToken:
		return keeper.ApproveToken(ctx, msg.Payload.Token.Symbol, msg.Payload.Token.TokenFees, msg.Payload.Token.MintLimit, msg.Payload.Token.TransferLimit, msg.Payload.Token.EndorserList, msg.Owner, msg.Payload.Token.Burnable, msg.Payload.Token.Transferable, msg.Payload.Token.Modifiable, msg.Payload.Token.Public)
	case RejectToken:
		return keeper.RejectToken(ctx, msg.Payload.Token.Symbol, msg.Owner)
	case FreezeToken:
		return keeper.FreezeToken(ctx, msg.Payload.Token.Symbol, msg.Owner)
	case UnfreezeToken:
		return keeper.UnfreezeToken(ctx, msg.Payload.Token.Symbol, msg.Owner)
	case ApproveTransferTokenOwnership:
		return keeper.ApproveTransferTokenOwnership(ctx, msg.Payload.Token.Symbol, msg.Owner)
	case RejectTransferTokenOwnership:
		return keeper.RejectTransferTokenOwnership(ctx, msg.Payload.Token.Symbol, msg.Owner)
	default:
		errMsg := fmt.Sprintf("Unrecognized status: %v", msg.Payload.Token.Status)
		return sdkTypes.ErrUnknownRequest(errMsg).Result()
	}

}

func handleMsgMintNonFungibleToken(ctx sdkTypes.Context, keeper *Keeper, msg MsgMintNonFungibleToken) sdkTypes.Result {
	return keeper.MintNonFungibleToken(ctx, msg.Symbol, msg.Owner, msg.To, msg.ItemID, msg.Properties, msg.Metadata)
}

func handleMsgTransferNonFungibleToken(ctx sdkTypes.Context, keeper *Keeper, msg MsgTransferNonFungibleToken) sdkTypes.Result {
	return keeper.TransferNonFungibleToken(ctx, msg.Symbol, msg.From, msg.To, msg.ItemID)
}

func handleMsgBurnNonFungibleToken(ctx sdkTypes.Context, keeper *Keeper, msg MsgBurnNonFungibleToken) sdkTypes.Result {
	return keeper.BurnNonFungibleToken(ctx, msg.Symbol, msg.From, msg.ItemID)
}

func handleMsgTransferNonFungibleTokenOwnership(ctx sdkTypes.Context, keeper *Keeper, msg MsgTransferNonFungibleTokenOwnership) sdkTypes.Result {
	return keeper.TransferTokenOwnership(ctx, msg.Symbol, msg.From, msg.To)
}

func handleMsgAcceptTokenOwnership(ctx sdkTypes.Context, keeper *Keeper, msg MsgAcceptNonFungibleTokenOwnership) sdkTypes.Result {
	return keeper.AcceptTokenOwnership(ctx, msg.Symbol, msg.From)
}

func handleMsgSetNonFungibleItemStatus(ctx sdkTypes.Context, keeper *Keeper, msg MsgSetNonFungibleItemStatus) sdkTypes.Result {

	signaturesErr := keeper.ValidateSignatures(ctx, msg)
	if signaturesErr != nil {
		return signaturesErr.Result()
	}

	//* token.metadata temporaily not in use.
	switch msg.ItemPayload.Item.Status {
	case FreezeItem:
		return keeper.FreezeNonFungibleItem(ctx, msg.ItemPayload.Item.Symbol, msg.Owner, msg.ItemPayload.Item.From, msg.ItemPayload.Item.ItemID, "")
	case UnfreezeItem:
		return keeper.UnfreezeNonFungibleItem(ctx, msg.ItemPayload.Item.Symbol, msg.Owner, msg.ItemPayload.Item.ItemID, "")
	default:
		errMsg := fmt.Sprintf("Unrecognized status: %v", msg.ItemPayload.Item.Status)
		return sdkTypes.ErrUnknownRequest(errMsg).Result()
	}

}

func handleMsgEndorsement(ctx sdkTypes.Context, keeper *Keeper, msg MsgEndorsement) sdkTypes.Result {
	// check token endorser list
	if !keeper.IsTokenEndorser(ctx, msg.Symbol, msg.From) {
		return sdkTypes.ErrInternal("Invalid endorser.").Result()
	}

	return keeper.MakeEndorsement(ctx, msg.Symbol, msg.From, msg.ItemID)
}

func handleMsgUpdateItemMetadata(ctx sdkTypes.Context, keeper *Keeper, msg MsgUpdateItemMetadata) sdkTypes.Result {
	return keeper.UpdateItemMetadata(ctx, msg.Symbol, msg.From, msg.ItemID, msg.Metadata)
}

func handleMsgUpdateNFTMetadata(ctx sdkTypes.Context, keeper *Keeper, msg MsgUpdateNFTMetadata) sdkTypes.Result {
	return keeper.UpdateNFTMetadata(ctx, msg.Symbol, msg.From, msg.Metadata)
}
