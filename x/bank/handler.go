package bank

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkBank "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/maxonrow/maxonrow-go/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k sdkBank.Keeper, accountKeeper auth.AccountKeeper) sdkTypes.Handler {
	return func(ctx sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		switch msg := msg.(type) {
		case MsgMxwSend:
			return handleMsgSend(ctx, k, msg, accountKeeper)
		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdkTypes.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdkTypes.Context, k sdkBank.Keeper, msg MsgMxwSend, accountKeeper auth.AccountKeeper) sdkTypes.Result {
	if !k.GetSendEnabled(ctx) {
		return sdkBank.ErrSendDisabled(k.Codespace()).Result()
	}
	err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	return MakeBankSendEvent(ctx, msg.FromAddress, msg.ToAddress, msg.Amount, accountKeeper)

}

func MakeBankSendEvent(ctx sdkTypes.Context, fromAddress sdkTypes.AccAddress, toAddress sdkTypes.AccAddress, amount sdkTypes.Coins, accountKeeper auth.AccountKeeper) sdkTypes.Result {

	ownerWalletAccount := accountKeeper.GetAccount(ctx, fromAddress)
	accountSequence := ownerWalletAccount.GetSequence()
	var log string
	if accountSequence == 0 {
		log = types.MakeResultLog(accountSequence, ctx.TxBytes())
	} else {
		log = types.MakeResultLog(accountSequence-1, ctx.TxBytes())
	}

	event :=
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute("From", fromAddress.String()),
			sdkTypes.NewAttribute("To", toAddress.String()),
			sdkTypes.NewAttribute("Amount", amount.AmountOf(types.CIN).String()),
		)

	eventParam := []string{fromAddress.String(), toAddress.String(), amount.AmountOf(types.CIN).String()}
	eventSignature := "Transferred(string,string,bignumber)"

	return sdkTypes.Result{
		Events: sdkTypes.Events.AppendEvents(sdkTypes.Events{event}, types.MakeMxwEvents(eventSignature, fromAddress.String(), eventParam)),
		Log:    log,
	}
}
