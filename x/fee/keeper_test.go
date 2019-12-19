package fee_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"github.com/maxonrow/maxonrow-go/app"
	"github.com/maxonrow/maxonrow-go/x/fee"
)

type AddressHolder []sdkTypes.AccAddress

// Init : create DB object
func defaultContext(key sdkTypes.StoreKey) sdkTypes.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdkTypes.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := sdkTypes.NewContext(cms, abci.Header{}, false, log.NewNopLogger())

	return ctx
}

func PrepareTest(t *testing.T) (sdkTypes.Context, fee.Keeper) {

	// Getting default codec for marshaling and unmarshaling
	cdc := app.MakeDefaultCodec()

	// Create key store for fee keeper
	key := sdkTypes.NewKVStoreKey("fee")

	// Getting context for fee
	ctx := defaultContext(key)

	// Creating fee keeper instance
	keeper := fee.NewKeeper(cdc, key)

	return ctx, keeper
}

func TestAuthorised(t *testing.T) {

	fmt.Printf("============\nStart Test : %s \n", "TestAuthorised")

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
	rs_addGet := keeper.GetAuthorisedAddresses(ctx)

	if assert.NotNil(t, rs_addGet) {
		for j, testAddr := range testAddrs {
			assert.Contains(t, rs_addGet, a[j])
			fmt.Printf("============\nAfter SetAuthorisedAddresses, Test Address : %d with value : %s \n", j+1, testAddr)

			rsAuthorisedAddr := keeper.IsAuthorised(ctx, a[j])
			assert.True(t, rsAuthorisedAddr)
			fmt.Printf("Is an Authorised Address\n")
		}

	}

	/// ------------------------------------case3: REMOVE - GET
	dropTestAddress := []sdkTypes.AccAddress{a[0]}
	keeper.RemoveAuthorisedAddresses(ctx, dropTestAddress)
	rs_removeGet := keeper.GetAuthorisedAddresses(ctx)

	if assert.NotNil(t, rs_removeGet) {

		for j, testAddr := range testAddrs {
			if j == 0 {
				fmt.Printf("============\nAfter RemoveAuthorisedAddresses, Test Address : %d with value : %s \n", j+1, testAddr)
				rsAuthorisedAddr := keeper.IsAuthorised(ctx, a[j])

				assert.False(t, rsAuthorisedAddr)
				fmt.Printf("Is not an Authorised Address\n")
			} else {
				assert.Contains(t, rs_removeGet, a[j])
				fmt.Printf("============\nAfter RemoveAuthorisedAddresses, Test Address : %d with value : %s \n", j+1, testAddr)
				rsAuthorisedAddr := keeper.IsAuthorised(ctx, a[j])

				assert.True(t, rsAuthorisedAddr)
				fmt.Printf("Is an Authorised Address\n")
			}

		}

	}

}

func TestFeeCollector(t *testing.T) {

	fmt.Printf("============\nStart Test : %s \n", "TestFeeCollector")

	ctx, keeper := PrepareTest(t)

	var testFeeCollectors = []string{
		"mxw1yw6mg7fty4mzcwupvzek53x5egm7tp2ldwaxq3",
		"mxw1yyz3h9calxmvjp4x05nnn70a8ex7fee3th7r7k",
		"mxw12vptxnfmt88rjyade6l3uzjeth64w005mdwhjj",
	}

	module := "fee"
	var a []sdkTypes.AccAddress
	/// ------------------------------------case1: IsFeeCollector ?
	for i, testFeeCollector := range testFeeCollectors {
		i = i + 1

		rsFeeCollector, err := sdkTypes.AccAddressFromBech32(testFeeCollector)
		assert.NoError(t, err)
		fmt.Printf("============\nTest FeeCollector : %d with value : %s \n", i, testFeeCollector)

		rsInvalidFeeCollector := keeper.IsFeeCollector(ctx, module, rsFeeCollector)
		assert.False(t, rsInvalidFeeCollector)
		fmt.Printf("Is not valid FeeCollector\n")

		a = append(a, rsFeeCollector)

	}

	/// ------------------------------------case2: ADD -> GET
	keeper.SetFeeCollectorAddresses(ctx, module, a)
	rs_addGet := keeper.GetFeeCollectorAddresses(ctx, module)

	if assert.NotNil(t, rs_addGet) {
		for j, testFeeCollector := range testFeeCollectors {
			assert.Contains(t, rs_addGet, a[j])
			fmt.Printf("============\nAfter SetFeeCollector, Test FeeCollector : %d with value : %s \n", j+1, testFeeCollector)

			rsValidFeeCollector := keeper.IsFeeCollector(ctx, module, a[j])
			assert.True(t, rsValidFeeCollector)
			fmt.Printf("Is valid FeeCollector\n")
		}

	}

	/// ------------------------------------case3: REMOVE - GET
	dropFeeCollector := a[0]
	keeper.RemoveFeeCollectorAddress(ctx, module, dropFeeCollector)
	rs_removeGet := keeper.GetFeeCollectorAddresses(ctx, module)

	if assert.NotNil(t, rs_removeGet) {

		for j, testFeeCollector := range testFeeCollectors {
			if j == 0 {
				fmt.Printf("============\nAfter RemoveFeeCollectorAddress, Test FeeCollector : %d with value : %s \n", j+1, testFeeCollector)
				rsValidFeeCollector := keeper.IsFeeCollector(ctx, module, a[j])

				assert.False(t, rsValidFeeCollector)
				fmt.Printf("Is not valid FeeCollector\n")
			} else {
				assert.Contains(t, rs_removeGet, a[j])
				fmt.Printf("============\nAfter RemoveFeeCollectorAddress, Test FeeCollector : %d with value : %s \n", j+1, testFeeCollector)
				rsValidFeeCollector := keeper.IsFeeCollector(ctx, module, a[j])

				assert.True(t, rsValidFeeCollector)
				fmt.Printf("Is valid FeeCollector\n")
			}

		}

	}

}

func TestGetMsgFeeSetting(t *testing.T) {

	fmt.Printf("============\nStart Test : %s \n", "TestGetMsgFeeSetting")

	ctx, keeper := PrepareTest(t)

	var testMsgFeeSettingTypes = []string{
		"",
		"default_9991",
		"msg fee setting returns",
	}

	for i, testMsgFeeSettingType := range testMsgFeeSettingTypes {
		rs_get, err := keeper.GetMsgFeeSetting(ctx, testMsgFeeSettingType)

		if i == 0 {
			if err == nil {
				assert.NotEqual(t, rs_get.Name, "default")
				fmt.Printf("============\nTest GetMsgFeeSetting Type : %d with value : %s \n", i+1, rs_get.Name)
			} else {
				assert.Error(t, err)
				fmt.Printf("============\nTest GetMsgFeeSetting Type : %s \n", err)
			}

		} else {
			if err == nil {
				if assert.NotNil(t, rs_get) {
					assert.NotContains(t, rs_get.Name, testMsgFeeSettingType)
					fmt.Printf("============\nTest GetMsgFeeSetting Type : %d with value : %v \n", i+1, testMsgFeeSettingType)
				} else {
					assert.Error(t, err)
					fmt.Printf("============\nTest GetMsgFeeSetting Type : %s \n", err)
				}
			} else {
				assert.Error(t, err)
				fmt.Printf("============\nTest GetMsgFeeSetting Type : %s \n", err)
			}
		}

	}

}

func TestFeeSetting(t *testing.T) {

	fmt.Printf("============\nStart Test : %s \n", "TestFeeSetting")

	ctx, keeper := PrepareTest(t)

	/// ------------------------------------case1: FeeSettingExists ?
	var testFeeSettingTypes = []string{
		"fee_default",
		"fee_admin",
		"fee_tx",
	}

	for i, testFeeSettingType := range testFeeSettingTypes {
		rsIsExist := keeper.FeeSettingExists(ctx, testFeeSettingType)
		assert.False(t, rsIsExist)
		fmt.Printf("============\nTest FeeSetting : %d with value : %s \n", i+1, testFeeSettingType)
		fmt.Printf("Is not existed\n")

	}

	/// ------------------------------------case2: ADD -> GET -> Existed
	var testIssuerAddrs = []string{
		"mxw1yw6mg7fty4mzcwupvzek53x5egm7tp2ldwaxq3",
		"mxw1yyz3h9calxmvjp4x05nnn70a8ex7fee3th7r7k",
		"mxw12vptxnfmt88rjyade6l3uzjeth64w005mdwhjj",
	}

	rsIssuerAddrs1, _ := sdkTypes.AccAddressFromBech32(testIssuerAddrs[0])
	rsIssuerAddrs2, _ := sdkTypes.AccAddressFromBech32(testIssuerAddrs[1])
	rsIssuerAddrs3, _ := sdkTypes.AccAddressFromBech32(testIssuerAddrs[2])

	var a []sdkTypes.AccAddress
	a = append(a, rsIssuerAddrs1)
	a = append(a, rsIssuerAddrs2)
	a = append(a, rsIssuerAddrs3)

	keeper.SetAuthorisedAddresses(ctx, a)
	rs_addGet := keeper.GetAuthorisedAddresses(ctx)

	if assert.NotNil(t, rs_addGet) {

		// 1. check if is authorised
		for j, item := range a {
			assert.Contains(t, rs_addGet, item)
			fmt.Printf("============\nAfter SetAuthorisedAddresses, Test Address : %d with value : %s \n", j+1, item)

			rsAuthorisedAddr := keeper.IsAuthorised(ctx, item)
			assert.True(t, rsAuthorisedAddr)
			fmt.Printf("Is an Authorised Address\n")

		}

		// 2. start Create FeeSetting
		processCreateFeeSetting(t, ctx, keeper, rsIssuerAddrs1, rsIssuerAddrs2, rsIssuerAddrs3)
		feeSettings := keeper.ListAllSysFeeSetting(ctx)
		fmt.Println(feeSettings)
	}

}

func processCreateFeeSetting(t *testing.T, ctx sdkTypes.Context, keeper fee.Keeper, rsIssuerAddrs1 sdkTypes.AccAddress, rsIssuerAddrs2 sdkTypes.AccAddress, rsIssuerAddrs3 sdkTypes.AccAddress) {
	amtMin_IssuerAddrs1 := sdkTypes.Coins{
		{
			Denom:  "cin",
			Amount: sdkTypes.NewInt(1000000),
		},
	}

	amtMax_IssuerAddrs1 := sdkTypes.Coins{
		{
			Denom:  "cin",
			Amount: sdkTypes.NewInt(5000000),
		},
	}

	amtMin_IssuerAddrs2 := sdkTypes.Coins{
		{
			Denom:  "cin",
			Amount: sdkTypes.NewInt(1000000),
		},
	}

	amtMax_IssuerAddrs2 := sdkTypes.Coins{
		{
			Denom:  "cin",
			Amount: sdkTypes.NewInt(1000000),
		},
	}

	amtMin_IssuerAddrs3 := sdkTypes.Coins{
		{
			Denom:  "cin",
			Amount: sdkTypes.NewInt(1000000),
		},
	}

	amtMax_IssuerAddrs3 := sdkTypes.Coins{
		{
			Denom:  "cin",
			Amount: sdkTypes.NewInt(900000),
		},
	}

	var testFeeSettings = []fee.MsgSysFeeSetting{
		{
			Name:       "fee_default",
			Min:        amtMin_IssuerAddrs1,
			Max:        amtMax_IssuerAddrs1,
			Percentage: "0",
			Issuer:     rsIssuerAddrs1,
		},
		{
			Name:       "fee_admin",
			Min:        amtMin_IssuerAddrs2,
			Max:        amtMax_IssuerAddrs2,
			Percentage: "3",
			Issuer:     rsIssuerAddrs2,
		},
		{
			Name:       "fee_tx",
			Min:        amtMin_IssuerAddrs3,
			Max:        amtMax_IssuerAddrs3,
			Percentage: "88",
			Issuer:     rsIssuerAddrs3,
		},
	}

	fs1, err1 := keeper.GetFeeSettingByName(ctx, "")
	assert.Nil(t, fs1)
	assert.Error(t, err1)

	for i, testFeeSetting := range testFeeSettings {
		keeper.CreateFeeSetting(ctx, testFeeSetting)
		rs_addGet, _ := keeper.GetFeeSettingByName(ctx, testFeeSetting.Name)

		if assert.NotNil(t, rs_addGet) {

			fmt.Printf("============\nAfter CreateFeeSetting, Test FeeSetting : %d with value : %s \n", i+1, rs_addGet.Name)

			rsIsExist := keeper.FeeSettingExists(ctx, rs_addGet.Name)
			assert.True(t, rsIsExist)
			fmt.Printf("Is existed\n")

		}

	}

}

func TestTxFeeSetting(t *testing.T) {

	fmt.Printf("============\nStart Test : %s \n", "TestTxFeeSetting")

	ctx, keeper := PrepareTest(t)

	/// ------------------------------------case2: ADD
	var testIssuerAddrs = []string{
		"mxw1yw6mg7fty4mzcwupvzek53x5egm7tp2ldwaxq3",
		"mxw1yyz3h9calxmvjp4x05nnn70a8ex7fee3th7r7k",
		"mxw12vptxnfmt88rjyade6l3uzjeth64w005mdwhjj",
	}

	rsIssuerAddrs1, _ := sdkTypes.AccAddressFromBech32(testIssuerAddrs[0])
	rsIssuerAddrs2, _ := sdkTypes.AccAddressFromBech32(testIssuerAddrs[1])
	rsIssuerAddrs3, _ := sdkTypes.AccAddressFromBech32(testIssuerAddrs[2])

	var testMsgAssignFeeToMsgs = []struct {
		FeeName string              `json:"fee_name"`
		MsgType string              `json:"msg_type"`
		Issuer  sdkTypes.AccAddress `json:"issuer"`
	}{
		{
			FeeName: "chargeAdmin",
			MsgType: "fee_admin",
			Issuer:  rsIssuerAddrs1,
		},
		{
			FeeName: "chargeTx",
			MsgType: "fee_tx",
			Issuer:  rsIssuerAddrs2,
		},
		{
			FeeName: "default",
			MsgType: "fee_default",
			Issuer:  rsIssuerAddrs3,
		},
	}

	var a []sdkTypes.AccAddress
	a = append(a, rsIssuerAddrs1)
	a = append(a, rsIssuerAddrs2)
	a = append(a, rsIssuerAddrs3)

	keeper.SetAuthorisedAddresses(ctx, a)
	rs_addGet := keeper.GetAuthorisedAddresses(ctx)

	if assert.NotNil(t, rs_addGet) {

		// 1. check if is authorised
		for j, item := range a {
			assert.Contains(t, rs_addGet, item)
			fmt.Printf("============\nAfter SetAuthorisedAddresses, Test Address : %d with value : %s \n", j+1, item)

			rsAuthorisedAddr := keeper.IsAuthorised(ctx, item)
			assert.True(t, rsAuthorisedAddr)
			fmt.Printf("Is an Authorised Address\n")

		}

		// 2. start Create FeeSetting
		processCreateFeeSetting(t, ctx, keeper, rsIssuerAddrs1, rsIssuerAddrs2, rsIssuerAddrs3)

		// 3. start Create TxFeeSetting
		for i, testMsgAssignFeeToMsg := range testMsgAssignFeeToMsgs {
			rs_add := keeper.AssignFeeToMsg(ctx, testMsgAssignFeeToMsg)

			if assert.NotNil(t, rs_add) {

				fmt.Printf("============\nAfter CreateTxFeeSetting, Test TxFeeSetting : %d with value : %s \n", i+1, rs_add.Log)

			}

		}

	}

}
