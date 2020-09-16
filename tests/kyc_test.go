package tests

import (
	"strconv"
	"testing"

	"github.com/maxonrow/maxonrow-go/x/kyc"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

type kycInfo struct {
	authorised string
	issuer     string
	provider   string
	action     string
	from       string
	signer     string
	data       string
	nonce      string
}

func makeKycTxs() []*testCase {

	tcs := []*testCase{

		//Whitelist
		{"kyc", true, true, "Doing kyc - INVALID SIGNER-1", "bob", "0cin", 0, kycInfo{"bob", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - INVALID SIGNER-2", "kyc-issuer-1", "0cin", 0, kycInfo{"kyc-issuer-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - MISMATCH SIGNER", "kyc-auth-2", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1234", "0"}, "", nil},
		{"kyc", false, false, "Doing kyc - HAPPY PATH-commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - EMPTY DATA", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "acc-30", "acc-30", "", "0"}, "", nil},

		// goh - last-time
		{"kyc", false, false, "Doing kyc - HAPPY PATH for dont-use-this-1", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "dont-use-this-1", "dont-use-this-1", "testKyc1251", "0"}, "", nil},

		{"kyc", true, true, "Doing kyc - KYC AGAIN!", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1234", "0"}, "", nil}, //code 1
		{"kyc", false, false, "Doing kyc - KYC AGAIN, DIFFERENT KYC DATA-commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc12399", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - DUPLICATE KYC DATA", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "acc-23", "acc-23", "testKyc12399", "0"}, "", nil},

		{"kyc", true, true, "Doing kyc - NOT ISSUER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "jeansoon", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc12345", "0"}, "", nil},                        //code 4
		{"kyc", true, true, "Doing kyc - NO ISSUER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "nope", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc12345", "0"}, "", nil},                             //code 4
		{"kyc", true, true, "Doing kyc - NOT AUTHORISED", "jeansoon", "0cin", 0, kycInfo{"jeansoon", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc12346", "0"}, "", nil},                    //code 1
		{"kyc", true, true, "Doing kyc - NOT AUTHORISED, BUT WHITELISTED ADDRESS", "alice", "0cin", 0, kycInfo{"alice", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc12346", "0"}, "", nil}, //code 1
		{"kyc", true, true, "Doing kyc - NOT PROVIDER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "jeansoon", "whitelist", "josephin", "josephin", "testKyc12347", "0"}, "", nil},                    //code 4
		{"kyc", true, true, "Doing kyc - NO PROVIDER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "nope", "whitelist", "josephin", "josephin", "testKyc12347", "0"}, "", nil},                         //code 4
		{"kyc", true, true, "Doing kyc - WRONG NONCE", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1238", "1"}, "", nil},
		{"kyc", true, true, "Doing kyc - NOT PROVIDER, NOT ISSUER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "jeansoon", "acc-19", "whitelist", "josephin", "josephin", "testKyc1239", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - TWO PROVIDER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-prov-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1239", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - TWO ISSUER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-issuer-1", "whitelist", "josephin", "josephin", "testKyc1239", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - EMPTY DATA", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "jeansoon", "acc-19", "whitelist", "acc-24", "acc-24", "", "0"}, "", nil},
		{"kyc", false, false, "Doing kyc - PAY FEE-commit", "kyc-auth-1", "1000cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc123456", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - WRONG GAS", "kyc-auth-1", "0cin", 1, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - NOPE ADDR", "kyc-auth-1", "0cin", 1, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "nope", "acc-23", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - INVALID FROM SIGNATURE", "kyc-auth-1", "0cin", 1, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "acc-23", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - INVALID AUTHORISED SIGNATURE-1", "kyc-auth-2", "0cin", 1, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "acc-23", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - INVALID AUTHORISED SIGNATURE-2", "kyc-issuer-1", "0cin", 1, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "acc-23", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - INVALID AUTHORISED SIGNATURE", "kyc-prov-1", "0cin", 1, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "acc-23", "testKyc1234", "0"}, "", nil},

		// Hints : here as Special-case, which is allowed after done the Whitelist as above
		// It should works now
		{"bank", false, false, "sending after whitelisting an account", "josephin", "200000000cin", 0, bankInfo{"josephin", "bob", "2cin"}, "", nil},

		//RevokeWhitelist
		// revoke whitelist needs provider signature(RevokeKycData)
		// revoke whitelist needs issuer signature(RevokeKycPayload)
		// revoke whitelist needs authorised signature(tx)
		{"kyc", true, true, "Undoing kyc - INVALID SIGNER", "bob", "0cin", 0, kycInfo{"bob", "kyc-issuer-1", "jeansoon", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - MISMATCH SIGNER", "kyc-auth-2", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "jeansoon", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - NOT PROVIDER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "jeansoon", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - NO PROVIDER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "nope", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - NOT ISSUER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "jeansoon", "kyc-prov-1", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - NO ISSUER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "nope", "kyc-prov-1", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - TWO ISSUER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-issuer-1", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - TWO PROVIDER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-prov-1", "kyc-prov-1", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - REVOKE SOMEONE IS NOT WHITELISTED", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "revokeWhitelist", "acc-23", "", "", ""}, "", nil},
		{"kyc", false, false, "Undoing kyc - HAPPY PATH-commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "revokeWhitelist", "josephin", "", "", ""}, "", nil},

		// Hints : here as Special-case, which is allowed after done the RevokeWhitelist as above
		// It should fail now
		{"bank", true, true, "sending after revoking an account", "josephin", "200000000cin", 0, bankInfo{"josephin", "bob", "1cin"}, "", nil},
		{"bank", false, false, "receiving after revoking an account", "alice", "200000000cin", 0, bankInfo{"alice", "josephin", "1cin"}, "", nil},

		//=============================================start : used by nft modules
		{"kyc", false, false, "Doing kyc - nft-yk - commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "nft-yk", "nft-yk", "testKyc123451111", "0"}, "", nil},
		{"kyc", false, false, "Doing kyc - nft-mostafa - commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "nft-mostafa", "nft-mostafa", "testKyc123452222", "0"}, "", nil},
	}

	return tcs
}

func makeKycWhitelistMsg(t *testing.T, authorised, issuer, provider, from, signer, data, nonce string) sdkTypes.Msg {
	// create new kyc data to be whitelisted
	kycData := kyc.NewKyc(tKeys[from].addr, nonce, data)

	// kyc signed by the address which want to be whitelisted
	kycDataBz, err := tCdc.MarshalJSON(kycData)
	require.NoError(t, err)
	signedKycDataBz, err := tKeys[signer].priv.Sign(sdkTypes.MustSortJSON(kycDataBz))
	require.NoError(t, err)

	// creating the kyc payload
	kycPayload := kyc.NewPayload(kycData, tKeys[from].pub, signedKycDataBz)

	// kycPayload to be signed by issuer and provider
	kycPayloadBz, err := tCdc.MarshalJSON(kycPayload)
	require.NoError(t, err)
	var signatures []kyc.Signature

	if issuer != "nope" {
		issuerSignedBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(kycPayloadBz))
		require.NoError(t, err)
		issuerSignature := kyc.NewSignature(tKeys[issuer].pub, issuerSignedBz)
		signatures = append(signatures, issuerSignature)
	}

	if provider != "nope" {
		providerSignedBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(kycPayloadBz))
		require.NoError(t, err)
		providerSignature := kyc.NewSignature(tKeys[provider].pub, providerSignedBz)
		signatures = append(signatures, providerSignature)
	}
	toWhitelistData := kyc.NewKycData(kycPayload, signatures)

	return kyc.NewMsgWhitelist(tKeys[authorised].addr, toWhitelistData)
}

func makeKycRevokeWhitelistMsg(t *testing.T, authorised, issuer, provider, to string) sdkTypes.Msg {
	// convert uint64 to string
	var seq uint64
	if provider == "nope" {
		seq = 0
	} else {
		acc := Account(tKeys[provider].addrStr)
		require.NotNil(t, acc)
		seq = acc.GetSequence()
	}

	providerNonceStr := strconv.FormatUint(seq, 10)

	// create new kyc data to be revoked
	revokeKycData := kyc.NewRevokeKycData(tKeys[provider].addr, providerNonceStr, tKeys[to].addr)

	var signedRevokeKycDataBz []byte

	if provider != "nope" {
		// revokeKycData signed by the provider
		revokeKycDataBz, err := tCdc.MarshalJSON(revokeKycData)
		require.NoError(t, err)
		signedRevokeKycDataBz, err = tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(revokeKycDataBz))
		require.NoError(t, err)
	}
	// creating the revoke kyc payload
	revokeKycPayload := kyc.NewRevokePayload(revokeKycData, tKeys[provider].pub, signedRevokeKycDataBz)

	var signatures []kyc.Signature

	if issuer != "nope" {
		// revokeKycPayload to be signed by issuer
		revokeKycPayloadBz, err := tCdc.MarshalJSON(revokeKycPayload)
		require.NoError(t, err)
		issuerSignedBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(revokeKycPayloadBz))
		require.NoError(t, err)
		issuerSignature := kyc.NewSignature(tKeys[issuer].pub, issuerSignedBz)
		signatures = append(signatures, issuerSignature)
	}
	return kyc.NewMsgRevokeWhitelist(sdkTypes.AccAddress(tKeys[authorised].addr), revokeKycPayload, signatures)
}
