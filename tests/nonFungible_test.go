package tests

import (
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	nonFungible "github.com/maxonrow/maxonrow-go/x/token/nonfungible"
	"github.com/stretchr/testify/require"
)

func makeCreateNonFungibleTokenMsg(t *testing.T, name, symbol, metadata, owner, applicationFee, tokenFeeCollector string) sdkTypes.Msg {

	// create new token
	ownerAddr := tKeys[owner].addr
	fee := nonFungible.Fee{
		To:    tKeys[tokenFeeCollector].addr,
		Value: applicationFee,
	}
	msgCreateNonFungibleToken := nonFungible.NewMsgCreateNonFungibleToken(symbol, ownerAddr, name, metadata, fee)

	return msgCreateNonFungibleToken
}

func makeApproveNonFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, status string, feeSettingName string, mintLimit, transferLimit string, endorserList []string) sdkTypes.Msg {

	providerAddr := tKeys[provider].addr

	var tokenFee = []nonFungible.TokenFee{
		{
			Action:  "transfer",
			FeeName: feeSettingName,
		},
		{
			Action:  "mint",
			FeeName: feeSettingName,
		},
		{
			Action:  "burn",
			FeeName: feeSettingName,
		},
		{
			Action:  "transferOwnership",
			FeeName: feeSettingName,
		},
		{
			Action:  "acceptOwnership",
			FeeName: feeSettingName,
		},
	}

	mintL := sdkTypes.NewUintFromString(mintLimit)
	transferL := sdkTypes.NewUintFromString(transferLimit)

	var endorsers []sdkTypes.AccAddress

	for _, v := range endorserList {
		addr := tKeys[v].addr
		endorsers = append(endorsers, addr)
	}

	tokenDoc := nonFungible.NewToken(providerAddr, providerNonce, status, symbol, transferL, mintL, tokenFee, endorsers)

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := nonFungible.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetFungibleTokenStatus := nonFungible.NewMsgSetNonFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetFungibleTokenStatus
}

//module of transfer
func makeTransferNonFungibleTokenMsg(t *testing.T, owner string, newOwner string, symbol string, itemID []byte) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr

	msgTransferPayload := nonFungible.NewMsgTransferNonFungibleToken(symbol, ownerAddr, newOwnerAddr, itemID)
	return msgTransferPayload

}

//module of mint
func makeMintNonFungibleTokenMsg(t *testing.T, owner string, newOwner string, symbol string, itemID []byte, properties, metadata []string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr

	msgMintPayload := nonFungible.NewMsgMintNonFungibleToken(ownerAddr, symbol, newOwnerAddr, itemID, properties, metadata)
	return msgMintPayload

}

//module of burn
func makeBurnNonFungibleTokenMsg(t *testing.T, owner string, symbol string, itemID []byte) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr

	msgBurnNonFungible := nonFungible.NewMsgBurnNonFungibleToken(symbol, ownerAddr, itemID)
	return msgBurnNonFungible

}

//moduel of transferOwnership
func makeTransferNonFungibleTokenOwnershipMsg(t *testing.T, owner string, newOwner string, symbol string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr

	msgTransferOwnershipPayload := nonFungible.NewMsgTransferNonFungibleTokenOwnership(symbol, ownerAddr, newOwnerAddr)
	return msgTransferOwnershipPayload

}

//module of acceptOwnership
func makeAcceptNonFungibleTokenOwnershipMsg(t *testing.T, newOwner string, symbol string) sdkTypes.Msg {

	fromAddr := tKeys[newOwner].addr

	msgAcceptOwnershipPayload := nonFungible.NewMsgAcceptNonFungibleTokenOwnership(symbol, fromAddr)
	return msgAcceptOwnershipPayload

}

func makeFreezeNonFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, burnable bool) sdkTypes.Msg {

	status := "FREEZE"
	providerAddr := tKeys[provider].addr

	tokenDoc := nonFungible.NewToken(providerAddr, providerNonce, status, symbol, sdkTypes.NewUint(0), sdkTypes.NewUint(0), nil, nil) // status : FREEZE / UNFREEZE

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := nonFungible.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetNonFungibleTokenStatus := nonFungible.NewMsgSetNonFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetNonFungibleTokenStatus
}

func makeUnfreezeNonFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, burnable bool) sdkTypes.Msg {

	status := "UNFREEZE"
	providerAddr := tKeys[provider].addr

	tokenDoc := nonFungible.NewToken(providerAddr, providerNonce, status, symbol, sdkTypes.NewUint(0), sdkTypes.NewUint(0), nil, nil) // status : FREEZE / UNFREEZE

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := nonFungible.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetNonFungibleTokenStatus := nonFungible.NewMsgSetNonFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetNonFungibleTokenStatus
}

func makeVerifyTransferNonFungibleTokenOwnership(t *testing.T, signer, provider, providerNonce, issuer, symbol, action string) sdkTypes.Msg {

	providerAddr := tKeys[provider].addr

	// burnable and tokenfees is not in used for verifying transfer token status, we just set it to false and leave it empty.
	verifyTransferTokenOwnershipDoc := nonFungible.NewToken(providerAddr, providerNonce, action, symbol, sdkTypes.NewUint(0), sdkTypes.NewUint(0), []nonFungible.TokenFee{}, nil)

	// provider sign
	verifyTransferTokenOwnershipDocBz, err := tCdc.MarshalJSON(verifyTransferTokenOwnershipDoc)
	require.NoError(t, err)
	signedVerifyTransferTokenOwnershipDoc, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(verifyTransferTokenOwnershipDocBz))
	require.NoError(t, err)

	verifyTransferTokenOwnershipPayload := nonFungible.NewPayload(*verifyTransferTokenOwnershipDoc, tKeys[provider].pub, signedVerifyTransferTokenOwnershipDoc)

	// issuer sign
	verifyTransferPayloadBz, err := tCdc.MarshalJSON(verifyTransferTokenOwnershipPayload)
	require.NoError(t, err)
	signedVerifyTransferPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(verifyTransferPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedVerifyTransferPayloadBz)
	signatures = append(signatures, signature)

	msgVerifyTransferNonFungibleTokenOwnership := nonFungible.NewMsgSetNonFungibleTokenStatus(tKeys[signer].addr, *verifyTransferTokenOwnershipPayload, signatures)

	return msgVerifyTransferNonFungibleTokenOwnership
}

func makeEndorsement(t *testing.T, signer, to, symbol string, itemID []byte) sdkTypes.Msg {

	signerAddr := tKeys[signer].addr
	toAddr := tKeys[to].addr

	return nonFungible.NewMsgEndorsement(symbol, signerAddr, toAddr, itemID)
}
