package maintenance

import (
	"encoding/json"
	"fmt"
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {

	msg1 := NewFeeeMaintainer("add", []sdkTypes.AccAddress{sdkTypes.AccAddress{1}}, []FeeCollector{FeeCollector{
		"module", sdkTypes.AccAddress{2},
	}})
	bz, err := json.Marshal(msg1)
	assert.NoError(t, err)
	fmt.Println(string(bz))

	msg2 := NewKycMaintainer("add", []sdkTypes.AccAddress{sdkTypes.AccAddress{1}}, []sdkTypes.AccAddress{sdkTypes.AccAddress{2}}, []sdkTypes.AccAddress{sdkTypes.AccAddress{2}})
	bz, err2 := json.Marshal(msg2)
	assert.NoError(t, err2)
	fmt.Println(string(bz))
}
