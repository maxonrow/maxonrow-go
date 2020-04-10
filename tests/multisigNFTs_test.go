package tests

import "github.com/maxonrow/maxonrow-go/utils"

func makeMultisigTxsNFTs() []*testCase {

	// Group addresses:
	// You can generate group address via `mxwcli` command. ex.
	// `mxwcli keys multisig-address "mxw1ydvzacxj0ws9jadxkmdzamc897jmln5dd90fzh" 1`
	tKeys["nft-grp-addr-1"] = &keyInfo{
		utils.MustGetAccAddressFromBech32("mxw10hu6zlsh4es6fnr3t6p884zzpdfnxfezau5khr"), nil, nil, "mxw10hu6zlsh4es6fnr3t6p884zzpdfnxfezau5khr",
	}
	// `mxwcli keys multisig-address "mxw1ydvzacxj0ws9jadxkmdzamc897jmln5dd90fzh" 2`
	tKeys["nft-grp-addr-2"] = &keyInfo{
		utils.MustGetAccAddressFromBech32("mxw1cq77f9vtl7ax9aant9sdjvyetsd5v4kryaxpwu"), nil, nil, "mxw1cq77f9vtl7ax9aant9sdjvyetsd5v4kryaxpwu",
	}
	// `mxwcli keys multisig-address "mxw1ydvzacxj0ws9jadxkmdzamc897jmln5dd90fzh" 3`
	tKeys["nft-grp-addr-3"] = &keyInfo{
		utils.MustGetAccAddressFromBech32("mxw1pamcuxj696gncrg5e3kg29k2e0gk4jc40uluht"), nil, nil, "mxw1pamcuxj696gncrg5e3kg29k2e0gk4jc40uluht",
	}

	// `mxwcli keys multisig-address "mxw1ydvzacxj0ws9jadxkmdzamc897jmln5dd90fzh" 4`
	tKeys["nft-grp-addr-4"] = &keyInfo{
		utils.MustGetAccAddressFromBech32("mxw1ukuenztvcqsytv0axrqcvrzx6f2cqmtzk7ruz3"), nil, nil, "mxw1ukuenztvcqsytv0axrqcvrzx6f2cqmtzk7ruz3",
	}

	// not exist account
	// `mxwcli keys multisig-address "mxw1ydvzacxj0ws9jadxkmdzamc897jmln5dd90fzh" 5`
	tKeys["nft-grp-addr-5"] = &keyInfo{
		utils.MustGetAccAddressFromBech32("mxw1300ry2de2uyxthzxwz2hymnahxmrw64h80p2gn"), nil, nil, "mxw1300ry2de2uyxthzxwz2hymnahxmrw64h80p2gn",
	}

	internalTx2x_createToken1 := &testCase{"nonFungibleToken", false, false, "CREATE non fungible token [TNFT-PUBLIC-FALSE-01] - Happy Path.  commit", "multisig-acc-nft-1", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-01", "nft-grp-addr-1", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx2x_createToken1", nil}
	internalTx2x_mintItem1 := &testCase{"nonFungibleToken", false, false, "MINT non fungible item [Item-PUBLIC-FALSE-01] - Happy path.  commit", "multisig-acc-nft-1", "100000000cin", 0, NonFungibleTokenInfo{"mint-item", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-01", "nft-grp-addr-1", "nft-mostafa", "token metadata", "Item-PUBLIC-FALSE-01", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx2x_mintItem1", nil}
	internalTx2x_burnItem1 := &testCase{"nonFungibleToken", false, false, "BURN non fungible item [Item-PUBLIC-FALSE-01] - Happy path.  commit", "multisig-acc-nft-1", "100000000cin", 0, NonFungibleTokenInfo{"burn-item", "", "", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-01", "nft-grp-addr-1", "", "token metadata", "Item-PUBLIC-FALSE-01", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx2x_burnItem1", nil}
	internalTx2x_transferTokenOwnership1 := &testCase{"nonFungibleToken", false, false, "TRANSFER non fungible token ownership [TNFT-PUBLIC-FALSE-01] - Happy path. commit", "multisig-acc-nft-1", "100000000cin", 0, NonFungibleTokenInfo{"transfer-token-ownership", "", "", "", "TNFT-PUBLIC-FALSE-01", "nft-grp-addr-1", "nft-grp-addr-3", "token metadata", "", "", "", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx2x_transferTokenOwnership1", nil}
	internalTx2x_acceptTokenOwnership1 := &testCase{"nonFungibleToken", false, false, "ACCEPT non fungible token ownership [TNFT-PUBLIC-FALSE-01] - Happy path. commit", "multisig-acc-nft-2", "100000000cin", 0, NonFungibleTokenInfo{"accept-token-ownership", "", "", "", "TNFT-PUBLIC-FALSE-01", "", "nft-grp-addr-3", "token metadata", "", "", "", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx2x_acceptTokenOwnership1", nil}

	internalTx3x_createToken1 := &testCase{"nonFungibleToken", false, false, "CREATE non fungible token [TNFT-PUBLIC-FALSE-02] - Happy Path.  commit", "multisig-acc-nft-2", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-02", "nft-grp-addr-2", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx3x_createToken1", nil}
	internalTx3x_mintItem1 := &testCase{"nonFungibleToken", false, false, "MINT non fungible item [Item-PUBLIC-FALSE-02] - Happy path.  commit", "multisig-acc-nft-2", "100000000cin", 0, NonFungibleTokenInfo{"mint-item", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-02", "nft-grp-addr-2", "nft-mostafa", "token metadata", "Item-PUBLIC-FALSE-02", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx3x_mintItem1", nil}
	internalTx3x_burnItem1 := &testCase{"nonFungibleToken", false, false, "BURN non fungible item [Item-PUBLIC-FALSE-02] - Happy path.  commit", "multisig-acc-nft-2", "100000000cin", 0, NonFungibleTokenInfo{"burn-item", "", "", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-02", "nft-grp-addr-2", "", "token metadata", "Item-PUBLIC-FALSE-02", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx3x_burnItem1", nil}
	internalTx3x_transferTokenOwnership1 := &testCase{"nonFungibleToken", false, false, "TRANSFER non fungible token ownership [TNFT-PUBLIC-FALSE-02] - Happy path. commit", "multisig-acc-nft-2", "100000000cin", 0, NonFungibleTokenInfo{"transfer-token-ownership", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-02", "nft-grp-addr-2", "nft-grp-addr-4", "token metadata", "", "", "", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx3x_transferTokenOwnership1", nil}
	internalTx3x_acceptTokenOwnership1 := &testCase{"nonFungibleToken", false, false, "ACCEPT non fungible token ownership [TNFT-PUBLIC-FALSE-02] - Happy path. commit", "multisig-acc-nft-2", "100000000cin", 0, NonFungibleTokenInfo{"accept-token-ownership", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-02", "", "nft-grp-addr-4", "token metadata", "", "", "", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx3x_acceptTokenOwnership1", nil}

	// refer NFTs: 3.
	internalTx3x_mintItem_err1 := &testCase{"nonFungibleToken", true, true, "MINT non fungible item [Item-PUBLIC-FALSE-02] - Invalid Token Symbol", "multisig-acc-nft-2", "100000000cin", 0, NonFungibleTokenInfo{"mint-item", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-02111", "nft-grp-addr-2", "nft-mostafa", "token metadata", "Item-PUBLIC-FALSE-02", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx3x_mintItem_err1", nil}
	internalTx3x_mintItem_err3 := &testCase{"nonFungibleToken", true, true, "MINT non fungible item [Item-PUBLIC-FALSE-02] - Re-mint Not allowed due to Token item id is in used.", "multisig-acc-nft-2", "100000000cin", 0, NonFungibleTokenInfo{"mint-item", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-02", "nft-grp-addr-2", "nft-mostafa", "token metadata", "Item-PUBLIC-FALSE-02", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx3x_mintItem_err3", nil}

	// refer NFTs: 6.
	internalTx3x_burnItem_err1 := &testCase{"nonFungibleToken", true, true, "BURN non fungible item [Item-PUBLIC-FALSE-02] - Invalid Token Symbol", "multisig-acc-nft-2", "100000000cin", 0, NonFungibleTokenInfo{"burn-item", "", "", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-029999", "nft-grp-addr-2", "", "token metadata", "Item-PUBLIC-FALSE-02", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx3x_burnItem_err1", nil}
	internalTx3x_burnItem_err3 := &testCase{"nonFungibleToken", true, true, "BURN non fungible item [Item-PUBLIC-FALSE-02] - Invalid Item ID", "multisig-acc-nft-2", "100000000cin", 0, NonFungibleTokenInfo{"burn-item", "", "", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-02", "nft-grp-addr-2", "", "token metadata", "Item-PUBLIC-FALSE-02111123", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx3x_burnItem_err3", nil}

	// refer NFTs: 7.
	internalTx3x_transferTokenOwnership_err1 := &testCase{"nonFungibleToken", true, true, "TRANSFER non fungible token ownership [TNFT-PUBLIC-FALSE-02] - Invalid token as yet to approved", "multisig-acc-nft-2", "100000000cin", 0, NonFungibleTokenInfo{"transfer-token-ownership", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-02111111123", "nft-grp-addr-2", "nft-grp-addr-4", "token metadata", "", "", "", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx3x_transferTokenOwnership_err1", nil}
	internalTx3x_transferTokenOwnership_err3 := &testCase{"nonFungibleToken", true, true, "Re-TRANSFER non fungible token ownership [TNFT-PUBLIC-FALSE-02] - Error due to Invalid TransferTokenOwnershipFlag", "multisig-acc-nft-2", "100000000cin", 0, NonFungibleTokenInfo{"transfer-token-ownership", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-02", "nft-grp-addr-2", "nft-grp-addr-4", "token metadata", "", "", "", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx3x_transferTokenOwnership_err3", nil}

	// refer NFTs: 9.
	internalTx3x_acceptTokenOwnership_err1 := &testCase{"nonFungibleToken", true, true, "ACCEPT non fungible token ownership [TNFT-PUBLIC-FALSE-02] - Invalid token as yet to approved", "multisig-acc-nft-2", "100000000cin", 0, NonFungibleTokenInfo{"accept-token-ownership", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-02777777777123", "", "nft-grp-addr-4", "token metadata", "", "", "", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx3x_acceptTokenOwnership_err1", nil}
	internalTx3x_acceptTokenOwnership_err3 := &testCase{"nonFungibleToken", true, true, "ACCEPT non fungible token ownership [TNFT-PUBLIC-FALSE-02] - Invalid token due to IsTokenOwnershipTransferrable == FALSE", "multisig-acc-nft-2", "100000000cin", 0, NonFungibleTokenInfo{"accept-token-ownership", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-02", "", "nft-grp-addr-4", "token metadata", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "MEMO: internalTx3x_acceptTokenOwnership_err3", nil}

	tcs := []*testCase{

		//create MultiSig Account1 : {"multisig-acc-nft-1"}, owner=="multisig-acc-nft-1"
		{"multiSig", false, false, "Create MultiSig Account1 - Happy Path - commit ", "multisig-acc-nft-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-nft-1", "", 1, []string{"multisig-acc-nft-1"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		//create MultiSig Account2 : {"multisig-acc-nft-2", "multisig-acc-nft-3"}, owner=="multisig-acc-nft-1"
		{"multiSig", false, false, "Create MultiSig Account2- Happy Path - commit  ", "multisig-acc-nft-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-nft-1", "", 2, []string{"multisig-acc-nft-2", "multisig-acc-nft-3"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		//create MultiSig Account3 : {"multisig-acc-nft-2", "multisig-acc-nft-3", "multisig-acc-nft-4"}, owner=="multisig-acc-nft-1"
		{"multiSig", false, false, "Create MultiSig Account3 - Happy Path - commit ", "multisig-acc-nft-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-nft-1", "", 2, []string{"multisig-acc-nft-2", "multisig-acc-nft-3", "multisig-acc-nft-4"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		{"multiSig", true, true, "Create MultiSig Account - nonkyc                ", "multisig-acc-nft-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-nft-1", "", 2, []string{"multisig-acc-nft-1", "multisig-nokyc"}, "", 0, nil}, "", nil},

		//create MultiSig Account4 : {"multisig-acc-nft-2", "multisig-acc-nft-3", "multisig-acc-nft-4"}, owner=="multisig-acc-nft-1"
		{"multiSig", false, false, "Create MultiSig Account4 - Happy Path - commit ", "multisig-acc-nft-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-nft-1", "", 3, []string{"multisig-acc-nft-2", "multisig-acc-nft-3", "multisig-acc-nft-4"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},

		//1. module : bank_test
		{"bank", false, false, "top-up Multisig Group-address1 - commit", "multisig-acc-nft-1", "800400000cin", 0, bankInfo{"multisig-acc-nft-1", "nft-grp-addr-1", "10000000000cin"}, "MEMO : top-up account", nil},
		{"bank", false, false, "top-up Multisig Group-address2 - commit", "multisig-acc-nft-2", "800400000cin", 0, bankInfo{"multisig-acc-nft-2", "nft-grp-addr-2", "10000000000cin"}, "MEMO : top-up account", nil},
		{"bank", false, false, "top-up Multisig Group-address3 - commit", "multisig-acc-nft-3", "800400000cin", 0, bankInfo{"multisig-acc-nft-3", "nft-grp-addr-3", "10000000000cin"}, "MEMO : top-up account", nil},
		{"bank", false, false, "top-up Multisig Group-address4 - commit", "multisig-acc-nft-4", "800400000cin", 0, bankInfo{"multisig-acc-nft-4", "nft-grp-addr-4", "10000000000cin"}, "MEMO : top-up account", nil},

		//2. module : kyc_test
		{"kyc", false, false, "Doing kyc - nft-mostafa - commit", "nft-kyc-auth-1", "0cin", 0, kycInfo{"nft-kyc-auth-1", "nft-kyc-issuer-1", "nft-kyc-prov-1", "whitelist", "nft-mostafa", "nft-mostafa", "testKyc123452222", "0"}, "", nil},

		//3. module : maintenance_test
		//add nft-mostafa as nonfungible authorised address
		{"maintenance", false, false, "1. Proposal, add nonfungible authorised address, Happy path. commit", "nft-maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add non fungible authorised address", "Add mostafa as non fungible authorised address", "nonFungible", "nft-mostafa", "", "", FeeCollector{}, "nft-maintainer-2", ""}, "", nil},
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

		//add NFTs fee collector with maintenance. (nft-mostafa is whitelisted.)
		{"maintenance", false, false, "4. Proposal, add NFTs fee collector address, Happy path. commit", "nft-maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add token fee collector", "Add mostafa as NFTs fee collector", "fee", "", "", "", FeeCollector{Module: "nonFungible", Address: "nft-mostafa"}, "nft-maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as NFTs fee collector, Happy path. commit", "nft-maintainer-1", "0cin", 0, CastAction{"nft-maintainer-1", "approve", 4}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as NFTs fee collector, Happy path. commit", "nft-maintainer-3", "0cin", 0, CastAction{"nft-maintainer-3", "approve", 4}, "", nil},

		//4. module : fee_test
		{"fee", false, false, "assign zero-fee to mostafa-commit", "nft-fee-auth", "0cin", 0, feeInfo{"assign-acc", "zero", "nft-mostafa", "", "", "", "", "nft-fee-auth"}, "", nil},

		// ============================================case-1.0 :  with ONE signer - HAPPY-PATH
		//create-token
		{"multiSig", false, false, "[case-1.0] Create MultiSig Tx for NFTs [Create-token] - submit counter+0 - Happy Path - wait-15-seconds.", "multisig-acc-nft-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-1", 0, internalTx2x_createToken1}, "MEMO : xxxxx", nil}, // nx "sequence":

		//Multisig Process : delete tx - before start 'Sign MultiSig Tx'
		{"multiSig", true, true, "[case-1.0] Delete MultiSig Tx - Error, due to Group address invalid. - wait-15-seconds.													", "multisig-acc-nft-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-nft-1", "", 0, nil, "nft-grp-addr-4", 0, internalTx2x_createToken1}, "MEMO : Delete MultiSig Tx", nil},
		{"multiSig", true, true, "[case-1.0] Delete MultiSig Tx - Error, due to Only group account owner can remove pending tx. - wait-15-seconds.", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-nft-2", "", 0, nil, "nft-grp-addr-1", 0, internalTx2x_createToken1}, "MEMO : Delete MultiSig Tx", nil},
		{"multiSig", true, true, "[case-1.0] Delete MultiSig Tx - Error, due to 'Pending tx is not found' which ID : 1. - wait-15-seconds.				", "multisig-acc-nft-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-nft-1", "", 0, nil, "nft-grp-addr-1", 1, internalTx2x_createToken1}, "MEMO : Delete MultiSig Tx", nil},

		{"nonFungibleToken", false, false, "[case-1.0] APPROVE nonfungible token [TNFT-PUBLIC-FALSE-01] - Happy path - wait-15-seconds.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-01", "nft-grp-addr-1", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "nft-jeansoon", "0", "nft-carlo", "default", "", "2", "2", []string{"nft-jeansoon", "nft-carlo"}}, "APPROVE nonfungible token [TNFT-PUBLIC-FALSE-01] - Happy path", nil},

		//mintItem - endorseItem - transferItem - burnItem
		{"multiSig", false, false, "[case-1.0] Create MultiSig Tx for NFTs [Mint-item] - submit counter+1 - Happy Path - wait-15-seconds.", "multisig-acc-nft-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-1", 1, internalTx2x_mintItem1}, "MEMO : xxxxx", nil}, // nx "sequence":
		{"nonFungibleToken", false, false, "[case-1.0] ENDORSE nonfungible item [Item-PUBLIC-FALSE-01] - Happy path - wait-15-seconds.", "nft-carlo", "100000000cin", 0, NonFungibleTokenInfo{"endorsement-item", "", "", "", "TNFT-PUBLIC-FALSE-01", "nft-carlo", "", "token metadata", "Item-PUBLIC-FALSE-01", "", "", true, false, true, true, false, "", "", "", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "ENDORSE nonfungible item [Item-PUBLIC-FALSE-01] - Happy path", nil},
		{"nonFungibleToken", false, false, "[case-1.0] TRANSFER nonfungible token item [Item-PUBLIC-FALSE-01] - Happy path - wait-15-seconds.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"transfer-item", "", "", "", "TNFT-PUBLIC-FALSE-01", "nft-mostafa", "nft-grp-addr-1", "token metadata", "Item-PUBLIC-FALSE-01", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "TRANSFER nonfungible token item [Item-PUBLIC-FALSE-01] - Happy path", nil},
		{"multiSig", false, false, "[case-1.0] Create MultiSig Tx for NFTs [Burn-item] - submit counter+2 - Happy Path - wait-15-seconds.", "multisig-acc-nft-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-1", 2, internalTx2x_burnItem1}, "MEMO : xxxxx", nil}, // nx "sequence":

		//transferTokenOwnership - Verify - Accept
		{"multiSig", false, false, "[case-1.0] Create MultiSig Tx for NFTs [Transfer-token-ownership] - submit counter+3 - Happy Path - wait-15-seconds.", "multisig-acc-nft-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-1", 3, internalTx2x_transferTokenOwnership1}, "MEMO : xxxxx", nil}, // nx "sequence":
		{"nonFungibleToken", false, false, "[case-1.0] VERIFY nonfungible token transfer ownership [TNFT-PUBLIC-FALSE-01] - Happy path - wait-15-seconds.", "nft-mostafa", "0cin", 0, NonFungibleTokenInfo{"verify-transfer-token-ownership", "", "", "", "TNFT-PUBLIC-FALSE-01", "nft-mostafa", "nft-grp-addr-1", "token metadata", "", "", "", true, false, true, true, false, "nft-jeansoon", "0", "nft-carlo", "", "APPROVE_TRANFER_TOKEN_OWNERSHIP", "", "", []string{""}}, "VERIFY nonfungible token transfer ownership [TNFT-PUBLIC-FALSE-01] - Happy path", nil},

		//Need signed by : {"multisig-acc-nft-2", "multisig-acc-nft-3", "multisig-acc-nft-4"}, owner=="multisig-acc-nft-1"
		//Threshold == 2, under 'nft-grp-addr-3'
		{"multiSig", false, false, "[case-1.0] Create MultiSig Tx for NFTs [Accept-token-ownership] - submit counter+0 - Happy Path - wait-15-seconds.", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-3", 0, internalTx2x_acceptTokenOwnership1}, "MEMO : xxxxx", nil}, // nx "sequence":
		{"multiSig", false, false, "[case-1.0] Sign MultiSig Tx for NFTs [Accept-token-ownership] - submit which signed by multisig-acc-nft-4 - wait-15-seconds.", "multisig-acc-nft-4", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "nft-grp-addr-3", 0, internalTx2x_acceptTokenOwnership1}, "MEMO : Sign MultiSig Tx", nil},

		//Multisig Process : Multisig-transfer-ownership
		{"multiSig", true, true, "[case-1.0] Transfer MultiSig Owner - Error, due to Group address invalid. - wait-15-seconds.                                                                ", "multisig-acc-nft-1", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-nft-1", "multisig-acc-nft-2", 0, nil, "nft-grp-addr-5", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},
		{"multiSig", true, true, "[case-1.0] Transfer MultiSig Owner - Error, due to Owner of group address invalid. - wait-15-seconds.                                                       ", "multisig-acc-nft-3", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-nft-3", "multisig-acc-nft-1", 0, nil, "nft-grp-addr-1", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},
		{"multiSig", true, true, "[case-1.0] Transfer MultiSig Owner - Error, due to without KYC. - wait-15-seconds.                                                                          ", "multisig-nokyc", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-nft-1", "multisig-nokyc", 0, nil, "nft-grp-addr-1", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},
		{"multiSig", false, false, "[case-1.0] Transfer MultiSig Owner - [from multisig-acc-nft-1 to multisig-acc-nft-2] - Happy Path - wait-15-seconds.                                   						", "multisig-acc-nft-1", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-nft-1", "multisig-acc-nft-2", 0, nil, "nft-grp-addr-1", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},
		{"multiSig", true, true, "[case-1.0] Re-transfer MultiSig Owner - Error, due to Owner of group address invalid as MultiSig-account already been transfer to others. - wait-15-seconds.", "multisig-acc-nft-1", "800400000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-nft-1", "multisig-acc-nft-2", 0, nil, "nft-grp-addr-1", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},

		// ============================================case-2.0 :  with TWO signers - with InternalTx Error-cases
		//create-token
		{"multiSig", false, false, "[case-2.0] Create MultiSig Tx for NFTs [Create-token] - submit - Happy Path - wait-15-seconds.", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-2", 0, internalTx3x_createToken1}, "MEMO : Create MultiSig Tx", nil},

		//Multisig Process : delete tx - before start 'Sign MultiSig Tx'
		{"multiSig", true, true, "[case-2.0] Delete MultiSig Tx - Error, due to Group address invalid. - wait-15-seconds.													", "multisig-acc-nft-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-nft-1", "", 0, nil, "nft-grp-addr-4", 0, internalTx3x_createToken1}, "MEMO : Delete MultiSig Tx", nil},
		{"multiSig", true, true, "[case-2.0] Delete MultiSig Tx - Error, due to Only group account owner can remove pending tx. - wait-15-seconds.", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-nft-2", "", 0, nil, "nft-grp-addr-2", 0, internalTx3x_createToken1}, "MEMO : Delete MultiSig Tx", nil},
		{"multiSig", true, true, "[case-2.0] Delete MultiSig Tx - Error, due to 'Pending tx is not found' which ID : 1. - wait-15-seconds.				", "multisig-acc-nft-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-nft-1", "", 0, nil, "nft-grp-addr-2", 1, internalTx3x_createToken1}, "MEMO : Delete MultiSig Tx", nil},

		{"multiSig", false, false, "[case-2.0] Sign MultiSig Tx for NFTs [Create-token] - submit which signed by multisig-acc-nft-3 - wait-15-seconds.", "multisig-acc-nft-3", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "nft-grp-addr-2", 0, internalTx3x_createToken1}, "MEMO : Sign MultiSig Tx", nil},
		{"nonFungibleToken", false, false, "[case-2.0] APPROVE nonfungible token [TNFT-PUBLIC-FALSE-02] - Happy path - wait-15-seconds.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-02", "nft-grp-addr-2", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "nft-jeansoon", "0", "nft-carlo", "default", "", "2", "2", []string{"nft-jeansoon", "nft-carlo"}}, "[case-2.0] APPROVE nonfungible token [TNFT-PUBLIC-FALSE-02] - Happy path", nil},

		//mintItem - endorseItem - transferItem - burnItem
		{"multiSig", false, false, "[case-2.0] Create MultiSig Tx for NFTs [Mint-item] - submit counter+1 - Happy Path - wait-15-seconds.", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-2", 1, internalTx3x_mintItem1}, "MEMO : xxxxx", nil},    // nx "sequence":
		{"multiSig", false, false, "[case-2.0] Create MultiSig Tx for NFTs [Mint-item] - [internalTx3x_mintItem_err1] - wait-15-seconds.", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-2", 2, internalTx3x_mintItem_err1}, "MEMO : xxxxx", nil}, // nx "sequence":
		{"multiSig", false, false, "[case-2.0] Create MultiSig Tx for NFTs [Mint-item] - [internalTx3x_mintItem_err3] - wait-15-seconds.", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-2", 3, internalTx3x_mintItem_err3}, "MEMO : xxxxx", nil}, // nx "sequence":
		{"multiSig", false, false, "[case-2.0] Sign MultiSig Tx for NFTs [Mint-item] - submit which signed by multisig-acc-nft-3 - wait-15-seconds.", "multisig-acc-nft-3", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "nft-grp-addr-2", 1, internalTx3x_mintItem1}, "MEMO : Sign MultiSig Tx", nil},

		{"nonFungibleToken", false, false, "[case-2.0] ENDORSE nonfungible item [Item-PUBLIC-FALSE-02] - Happy path - wait-15-seconds.", "nft-carlo", "100000000cin", 0, NonFungibleTokenInfo{"endorsement-item", "", "", "", "TNFT-PUBLIC-FALSE-02", "nft-carlo", "", "token metadata", "Item-PUBLIC-FALSE-02", "", "", true, false, true, true, false, "", "", "", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "[case-2.0] ENDORSE nonfungible item [Item-PUBLIC-FALSE-02] - Happy path", nil},
		{"nonFungibleToken", false, false, "[case-2.0] TRANSFER nonfungible token item [Item-PUBLIC-FALSE-02] - Happy path - wait-15-seconds.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"transfer-item", "", "", "", "TNFT-PUBLIC-FALSE-02", "nft-mostafa", "nft-grp-addr-2", "token metadata", "Item-PUBLIC-FALSE-02", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "[case-2.0] TRANSFER nonfungible token item [Item-PUBLIC-FALSE-02] - Happy path", nil},
		{"multiSig", false, false, "[case-2.0] Create MultiSig Tx for NFTs [Burn-item] - submit counter+2 - Happy Path - wait-15-seconds.", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-2", 4, internalTx3x_burnItem1}, "MEMO : xxxxx", nil},    // nx "sequence":
		{"multiSig", false, false, "[case-2.0] Create MultiSig Tx for NFTs [Burn-item] - [internalTx3x_burnItem_err1] - wait-15-seconds.", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-2", 5, internalTx3x_burnItem_err1}, "MEMO : xxxxx", nil}, // nx "sequence":
		{"multiSig", false, false, "[case-2.0] Create MultiSig Tx for NFTs [Burn-item] - [internalTx3x_burnItem_err3] - wait-15-seconds.", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-2", 6, internalTx3x_burnItem_err3}, "MEMO : xxxxx", nil}, // nx "sequence":
		{"multiSig", false, false, "[case-2.0] Sign MultiSig Tx for NFTs [Burn-item] - submit which signed by multisig-acc-nft-3 - wait-15-seconds.", "multisig-acc-nft-3", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "nft-grp-addr-2", 4, internalTx3x_burnItem1}, "MEMO : Sign MultiSig Tx", nil},

		//transferTokenOwnership - Verify - Accept
		{"multiSig", false, false, "[case-2.0] Create MultiSig Tx for NFTs [Transfer-token-ownership] - submit counter+3 - Happy Path - wait-15-seconds.", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-2", 7, internalTx3x_transferTokenOwnership1}, "MEMO : xxxxx", nil},                  // nx "sequence":
		{"multiSig", false, false, "[case-2.0] Create MultiSig Tx for NFTs [Transfer-token-ownership] - [internalTx3x_transferTokenOwnership_err1] - wait-15-seconds.", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-2", 8, internalTx3x_transferTokenOwnership_err1}, "MEMO : xxxxx", nil}, // nx "sequence":
		{"multiSig", false, false, "[case-2.0] Create MultiSig Tx for NFTs [Transfer-token-ownership] - [internalTx3x_transferTokenOwnership_err3] - wait-15-seconds.", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-2", 9, internalTx3x_transferTokenOwnership_err3}, "MEMO : xxxxx", nil}, // nx "sequence":
		{"multiSig", false, false, "[case-2.0] Sign MultiSig Tx for NFTs [Transfer-token-ownership] - submit which signed by multisig-acc-nft-3 - wait-15-seconds.", "multisig-acc-nft-3", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "nft-grp-addr-2", 7, internalTx3x_transferTokenOwnership1}, "MEMO : Sign MultiSig Tx", nil},
		{"nonFungibleToken", false, false, "[case-2.0] VERIFY nonfungible token transfer ownership [TNFT-PUBLIC-FALSE-02] - Happy path - wait-15-seconds.", "nft-mostafa", "0cin", 0, NonFungibleTokenInfo{"verify-transfer-token-ownership", "", "", "", "TNFT-PUBLIC-FALSE-02", "nft-mostafa", "nft-grp-addr-2", "token metadata", "", "", "", true, false, true, true, false, "nft-jeansoon", "0", "nft-carlo", "", "APPROVE_TRANFER_TOKEN_OWNERSHIP", "", "", []string{""}}, "[case-2.0] VERIFY nonfungible token transfer ownership [TNFT-PUBLIC-FALSE-02] - Happy path", nil},

		// Need signed by : {"multisig-acc-nft-2", "multisig-acc-nft-3", "multisig-acc-nft-4"}, owner=="multisig-acc-nft-1"
		// Threshold == 3, under 'nft-grp-addr-4'
		{"multiSig", false, false, "[case-2.0] Create MultiSig Tx for NFTs [Accept-token-ownership] - submit counter+0 - Happy Path - wait-15-seconds.", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-4", 0, internalTx3x_acceptTokenOwnership1}, "MEMO : xxxxx", nil},                // nx "sequence":
		{"multiSig", false, false, "[case-2.0] Create MultiSig Tx for NFTs [Accept-token-ownership] - [internalTx3x_acceptTokenOwnership_err1] - wait-15-seconds.", "multisig-acc-nft-4", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-4", 1, internalTx3x_acceptTokenOwnership_err1}, "MEMO : xxxxx", nil}, // nx "sequence":
		{"multiSig", false, false, "[case-2.0] Create MultiSig Tx for NFTs [Accept-token-ownership] - [internalTx3x_acceptTokenOwnership_err3] - wait-15-seconds.", "multisig-acc-nft-4", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "nft-grp-addr-4", 2, internalTx3x_acceptTokenOwnership_err3}, "MEMO : xxxxx", nil}, // nx "sequence":
		{"multiSig", false, false, "[case-2.0] Sign MultiSig Tx for NFTs [Accept-token-ownership] - submit which signed by multisig-acc-nft-4 - wait-15-seconds.", "multisig-acc-nft-4", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "nft-grp-addr-4", 0, internalTx3x_acceptTokenOwnership1}, "MEMO : Sign MultiSig Tx", nil},
		{"multiSig", false, false, "[case-2.0] Sign MultiSig Tx for NFTs [Accept-token-ownership] - submit which signed by multisig-acc-nft-3 - wait-15-seconds.", "multisig-acc-nft-3", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "nft-grp-addr-4", 0, internalTx3x_acceptTokenOwnership1}, "MEMO : Sign MultiSig Tx", nil},

		//Multisig Process : Multisig-transfer-ownership
		{"multiSig", true, true, "[case-2.0] Transfer MultiSig Owner - Error, due to Group address invalid. - wait-15-seconds.																															 ", "multisig-acc-nft-1", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-nft-1", "multisig-acc-nft-4", 0, nil, "nft-grp-addr-5", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},
		{"multiSig", true, true, "[case-2.0] Transfer MultiSig Owner - Error, due to Owner of group address invalid. - wait-15-seconds.																											 ", "multisig-acc-nft-2", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-nft-2", "multisig-acc-nft-4", 0, nil, "nft-grp-addr-3", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},
		{"multiSig", true, true, "[case-2.0] Transfer MultiSig Owner - Error, due to without KYC - wait-15-seconds.																																					 ", "multisig-nokyc", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-nft-1", "multisig-nokyc", 0, nil, "nft-grp-addr-5", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},
		{"multiSig", false, false, "[case-2.0] Transfer MultiSig Owner - [from multisig-acc-nft-1 to multisig-acc-nft-4] Happy Path - commit. - wait-15-seconds.																		 ", "multisig-acc-nft-1", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-nft-1", "multisig-acc-nft-4", 0, nil, "nft-grp-addr-3", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},
		{"multiSig", true, true, "[case-2.0] Re-transfer MultiSig Owner - Error, due to Owner of group address invalid [MultiSig-account already been transfer to others]. - wait-15-seconds.", "multisig-acc-nft-1", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-nft-1", "multisig-acc-nft-4", 0, nil, "nft-grp-addr-3", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil},

		// signer without through KYC
		{"multiSig", true, true, "Create MultiSig Account - Error, due to without KYC            ", "multisig-nokyc", "800400000cin", 0, MultisigInfo{"create", "multisig-nokyc", "", 2, []string{"multisig-acc-nft-1", "multisig-nokyc"}, "", 0, nil}, "", nil},
	}

	return tcs
}
