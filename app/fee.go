package app

import (
	"fmt"

	"github.com/maxonrow/maxonrow-go/x/bank"
	"github.com/maxonrow/maxonrow-go/x/fee"
	token "github.com/maxonrow/maxonrow-go/x/token/fungible"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/maxonrow/maxonrow-go/types"
)

func (app *mxwApp) CheckFee(ctx sdkTypes.Context, tx sdkTypes.Tx) sdkTypes.Error {
	stdTx, _ := tx.(sdkAuth.StdTx)
	fee, err := app.CalculateFee(ctx, tx)
	if err != nil {
		return err
	}
	if stdTx.Fee.Amount.IsAllLT(fee) {
		return sdkTypes.ErrInsufficientFee(fmt.Sprintf("Insufficient fee amount, need: %s", fee))
	}
	return nil
}

func (app *mxwApp) CalculateFee(ctx sdkTypes.Context, tx sdkTypes.Tx) (sdkTypes.Coins, sdkTypes.Error) {

	msgs := tx.GetMsgs()

	var fees sdkTypes.Coins
	for _, msg := range msgs {
		signer := msg.GetSigners()[0]

		var amt sdkTypes.Coins
		var multiplier string
		var multiplierErr sdkTypes.Error

		// 1- try to get fee-setting by account
		feeSetting, _ := app.feeKeeper.GetAccFeeSetting(ctx, signer)
		if feeSetting == nil {

			var isCustomAction = func(msg string) bool {
				return msg == token.MsgTypeTransferFungibleToken ||
					msg == token.MsgTypeMintFungibleToken ||
					msg == token.MsgTypeBurnFungibleToken ||
					msg == token.MsgTypeTransferFungibleTokenOwnership ||
					msg == token.MsgTypeAcceptFungibleTokenOwnership
			}
			r := msg.Route()
			t := msg.Type()
			if r == token.MsgRoute &&
				isCustomAction(t) {

				tokenFeeSetting, tokenAmt, feeSettingErr := app.getTokenFeeSetting(msg, ctx)
				if feeSettingErr != nil {
					return nil, feeSettingErr
				}

				multiplier, multiplierErr = app.feeKeeper.GetTokenFeeMultiplier(ctx)
				if multiplierErr != nil {
					return nil, sdkTypes.ErrInternal("Get fee multiplier failed.")
				}

				amt = tokenAmt
				feeSetting = tokenFeeSetting

			} else {

				// 2- try to get fee-setting by msg-type
				feeSetting, _ = app.feeKeeper.GetMsgFeeSetting(ctx, msg.Route()+"-"+msg.Type())
				multiplier, multiplierErr = app.feeKeeper.GetFeeMultiplier(ctx)
				if multiplierErr != nil {
					return nil, sdkTypes.ErrInternal("Get fee multiplier failed.")
				}

				bankMsg, ok := msg.(bank.MsgMxwSend)
				if ok {
					amt = bankMsg.Amount
				}
			}
		}

		fee, _ := calculateFee(ctx, feeSetting, multiplier, amt)
		fees = fees.Add(fee)
	}

	return fees, nil
}

func calculateFee(ctx sdkTypes.Context, feeSetting *fee.FeeSetting, mul string, amt sdkTypes.Coins) (sdkTypes.Coins, sdkTypes.Error) {

	if feeSetting == nil {
		panic("Fee setting should not be empty.")
	}

	amount := amt.AmountOf(types.CIN)
	if amount.IsZero() {
		return feeSetting.Min, nil
	}
	minFee := feeSetting.Min.AmountOf(types.CIN)
	maxFee := feeSetting.Max.AmountOf(types.CIN)
	percentage := sdkTypes.MustNewDecFromStr(feeSetting.Percentage)
	multiplier := sdkTypes.MustNewDecFromStr(mul)

	feeD := amount.ToDec().Mul(percentage).Mul(multiplier)
	feeD = feeD.Quo(sdkTypes.MustNewDecFromStr("100.0"))

	fee := feeD.RoundInt()
	if fee.LT(minFee) {
		fee = minFee
	}
	if fee.GT(maxFee) {
		fee = maxFee
	}

	return sdkTypes.Coins{sdkTypes.NewCoin(types.CIN, fee)}, nil
}

func (app *mxwApp) getTokenFeeSetting(msg sdkTypes.Msg, ctx sdkTypes.Context) (*fee.FeeSetting, sdkTypes.Coins, sdkTypes.Error) {

	var amt sdkTypes.Coins
	var feeSetting *fee.FeeSetting
	var feeSettingErr sdkTypes.Error

	switch msgType := msg.(type) {
	case token.MsgTransferFungibleToken:

		feeSetting, feeSettingErr = app.feeKeeper.GetTokenFeeSetting(ctx, msgType.Symbol, fee.TransferFungibleToken)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}

		transferAmt := msgType.Value.String() + types.CIN

		transferAmtCoins, parseErr := sdkTypes.ParseCoins(transferAmt)
		if parseErr != nil {
			return nil, nil, sdkTypes.ErrUnknownRequest("Parsing value failed.")
		}

		amt = transferAmtCoins

	case token.MsgMintFungibleToken:
		feeSetting, feeSettingErr = app.feeKeeper.GetTokenFeeSetting(ctx, msgType.Symbol, fee.MintFungibleToken)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}

		mintAmt := msgType.Value.String() + types.CIN
		mintAmtCoins, parseMintAmtErr := sdkTypes.ParseCoins(mintAmt)
		if parseMintAmtErr != nil {
			return nil, nil, sdkTypes.ErrUnknownRequest("Parsing value failed.")
		}

		amt = mintAmtCoins

	case token.MsgBurnFungibleToken:
		feeSetting, feeSettingErr = app.feeKeeper.GetTokenFeeSetting(ctx, msgType.Symbol, fee.BurnFungibleToken)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}

		burnAmt := msgType.Value.String() + types.CIN
		burnAmtCoins, burnAmtCoinsErr := sdkTypes.ParseCoins(burnAmt)
		if burnAmtCoinsErr != nil {
			return nil, nil, sdkTypes.ErrUnknownRequest("Parsing value failed.")
		}

		amt = burnAmtCoins

	case token.MsgTransferFungibleTokenOwnership:
		feeSetting, feeSettingErr = app.feeKeeper.GetTokenFeeSetting(ctx, msgType.Symbol, fee.TransferTokenOwnership)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}

	case token.MsgAcceptFungibleTokenOwnership:
		feeSetting, feeSettingErr = app.feeKeeper.GetTokenFeeSetting(ctx, msgType.Symbol, fee.AcceptTokenOwnership)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}
	}

	return feeSetting, amt, nil
}
