package app

import (
	"bytes"
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/maxonrow/maxonrow-go/types"
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

		params := app.accountKeeper.GetParams(ctx)

		signer := stdTx.GetMsgs()[0].GetSigners()[0]
		signerAcc := app.accountKeeper.GetAccount(ctx, signer)
		if signerAcc.GetMultiSig() == nil {
			if err := tx.ValidateBasic(); err != nil {
				return ctx, err
			}

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
		signerAddrs := stdTx.GetSigners()
		if len(signerAddrs) != 1 {
			return ctx, sdkTypes.ErrInternal(fmt.Sprintf("MXW transactions accept only one signer. it has %v senders", len(signerAddrs)))
		}

		stdSigs := stdTx.Signatures
		if signerAcc.GetMultiSig() == nil {
			if len(stdSigs) != 1 {
				return ctx, sdkTypes.ErrInternal(fmt.Sprintf("MXW transactions accept only one signature. it has %v signatures", len(stdSigs)))
			}
		}

		signerAcc, err := sdkAuth.GetSignerAcc(ctx, app.accountKeeper, signerAddrs[0])
		if err != nil {
			return ctx, err
		}

		// removing this condition will cause app-hash change
		if !stdTx.Fee.Amount.IsZero() {
			err = sdkAuth.DeductFees(app.supplyKeeper, ctx, signerAcc, stdTx.Fee.Amount)
			if err != nil {
				return ctx, err
			}

			signerAcc = app.accountKeeper.GetAccount(ctx, signerAcc.GetAddress())
		}

		if signerAcc.IsMultiSig() {
			signerMultiSig := signerAcc.GetMultiSig()
			txID, validate := signerMultiSig.ValidateMultiSigTx(stdTx)
			if !validate {
				return ctx, sdkTypes.ErrUnknownAddress("Invalid multisig tx.")
			}

			for _, v := range stdSigs {
				accAddress, err := sdkTypes.AccAddressFromHex(string(v.PubKey.Address()))
				if err != nil {
					return ctx, err
				}
				if !signerMultiSig.IsSigner(accAddress) {
					return ctx, sdkTypes.ErrUnauthorized("Invalid multisig account signer.")
				}
				signer := app.accountKeeper.GetAccount(ctx, accAddress)
				signBytes := types.GetSignBytes(ctx, stdTx, signer)
				_, err = processSig(ctx, signer, v, signBytes, simulate, params)
				if err != nil {
					return ctx, err
				}
			}
			// after validating everything, delete the pendingTx
			isDeleted := signerMultiSig.RemoveTx(txID)
			if !isDeleted {
				return ctx, sdkTypes.ErrUnknownRequest("Delete failed.")
			}
			signerAcc.SetMultiSig(signerMultiSig)
			app.accountKeeper.SetAccount(ctx, signerAcc)
		} else {

			signBytes := types.GetSignBytes(ctx, stdTx, signerAcc)

			stdSig := stdSigs[0]
			signerAcc, err = processSig(ctx, signerAcc, stdSig, signBytes, simulate, params)
			if err != nil {
				return ctx, err
			}

			app.accountKeeper.SetAccount(ctx, signerAcc)
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

// verify the signature and increment the sequence. If the account doesn't have
// a pubkey, set it.
func processSig(
	ctx sdkTypes.Context, acc exported.Account, sig sdkAuth.StdSignature, signBytes []byte, _ bool, params sdkAuth.Params,
) (updatedAcc exported.Account, err error) {

	pubKey := acc.GetPubKey()
	if pubKey == nil {
		pubKey = sig.PubKey
		if pubKey == nil {
			return nil, sdkTypes.ErrInvalidPubKey("PubKey not found")
		}

		if !bytes.Equal(pubKey.Address(), acc.GetAddress()) {
			return nil, sdkTypes.ErrInvalidPubKey(
				fmt.Sprintf("PubKey does not match Signer address %s", acc.GetAddress()))
		}

		err = acc.SetPubKey(pubKey)
		if err != nil {
			return nil, sdkTypes.ErrInternal("setting PubKey on signer's account")
		}
	}

	if !pubKey.VerifyBytes(signBytes, sig.Signature) {
		return nil, sdkTypes.ErrUnauthorized("signature verification failed; verify correct account sequence and chain-id" + string(signBytes))
	}

	if acc.IsMultiSig() {
		acc.GetMultiSig().IncCounter()
	} else {
		if err := acc.SetSequence(acc.GetSequence() + 1); err != nil {
			panic(err)
		}
	}

	return acc, err
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
