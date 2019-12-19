package types

import (
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestAddressHolder(t *testing.T) {
	var a AddressHolder

	addr1 := sdkTypes.AccAddress{1}
	addr2 := sdkTypes.AccAddress{2}
	addr3 := sdkTypes.AccAddress{3}
	addr4 := sdkTypes.AccAddress{4}
	addr5 := sdkTypes.AccAddress{5}

	s:= []sdkTypes.AccAddress{addr3,addr4,addr5}

	assert.True(t, a.Append(addr1))
	assert.True(t, a.Append(addr2))
	assert.True(t, a.Append(addr3))
	assert.False(t, a.Append(addr2))
	assert.Equal(t, a.Size(), 3)
	assert.True(t, a.Remove(addr1))
	assert.Equal(t, a.Size(), 2)
	i, ok := a.Contains(addr2)
	assert.Equal(t, i, 0)
	assert.Equal(t, ok, true)
	i, ok = a.Contains(addr4)
	assert.Equal(t, i, -1)
	assert.Equal(t, ok, false)

	a.AppendAccAddrs(s)
	assert.Equal(t, a.Size(), 4)

}
