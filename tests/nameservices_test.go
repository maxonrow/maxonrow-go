package tests

import (
	"testing"

	"github.com/maxonrow/maxonrow-go/x/nameservice"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

type NameServiceInfo struct {
	Action         string
	Name           string
	From           string
	ApplicationFee string
	FeeCollector   string
	Provider       string
	ProviderNonce  string
	Issuer         string
	approved       string
}

func makeNameservicesTxs() []*testCase {

	tcs := []*testCase{

		{"nameservice", true, true, "creating alias with name mxw-alias-invalid signer", "yk", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias", "nago", "10000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", false, false, "creating alias with name mxw-alias", "nago", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias", "nago", "10000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", false, true, "creating alias with name mxw-alias-again", "nago", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias2", "nago", "10000", "ns-feecollector", "", "", "", ""}, "", nil},

		{"nameservice", false, false, "creating alias with name mxw-alias-1", "acc-40", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias-1", "acc-40", "1000000000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", false, false, "creating alias with name mxw-alias-2", "carlo", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias-2", "carlo", "1000000000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", false, false, "creating alias with name MXW-ALIAS", "yk", "100000000cin", 0, NameServiceInfo{"create", "MXW-ALIAS", "yk", "10000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", false, false, "creating alias with name MXW-@LI@S", "dont-use-this-1", "100000000cin", 0, NameServiceInfo{"create", "MXW-@LI@S", "dont-use-this-1", "10000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", true, true, "creating alias with wrong free collector", "jeansoon", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias", "jeansoon", "1000000000", "acc-40", "", "", "", ""}, "", nil},
		{"nameservice", false, true, "Already existed  name mxw-alias", "nago", "100000000cin", 0, NameServiceInfo{"create", "new-alias", "nago", "100000000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", true, true, "create alias with zero system fee", "nago", "0cin", 0, NameServiceInfo{"create", "mxw-alias", "nago", "100000000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", false, true, "create alias with zero application-fee, have pending", "nago", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias-nago", "nago", "0", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", false, false, "create alias with zero application-fee", "mostafa", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias-mostafa", "mostafa", "0", "ns-feecollector", "", "", "", ""}, "", nil},
		///{"nameservice", true, true, "create alias with non-kyc account", "gohck", "100000000cin", 0, NameServiceInfo{"create", "mxwone", "gohck", "100000000", "ns-feecollector", "", "", "", ""}, "", nil}, //kiv
		{"nameservice", true, true, "creating alias with signer and owner are different", "nago", "100000000cin", 0, NameServiceInfo{"create", "mxwone", "gohck", "100000000", "ns-feecollector", "", "", "", ""}, "", nil},

		//NameServie
		//Approve
		{"nameservice", true, true, "Non nameservice authorizer account try to authorize", "gohck", "0cin", 0, NameServiceInfo{"approve", "mxw-alias", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", true, true, "all account different than nameservice provider,issuer,auth account try to sign", "mostafa", "0cin", 0, NameServiceInfo{"approve", "mxw-alias", "", "", " ", "nago", "0", "acc-40", "true"}, "", nil},

		{"nameservice", true, true, "Approve name mxw-alias without authorizer", "nago", "0cin", 0, NameServiceInfo{"approve", "mxw-alias-2", "", "", "", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", false, false, "Approve name mxw-alias", "ns-auth", "0cin", 0, NameServiceInfo{"approve", "mxw-alias", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", false, true, "Approve name mxw-alias-repeated", "ns-auth", "0cin", 0, NameServiceInfo{"approve", "mxw-alias", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},

		//NameService
		//Revoke
		{"nameservice", true, true, "Non Authorizer revoke name mxw-alias", "acc-40", "0cin", 0, NameServiceInfo{"revoke", "mxw-alias", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", true, true, "All  different accounts provider,issuer,auth in nameservice try to sign and revoke", "acc-40", "0cin", 0, NameServiceInfo{"revoke", "mxw-alias", "", "", " ", "nago", "0", "acc-40", "true"}, "", nil},
		{"nameservice", true, true, "Authorizer revoke name mxw-alias without issuer", "ns-auth", "0cin", 0, NameServiceInfo{"revoke", "mxw-alias-1", "", "0", "", "ns-provider", "0", "acc-40", "true"}, "", nil},
		{"nameservice", false, false, "Authorizer revoke  name mxw-alias", "ns-auth", "0cin", 0, NameServiceInfo{"revoke", "mxw-alias", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", false, true, "Authorizer revoke  name mxw-alias-again", "ns-auth", "0cin", 0, NameServiceInfo{"revoke", "mxw-alias", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},

		//NameService
		//Reject
		{"nameservice", true, true, "Non Authorizer reject name mxw-alias", "acc-40", "0cin", 0, NameServiceInfo{"reject", "mxw-alias", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", true, true, "all account different than nameservice provider,issuer,auth account try to sign", "acc-40", "0cin", 0, NameServiceInfo{"reject", "mxw-alias", "", "", " ", "nago", "0", "acc-40", "true"}, "", nil},
		{"nameservice", false, false, "Authorizer reject name mxw-alias-1", "ns-auth", "0cin", 0, NameServiceInfo{"reject", "mxw-alias-1", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", false, false, "Authorizer reject name mxw-alias-2", "ns-auth", "0cin", 0, NameServiceInfo{"reject", "mxw-alias-2", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", false, true, "Authorizer reject name mxw-alias-2-again", "ns-auth", "0cin", 0, NameServiceInfo{"reject", "mxw-alias-2", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
	}

	return tcs
}

func makeCreateNameServiceMsg(t *testing.T, name, owner, applicationFee, aliasFeeCollector string) sdkTypes.Msg {

	// create new alias
	ownerAddr := tKeys[owner].addr
	fee := nameservice.Fee{
		To:    tKeys[aliasFeeCollector].addr,
		Value: applicationFee,
	}
	msgCreateAlias := nameservice.NewMsgCreateAlias(name, ownerAddr, fee)

	return msgCreateAlias
}

func setStatusAliasMsg(t *testing.T, signer, provider, providerNonce, issuer, name, status string) sdkTypes.Msg {

	providerAddr := tKeys[provider].addr

	nameserviceDoc := nameservice.NewAlias(providerAddr, providerNonce, status, name)

	// provider sign the nameservice
	nsProvider, err := tCdc.MarshalJSON(nameserviceDoc)
	require.NoError(t, err)
	signedAlias, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(nsProvider))
	require.NoError(t, err)

	nameservicePayload := nameservice.NewPayload(*nameserviceDoc, tKeys[provider].pub, signedAlias)

	// issuer sign the nameservice
	aliasPayload, err := tCdc.MarshalJSON(nameservicePayload)
	require.NoError(t, err)
	signedAliasPayload, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(aliasPayload))
	require.NoError(t, err)

	var signatures []nameservice.Signature
	signature := nameservice.NewSignature(tKeys[issuer].pub, signedAliasPayload)
	signatures = append(signatures, signature)

	msgSetFungibleTokenStatus := nameservice.NewMsgSetAliasStatus(tKeys[signer].addr, *nameservicePayload, signatures)

	return msgSetFungibleTokenStatus
}
