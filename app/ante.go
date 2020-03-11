package app

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/utils"
	fungible "github.com/maxonrow/maxonrow-go/x/token/fungible"
	nonFungible "github.com/maxonrow/maxonrow-go/x/token/nonfungible"
	rpc "github.com/tendermint/tendermint/rpc/core"
)

func (app *mxwApp) NewAnteHandler() sdkTypes.AnteHandler {
	return func(
		ctx sdkTypes.Context, tx sdkTypes.Tx, simulate bool,
	) (sdkTypes.Context, error) {

		// HACK:
		// When the create-non-empty-block is set to false, and the node restarts;
		// signature verification fails, because chain-id doesn't set properly, it's empty.
		// Two ways to set the chain-id:
		// 1- Through genesis block
		// 2- Creating a new block
		// Set chain-id for verify the signature when non-empty blocks in false
		// chain := app.deliverState.ctx.ChainID()

		var chainID = ctx.ChainID()
		if chainID == "" {
			chainID, blockHeight := app.retrieveChainID()

			ctx = ctx.WithChainID(chainID)
			ctx = ctx.WithBlockHeight(blockHeight)
		}

		params := app.accountKeeper.GetParams(ctx)
		stdTx, ok := tx.(sdkAuth.StdTx)
		if !ok {
			return ctx, sdkTypes.ErrInternal("tx must be StdTx")
		}

		if len(stdTx.Fee.Amount) != 1 {
			return ctx, sdkTypes.ErrInternal(fmt.Sprintf("Fee is missed. Make sure you have entered the fee amount in cin"))
		}

		if stdTx.Fee.Amount[0].Denom != types.CIN {
			return ctx, sdkTypes.ErrInternal(fmt.Sprintf("Invalid denom for fee. Fee should pay by cin."))
		}

		if stdTx.Fee.Gas != 0 {
			return ctx, sdkTypes.ErrInternal(fmt.Sprintf("MXW transactions are gas free. tx_gas: %v", stdTx.Fee.Gas))
		}

		if err := tx.ValidateBasic(); err != nil {
			return ctx, err
		}

		if err := app.ValidateMemo(stdTx, params); err != nil {
			return ctx, err
		}

		isGenesis := ctx.BlockHeight() == 0

		if !isGenesis {
			checkFeeErr := app.CheckFee(ctx, tx)
			if checkFeeErr != nil {
				return ctx, checkFeeErr
			}
		}

		// -------------------------------------------------------------------------
		// Check signatures and updat sequence and public_keys (if it hasn't set yet)
		// It may update account's PublicKey
		acc, err := utils.CheckTxSig(ctx, stdTx, app.accountKeeper, app.kycKeeper)
		if err != nil {
			return ctx, err
		}

		// Try to delete it from pending list
		if acc.IsMultiSig() {
			multisig := acc.GetMultiSig()
			txID, _ := multisig.CheckTx(stdTx)
			isMetric := multisig.IsMetric(txID)
			if !isMetric {
				return ctx, sdkTypes.ErrUnknownRequest("Multisig Transaction is not valid.")
			}
			// after validating everything, delete the pendingTx
			isDeleted := multisig.RemoveTx(txID)
			if !isDeleted {
				return ctx, sdkTypes.ErrUnknownRequest("Delete failed.")
			}
			acc.SetMultiSig(multisig)
		}

		if err := acc.SetSequence(acc.GetSequence() + 1); err != nil {
			panic(err)
		}
		app.accountKeeper.SetAccount(ctx, acc)
		// -------------------------------------------------------------------------

		// removing this condition will cause app-hash change
		if !stdTx.Fee.Amount.IsZero() {
			err := sdkAuth.DeductFees(app.supplyKeeper, ctx, acc, stdTx.Fee.Amount)
			if err != nil {
				return ctx, err
			}

			// TODO: WHY we get account again here? Try to add some test cases for this
			acc = app.accountKeeper.GetAccount(ctx, acc.GetAddress())
		}

		for _, msg := range stdTx.GetMsgs() {

			if !ok {
				ctx.Logger().Info("This message is not MXW message.", "msg", msg.Type())
			} else {
				validateMsgErr := app.validateMsg(ctx, msg)
				if validateMsgErr != nil {
					return ctx, validateMsgErr
				}
			}

			// Create fungible token, pay application fee
			createFungibleTokenMsg, ok := msg.(fungible.MsgCreateFungibleToken)
			if ok {
				amt, parseErr := sdkTypes.ParseCoins(createFungibleTokenMsg.Fee.Value + types.CIN)
				if parseErr != nil {
					return ctx, sdkTypes.ErrInvalidCoins("Parse value to coins failed.")
				}

				sendCoinsErr := app.bankKeeper.SendCoins(ctx, createFungibleTokenMsg.Owner, createFungibleTokenMsg.Fee.To, amt)
				if sendCoinsErr != nil {
					return ctx, sendCoinsErr
				}
			}
			// Create non fungible token, pay application fee
			createNonFungibleTokenMsg, ok := msg.(nonFungible.MsgCreateNonFungibleToken)
			if ok {
				amt, parseErr := sdkTypes.ParseCoins(createNonFungibleTokenMsg.Fee.Value + types.CIN)
				if parseErr != nil {
					return ctx, sdkTypes.ErrInvalidCoins("Parse value to coins failed.")
				}

				sendCoinsErr := app.bankKeeper.SendCoins(ctx, createNonFungibleTokenMsg.Owner, createNonFungibleTokenMsg.Fee.To, amt)
				if sendCoinsErr != nil {
					return ctx, sendCoinsErr
				}
			}

		}

		return ctx, nil // continue...
	}
}

func (app *mxwApp) retrieveChainID() (string, int64) {
	if app.chainID == "" {
		// calling status caused having deadlock since it try to get information from consensus engine
		//status, err := rpc.Status(nil)
		//if err != nil {
		//	panic(err)
		//}
		//num := int64(status.SyncInfo.LatestBlockHeight - 1)

		num := int64(1)
		res, err := rpc.Block(nil, &num)
		if err != nil {
			panic(err)
		}
		if res == nil {
			panic(err)
		}

		app.blockHeight = num
		app.chainID = string(res.Block.ChainID)
	}

	return app.chainID, app.blockHeight
}

func (app *mxwApp) ValidateMemo(tx sdkAuth.StdTx, params sdkAuth.Params) error {
	memo := tx.GetMemo()

	memoLength := len(memo)
	if uint64(memoLength) > params.MaxMemoCharacters {
		return sdkTypes.ErrMemoTooLarge(fmt.Sprintf(
			"maximum number of characters is %d but received %d characters",
			params.MaxMemoCharacters, memoLength))
	}

	return nil
}
