package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateAlias{}, "nameservice/createAlias", nil)
	cdc.RegisterConcrete(MsgSetAliasStatus{}, "nameservice/setAliasStatus", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
	codec.RegisterCrypto(msgCdc)
}
