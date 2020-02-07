package nonfungible

import "github.com/cosmos/cosmos-sdk/codec"

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateNonFungibleToken{}, "nonFungible/"+MsgTypeCreateNonFungibleToken, nil)
	cdc.RegisterConcrete(MsgSetNonFungibleTokenStatus{}, "nonFungible/"+MsgTypeSetNonFungibleTokenStatus, nil)
	cdc.RegisterConcrete(MsgTransferNonFungibleToken{}, "nonFungible/"+MsgTypeTransferNonFungibleToken, nil)
	cdc.RegisterConcrete(MsgMintNonFungibleToken{}, "nonFungible/"+MsgTypeMintNonFungibleToken, nil)
	cdc.RegisterConcrete(MsgBurnNonFungibleToken{}, "nonFungible/"+MsgTypeBurnNonFungibleToken, nil)
	cdc.RegisterConcrete(MsgTransferNonFungibleTokenOwnership{}, "nonFungible/"+MsgTypeTransferNonFungibleTokenOwnership, nil)
	cdc.RegisterConcrete(MsgAcceptNonFungibleTokenOwnership{}, "nonFungible/"+MsgTypeAcceptNonFungibleTokenOwnership, nil)
	cdc.RegisterConcrete(MsgSetNonFungibleItemStatus{}, "nonFungible/"+MsgTypeSetNonFungibleItemStatus, nil)
	cdc.RegisterConcrete(MsgEndorsement{}, "nonFungible/"+MsgTypeEndorsement, nil)
	cdc.RegisterConcrete(MsgUpdateItemMetadata{}, "nonFungible/"+MsgTypeUpdateItemMetadata, nil)
	cdc.RegisterConcrete(MsgUpdateNFTMetadata{}, "nonFungible/"+MsgTypeUpdateNFTMetadata, nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
	codec.RegisterCrypto(msgCdc)
}
