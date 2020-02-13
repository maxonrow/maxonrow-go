package tests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/maxonrow/maxonrow-go/types"

	"github.com/maxonrow/maxonrow-go/x/bank"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	msgType       string
	checkFailed   bool
	deliverFailed bool
	desc          string
	signer        string
	fees          string
	gas           uint64
	msgInfo       interface{}
	memo          string
	hash          []byte
}

func TestTxs(t *testing.T) {

	val1 := Validator(tValidator)
	assert.Equal(t, val1.OperatorAddress.String(), tValidator)
	acc1 := Account(tKeys["alice"].addrStr)
	bal1 := acc1.GetCoins()[0]

	_, err := time.ParseDuration("60s")
	if err != nil {
		fmt.Println("timeout", err)
	}

	var tcs []*testCase

	tcs = append(tcs, makeBankTxs()...)
	tcs = append(tcs, makeKycTxs()...)
	tcs = append(tcs, makeMaintenaceTxs()...)
	tcs = append(tcs, makeFeeTxs()...)
	tcs = append(tcs, makeNonFungibleTokenTxs()...)
	tcs = append(tcs, makeFungibleTokenTxs()...)
	tcs = append(tcs, makeNameservicesTxs()...)

	var totalFee = sdkTypes.NewInt64Coin("cin", 0)
	var totalAmt = sdkTypes.NewInt64Coin("cin", 0)

	for n, tc := range tcs {

		fees, err := types.ParseCoins(tc.fees)
		assert.NoError(t, err, tc.fees)

		var msg sdkTypes.Msg
		switch tc.msgType {
		case "bank":
			{
				i := tc.msgInfo.(bankInfo)
				amt, err := types.ParseCoins(i.amount)
				assert.NoError(t, err, i.amount)

				msg = bank.NewMsgSend(tKeys[i.from].addr, tKeys[i.to].addr, amt)

				if !tc.deliverFailed {
					if i.from == "alice" {
						totalFee = totalFee.Add(fees[0])
						totalAmt = totalAmt.Add(amt[0])
					}
				}
			}
		case "kyc":
			{
				i := tc.msgInfo.(kycInfo)

				switch i.action {
				case "whitelist":
					msg = makeKycWhitelistMsg(t, i.authorised, i.issuer, i.provider, i.from, i.signer, i.data, i.nonce)

				case "revokeWhitelist":
					msg = makeKycRevokeWhitelistMsg(t, i.authorised, i.issuer, i.provider, i.from)
				}
			}

		case "fee":
			{
				i := tc.msgInfo.(feeInfo)

				msg = makeFeeMsg(t, i.function, i.name, i.assignee, i.multiplier, i.min, i.max, i.percentage, i.issuer)
			}

		case "nameservice":
			{
				i := tc.msgInfo.(NameServiceInfo)

				switch i.Action {
				case "create":
					msg = makeCreateNameServiceMsg(t, i.Name, i.From, i.ApplicationFee, i.FeeCollector)
				case "approve":
					msg = setStatusAliasMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Name, "APPROVE")
				case "reject":
					msg = setStatusAliasMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Name, "REJECT")
				case "revoke":
					msg = setStatusAliasMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Name, "REVOKE")
				}

			}
		case "token":
			{
				i := tc.msgInfo.(TokenInfo)

				switch i.Action {
				case "create":
					msg = makeCreateFungibleTokenMsg(t, i.Name, i.Symbol, i.Metadata, i.Owner, i.MaxSupply, i.ApplicationFee, i.FeeCollector, i.Decimals, i.FixedSupply)
				case "approve":
					msg = makeApproveFungibleTokenMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, "APPROVE", i.Burnable, i.FeeSettingName)
				case "transfer":
					msg = makeTransferFungibleTokenMsg(t, i.Owner, i.NewOwner, i.Symbol, i.AmountOfToken)
				case "mint":
					msg = makeMintFungibleTokenMsg(t, i.Owner, i.NewOwner, i.Symbol, i.AmountOfToken)
				case "burn":
					msg = makeBurnFungibleTokenMsg(t, i.Owner, i.Symbol, i.AmountOfToken)
				case "transfer-ownership":
					msg = makeTransferFungibleTokenOwnershipMsg(t, i.Owner, i.NewOwner, i.Symbol)
				case "accept-ownership":
					msg = makeAcceptFungibleTokenOwnershipMsg(t, i.NewOwner, i.Symbol)
				case "freeze":
					msg = makeFreezeFungibleTokenMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.Burnable)
				case "unfreeze":
					msg = makeUnfreezeFungibleTokenMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.Burnable)
				case "verify-transfer-tokenOwnership":
					msg = makeVerifyTransferTokenOwnership(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.VerifyTransferTokenOwnership)
				}
			}
		case "nonFungibleToken":
			{
				i := tc.msgInfo.(NonFungibleTokenInfo)

				switch i.Action {
				case "create":
					msg = makeCreateNonFungibleTokenMsg(t, i.Name, i.Symbol, i.TokenMetadata, i.Owner, i.ApplicationFee, i.FeeCollector)
				case "approve":
					msg = makeApproveNonFungibleTokenMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, "APPROVE", i.FeeSettingName, i.MintLimit, i.TransferLimit, i.EndorserList, i.Burnable, i.Modifiable, i.Public)
				case "transfer":
					msg = makeTransferNonFungibleTokenMsg(t, i.Owner, i.NewOwner, i.Symbol, i.ItemID)
				case "mint":
					msg = makeMintNonFungibleTokenMsg(t, i.Owner, i.NewOwner, i.Symbol, i.ItemID, i.Properties, i.Metadata)
				case "burn":
					msg = makeBurnNonFungibleTokenMsg(t, i.Owner, i.Symbol, i.ItemID)
				case "transfer-ownership":
					msg = makeTransferNonFungibleTokenOwnershipMsg(t, i.Owner, i.NewOwner, i.Symbol)
				case "accept-ownership":
					msg = makeAcceptNonFungibleTokenOwnershipMsg(t, i.NewOwner, i.Symbol)
				case "freeze-item":
					msg = makeFreezeNonFungibleItemMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.ItemID)
				case "unfreeze-item":
					msg = makeUnfreezeNonFungibleItemMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.ItemID)
				case "freeze":
					msg = makeFreezeNonFungibleTokenMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.Burnable, i.Modifiable, i.Public)
				case "unfreeze":
					msg = makeUnfreezeNonFungibleTokenMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.Burnable, i.Modifiable, i.Public)
				case "verify-transfer-tokenOwnership":
					msg = makeVerifyTransferNonFungibleTokenOwnership(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.VerifyTransferTokenOwnership, i.Burnable, i.Modifiable, i.Public)
				case "endorsement":
					msg = makeEndorsement(t, tc.signer, i.Owner, i.Symbol, i.ItemID)
				case "update-item-metadata":
					msg = makeUpdateItemMetadataMsg(t, i.Symbol, i.Owner, i.ItemID, i.Metadata)
				case "update-nft-metadata":
					msg = makeUpdateNFTMetadataMsg(t, i.Symbol, i.Owner, i.TokenMetadata)
				}
			}
		case "maintenance":
			{
				i := tc.msgInfo.(MaintenanceInfo)
				msg = makeMaintenanceMsg(t, i.Action, i.Title, i.Description, i.ProposalType, i.AuthorisedAddress, i.IssuerAddress, i.ProviderAddress, i.Proposer, i.ValidatorPubKey, i.FeeCollector)
			}

		case "maintenance-cast-action":
			{
				i := tc.msgInfo.(CastAction)
				msg = makeCastActionMsg(t, i.Action, i.Caster, i.ProposalId)
			}
		}

		tx, bz := MakeSignedTx(t, tc.signer, tc.gas, fees, tc.memo, msg)
		txb, err := tCdc.MarshalJSON(tx)
		assert.NoError(t, err)
		fmt.Printf("============\nBroadcasting %d tx (%s): %s\n", n+1, tc.desc, string(txb))

		res := BroadcastTxCommit(bz)
		tc.hash = res.Hash.Bytes()

		if !tc.checkFailed {
			tKeys[tc.signer].seq++
			require.Zero(t, res.CheckTx.Code, "test case %v(%v) check failed: %v", n+1, tc.desc, res.CheckTx.Log)

			if !tc.deliverFailed {
				require.Zero(t, res.DeliverTx.Code, "test case %v(%v) deliver failed: %v", n+1, tc.desc, res.DeliverTx.Log)
			} else {
				require.NotZero(t, res.DeliverTx.Code, "test case %v(%v) deliver failed: %v", n+1, tc.desc, res.DeliverTx.Log)
			}

		} else {
			require.NotZero(t, res.CheckTx.Code, "test case %v(%v) check failed: %v", n+1, tc.desc, res.CheckTx.Log)
		}

		if strings.Contains(tc.desc, "commit") {
			WaitForNextHeightTM(tPort)
		}
	}

	WaitForNextHeightTM(tPort)

	/// check results
	for i, tc := range tcs {
		res := Tx(tc.hash)
		if tc.checkFailed {
			assert.Nil(t, res, "test case %v(%v) failed", i, tc.desc)
		} else {
			if tc.deliverFailed {
				assert.NotZero(t, res.TxResult.Code, "test case %v(%v) failed", i, tc.desc)
			} else {
				assert.Zero(t, res.TxResult.Code, "test case %v(%v) failed: %v", i, tc.desc, res.TxResult.Log)
			}
		}
	}

	acc2 := Account(tKeys["alice"].addrStr)
	bal2 := acc2.GetCoins()[0]
	diff := bal1.Sub(bal2)

	total := totalAmt.Add(totalFee)
	require.Equal(t, diff, total)

	accGohck := Account(tKeys["gohck"].addrStr)
	require.Empty(t, accGohck.GetCoins())

	fmt.Println(totalFee)
	//mxwcli query distribution validator-outstanding-rewards mxwvaloper1rjgjjkkjqtd676ydahysmnfsg0v4yvwfp2n965
}

func MakeSignedTx(t *testing.T, name string, gas uint64, fees sdkTypes.Coins, memo string, msg sdkTypes.Msg) (sdkAuth.StdTx, []byte) {
	acc := Account(tKeys[name].addrStr)
	require.NotNil(t, acc)

	signMsg := authTypes.StdSignMsg{
		AccountNumber: acc.GetAccountNumber(),
		ChainID:       "maxonrow-chain",
		Fee:           authTypes.NewStdFee(gas, fees),
		Memo:          memo,
		Msgs:          []sdkTypes.Msg{msg},
		Sequence:      tKeys[name].seq,
	}

	signBz, signBzErr := tCdc.MarshalJSON(signMsg)
	if signBzErr != nil {
		panic(signBzErr)
	}

	sig, err := tKeys[name].priv.Sign(sdkTypes.MustSortJSON(signBz))
	if err != nil {
		panic(err)
	}

	pub := tKeys[name].priv.PubKey()
	stdSig := sdkAuth.StdSignature{
		PubKey:    pub,
		Signature: sig,
	}

	sdtTx := authTypes.NewStdTx(signMsg.Msgs, signMsg.Fee, []authTypes.StdSignature{stdSig}, signMsg.Memo)
	bz, err := tCdc.MarshalBinaryLengthPrefixed(sdtTx)
	if err != nil {
		panic(err)
	}
	return sdtTx, bz
}
