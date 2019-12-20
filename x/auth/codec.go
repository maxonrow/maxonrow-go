package auth

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateMultiSigAccount{}, "mxw/msgCreateMultiSigAccount", nil)
	cdc.RegisterConcrete(MsgUpdateMultiSigAccount{}, "mxw/msgUpdateMultiSigAccount", nil)
	cdc.RegisterConcrete(MsgTransferMultiSigOwner{}, "mxw/msgTransferMultiSigOwner", nil)
	cdc.RegisterConcrete(MsgCreateMultiSigTx{}, "mxw/msgCreateMultiSigTx", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
