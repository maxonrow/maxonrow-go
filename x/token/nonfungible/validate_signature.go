package nonfungible

import (
	"bytes"
	"fmt"
	"strconv"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/tendermint/tendermint/crypto"
)

func (k Keeper) ValidateSignatures(ctx sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Error {

	var fromSignature Signature
	var issuerSignatures []Signature
	var fromSignBytes, issuerSignBytes []byte
	var fromAddr sdkTypes.AccAddress
	var fromAccountNonce string

	switch msg := msg.(type) {
	case MsgSetNonFungibleTokenStatus:
		fromSignature = NewSignature(msg.Payload.PubKey, msg.Payload.Signature)
		// * from sign bytes
		fromSignBytes = msg.Payload.Token.GetFromSignBytes()
		fromAddr = msg.Payload.Token.From
		fromAccountNonce = msg.Payload.Token.Nonce

		//* issuer sign bytes
		issuerSignBytes = msg.Payload.GetIssuerSignBytes()
		issuerSignatures = msg.Signatures
	case MsgSetNonFungibleItemStatus:
		fromSignature = NewSignature(msg.ItemPayload.PubKey, msg.ItemPayload.Signature)
		// * from sign bytes
		fromSignBytes = msg.ItemPayload.Item.GetAccountStatusSettingFromSignBytes()
		fromAddr = msg.ItemPayload.Item.From
		fromAccountNonce = msg.ItemPayload.Item.Nonce

		//* issuer sign bytes
		issuerSignBytes = msg.ItemPayload.GetAccountStatusSettingSignBytes()
		issuerSignatures = msg.Signatures
	default:
		errMsg := fmt.Sprintf("Invalid signature for non-fungible token: %v", msg.Type())
		return sdkTypes.ErrUnknownRequest(errMsg)
	}

	//* check account sequence with passed in nonce
	acc := k.accountKeeper.GetAccount(ctx, fromAddr)
	if acc == nil {
		nonce, nonceErr := strconv.ParseUint(fromAccountNonce, 10, 64)
		if nonceErr != nil {
			return sdkTypes.ErrInvalidSequence("Wallet signature is invalid.")
		}

		if nonce != 0 {
			return sdkTypes.ErrInvalidSequence("Wallet signature is invalid.")
		}

	} else {
		nonce, nonceErr := strconv.ParseUint(fromAccountNonce, 10, 64)
		sequence := acc.GetSequence()
		if nonceErr != nil {
			return sdkTypes.ErrInvalidSequence("Wallet signature is invalid.")
		}

		if nonce != sequence {
			return sdkTypes.ErrInvalidSequence("Wallet signature is invalid.")
		}

		acc.SetSequence(sequence + 1)
	}

	if !k.IsProvider(ctx, fromAddr) {
		return sdkTypes.ErrUnauthorized("Insufficient provider signature.")
	}

	//* verify from sign
	if !(processSig(acc, fromSignature, fromSignBytes)) {

		return sdkTypes.ErrUnauthorized("From signature verification failed.")
	}

	//* at least one issuer
	issuerCounter := 0

	//* verify issuer sign
	for i := 0; i < len(issuerSignatures); i++ {
		issuerAddr := sdkTypes.AccAddress(issuerSignatures[i].PubKey.Address())
		issuerAcc := k.accountKeeper.GetAccount(ctx, issuerAddr)

		if k.IsIssuer(ctx, issuerAddr) {
			issuerCounter++
		} else {
			return sdkTypes.ErrUnauthorized("Unauthorized signature.")
		}

		if !(processSig(issuerAcc, issuerSignatures[i], issuerSignBytes)) {

			return sdkTypes.ErrUnauthorized("Signature verification failed.")
		}

		issuerAcc.SetSequence(issuerAcc.GetSequence() + 1)
	}

	if issuerCounter < 1 {
		return sdkTypes.ErrUnauthorized("Insufficient issuer signature.")
	}

	return nil
}

func processSig(
	signerAcc exported.Account, signature Signature, signBytes []byte) bool {

	pubKey, res := ProcessPubKey(signerAcc, signature)
	if !res.IsOK() {

		return false
	}

	if signerAcc != nil {
		err := signerAcc.SetPubKey(pubKey)
		if err != nil {

			return false
		}
	}

	return pubKey.VerifyBytes(signBytes, signature.Signature)
}

// ProcessPubKey verifies that the given account address matches that of the
// StdSignature. In addition, it will set the public key of the account if it
// has not been set.
func ProcessPubKey(acc exported.Account, sig Signature) (crypto.PubKey, sdkTypes.Result) {
	var pubKey crypto.PubKey
	if acc != nil {
		pubKey = acc.GetPubKey()
	}

	if pubKey == nil {

		return sig.PubKey, sdkTypes.Result{}
	}
	if acc != nil {
		cryptoPubKey := pubKey
		if cryptoPubKey == nil {
			return nil, sdkTypes.ErrInvalidPubKey("PubKey not found").Result()
		}

		if !bytes.Equal(cryptoPubKey.Address(), acc.GetAddress()) {
			return nil, sdkTypes.ErrInvalidPubKey(
				fmt.Sprintf("PubKey does not match Signer address %s", acc.GetAddress())).Result()
		}

		return cryptoPubKey, sdkTypes.Result{}
	}
	return nil, sdkTypes.Result{}
}
