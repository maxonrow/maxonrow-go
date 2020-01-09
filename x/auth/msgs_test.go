package auth

import (
	"fmt"
	"testing"

	//"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"

	//sdkBank "github.com/cosmos/cosmos-sdk/x/bank"
	bank "github.com/maxonrow/maxonrow-go/x/bank"
)

func MakeTestCodec() *codec.Codec {

	var cdc = msgCdc

	return cdc
}
func TestEncode(t *testing.T) {

	cdc := MakeTestCodec()

	// msg1 := NewMsgCreateMultiSigAccount(sdkTypes.AccAddress{1}, 1, []sdkTypes.AccAddress{sdkTypes.AccAddress{1}, sdkTypes.AccAddress{2}})
	// stdtx := sdkAuth.StdTx{[]sdkTypes.Msg{msg1}, sdkAuth.StdFee{}, []sdkAuth.StdSignature{}, ""}
	// bz, err := cdc.MarshalJSON(stdtx)
	// assert.NoError(t, err)
	// fmt.Println(string(bz))

	// msg2 := NewMsgUpdateMultiSigAccount(sdkTypes.AccAddress{1}, sdkTypes.AccAddress{1}, 1, []sdkTypes.AccAddress{sdkTypes.AccAddress{1}, sdkTypes.AccAddress{2}})
	// stdtx = sdkAuth.StdTx{[]sdkTypes.Msg{msg2}, sdkAuth.StdFee{}, []sdkAuth.StdSignature{}, ""}
	// bz, err2 := cdc.MarshalJSON(stdtx)
	// assert.NoError(t, err2)
	// fmt.Println(string(bz))

	// msg3 := NewMsgTransferMultiSigOwner(sdkTypes.AccAddress{1}, sdkTypes.AccAddress{1}, sdkTypes.AccAddress{2})
	// stdtx = sdkAuth.StdTx{[]sdkTypes.Msg{msg3}, sdkAuth.StdFee{}, []sdkAuth.StdSignature{}, ""}
	// bz, err3 := cdc.MarshalJSON(stdtx)
	// assert.NoError(t, err3)
	// fmt.Println(string(bz))

	multiSigMsg := bank.NewMsgSend(sdkTypes.AccAddress{1}, sdkTypes.AccAddress{2}, sdkTypes.Coins{sdkTypes.Coin{}})
	multiSigStdTx := sdkAuth.NewStdTx([]sdkTypes.Msg{multiSigMsg}, sdkAuth.StdFee{}, []sdkAuth.StdSignature{}, "")

	msg4 := NewMsgCreateMultiSigTx(sdkTypes.AccAddress{1}, multiSigStdTx, sdkTypes.AccAddress{2})
	//stdtx = sdkAuth.NewStdTx([]sdkTypes.Msg{msg4}, sdkAuth.StdFee{}, []sdkAuth.StdSignature{}, "")

	bz1 := sdkTypes.MustSortJSON(cdc.MustMarshalJSON(msg4))
	bz2 := msg4.GetSignBytes()
	//assert.NoError(t, err4)
	fmt.Println(string(bz1))
	fmt.Println(string(bz2))

	a := 1
	b := a

	if a == 2 {
		fmt.Println("avc")
	} else if b == 1 {
		fmt.Println("aaa")
	}

}
