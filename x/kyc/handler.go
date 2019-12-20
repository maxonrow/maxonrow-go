package kyc

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/types"
)

func NewHandler(keeper *Keeper) sdkTypes.Handler {
	return func(ctx sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		switch msg := msg.(type) {
		case MsgWhitelist:
			return handleMsgWhitelist(ctx, keeper, msg)
		case MsgRevokeWhitelist:
			return handleMsgRevokeWhitelist(ctx, keeper, msg)
		case MsgKycBind:
			return handleMsgKycBind(ctx, keeper, msg)
		case MsgKycUnbind:
			return handleMsgKycUnbind(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized whitelist Msg type: %v", msg.Type())
			return sdkTypes.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgWhitelist(ctx sdkTypes.Context, keeper *Keeper, msg MsgWhitelist) sdkTypes.Result {

	if !keeper.IsAuthorised(ctx, msg.Owner) {
		return sdkTypes.ErrUnauthorized("Not authorized to whitelist").Result()
	}

	signaturesErr := keeper.ValidateSignatures(ctx, msg)

	if signaturesErr != nil {
		return signaturesErr.Result()
	}

	keeper.Whitelist(ctx, msg.KycData.Payload.Kyc.From, msg.KycData.Payload.Kyc.KycAddress)

	ownerWalletAccount := keeper.accountKeeper.GetAccount(ctx, msg.Owner)
	accountSequence := ownerWalletAccount.GetSequence()
	var log string
	if accountSequence == 0 {
		log = types.MakeResultLog(accountSequence, ctx.TxBytes())
	} else {
		log = types.MakeResultLog(accountSequence-1, ctx.TxBytes())
	}

	eventParam := []string{msg.KycData.Payload.Kyc.From.String(), msg.KycData.Payload.Kyc.KycAddress}
	eventSignature := "KycWhitelisted(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msg.GetSigners()[0].String(), eventParam),
		Log:    log,
	}
}

func handleMsgRevokeWhitelist(ctx sdkTypes.Context, keeper *Keeper, msg MsgRevokeWhitelist) sdkTypes.Result {

	if !keeper.IsAuthorised(ctx, msg.Owner) {
		return sdkTypes.ErrUnauthorized("Not authorized to whitelist").Result()
	}

	signaturesErr := keeper.ValidateRevokeWhitelistSignatures(ctx, msg)
	if signaturesErr != nil {
		return signaturesErr.Result()
	}

	return keeper.RevokeWhitelist(ctx, msg.RevokePayload.RevokeKycData.To, msg.Owner)

}

func handleMsgKycBind(ctx sdkTypes.Context, keeper *Keeper, msg MsgKycBind) sdkTypes.Result {

	if !keeper.IsWhitelisted(ctx, msg.From) {
		return sdkTypes.ErrUnauthorized("Singer is not whitelisted.").Result()
	}

	ownerWalletAccount := keeper.accountKeeper.GetAccount(ctx, msg.From)
	accountSequence := ownerWalletAccount.GetSequence()
	var log string
	if accountSequence == 0 {
		log = types.MakeResultLog(accountSequence, ctx.TxBytes())
	} else {
		log = types.MakeResultLog(accountSequence-1, ctx.TxBytes())
	}

	eventParam := []string{msg.From.String(), msg.To.String(), msg.KycAddress}
	eventSignature := "KycBinded(string,string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msg.GetSigners()[0].String(), eventParam),
		Log:    log,
	}
}

func handleMsgKycUnbind(ctx sdkTypes.Context, keeper *Keeper, msg MsgKycUnbind) sdkTypes.Result {

	if !keeper.IsWhitelisted(ctx, msg.From) {
		return sdkTypes.ErrUnauthorized("Singer is not whitelisted.").Result()
	}

	ownerWalletAccount := keeper.accountKeeper.GetAccount(ctx, msg.From)
	accountSequence := ownerWalletAccount.GetSequence()
	var log string
	if accountSequence == 0 {
		log = types.MakeResultLog(accountSequence, ctx.TxBytes())
	} else {
		log = types.MakeResultLog(accountSequence-1, ctx.TxBytes())
	}

	eventParam := []string{msg.From.String(), msg.To.String(), msg.KycAddress}
	eventSignature := "KycUnbinded(string,string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msg.GetSigners()[0].String(), eventParam),
		Log:    log,
	}
}
