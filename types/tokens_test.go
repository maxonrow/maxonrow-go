package types

import (
	"fmt"
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestCoinConversion(t *testing.T) {

	c1 := sdkTypes.Coin{
		Amount: sdkTypes.NewInt(123),
		Denom:  MXW,
	}
	c2, err := sdkTypes.ConvertCoin(c1, CIN)

	c3 := MXWtoCIN(c1)
	assert.NoError(t, err)
	assert.Equal(t, c2, c3)

	fmt.Println(c2.String())
	fmt.Println(c1.String())
}
