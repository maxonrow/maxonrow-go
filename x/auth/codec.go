package auth

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	bank "github.com/maxonrow/maxonrow-go/x/bank"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateMultiSigAccount{}, "mxw/msgCreateMultiSigAccount", nil)
	cdc.RegisterConcrete(MsgUpdateMultiSigAccount{}, "mxw/msgUpdateMultiSigAccount", nil)
	cdc.RegisterConcrete(MsgTransferMultiSigOwner{}, "mxw/msgTransferMultiSigOwner", nil)
	cdc.RegisterConcrete(MsgCreateMultiSigTx{}, "mxw/msgCreateMultiSigTx", nil)

	cdc.RegisterConcrete(MsgSignMultiSigTx{}, "mxw/msgSignMultiSigTx", nil)

}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)

	sdkTypes.RegisterCodec(msgCdc)
	sdkAuth.RegisterCodec(msgCdc)
	codec.RegisterCrypto(msgCdc)

	bank.RegisterCodec(msgCdc)
	//fungible.RegisterCodec(cdc)
	//nonFungible.RegisterCodec(cdc)
	//fee.RegisterCodec(cdc)
	//maintenance.RegisterCodec(cdc)
	//auth.RegisterCodec(cdc)
	
	// To register codec for internal transaction for multi-siog account (cosmos-sdk)
	bank.RegisterCodec(sdkAuth.ModuleCdc)
}
