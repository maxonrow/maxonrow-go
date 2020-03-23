package tests

import (
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/maxonrow/maxonrow-go/types"
	multisig "github.com/maxonrow/maxonrow-go/x/auth"
	"github.com/stretchr/testify/assert"
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
		mustGetAccAddressFromBech32("mxw1z8r356ll7aum0530xve2upx74ed8ffavyxy503"), nil, nil, "mxw1z8r356ll7aum0530xve2upx74ed8ffavyxy503",
	}
	// `mxwcli keys multisig-address "mxw1ld3stcsk5l8xjngw2ucuazux895rk2hxve69gr" 2`
	tKeys["grp-addr-2"] = &keyInfo{
		mustGetAccAddressFromBech32("mxw1q6nmfejarl5e4xzceqxcygner7a6llgwnrdtl6"), nil, nil, "mxw1q6nmfejarl5e4xzceqxcygner7a6llgwnrdtl6",
	}
	// `mxwcli keys multisig-address "mxw1ld3stcsk5l8xjngw2ucuazux895rk2hxve69gr" 3`
	tKeys["grp-addr-3"] = &keyInfo{
		mustGetAccAddressFromBech32("mxw1hkm4p04nsmv9q0hg4m9eeuapfdr7n4rfl04vh9"), nil, nil, "mxw1hkm4p04nsmv9q0hg4m9eeuapfdr7n4rfl04vh9",
	}
	// not exist account
	// `mxwcli keys multisig-address "mxw1ld3stcsk5l8xjngw2ucuazux895rk2hxve69gr" 4`
	tKeys["grp-addr-4"] = &keyInfo{
		mustGetAccAddressFromBech32("mxw1szm87m362urkvj833jd7nekwdjh7s8p4q3f25f"), nil, nil, "mxw1szm87m362urkvj833jd7nekwdjh7s8p4q3f25f",
	}

	internalTx1 := &testCase{"bank", true, true, "sending 1 cin  ", "multisig-acc-1", "800400000cin", 0, bankInfo{"mostafa", "bob", "1cin"}, "tx1", nil}
	internalTx2 := &testCase{"bank", true, true, "sending 1 cin  ", "multisig-acc-1", "800400000cin", 0, bankInfo{"grp-addr-3", "bob", "1cin"}, "tx2", nil}
	internalTx3 := &testCase{"bank", false, false, "sending 1 cin", "multisig-acc-1", "800400000cin", 0, bankInfo{"grp-addr-1", "bob", "1cin"}, "tx3", nil}
	internalTx4 := &testCase{"bank", false, false, "sending 1 cin", "multisig-acc-2", "800400000cin", 0, bankInfo{"grp-addr-2", "bob", "1cin"}, "tx4", nil}
	internalTx5 := &testCase{"bank", false, false, "sending 1 cin", "multisig-acc-3", "800400000cin", 0, bankInfo{"grp-addr-2", "bob", "1cin"}, "tx5", nil}

	tcs := []*testCase{

		//create MultiSig Account1 : {"multisig-acc-1"}, owner=="multisig-acc-1"
		{"multiSig", false, false, "Create MultiSig Account1 - Happy Path - commit", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 1, []string{"multisig-acc-1"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		//create MultiSig Account2 : {"multisig-acc-2", "multisig-acc-3"}, owner=="multisig-acc-1"
		{"multiSig", false, false, "Create MultiSig Account2- Happy Path - commit ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 2, []string{"multisig-acc-2", "multisig-acc-3"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		//create MultiSig Account3 : {"multisig-acc-2", "multisig-acc-3", "multisig-acc-4"}, owner=="multisig-acc-1"
		{"multiSig", false, false, "Create MultiSig Account3 - Happy Path - commit", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 2, []string{"multisig-acc-2", "multisig-acc-3", "multisig-acc-4"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		{"multiSig", true, true, "Create MultiSig Account - non-kyc               ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 2, []string{"multisig-acc-1", "multisig-acc-no-kyc"}, "", 0, nil}, "", nil},

		{"bank", false, false, "top-up Multisig Group-address1 - commit", "multisig-acc-1", "800400000cin", 0, bankInfo{"multisig-acc-1", "grp-addr-1", "10000000000cin"}, "MEMO : top-up account", nil},
		{"bank", false, false, "top-up Multisig Group-address2 - commit", "multisig-acc-2", "800400000cin", 0, bankInfo{"multisig-acc-2", "grp-addr-2", "10000000000cin"}, "MEMO : top-up account", nil},
		{"bank", false, false, "top-up Multisig Group-address3 - commit", "multisig-acc-3", "800400000cin", 0, bankInfo{"multisig-acc-3", "grp-addr-3", "10000000000cin"}, "MEMO : top-up account", nil},

		//====================start : case-1.1
		//-- Using : 'MultiSig Account1' which owner=="multisig-acc-1"
		//-- Scenario : using 'grp-addr-1' which only with ONE signer {'MultiSig Account1'}, should broadcast immediately

		{"multiSig", true, true, "MultiSig-create-tx-bank - Invalid sequence   ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 5, internalTx3}, "MEMO : MultiSig-create-tx-bank", nil},

		// multiSig-create-tx-bank with one signer, should broadcast immediately
		{"multiSig", true, true, "MultiSig-create-tx-bank - wrong signer       ", "multisig-acc-4", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx1}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank - invalid sender     ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx1}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank - invalid sender2    ", "mostafa", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx1}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank - invalid tx-id      ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 1, internalTx1}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank - Invalid internal_tx", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx1}, "MEMO : MultiSig-create-tx-bank", nil},
		// this case will print an error messsage which can't find pending tx.
		{"multiSig", true, true, "MultiSig-create-tx-bank - Invalid internal_tx", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx2}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", false, false, "MultiSig-create-tx-bank - Happy path - commit", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx3}, "MEMO : MultiSig-create-tx-bank", nil},

		//Topic : create tx / sign tx
		//Remarks : Since with one signer, and should broadcast immediately,
		//          so this 'create tx' which include process of checkIsMetric will not left any pending Tx (as Broadcasted internal transaction successfully).
		//					Thus this case-1.1 should not include cases of [false, false] under 'Delete MultiSig Tx'.
		{"multiSig", true, true, "case-1.1-Sign MultiSig Tx - Errr, due to All signers must pass kyc.																							 ", "multisig-acc-no-kyc", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx3}, "MEMO : Sign MultiSig Tx", nil}, //ok-20200316
		{"multiSig", true, true, "case-1.1-Re-sign MultiSig Tx - Error for counter+0, due to already signed by multisig-acc-1 while create-tx-bank.", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx3}, "MEMO : Sign MultiSig Tx", nil}, //ok-20200316
		{"multiSig", false, false, "case-1.1-Create MultiSig Tx - submit counter+1.															                                   ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 1, internalTx3}, "MEMO : Create MultiSig Tx", nil}, //ok-20200316
		{"multiSig", true, true, "case-1.1-Re-create MultiSig Tx - submit counter+1 - Error, due to Re-create Tx with same sequence			           ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 1, internalTx3}, "MEMO : Create MultiSig Tx", nil}, //ok-20200316
		{"multiSig", true, true, "case-1.1-Re-sign MultiSig Tx - Error for counter+1, due to already signed by multisig-acc-1 while create-tx-bank.", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-1", 1, internalTx3}, "MEMO : Sign MultiSig Tx", nil}, //ok-20200316
		{"multiSig", false, false, "case-1.1-Create MultiSig Tx - submit counter+2.															                                   ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 2, internalTx3}, "MEMO : Create MultiSig Tx", nil}, //ok-20200316
		{"multiSig", true, true, "case-1.1-Re-create MultiSig Tx - submit counter+2 - Error, due to Re-create Tx with same sequence			           ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 2, internalTx3}, "MEMO : Create MultiSig Tx", nil}, //ok-20200316
		{"multiSig", true, true, "case-1.1-Re-sign MultiSig Tx - Error for counter+2, due to already signed by multisig-acc-1 while create-tx-bank.", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-1", 2, internalTx3}, "MEMO : Sign MultiSig Tx", nil}, //ok-20200316

		//Topic : delete tx - after 'Sign MultiSig Tx'
		{"multiSig", true, true, "case-1.1-Delete MultiSig Tx - Error, due to Group address invalid.                         ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-4", 0, internalTx3}, "MEMO : Delete MultiSig Tx", nil}, //ok-20200316
		{"multiSig", true, true, "case-1.1-Delete MultiSig Tx - Error, due to Owner address invalid.                         ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-2", "", 0, nil, "grp-addr-1", 0, internalTx3}, "MEMO : Delete MultiSig Tx", nil}, //ok-20200316
		{"multiSig", true, true, "case-1.1-Delete MultiSig Tx - Error, due to 'Pending tx is not found' which ID : 3.       ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-1", 3, internalTx3}, "MEMO : Delete MultiSig Tx", nil},  //ok-20200316

		//Topic : transfer ownership
		{"multiSig", true, true, "case-1.1-Transfer MultiSig Owner - Error, due to Group address invalid.                                                                ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-4", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil}, // ok-20200317
		{"multiSig", true, true, "case-1.1-Transfer MultiSig Owner - Error, due to Owner of group address invalid.                                                       ", "multisig-acc-3", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-3", "multisig-acc-1", 0, nil, "grp-addr-1", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil}, // ok-20200317
		{"multiSig", true, true, "case-1.1-Transfer MultiSig Owner - Error, due to without KYC																																					 ", "multisig-acc-no-kyc", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-no-kyc", 0, nil, "grp-addr-1", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
		{"multiSig", false, false, "case-1.1-Transfer MultiSig Owner - [from multisig-acc-1 to multisig-acc-2] - Happy Path - commit.                                    ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-1", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil}, // ok-20200317
		{"multiSig", true, true, "case-1.1-Re-transfer MultiSig Owner - Error, due to Owner of group address invalid as MultiSig-account already been transfer to others.", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-1", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil}, // ok-20200317

		//====================start : case-1.2
		//-- Scenario : using 'grp-addr-2' with 'internalTx4'

		{"multiSig", true, true, "MultiSig-create-tx-bank - counter+5           ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 5, internalTx4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank - counter+1           ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 1, internalTx4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", false, false, "MultiSig-create-tx-bank - Happy Path        ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank - Resubmit fail       ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank - Resubmit counter+2  ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", false, false, "MultiSig-create-tx-bank - Resubmit counter+1", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 1, internalTx4}, "MEMO : MultiSig-create-tx-bank", nil},

		{"multiSig", true, true, "MultiSig-sign-tx-bank - wrong signer       ", "multisig-acc-4", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-sign-tx-bank - owner can't sign   ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx4}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-sign-tx-bank - wrong sender       ", "mostafa", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx4}, "MEMO : MultiSig-sign-tx-bank", nil},
		{"multiSig", false, false, "MultiSig-sign-tx-bank - Happy Path-commit", "multisig-acc-3", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx4}, "MEMO : MultiSig-sign-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-sign-tx-bank - resubmit           ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx4}, "MEMO : MultiSig-sign-tx-bank", nil},

		//====================start : case-2
		//-- Using : 'MultiSig Account2' which owner=="multisig-acc-1", signer-list == {"multisig-acc-2", "multisig-acc-3"}
		//-- Scenario : using 'grp-addr-2' which 'msg signed by acc-2, internal_tx signed by acc-3. acc_2 sends both acc_2 and acc_3 signatures'

		//Topic : update
		{"multiSig", true, true, "case-2-Update MultiSig Account - Error, due to Group address invalid.                     ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"update", "multisig-acc-1", "", 2, []string{"multisig-acc-2", "multisig-acc-3"}, "grp-addr-4", 0, nil}, "MEMO : Update MultiSig Account", nil}, // ok-20200316
		{"multiSig", true, true, "case-2-Update MultiSig Account - Error, due to MultiSig Account's 'Owner address invalid.'", "multisig-acc-4", "800400000cin", 0, MultisigInfo{"update", "multisig-acc-4", "", 2, []string{"multisig-acc-2", "multisig-acc-3"}, "grp-addr-2", 0, nil}, "MEMO : Update MultiSig Account", nil}, // ok-20200316

		{"multiSig", true, true, "MultiSig-create-tx-bank2 - unknown signer                 ", "multisig-acc-4", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx5}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank2 - owner can't sign               ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx5}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank2 - Invalid signer                 ", "mostafa", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx5}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank2 - Owner can't create internal-tx ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx5}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", false, false, "MultiSig-create-tx-bank2 - Happy path                   ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx5}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-create-tx-bank2 - Resubmit with same tx_id       ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 2, internalTx5}, "MEMO : MultiSig-create-tx-bank", nil},

		//Topic : delete tx - before 'Sign MultiSig Tx'
		{"multiSig", false, false, "case-2-Delete MultiSig Tx - Happy-path, due to the Valid PendingTx ID: 2. been found ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-2", 2, internalTx5}, "MEMO : Delete MultiSig Tx", nil},

		//Topic : create tx with nx sequence again
		{"multiSig", false, false, "case-2-Create MultiSig Tx - Happy path", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 3, internalTx5}, "MEMO : MultiSig-create-tx-bank", nil},

		{"multiSig", false, false, "MultiSig-sign-tx-bank - Happy Path", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 3, internalTx5}, "MEMO : MultiSig-sign-tx-bank", nil},
		{"multiSig", true, true, "MultiSig-sign-tx-bank - resubmit    ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 3, internalTx5}, "MEMO : MultiSig-sign-tx-bank", nil},

		//Topic : transfer ownership
		{"multiSig", true, true, "case-2-Transfer MultiSig Owner - Error, due to Group address invalid.                                                               ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-4", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},           // ok-20200317
		{"multiSig", true, true, "case-2-Transfer MultiSig Owner - Error, due to Owner of group address invalid.                                                      ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-2", "multisig-acc-1", 0, nil, "grp-addr-2", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},           // ok-20200317
		{"multiSig", true, true, "case-2-Transfer MultiSig Owner - Error, due to without KYC                                                                          ", "multisig-acc-no-kyc", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-no-kyc", 0, nil, "grp-addr-4", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
		{"multiSig", false, false, "case-2-Transfer MultiSig Owner - [from multisig-acc-1 to multisig-acc-2] Happy Path - commit.                                     ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-2", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},           // ok-20200317
		{"multiSig", true, true, "case-2-Re-transfer MultiSig Owner - Error, due to Owner of group address invalid [MultiSig-account already been transfer to others].", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-2", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},           // ok-20200317

		//Topic : delete tx - after 'Sign MultiSig Tx'
		{"multiSig", true, true, "case-2-Delete MultiSig Tx - Error, due to Group address invalid.                         ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-4", 3, internalTx5}, "MEMO : Delete MultiSig Tx", nil}, // ok-20200317
		{"multiSig", true, true, "case-2-Delete MultiSig Tx - Error, due to Only group account owner can remove pending tx.", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-2", 3, internalTx5}, "MEMO : Delete MultiSig Tx", nil}, // ok-20200317
		{"multiSig", true, true, "case-2-Delete MultiSig Tx - Error, due to 'Pending tx is not found' which ID : 9.        ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-2", "", 0, nil, "grp-addr-2", 9, internalTx5}, "MEMO : Delete MultiSig Tx", nil}, // ok-20200317

		//====================start : case-3
		//-- Using : 'MultiSig Account3' which owner=="multisig-acc-1", signer-list == {"multisig-acc-2", "multisig-acc-3", "multisig-acc-4"}
		//-- Scenario : using 'grp-addr-3' which 'Early with THREE signers "multisig-acc-2", "multisig-acc-3", "multisig-acc-4"}, then after update with THREE signers {"multisig-acc-1", "multisig-acc-3", "multisig-acc-4"}'

		//Topic : update
		{"multiSig", true, true, "case-3-Update MultiSig Account - Error, due to number of thresholds bigger than signers list       ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"update", "multisig-acc-1", "", 3, []string{"multisig-acc-2", "multisig-acc-3"}, "grp-addr-3", 0, nil}, "MEMO : Update MultiSig Account", nil}, // ok-20200317
		{"multiSig", false, false, "case-3-Update MultiSig Account - Happy Path - commit.                                 					 ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"update", "multisig-acc-1", "", 2, []string{"multisig-acc-1", "multisig-acc-3", "multisig-acc-4"}, "grp-addr-3", 0, nil}, "MEMO : Update MultiSig Account", nil}, // ok-20200317

		//Topic : create-tx without any sign-tx, but try to delete-tx later
		{"multiSig", false, false, "case-3-Create MultiSig Tx - submit counter+0 - Happy Path commit.                            ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-3", 0, internalTx2}, "MEMO : Create MultiSig Tx", nil},
		{"multiSig", true, true, "case-3-Re-create MultiSig Tx - submit counter+0 - Error, due to Re-create Tx with same sequence", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-3", 0, internalTx2}, "MEMO : Create MultiSig Tx", nil}, //ok-20200316

		//Topic : delete tx - before start 'Sign MultiSig Tx'
		{"multiSig", true, true, "case-3-Delete MultiSig Tx - Error, due to Group address invalid.													", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-4", 0, internalTx2}, "MEMO : Delete MultiSig Tx", nil}, // ok-20200317
		{"multiSig", true, true, "case-3-Delete MultiSig Tx - Error, due to Only group account owner can remove pending tx. ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-2", "", 0, nil, "grp-addr-3", 0, internalTx2}, "MEMO : Delete MultiSig Tx", nil}, // ok-20200317
		{"multiSig", true, true, "case-3-Delete MultiSig Tx - Error, due to 'Pending tx is not found' which ID : 1.         ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-3", 1, internalTx2}, "MEMO : Delete MultiSig Tx", nil},

		//Topic : create tx & sign tx
		//Remarks : Since after all the signers signed the tx, 'sign-tx' which include process of checkIsMetric will not left any pending Tx (as Broadcasted internal transaction successfully).
		//					Thus this case-2 should not include cases of [false, false] under 'Delete MultiSig Tx'.
		{"multiSig", true, true, "case-3-Create MultiSig Tx - Error, due to 'Internal transaction signature error' while submit counter+0.", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-3", 0, internalTx2}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", false, false, "case-3-Create MultiSig Tx - submit counter+1 - Happy Path commit.                            					", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-3", 1, internalTx2}, "MEMO : MultiSig-create-tx-bank", nil},
		{"multiSig", true, true, "case-3-Re-create MultiSig Tx - submit counter+1 - Error, due to Re-create Tx with same sequence					", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-3", 1, internalTx2}, "MEMO : MultiSig-create-tx-bank", nil}, //ok-20200316

		{"multiSig", true, true, "case-3-Sign MultiSig Tx - Error, due to All signers must pass kyc.                             ", "multisig-acc-no-kyc", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-3", 1, internalTx2}, "MEMO : Sign MultiSig Tx", nil}, // ok-20200317
		{"multiSig", true, true, "case-3-Sign MultiSig Tx - Error, due to Sender is not group account's signer.                  ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-3", 1, internalTx2}, "MEMO : Sign MultiSig Tx", nil},      // ok-20200317

		{"multiSig", true, true, "case-3-Sign MultiSig Tx - Error, due to already signed by multisig-acc-1                           																				", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-3", 1, internalTx2}, "MEMO : Sign MultiSig Tx", nil}, // ok-20200317
		{"multiSig", false, false, "case-3-Sign MultiSig Tx - submit counter+1 which signed by multisig-acc-3 - commit.              																				", "multisig-acc-3", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-3", 1, internalTx2}, "MEMO : Sign MultiSig Tx", nil}, // ok-20200317
		{"multiSig", true, true, "case-3-Re-sign MultiSig Tx - Error, due to counter+1 already signed by multisig-acc-3              																				", "multisig-acc-3", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-3", 1, internalTx2}, "MEMO : Sign MultiSig Tx", nil}, // ok-20200317
		{"multiSig", false, false, "case-3-Sign MultiSig Tx - submit counter+1 which signed by multisig-acc-4 - commit.              																			  ", "multisig-acc-4", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-3", 1, internalTx2}, "MEMO : Sign MultiSig Tx", nil}, // ok-20200317
		{"multiSig", true, true, "case-3-Re-sign MultiSig Tx - Error, due to counter+1 already signed by multisig-acc-4              																				", "multisig-acc-4", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-3", 1, internalTx2}, "MEMO : Sign MultiSig Tx", nil}, // ok-20200317
		{"multiSig", true, true, "case-3-Sign MultiSig Tx - Error, is Invalid as Sender multisig-acc-2 is not group account's signer after done the Update-MultiSig-Account.", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-3", 1, internalTx2}, "MEMO : Sign MultiSig Tx", nil}, // ok-20200317

		//Topic : transfer ownership
		{"multiSig", true, true, "case-3-Transfer MultiSig Owner - Error, due to Group address invalid.                                                               ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-4", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},           // ok-20200317
		{"multiSig", true, true, "case-3-Transfer MultiSig Owner - Error, due to Owner of group address invalid.                                                      ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-2", "multisig-acc-1", 0, nil, "grp-addr-3", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},           // ok-20200317
		{"multiSig", true, true, "case-3-Transfer MultiSig Owner - Error, due to without KYC                                                                          ", "multisig-acc-no-kyc", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-no-kyc", 0, nil, "grp-addr-4", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
		{"multiSig", false, false, "case-3-Transfer MultiSig Owner - [from multisig-acc-1 to multisig-acc-2] Happy Path - commit.                                     ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-3", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},           // ok-20200317
		{"multiSig", true, true, "case-3-Re-transfer MultiSig Owner - Error, due to Owner of group address invalid [MultiSig-account already been transfer to others].", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-3", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},           // ok-20200317

		// signer without through KYC
		{"multiSig", true, true, "Create MultiSig Account - Error, due to without KYC", "multisig-acc-no-kyc", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-no-kyc", "", 2, []string{"multisig-acc-1", "multisig-acc-no-kyc"}, "", 0, nil}, "MEMO : Create MultiSig Account", nil},
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
	groupAddr := tKeys[groupAddress].addr
	assert.NotNil(t, groupAddr)

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
	assert.NotNil(t, groupAddr)

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
	groupAddr := tKeys[groupAddress].addr
	assert.NotNil(t, groupAddr)

	msgTransferMultiSigOwnerPayload := multisig.NewMsgTransferMultiSigOwner(groupAddr, newOwnerAddr, ownerAddr)

	return msgTransferMultiSigOwnerPayload
}

//makeDeleteMultiSigTxMsg
func makeDeleteMultiSigTxMsg(t *testing.T, groupAddress string, txID uint64, senderAddress string) sdkTypes.Msg {

	senderAddr := tKeys[senderAddress].addr
	groupAddr := tKeys[groupAddress].addr
	assert.NotNil(t, groupAddr)

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
