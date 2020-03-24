package tests

func makeMultisigTxsNFTs() []*testCase {

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

	// internalTx1 := &testCase{"nonFungibleToken", true, true, "Create non fungible token [TNFT-PUBLIC-FALSE] - Happy Path.  commit", "multisig-acc-1", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE", "mostafa", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token.", nil}
	// internalTx2 := &testCase{"nonFungibleToken", true, true, "Create non fungible token [TNFT-PUBLIC-FALSE] - Happy Path.  commit", "multisig-acc-1", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE", "grp-addr-3", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token.", nil}
	internalTx3_createToken_00 := &testCase{"nonFungibleToken", false, false, "Create non fungible token [TNFT-PUBLIC-FALSE-00] - Happy Path.  commit", "multisig-acc-1", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-00", "grp-addr-1", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token.", nil}
	// internalTx3_createToken_01 := &testCase{"nonFungibleToken", false, false, "Create non fungible token [TNFT-PUBLIC-FALSE-01] - Happy Path.  commit", "multisig-acc-1", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-01", "grp-addr-1", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token.", nil}
	// internalTx3_createToken_02 := &testCase{"nonFungibleToken", false, false, "Create non fungible token [TNFT-PUBLIC-FALSE-02] - Happy Path.  commit", "multisig-acc-1", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-02", "grp-addr-1", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token.", nil}
	// internalTx3_createToken_03 := &testCase{"nonFungibleToken", false, false, "Create non fungible token [TNFT-PUBLIC-FALSE-03] - Happy Path.  commit", "multisig-acc-1", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-03", "grp-addr-1", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token.", nil}
	// internalTx3_createToken_04 := &testCase{"nonFungibleToken", false, false, "Create non fungible token [TNFT-PUBLIC-FALSE-04] - Happy Path.  commit", "multisig-acc-1", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-04", "grp-addr-1", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token.", nil}
	// internalTx3_createToken_05 := &testCase{"nonFungibleToken", false, false, "Create non fungible token [TNFT-PUBLIC-FALSE-05] - Happy Path.  commit", "multisig-acc-1", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-05", "grp-addr-1", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token.", nil}

	internalTx3_approveToken_00 := &testCase{"nonFungibleToken", false, false, "APPROVE non fungible token [TNFT-PUBLIC-FALSE-00] - Happy Path.  commit", "multisig-acc-1", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "multisig-acc-1", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-00", "grp-addr-1", "", "token metadata", "", "item-properties", "item-metadata", true, false, true, true, false, "nft-jeansoon", "0", "nft-carlo", "default", "", "2", "2", []string{"nft-jeansoon", "nft-carlo"}}, "", nil}
	internalTx3_mintItem_00 := &testCase{"nonFungibleToken", false, false, "MINT non fungible item [TNFT-PUBLIC-FALSE-00] - Happy Path.  commit", "multisig-acc-1", "100000000cin", 0, NonFungibleTokenInfo{"mint-item", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-PUBLIC-FALSE-00", "grp-addr-1", "nft-mostafa", "token metadata", "Item-PUBLIC-FALSE-00", "item-properties", "item-metadata", true, false, true, true, false, "nft-jeansoon", "0", "nft-carlo", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "", nil}
	// internalTx3_burnItem_00 := &testCase{"nonFungibleToken", false, false, "Burn non fungible token item [TNFT-PUBLIC-FALSE-00] - Happy Path.  commit", "multisig-acc-1", "100000000cin", 0, NonFungibleTokenInfo{"burn-item", "", "", "", "TNFT-PUBLIC-FALSE-00", "grp-addr-1", "", "token metadata", "Item-PUBLIC-FALSE-00", "item-properties", "item-metadata", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "", nil}
	// internalTx3_transferTokenOwnership_00 := &testCase{"nonFungibleToken", false, false, "Transfer non fungible token ownership [TNFT-PUBLIC-FALSE-00] - Happy path", "multisig-acc-1", "100000000cin", 0, NonFungibleTokenInfo{"transfer-token-ownership", "", "", "", "TNFT-PUBLIC-FALSE-00", "grp-addr-1", "grp-addr-2", "token metadata", "", "", "", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "", nil}
	// internalTx3_acceptTokenOwnership_00 := &testCase{"nonFungibleToken", false, false, "Accept non fungible token ownership [TNFT-PUBLIC-FALSE-00] - Happy path. commit", "multisig-acc-2", "100000000cin", 0, NonFungibleTokenInfo{"accept-token-ownership", "", "", "", "TNFT-PUBLIC-FALSE-00", "", "grp-addr-2", "token metadata", "", "", "", true, false, true, true, false, "", "", "", "", "", "", "", []string{""}}, "", nil}

	tcs := []*testCase{

		//create MultiSig Account1 : {"multisig-acc-1"}, owner=="multisig-acc-1"
		{"multiSig", false, false, "Create MultiSig Account1 - Happy Path - commit ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 1, []string{"multisig-acc-1"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		//create MultiSig Account2 : {"multisig-acc-2", "multisig-acc-3"}, owner=="multisig-acc-1"
		{"multiSig", false, false, "Create MultiSig Account2- Happy Path - commit  ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 2, []string{"multisig-acc-2", "multisig-acc-3"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		//create MultiSig Account3 : {"multisig-acc-2", "multisig-acc-3", "multisig-acc-4"}, owner=="multisig-acc-1"
		{"multiSig", false, false, "Create MultiSig Account3 - Happy Path - commit ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 2, []string{"multisig-acc-2", "multisig-acc-3", "multisig-acc-4"}, "", 0, nil}, "MEMO : Create MultiSig Account - Happy Path", nil},
		{"multiSig", true, true, "Create MultiSig Account - non-kyc                ", "multisig-acc-1", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-1", "", 2, []string{"multisig-acc-1", "multisig-acc-no-kyc"}, "", 0, nil}, "", nil},

		//1. module : bank_test (mtf ORIG)
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

		// NFTs testing :
		{"multiSig", false, false, "Create MultiSig Tx for NFTs [Create-token] - submit counter+0 - Happy Path commit.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx3_createToken_00}, "MEMO : xxxxx", nil},
		{"multiSig", false, false, "Create MultiSig Tx for NFTs [Approve-token] - Happy path -  commit.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 1, internalTx3_approveToken_00}, "MEMO : xxxxx", nil},
		{"multiSig", false, false, "Create MultiSig Tx for NFTs [Mint-item] - Happy path -  commit.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 2, internalTx3_mintItem_00}, "MEMO : xxxxx", nil},

		/* -- Done Goh, yet to finalise.
		//=========================================== start : case-1.1 Create-token (CREATE-TX, SIGN-TX, DELETE-TX, TRANSFER-OWNERSHIP-TX)===========================================
		//Remarks : Using 'MultiSig Account1' which owner=="multisig-acc-1"
		//Scenario : using 'grp-addr-1' which only with ONE signer {'MultiSig Account1'}, should broadcast immediately
		{"multiSig", true, true, "Create MultiSig Tx for NFTs [Create-token] - Error, Invalid sequence   ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 5, internalTx3_createToken_00}, "MEMO : xxxx", nil},

		// Create-token : multiSig-create
		// Remarks: with one signer, should broadcast immediately
		{"multiSig", true, true, "Create MultiSig Tx for NFTs [Create-token] - Error, Invalid signer due to Sender is not group account's signer.       ", "multisig-acc-4", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx1}, "MEMO : xxxx", nil},
		{"multiSig", true, true, "Create MultiSig Tx for NFTs [Create-token] - Error, Invalid sender due to Sender is not group account's signer.    ", "multisig-acc-2", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx1}, "MEMO : xxxx", nil},
		{"multiSig", true, true, "Create MultiSig Tx for NFTs [Create-token] - Error, Invalid sender2 due to due to Sender is not group account's signer.   ", "mostafa", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx1}, "MEMO : xxxx", nil},
		{"multiSig", true, true, "Create MultiSig Tx for NFTs [Create-token] - Error, Invalid tx-id due to Internal transaction signature error.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 1, internalTx1}, "MEMO : xxxx", nil},
		{"multiSig", true, true, "Create MultiSig Tx for NFTs [Create-token] - Error, Invalid internal_tx due to Internal transaction signature error.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx2}, "MEMO : xxxx", nil},
		{"multiSig", false, false, "Create MultiSig Tx for NFTs [Create-token] - submit counter+0 - Happy Path commit.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 0, internalTx3_createToken_00}, "MEMO : xxxxx", nil},
		{"multiSig", false, false, "Create MultiSig Tx for NFTs [Create-token] - submit counter+1 - Happy Path commit.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 1, internalTx3_createToken_01}, "MEMO : xxxxx", nil},
		{"multiSig", false, false, "Create MultiSig Tx for NFTs [Create-token] - submit counter+2 - Happy Path commit.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 2, internalTx3_createToken_02}, "MEMO : xxxxx", nil},
		{"multiSig", false, false, "Create MultiSig Tx for NFTs [Create-token] - submit counter+3 - Happy Path commit.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 3, internalTx3_createToken_03}, "MEMO : xxxxx", nil},
		{"multiSig", false, false, "Create MultiSig Tx for NFTs [Create-token] - submit counter+4 - Happy Path commit.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 4, internalTx3_createToken_04}, "MEMO : xxxxx", nil},
		{"multiSig", false, false, "Create MultiSig Tx for NFTs [Create-token] - submit counter+5 - Happy Path commit.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 5, internalTx3_createToken_05}, "MEMO : xxxxx", nil},

		//Create-token : multiSig-delete
		{"multiSig", true, true, "Delete MultiSig Tx for NFTs [Create-token] -  Error, due to Group address invalid.                         ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-4", 0, internalTx3_createToken_00}, "MEMO : Delete MultiSig Tx", nil},
		{"multiSig", true, true, "Delete MultiSig Tx for NFTs [Create-token] -  Error, Invalid Owner address due to Only group account owner can remove pending tx.                         ", "multisig-acc-2", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-2", "", 0, nil, "grp-addr-1", 0, internalTx3_createToken_00}, "MEMO : Delete MultiSig Tx", nil},
		{"multiSig", true, true, "Delete MultiSig Tx for NFTs [Create-token] -  Error, due to Pending tx is not found which ID : 33.       ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-1", 33, internalTx3_createToken_00}, "MEMO : Delete MultiSig Tx", nil},

		//Create-token : multiSig-transfer-ownership
		{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [Create-token] -  Error, due to Group address invalid.                                                                ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-4", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
		{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [Create-token] -  Error, due to Owner of group address invalid.                                                       ", "multisig-acc-3", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-3", "multisig-acc-1", 0, nil, "grp-addr-1", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
		{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [Create-token] -  Error, due to without KYC																																					 ", "multisig-acc-no-kyc", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-no-kyc", 0, nil, "grp-addr-1", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //

		//--dx=========================================== start : case-1.1 approve-token (CREATE-TX, SIGN-TX, DELETE-TX, TRANSFER-OWNERSHIP-TX)===========================================
		//Remarks :
		{"multiSig", false, false, "Create MultiSig Tx for NFTs [Approve-token] - Happy path -  commit.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 1, internalTx3_approveToken_00}, "MEMO : xxxxx", nil},


			//OK=========================================== start : case-1.1 mint-item (CREATE-TX, SIGN-TX, DELETE-TX, TRANSFER-OWNERSHIP-TX)===========================================
			//Remarks : before transfer is [multisig-acc-1], after transfer is [multisig-acc-2]

			//Mint-item : multiSig-create & sign
			{"multiSig", false, false, "Create MultiSig Tx for NFTs [Mint-item] - Happy path -  commit.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 6, internalTx3_mintItem_00}, "MEMO : xxxxx", nil},
			{"multiSig", true, true, "Sign MultiSig Tx for NFTs [Mint-item] - Errr, due to All signers must pass kyc.																							 ", "multisig-acc-no-kyc", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-1", 6, internalTx3_mintItem_00}, "MEMO : xxxxx", nil},
			{"multiSig", true, true, "Re-sign MultiSig Tx for NFTs [Mint-item] - Error for counter+0, due to already signed by multisig-acc-1 while Create MultiSig Tx for NFTs [Mint-item].", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-1", 6, internalTx3_mintItem_00}, "MEMO : xxxxx", nil},

			//Mint-item : multiSig-delete
			{"multiSig", true, true, "Delete MultiSig Tx for NFTs [Mint-item] -  Error, due to Group address invalid.                         ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-4", 6, internalTx3_mintItem_00}, "MEMO : Delete MultiSig Tx", nil},
			{"multiSig", true, true, "Delete MultiSig Tx for NFTs [Mint-item] -  Error, Invalid Owner address due to Only group account owner can remove pending tx.                         ", "multisig-acc-2", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-2", "", 0, nil, "grp-addr-1", 6, internalTx3_mintItem_00}, "MEMO : Delete MultiSig Tx", nil},
			{"multiSig", true, true, "Delete MultiSig Tx for NFTs [Mint-item] -  Error, due to Pending tx is not found which ID : 33.       ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-1", 33, internalTx3_mintItem_00}, "MEMO : Delete MultiSig Tx", nil},

			//Mint-item : multiSig-transfer-ownership
			{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [Mint-item] -  Error, due to Group address invalid.                                                                ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-4", 6, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
			{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [Mint-item] -  Error, due to Owner of group address invalid.                                                       ", "multisig-acc-3", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-3", "multisig-acc-1", 0, nil, "grp-addr-1", 6, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
			{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [Mint-item] -  Error, due to without KYC																																					 ", "multisig-acc-no-kyc", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-no-kyc", 0, nil, "grp-addr-1", 6, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
			{"multiSig", false, false, "Transfer MultiSig Owner for NFTs [Mint-item] -  [from multisig-acc-1 to multisig-acc-2] - Happy Path - commit.                                    ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-1", 6, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
			{"multiSig", true, true, "Re-transfer MultiSig Owner for NFTs [Mint-item] -  Error, due to Owner of group address invalid as MultiSig-account already been transfer to others.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-1", 6, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //

			//OK=========================================== start : case-1.1 burn-item (CREATE-TX, SIGN-TX, DELETE-TX, TRANSFER-OWNERSHIP-TX)===========================================
			//Remarks : before transfer is [multisig-acc-2], after transfer is [multisig-acc-1]

			//Burn-item : multiSig-create & sign
			{"multiSig", false, false, "Create MultiSig Tx for NFTs [Burn-item] - Happy path -  commit.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 7, internalTx3_burnItem_00}, "MEMO : xxxxx", nil},
			{"multiSig", true, true, "Sign MultiSig Tx for NFTs [Burn-item] - Errr, due to All signers must pass kyc.																							 ", "multisig-acc-no-kyc", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-1", 7, internalTx3_burnItem_00}, "MEMO : xxxxx", nil},
			{"multiSig", true, true, "Re-sign MultiSig Tx for NFTs [Burn-item] - Error for counter+0, due to already signed by multisig-acc-1 while Create MultiSig Tx for NFTs [Burn-item].", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-1", 7, internalTx3_burnItem_00}, "MEMO : xxxxx", nil},

			//Burn-item : multiSig-delete
			{"multiSig", true, true, "Delete MultiSig Tx for NFTs [Burn-item] -  Error, due to Group address invalid.                         ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-4", 7, internalTx3_burnItem_00}, "MEMO : Delete MultiSig Tx", nil},
			///-- start : due to above [//------------------------------- start : mint-item]
			{"multiSig", true, true, "Delete MultiSig Tx for NFTs [Burn-item] -  Error, Invalid Owner address due to Only group account owner can remove pending tx.                         ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-1", 7, internalTx3_burnItem_00}, "MEMO : Delete MultiSig Tx", nil},
			{"multiSig", true, true, "Delete MultiSig Tx for NFTs [Burn-item] -  Error, due to Pending tx is not found which ID : 37.      ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-1", 37, internalTx3_burnItem_00}, "MEMO : Delete MultiSig Tx", nil},

			//Burn-item : multiSig-transfer-ownership
			{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [Burn-item] -  Error, due to Group address invalid.                                                                ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-4", 7, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
			{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [Burn-item] -  Error, due to Owner of group address invalid.                                                       ", "multisig-acc-3", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-3", "multisig-acc-1", 0, nil, "grp-addr-1", 7, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
			{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [Burn-item] -  Error, due to without KYC																																					 ", "multisig-acc-no-kyc", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-no-kyc", 0, nil, "grp-addr-1", 7, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
			///--- start : due to above [//------------------------------- start : mint-item]
			{"multiSig", false, false, "Transfer MultiSig Owner for NFTs [Burn-item] -  [from multisig-acc-2 to multisig-acc-1] - Happy Path - commit.                                    ", "multisig-acc-2", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-2", "multisig-acc-1", 0, nil, "grp-addr-1", 7, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
			///--- start : due to above [//------------------------------- start : mint-item]
			{"multiSig", true, true, "Re-transfer MultiSig Owner for NFTs [Burn-item] -  Error, due to Owner of group address invalid as MultiSig-account already been transfer to others.", "multisig-acc-2", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-2", "multisig-acc-1", 0, nil, "grp-addr-1", 7, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //

			//OK=========================================== start : case-1.1 Transfer-ownership of Token (CREATE-TX, SIGN-TX, DELETE-TX, TRANSFER-OWNERSHIP-TX)===========================================
			//Remarks : before transfer is [multisig-acc-1], after transfer is [multisig-acc-2]

			// Transfer-ownership of Token : multiSig-create & sign
			{"multiSig", false, false, "Create MultiSig Tx for NFTs [transfer-token-ownership] - Happy path -  commit.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-1", 8, internalTx3_transferTokenOwnership_00}, "MEMO : xxxxx", nil},
			{"multiSig", true, true, "Sign MultiSig Tx for NFTs [transfer-token-ownership] - Errr, due to All signers must pass kyc.																							 ", "multisig-acc-no-kyc", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-1", 8, internalTx3_transferTokenOwnership_00}, "MEMO : xxxxx", nil},
			{"multiSig", true, true, "Re-sign MultiSig Tx for NFTs [transfer-token-ownership] - Error for counter+0, due to already signed by multisig-acc-1 while Create MultiSig Tx for NFTs [transfer-token-ownership].", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-1", 8, internalTx3_transferTokenOwnership_00}, "MEMO : xxxxx", nil},

			//transfer-token-ownership : multiSig-delete
			{"multiSig", true, true, "Delete MultiSig Tx for NFTs [transfer-token-ownership] -  Error, due to Group address invalid.                         ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-4", 8, internalTx3_transferTokenOwnership_00}, "MEMO : Delete MultiSig Tx", nil},
			///--- start : due to above [//------------------------------- start : burn-item]
			{"multiSig", true, true, "Delete MultiSig Tx for NFTs [transfer-token-ownership] -  Error, Invalid Owner address due to Only group account owner can remove pending tx.                         ", "multisig-acc-2", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-2", "", 0, nil, "grp-addr-1", 8, internalTx3_transferTokenOwnership_00}, "MEMO : Delete MultiSig Tx", nil},
			{"multiSig", true, true, "Delete MultiSig Tx for NFTs [transfer-token-ownership] -  Error, due to Pending tx is not found which ID : 41.      ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-1", 41, internalTx3_transferTokenOwnership_00}, "MEMO : Delete MultiSig Tx", nil},

			//transfer-token-ownership : multiSig-transfer-ownership
			{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [transfer-token-ownership] -  Error, due to Group address invalid.                                                                ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-4", 8, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
			{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [transfer-token-ownership] -  Error, due to Owner of group address invalid.                                                       ", "multisig-acc-3", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-3", "multisig-acc-1", 0, nil, "grp-addr-1", 8, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
			{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [transfer-token-ownership] -  Error, due to without KYC																																					 ", "multisig-acc-no-kyc", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-no-kyc", 0, nil, "grp-addr-1", 8, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
			///--- start : due to above [//------------------------------- start : mint-item]
			{"multiSig", false, false, "Transfer MultiSig Owner for NFTs [transfer-token-ownership] -  [from multisig-acc-1 to multisig-acc-2] - Happy Path - commit.                                    ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-1", 8, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
			///--- start : due to above [//------------------------------- start : mint-item]
			{"multiSig", true, true, "Re-transfer MultiSig Owner for NFTs [transfer-token-ownership] -  Error, due to Owner of group address invalid as MultiSig-account already been transfer to others.", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-1", 8, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //

			//OK=========================================== start : case-1.1 Accept-ownership of Token (CREATE-TX, SIGN-TX, DELETE-TX, TRANSFER-OWNERSHIP-TX)===========================================
			//Remarks : before transfer is [multisig-acc-2]

			// Accept-ownership of Token : multiSig-create & sign
			{"multiSig", false, false, "Create MultiSig Tx for NFTs [accept-token-ownership] - Happy path -  commit.", "multisig-acc-2", "100000000cin", 0, MultisigInfo{"create-internal-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx3_acceptTokenOwnership_00}, "MEMO : xxxxx", nil},
			{"multiSig", true, true, "Sign MultiSig Tx for NFTs [accept-token-ownership] - Errr, due to All signers must pass kyc.																							 ", "multisig-acc-no-kyc", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx3_acceptTokenOwnership_00}, "MEMO : xxxxx", nil},
			{"multiSig", true, true, "Re-sign MultiSig Tx for NFTs [accept-token-ownership] - Error for counter+0, due to already signed by multisig-acc-2 while Create MultiSig Tx for NFTs [accept-token-ownership].", "multisig-acc-2", "100000000cin", 0, MultisigInfo{"multiSig-sign-tx", "", "", 0, nil, "grp-addr-2", 0, internalTx3_acceptTokenOwnership_00}, "MEMO : xxxxx", nil},

			//Accept-token-ownership : multiSig-delete
			{"multiSig", true, true, "Delete MultiSig Tx for NFTs [accept-token-ownership] -  Error, due to Group address invalid.                         ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-4", 0, internalTx3_acceptTokenOwnership_00}, "MEMO : Delete MultiSig Tx", nil},
			///--- start : due to above [//------------------------------- start : Transfer-ownership of Token]
			{"multiSig", true, true, "Delete MultiSig Tx for NFTs [accept-token-ownership] -  Error, Invalid Owner address due to Only group account owner can remove pending tx.                         ", "multisig-acc-4", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-4", "", 0, nil, "grp-addr-2", 0, internalTx3_acceptTokenOwnership_00}, "MEMO : Delete MultiSig Tx", nil},
			{"multiSig", true, true, "Delete MultiSig Tx for NFTs [accept-token-ownership] -  Error, due to Pending tx is not found which ID : 42.      ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"multiSig-delete-tx", "multisig-acc-1", "", 0, nil, "grp-addr-2", 42, internalTx3_acceptTokenOwnership_00}, "MEMO : Delete MultiSig Tx", nil},

			//Accept-token-ownership : multiSig-transfer-ownership
			{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [accept-token-ownership] -  Error, due to Group address invalid.                                                                ", "multisig-acc-1", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-2", 0, nil, "grp-addr-4", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
			{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [accept-token-ownership] -  Error, due to Owner of group address invalid.                                                       ", "multisig-acc-3", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-3", "multisig-acc-1", 0, nil, "grp-addr-2", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //
			{"multiSig", true, true, "Transfer MultiSig Owner for NFTs [accept-token-ownership] -  Error, due to without KYC																																					 ", "multisig-acc-no-kyc", "100000000cin", 0, MultisigInfo{"transfer-ownership", "multisig-acc-1", "multisig-acc-no-kyc", 0, nil, "grp-addr-2", 0, nil}, "MEMO : Transfer MultiSig Owner.", nil}, //

		*/

		//====================start : case-xxxx
		// signer without through KYC
		{"multiSig", true, true, "Create MultiSig Account - Error, due to without KYC            ", "multisig-acc-no-kyc", "800400000cin", 0, MultisigInfo{"create", "multisig-acc-no-kyc", "", 2, []string{"multisig-acc-1", "multisig-acc-no-kyc"}, "", 0, nil}, "", nil},
	}

	return tcs
}
