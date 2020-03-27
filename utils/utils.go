package utils

import (
	"bytes"
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/maxonrow/maxonrow-go/x/kyc"
	"github.com/tendermint/tendermint/crypto"
)

func GetSignBytes(ctx sdkTypes.Context, tx sdkAuth.StdTx, acc exported.Account) []byte {
	genesis := ctx.BlockHeight() == 0
	chainID := ctx.ChainID()
	var accNum uint64
	if !genesis {
		accNum = acc.GetAccountNumber()
	}

	if acc.IsMultiSig() {
		// 1- check if tx exists in pending tx, then tx_id is same as pending.id
		// 2. if not, seq = acc.GetCounter()
		multisig := acc.GetMultiSig()
		txID, exist := multisig.CheckTx(tx)
		if exist {
			return sdkAuth.StdSignBytes(
				chainID, accNum, txID, tx.Fee, tx.Msgs, tx.Memo,
			)
		} else {
			ctx.Logger().Error("Unable to find transaction from pending list", "address", acc.GetAddress(), "Tx", tx)
			return sdkAuth.StdSignBytes(
				chainID, accNum, multisig.GetCounter(), tx.Fee, tx.Msgs, tx.Memo,
			)
		}
	} else {
		seq := acc.GetSequence()
		return sdkAuth.StdSignBytes(
			chainID, accNum, seq, tx.Fee, tx.Msgs, tx.Memo,
		)
	}
}

func CheckTxSig(ctx sdkTypes.Context, tx sdkAuth.StdTx, accountKeeper sdkAuth.AccountKeeper, kycKeeper kyc.Keeper) (exported.Account, error) {
	params := accountKeeper.GetParams(ctx)
	msgs := tx.GetMsgs()
	if len(msgs) != 1 {
		return nil, sdkTypes.ErrInternal(fmt.Sprintf("MXW transactions accept only one message per transaction. it has %v messages", len(msgs)))
	}

	signers := msgs[0].GetSigners()
	if len(signers) != 1 {
		return nil, sdkTypes.ErrInternal(fmt.Sprintf("MXW transactions accept only one signature per message. it has %v messages", len(signers)))
	}
	signer := signers[0]

	if !kycKeeper.IsWhitelisted(ctx, signer) {
		return nil, sdkTypes.ErrInternal(fmt.Sprintf("Message signer is not whitelisted: %v", signer))
	}

	signerAcc := GetAccount(ctx, accountKeeper, signer)
	stdSigs := tx.Signatures
	if len(stdSigs) == 0 {
		return nil, sdkTypes.ErrInternal(fmt.Sprintf("No signature found. You should sign the transaction"))
	}

	var pubKeys []crypto.PubKey
	if signerAcc.IsMultiSig() {
		if len(stdSigs) > int(params.TxSigLimit) {
			return nil, sdkTypes.ErrInternal(fmt.Sprintf("Maximum signatures should be %v. It has %v signatures.", params.TxSigLimit, len(stdSigs)))
		}

		multisig := signerAcc.GetMultiSig()
		signers := multisig.GetSigners()
		for _, signer := range signers {
			acc := accountKeeper.GetAccount(ctx, signer)
			pubKeys = append(pubKeys, acc.GetPubKey())
		}
	} else {
		if len(stdSigs) != 1 {
			return nil, sdkTypes.ErrInternal(fmt.Sprintf("One signature per transaction for normal transactions. It has %v signatures.", len(stdSigs)))
		}

		sig := stdSigs[0]
		pubKey := signerAcc.GetPubKey()
		if pubKey == nil {
			pubKey = sig.PubKey
			if pubKey == nil {
				return nil, sdkTypes.ErrInvalidPubKey("PubKey not found")
			}

			if !bytes.Equal(pubKey.Address(), signerAcc.GetAddress()) {
				return nil, sdkTypes.ErrInvalidPubKey(
					fmt.Sprintf("PubKey does not match Signer address %s", signerAcc.GetAddress()))
			}

			// Set public key for the first time
			signerAcc.SetPubKey(pubKey)
		}

		pubKeys = append(pubKeys, pubKey)
	}

	for _, stdSig := range stdSigs {
		// signerAcc is groupAccount
		signBytes := GetSignBytes(ctx, tx, signerAcc)
		matched := false
		for i, pubKey := range pubKeys {
			if pubKey == nil {
				continue
			}

			if pubKey.VerifyBytes(signBytes, stdSig.Signature) {
				signedBy, err := sdkTypes.AccAddressFromHex(pubKey.Address().String())
				if err != nil {
					return nil, err
				}
				if !signerAcc.IsSigner(signedBy) {
					return nil, sdkTypes.ErrUnauthorized("Unauthorized signer: %v")
				}

				if !kycKeeper.IsWhitelisted(ctx, signedBy) {
					return nil, sdkTypes.ErrInternal(fmt.Sprintf("Transaction signer is not whitelisted: %v", signer))
				}

				//
				pubKeys[i] = nil

				matched = true
				break
			}
		}

		if !matched {
			return nil, sdkTypes.ErrUnauthorized("signature verification failed; verify correct account sequence and chain-id" + string(signBytes))
		}
	}

	return signerAcc, nil
}

func GetAccount(ctx sdkTypes.Context, keeper sdkAuth.AccountKeeper, address sdkTypes.AccAddress) exported.Account {
	acc := keeper.GetAccount(ctx, address)
	if acc == nil {
		ctx.Logger().Error("Invalid or non-exist address", "address", address)
	}

	return acc
}
