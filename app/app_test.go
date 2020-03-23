package app

import (
	"fmt"
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/x/auth"
	"github.com/stretchr/testify/assert"
)

func TestDerive(t *testing.T) {

	delAddr1, _ := sdkTypes.AccAddressFromBech32("mxw1dk6d4g8gxcy3terzwfffn7cx5hn6sukg8xu6np")

	derivedRes1, _ := sdkTypes.AccAddressFromBech32("mxw18wthtgrc32p372uq9d8pue7xyqhvzvjvfhw0mt")
	derivedAddr := auth.DeriveMultiSigAddress(delAddr1, 0x987654321)
	fmt.Println(derivedAddr.String())
	assert.Equal(t, derivedAddr, derivedRes1)

	derivedRes2, _ := sdkTypes.AccAddressFromBech32("mxw1c9galkrvda7x4gd744c4fp0de5ymaql5e8lcdr")
	derivedAddr2 := auth.DeriveMultiSigAddress(delAddr1, 0x98765432100)
	fmt.Println(derivedAddr2.String())
	assert.Equal(t, derivedAddr2, derivedRes2)

	assert.NotEqual(t, derivedAddr, derivedAddr2)

}
