package tests

import "github.com/maxonrow/maxonrow-go/utils"

func makeMultisigTxsNFTs() []*testCase {

	// Group addresses:
	// You can generate group address via `mxwcli` command. ex.
	// `mxwcli keys multisig-address "mxw1ld3stcsk5l8xjngw2ucuazux895rk2hxve69gr" 1`
	tKeys["grp-addr-1"] = &keyInfo{
		utils.MustGetAccAddressFromBech32("mxw1z8r356ll7aum0530xve2upx74ed8ffavyxy503"), nil, nil, "mxw1z8r356ll7aum0530xve2upx74ed8ffavyxy503",
	}
	// `mxwcli keys multisig-address "mxw1ld3stcsk5l8xjngw2ucuazux895rk2hxve69gr" 2`
	tKeys["grp-addr-2"] = &keyInfo{
		utils.MustGetAccAddressFromBech32("mxw1q6nmfejarl5e4xzceqxcygner7a6llgwnrdtl6"), nil, nil, "mxw1q6nmfejarl5e4xzceqxcygner7a6llgwnrdtl6",
	}
	// `mxwcli keys multisig-address "mxw1ld3stcsk5l8xjngw2ucuazux895rk2hxve69gr" 3`
	tKeys["grp-addr-3"] = &keyInfo{
		utils.MustGetAccAddressFromBech32("mxw1hkm4p04nsmv9q0hg4m9eeuapfdr7n4rfl04vh9"), nil, nil, "mxw1hkm4p04nsmv9q0hg4m9eeuapfdr7n4rfl04vh9",
	}
	// not exist account
	// `mxwcli keys multisig-address "mxw1ld3stcsk5l8xjngw2ucuazux895rk2hxve69gr" 4`
	tKeys["grp-addr-4"] = &keyInfo{
		utils.MustGetAccAddressFromBech32("mxw1szm87m362urkvj833jd7nekwdjh7s8p4q3f25f"), nil, nil, "mxw1szm87m362urkvj833jd7nekwdjh7s8p4q3f25f",
	}

	internalTx4 := &testCase{"nonFungibleToken", false, false, "Create non fungible token [TNFT-PUBLIC-FALSE-00] - Happy Path.  commit", "multisig-acc-1", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-00", "grp-addr-1", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token.", nil}
	internalTx5 := &testCase{"nonFungibleToken", false, false, "Create non fungible token [TNFT-PUBLIC-FALSE-00] - Happy Path.  commit", "multisig-acc-2", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-00", "grp-addr-2", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token.", nil}

	tcs := []*testCase{

		//create MultiSig Account1 : {"multisig-acc-1"}, owner=="multisig-acc-1"
		{"multiSig", false, false, "Create MultiSig Account1 - Happy Path - commit ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 1, []string{"multisig-acc-1"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		//create MultiSig Account2 : {"multisig-acc-2", "multisig-acc-3"}, owner=="multisig-acc-1"
		{"multiSig", false, false, "Create MultiSig Account2- Happy Path - commit  ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 2, []string{"multisig-acc-2", "multisig-acc-3"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		//create MultiSig Account3 : {"multisig-acc-2", "multisig-acc-3", "multisig-acc-4"}, owner=="multisig-acc-1"
		{"multiSig", false, false, "Create MultiSig Account3 - Happy Path - commit ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 2, []string{"multisig-acc-2", "multisig-acc-3", "multisig-acc-4"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		{"multiSig", true, true, "Create MultiSig Account - non-kyc                ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 2, []string{"multisig-acc-1", "multisig-nokyc"}, "", 0, nil}, "", nil},

		//1. module : bank_test
		{"bank", false, false, "top-up Multisig Group-address1 - commit", "multisig-acc-1", "800400000cin", 0, bankInfo{"multisig-acc-1", "grp-addr-1", "10000000000cin"}, "MEMO : top-up account", nil},
		{"bank", false, false, "top-up Multisig Group-address2 - commit", "multisig-acc-2", "800400000cin", 0, bankInfo{"multisig-acc-2", "grp-addr-2", "10000000000cin"}, "MEMO : top-up account", nil},
		{"bank", false, false, "top-up Multisig Group-address3 - commit", "multisig-acc-3", "800400000cin", 0, bankInfo{"multisig-acc-3", "grp-addr-3", "10000000000cin"}, "MEMO : top-up account", nil},

		//2. module : kyc_test
		{"kyc", false, false, "Doing kyc - nft-mostafa - commit", "nft-kyc-auth-1", "0cin", 0, kycInfo{"nft-kyc-auth-1", "nft-kyc-issuer-1", "nft-kyc-prov-1", "whitelist", "nft-mostafa", "nft-mostafa", "testKyc123452222", "0"}, "", nil},

		//3. module : maintenance_test
		//add multisig-acc-1 as nonfungible authorised address
		{"maintenance", false, false, "1. Proposal, add nonfungible authorised address, Happy path. commit", "nft-maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add non fungible authorised address", "Add mostafa as non fungible authorised address", "nonFungible", "multisig-acc-1", "", "", FeeCollector{}, "nft-maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as non fungible token authorised address, Happy path. commit", "nft-maintainer-1", "0cin", 0, CastAction{"nft-maintainer-1", "approve", 1}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as non fungible token authorised address, Happy path. commit", "nft-maintainer-3", "0cin", 0, CastAction{"nft-maintainer-3", "approve", 1}, "", nil},

		//add nft-carlo as nonfungible issuer address
		{"maintenance", false, false, "2. Proposal, add nonfungible issuer address, Happy path. commit", "nft-maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add non fungible fee issuer address", "Add carlo as non fungible issuer address", "nonFungible", "", "nft-carlo", "", FeeCollector{}, "nft-maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve carlo as non fungible token issuer address, Happy path. commit", "nft-maintainer-1", "0cin", 0, CastAction{"nft-maintainer-1", "approve", 2}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve carlo as non fungible token issuer address, Happy path. commit", "nft-maintainer-3", "0cin", 0, CastAction{"nft-maintainer-3", "approve", 2}, "", nil},

		//add nft-jeansoon as nonfungible provider address
		{"maintenance", false, false, "3. Proposal, add nonfungible provider address, Happy path. commit", "nft-maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add non fungible fee provider address", "Add jeansoon as non fungible provider address", "nonFungible", "", "", "nft-jeansoon", FeeCollector{}, "nft-maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "Cast action to approve jeansoon as non fungible token provider address, Happy path. commit", "nft-maintainer-1", "0cin", 0, CastAction{"nft-maintainer-1", "approve", 3}, "", nil},
		{"maintenance-cast-action", false, false, "Cast action to approve jeansoon as non fungible token provider address, Happy path. commit", "nft-maintainer-3", "0cin", 0, CastAction{"nft-maintainer-3", "approve", 3}, "", nil},

		//add nameservice fee collector with maintenance. (multisig-acc-1 is whitelisted.)
		{"maintenance", false, false, "4. Proposal, add token fee collector address, Happy path. commit", "nft-maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add token fee collector", "Add mostafa as nameservice fee collector", "fee", "", "", "", FeeCollector{Module: "nonFungible", Address: "multisig-acc-1"}, "nft-maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as nameservice fee collector, Happy path. commit", "nft-maintainer-1", "0cin", 0, CastAction{"nft-maintainer-1", "approve", 4}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as nameservice fee collector, Happy path. commit", "nft-maintainer-3", "0cin", 0, CastAction{"nft-maintainer-3", "approve", 4}, "", nil},

		//4. module : fee_test
		{"fee", false, false, "assign zero-fee to mostafa-commit", "nft-fee-auth", "0cin", 0, feeInfo{"assign-acc", "zero", "nft-mostafa", "", "", "", "", "nft-fee-auth"}, "", nil},

		// case-1 :  with ONE signer, should broadcast immediately
		{"multiSig", false, false, "Create MultiSig Tx for NFTs [Create-token] - submit counter+0 - Happy Path commit.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx4}, "MEMO : xxxxx", nil},

		// case-2 :  with TWO signers
		{"multiSig", false, false, "Create MultiSig Tx for NFTs [Create-token] - submit - Happy Path commit.                        ", "multisig-acc-2", "800400000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx5}, "MEMO : Create MultiSig Tx", nil},
		{"multiSig", false, false, "Sign MultiSig Tx for NFTs [Create-token] - submit which signed by multisig-acc-3 - commit.      ", "multisig-acc-3", "800400000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx5}, "MEMO : Sign MultiSig Tx", nil},

		//====================start : case-xxxx
		// signer without through KYC
		{"multiSig", true, true, "Create MultiSig Account - Error, due to without KYC            ", "multisig-nokyc", "800400000cin", 0, MultisigInfo{"create", "multisig-nokyc", "", 2, []string{"multisig-acc-1", "multisig-nokyc"}, "", 0, nil}, "", nil},
	}

	return tcs
}
