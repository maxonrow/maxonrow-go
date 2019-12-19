package fungible

import "github.com/cosmos/cosmos-sdk/codec"

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateFungibleToken{}, "token/"+MsgTypeCreateFungibleToken, nil)
	cdc.RegisterConcrete(MsgSetFungibleTokenStatus{}, "token/"+MsgTypeSetFungibleTokenStatus, nil)
	cdc.RegisterConcrete(MsgTransferFungibleToken{}, "token/"+MsgTypeTransferFungibleToken, nil)
	cdc.RegisterConcrete(MsgMintFungibleToken{}, "token/"+MsgTypeMintFungibleToken, nil)
	cdc.RegisterConcrete(MsgBurnFungibleToken{}, "token/"+MsgTypeBurnFungibleToken, nil)
	cdc.RegisterConcrete(MsgTransferFungibleTokenOwnership{}, "token/"+MsgTypeTransferFungibleTokenOwnership, nil)
	cdc.RegisterConcrete(MsgAcceptFungibleTokenOwnership{}, "token/"+MsgTypeAcceptFungibleTokenOwnership, nil)
	cdc.RegisterConcrete(MsgSetFungibleTokenAccountStatus{}, "token/"+MsgTypeSetFungibleTokenAccountStatus, nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
	codec.RegisterCrypto(msgCdc)
}
