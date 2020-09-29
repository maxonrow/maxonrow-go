package tests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/utils"

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

	stdOut, _, _ := utils.RunProcess("", "mxwcli", []string{"query", "distribution", "validator-outstanding-rewards", "mxwvaloper1rjgjjkkjqtd676ydahysmnfsg0v4yvwfp2n965"})
	s := "0"
	if len(stdOut) > 24 {
		s = strings.TrimSpace(strings.Replace(stdOut[24:], "\"", "", -1))
	}
	initialRewards, _ := sdkTypes.NewDecFromStr(s)

	acc1 := Account(tKeys["alice"].addrStr)
	bal1 := acc1.GetCoins()[0]

	var tcs []*testCase

	tcs = append(tcs, makeBankTxs()...)
	tcs = append(tcs, makeKycTxs()...)
	tcs = append(tcs, makeMaintenaceTxs()...)
	tcs = append(tcs, makeFeeTxs()...)
	tcs = append(tcs, makeNonFungibleTokenTxs()...)
	tcs = append(tcs, makeFungibleTokenTxs()...)
	// tcs = append(tcs, makeNameservicesTxs()...)
	//tcs = append(tcs, makeMultisigTxs()...)
	//tcs = append(tcs, makeMultisigTxsNFTs()...)
	//tcs = append(tcs, makeMultisigTxsFTs()...)

	var totalFee = sdkTypes.NewInt64Coin("cin", 0)
	var totalFeeAlice = sdkTypes.NewInt64Coin("cin", 0)
	var totalAmtAlice = sdkTypes.NewInt64Coin("cin", 0)

	seqs := make(map[string]uint64)

	for n, tc := range tcs {

		fees, err := types.ParseCoins(tc.fees)
		assert.NoError(t, err, tc.fees)

		msg := makeMsg(t, tc.msgType, tc.signer, tc.msgInfo)
		switch tc.msgType {
		case "bank":
			{
				i := tc.msgInfo.(bankInfo)
				amt, err := types.ParseCoins(i.amount)
				assert.NoError(t, err, i.amount)

				if !tc.deliverFailed {
					if i.from == "alice" {
						totalFeeAlice = totalFeeAlice.Add(fees[0])
						totalAmtAlice = totalAmtAlice.Add(amt[0])
					}
				}
			}
		}

		tx, bz := makeSignedTx(t, tc.signer, tc.signer, 0, tc.gas, fees, tc.memo, msg)
		txb, err := tCdc.MarshalJSON(tx)
		assert.NoError(t, err)
		fmt.Printf("============\nBroadcasting %d tx (%s): %s\n", n+1, tc.desc, string(txb))

		res := BroadcastTxCommit(bz)
		tc.hash = res.Hash.Bytes()

		if !tc.checkFailed {
			seqs[tc.signer] = seqs[tc.signer] + 1
			require.Zero(t, res.CheckTx.Code, "test case %v(%v) check should not fail: %v", n+1, tc.desc, res.CheckTx.Log)

			if tc.deliverFailed {
				require.NotZero(t, res.DeliverTx.Code, "test case %v(%v) deliver should fail: %v", n+1, tc.desc, res.DeliverTx.Log)
			} else {
				require.Zero(t, res.DeliverTx.Code, "test case %v(%v) deliver should not fail: %v", n+1, tc.desc, res.DeliverTx.Log)
			}

			totalFee = totalFee.Add(fees[0])
		} else {
			require.NotZero(t, res.CheckTx.Code, "test case %v(%v) check failed: %v", n+1, tc.desc, res.CheckTx.Log)
		}

		if strings.Contains(tc.desc, "commit") {
			WaitForNextHeightTM(tPort)
		} else if strings.Contains(tc.desc, "wait-5-seconds") {
			time.Sleep(5 * time.Second) // need wait for 5 seconds, due to Blockchain consensus concern
		}
	}

	WaitForNextHeightTM(tPort)

	/// check results
	for i, tc := range tcs {
		res := Tx(tc.hash)
		if tc.checkFailed {
			require.Nil(t, res, "test case %v(%v) should fail", i, tc.desc)
		} else {
			if tc.deliverFailed {
				require.NotZero(t, res.TxResult.Code, "test case %v(%v) should fail", i, tc.desc)
			} else {
				require.Zero(t, res.TxResult.Code, "test case %v(%v) should not fail: %v", i, tc.desc, res.TxResult.Log)

				// Check status of the internal transaction
				log := types.ResultLogFromTMLog(res.TxResult.Log)
				if log.InternalHash != nil {
					internalTC := tc.msgInfo.(MultisigInfo).InternalTx
					// wait for 5 blocks until we can find the internal tx by its hash
					found := false
					for i := 0; i < 3; i++ {
						time.Sleep(1 * time.Second)
						WaitForNextHeightTM(tPort)
						internalTx := Tx(log.InternalHash)
						if internalTx != nil {
							found = true
							if internalTC.deliverFailed {
								require.NotZero(t, internalTx.TxResult.Code, "test case %v(%v) should fail", i, internalTC.desc)
							} else {
								require.Zero(t, internalTx.TxResult.Code, "test case %v(%v) should not fail: %v", i, internalTC.desc, internalTx.TxResult.Log)
							}
							break
						}
					}

					if internalTC.checkFailed {
						require.False(t, found)
					} else {
						require.True(t, found, "Unable to find internal transaction: %v", tc.desc)
					}
				}
			}
		}
	}

	WaitForNextHeightTM(tPort)

	// check account sequences to be increased properly
	for name, seq1 := range seqs {
		seq2 := AccSequence(tKeys[name].addrStr)
		assert.Equal(t, seq1, seq2)
	}

	acc2 := Account(tKeys["alice"].addrStr)
	bal2 := acc2.GetCoins()[0]
	diff := bal1.Sub(bal2)

	total := totalAmtAlice.Add(totalFeeAlice)
	require.True(t, diff.IsEqual(total))

	// Check all the fees are distributed to the validator
	stdOut, _, _ = utils.RunProcess("", "mxwcli", []string{"query", "distribution", "validator-outstanding-rewards", "mxwvaloper1rjgjjkkjqtd676ydahysmnfsg0v4yvwfp2n965"})
	s = strings.TrimSpace(strings.Replace(stdOut[24:], "\"", "", -1))
	dec, _ := sdkTypes.NewDecFromStr(s)
	dec = dec.Sub(initialRewards)
	distributedFee := sdkTypes.NewCoin("cin", dec.RoundInt())
	require.True(t, totalFee.IsEqual(distributedFee), "The validator rewards is not matched with accumulated fee, It should be %v but got %v", totalFee.Amount, distributedFee.Amount)
}

func makeMsg(t *testing.T, msgType string, signer string, msgInfo interface{}) sdkTypes.Msg {
	var msg sdkTypes.Msg

	switch msgType {
	case "bank":
		{
			i := msgInfo.(bankInfo)
			amt, err := types.ParseCoins(i.amount)
			assert.NoError(t, err, i.amount)

			msg = bank.NewMsgSend(tKeys[i.from].addr, tKeys[i.to].addr, amt)
		}
	case "kyc":
		{
			i := msgInfo.(kycInfo)

			switch i.action {
			case "whitelist":
				msg = makeKycWhitelistMsg(t, i.authorised, i.issuer, i.provider, i.from, i.signer, i.data, i.nonce)

			case "revokeWhitelist":
				msg = makeKycRevokeWhitelistMsg(t, i.authorised, i.issuer, i.provider, i.from)
			}
		}

	case "fee":
		{
			i := msgInfo.(feeInfo)

			msg = makeFeeMsg(t, i.function, i.name, i.assignee, i.multiplier, i.min, i.max, i.percentage, i.issuer)
		}

	case "nameservice":
		{
			i := msgInfo.(NameServiceInfo)

			switch i.Action {
			case "create":
				msg = makeCreateNameServiceMsg(t, i.Name, i.From, i.ApplicationFee, i.FeeCollector)
			case "approve":
				msg = setStatusAliasMsg(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Name, "APPROVE")
			case "reject":
				msg = setStatusAliasMsg(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Name, "REJECT")
			case "revoke":
				msg = setStatusAliasMsg(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Name, "REVOKE")
			}

		}
	case "token":
		{
			i := msgInfo.(TokenInfo)

			switch i.Action {
			case "create":
				msg = makeCreateFungibleTokenMsg(t, i.Name, i.Symbol, i.Metadata, i.Owner, i.MaxSupply, i.ApplicationFee, i.FeeCollector, i.Decimals, i.FixedSupply)
			case "approve":
				msg = makeApproveFungibleTokenMsg(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, "APPROVE", i.Burnable, i.FeeSettingName)
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
				msg = makeFreezeFungibleTokenMsg(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.Burnable)
			case "unfreeze":
				msg = makeUnfreezeFungibleTokenMsg(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.Burnable)
			case "verify-transfer-tokenOwnership":
				msg = makeVerifyTransferTokenOwnership(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.VerifyTransferTokenOwnership)
			}
		}
	case "nonFungibleToken":
		{
			i := msgInfo.(NonFungibleTokenInfo)

			switch i.Action {
			case "create":
				msg = makeCreateNonFungibleTokenMsg(t, i.Name, i.Symbol, i.TokenMetadata, i.Owner, i.ApplicationFee, i.FeeCollector)
			case "approve":
				msg = makeApproveNonFungibleTokenMsg(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, "APPROVE", i.FeeSettingName, i.MintLimit, i.TransferLimit, i.EndorserList, i.Burnable, i.Modifiable, i.Public, i.EndorserListLimit)
			case "reject":
				msg = makeRejectNonFungibleTokenMsg(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, "REJECT")
			case "transfer-item":
				msg = makeTransferNonFungibleTokenMsg(t, i.Owner, i.NewOwner, i.Symbol, i.ItemID)
			case "mint-item":
				msg = makeMintNonFungibleTokenMsg(t, i.Owner, i.NewOwner, i.Symbol, i.ItemID, i.Properties, i.Metadata)
			case "burn-item":
				msg = makeBurnNonFungibleTokenMsg(t, i.Owner, i.Symbol, i.ItemID)
			case "transfer-token-ownership":
				msg = makeTransferNonFungibleTokenOwnershipMsg(t, i.Owner, i.NewOwner, i.Symbol)
			case "accept-token-ownership":
				msg = makeAcceptNonFungibleTokenOwnershipMsg(t, i.NewOwner, i.Symbol)
			case "verify-transfer-token-ownership":
				msg = makeVerifyTransferNonFungibleTokenOwnershipMsg(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.VerifyTransferTokenOwnership, i.Burnable, i.Modifiable, i.Public)
			case "reject-transfer-token-ownership":
				msg = makeRejectTransferTokenOwnershipMsg(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, "REJECT_TRANSFER_TOKEN_OWNERSHIP")
			case "freeze-item":
				msg = makeFreezeNonFungibleItemMsg(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.ItemID)
			case "unfreeze-item":
				msg = makeUnfreezeNonFungibleItemMsg(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.ItemID)
			case "freeze":
				msg = makeFreezeNonFungibleTokenMsg(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.Burnable, i.Modifiable, i.Public)
			case "unfreeze":
				msg = makeUnfreezeNonFungibleTokenMsg(t, signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.Burnable, i.Modifiable, i.Public)
			case "endorsement-item":
				msg = makeEndorsementMsg(t, signer, i.Owner, i.Symbol, i.ItemID, i.Metadata)
			case "update-item-metadata":
				msg = makeUpdateItemMetadataMsg(t, i.Symbol, i.Owner, i.ItemID, i.Metadata)
			case "update-nft-metadata":
				msg = makeUpdateNFTMetadataMsg(t, i.Symbol, i.Owner, i.TokenMetadata)
			case "update-nft-endorserlist":
				msg = makeUpdateNFTEndorserListMsg(t, signer, i.Owner, i.Symbol, i.EndorserList)
			}
		}
	case "maintenance":
		{
			i := msgInfo.(MaintenanceInfo)
			msg = makeMaintenanceMsg(t, i.Action, i.Title, i.Description, i.ProposalType, i.AuthorisedAddress, i.IssuerAddress, i.ProviderAddress, i.Proposer, i.ValidatorPubKey, i.FeeCollector)
		}

	case "maintenance-cast-action":
		{
			i := msgInfo.(CastAction)
			msg = makeCastActionMsg(t, i.Action, i.Caster, i.ProposalId)
		}
	case "multiSig":
		{
			i := msgInfo.(MultisigInfo)

			switch i.Action {
			case "create":
				msg = makeCreateMultiSigAccountMsg(t, i.Owner, i.Threshold, i.Signers)
			case "update":
				msg = makeUpdateMultiSigAccountMsg(t, i.Owner, i.GroupAddress, i.Threshold, i.Signers)
			case "transfer-ownership":
				msg = makeTransferMultiSigOwnerMsg(t, i.GroupAddress, i.NewOwner, i.Owner)
			case "create-internal-tx":
				msg = createInternalTx(t, signer, i.GroupAddress, i.TxID, i.InternalTx)
			case "multiSig-sign-tx":
				msg = makeSignMultiSigTxMsg(t, signer, i.GroupAddress, i.TxID)
			case "multiSig-delete-tx":
				msg = makeDeleteMultiSigTxMsg(t, i.GroupAddress, i.TxID, i.Owner)

			}
		}
	}

	return msg
}

// for most of transactions, sender is same as signer.
// only for multi-sig transactions sender and signer are different.
func makeSignedTx(t *testing.T, sender, signer string, seq, gas uint64, fees sdkTypes.Coins, memo string, msg sdkTypes.Msg) (sdkAuth.StdTx, []byte) {
	acc := Account(tKeys[sender].addrStr)
	require.NotNil(t, acc, "alias:%s", sender)

	if !acc.IsMultiSig() {
		seq = acc.GetSequence()
	}

	signMsg := authTypes.StdSignMsg{
		AccountNumber: acc.GetAccountNumber(),
		ChainID:       "maxonrow-chain",
		Fee:           authTypes.NewStdFee(gas, fees),
		Memo:          memo,
		Msgs:          []sdkTypes.Msg{msg},
		Sequence:      seq,
	}

	signBz, signBzErr := tCdc.MarshalJSON(signMsg)
	if signBzErr != nil {
		panic(signBzErr)
	}

	sig, err := tKeys[signer].priv.Sign(sdkTypes.MustSortJSON(signBz))
	if err != nil {
		panic(err)
	}

	pub := tKeys[signer].priv.PubKey()
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
