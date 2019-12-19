package fungible

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	multiStore "github.com/cosmos/cosmos-sdk/store/types"
	dbm "github.com/tendermint/tm-db"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/fee"
)

// Due to how the framework is set, most of the basic input validation is in the Validate() methods
// That means that all of the data that comes to Keeper is already validated, and we do not
// test it here

var (
	delPk1   = ed25519.GenPrivKey().PubKey()
	delPk2   = ed25519.GenPrivKey().PubKey()
	delPk3   = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdkTypes.AccAddress(delPk1.Address())
	delAddr2 = sdkTypes.AccAddress(delPk2.Address())
	delAddr3 = sdkTypes.AccAddress(delPk3.Address())

	// TODO move to common testing package for all modules
	// test addresses
	TestAddrs = []sdkTypes.AccAddress{
		delAddr1, delAddr2, delAddr3,
	}
)

func TestTokenKeeper(t *testing.T) {
	var tests = []struct {
		Name string
	}{}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
		})
	}
}

func MakeTestCodec() *codec.Codec {

	var cdc = codec.New()

	RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdkTypes.RegisterCodec(cdc)

	codec.RegisterCrypto(cdc)
	return cdc

}

func defaultContext(t *testing.T, tokenKey sdkTypes.StoreKey, feeKey sdkTypes.StoreKey, keyAcc sdkTypes.StoreKey,
	keyParams sdkTypes.StoreKey, tkeyParams sdkTypes.StoreKey) multiStore.CommitMultiStore {

	sdkTypes.GetConfig().SetBech32PrefixForAccount(types.Bech32PrefixAccAddr, types.Bech32PrefixAccPub)
	sdkTypes.GetConfig().SetBech32PrefixForValidator(types.Bech32PrefixValAddr, types.Bech32PrefixValPub)
	sdkTypes.GetConfig().SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr, types.Bech32PrefixConsPub)

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)

	cms.MountStoreWithDB(tokenKey, sdkTypes.StoreTypeIAVL, nil)
	cms.MountStoreWithDB(keyAcc, sdkTypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(feeKey, sdkTypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(keyParams, sdkTypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkeyParams, sdkTypes.StoreTypeTransient, db)

	err := cms.LoadLatestVersion()
	require.Nil(t, err)

	return cms
}

func PrepareTest(t *testing.T) (sdkTypes.Context, *Keeper) {

	//1. Getting default codec for marshaling and unmarshaling
	cdc := MakeTestCodec()

	//2. Create key store for fee-keeper, token-keeper, auth.StoreKey, params.StoreKey
	tokenKey := sdkTypes.NewKVStoreKey("token")
	feeKey := sdkTypes.NewKVStoreKey("fee")
	keyAcc := sdkTypes.NewKVStoreKey(auth.StoreKey)
	keyParams := sdkTypes.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdkTypes.NewTransientStoreKey(params.TStoreKey)

	//3. Getting context for fee
	cms := defaultContext(t, tokenKey, feeKey, keyAcc, keyParams, tkeyParams)
	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ctx := sdkTypes.NewContext(cms, abci.Header{ChainID: "foochainid"}, false, log.NewNopLogger())

	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, nil)
	feeKeeper := fee.NewKeeper(cdc, feeKey)

	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)

	//4. Creating instance base on fee-keeper, account-keeper, bank-keeper instance
	keeper := NewKeeper(cdc, &accountKeeper, &feeKeeper, tokenKey)
	amt, _ := sdkTypes.NewIntFromString("100000000000000000000000000000000000000")
	initCoins := sdkTypes.NewCoins(sdkTypes.NewCoin("cin", amt))
	feeKeeper.SetFeeCollectorAddresses(ctx, "token", []sdkTypes.AccAddress{delAddr1})

	processFeeSetting(t, ctx, feeKeeper)

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	// set the account as well
	for _, addr := range TestAddrs {

		addracc := keeper.accountKeeper.NewAccountWithAddress(ctx, addr)
		addracc.SetAccountNumber(keeper.accountKeeper.GetNextAccountNumber(ctx))
		keeper.accountKeeper.SetAccount(ctx, addracc)
		_, err := bankKeeper.AddCoins(ctx, addracc.GetAddress(), initCoins)
		fmt.Println(addracc.GetCoins())

		require.Nil(t, err)
	}

	return ctx, &keeper

}

func processFeeSetting(t *testing.T, ctx sdkTypes.Context, keeper fee.Keeper) {

	/// ------------------------------------case1: FeeSettingExists ?
	var testFeeSettingTypes = []string{
		"default",
		"zero",
	}

	for i, testFeeSettingType := range testFeeSettingTypes {
		rsIsExist := keeper.FeeSettingExists(ctx, testFeeSettingType)
		assert.False(t, rsIsExist)
		fmt.Printf("============\nTest FeeSetting : %d with value : %s \n", i+1, testFeeSettingType)

	}

	/// ------------------------------------case2: ADD -> GET -> Existed
	var testIssuerAddrs = []string{
		"mxw1yw6mg7fty4mzcwupvzek53x5egm7tp2ldwaxq3",
		"mxw1yyz3h9calxmvjp4x05nnn70a8ex7fee3th7r7k",
	}

	rsIssuerAddrs1, _ := sdkTypes.AccAddressFromBech32(testIssuerAddrs[0])
	rsIssuerAddrs2, _ := sdkTypes.AccAddressFromBech32(testIssuerAddrs[1])

	var a []sdkTypes.AccAddress
	a = append(a, rsIssuerAddrs1)
	a = append(a, rsIssuerAddrs2)

	keeper.SetAuthorisedAddresses(ctx, a)
	rsAddGet := keeper.GetAuthorisedAddresses(ctx)

	if assert.NotNil(t, rsAddGet) {

		// 1. check if is authorised
		for j, item := range a {
			assert.Contains(t, rsAddGet, item)
			fmt.Printf("============\nAfter SetAuthorisedAddresses, Test Address : %d with value : %s \n", j+1, item)

			rsAuthorisedAddr := keeper.IsAuthorised(ctx, item)
			assert.True(t, rsAuthorisedAddr)

		}

		// 2. start Create FeeSetting
		processCreateFeeSetting(t, ctx, keeper, rsIssuerAddrs1, rsIssuerAddrs2)
		feeSettings := keeper.ListAllSysFeeSetting(ctx)
		fmt.Printf("\n============feeSettings : %v\n", feeSettings)

	}

}

func processCreateFeeSetting(t *testing.T, ctx sdkTypes.Context, keeper fee.Keeper, rsIssuerAddrs1 sdkTypes.AccAddress, rsIssuerAddrs2 sdkTypes.AccAddress) {

	mindefault := sdkTypes.Coins{
		{
			Denom:  "cin",
			Amount: sdkTypes.NewInt(100000000),
		},
	}
	maxdefault := sdkTypes.Coins{
		{
			Denom:  "cin",
			Amount: sdkTypes.NewInt(1000000000),
		},
	}
	minzero := sdkTypes.Coins{
		{
			Denom:  "cin",
			Amount: sdkTypes.NewInt(0),
		},
	}
	maxzero := sdkTypes.Coins{
		{
			Denom:  "cin",
			Amount: sdkTypes.NewInt(0),
		},
	}

	var caseFeeSettings = []fee.MsgSysFeeSetting{
		{
			Name:       "default",
			Min:        mindefault,
			Max:        maxdefault,
			Percentage: "0.5",
			Issuer:     rsIssuerAddrs1,
		},
		{
			Name:       "zero",
			Min:        minzero,
			Max:        maxzero,
			Percentage: "0",
			Issuer:     rsIssuerAddrs2,
		},
	}

	for i, caseFeeSetting := range caseFeeSettings {
		keeper.CreateFeeSetting(ctx, caseFeeSetting)
		rsAddGet, _ := keeper.GetFeeSettingByName(ctx, caseFeeSetting.Name)

		if assert.NotNil(t, rsAddGet) {

			fmt.Printf("============\nAfter CreateFeeSetting, Test FeeSetting : %d with value : %v \n", i+1, rsAddGet.Name)

			rsIsExist := keeper.FeeSettingExists(ctx, rsAddGet.Name)
			assert.True(t, rsIsExist)

		}

	}

}

func TestInitGenesisAndApproverCheck(t *testing.T) {
	ctx, keeper := PrepareTest(t)

	approver1, err := sdkTypes.AccAddressFromBech32("mxw1zzrutc6x9kc7ttafwawv9ve3jm89zxpeedl8ap")
	if err != nil {
		t.Fatal(err)
	}

	approver2, err := sdkTypes.AccAddressFromBech32("mxw1jdpaz5fppzt726w44hx4q9yerlfc3ldgpj83fe")
	if err != nil {
		t.Fatal(err)
	}

	genesisState := &GenesisState{
		AuthorizedAddresses: []sdkTypes.AccAddress{approver1, approver2},
	}

	InitGenesis(ctx, keeper, *genesisState)

	if len(genesisState.AuthorizedAddresses) != 2 {
		t.Fatalf("Has %d approvers. Expected 2", len(genesisState.AuthorizedAddresses))
	}

	if !genesisState.AuthorizedAddresses[0].Equals(approver1) {
		t.Fatalf("Invalid first approver. Got %s", genesisState.AuthorizedAddresses[0])
	}

	if !genesisState.AuthorizedAddresses[1].Equals(approver2) {
		t.Fatalf("Invalid second approver. Got %s", genesisState.AuthorizedAddresses[1])
	}

	if !keeper.IsAuthorised(ctx, approver1) {
		t.Fatalf("Expected %s to be approver", approver1.String())
	}

	if !keeper.IsAuthorised(ctx, approver2) {
		t.Fatalf("Expected %s to be approver", approver2.String())
	}

}

func TestCreateFungibleFungibleToken(t *testing.T) {

	ctx, keeper := PrepareTest(t)

	approver1, err := sdkTypes.AccAddressFromBech32("mxw1zzrutc6x9kc7ttafwawv9ve3jm89zxpeedl8ap")
	if err != nil {
		t.Fatal(err)
	}

	approver2, err := sdkTypes.AccAddressFromBech32("mxw1jdpaz5fppzt726w44hx4q9yerlfc3ldgpj83fe")
	if err != nil {
		t.Fatal(err)
	}

	approver1Acc := keeper.accountKeeper.NewAccountWithAddress(ctx, approver1)
	approver1Acc.SetAccountNumber(keeper.accountKeeper.GetNextAccountNumber(ctx))
	keeper.accountKeeper.SetAccount(ctx, approver1Acc)

	genesisState := &GenesisState{
		AuthorizedAddresses: []sdkTypes.AccAddress{approver1, approver2},
	}

	var applicationFee = struct {
		To    sdkTypes.AccAddress `json:"to"`
		Value string              `json:"value"`
	}{
		To:    delAddr1,
		Value: "10",
	}

	var testCases = []struct {
		Flags        types.Bitmask
		Name         string
		Symbol       string
		Decimals     int
		Owner        sdkTypes.AccAddress
		NewOwner     sdkTypes.AccAddress
		Metadata     string
		TotalSupply  sdkTypes.Uint
		MaxSupply    sdkTypes.Uint
		ExpectedCode sdkTypes.CodeType
	}{
		{
			Flags:        FixedSupplyNotBurnableFungibleTokenMask, // FungibleFlag
			Name:         "Happy path fixed supply - Not Burnable Fungible Token",
			Symbol:       "TST-1",
			Decimals:     18,
			Owner:        delAddr1,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(500000000),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
		},
		{
			Flags:       FixedSupplyBurnableFungibleTokenMask, // FungibleFlag + BurnFlag
			Name:        "Happy path fixed supply - Burnable Fungible Token",
			Symbol:      "TST-2",
			Decimals:    18,
			Owner:       delAddr1,
			TotalSupply: sdkTypes.NewUint(0),
			MaxSupply:   sdkTypes.NewUint(500000000),
			Metadata:    "ipfs-hash-link",
		},
		{
			Flags:        DynamicFungibleTokenMask, // FungibleFlag + MintFlag + BurnFlag
			Name:         "Happy path dynamic supply",
			Symbol:       "tsttttt",
			Decimals:     18,
			Owner:        delAddr2,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			fmt.Println(testCase)
			InitGenesis(ctx, keeper, *genesisState)

			var fixedSupply bool
			fixedSupply = true
			if testCase.Flags.HasFlag(DynamicFungibleTokenMask) {
				fixedSupply = false
			}
			res := keeper.CreateFungibleToken(ctx, testCase.Name, testCase.Symbol, testCase.Decimals, testCase.Owner,
				fixedSupply, testCase.MaxSupply, testCase.Metadata, applicationFee)
			//fmt.Println(res)

			if res.Code != testCase.ExpectedCode {
				t.Fatalf("Expected code %d. Got code %d", testCase.ExpectedCode, res.Code)
			}

			if !keeper.TokenExists(ctx, testCase.Symbol) {
				t.Fatal("After creating, token does not exist")
			}

			var token = new(Token)
			keeper.mustGetTokenData(ctx, testCase.Symbol, token)

			if !token.Owner.Equals(testCase.Owner) {
				t.Fatalf("Expected %s owner. Got %s", testCase.Owner.String(), token.Owner.String())
			}

			if token.Symbol != testCase.Symbol {
				t.Fatalf("Expected %s id. Got %s", testCase.Symbol, token.Symbol)
			}

			if token.Flags.HasFlag(FixedSupplyBurnableFungibleTokenMask) != testCase.Flags.HasFlag(FixedSupplyBurnableFungibleTokenMask) {
				t.Fatalf("Expected %v fixed supply. Got %v", testCase.Flags.HasFlag(FixedSupplyBurnableFungibleTokenMask), token.Flags.HasFlag(FixedSupplyBurnableFungibleTokenMask))
			}

			var expectedTotalSupply sdkTypes.Uint
			if !testCase.Flags.HasFlag(DynamicFungibleTokenMask) {
				expectedTotalSupply = testCase.TotalSupply
			} else {
				expectedTotalSupply = testCase.TotalSupply
			}

			if !token.TotalSupply.Equal(expectedTotalSupply) {
				t.Fatalf("Expected %v total supply. Got %v", expectedTotalSupply.String(), token.TotalSupply.String())
			}

			if token.Flags.HasFlag(ApprovedFlag) {
				t.Fatal("Token must start as unapproved")
			}

			if token.Flags.HasFlag(FrozenFlag) {
				t.Fatalf("Token must start as non-frozen")
			}

			if token.Metadata != testCase.Metadata {
				t.Fatalf("Expected %s data link. Got %s", testCase.Metadata, token.Metadata)
			}

			if res := keeper.CreateFungibleToken(ctx, testCase.Name, testCase.Symbol, testCase.Decimals, testCase.Owner, fixedSupply, testCase.TotalSupply, testCase.Metadata, applicationFee); res.Code != types.CodeTokenDuplicated {
				t.Fatalf("Expected error when creating duplicate asset class")
			}
		})
	}
}

func TestApproveToken(t *testing.T) {

	ctx, keeper := PrepareTest(t)

	approver1, err := sdkTypes.AccAddressFromBech32("mxw1zzrutc6x9kc7ttafwawv9ve3jm89zxpeedl8ap")
	if err != nil {
		t.Fatal(err)
	}

	approver2, err := sdkTypes.AccAddressFromBech32("mxw1jdpaz5fppzt726w44hx4q9yerlfc3ldgpj83fe")
	if err != nil {
		t.Fatal(err)
	}

	approver1Acc := keeper.accountKeeper.NewAccountWithAddress(ctx, approver1)
	approver1Acc.SetAccountNumber(keeper.accountKeeper.GetNextAccountNumber(ctx))
	keeper.accountKeeper.SetAccount(ctx, approver1Acc)

	genesisState := &GenesisState{
		AuthorizedAddresses: []sdkTypes.AccAddress{approver1, approver2},
	}

	var applicationFee = struct {
		To    sdkTypes.AccAddress `json:"to"`
		Value string              `json:"value"`
	}{
		To:    delAddr2,
		Value: "100000",
	}

	var tokenFee1 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "zero",
		},
		{
			Action:  "mint",
			FeeName: "zero",
		},
		{
			Action:  "burn",
			FeeName: "zero",
		},
		{
			Action:  "transferOwnership",
			FeeName: "zero",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "zero",
		},
	}

	var tokenFee2 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "default",
		},
		{
			Action:  "mint",
			FeeName: "default",
		},
		{
			Action:  "burn",
			FeeName: "default",
		},
		{
			Action:  "transferOwnership",
			FeeName: "default",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "default",
		},
	}

	var testCases = []struct {
		Flags        types.Bitmask
		Name         string
		Symbol       string
		Decimals     int
		Owner        sdkTypes.AccAddress
		NewOwner     sdkTypes.AccAddress
		Metadata     string
		TotalSupply  sdkTypes.Uint
		MaxSupply    sdkTypes.Uint
		TokenFees    []TokenFee
		ExpectedCode sdkTypes.CodeType
	}{
		{
			Flags:        FixedSupplyNotBurnableFungibleTokenMask, // FungibleFlag
			Name:         "Happy path fixed supply - Not Burnable Fungible Token",
			Symbol:       "TST-1",
			Decimals:     18,
			Owner:        delAddr1,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			TokenFees:    tokenFee1,
			ExpectedCode: sdkTypes.CodeOK,
		},
		{
			Flags:        FixedSupplyBurnableFungibleTokenMask, // FungibleFlag + BurnFlag
			Name:         "Happy path fixed supply - Burnable Fungible Token",
			Symbol:       "TST-2",
			Decimals:     18,
			Owner:        delAddr1,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			TokenFees:    tokenFee1,
			ExpectedCode: sdkTypes.CodeOK,
		},
		{
			Flags:        DynamicFungibleTokenMask, // FungibleFlag + MintFlag + BurnFlag
			Name:         "Happy path dynamic supply",
			Symbol:       "tsttttt",
			Decimals:     18,
			Owner:        delAddr1,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			TokenFees:    tokenFee2,
			ExpectedCode: sdkTypes.CodeOK,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			InitGenesis(ctx, keeper, *genesisState)

			var fixedSupply bool
			fixedSupply = true
			if testCase.Flags.HasFlag(DynamicFungibleTokenMask) {
				fixedSupply = false
			}

			res := keeper.CreateFungibleToken(ctx, testCase.Name, testCase.Symbol, testCase.Decimals, testCase.Owner,
				fixedSupply, testCase.MaxSupply, testCase.Metadata, applicationFee)
			fmt.Println(testCase)
			if res.Code != testCase.ExpectedCode {
				t.Fatalf("Expected code %d. Got code %d", testCase.ExpectedCode, res.Code)
			}

			var burnable bool
			burnable = false
			if testCase.Flags.HasFlag(BurnFlag) {
				burnable = true
			}

			resApproval := keeper.ApproveToken(ctx, testCase.Symbol, testCase.TokenFees, burnable, approver1, "ipfs-hash-link")
			if resApproval.Code != testCase.ExpectedCode {
				t.Fatalf("Expected approval code %d. Got %d", testCase.ExpectedCode, resApproval.Code)
			}

			var fungibleToken = new(Token)
			keeper.mustGetTokenData(ctx, testCase.Symbol, fungibleToken)

			if testCase.ExpectedCode == sdkTypes.CodeOK {
				if !fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("Token not approved after successful approval")
				}

				if testCase.Metadata != fungibleToken.Metadata {
					t.Fatalf("Expected approval metada %s. Got %s", testCase.Metadata, fungibleToken.Metadata)
				}

				account := keeper.getFungibleAccount(ctx, testCase.Symbol, testCase.Owner)
				if account == nil {
					t.Fatal("After approval issuer account does not exist")
				}

				if fixedSupply {
					if !account.Balance.Equal(testCase.TotalSupply) {
						t.Fatalf("Expected token owner to have %v balance. Got %v", testCase.TotalSupply.String(), account.Balance.String())
					}
				} else {
					if !account.Balance.Equal(testCase.TotalSupply) {
						t.Fatalf("Expected token owner to have %v balance. Got %v", testCase.TotalSupply.String(), account.Balance.String())
					}
				}

				if !account.Balance.Equal(fungibleToken.TotalSupply) {
					t.Fatalf("After approving, issuer balance %v not equal to total supply %v", account.Balance.String(), fungibleToken.TotalSupply.String())
				}
				if account.Frozen {
					t.Fatal("Issuer started as frozen")
				}

				resApprovalAgain := keeper.ApproveToken(ctx, testCase.Symbol, testCase.TokenFees, burnable, approver1, testCase.Metadata)
				if resApprovalAgain.Code != types.CodeTokenApproved {
					t.Fatal("Did double approval")
				}

				// Trying to reject after approval - not sure in which test suite it belongs
				resRejectal := keeper.RejectToken(ctx, testCase.Symbol, approver1)
				if resRejectal.Code == sdkTypes.CodeOK {
					t.Fatal("Managed to reject after approval")
				}
			} else {
				if fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("token approved after unsucessful approval")
				}

				if fungibleToken.Metadata != "" {
					t.Fatal("Fungible token approval metada not empty after unsuccessful approval")
				}
			}
		})
	}
}

func TestRejectAssetClass(t *testing.T) {
	ctx, keeper := PrepareTest(t)

	approver1, err := sdkTypes.AccAddressFromBech32("mxw1zzrutc6x9kc7ttafwawv9ve3jm89zxpeedl8ap")
	if err != nil {
		t.Fatal(err)
	}

	approver2, err := sdkTypes.AccAddressFromBech32("mxw1jdpaz5fppzt726w44hx4q9yerlfc3ldgpj83fe")
	if err != nil {
		t.Fatal(err)
	}

	approver1Acc := keeper.accountKeeper.NewAccountWithAddress(ctx, approver1)
	approver1Acc.SetAccountNumber(keeper.accountKeeper.GetNextAccountNumber(ctx))
	keeper.accountKeeper.SetAccount(ctx, approver1Acc)

	genesisState := &GenesisState{
		AuthorizedAddresses: []sdkTypes.AccAddress{approver1, approver2},
	}

	var applicationFee = struct {
		To    sdkTypes.AccAddress `json:"to"`
		Value string              `json:"value"`
	}{
		To:    delAddr1,
		Value: "100000",
	}

	var testCases = []struct {
		Flags        types.Bitmask
		Name         string
		Symbol       string
		Decimals     int
		Owner        sdkTypes.AccAddress
		NewOwner     sdkTypes.AccAddress
		Metadata     string
		TotalSupply  sdkTypes.Uint
		MaxSupply    sdkTypes.Uint
		ExpectedCode sdkTypes.CodeType
	}{
		{
			Flags:        FixedSupplyNotBurnableFungibleTokenMask, // FungibleFlag
			Name:         "Happy path fixed supply - Not Burnable Fungible Token",
			Symbol:       "TST-1",
			Decimals:     18,
			Owner:        delAddr2,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
		},
		{
			Flags:        FixedSupplyBurnableFungibleTokenMask, // FungibleFlag + BurnFlag
			Name:         "Happy path fixed supply - Burnable Fungible Token",
			Symbol:       "TST-2",
			Decimals:     18,
			Owner:        delAddr1,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
		},
		{
			Flags:        DynamicFungibleTokenMask, // FungibleFlag + MintFlag + BurnFlag
			Name:         "Happy path dynamic supply",
			Symbol:       "tsttttt",
			Decimals:     18,
			Owner:        delAddr2,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {

			InitGenesis(ctx, keeper, *genesisState)

			var fixedSupply bool
			fixedSupply = true
			if testCase.Flags.HasFlag(DynamicFungibleTokenMask) {
				fixedSupply = false
			}

			res := keeper.CreateFungibleToken(ctx, testCase.Name, testCase.Symbol, testCase.Decimals, testCase.Owner,
				fixedSupply, testCase.TotalSupply, testCase.Metadata, applicationFee)
			if res.Code != testCase.ExpectedCode {
				t.Fatalf("Expected code %d. Got code %d", testCase.ExpectedCode, res.Code)
			}

			resRejectal := keeper.RejectToken(ctx, testCase.Symbol, approver1)
			if resRejectal.Code != testCase.ExpectedCode {
				t.Fatalf("Expected approval code %d. Got %d", testCase.ExpectedCode, resRejectal.Code)
			}

			if resRejectal.Code == sdkTypes.CodeOK {
				if keeper.TokenExists(ctx, testCase.Symbol) {
					t.Fatal("Token still exists after rejectal")
				}
			} else {
				var fungibleToken = new(Token)
				keeper.mustGetTokenData(ctx, testCase.Symbol, fungibleToken)

				if fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("Accidentally approved")
				}
			}
		})
	}
}

func TestFreezeFungibleToken(t *testing.T) {
	ctx, keeper := PrepareTest(t)

	approver1, err := sdkTypes.AccAddressFromBech32("mxw1zzrutc6x9kc7ttafwawv9ve3jm89zxpeedl8ap")
	if err != nil {
		t.Fatal(err)
	}

	approver2, err := sdkTypes.AccAddressFromBech32("mxw1jdpaz5fppzt726w44hx4q9yerlfc3ldgpj83fe")
	if err != nil {
		t.Fatal(err)
	}

	approver1Acc := keeper.accountKeeper.NewAccountWithAddress(ctx, approver1)
	approver1Acc.SetAccountNumber(keeper.accountKeeper.GetNextAccountNumber(ctx))
	keeper.accountKeeper.SetAccount(ctx, approver1Acc)

	genesisState := &GenesisState{
		AuthorizedAddresses: []sdkTypes.AccAddress{approver1, approver2},
	}

	var applicationFee = struct {
		To    sdkTypes.AccAddress `json:"to"`
		Value string              `json:"value"`
	}{
		To:    delAddr1,
		Value: "100000",
	}

	var testCases = []struct {
		Flags        types.Bitmask
		Name         string
		Symbol       string
		Decimals     int
		Owner        sdkTypes.AccAddress
		NewOwner     sdkTypes.AccAddress
		Metadata     string
		TotalSupply  sdkTypes.Uint
		MaxSupply    sdkTypes.Uint
		ExpectedCode sdkTypes.CodeType
	}{
		{
			Flags:        FixedSupplyNotBurnableFungibleTokenMask, // FungibleFlag
			Name:         "Happy path fixed supply",
			Symbol:       "TST-1",
			Decimals:     18,
			Owner:        delAddr3,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
		},
		{
			Flags:       FixedSupplyBurnableFungibleTokenMask, // FungibleFlag + BurnFlag
			Name:        "Happy path fixed supply - Burnable Fungible Token",
			Symbol:      "TST-2",
			Decimals:    18,
			Owner:       delAddr1,
			TotalSupply: sdkTypes.NewUint(0),
			MaxSupply:   sdkTypes.NewUint(0),
			Metadata:    "ipfs-hash-link",
		},
		{
			Flags:        DynamicFungibleTokenMask, // FungibleFlag + MintFlag + BurnFlag
			Name:         "Happy path dynamic supply",
			Symbol:       "tsttttt",
			Decimals:     18,
			Owner:        delAddr3,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {

			InitGenesis(ctx, keeper, *genesisState)

			var fixedSupply bool
			fixedSupply = true
			if testCase.Flags.HasFlag(DynamicFungibleTokenMask) {
				fixedSupply = false
			}

			res := keeper.CreateFungibleToken(ctx, testCase.Name, testCase.Symbol, testCase.Decimals, testCase.Owner,
				fixedSupply, testCase.MaxSupply, testCase.Metadata, applicationFee)
			if res.Code != testCase.ExpectedCode {
				t.Fatalf("Expected code %d. Got code %d", testCase.ExpectedCode, res.Code)
			}

			resFreeze := keeper.FreezeToken(ctx, testCase.Symbol, approver1, testCase.Metadata)
			if resFreeze.Code != testCase.ExpectedCode {
				t.Fatalf("Expected approval code %d. Got %d", testCase.ExpectedCode, resFreeze.Code)
			}

			var fungibleToken = new(Token)
			keeper.mustGetTokenData(ctx, testCase.Symbol, fungibleToken)

			if resFreeze.Code == sdkTypes.CodeOK {
				if !fungibleToken.Flags.HasFlag(FrozenFlag) {
					t.Fatal("Token not frozen after successful freeze")
				}

				if fungibleToken.Metadata != testCase.Metadata {
					t.Fatalf("Expected freeze metadata %s. Got %s", testCase.Metadata, fungibleToken.Metadata)
				}

				// If already frozen, should fail
				if res := keeper.FreezeToken(ctx, testCase.Symbol, approver1, testCase.Metadata); res.Code != 2004 {
					t.Fatal("Did not fail when freezing again")
				}
			} else {
				if fungibleToken.Flags.HasFlag(FrozenFlag) {
					t.Fatal("token frozen after unsuccessful freeze")
				}

				if fungibleToken.Metadata != "" {
					t.Fatalf("Freeze metada %s after unsuccessful freeze", fungibleToken.Metadata)
				}
			}
		})
	}
}

func TestUnfreezeToken(t *testing.T) {
	ctx, keeper := PrepareTest(t)

	approver1, err := sdkTypes.AccAddressFromBech32("mxw1zzrutc6x9kc7ttafwawv9ve3jm89zxpeedl8ap")
	if err != nil {
		t.Fatal(err)
	}

	approver2, err := sdkTypes.AccAddressFromBech32("mxw1jdpaz5fppzt726w44hx4q9yerlfc3ldgpj83fe")
	if err != nil {
		t.Fatal(err)
	}

	approver1Acc := keeper.accountKeeper.NewAccountWithAddress(ctx, approver1)
	approver1Acc.SetAccountNumber(keeper.accountKeeper.GetNextAccountNumber(ctx))
	keeper.accountKeeper.SetAccount(ctx, approver1Acc)

	genesisState := &GenesisState{
		AuthorizedAddresses: []sdkTypes.AccAddress{approver1, approver2},
	}

	var applicationFee = struct {
		To    sdkTypes.AccAddress `json:"to"`
		Value string              `json:"value"`
	}{
		To:    delAddr1,
		Value: "100000",
	}

	var testCases = []struct {
		Flags          types.Bitmask
		Name           string
		Symbol         string
		Decimals       int
		Owner          sdkTypes.AccAddress
		NewOwner       sdkTypes.AccAddress
		Metadata       string
		TotalSupply    sdkTypes.Uint
		MaxSupply      sdkTypes.Uint
		ExpectedCode   sdkTypes.CodeType
		ShouldBeFrozen bool
	}{
		{
			Flags:          FixedSupplyBurnableFungibleTokenMask, // FungibleFlag + BurnFlag
			Name:           "Happy path fixed supply - Burnable Fungible Token",
			Symbol:         "TST-1",
			Decimals:       18,
			Owner:          delAddr2,
			TotalSupply:    sdkTypes.NewUint(0),
			MaxSupply:      sdkTypes.NewUint(0),
			Metadata:       "",
			ExpectedCode:   sdkTypes.CodeOK,
			ShouldBeFrozen: true,
		},
		{
			Flags:          FixedSupplyBurnableFungibleTokenMask, // FungibleFlag + BurnFlag
			Name:           "Happy path fixed supply - Burnable Fungible Token",
			Symbol:         "TST-2",
			Decimals:       18,
			Owner:          delAddr1,
			TotalSupply:    sdkTypes.NewUint(0),
			MaxSupply:      sdkTypes.NewUint(0),
			Metadata:       "",
			ExpectedCode:   sdkTypes.CodeOK,
			ShouldBeFrozen: true,
		},
		{
			Flags:          DynamicFungibleTokenMask, // FungibleFlag + MintFlag + BurnFlag
			Name:           "Happy path dynamic supply",
			Symbol:         "tsttttt",
			Decimals:       18,
			Owner:          delAddr3,
			TotalSupply:    sdkTypes.NewUint(0),
			MaxSupply:      sdkTypes.NewUint(0),
			Metadata:       "",
			ExpectedCode:   sdkTypes.CodeOK,
			ShouldBeFrozen: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {

			InitGenesis(ctx, keeper, *genesisState)

			var fixedSupply bool
			fixedSupply = true
			if testCase.Flags.HasFlag(DynamicFungibleTokenMask) {
				fixedSupply = false
			}

			res := keeper.CreateFungibleToken(ctx, testCase.Name, testCase.Symbol, testCase.Decimals, testCase.Owner,
				fixedSupply, testCase.MaxSupply, testCase.Metadata, applicationFee)
			if res.Code != testCase.ExpectedCode {
				t.Fatalf("Expected code %d. Got code %d", testCase.ExpectedCode, res.Code)
			}

			if testCase.ShouldBeFrozen {
				resFreeze := keeper.FreezeToken(ctx, testCase.Symbol, approver1, "freeze-data-link")
				if resFreeze.Code != sdkTypes.CodeOK {
					t.Fatalf("Expected approval code %d. Got %d", testCase.ExpectedCode, resFreeze.Code)
				}
			}

			resUnfreeze := keeper.UnfreezeToken(ctx, testCase.Symbol, approver1, testCase.Metadata)

			var fungibleToken = new(Token)
			keeper.mustGetTokenData(ctx, testCase.Symbol, fungibleToken)

			if resUnfreeze.Code == sdkTypes.CodeOK {
				if fungibleToken.Flags.HasFlag(FrozenFlag) {
					t.Fatal("Token not unfrozen after successful unfreeze")
				}

				if fungibleToken.Metadata != testCase.Metadata {
					t.Fatalf("Expected unfreeze metadata %s. Got %s", testCase.Metadata, fungibleToken.Metadata)
				}

				// If already frozen, should fail
				if res := keeper.UnfreezeToken(ctx, testCase.Symbol, approver1, testCase.Metadata); res.Code != sdkTypes.CodeUnknownRequest {
					t.Fatal("Did not fail when unfreezing again")
				}
			} else {
				if testCase.ShouldBeFrozen && !fungibleToken.Flags.HasFlag(FrozenFlag) {
					t.Fatal("Token unfrozen after unsuccessful unfreeze")
				}

				if fungibleToken.Metadata != "" {
					t.Fatalf("Unfreeze metadata %s after unsuccessful unfreeze", fungibleToken.Metadata)
				}
			}
		})
	}
}

func TestTransferFungibleToken(t *testing.T) {

	fmt.Printf("============\nStart Test : %s \n", "TestTransferFungibleToken")
	ctx, keeper := PrepareTest(t)

	approver1, err := sdkTypes.AccAddressFromBech32("mxw1zzrutc6x9kc7ttafwawv9ve3jm89zxpeedl8ap")
	if err != nil {
		t.Fatal(err)
	}

	approver2, err := sdkTypes.AccAddressFromBech32("mxw1jdpaz5fppzt726w44hx4q9yerlfc3ldgpj83fe")
	if err != nil {
		t.Fatal(err)
	}

	approver1Acc := keeper.accountKeeper.NewAccountWithAddress(ctx, approver1)
	approver1Acc.SetAccountNumber(keeper.accountKeeper.GetNextAccountNumber(ctx))
	keeper.accountKeeper.SetAccount(ctx, approver1Acc)

	genesisState := &GenesisState{
		AuthorizedAddresses: []sdkTypes.AccAddress{approver1, approver2},
	}

	var applicationFee = struct {
		To    sdkTypes.AccAddress `json:"to"`
		Value string              `json:"value"`
	}{
		To:    delAddr1,
		Value: "100000",
	}

	var tokenFee1 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "zero",
		},
		{
			Action:  "mint",
			FeeName: "zero",
		},
		{
			Action:  "burn",
			FeeName: "zero",
		},
		{
			Action:  "transferOwnership",
			FeeName: "zero",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "zero",
		},
	}

	var tokenFee2 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "default",
		},
		{
			Action:  "mint",
			FeeName: "default",
		},
		{
			Action:  "burn",
			FeeName: "default",
		},
		{
			Action:  "transferOwnership",
			FeeName: "default",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "default",
		},
	}

	var testCases = []struct {
		Flags                 types.Bitmask
		Name                  string
		Symbol                string
		Decimals              int
		Owner                 sdkTypes.AccAddress
		NewOwner              sdkTypes.AccAddress
		Metadata              string
		TotalSupply           sdkTypes.Uint
		MaxSupply             sdkTypes.Uint
		TransferAmountOfToken sdkTypes.Uint
		TokenFees             []TokenFee
		ExpectedCode          sdkTypes.CodeType
	}{
		{
			Flags:                 FixedSupplyNotBurnableFungibleTokenMask, // FungibleFlag
			Name:                  "Happy path fixed supply - Not Burnable Fungible Token",
			Symbol:                "TST-1",
			Decimals:              18,
			Owner:                 delAddr3,
			NewOwner:              delAddr2,
			TotalSupply:           sdkTypes.NewUint(0),
			MaxSupply:             sdkTypes.NewUint(0),
			Metadata:              "ipfs-hash-link",
			TransferAmountOfToken: sdkTypes.NewUint(100),
			TokenFees:             tokenFee1,
			ExpectedCode:          sdkTypes.CodeOK,
		},
		{
			Flags:                 FixedSupplyBurnableFungibleTokenMask, // FungibleFlag + BurnFlag
			Name:                  "Happy path fixed supply - Burnable Fungible Token",
			Symbol:                "TST-2",
			Decimals:              18,
			Owner:                 delAddr1,
			TotalSupply:           sdkTypes.NewUint(0),
			MaxSupply:             sdkTypes.NewUint(0),
			Metadata:              "ipfs-hash-link",
			TransferAmountOfToken: sdkTypes.NewUint(100),
			TokenFees:             tokenFee1,
			ExpectedCode:          sdkTypes.CodeOK,
		},
		{
			Flags:                 DynamicFungibleTokenMask, // FungibleFlag + MintFlag + BurnFlag
			Name:                  "Happy path dynamic supply",
			Symbol:                "tsttttt",
			Decimals:              18,
			Owner:                 delAddr3,
			NewOwner:              delAddr2,
			TotalSupply:           sdkTypes.NewUint(0),
			MaxSupply:             sdkTypes.NewUint(0),
			Metadata:              "ipfs-hash-link",
			TransferAmountOfToken: sdkTypes.NewUint(10),
			TokenFees:             tokenFee2,
			ExpectedCode:          sdkTypes.CodeOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			fmt.Println(testCase)
			InitGenesis(ctx, keeper, *genesisState)

			var fixedSupply bool
			fixedSupply = true
			if testCase.Flags.HasFlag(DynamicFungibleTokenMask) {
				fixedSupply = false
			}

			resCreateToken := keeper.CreateFungibleToken(ctx, testCase.Name, testCase.Symbol, testCase.Decimals, testCase.Owner,
				fixedSupply, testCase.MaxSupply, testCase.Metadata, applicationFee)
			//fmt.Println(resCreateToken)
			if resCreateToken.Code != testCase.ExpectedCode {
				t.Fatalf("Expected code %d. Got code %d", testCase.ExpectedCode, resCreateToken.Code)
			}

			var burnable bool
			burnable = false
			if testCase.Flags.HasFlag(BurnFlag) {
				burnable = true
			}

			resApproval := keeper.ApproveToken(ctx, testCase.Symbol, testCase.TokenFees, burnable, approver1, "ipfs-hash-link")
			//fmt.Println(resApproval)
			if resApproval.Code != testCase.ExpectedCode {
				t.Fatalf("Expected approval code %d. Got %d", testCase.ExpectedCode, resApproval.Code)
			}

			var fungibleToken = new(Token)
			keeper.mustGetTokenData(ctx, testCase.Symbol, fungibleToken)

			if testCase.ExpectedCode == sdkTypes.CodeOK {
				if !fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("Token not approved after successful approval")
				}

				if testCase.Metadata != fungibleToken.Metadata {
					t.Fatalf("Expected approval metada %s. Got %s", testCase.Metadata, fungibleToken.Metadata)
				}

				account := keeper.getFungibleAccount(ctx, testCase.Symbol, testCase.Owner)
				if account == nil {
					t.Fatal("After approval issuer account does not exist")
				}

				if !account.Balance.Equal(testCase.TotalSupply) {
					t.Fatalf("Expected token owner to have %v balance. Got %v", testCase.TotalSupply.String(), account.Balance.String())
				}

				if !account.Balance.Equal(fungibleToken.TotalSupply) {
					t.Fatalf("After approving, issuer balance %v not equal to total supply %v", account.Balance.String(), fungibleToken.TotalSupply.String())
				}
				if account.Frozen {
					t.Fatal("Issuer started as frozen")
				}

				// Trying to transfer after approval
				resTransferToken := keeper.TransferFungibleToken(ctx, testCase.Symbol, testCase.Owner, testCase.NewOwner, testCase.TransferAmountOfToken)
				//fmt.Println(resTransferToken)

				if resTransferToken.Code == sdkTypes.CodeOK {
					fmt.Printf("Managed to transfer token after approval\n")

				}

			} else {
				if fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("token approved after unsucessful approval")
				}

				if fungibleToken.Metadata != "" {
					t.Fatal("Fungible token approval Metadata not empty after unsuccessful approval")
				}
			}

		}) // end for
	}

}

func TestMintFungibleToken(t *testing.T) {
	fmt.Printf("============\nStart Test : %s \n", "TestMintFungibleToken")
	ctx, keeper := PrepareTest(t)

	approver1, err := sdkTypes.AccAddressFromBech32("mxw1zzrutc6x9kc7ttafwawv9ve3jm89zxpeedl8ap")
	if err != nil {
		t.Fatal(err)
	}

	approver2, err := sdkTypes.AccAddressFromBech32("mxw1jdpaz5fppzt726w44hx4q9yerlfc3ldgpj83fe")
	if err != nil {
		t.Fatal(err)
	}

	approver1Acc := keeper.accountKeeper.NewAccountWithAddress(ctx, approver1)
	approver1Acc.SetAccountNumber(keeper.accountKeeper.GetNextAccountNumber(ctx))
	keeper.accountKeeper.SetAccount(ctx, approver1Acc)

	genesisState := &GenesisState{
		AuthorizedAddresses: []sdkTypes.AccAddress{approver1, approver2},
	}

	var applicationFee = struct {
		To    sdkTypes.AccAddress `json:"to"`
		Value string              `json:"value"`
	}{
		To:    delAddr1,
		Value: "100000",
	}

	var tokenFee1 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "zero",
		},
		{
			Action:  "mint",
			FeeName: "zero",
		},
		{
			Action:  "burn",
			FeeName: "zero",
		},
		{
			Action:  "transferOwnership",
			FeeName: "zero",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "zero",
		},
	}

	var tokenFee2 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "default",
		},
		{
			Action:  "mint",
			FeeName: "default",
		},
		{
			Action:  "burn",
			FeeName: "default",
		},
		{
			Action:  "transferOwnership",
			FeeName: "default",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "default",
		},
	}

	var testCases = []struct {
		Flags             types.Bitmask
		Name              string
		Symbol            string
		Decimals          int
		Owner             sdkTypes.AccAddress
		NewOwner          sdkTypes.AccAddress
		Metadata          string
		TotalSupply       sdkTypes.Uint
		MaxSupply         sdkTypes.Uint
		MintAmountOfToken sdkTypes.Uint
		TokenFees         []TokenFee
		ExpectedCode      sdkTypes.CodeType
	}{
		{
			Flags:             FixedSupplyNotBurnableFungibleTokenMask, // FungibleFlag
			Name:              "Happy path fixed supply - Not Burnable Fungible Token",
			Symbol:            "TST-1",
			Decimals:          18,
			Owner:             delAddr3,
			NewOwner:          delAddr2,
			TotalSupply:       sdkTypes.NewUint(0),
			MaxSupply:         sdkTypes.NewUint(0),
			Metadata:          "ipfs-hash-link",
			MintAmountOfToken: sdkTypes.NewUint(100),
			TokenFees:         tokenFee1,
			ExpectedCode:      sdkTypes.CodeOK,
		},
		{
			Flags:             FixedSupplyBurnableFungibleTokenMask, // FungibleFlag + BurnFlag
			Name:              "Happy path fixed supply - Burnable Fungible Token",
			Symbol:            "TST-2",
			Decimals:          18,
			Owner:             delAddr3,
			NewOwner:          delAddr2,
			TotalSupply:       sdkTypes.NewUint(0),
			MaxSupply:         sdkTypes.NewUint(0),
			Metadata:          "ipfs-hash-link",
			MintAmountOfToken: sdkTypes.NewUint(100),
			TokenFees:         tokenFee1,
			ExpectedCode:      sdkTypes.CodeOK,
		},
		{
			Flags:             DynamicFungibleTokenMask, // FungibleFlag + MintFlag + BurnFlag
			Name:              "Happy path dynamic supply",
			Symbol:            "tsttttt",
			Decimals:          18,
			Owner:             delAddr3,
			NewOwner:          delAddr2,
			TotalSupply:       sdkTypes.NewUint(0),
			MaxSupply:         sdkTypes.NewUint(0),
			Metadata:          "ipfs-hash-link",
			MintAmountOfToken: sdkTypes.NewUint(1),
			TokenFees:         tokenFee2,
			ExpectedCode:      sdkTypes.CodeOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			fmt.Println(testCase)
			InitGenesis(ctx, keeper, *genesisState)

			var fixedSupply bool
			fixedSupply = true
			if testCase.Flags.HasFlag(DynamicFungibleTokenMask) {
				fixedSupply = false
			}

			resCreateToken := keeper.CreateFungibleToken(ctx, testCase.Name, testCase.Symbol, testCase.Decimals, testCase.Owner,
				fixedSupply, testCase.MaxSupply, testCase.Metadata, applicationFee)
			//fmt.Printf("============\n")
			//fmt.Println(resCreateToken)
			if resCreateToken.Code != testCase.ExpectedCode {
				t.Fatalf("Expected code %d. Got code %d", testCase.ExpectedCode, resCreateToken.Code)
			}

			var burnable bool
			burnable = false
			if testCase.Flags.HasFlag(BurnFlag) {
				burnable = true
			}

			resApproval := keeper.ApproveToken(ctx, testCase.Symbol, testCase.TokenFees, burnable, approver1, "ipfs-hash-link")
			//fmt.Printf("============\n")
			//fmt.Println(resApproval)
			if resApproval.Code != testCase.ExpectedCode {
				t.Fatalf("Expected approval code %d. Got %d", testCase.ExpectedCode, resApproval.Code)
			}

			var fungibleToken = new(Token)
			keeper.mustGetTokenData(ctx, testCase.Symbol, fungibleToken)

			if testCase.ExpectedCode == sdkTypes.CodeOK {
				if !fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("Token not approved after successful approval")
				}

				if testCase.Metadata != fungibleToken.Metadata {
					t.Fatalf("Expected approval metada %s. Got %s", testCase.Metadata, fungibleToken.Metadata)
				}

				account := keeper.getFungibleAccount(ctx, testCase.Symbol, testCase.Owner)
				if account == nil {
					t.Fatal("After approval issuer account does not exist")
				}

				if !account.Balance.Equal(testCase.TotalSupply) {
					t.Fatalf("Expected token owner to have %v balance. Got %v", testCase.TotalSupply.String(), account.Balance.String())
				}

				if !account.Balance.Equal(fungibleToken.TotalSupply) {
					t.Fatalf("After approving, issuer balance %v not equal to total supply %v", account.Balance.String(), fungibleToken.TotalSupply.String())
				}
				if account.Frozen {
					t.Fatal("Issuer started as frozen")
				}

				// Trying to mint after approval
				resMintToken := keeper.MintFungibleToken(ctx, testCase.Symbol, testCase.Owner, testCase.NewOwner, testCase.MintAmountOfToken)
				//fmt.Println(resMintToken)
				//fmt.Printf("============\n\n")
				if resMintToken.Code == sdkTypes.CodeOK {
					fmt.Printf("\nManaged to mint token after approval\n")

				}

			} else {
				if fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("token approved after unsucessful approval")
				}

				if fungibleToken.Metadata != "" {
					t.Fatal("Fungible token approval Metadata not empty after unsuccessful approval")
				}
			}

		}) // end for
	}

}

func TestBurnFungibleToken(t *testing.T) {
	fmt.Printf("============\nStart Test : %s \n", "TestBurnFungibleToken")
	ctx, keeper := PrepareTest(t)

	approver1, err := sdkTypes.AccAddressFromBech32("mxw1zzrutc6x9kc7ttafwawv9ve3jm89zxpeedl8ap")
	if err != nil {
		t.Fatal(err)
	}

	approver2, err := sdkTypes.AccAddressFromBech32("mxw1jdpaz5fppzt726w44hx4q9yerlfc3ldgpj83fe")
	if err != nil {
		t.Fatal(err)
	}

	approver1Acc := keeper.accountKeeper.NewAccountWithAddress(ctx, approver1)
	approver1Acc.SetAccountNumber(keeper.accountKeeper.GetNextAccountNumber(ctx))
	keeper.accountKeeper.SetAccount(ctx, approver1Acc)

	genesisState := &GenesisState{
		AuthorizedAddresses: []sdkTypes.AccAddress{approver1, approver2},
	}

	var applicationFee = struct {
		To    sdkTypes.AccAddress `json:"to"`
		Value string              `json:"value"`
	}{
		To:    delAddr1,
		Value: "100000",
	}

	var tokenFee1 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "zero",
		},
		{
			Action:  "mint",
			FeeName: "zero",
		},
		{
			Action:  "burn",
			FeeName: "zero",
		},
		{
			Action:  "transferOwnership",
			FeeName: "zero",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "zero",
		},
	}

	var tokenFee2 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "default",
		},
		{
			Action:  "mint",
			FeeName: "default",
		},
		{
			Action:  "burn",
			FeeName: "default",
		},
		{
			Action:  "transferOwnership",
			FeeName: "default",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "default",
		},
	}

	var testCases = []struct {
		Flags             types.Bitmask
		Name              string
		Symbol            string
		Decimals          int
		Owner             sdkTypes.AccAddress
		NewOwner          sdkTypes.AccAddress
		Metadata          string
		TotalSupply       sdkTypes.Uint
		MaxSupply         sdkTypes.Uint
		BurnAmountOfToken sdkTypes.Uint
		TokenFees         []TokenFee
		ExpectedCode      sdkTypes.CodeType
	}{
		{
			Flags:             FixedSupplyNotBurnableFungibleTokenMask, // FungibleFlag
			Name:              "Happy path fixed supply - Not Burnable Fungible Token",
			Symbol:            "TST-1",
			Decimals:          18,
			Owner:             delAddr3,
			TotalSupply:       sdkTypes.NewUint(0),
			MaxSupply:         sdkTypes.NewUint(0),
			Metadata:          "ipfs-hash-link",
			BurnAmountOfToken: sdkTypes.NewUint(10000),
			TokenFees:         tokenFee1,
			ExpectedCode:      sdkTypes.CodeOK,
		},
		{
			Flags:             FixedSupplyBurnableFungibleTokenMask, // FungibleFlag + BurnFlag
			Name:              "Happy path fixed supply - Burnable Fungible Token",
			Symbol:            "TST-2",
			Decimals:          18,
			Owner:             delAddr1,
			TotalSupply:       sdkTypes.NewUint(0),
			MaxSupply:         sdkTypes.NewUint(0),
			Metadata:          "ipfs-hash-link",
			BurnAmountOfToken: sdkTypes.NewUint(10000),
			TokenFees:         tokenFee1,
			ExpectedCode:      sdkTypes.CodeOK,
		},
		{
			Flags:             DynamicFungibleTokenMask, // FungibleFlag + MintFlag + BurnFlag
			Name:              "Happy path dynamic supply",
			Symbol:            "tsttttt",
			Decimals:          18,
			Owner:             delAddr3,
			TotalSupply:       sdkTypes.NewUint(0),
			MaxSupply:         sdkTypes.NewUint(0),
			Metadata:          "ipfs-hash-link",
			BurnAmountOfToken: sdkTypes.NewUint(10),
			TokenFees:         tokenFee2,
			ExpectedCode:      sdkTypes.CodeOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			fmt.Println(testCase)
			InitGenesis(ctx, keeper, *genesisState)

			var fixedSupply bool
			fixedSupply = true
			if testCase.Flags.HasFlag(DynamicFungibleTokenMask) {
				fixedSupply = false
			}

			resCreateToken := keeper.CreateFungibleToken(ctx, testCase.Name, testCase.Symbol, testCase.Decimals, testCase.Owner,
				fixedSupply, testCase.MaxSupply, testCase.Metadata, applicationFee)
			//fmt.Println(resCreateToken)
			if resCreateToken.Code != testCase.ExpectedCode {
				t.Fatalf("Expected code %d. Got code %d", testCase.ExpectedCode, resCreateToken.Code)
			}

			var burnable bool
			burnable = false
			if testCase.Flags.HasFlag(BurnFlag) {
				burnable = true
			}

			resApproval := keeper.ApproveToken(ctx, testCase.Symbol, testCase.TokenFees, burnable, approver1, "ipfs-hash-link")
			//fmt.Println(resApproval)
			if resApproval.Code != testCase.ExpectedCode {
				t.Fatalf("Expected approval code %d. Got %d", testCase.ExpectedCode, resApproval.Code)
			}

			var fungibleToken = new(Token)
			keeper.mustGetTokenData(ctx, testCase.Symbol, fungibleToken)

			if testCase.ExpectedCode == sdkTypes.CodeOK {
				if !fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("Token not approved after successful approval")
				}

				if testCase.Metadata != fungibleToken.Metadata {
					t.Fatalf("Expected approval metada %s. Got %s", testCase.Metadata, fungibleToken.Metadata)
				}

				account := keeper.getFungibleAccount(ctx, testCase.Symbol, testCase.Owner)
				if account == nil {
					t.Fatal("After approval issuer account does not exist")
				}

				if !account.Balance.Equal(testCase.TotalSupply) {
					t.Fatalf("Expected token owner to have %v balance. Got %v", testCase.TotalSupply.String(), account.Balance.String())
				}

				if !account.Balance.Equal(fungibleToken.TotalSupply) {
					t.Fatalf("After approving, issuer balance %v not equal to total supply %v", account.Balance.String(), fungibleToken.TotalSupply.String())
				}
				if account.Frozen {
					t.Fatal("Issuer started as frozen")
				}

				// Trying to burn after approval
				resBurnToken := keeper.BurnFungibleToken(ctx, testCase.Symbol, testCase.Owner, account.Balance)
				//fmt.Println(resBurnToken)

				if resBurnToken.Code == sdkTypes.CodeOK {
					fmt.Printf("Managed to burn token after approval\n")

				}

			} else {
				if fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("token approved after unsucessful approval")
				}

				if fungibleToken.Metadata != "" {
					t.Fatal("Fungible token approval Metadata not empty after unsuccessful approval")
				}
			}

		}) // end for
	}

}

func TestTransferTokenOwnership(t *testing.T) {

	fmt.Printf("============\nStart Test : %s \n", "TestTransferTokenOwnership")
	ctx, keeper := PrepareTest(t)

	approver1, err := sdkTypes.AccAddressFromBech32("mxw1zzrutc6x9kc7ttafwawv9ve3jm89zxpeedl8ap")
	if err != nil {
		t.Fatal(err)
	}

	approver2, err := sdkTypes.AccAddressFromBech32("mxw1jdpaz5fppzt726w44hx4q9yerlfc3ldgpj83fe")
	if err != nil {
		t.Fatal(err)
	}

	approver1Acc := keeper.accountKeeper.NewAccountWithAddress(ctx, approver1)
	approver1Acc.SetAccountNumber(keeper.accountKeeper.GetNextAccountNumber(ctx))
	keeper.accountKeeper.SetAccount(ctx, approver1Acc)

	genesisState := &GenesisState{
		AuthorizedAddresses: []sdkTypes.AccAddress{approver1, approver2},
	}

	var applicationFee = struct {
		To    sdkTypes.AccAddress `json:"to"`
		Value string              `json:"value"`
	}{
		To:    delAddr1,
		Value: "100000",
	}

	var tokenFee1 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "zero",
		},
		{
			Action:  "mint",
			FeeName: "zero",
		},
		{
			Action:  "burn",
			FeeName: "zero",
		},
		{
			Action:  "transferOwnership",
			FeeName: "zero",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "zero",
		},
	}

	var tokenFee2 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "default",
		},
		{
			Action:  "mint",
			FeeName: "default",
		},
		{
			Action:  "burn",
			FeeName: "default",
		},
		{
			Action:  "transferOwnership",
			FeeName: "default",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "default",
		},
	}

	var testCases = []struct {
		Flags        types.Bitmask
		Name         string
		Symbol       string
		Decimals     int
		Owner        sdkTypes.AccAddress
		NewOwner     sdkTypes.AccAddress
		Metadata     string
		TotalSupply  sdkTypes.Uint
		MaxSupply    sdkTypes.Uint
		TransferFee  sdkTypes.Uint
		ExpectedCode sdkTypes.CodeType
		TokenFees    []TokenFee
	}{
		{
			Flags:        FixedSupplyNotBurnableFungibleTokenMask, // FungibleFlag
			Name:         "Happy path fixed supply - Not Burnable Fungible Token",
			Symbol:       "TST-1",
			Decimals:     18,
			Owner:        delAddr3,
			NewOwner:     delAddr2,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			TransferFee:  sdkTypes.NewUint(100),
			ExpectedCode: sdkTypes.CodeOK,
			TokenFees:    tokenFee1,
		},
		{
			Flags:        FixedSupplyBurnableFungibleTokenMask, // FungibleFlag + BurnFlag
			Name:         "Happy path fixed supply - Burnable Fungible Token",
			Symbol:       "TST-2",
			Decimals:     18,
			Owner:        delAddr1,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			TransferFee:  sdkTypes.NewUint(100),
			ExpectedCode: sdkTypes.CodeOK,
			TokenFees:    tokenFee1,
		},
		{
			Flags:        DynamicFungibleTokenMask, // FungibleFlag + MintFlag + BurnFlag
			Name:         "Happy path dynamic supply",
			Symbol:       "tsttttt",
			Decimals:     18,
			Owner:        delAddr3,
			NewOwner:     delAddr2,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			TransferFee:  sdkTypes.NewUint(10),
			ExpectedCode: sdkTypes.CodeOK,
			TokenFees:    tokenFee2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			fmt.Println(testCase)
			InitGenesis(ctx, keeper, *genesisState)

			var fixedSupply bool
			fixedSupply = true
			if testCase.Flags.HasFlag(DynamicFungibleTokenMask) {
				fixedSupply = false
			}

			resCreateToken := keeper.CreateFungibleToken(ctx, testCase.Name, testCase.Symbol, testCase.Decimals, testCase.Owner,
				fixedSupply, testCase.MaxSupply, testCase.Metadata, applicationFee)
			//fmt.Println(resCreateToken)
			if resCreateToken.Code != testCase.ExpectedCode {
				t.Fatalf("Expected code %d. Got code %d", testCase.ExpectedCode, resCreateToken.Code)
			}

			var burnable bool
			burnable = false
			if testCase.Flags.HasFlag(BurnFlag) {
				burnable = true
			}

			resApproval := keeper.ApproveToken(ctx, testCase.Symbol, testCase.TokenFees, burnable, approver1, "ipfs-hash-link")
			//fmt.Println(resApproval)
			if resApproval.Code != testCase.ExpectedCode {
				t.Fatalf("Expected approval code %d. Got %d", testCase.ExpectedCode, resApproval.Code)
			}

			var fungibleToken = new(Token)
			keeper.mustGetTokenData(ctx, testCase.Symbol, fungibleToken)

			if testCase.ExpectedCode == sdkTypes.CodeOK {
				if !fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("Token not approved after successful approval")
				}

				if testCase.Metadata != fungibleToken.Metadata {
					t.Fatalf("Expected approval metada %s. Got %s", testCase.Metadata, fungibleToken.Metadata)
				}

				account := keeper.getFungibleAccount(ctx, testCase.Symbol, testCase.Owner)
				if account == nil {
					t.Fatal("After approval issuer account does not exist")
				}

				if !account.Balance.Equal(testCase.TotalSupply) {
					t.Fatalf("Expected token owner to have %v balance. Got %v", testCase.TotalSupply.String(), account.Balance.String())
				}

				if !account.Balance.Equal(fungibleToken.TotalSupply) {
					t.Fatalf("After approving, issuer balance %v not equal to total supply %v", account.Balance.String(), fungibleToken.TotalSupply.String())
				}
				if account.Frozen {
					t.Fatal("Issuer started as frozen")
				}

				// Trying to transfer ownership after approval
				// from : Original Owner
				// to : New Owner
				resTransferTokenOwnership := keeper.transferFungibleTokenOwnership(ctx, testCase.Owner, testCase.NewOwner, fungibleToken, testCase.Metadata)
				//fmt.Println(resTransferTokenOwnership)

				if resTransferTokenOwnership.Code == sdkTypes.CodeOK {
					fmt.Printf("Managed to transfer-ownership base on approval token\n")

				}

			} else {
				if fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("token approved after unsucessful approval")
				}

				if fungibleToken.Metadata != "" {
					t.Fatal("Fungible token approval Metadata not empty after unsuccessful approval")
				}
			}

		}) // end for
	}

}

func TestAcceptTokenOwnership(t *testing.T) {

	fmt.Printf("============\nStart Test : %s \n", "TestAcceptTokenOwnership")
	ctx, keeper := PrepareTest(t)

	approver1, err := sdkTypes.AccAddressFromBech32("mxw1zzrutc6x9kc7ttafwawv9ve3jm89zxpeedl8ap")
	if err != nil {
		t.Fatal(err)
	}

	approver2, err := sdkTypes.AccAddressFromBech32("mxw1jdpaz5fppzt726w44hx4q9yerlfc3ldgpj83fe")
	if err != nil {
		t.Fatal(err)
	}

	approver1Acc := keeper.accountKeeper.NewAccountWithAddress(ctx, approver1)
	approver1Acc.SetAccountNumber(keeper.accountKeeper.GetNextAccountNumber(ctx))
	keeper.accountKeeper.SetAccount(ctx, approver1Acc)

	genesisState := &GenesisState{
		AuthorizedAddresses: []sdkTypes.AccAddress{approver1, approver2},
	}

	var applicationFee = struct {
		To    sdkTypes.AccAddress `json:"to"`
		Value string              `json:"value"`
	}{
		To:    delAddr1,
		Value: "100000",
	}

	var tokenFee1 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "zero",
		},
		{
			Action:  "mint",
			FeeName: "zero",
		},
		{
			Action:  "burn",
			FeeName: "zero",
		},
		{
			Action:  "transferOwnership",
			FeeName: "zero",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "zero",
		},
	}

	var tokenFee2 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "default",
		},
		{
			Action:  "mint",
			FeeName: "default",
		},
		{
			Action:  "burn",
			FeeName: "default",
		},
		{
			Action:  "transferOwnership",
			FeeName: "default",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "default",
		},
	}

	var testCases = []struct {
		Flags        types.Bitmask
		Name         string
		Symbol       string
		Decimals     int
		Owner        sdkTypes.AccAddress
		NewOwner     sdkTypes.AccAddress
		Metadata     string
		TotalSupply  sdkTypes.Uint
		MaxSupply    sdkTypes.Uint
		TransferFee  sdkTypes.Uint
		ExpectedCode sdkTypes.CodeType
		TokenFees    []TokenFee
	}{
		{
			Flags:        FixedSupplyNotBurnableFungibleTokenMask, // FungibleFlag
			Name:         "Happy path fixed supply - Not Burnable Fungible Token",
			Symbol:       "TST-1",
			Decimals:     18,
			Owner:        delAddr3,
			NewOwner:     delAddr2,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			TransferFee:  sdkTypes.NewUint(100),
			ExpectedCode: sdkTypes.CodeOK,
			TokenFees:    tokenFee1,
		},
		{
			Flags:        FixedSupplyBurnableFungibleTokenMask, // FungibleFlag + BurnFlag
			Name:         "Happy path fixed supply - Burnable Fungible Token",
			Symbol:       "TST-2",
			Decimals:     18,
			Owner:        delAddr1,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			TransferFee:  sdkTypes.NewUint(100),
			ExpectedCode: sdkTypes.CodeOK,
			TokenFees:    tokenFee1,
		},
		{
			Flags:        DynamicFungibleTokenMask, // FungibleFlag + MintFlag + BurnFlag
			Name:         "Happy path dynamic supply",
			Symbol:       "tsttttt",
			Decimals:     18,
			Owner:        delAddr3,
			NewOwner:     delAddr2,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			TransferFee:  sdkTypes.NewUint(10),
			ExpectedCode: sdkTypes.CodeOK,
			TokenFees:    tokenFee2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			fmt.Println(testCase)
			InitGenesis(ctx, keeper, *genesisState)

			var fixedSupply bool
			fixedSupply = true
			if testCase.Flags.HasFlag(DynamicFungibleTokenMask) {
				fixedSupply = false
			}

			resCreateToken := keeper.CreateFungibleToken(ctx, testCase.Name, testCase.Symbol, testCase.Decimals, testCase.Owner,
				fixedSupply, testCase.MaxSupply, testCase.Metadata, applicationFee)
			//fmt.Println(resCreateToken)
			if resCreateToken.Code != testCase.ExpectedCode {
				t.Fatalf("Expected code %d. Got code %d", testCase.ExpectedCode, resCreateToken.Code)
			}

			var burnable bool
			burnable = false
			if testCase.Flags.HasFlag(BurnFlag) {
				burnable = true
			}

			resApproval := keeper.ApproveToken(ctx, testCase.Symbol, testCase.TokenFees, burnable, approver1, "ipfs-hash-link")
			//fmt.Println(resApproval)
			if resApproval.Code != testCase.ExpectedCode {
				t.Fatalf("Expected approval code %d. Got %d", testCase.ExpectedCode, resApproval.Code)
			}

			var fungibleToken = new(Token)
			keeper.mustGetTokenData(ctx, testCase.Symbol, fungibleToken)

			if testCase.ExpectedCode == sdkTypes.CodeOK {
				if !fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("Token not approved after successful approval")
				}

				if testCase.Metadata != fungibleToken.Metadata {
					t.Fatalf("Expected approval metada %s. Got %s", testCase.Metadata, fungibleToken.Metadata)
				}

				account := keeper.getFungibleAccount(ctx, testCase.Symbol, testCase.Owner)
				if account == nil {
					t.Fatal("After approval issuer account does not exist")
				}

				if !account.Balance.Equal(testCase.TotalSupply) {
					t.Fatalf("Expected token owner to have %v balance. Got %v", testCase.TotalSupply.String(), account.Balance.String())
				}

				if !account.Balance.Equal(fungibleToken.TotalSupply) {
					t.Fatalf("After approving, issuer balance %v not equal to total supply %v", account.Balance.String(), fungibleToken.TotalSupply.String())
				}
				if account.Frozen {
					t.Fatal("Issuer started as frozen")
				}

				// Trying to transfer ownership after approval
				// from : Original Owner
				// to : New Owner
				resTransferTokenOwnership := keeper.transferFungibleTokenOwnership(ctx, testCase.Owner, testCase.NewOwner, fungibleToken, testCase.Metadata)
				//fmt.Println(resTransferTokenOwnership)

				if resTransferTokenOwnership.Code == sdkTypes.CodeOK {
					fmt.Printf("Managed to transfer-ownership base on approval token\n")

					// Trying to accept ownership after transfer-token-ownership
					// from : New Owner who accepted the ownership from original party
					// tokenOwner : New Owner
					resAcceptTokenOwnership := keeper.acceptFungibleTokenOwnership(ctx, testCase.NewOwner, fungibleToken, testCase.Metadata)
					//fmt.Println(resAcceptTokenOwnership)

					if resAcceptTokenOwnership.Code == sdkTypes.CodeOK {
						fmt.Printf("Managed to accept ownership base on transfer-ownership\n")
					} else {
						fmt.Printf("AcceptTokenOwnership Error : %x \n", resAcceptTokenOwnership.Code)
					}

				}

			} else {
				if fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("token approved after unsucessful approval")
				}

				if fungibleToken.Metadata != "" {
					t.Fatal("Fungible token approval Metadata not empty after unsuccessful approval")
				}
			}

		}) // end for
	}

}

func TestFreezeFungibleTokenAccount(t *testing.T) {
	ctx, keeper := PrepareTest(t)

	fmt.Printf("============\nStart Test : %s \n", "TestFreezeFungibleTokenAccount")
	approver1, err := sdkTypes.AccAddressFromBech32("mxw1zzrutc6x9kc7ttafwawv9ve3jm89zxpeedl8ap")
	if err != nil {
		t.Fatal(err)
	}

	approver2, err := sdkTypes.AccAddressFromBech32("mxw1jdpaz5fppzt726w44hx4q9yerlfc3ldgpj83fe")
	if err != nil {
		t.Fatal(err)
	}

	approver1Acc := keeper.accountKeeper.NewAccountWithAddress(ctx, approver1)
	approver1Acc.SetAccountNumber(keeper.accountKeeper.GetNextAccountNumber(ctx))
	keeper.accountKeeper.SetAccount(ctx, approver1Acc)

	genesisState := &GenesisState{
		AuthorizedAddresses: []sdkTypes.AccAddress{approver1, approver2},
	}

	var applicationFee = struct {
		To    sdkTypes.AccAddress `json:"to"`
		Value string              `json:"value"`
	}{
		To:    delAddr1,
		Value: "100000",
	}

	var tokenFee1 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "zero",
		},
		{
			Action:  "mint",
			FeeName: "zero",
		},
		{
			Action:  "burn",
			FeeName: "zero",
		},
		{
			Action:  "transferOwnership",
			FeeName: "zero",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "zero",
		},
	}

	var tokenFee2 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "default",
		},
		{
			Action:  "mint",
			FeeName: "default",
		},
		{
			Action:  "burn",
			FeeName: "default",
		},
		{
			Action:  "transferOwnership",
			FeeName: "default",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "default",
		},
	}

	var testCases = []struct {
		Flags        types.Bitmask
		Name         string
		Symbol       string
		Decimals     int
		Owner        sdkTypes.AccAddress
		NewOwner     sdkTypes.AccAddress
		Metadata     string
		TotalSupply  sdkTypes.Uint
		MaxSupply    sdkTypes.Uint
		TransferFee  sdkTypes.Uint
		ExpectedCode sdkTypes.CodeType
		TokenFees    []TokenFee
	}{
		{
			Flags:        FixedSupplyNotBurnableFungibleTokenMask, // FungibleFlag
			Name:         "Happy path fixed supply - Not Burnable Fungible Token",
			Symbol:       "TST-1",
			Decimals:     18,
			Owner:        delAddr3,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
			TokenFees:    tokenFee1,
		},
		{
			Flags:        FixedSupplyBurnableFungibleTokenMask, // FungibleFlag + BurnFlag
			Name:         "Happy path fixed supply - Burnable Fungible Token",
			Symbol:       "TST-2",
			Decimals:     18,
			Owner:        delAddr1,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
			TokenFees:    tokenFee1,
		},
		{
			Flags:        DynamicFungibleTokenMask, // FungibleFlag + MintFlag + BurnFlag
			Name:         "Happy path dynamic supply",
			Symbol:       "tsttttt",
			Decimals:     18,
			Owner:        delAddr3,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
			TokenFees:    tokenFee2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			fmt.Println(testCase)
			InitGenesis(ctx, keeper, *genesisState)

			var fixedSupply bool
			fixedSupply = true
			if testCase.Flags.HasFlag(DynamicFungibleTokenMask) {
				fixedSupply = false
			}

			res := keeper.CreateFungibleToken(ctx, testCase.Name, testCase.Symbol, testCase.Decimals, testCase.Owner,
				fixedSupply, testCase.MaxSupply, testCase.Metadata, applicationFee)
			if res.Code != testCase.ExpectedCode {
				t.Fatalf("Expected CreateFungibleToken code %d. Got code %d", testCase.ExpectedCode, res.Code)
			}

			var burnable bool
			burnable = false
			if testCase.Flags.HasFlag(BurnFlag) {
				burnable = true
			}

			resApproval := keeper.ApproveToken(ctx, testCase.Symbol, testCase.TokenFees, burnable, approver1, testCase.Metadata)
			if resApproval.Code != testCase.ExpectedCode {
				t.Fatalf("Expected ApproveToken code %d. Got %d", testCase.ExpectedCode, resApproval.Code)
			}

			var fungibleToken = new(Token)
			keeper.mustGetTokenData(ctx, testCase.Symbol, fungibleToken)

			if testCase.ExpectedCode == sdkTypes.CodeOK {
				if !fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("Token not approved after successful approval")
				}

				if testCase.Metadata != fungibleToken.Metadata {
					t.Fatalf("Expected approval metada %s. Got %s", testCase.Metadata, fungibleToken.Metadata)
				}

				account := keeper.getFungibleAccount(ctx, testCase.Symbol, testCase.Owner)
				if account == nil {
					t.Fatal("After approval issuer account does not exist")
				}

				if !account.Balance.Equal(testCase.TotalSupply) {
					t.Fatalf("Expected token owner to have %v balance. Got %v", testCase.TotalSupply.String(), account.Balance.String())
				}

				if !account.Balance.Equal(fungibleToken.TotalSupply) {
					t.Fatalf("After approving, issuer balance %v not equal to total supply %v", account.Balance.String(), fungibleToken.TotalSupply.String())
				}
				if account.Frozen {
					t.Fatal("Issuer started as frozen")
				}

				// Trying to Freeze TokenAccount after approval
				resFreezeFungibleTokenAccount := keeper.FreezeFungibleTokenAccount(ctx, testCase.Symbol, approver1, testCase.Owner, testCase.Metadata)
				if resFreezeFungibleTokenAccount.Code != sdkTypes.CodeOK {
					t.Fatalf("Expected FreezeFungibleTokenAccount code %d. Got %d", sdkTypes.CodeOK, resFreezeFungibleTokenAccount.Code)
				} else {
					fmt.Printf("OK Result")
				}

			} else {
				if fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("token approved after unsucessful approval")
				}

				if fungibleToken.Metadata != "" {
					t.Fatal("Fungible token approval Metadata not empty after unsuccessful approval")
				}
			}

		})
	}

}

func TestUnfreezeFungibleTokenAccount(t *testing.T) {
	ctx, keeper := PrepareTest(t)

	fmt.Printf("============\nStart Test : %s \n", "TestUnfreezeFungibleTokenAccount")
	approver1, err := sdkTypes.AccAddressFromBech32("mxw1zzrutc6x9kc7ttafwawv9ve3jm89zxpeedl8ap")
	if err != nil {
		t.Fatal(err)
	}

	approver2, err := sdkTypes.AccAddressFromBech32("mxw1jdpaz5fppzt726w44hx4q9yerlfc3ldgpj83fe")
	if err != nil {
		t.Fatal(err)
	}

	approver1Acc := keeper.accountKeeper.NewAccountWithAddress(ctx, approver1)
	approver1Acc.SetAccountNumber(keeper.accountKeeper.GetNextAccountNumber(ctx))
	keeper.accountKeeper.SetAccount(ctx, approver1Acc)

	genesisState := &GenesisState{
		AuthorizedAddresses: []sdkTypes.AccAddress{approver1, approver2, delAddr3}, // include : token-owner
	}

	var applicationFee = struct {
		To    sdkTypes.AccAddress `json:"to"`
		Value string              `json:"value"`
	}{
		To:    delAddr1,
		Value: "100000",
	}

	var tokenFee1 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "zero",
		},
		{
			Action:  "mint",
			FeeName: "zero",
		},
		{
			Action:  "burn",
			FeeName: "zero",
		},
		{
			Action:  "transferOwnership",
			FeeName: "zero",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "zero",
		},
	}

	var tokenFee2 = []TokenFee{
		{
			Action:  "transfer",
			FeeName: "default",
		},
		{
			Action:  "mint",
			FeeName: "default",
		},
		{
			Action:  "burn",
			FeeName: "default",
		},
		{
			Action:  "transferOwnership",
			FeeName: "default",
		},
		{
			Action:  "acceptOwnership",
			FeeName: "default",
		},
	}

	var testCases = []struct {
		Flags        types.Bitmask
		Name         string
		Symbol       string
		Decimals     int
		Owner        sdkTypes.AccAddress
		NewOwner     sdkTypes.AccAddress
		Metadata     string
		TotalSupply  sdkTypes.Uint
		MaxSupply    sdkTypes.Uint
		TransferFee  sdkTypes.Uint
		ExpectedCode sdkTypes.CodeType
		TokenFees    []TokenFee
	}{
		{
			Flags:        FixedSupplyNotBurnableFungibleTokenMask, // FungibleFlag
			Name:         "Happy path fixed supply - Not Burnable Fungible Token",
			Symbol:       "TST-1",
			Decimals:     18,
			Owner:        delAddr3,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
			TokenFees:    tokenFee1,
		},
		{
			Flags:        FixedSupplyBurnableFungibleTokenMask, // FungibleFlag + BurnFlag
			Name:         "Happy path fixed supply - Burnable Fungible Token",
			Symbol:       "TST-2",
			Decimals:     18,
			Owner:        delAddr1,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
			TokenFees:    tokenFee1,
		},
		{
			Flags:        DynamicFungibleTokenMask, // FungibleFlag + MintFlag + BurnFlag
			Name:         "Happy path dynamic supply",
			Symbol:       "tsttttt",
			Decimals:     18,
			Owner:        delAddr3,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
			TokenFees:    tokenFee2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			fmt.Println(testCase)
			InitGenesis(ctx, keeper, *genesisState)

			var fixedSupply bool
			fixedSupply = true
			if testCase.Flags.HasFlag(DynamicFungibleTokenMask) {
				fixedSupply = false
			}

			res := keeper.CreateFungibleToken(ctx, testCase.Name, testCase.Symbol, testCase.Decimals, testCase.Owner,
				fixedSupply, testCase.MaxSupply, testCase.Metadata, applicationFee)
			if res.Code != testCase.ExpectedCode {
				t.Fatalf("Expected CreateFungibleToken code %d. Got code %d", testCase.ExpectedCode, res.Code)
			}

			var burnable bool
			burnable = false
			if testCase.Flags.HasFlag(BurnFlag) {
				burnable = true
			}

			resApproval := keeper.ApproveToken(ctx, testCase.Symbol, testCase.TokenFees, burnable, approver1, testCase.Metadata)
			if resApproval.Code != testCase.ExpectedCode {
				t.Fatalf("Expected ApproveToken code %d. Got %d", testCase.ExpectedCode, resApproval.Code)
			}

			var fungibleToken = new(Token)
			keeper.mustGetTokenData(ctx, testCase.Symbol, fungibleToken)

			if testCase.ExpectedCode == sdkTypes.CodeOK {
				if !fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("Token not approved after successful approval")
				}

				if testCase.Metadata != fungibleToken.Metadata {
					t.Fatalf("Expected approval metada %s. Got %s", testCase.Metadata, fungibleToken.Metadata)
				}

				account := keeper.getFungibleAccount(ctx, testCase.Symbol, testCase.Owner)
				if account == nil {
					t.Fatal("After approval issuer account does not exist")
				}

				if !account.Balance.Equal(testCase.TotalSupply) {
					t.Fatalf("Expected token owner to have %v balance. Got %v", testCase.TotalSupply.String(), account.Balance.String())
				}

				if !account.Balance.Equal(fungibleToken.TotalSupply) {
					t.Fatalf("After approving, issuer balance %v not equal to total supply %v", account.Balance.String(), fungibleToken.TotalSupply.String())
				}
				if account.Frozen {
					t.Fatal("Issuer started as frozen")
				}

				// Trying to Freeze TokenAccount after approval
				resFreezeFungibleTokenAccount := keeper.FreezeFungibleTokenAccount(ctx, testCase.Symbol, approver1, testCase.Owner, testCase.Metadata)
				if resFreezeFungibleTokenAccount.Code == sdkTypes.CodeOK {
					fmt.Printf("Managed to freeze token account\n")

					// Trying to unfreeze TokenAccount after freezed
					resUnfreezeFungibleTokenAccount := keeper.UnfreezeFungibleTokenAccount(ctx, testCase.Symbol, approver1, testCase.Owner, testCase.Metadata)
					fmt.Println(resUnfreezeFungibleTokenAccount)

					if resUnfreezeFungibleTokenAccount.Code == sdkTypes.CodeOK {
						fmt.Printf("Managed to unfreeze fungible token account\n")
					} else {
						t.Fatalf("Expected UnfreezeFungibleTokenAccount code %d. Got %d", testCase.ExpectedCode, resUnfreezeFungibleTokenAccount.Code)
					}

				} else {

					t.Fatalf("Expected FreezeFungibleTokenAccount code %d. Got %d", testCase.ExpectedCode, resFreezeFungibleTokenAccount.Code)
				}

			} else {
				if fungibleToken.Flags.HasFlag(ApprovedFlag) {
					t.Fatal("token approved after unsucessful approval")
				}

				if fungibleToken.Metadata != "" {
					t.Fatal("Fungible token approval Metadata not empty after unsuccessful approval")
				}
			}

		}) // end for
	}

}

func TestListTokenData(t *testing.T) {

	fmt.Printf("============\nStart Test : %s \n", "TestListTokenData")
	ctx, keeper := PrepareTest(t)

	approver1, err := sdkTypes.AccAddressFromBech32("mxw1zzrutc6x9kc7ttafwawv9ve3jm89zxpeedl8ap")
	if err != nil {
		t.Fatal(err)
	}

	approver2, err := sdkTypes.AccAddressFromBech32("mxw1jdpaz5fppzt726w44hx4q9yerlfc3ldgpj83fe")
	if err != nil {
		t.Fatal(err)
	}

	approver1Acc := keeper.accountKeeper.NewAccountWithAddress(ctx, approver1)
	approver1Acc.SetAccountNumber(keeper.accountKeeper.GetNextAccountNumber(ctx))
	keeper.accountKeeper.SetAccount(ctx, approver1Acc)

	genesisState := &GenesisState{
		AuthorizedAddresses: []sdkTypes.AccAddress{approver1, approver2},
	}

	var applicationFee = struct {
		To    sdkTypes.AccAddress `json:"to"`
		Value string              `json:"value"`
	}{
		To:    delAddr1,
		Value: "0",
	}

	var testCases = []struct {
		Flags        types.Bitmask
		Name         string
		Symbol       string
		Decimals     int
		Owner        sdkTypes.AccAddress
		NewOwner     sdkTypes.AccAddress
		Metadata     string
		TotalSupply  sdkTypes.Uint
		MaxSupply    sdkTypes.Uint
		TransferFee  sdkTypes.Uint
		ExpectedCode sdkTypes.CodeType
	}{
		{
			Flags:        FixedSupplyNotBurnableFungibleTokenMask, // FungibleFlag
			Name:         "Happy path fixed supply - Not Burnable Fungible Token",
			Symbol:       "TST-1",
			Decimals:     18,
			Owner:        delAddr1,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
		},
		{
			Flags:        FixedSupplyBurnableFungibleTokenMask, // FungibleFlag + BurnFlag
			Name:         "Happy path fixed supply - Burnable Fungible Token",
			Symbol:       "TST-2",
			Decimals:     18,
			Owner:        delAddr1,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
		},
		{
			Flags:        DynamicFungibleTokenMask, // FungibleFlag + MintFlag + BurnFlag
			Name:         "Happy path dynamic supply",
			Symbol:       "tsttttt",
			Decimals:     18,
			Owner:        delAddr2,
			TotalSupply:  sdkTypes.NewUint(0),
			MaxSupply:    sdkTypes.NewUint(0),
			Metadata:     "ipfs-hash-link",
			ExpectedCode: sdkTypes.CodeOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			fmt.Println(testCase)
			InitGenesis(ctx, keeper, *genesisState)

			var fixedSupply bool
			fixedSupply = true
			if testCase.Flags.HasFlag(DynamicFungibleTokenMask) {
				fixedSupply = false
			}

			res := keeper.CreateFungibleToken(ctx, testCase.Name, testCase.Symbol, testCase.Decimals, testCase.Owner,
				fixedSupply, testCase.TotalSupply, testCase.Metadata, applicationFee)
			//fmt.Println(res)

			if res.Code != testCase.ExpectedCode {
				t.Fatalf("Expected code %d. Got code %d", testCase.ExpectedCode, res.Code)
			}

			if !keeper.TokenExists(ctx, testCase.Symbol) {
				t.Fatal("After creating, token does not exist")
			}

		})
	}

	resFungibleTokens := keeper.ListTokens(ctx)
	fmt.Printf("Data of Fungible Token : %v\n", resFungibleTokens)
}

func TestTokenAuthorisedAddresses(t *testing.T) {

	fmt.Printf("============\nStart Test : %s \n", "TestTokenAuthorisedAddresses")

	ctx, keeper := PrepareTest(t)

	var testAddrs = []string{
		"mxw1yw6mg7fty4mzcwupvzek53x5egm7tp2ldwaxq3",
		"mxw1yyz3h9calxmvjp4x05nnn70a8ex7fee3th7r7k",
		"mxw12vptxnfmt88rjyade6l3uzjeth64w005mdwhjj",
	}

	var a []sdkTypes.AccAddress
	/// ------------------------------------case1: IsAuthorised ?
	for i, testAddr := range testAddrs {
		i = i + 1

		rsAddr, err := sdkTypes.AccAddressFromBech32(testAddr)
		assert.NoError(t, err)
		fmt.Printf("============\nTest Address : %d with value : %s \n", i, testAddr)

		rsUnauthorisedAddr := keeper.IsAuthorised(ctx, rsAddr)
		assert.False(t, rsUnauthorisedAddr)
		fmt.Printf("Is not an Authorised Address\n")

		a = append(a, rsAddr)

	}

	/// ------------------------------------case2: ADD -> GET
	keeper.SetAuthorisedAddresses(ctx, a)
	rsAddGet := keeper.GetAuthorisedAddresses(ctx)

	if assert.NotNil(t, rsAddGet) {
		for j, testAddr := range testAddrs {
			assert.Contains(t, rsAddGet, a[j])
			fmt.Printf("============\nAfter SetAuthorisedAddresses, Test Address : %d with value : %s \n", j+1, testAddr)

			rsAuthorisedAddr := keeper.IsAuthorised(ctx, a[j])
			assert.True(t, rsAuthorisedAddr)
			fmt.Printf("Is an Authorised Address\n")
		}

	}

	/// ------------------------------------case3: REMOVE - GET
	dropTestAddress := []sdkTypes.AccAddress{a[0]}
	keeper.RemoveAuthorisedAddresses(ctx, dropTestAddress)
	rsRemoveGet := keeper.GetAuthorisedAddresses(ctx)

	if assert.NotNil(t, rsRemoveGet) {

		for j, testAddr := range testAddrs {
			if j == 0 {
				fmt.Printf("============\nAfter RemoveAuthorisedAddresses, Test Address : %d with value : %s \n", j+1, testAddr)
				rsAuthorisedAddr := keeper.IsAuthorised(ctx, a[j])

				assert.False(t, rsAuthorisedAddr)
				fmt.Printf("Is not an Authorised Address\n")
			} else {
				assert.Contains(t, rsRemoveGet, a[j])
				fmt.Printf("============\nAfter RemoveAuthorisedAddresses, Test Address : %d with value : %s \n", j+1, testAddr)
				rsAuthorisedAddr := keeper.IsAuthorised(ctx, a[j])

				assert.True(t, rsAuthorisedAddr)
				fmt.Printf("Is an Authorised Address\n")
			}

		}

	}

}

func TestTokenIssuerAddresses(t *testing.T) {

	fmt.Printf("============\nStart Test : %s \n", "TestTokenIssuerAddresses")

	ctx, keeper := PrepareTest(t)

	var testAddrs = []string{
		"mxw1yw6mg7fty4mzcwupvzek53x5egm7tp2ldwaxq3",
		"mxw1yyz3h9calxmvjp4x05nnn70a8ex7fee3th7r7k",
		"mxw12vptxnfmt88rjyade6l3uzjeth64w005mdwhjj",
	}

	var a []sdkTypes.AccAddress
	/// ------------------------------------case1: IsIssuer ?
	for i, testAddr := range testAddrs {
		i = i + 1

		rsAddr, err := sdkTypes.AccAddressFromBech32(testAddr)
		assert.NoError(t, err)
		fmt.Printf("============\nTest Address : %d with value : %s \n", i, testAddr)

		rsUnissuerAddr := keeper.IsIssuer(ctx, rsAddr)
		assert.False(t, rsUnissuerAddr)
		fmt.Printf("Is not an Issuer Address\n")

		a = append(a, rsAddr)

	}

	/// ------------------------------------case2: ADD -> GET
	keeper.SetIssuerAddresses(ctx, a)
	rsAddGet := keeper.GetIssuerAddresses(ctx)

	if assert.NotNil(t, rsAddGet) {
		for j, testAddr := range testAddrs {
			assert.Contains(t, rsAddGet, a[j])
			fmt.Printf("============\nAfter SetIssuerAddresses, Test Address : %d with value : %s \n", j+1, testAddr)

			rsIssuerAddr := keeper.IsIssuer(ctx, a[j])
			assert.True(t, rsIssuerAddr)
			fmt.Printf("Is an Issuer Address\n")
		}

	}

	/// ------------------------------------case3: REMOVE - GET
	dropTestAddress := []sdkTypes.AccAddress{a[0]}
	keeper.RemoveIssuerAddresses(ctx, dropTestAddress)
	rsRemoveGet := keeper.GetIssuerAddresses(ctx)

	if assert.NotNil(t, rsRemoveGet) {

		for j, testAddr := range testAddrs {
			if j == 0 {
				fmt.Printf("============\nAfter RemoveIssuerAddresses, Test Address : %d with value : %s \n", j+1, testAddr)
				rsIssuerAddr := keeper.IsIssuer(ctx, a[j])

				assert.False(t, rsIssuerAddr)
				fmt.Printf("Is not an Issuer Address\n")
			} else {
				assert.Contains(t, rsRemoveGet, a[j])
				fmt.Printf("============\nAfter RemoveIssuerAddresses, Test Address : %d with value : %s \n", j+1, testAddr)
				rsIssuerAddr := keeper.IsIssuer(ctx, a[j])

				assert.True(t, rsIssuerAddr)
				fmt.Printf("Is an Issuer Address\n")
			}

		}

	}

}

func TestTokenProviderAddresses(t *testing.T) {

	fmt.Printf("============\nStart Test : %s \n", "TestTokenProviderAddresses")

	ctx, keeper := PrepareTest(t)

	var testAddrs = []string{
		"mxw1yw6mg7fty4mzcwupvzek53x5egm7tp2ldwaxq3",
		"mxw1yyz3h9calxmvjp4x05nnn70a8ex7fee3th7r7k",
		"mxw12vptxnfmt88rjyade6l3uzjeth64w005mdwhjj",
	}

	var a []sdkTypes.AccAddress
	/// ------------------------------------case1: IsIssuer ?
	for i, testAddr := range testAddrs {
		i = i + 1

		rsAddr, err := sdkTypes.AccAddressFromBech32(testAddr)
		assert.NoError(t, err)
		fmt.Printf("============\nTest Address : %d with value : %s \n", i, testAddr)

		rsUnproviderAddr := keeper.IsProvider(ctx, rsAddr)
		assert.False(t, rsUnproviderAddr)
		fmt.Printf("Is not an Provider Address\n")

		a = append(a, rsAddr)

	}

	/// ------------------------------------case2: ADD -> GET
	keeper.SetProviderAddresses(ctx, a)
	rsAddGet := keeper.GetProviderAddresses(ctx)

	if assert.NotNil(t, rsAddGet) {
		for j, testAddr := range testAddrs {
			assert.Contains(t, rsAddGet, a[j])
			fmt.Printf("============\nAfter SetProviderAddresses, Test Address : %d with value : %s \n", j+1, testAddr)

			rsProviderAddr := keeper.IsProvider(ctx, a[j])
			assert.True(t, rsProviderAddr)
			fmt.Printf("Is an Provider Address\n")
		}

	}

	/// ------------------------------------case3: REMOVE - GET
	dropTestAddress := []sdkTypes.AccAddress{a[0]}
	keeper.RemoveProviderAddresses(ctx, dropTestAddress)
	rsRemoveGet := keeper.GetProviderAddresses(ctx)

	if assert.NotNil(t, rsRemoveGet) {

		for j, testAddr := range testAddrs {
			if j == 0 {
				fmt.Printf("============\nAfter RemoveProviderAddresses, Test Address : %d with value : %s \n", j+1, testAddr)
				rsProviderAddr := keeper.IsProvider(ctx, a[j])

				assert.False(t, rsProviderAddr)
				fmt.Printf("Is not an Provider Address\n")
			} else {
				assert.Contains(t, rsRemoveGet, a[j])
				fmt.Printf("============\nAfter RemoveProviderAddresses, Test Address : %d with value : %s \n", j+1, testAddr)
				rsProviderAddr := keeper.IsProvider(ctx, a[j])

				assert.True(t, rsProviderAddr)
				fmt.Printf("Is an Provider Address\n")
			}

		}

	}

}
