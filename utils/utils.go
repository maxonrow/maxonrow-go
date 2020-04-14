package utils

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/maxonrow/maxonrow-go/x/kyc"
	"golang.org/x/crypto/ripemd160"
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
			ctx.Logger().Info("Unable to find transaction from pending list", "address", acc.GetAddress(), "Tx", tx)
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

func CheckTxSig(ctx sdkTypes.Context, tx sdkAuth.StdTx, accountKeeper sdkAuth.AccountKeeper, kycKeeper kyc.Keeper) (exported.Account, bool, error) {
	params := accountKeeper.GetParams(ctx)
	msgs := tx.GetMsgs()
	if len(msgs) != 1 {
		return nil, false, sdkTypes.ErrInternal(fmt.Sprintf("MXW transactions accept only one message per transaction. it has %v messages", len(msgs)))
	}

	signers := msgs[0].GetSigners()
	if len(signers) != 1 {
		return nil, false, sdkTypes.ErrInternal(fmt.Sprintf("MXW transactions accept only one signature per message. it has %v messages", len(signers)))
	}
	signer := signers[0]

	if !kycKeeper.IsWhitelisted(ctx, signer) {
		return nil, false, sdkTypes.ErrInternal(fmt.Sprintf("Message signer is not whitelisted: %s", signer))
	}

	signerAcc := GetAccount(ctx, accountKeeper, signer)
	stdSigs := tx.Signatures
	if len(stdSigs) == 0 {
		return nil, false, sdkTypes.ErrInternal(fmt.Sprintf("No signature found. You should sign the transaction"))
	}

	if signerAcc.IsMultiSig() {
		if len(stdSigs) > int(params.TxSigLimit) {
			return nil, false, sdkTypes.ErrInternal(fmt.Sprintf("Maximum signatures should be %v. It has %v signatures.", params.TxSigLimit, len(stdSigs)))
		}
	} else {
		if len(stdSigs) != 1 {
			return nil, false, sdkTypes.ErrInternal(fmt.Sprintf("One signature per transaction for normal transactions. It has %v signatures.", len(stdSigs)))
		}

		sig := stdSigs[0]
		pubKey := signerAcc.GetPubKey()
		if pubKey == nil {
			pubKey = sig.PubKey
			if pubKey == nil {
				return nil, false, sdkTypes.ErrInvalidPubKey("PubKey not found")
			}

			if !bytes.Equal(pubKey.Address(), signerAcc.GetAddress()) {
				return nil, false, sdkTypes.ErrInvalidPubKey(
					fmt.Sprintf("PubKey does not match Signer address %s", signerAcc.GetAddress()))
			}

			// Set public key for the first time
			signerAcc.SetPubKey(pubKey)
		}
	}

	validSignatures := 0
	signBytes := GetSignBytes(ctx, tx, signerAcc)
	isMetric := checkSigsRecursively(ctx, accountKeeper, kycKeeper, signerAcc, stdSigs, signBytes, &validSignatures)
	if validSignatures < len(stdSigs) {
		return nil, false, sdkTypes.ErrUnauthorized("signature verification failed; verify correct account sequence and chain-id" + string(signBytes))
	}
	return signerAcc, isMetric, nil
}

func GetAccount(ctx sdkTypes.Context, keeper sdkAuth.AccountKeeper, address sdkTypes.AccAddress) exported.Account {
	acc := keeper.GetAccount(ctx, address)
	if acc == nil {
		ctx.Logger().Error("Invalid or non-exist address", "address", address)
	}

	return acc
}

func DeriveMultiSigAddress(addr sdkTypes.AccAddress, sequence uint64) sdkTypes.AccAddress {

	addrBz := addr.Bytes()
	sequenceBz := sdkTypes.Uint64ToBigEndian(sequence)
	sequenceBz = bytes.TrimLeft(sequenceBz, "\x00")
	temp := append(addrBz[:], sequenceBz[:]...)

	hasherSHA256 := sha256.New()
	hasherSHA256.Write(temp[:]) // does not error
	sha := hasherSHA256.Sum(nil)

	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha) // does not error

	return sdkTypes.AccAddress(hasherRIPEMD160.Sum(nil))
}

func MustGetAccAddressFromBech32(bech32 string) sdkTypes.AccAddress {
	addr, _ := sdkTypes.AccAddressFromBech32(bech32)
	return addr
}

func checkSigsRecursively(ctx sdkTypes.Context, accountKeeper sdkAuth.AccountKeeper, kycKeeper kyc.Keeper, acc exported.Account, stdSigs []sdkAuth.StdSignature, signBytes []byte, validSignatures *int) bool {
	if acc.IsMultiSig() {
		ms := acc.GetMultiSig()
		signers := ms.GetSigners()
		threshold := ms.GetThreshold()
		for _, signer := range signers {
			signerAcc := accountKeeper.GetAccount(ctx, signer)
			ok := checkSigsRecursively(ctx, accountKeeper, kycKeeper, signerAcc, stdSigs, signBytes, validSignatures)
			if ok {
				threshold = threshold - 1
				if threshold == 0 {
					return true
				}
			}
		}
	} else {
		pubKey := acc.GetPubKey()
		if pubKey == nil {
			ctx.Logger().Error("Public key is not set ", "Address", acc.GetAddress())
			return false
		}

		for _, stdSig := range stdSigs {
			if pubKey.VerifyBytes(signBytes, stdSig.Signature) {
				signedBy, err := sdkTypes.AccAddressFromHex(pubKey.Address().String())
				if err != nil {
					return false
				}

				if !kycKeeper.IsWhitelisted(ctx, signedBy) {
					ctx.Logger().Error("Signer is not whitelisted", "Address", signedBy)
					return false
				}

				// Extra checks
				if !acc.GetAddress().Equals(signedBy) {
					ctx.Logger().Error("Public is match with the address ", "Address", signedBy)
					return false
				}

				if !acc.IsSigner(signedBy) {
					ctx.Logger().Error("Unauthorized signer", "Address", signedBy)
					return false
				}

				*validSignatures = *validSignatures + 1
				return true
			}
		}
	}

	return false
}
