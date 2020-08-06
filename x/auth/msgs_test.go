package auth

import (
	"fmt"
	"testing"

	//"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"

	//sdkBank "github.com/cosmos/cosmos-sdk/x/bank"
	bank "github.com/maxonrow/maxonrow-go/x/bank"
)

func MakeTestCodec() *codec.Codec {

	var cdc = msgCdc

	return cdc
}
func TestEncodeJSON(t *testing.T) {

	cdc := MakeTestCodec()

	delPk1 := ed25519.GenPrivKey().PubKey()
	delPk2 := ed25519.GenPrivKey().PubKey()
	delAddr1 := sdkTypes.AccAddress(delPk1.Address())
	delAddr2 := sdkTypes.AccAddress(delPk2.Address())

	addr1 := sdkTypes.AccAddress(delAddr1)
	addr2 := sdkTypes.AccAddress(delAddr2)

	msg1 := NewMsgCreateMultiSigAccount(addr1, 1, []sdkTypes.AccAddress{addr1, addr2})
	stdtx := sdkAuth.StdTx{[]sdkTypes.Msg{msg1}, sdkAuth.StdFee{}, nil, ""}
	bz, err := cdc.MarshalJSON(stdtx)
	assert.NoError(t, err)
	msg11 := new(sdkAuth.StdTx)
	err = cdc.UnmarshalJSON(bz, msg11)
	fmt.Println(string(bz))
	assert.NoError(t, err)
	assert.Equal(t, stdtx, *msg11)

	msg2 := NewMsgUpdateMultiSigAccount(addr1, addr1, 1, []sdkTypes.AccAddress{addr1, addr2})
	stdtx = sdkAuth.StdTx{[]sdkTypes.Msg{msg2}, sdkAuth.StdFee{}, nil, ""}
	bz, err2 := cdc.MarshalJSON(stdtx)
	assert.NoError(t, err2)
	msg22 := new(sdkAuth.StdTx)
	cdc.UnmarshalJSON(bz, msg22)
	assert.Equal(t, stdtx, *msg22)

	msg3 := NewMsgTransferMultiSigOwner(addr1, addr1, addr2)
	stdtx = sdkAuth.StdTx{[]sdkTypes.Msg{msg3}, sdkAuth.StdFee{}, nil, ""}
	bz, err3 := cdc.MarshalJSON(stdtx)
	assert.NoError(t, err3)
	msg33 := new(sdkAuth.StdTx)
	cdc.UnmarshalJSON(bz, msg33)
	assert.Equal(t, stdtx, *msg33)

	multiSigMsg := bank.NewMsgSend(addr1, addr2, sdkTypes.Coins{sdkTypes.NewInt64Coin("cin", 1)})
	multiSigStdTx := sdkAuth.NewStdTx([]sdkTypes.Msg{multiSigMsg}, sdkAuth.StdFee{}, nil, "")
	bz1, err4 := cdc.MarshalJSON(multiSigStdTx)
	assert.NoError(t, err4)
	msg44 := new(sdkAuth.StdTx)
	cdc.UnmarshalJSON(bz1, msg44)
	assert.Equal(t, multiSigStdTx, *msg44)

}

func TestEncodeAmino(t *testing.T) {

	cdc := MakeTestCodec()
	addr1, _ := sdkTypes.AccAddressFromHex("0x1234")
	addr2, _ := sdkTypes.AccAddressFromHex("0x5678")

	msg1 := NewMsgCreateMultiSigAccount(addr1, 1, []sdkTypes.AccAddress{addr1, addr2})
	stdtx := sdkAuth.StdTx{[]sdkTypes.Msg{msg1}, sdkAuth.StdFee{}, nil, ""}
	bz, err := cdc.MarshalBinaryBare(stdtx)
	assert.NoError(t, err)
	msg11 := new(sdkAuth.StdTx)
	cdc.UnmarshalBinaryBare(bz, msg11)
	assert.Equal(t, stdtx, *msg11)

	msg2 := NewMsgUpdateMultiSigAccount(addr1, addr1, 1, []sdkTypes.AccAddress{addr1, addr2})
	stdtx = sdkAuth.StdTx{[]sdkTypes.Msg{msg2}, sdkAuth.StdFee{}, nil, ""}
	bz, err2 := cdc.MarshalBinaryBare(stdtx)
	assert.NoError(t, err2)
	msg22 := new(sdkAuth.StdTx)
	cdc.UnmarshalBinaryBare(bz, msg22)
	assert.Equal(t, stdtx, *msg22)

	msg3 := NewMsgTransferMultiSigOwner(addr1, addr1, addr2)
	stdtx = sdkAuth.StdTx{[]sdkTypes.Msg{msg3}, sdkAuth.StdFee{}, nil, ""}
	bz, err3 := cdc.MarshalBinaryBare(stdtx)
	assert.NoError(t, err3)
	msg33 := new(sdkAuth.StdTx)
	cdc.UnmarshalBinaryBare(bz, msg33)
	assert.Equal(t, stdtx, *msg33)

	multiSigMsg := bank.NewMsgSend(addr1, addr2, sdkTypes.Coins{sdkTypes.NewInt64Coin("cin", 1)})
	multiSigStdTx := sdkAuth.NewStdTx([]sdkTypes.Msg{multiSigMsg}, sdkAuth.StdFee{}, nil, "")
	bz1, err4 := cdc.MarshalBinaryBare(multiSigStdTx)
	assert.NoError(t, err4)

	msg44 := new(sdkAuth.StdTx)
	cdc.UnmarshalBinaryBare(bz1, msg44)
	assert.Equal(t, multiSigStdTx, *msg44)

}
