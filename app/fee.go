package app

import (
	"fmt"

	"github.com/maxonrow/maxonrow-go/x/bank"
	"github.com/maxonrow/maxonrow-go/x/fee"
	ft "github.com/maxonrow/maxonrow-go/x/token/fungible"
	nft "github.com/maxonrow/maxonrow-go/x/token/nonfungible"

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
		var feeSetting *fee.FeeSetting

		var isCustomAction = func(msg string) bool {
			return msg == ft.MsgTypeTransferFungibleToken ||
				msg == ft.MsgTypeMintFungibleToken ||
				msg == ft.MsgTypeBurnFungibleToken ||
				msg == ft.MsgTypeTransferFungibleTokenOwnership ||
				msg == ft.MsgTypeAcceptFungibleTokenOwnership ||
				// nft's
				msg == nft.MsgTypeTransferNonFungibleItem ||
				msg == nft.MsgTypeMintNonFungibleItem ||
				msg == nft.MsgTypeBurnNonFungibleItem ||
				msg == nft.MsgTypeTransferNonFungibleTokenOwnership ||
				msg == nft.MsgTypeAcceptNonFungibleTokenOwnership ||
				msg == nft.MsgTypeEndorsement ||
				msg == nft.MsgTypeUpdateEndorserList ||
				msg == nft.MsgTypeUpdateItemMetadata
		}
		r := msg.Route()
		t := msg.Type()
		if r == ft.MsgRoute &&
			isCustomAction(t) {

			tokenFeeSetting, tokenAmt, feeSettingErr := app.getTokenFeeSetting(msg, ctx)
			if feeSettingErr != nil {
				return nil, feeSettingErr
			}

			multiplier, multiplierErr = app.feeKeeper.GetFungibleTokenFeeMultiplier(ctx)
			if multiplierErr != nil {
				return nil, sdkTypes.ErrInternal("Get fee multiplier failed.")
			}

			amt = tokenAmt
			feeSetting = tokenFeeSetting

		} else if r == nft.MsgRoute &&
			isCustomAction(t) {
			tokenFeeSetting, tokenAmt, feeSettingErr := app.getTokenFeeSetting(msg, ctx)
			if feeSettingErr != nil {
				return nil, feeSettingErr
			}

			multiplier, multiplierErr = app.feeKeeper.GetNonFungibleTokenFeeMultiplier(ctx)
			if multiplierErr != nil {
				return nil, sdkTypes.ErrInternal("Get fee multiplier failed.")
			}

			amt = tokenAmt
			feeSetting = tokenFeeSetting

		} else {

			// try to get fee-setting by msg-type
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

		// try to get fee-setting by account.
		// if account have fee setting overwrite it.
		accFeeSetting, _ := app.feeKeeper.GetAccFeeSetting(ctx, signer)
		if accFeeSetting != nil {
			feeSetting = accFeeSetting
		}

		fee, err := calculateFee(ctx, feeSetting, multiplier, amt)
		if err != nil {
			return nil, err
		}
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

	feeD := amount.ToDec().Mul(percentage)
	feeD = feeD.Quo(sdkTypes.MustNewDecFromStr("100.0"))

	fee := feeD.RoundInt()
	if fee.LT(minFee) {
		fee = minFee
	}
	if fee.GT(maxFee) {
		fee = maxFee
	}

	feeD = fee.ToDec().Mul(multiplier)
	fee = feeD.RoundInt()

	return sdkTypes.Coins{sdkTypes.NewCoin(types.CIN, fee)}, nil
}

func (app *mxwApp) getTokenFeeSetting(msg sdkTypes.Msg, ctx sdkTypes.Context) (*fee.FeeSetting, sdkTypes.Coins, sdkTypes.Error) {

	var amt sdkTypes.Coins
	var feeSetting *fee.FeeSetting
	var feeSettingErr sdkTypes.Error

	switch msgType := msg.(type) {
	case ft.MsgTransferFungibleToken:
		feeSetting, feeSettingErr = app.feeKeeper.GetFungibleTokenFeeSetting(ctx, msgType.Symbol, fee.TransferToken)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}

		transferAmt := msgType.Value.String() + types.CIN

		transferAmtCoins, parseErr := sdkTypes.ParseCoins(transferAmt)
		if parseErr != nil {
			return nil, nil, sdkTypes.ErrUnknownRequest("Parsing value failed.")
		}

		amt = transferAmtCoins

	case ft.MsgMintFungibleToken:
		feeSetting, feeSettingErr = app.feeKeeper.GetFungibleTokenFeeSetting(ctx, msgType.Symbol, fee.MintToken)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}

		mintAmt := msgType.Value.String() + types.CIN
		mintAmtCoins, parseMintAmtErr := sdkTypes.ParseCoins(mintAmt)
		if parseMintAmtErr != nil {
			return nil, nil, sdkTypes.ErrUnknownRequest("Parsing value failed.")
		}

		amt = mintAmtCoins

	case ft.MsgBurnFungibleToken:
		feeSetting, feeSettingErr = app.feeKeeper.GetFungibleTokenFeeSetting(ctx, msgType.Symbol, fee.BurnToken)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}

		burnAmt := msgType.Value.String() + types.CIN
		burnAmtCoins, burnAmtCoinsErr := sdkTypes.ParseCoins(burnAmt)
		if burnAmtCoinsErr != nil {
			return nil, nil, sdkTypes.ErrUnknownRequest("Parsing value failed.")
		}

		amt = burnAmtCoins

	case ft.MsgTransferFungibleTokenOwnership:
		feeSetting, feeSettingErr = app.feeKeeper.GetFungibleTokenFeeSetting(ctx, msgType.Symbol, fee.TransferTokenOwnership)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}

	case ft.MsgAcceptFungibleTokenOwnership:
		feeSetting, feeSettingErr = app.feeKeeper.GetFungibleTokenFeeSetting(ctx, msgType.Symbol, fee.AcceptTokenOwnership)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}

	// nft's
	case nft.MsgTransferNonFungibleItem:
		feeSetting, feeSettingErr = app.feeKeeper.GetNonFungibleTokenFeeSetting(ctx, msgType.Symbol, fee.TransferToken)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}

	case nft.MsgMintNonFungibleItem:
		feeSetting, feeSettingErr = app.feeKeeper.GetNonFungibleTokenFeeSetting(ctx, msgType.Symbol, fee.MintToken)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}

	case nft.MsgBurnNonFungibleItem:
		feeSetting, feeSettingErr = app.feeKeeper.GetNonFungibleTokenFeeSetting(ctx, msgType.Symbol, fee.BurnToken)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}

	case nft.MsgTransferNonFungibleTokenOwnership:
		feeSetting, feeSettingErr = app.feeKeeper.GetNonFungibleTokenFeeSetting(ctx, msgType.Symbol, fee.TransferTokenOwnership)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}

	case nft.MsgAcceptNonFungibleTokenOwnership:
		feeSetting, feeSettingErr = app.feeKeeper.GetNonFungibleTokenFeeSetting(ctx, msgType.Symbol, fee.AcceptTokenOwnership)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}

	case nft.MsgEndorsement:
		feeSetting, feeSettingErr = app.feeKeeper.GetNonFungibleTokenFeeSetting(ctx, msgType.Symbol, fee.Endorse)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}
	case nft.MsgUpdateEndorserList:
		feeSetting, feeSettingErr = app.feeKeeper.GetNonFungibleTokenFeeSetting(ctx, msgType.Symbol, fee.UpdateNFTEndorserList)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}
	case nft.MsgUpdateItemMetadata:
		feeSetting, feeSettingErr = app.feeKeeper.GetNonFungibleTokenFeeSetting(ctx, msgType.Symbol, fee.UpdateNFTItemMetadata)
		if feeSettingErr != nil {
			return nil, nil, feeSettingErr
		}
	}

	return feeSetting, amt, nil
}
