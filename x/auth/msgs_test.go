package auth

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func MakeTestCodec() *codec.Codec {

	var cdc = codec.New()

	RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	sdkTypes.RegisterCodec(cdc)

	codec.RegisterCrypto(cdc)
	return cdc
}
func TestEncode(t *testing.T) {

	cdc := MakeTestCodec()

	msg1 := NewMsgCreateMultiSigAccount(sdkTypes.AccAddress{1}, 1, []sdkTypes.AccAddress{sdkTypes.AccAddress{1}, sdkTypes.AccAddress{2}})
	stdtx := auth.StdTx{[]sdkTypes.Msg{msg1}, auth.StdFee{}, []auth.StdSignature{}, ""}
	bz, err := cdc.MarshalJSON(stdtx)
	assert.NoError(t, err)
	fmt.Println(string(bz))

	msg2 := NewMsgUpdateMultiSigAccount(sdkTypes.AccAddress{1}, sdkTypes.AccAddress{1}, 1, []sdkTypes.AccAddress{sdkTypes.AccAddress{1}, sdkTypes.AccAddress{2}})
	stdtx = auth.StdTx{[]sdkTypes.Msg{msg2}, auth.StdFee{}, []auth.StdSignature{}, ""}
	bz, err2 := cdc.MarshalJSON(stdtx)
	assert.NoError(t, err2)
	fmt.Println(string(bz))

	msg3 := NewMsgTransferMultiSigOwner(sdkTypes.AccAddress{1}, sdkTypes.AccAddress{1}, sdkTypes.AccAddress{2})
	stdtx = auth.StdTx{[]sdkTypes.Msg{msg3}, auth.StdFee{}, []auth.StdSignature{}, ""}
	bz, err3 := cdc.MarshalJSON(stdtx)
	assert.NoError(t, err3)
	fmt.Println(string(bz))

	multiSigMsg := bank.NewMsgSend(sdkTypes.AccAddress{1}, sdkTypes.AccAddress{2}, sdkTypes.Coins{sdkTypes.Coin{}})
	multiSigStdTx := auth.StdTx{[]sdkTypes.Msg{multiSigMsg}, auth.StdFee{}, []auth.StdSignature{}, ""}

	msg4 := NewMsgCreateMultiSigTx(sdkTypes.AccAddress{1}, 1, multiSigStdTx, sdkTypes.AccAddress{2})
	stdtx = auth.StdTx{[]sdkTypes.Msg{msg4}, auth.StdFee{}, []auth.StdSignature{}, ""}
	bz, err4 := cdc.MarshalJSON(stdtx)
	assert.NoError(t, err4)
	fmt.Println(string(bz))
}
