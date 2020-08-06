package utils

import (
	"fmt"
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/stretchr/testify/assert"
)

func TestDerive(t *testing.T) {

	// ---------------
	config := sdkTypes.GetConfig()
	config.SetBech32PrefixForAccount(types.Bech32PrefixAccAddr, types.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(types.Bech32PrefixValAddr, types.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr, types.Bech32PrefixConsPub)
	config.SetCoinType(376)
	config.SetFullFundraiserPath("44'/376'/0'/0/0")
	config.SetKeyringServiceName("mxw")

	config.Seal()
	// ---------------

	delAddr1, _ := sdkTypes.AccAddressFromBech32("mxw1dk6d4g8gxcy3terzwfffn7cx5hn6sukg8xu6np")

	derivedRes1, _ := sdkTypes.AccAddressFromBech32("mxw18wthtgrc32p372uq9d8pue7xyqhvzvjvfhw0mt")
	derivedAddr := DeriveMultiSigAddress(delAddr1, 0x987654321)
	fmt.Println(derivedAddr.String())
	assert.Equal(t, derivedAddr, derivedRes1)

	derivedRes2, _ := sdkTypes.AccAddressFromBech32("mxw1c9galkrvda7x4gd744c4fp0de5ymaql5e8lcdr")
	derivedAddr2 := DeriveMultiSigAddress(delAddr1, 0x98765432100)
	fmt.Println(derivedAddr2.String())
	assert.Equal(t, derivedAddr2, derivedRes2)

	assert.NotEqual(t, derivedAddr, derivedAddr2)

}
