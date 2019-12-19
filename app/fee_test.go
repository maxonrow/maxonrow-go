package app

import (
	"fmt"
	"testing"
	"time"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/maxonrow/maxonrow-go/x/fee"
)

func TestFeeCalc(t *testing.T) {
	ctx := sdkTypes.Context{}
	min1, _ := sdkTypes.ParseCoins("300000000000000000cin")
	max1, _ := sdkTypes.ParseCoins("1598765330000000000cin")
	feeSetting1 := &fee.FeeSetting{
		Name: "test1",

		Min:        min1,
		Max:        max1,
		Percentage: "0.5012345643",
	}

	amt1, _ := sdkTypes.ParseCoins("88123455432100000000cin")
	expectedFee1, _ := sdkTypes.ParseCoins("4369016736367332cin")
	fee1, err := calculateFee(ctx, feeSetting1, "0.00989125", amt1)
	fmt.Println(fee1)
	assert.NoError(t, err)
	assert.Equal(t, fee1, expectedFee1)

	min2, _ := sdkTypes.ParseCoins("50cin")
	max2, _ := sdkTypes.ParseCoins("100cin")
	feeSetting2 := &fee.FeeSetting{
		Name: "test2",

		Min:        min2,
		Max:        max2,
		Percentage: "0.50",
	}

	amt2, _ := sdkTypes.ParseCoins("11111cin")
	expectedFee2, _ := sdkTypes.ParseCoins("39cin")
	fee2, err := calculateFee(ctx, feeSetting2, "0.69", amt2)
	fmt.Println(fee2)
	assert.NoError(t, err)
	assert.Equal(t, fee2, expectedFee2)
}

func TestFeeCalcBadPercentage(t *testing.T) {
	ctx := sdkTypes.Context{}
	min, _ := sdkTypes.ParseCoins("1000000000000000cin,")
	max, _ := sdkTypes.ParseCoins("1000000000000000000000000cin")
	feeSetting1 := &fee.FeeSetting{
		Name: "test",

		Min:        min,
		Max:        max,
		Percentage: "0.5123",
	}

	amt1, _ := sdkTypes.ParseCoins("11111111111111111cin")
	expectedFee1, _ := sdkTypes.ParseCoins("56922222222222cin")

	fee1, err := calculateFee(ctx, feeSetting1, "1", amt1)
	fmt.Println(fee1)
	assert.NoError(t, err)
	assert.Equal(t, fee1, expectedFee1)
}

func TestFeeCalcBadMax(t *testing.T) {
	ctx := sdkTypes.Context{}
	min, _ := sdkTypes.ParseCoins("1000000000000000cin,")
	max, _ := sdkTypes.ParseCoins("100000cin")
	feeSetting1 := &fee.FeeSetting{
		Name: "test",

		Min:        min,
		Max:        max,
		Percentage: "0.5123",
	}

	amt1, _ := sdkTypes.ParseCoins("11111111111111111cin")
	expectedFee1, _ := sdkTypes.ParseCoins("100000cin")

	fee1, err := calculateFee(ctx, feeSetting1, "1", amt1)
	fmt.Println(fee1)
	assert.NoError(t, err)
	assert.Equal(t, fee1, expectedFee1)
}

func TestFeeCalc100mxw(t *testing.T) {
	ctx := sdkTypes.Context{}
	min, _ := sdkTypes.ParseCoins("1000000000000000cin,")
	max, _ := sdkTypes.ParseCoins("1000000000000000000000000cin")
	feeSetting1 := &fee.FeeSetting{
		Name: "test",

		Min:        min,
		Max:        max,
		Percentage: "0.12587",
	}

	amt1, _ := sdkTypes.ParseCoins("100000000000000000000cin")
	expectedFee1, _ := sdkTypes.ParseCoins("124316046741000cin")
	expectedFee2, _ := sdkTypes.ParseCoins("251740000000000000cin")

	fee1, err := calculateFee(ctx, feeSetting1, "0.0009876543", amt1)
	assert.NoError(t, err)
	assert.Equal(t, fee1, expectedFee1)

	fee2, err := calculateFee(ctx, feeSetting1, "2", amt1)
	assert.NoError(t, err)
	assert.Equal(t, fee2, expectedFee2)
}

func TestFeeCalc1000890000000000cin(t *testing.T) {
	ctx := sdkTypes.Context{}
	min, _ := sdkTypes.ParseCoins("1000000000000000cin,")
	max, _ := sdkTypes.ParseCoins("1000000000000000000000000cin")
	feeSetting1 := &fee.FeeSetting{
		Name: "test",

		Min:        min,
		Max:        max,
		Percentage: "0.12587",
	}

	amt1, _ := sdkTypes.ParseCoins("1000890000000000cin") // 100mxw // 100000000000000000000cin
	expectedFee1, _ := sdkTypes.ParseCoins("1244266880cin")
	expectedFee2, _ := sdkTypes.ParseCoins("2519640486000cin")

	fee1, err := calculateFee(ctx, feeSetting1, "0.0009876543", amt1)
	assert.NoError(t, err)
	assert.Equal(t, fee1, expectedFee1)

	fee2, err := calculateFee(ctx, feeSetting1, "2", amt1)
	assert.NoError(t, err)
	assert.Equal(t, fee2, expectedFee2)
}

func TestCalculateBigInt(t *testing.T) {
	fmt.Println(time.Now())

	var min, max, amt, percentage, adjuster sdkTypes.Int
	min, _ = sdkTypes.NewIntFromString("10")
	max, _ = sdkTypes.NewIntFromString("10")
	amt, _ = sdkTypes.NewIntFromString("10")
	percentage, _ = sdkTypes.NewIntFromString("50")
	adjuster, _ = sdkTypes.NewIntFromString("10000")

	var fee sdkTypes.Int
	for i := 0; i < 1000000; i++ {
		fee = amt.Mul(percentage).Quo(adjuster)
	}
	if fee.LT(min) {
		fee = min
	}

	if fee.GT(max) {
		fee = max
	}

	fmt.Println(fee)
	fmt.Println(time.Now())

	assert.Equal(t, fee, fee)

}

func TestCalculateInt64(t *testing.T) {
	fmt.Println(time.Now())

	var min, max, amt, percentage, adjuster, fee int64
	min = 10
	max = 10
	amt = 10
	percentage = 50
	adjuster = 10000

	for i := 0; i < 1000000; i++ {
		fee = amt * percentage / adjuster
	}
	if fee < min {
		fee = min
	}

	if fee > max {
		fee = max
	}

	fmt.Println(fee)
	fmt.Println(time.Now())
}
