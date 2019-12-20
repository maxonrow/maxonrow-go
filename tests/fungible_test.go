package tests

import (
	"testing"

	token "github.com/maxonrow/maxonrow-go/x/token/fungible"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func makeCreateFungibleTokenMsg(t *testing.T, name, symbol, metadata, owner, maxSupply, applicationFee, tokenFeeCollector string, decimals int, fixedSupply bool) sdkTypes.Msg {

	// create new token
	ownerAddr := tKeys[owner].addr
	maxSupplyUint := sdkTypes.NewUintFromString(maxSupply)
	fee := token.Fee{
		To:    tKeys[tokenFeeCollector].addr,
		Value: applicationFee,
	}
	msgCreateFungibleToken := token.NewMsgCreateFungibleToken(symbol, decimals, ownerAddr, name, fixedSupply, maxSupplyUint, metadata, fee)

	return msgCreateFungibleToken
}

func makeApproveFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, status string, burnable bool, feeSettingName string) sdkTypes.Msg {

	providerAddr := tKeys[provider].addr

	var tokenFee = []token.TokenFee{
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
	tokenDoc := token.NewToken(providerAddr, providerNonce, status, symbol, burnable, tokenFee)

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := token.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []token.Signature
	signature := token.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetFungibleTokenStatus := token.NewMsgSetFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetFungibleTokenStatus
}

//module of transfer
func makeTransferFungibleTokenMsg(t *testing.T, owner string, newOwner string, symbol string, transferAmountOfToken string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr
	transferAmountOfTokenUint := sdkTypes.NewUintFromString(transferAmountOfToken)

	msgTransferPayload := token.NewMsgTransferFungibleToken(symbol, transferAmountOfTokenUint, ownerAddr, newOwnerAddr)
	return msgTransferPayload

}

//module of mint
func makeMintFungibleTokenMsg(t *testing.T, owner string, newOwner string, symbol string, mintAmountOfToken string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr
	mintAmountOfTokenUint := sdkTypes.NewUintFromString(mintAmountOfToken)

	msgMintPayload := token.NewMsgIssueFungibleAsset(ownerAddr, symbol, newOwnerAddr, mintAmountOfTokenUint)
	return msgMintPayload

}

//module of burn
func makeBurnFungibleTokenMsg(t *testing.T, owner string, symbol string, burnAmountOfToken string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	burnAmountOfTokenUint := sdkTypes.NewUintFromString(burnAmountOfToken)

	msgTransferPayload := token.NewMsgBurnFungibleToken(symbol, burnAmountOfTokenUint, ownerAddr)
	return msgTransferPayload

}

//moduel of transferOwnership
func makeTransferFungibleTokenOwnershipMsg(t *testing.T, owner string, newOwner string, symbol string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr

	msgTransferOwnershipPayload := token.NewMsgTransferFungibleTokenOwnership(symbol, ownerAddr, newOwnerAddr)
	return msgTransferOwnershipPayload

}

//module of acceptOwnership
func makeAcceptFungibleTokenOwnershipMsg(t *testing.T, newOwner string, symbol string) sdkTypes.Msg {

	fromAddr := tKeys[newOwner].addr

	msgAcceptOwnershipPayload := token.NewMsgAcceptFungibleTokenOwnership(symbol, fromAddr)
	return msgAcceptOwnershipPayload

}

func makeFreezeFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, burnable bool) sdkTypes.Msg {

	status := "FREEZE"
	providerAddr := tKeys[provider].addr

	tokenDoc := token.NewToken(providerAddr, providerNonce, status, symbol, burnable, nil) // status : FREEZE / UNFREEZE

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := token.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []token.Signature
	signature := token.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetFungibleTokenStatus := token.NewMsgSetFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetFungibleTokenStatus
}

func makeUnfreezeFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, burnable bool) sdkTypes.Msg {

	status := "UNFREEZE"
	providerAddr := tKeys[provider].addr

	tokenDoc := token.NewToken(providerAddr, providerNonce, status, symbol, burnable, nil) // status : FREEZE / UNFREEZE

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := token.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []token.Signature
	signature := token.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetFungibleTokenStatus := token.NewMsgSetFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetFungibleTokenStatus
}

func makeVerifyTransferTokenOwnership(t *testing.T, signer, provider, providerNonce, issuer, symbol, action string) sdkTypes.Msg {

	providerAddr := tKeys[provider].addr

	// burnable and tokenfees is not in used for verifying transfer token status, we just set it to false and leave it empty.
	verifyTransferTokenOwnershipDoc := token.NewToken(providerAddr, providerNonce, action, symbol, false, []token.TokenFee{})

	// provider sign
	verifyTransferTokenOwnershipDocBz, err := tCdc.MarshalJSON(verifyTransferTokenOwnershipDoc)
	require.NoError(t, err)
	signedVerifyTransferTokenOwnershipDoc, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(verifyTransferTokenOwnershipDocBz))
	require.NoError(t, err)

	verifyTransferTokenOwnershipPayload := token.NewPayload(*verifyTransferTokenOwnershipDoc, tKeys[provider].pub, signedVerifyTransferTokenOwnershipDoc)

	// issuer sign
	verifyTransferPayloadBz, err := tCdc.MarshalJSON(verifyTransferTokenOwnershipPayload)
	require.NoError(t, err)
	signedVerifyTransferPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(verifyTransferPayloadBz))
	require.NoError(t, err)

	var signatures []token.Signature
	signature := token.NewSignature(tKeys[issuer].pub, signedVerifyTransferPayloadBz)
	signatures = append(signatures, signature)

	msgVerifyTransferTokenOwnership := token.NewMsgSetFungibleTokenStatus(tKeys[signer].addr, *verifyTransferTokenOwnershipPayload, signatures)

	return msgVerifyTransferTokenOwnership
}
