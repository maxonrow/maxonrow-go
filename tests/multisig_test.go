package tests

import (
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/maxonrow/maxonrow-go/types"
	multisig "github.com/maxonrow/maxonrow-go/x/auth"
	"github.com/stretchr/testify/require"
)

type MultisigInfo struct {
	Action       string
	Owner        string
	NewOwner     string
	Threshold    int
	Signers      []string
	GroupAddress string
	TxID         uint64
	InternalTx   *testCase
}

func mustGetAccAddressFromBech32(bech32 string) sdkTypes.AccAddress {
	addr, _ := sdkTypes.AccAddressFromBech32(bech32)
	return addr
}

func makeMultisigTxs() []*testCase {

	// Group addresses:
	// You can generate group address via `mxwcli` command. ex.
	// `mxwcli keys multisig-address "mxw1ld3stcsk5l8xjngw2ucuazux895rk2hxve69gr" 1`
	tKeys["grp-addr-1"] = &keyInfo{
		mustGetAccAddressFromBech32("mxw1n94cgykexyzjuxmt97eaz2q2jecak2n84my0ny"), nil, nil, "mxw1n94cgykexyzjuxmt97eaz2q2jecak2n84my0ny",
	}
	// `mxwcli keys multisig-address "mxw1ld3stcsk5l8xjngw2ucuazux895rk2hxve69gr" 2`
	tKeys["grp-addr-2"] = &keyInfo{
		mustGetAccAddressFromBech32("mxw1gs8fq6sd5nd4vppnancpkjjh0gycdfr2g9dw0f"), nil, nil, "mxw1gs8fq6sd5nd4vppnancpkjjh0gycdfr2g9dw0f",
	}
	// `mxwcli keys multisig-address "mxw1ld3stcsk5l8xjngw2ucuazux895rk2hxve69gr" 3`
	tKeys["grp-addr-3"] = &keyInfo{
		mustGetAccAddressFromBech32("mxw1je73yfjvpmms68jswd42ceq6k8dl92uz45qztp"), nil, nil, "mxw1je73yfjvpmms68jswd42ceq6k8dl92uz45qztp",
	}
	// not exist account
	// `mxwcli keys multisig-address "mxw1ld3stcsk5l8xjngw2ucuazux895rk2hxve69gr" 4`
	tKeys["grp-addr-4"] = &keyInfo{
		mustGetAccAddressFromBech32("mxw1nnem49mcz6532fhn69nrwmp7p7kl65s397hvs6"), nil, nil, "mxw1nnem49mcz6532fhn69nrwmp7p7kl65s397hvs6",
	}

	internalTx_1 := &testCase{"bank", true, true, "sending 1 cin  ", "multisig-acc-1", "800400000cin", 0, bankInfo{"mostafa", "bob", "1cin"}, "tx1", nil}
	internalTx_2 := &testCase{"bank", true, true, "sending 1 cin  ", "multisig-acc-1", "800400000cin", 0, bankInfo{"grp-addr-3", "bob", "1cin"}, "tx2", nil}
	internalTx_3 := &testCase{"bank", false, false, "sending 1 cin", "multisig-acc-1", "800400000cin", 0, bankInfo{"grp-addr-1", "bob", "1cin"}, "tx3", nil}
	internalTx_4 := &testCase{"bank", false, false, "sending 1 cin", "multisig-acc-2", "800400000cin", 0, bankInfo{"grp-addr-2", "bob", "1cin"}, "tx4", nil}
	internalTx_5 := &testCase{"bank", false, false, "sending 1 cin", "multisig-acc-3", "800400000cin", 0, bankInfo{"grp-addr-2", "bob", "1cin"}, "tx5", nil}

	tcs := []*testCase{

		//create MultiSig Account
		{"multiSig", false, false, "Create MultiSig Account1 - Happy Path - commit ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 1, []string{"multisig-acc-1"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		{"multiSig", false, false, "Create MultiSig Account2- Happy Path - commit  ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 2, []string{"multisig-acc-2", "multisig-acc-3"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		{"multiSig", false, false, "Create MultiSig Account3 - Happy Path - commit ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 2, []string{"multisig-acc-2", "multisig-acc-3", "multisig-acc-4"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		{"multiSig", true, true, "Create MultiSig Account - non-kyc               ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 2, []string{"multisig-acc-1", "multisig-acc-no-kyc"}, "", 0, nil}, "", nil},

		{"bank", false, false, "top-up Multisig Group-address1 - commit", "multisig-acc-1", "800400000cin", 0, bankInfo{"multisig-acc-1", "grp-addr-1", "10000000000cin"}, "MEMO : top-up account", nil},
		{"bank", false, false, "top-up Multisig Group-address2 - commit", "multisig-acc-2", "800400000cin", 0, bankInfo{"multisig-acc-2", "grp-addr-2", "10000000000cin"}, "MEMO : top-up account", nil},
		{"bank", false, false, "top-up Multisig Group-address3 - commit", "multisig-acc-3", "800400000cin", 0, bankInfo{"multisig-acc-3", "grp-addr-3", "10000000000cin"}, "MEMO : top-up account", nil},

		{"multiSig", true, true, "MultiSig-create-tx-bank - Invalid sequence", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 5, internalTx_3}, "MEMO : MultiSig-create-tx-bank", nil},

		// multiSig-create-tx-bank with one signer, should broadcast immediately
		{"multiSig", true, true, "MultiSig-create-tx-bank - wrong signer       ", "multisig-acc-4", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx_1}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank - invalid sender     ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx_1}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank - invalid sender2           ", "mostafa", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx_1}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank - invalid tx-id      ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 1, internalTx_1}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank - Invalid internal_tx", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx_1}, "MEMO : MultiSig-create-tx-bank", nil},
		// this case will print an error messsage which can't find pending tx.
		{"multiSig", true, true, "MultiSig-create-tx-bank - Invalid internal_tx", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx_2}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", false, false, "MultiSig-create-tx-bank - Happy path        ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx_3}, "MEMO : MultiSig-create-tx-bank", nil},

		//
		{"multiSig", true, true, "MultiSig-create-tx-bank - counter+5            ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 5, internalTx_4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank - counter+1            ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 1, internalTx_4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", false, false, "MultiSig-create-tx-bank - Happy Path          ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx_4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank - Resubmit fail        ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx_4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank - Resubmit counter+2   ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx_4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", false, false, "MultiSig-create-tx-bank - Resubmit counter+1  ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 1, internalTx_4}, "MEMO : MultiSig-create-tx-bank", nil},

		{"multiSig", true, true, "MultiSig-sign-tx-bank - wrong signer      ", "multisig-acc-4", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx_4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-sign-tx-bank - owner can't sign  ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx_4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-sign-tx-bank - wrong sender             ", "mostafa", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx_4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", false, false, "MultiSig-sign-tx-bank - Happy Path-commit", "multisig-acc-3", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx_4}, "MEMO : MultiSig-sign-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-sign-tx-bank - resubmit          ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx_4}, "MEMO : MultiSig-sign-tx-bank", nil},

		// Note: msg signed by acc-2, internal_tx signed by acc-3. acc_2 sends both acc_2 and acc_3 signatures
		{"multiSig", true, true, "MultiSig-create-tx-bank2 - unknown signer                 ", "multisig-acc-4", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx_5}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank2 - owner can't sign               ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx_5}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank2 - Invalid signer                        ", "mostafa", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx_5}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank2 - Owner can't create internal-tx ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx_5}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", false, false, "MultiSig-create-tx-bank2 - Happy path                    ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx_5}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank2 - Resubmit with same tx_id       ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx_5}, "MEMO : MultiSig-create-tx-bank", nil},

		{"multiSig", false, false, "MultiSig-sign-tx-bank - Happy Path", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx_5}, "MEMO : MultiSig-sign-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-sign-tx-bank - resubmit   ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx_5}, "MEMO : MultiSig-sign-tx-bank", nil},
	}

	return tcs
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

func createInternalTx(t *testing.T, sender, groupAddress string, counter uint64, internalTxTemplate *testCase) sdkTypes.Msg {
	senderAddr := tKeys[sender].addr
	groupAddr := tKeys[groupAddress].addr

	internalTxMsg := makeMsg(t, internalTxTemplate.msgType, internalTxTemplate.signer, internalTxTemplate.msgInfo)
	fees, _ := types.ParseCoins(internalTxTemplate.fees)
	internalTx, _ := makeSignedTx(t, groupAddress, internalTxTemplate.signer, counter, internalTxTemplate.gas, fees, internalTxTemplate.memo, internalTxMsg)

	msgCreateMultiSigTx := multisig.NewMsgCreateMultiSigTx(groupAddr, internalTx, senderAddr)
	return msgCreateMultiSigTx

}

//moduel of TransferMultiSigOwner
func makeTransferMultiSigOwnerMsg(t *testing.T, groupAddress string, newOwner string, owner string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr
	groupAddr, _ := sdkTypes.AccAddressFromBech32(groupAddress)

	msgTransferMultiSigOwnerPayload := multisig.NewMsgTransferMultiSigOwner(groupAddr, newOwnerAddr, ownerAddr)

	return msgTransferMultiSigOwnerPayload
}

//makeDeleteMultiSigTxMsg
func makeDeleteMultiSigTxMsg(t *testing.T, groupAddress string, txID uint64, senderAddress string) sdkTypes.Msg {

	senderAddr := tKeys[senderAddress].addr
	groupAddr, _ := sdkTypes.AccAddressFromBech32(groupAddress)

	msgDeleteMultiSigTx := multisig.NewMsgDeleteMultiSigTx(groupAddr, txID, senderAddr)
	return msgDeleteMultiSigTx
}

func makeSignMultiSigTxMsg(t *testing.T, signer, groupAddress string, txID uint64) sdkTypes.Msg {
	groupAcc := Account(tKeys[groupAddress].addrStr)
	groupMultisig := groupAcc.GetMultiSig()
	require.NotNil(t, groupMultisig)
	ptx := groupMultisig.GetPendingTx(txID)
	require.NotNil(t, ptx)
	internalTx := ptx.GetTx().(sdkAuth.StdTx)
	require.NotNil(t, internalTx)

	signedTx, _ := makeSignedTx(t, groupAddress, signer, txID, internalTx.Fee.Gas, internalTx.Fee.Amount, internalTx.Memo, internalTx.Msgs[0])

	msgSignMultiSigTx := multisig.NewMsgSignMultiSigTx(groupAcc.Address, txID, signedTx.Signatures[0], tKeys[signer].addr)
	return msgSignMultiSigTx
}
