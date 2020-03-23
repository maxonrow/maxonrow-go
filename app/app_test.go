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

	derivedAddr := auth.DeriveMultiSigAddress(delAddr1, 0x987654321)
	fmt.Println(derivedAddr.String())

	derivedAddr2 := auth.DeriveMultiSigAddress(delAddr1, 0x98765432100)
	fmt.Println(derivedAddr2.String())

	assert.NotEqual(t, derivedAddr, derivedAddr2)
}
