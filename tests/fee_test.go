package tests

import (
	"testing"

	"github.com/maxonrow/maxonrow-go/x/fee"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

type feeInfo struct {
	function   string
	name       string
	assignee   string
	multiplier string
	min        string
	max        string
	percentage string
	issuer     string
}

func makeFeeTxs() []*testCase {

	tcs := []*testCase{

		{"fee", true, true, "Updating default fee settings-invalid signer", "alice", "0cin", 0, feeInfo{"sys-fee", "default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", true, true, "Updating default fee settings-wrong issuer", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "default", "", "", "400000000cin", "2000000000cin", "0.001", "alice"}, "", nil},
		{"fee", false, false, "Updating default fee settings-commit", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", true, true, "Updating default fee settings-wrong fee", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "wrong", "", "", "888888888cin", "111111111cin", "0.001", "fee-auth"}, "", nil},
		{"fee", true, true, "Add double fee settings-wrong signer", "alice", "0cin", 0, feeInfo{"sys-fee", "double", "", "", "800000000cin", "4000000000cin", "0.002", "fee-auth"}, "", nil},
		{"fee", false, false, "Add double fee settings-commit", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "double", "", "", "800000000cin", "4000000000cin", "0.002", "fee-auth"}, "", nil},
		{"fee", false, false, "Add double fee settings-commit - extra fee", "fee-auth", "100000cin", 0, feeInfo{"sys-fee", "double1", "", "", "800000000cin", "4000000000cin", "0.002", "fee-auth"}, "", nil},
		{"fee", true, true, "assign double to bank txs-wrong issuer", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "double", "bank-send", "", "", "", "", "bob"}, "", nil},
		{"fee", true, true, "assign double to bank txs-wrong signer", "bob", "0cin", 0, feeInfo{"assign-msg", "double", "bank-send", "", "", "", "", "fee-auth"}, "", nil},
		{"fee", false, true, "assign double to bank txs-wrong name", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "wrong", "bank-send", "", "", "", "", "fee-auth"}, "", nil},
		{"fee", true, true, "assign double to bank txs-empty name", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "", "bank-send", "", "", "", "", "fee-auth"}, "", nil},
		//??? No way to check msg-type is valid?
		{"fee", false, false, "assign double to bank txs-wrong msg", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "double", "wrong-send", "", "", "", "", "fee-auth"}, "", nil},
		{"fee", true, true, "assign double to bank txs-empty msg", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "double", "", "", "", "", "", "fee-auth"}, "", nil},
		{"fee", false, false, "assign double to bank txs-commit", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "double", "bank-send", "", "", "", "", "fee-auth"}, "", nil},

		// Hints : here as Special-case, which is allowed after done as above
		// should pay double
		{"kyc", false, false, "Doing kyc - yk - commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "yk", "yk", "testKyc12345678", "0"}, "", nil},
		{"bank", true, true, "sending after updating fee- wrong fee", "yk", "500000000cin", 0, bankInfo{"yk", "gohck", "1cin"}, "", nil},
		{"bank", false, false, "sending after updating fee", "yk", "800000000cin", 0, bankInfo{"yk", "mostafa", "1000cin"}, "", nil},

		// assign zero fee to an account
		{"fee", true, true, "assign zero-fee to nil address", "fee-auth", "0cin", 0, feeInfo{"assign-acc", "zero", "nope", "", "", "", "", "fee-auth"}, "", nil},
		{"fee", false, false, "assign zero-fee to mostafa-commit", "fee-auth", "0cin", 0, feeInfo{"assign-acc", "zero", "mostafa", "", "", "", "", "fee-auth"}, "", nil},
		// zero fee
		{"bank", true, true, "receiving after updating acc-fee", "mostafa", "0cin", 0, bankInfo{"bob", "mostafa", "1cin"}, "", nil},
		{"bank", false, false, "sending after updating acc-fee", "mostafa", "0cin", 0, bankInfo{"mostafa", "bob", "1cin"}, "", nil},

		// updating multiplier
		{"fee", true, true, "invalid signer", "eve", "0cin", 0, feeInfo{"multiplier", "", "", "2", "", "", "", "fee-auth"}, "", nil},
		{"fee", true, true, "mismatch signer", "fee-auth", "0cin", 0, feeInfo{"multiplier", "", "", "2", "", "", "", "eve"}, "", nil},
		{"fee", true, true, "invalid multiplier", "fee-auth", "0cin", 0, feeInfo{"multiplier", "", "", "1.*", "", "", "", "fee-auth"}, "", nil},
		{"fee", false, false, "updating multiplier", "fee-auth", "0cin", 0, feeInfo{"multiplier", "", "", "2", "", "", "", "fee-auth"}, "", nil},
		{"fee", false, false, "updating multiplier", "fee-auth", "0cin", 0, feeInfo{"multiplier", "", "", "1.0005", "", "", "", "fee-auth"}, "", nil},
		// Reset default fee setting to genesis default
		{"fee", false, false, "Updating default fee settings-commit", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "default", "", "", "100000000cin", "1000000000cin", "0.001", "fee-auth"}, "", nil},

		//=============================================start : used by nameservices modules
		{"fee", false, false, "assign zero-fee to bank alias msg-commit", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "zero", "nameservice-setAliasStatus", "", "", "", "", "fee-auth"}, "", nil},

		//=============================================start : used by nft modules
		{"fee", false, false, "assign zero-fee to mostafa-commit", "nft-fee-auth", "0cin", 0, feeInfo{"assign-acc", "zero", "nft-mostafa", "", "", "", "", "nft-fee-auth"}, "", nil},

		// assign zero fee to sdk's airdrop address
		{"fee", false, false, "assign zero-fee to sdk airdrop address", "nft-fee-auth", "0cin", 0, feeInfo{"assign-acc", "zero", "airdrop", "", "", "", "", "nft-fee-auth"}, "", nil},

		//=============================================start : used by fungible token modules
		//set fee fore msgFungibleTokenMultiplier to fee 0cin
		{"fee", false, false, "assign msgFungibleTokenMultiplier to fee 0cin. commit", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "zero", "fee-updateFungibleTokenMultiplier", "", "", "", "", "fee-auth"}, "", nil},
		{"fee", false, false, "assign msgNonFungibleTokenMultiplier to fee 0cin. commit", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "zero", "fee-updateNonFungibleTokenMultiplier", "", "", "", "", "fee-auth"}, "", nil},

		//create token action fee setting
		// ft
		{"fee", false, false, "Create ft transfer action fee setting", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "ft_transfer_default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", false, false, "Create ft mint action fee setting", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "ft_mint_default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", false, false, "Create ft burn action fee setting", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "ft_burn_default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", false, false, "Create ft transferOwnership action fee setting", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "ft_transferOwnership_default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", false, false, "Create ft acceptOwnership action fee setting", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "ft_acceptOwnership_default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},

		// nft
		{"fee", false, false, "Create nft transfer action fee setting", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "nft_transfer_default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", false, false, "Create nft mint action fee setting", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "nft_mint_default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", false, false, "Create nft burn action fee setting", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "nft_burn_default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", false, false, "Create nft transferOwnership action fee setting", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "nft_transferOwnership_default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", false, false, "Create nft acceptOwnership action fee setting", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "nft_acceptOwnership_default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", false, false, "Create nft endorse action fee setting", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "nft_endorse_default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", false, false, "Create nft updateEndorserList action fee setting", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "nft_updateNFTEndorserList_default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", false, false, "Create nft updateNFTItemMetadata action fee setting", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "nft_updateNFTItemMetadata_default", "", "", "100000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},

		//add token fee multiplier
		{"fee", false, false, "create fungible token fee multiplier. commit", "fee-auth", "0cin", 0, feeInfo{"fungible-token-fee-multiplier", "", "", "1", "", "", "", "fee-auth"}, "", nil},
		{"fee", false, false, "create nonFungible token fee multiplier. commit", "fee-auth", "0cin", 0, feeInfo{"nonFungible-token-fee-multiplier", "", "", "1", "", "", "", "fee-auth"}, "", nil},
	}

	return tcs
}

func makeFeeMsg(t *testing.T, function, name, assignee, multiplier, _min, _max, percentage, issuer string) sdkTypes.Msg {
	var msg sdkTypes.Msg
	switch function {
	case "sys-fee":
		min, err := sdkTypes.ParseCoins(_min)
		require.NoError(t, err)
		max, err := sdkTypes.ParseCoins(_max)
		require.NoError(t, err)

		msg = fee.NewMsgSysFeeSetting(name, min, max, percentage, tKeys[issuer].addr)
	case "assign-msg":
		msg = fee.NewMsgAssignFeeToMsg(name, assignee, tKeys[issuer].addr)
	case "assign-acc":
		msg = fee.NewMsgAssignFeeToAcc(name, tKeys[assignee].addr, tKeys[issuer].addr)
	case "multiplier":
		msg = fee.NewMsgMultiplier(multiplier, tKeys[issuer].addr)
	case "fungible-token-fee-multiplier":
		msg = fee.NewMsgFungibleTokenMultiplier(multiplier, tKeys[issuer].addr)
	case "nonFungible-token-fee-multiplier":
		msg = fee.NewMsgNonFungibleTokenMultiplier(multiplier, tKeys[issuer].addr)
	}

	return msg
}
