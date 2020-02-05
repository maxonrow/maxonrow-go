package tests

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/maxonrow/maxonrow-go/types"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth"
	multisig "github.com/maxonrow/maxonrow-go/x/auth"
	"github.com/maxonrow/maxonrow-go/x/bank"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MultisigInfo struct {
	Action         string
	ApplicationFee string
	FeeCollector   string
	Owner          string
	NewOwner       string
	Threshold      int
	Signers        []string
	GroupAddress   string
	TransactionId  string
}

func MakeMultisigAccsTxs(t *testing.T) {
	val1 := Validator(tValidator)
	fmt.Println(val1)

	_, err := time.ParseDuration("60s")
	if err != nil {
		fmt.Println("timeout", err)
	}

	tcs := []*testCase{

		// assign zero fee to an account
		{"kyc", false, false, "Doing kyc - mostafa - commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "mostafa", "mostafa", "testKyc123456789", "0"}, "MEMO : Doing kyc - mostafa - commit", nil},

		//goh123 - prepare for MultiSig module
		{"kyc", false, false, "Doing kyc - acc-21 - commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "acc-21", "acc-21", "testKyc122222222", "0"}, "MEMO : Doing kyc - acc-21 - commit", nil},
		{"kyc", false, false, "Doing kyc - acc-23 - commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "acc-23", "acc-23", "testKyc111111111", "0"}, "MEMO : Doing kyc - acc-23 - commit", nil},
		{"kyc", false, false, "Doing kyc - acc-24 - commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "acc-24", "acc-24", "testKyc133333333", "0"}, "MEMO : Doing kyc - acc-24 - commit", nil},
		{"bank", false, false, "sending 10000000000 cin", "alice", "200000000cin", 0, bankInfo{"alice", "acc-21", "10000000000cin"}, "MEMO : alice sending 10000000000cin to acc-21", nil},
		{"bank", false, false, "sending 10000000000 cin", "alice", "200000000cin", 0, bankInfo{"alice", "acc-23", "10000000000cin"}, "MEMO : alice sending 10000000000cin to acc-23", nil},
		{"bank", false, false, "sending 10000000000 cin", "alice", "200000000cin", 0, bankInfo{"alice", "acc-24", "10000000000cin"}, "MEMO : alice sending 10000000000cin to acc-24", nil},

		//create MultiSig Account
		{"multiSig", false, false, "Create MultiSig Account - Happy Path", "acc-21", "200000000cin", 0, MultisigInfo{"create", "10000000", "mostafa", "acc-21", "", 2, []string{"acc-21", "acc-24"}, "", ""}, "MEMO : Create MultiSig Account - Happy Path", nil},

		//NOTE : THIS for GroupAddress send BankTx
		{"bank-multisig", false, false, "BankTx sending 70000000000cin from acc-alice to Multisig Group-address", "alice", "200000000cin", 0, bankInfo{"alice", "mxw14fr3w8ffacdtkn6cmeg2ndpe7lxdzwt453crce", "10000000000cin"}, "MEMO : alice sending 70000000000cin to multisig-GroupAddress", nil},

		// multiSig-create-tx-bank
		{"multiSig-create-tx-bank", false, false, "MultiSig-create-tx-bank - Happy Path", "acc-21", "200000000cin", 0, bankInfo{"mxw14fr3w8ffacdtkn6cmeg2ndpe7lxdzwt453crce", "acc-24", "20cin"}, "MEMO : MultiSig-create-tx-bank - Happy Path", nil}, // OK

		// // multiSig-sign-tx-bank
		// {"multiSig", false, false, "MultiSig-sign-tx-bank - Happy Path", "acc-24", "200000000cin", 0, MultisigInfo{"multiSig-sign-tx-bank", "10000000", "mostafa", "acc-24", "", 2, []string{"acc-21", "acc-24"}, "mxw14fr3w8ffacdtkn6cmeg2ndpe7lxdzwt453crce", "1"}, "MEMO : MultiSig-sign-tx-bank - Happy Path", nil}, // OK

	}

	var totalFee = sdkTypes.NewInt64Coin("cin", 0)
	var totalAmt = sdkTypes.NewInt64Coin("cin", 0)

	for n, tc := range tcs {

		fees, err := types.ParseCoins(tc.fees)
		assert.NoError(t, err, tc.fees)

		var msg sdkTypes.Msg
		switch tc.msgType {
		case "bank-multisig":
			{
				i := tc.msgInfo.(bankInfo)
				amt, err := types.ParseCoins(i.amount)
				assert.NoError(t, err, i.amount)

				groupAddr, _ := sdkTypes.AccAddressFromBech32(i.to)
				msg = bank.NewMsgSend(tKeys[i.from].addr, groupAddr, amt)

				if !tc.deliverFailed {
					if i.from == "alice" {
						totalFee = totalFee.Add(fees[0])
						totalAmt = totalAmt.Add(amt[0])
					}
				}
			}
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

		case "multiSig":
			{
				i := tc.msgInfo.(MultisigInfo)

				switch i.Action {
				case "create":
					msg = makeCreateMultiSigAccountMsg(t, i.Owner, i.Threshold, i.Signers)
				case "update":
					msg = makeUpdateMultiSigAccountMsg(t, i.Owner, i.GroupAddress, i.Threshold, i.Signers)
				case "transfer-ownership":
					msg = makeTransferMultiSigOwnerMsg(t, i.GroupAddress, i.NewOwner, i.Owner)
				case "multiSig-sign-tx-bank":
					msg = makeSignMultiSigTxMsg(t, i.GroupAddress, i.TransactionId, i.Owner)
				case "multiSig-delete-tx-bank":
					msg = makeDeleteMultiSigTxMsg(t, i.GroupAddress, i.TransactionId, i.Owner)

				}
			}
		case "multiSig-create-tx-bank":
			{
				i := tc.msgInfo.(bankInfo)
				msg = makeBanksendMsg(t, i.from, i.to, i.amount)
			}
		}

		isMultiSig := checkMultiSig(tc.msgType)
		tx, bz := MakeMultiSigOrSingleSigTx(t, tc.signer, tc.gas, fees, tc.memo, msg, isMultiSig)

		txb, err := tCdc.MarshalJSON(tx)
		assert.NoError(t, err)
		fmt.Printf("============\nBroadcasting %d tx (%s): %s\n", n+1, tc.desc, string(txb))

		res := BroadcastTxCommit(bz)
		tc.hash = res.Hash.Bytes()

		if !tc.checkFailed {
			tKeys[tc.signer].seq++
			require.Zero(t, res.CheckTx.Code, "test case %v(%v) check failed--111: %v", n+1, tc.desc, res.CheckTx.GetLog())

			if !tc.deliverFailed {
				require.Zero(t, res.DeliverTx.Code, "test case %v(%v) deliver failed - 001: %v", n+1, tc.desc, res.DeliverTx.Log)
			} else {
				require.NotZero(t, res.DeliverTx.Code, "test case %v(%v) deliver failed - 002: %v", n+1, tc.desc, res.DeliverTx.Log)
			}

		} else {
			require.NotZero(t, res.CheckTx.Code, "test case %v(%v) check failed--222: %v", n+1, tc.desc, res.CheckTx.Log)
		}

		if strings.Contains(tc.desc, "commit") {
			WaitForNextHeightTM(tPort)
		}
	}

	WaitForNextHeightTM(tPort)

	/// check results
	for i, tc := range tcs {
		res := Tx(tc.hash)
		_ = Tx(tc.hash)
		if tc.checkFailed {
			assert.Nil(t, res, "test case %v(%v) failed", i, tc.desc)
		} else {
			if tc.deliverFailed {
				assert.NotZero(t, res.TxResult.Code, "test case %v(%v) failed", i, tc.desc)

			} else {
				//fmt.Printf("get Group-Address : %v \n", res.TxResult.Log)
				assert.Zero(t, res.TxResult.Code, "test case %v(%v) failed: %v", i, tc.desc, res.TxResult.Log)
			}
		}
	}

	val2 := Validator(tValidator)
	fmt.Println(val2)
}

func checkMultiSig(msgType string) bool {

	multiSigArray := [1]string{"multiSig-create-tx-bank"}

	for _, item := range multiSigArray {
		if item == msgType {
			return true
		}
	}

	return false
}

//
func MakeMultiSigOrSingleSigTx(t *testing.T, signer string, gas uint64, fees sdkTypes.Coins, memo string, msg sdkTypes.Msg,
	isMultiSig bool) (authTypes.StdTx, []byte) {

	//1. SingleSig
	if isMultiSig == false {
		acc := Account(tKeys[signer].addrStr)
		require.NotNil(t, acc)

		signMsg := authTypes.StdSignMsg{
			AccountNumber: acc.GetAccountNumber(),
			ChainID:       "maxonrow-chain",
			Fee:           authTypes.NewStdFee(gas, fees),
			Memo:          memo,
			Msgs:          []sdkTypes.Msg{msg},
			Sequence:      tKeys[signer].seq,
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
		stdSig := authTypes.StdSignature{
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

	acc := Account(tKeys[signer].addrStr)
	require.NotNil(t, acc)

	// 2. MultiSig
	groupAddress := "mxw14fr3w8ffacdtkn6cmeg2ndpe7lxdzwt453crce"

	pendingTx := authTypes.NewStdTx([]sdkTypes.Msg{msg}, authTypes.NewStdFee(gas, fees), nil, memo) // no need signatures
	msgCreateMultiSigTx := makeCreateMultiSigTxMsg(t, groupAddress, pendingTx, signer)

	signMsg := authTypes.StdSignMsg{
		AccountNumber: acc.GetAccountNumber(),
		ChainID:       "maxonrow-chain",
		Fee:           authTypes.NewStdFee(gas, fees),
		Memo:          memo,
		Msgs:          []sdkTypes.Msg{msgCreateMultiSigTx},
		Sequence:      1,
	}

	signBz := sdkTypes.MustSortJSON(tCdc.MustMarshalJSON(signMsg))

	sig, err := tKeys[signer].priv.Sign(sdkTypes.MustSortJSON(signBz))
	if err != nil {
		panic(err)
	}

	pub := tKeys[signer].priv.PubKey()
	stdSig := authTypes.StdSignature{
		PubKey:    pub,
		Signature: sig,
	}

	stdMultiSigtx := authTypes.NewStdTx([]sdkTypes.Msg{msgCreateMultiSigTx}, signMsg.Fee, []authTypes.StdSignature{stdSig}, memo)
	bz, err := tCdc.MarshalBinaryLengthPrefixed(stdMultiSigtx)

	if err != nil {
		panic(err)
	}

	return stdMultiSigtx, bz

}

//module of CreateMultiSigAccount
func makeCreateMultiSigAccountMsg(t *testing.T, owner string, threshold int, signers []string) sdkTypes.Msg {

	// create new multisig account
	ownerAddr := tKeys[owner].addr

	// convert item
	var signersAddr []sdkTypes.AccAddress
	for i := 0; i < len(signers); i++ {
		signerStr := signers[i]
		ownerAddr := tKeys[signerStr].addr

		signersAddr = append(signersAddr, ownerAddr)
	}

	msgCreateMultiSigAccountPayload := multisig.NewMsgCreateMultiSigAccount(ownerAddr, threshold, signersAddr)

	return msgCreateMultiSigAccountPayload
}

//module of UpdateMultiSigAccount
func makeUpdateMultiSigAccountMsg(t *testing.T, owner string, groupAddress string, threshold int, signers []string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	groupAddr, _ := sdkTypes.AccAddressFromBech32(groupAddress)

	var signersAddr []sdkTypes.AccAddress
	for i := 0; i < len(signers); i++ {
		signerStr := signers[i]
		ownerAddr := tKeys[signerStr].addr

		signersAddr = append(signersAddr, ownerAddr)
	}

	msgUpdateMultiSigAccountPayload := multisig.NewMsgUpdateMultiSigAccount(ownerAddr, groupAddr, threshold, signersAddr)
	return msgUpdateMultiSigAccountPayload

}

//moduel of TransferMultiSigOwner
func makeTransferMultiSigOwnerMsg(t *testing.T, groupAddress string, newOwner string, owner string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr
	groupAddr, _ := sdkTypes.AccAddressFromBech32(groupAddress)

	msgTransferMultiSigOwnerPayload := multisig.NewMsgTransferMultiSigOwner(groupAddr, newOwnerAddr, ownerAddr)

	return msgTransferMultiSigOwnerPayload
}

//makeCreateMultiSigTxMsg
func makeCreateMultiSigTxMsg(t *testing.T, groupAddress string, tx authTypes.StdTx, senderAddress string) sdkTypes.Msg {

	senderAddr := tKeys[senderAddress].addr
	groupAddr, _ := sdkTypes.AccAddressFromBech32(groupAddress)

	msgCreateMultiSigTx := multisig.NewMsgCreateMultiSigTx(groupAddr, tx, senderAddr)
	return msgCreateMultiSigTx
}

//makeSignMultiSigTxMsg : as Acknowledgement
func makeSignMultiSigTxMsg(t *testing.T, groupAddress string, txCode string, senderAddress string) sdkTypes.Msg {

	senderAddr := tKeys[senderAddress].addr
	txID, _ := strconv.ParseUint(txCode, 10, 64)
	groupAddr, _ := sdkTypes.AccAddressFromBech32(groupAddress)

	msgSignMultiSigTx := multisig.NewMsgSignMultiSigTx(groupAddr, txID, senderAddr)
	return msgSignMultiSigTx
}

//makeDeleteMultiSigTxMsg
func makeDeleteMultiSigTxMsg(t *testing.T, groupAddress string, txCode string, senderAddress string) sdkTypes.Msg {

	senderAddr := tKeys[senderAddress].addr
	txID, _ := strconv.ParseUint(txCode, 10, 64)
	groupAddr, _ := sdkTypes.AccAddressFromBech32(groupAddress)

	msgDeleteMultiSigTx := multisig.NewMsgDeleteMultiSigTx(groupAddr, txID, senderAddr)
	return msgDeleteMultiSigTx
}

// Tx
//module of CreateMultiSigTx
func makeBanksendMsg(t *testing.T, from string, to string, amount string) sdkTypes.Msg {

	var msgBanksendPayload sdkTypes.Msg
	amt, err := types.ParseCoins(amount)
	assert.NoError(t, err, amount)

	fromAddr, _ := sdkTypes.AccAddressFromBech32(from)

	msgBanksendPayload = bank.NewMsgSend(fromAddr, tKeys[to].addr, amt)
	return msgBanksendPayload
}
