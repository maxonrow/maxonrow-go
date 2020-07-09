package nonfungible

import "github.com/cosmos/cosmos-sdk/codec"

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateNonFungibleToken{}, "nonFungible/"+MsgTypeCreateNonFungibleToken, nil)
	cdc.RegisterConcrete(MsgSetNonFungibleTokenStatus{}, "nonFungible/"+MsgTypeSetNonFungibleTokenStatus, nil)
	cdc.RegisterConcrete(MsgTransferNonFungibleItem{}, "nonFungible/"+MsgTypeTransferNonFungibleItem, nil)
	cdc.RegisterConcrete(MsgMintNonFungibleItem{}, "nonFungible/"+MsgTypeMintNonFungibleItem, nil)
	cdc.RegisterConcrete(MsgBurnNonFungibleItem{}, "nonFungible/"+MsgTypeBurnNonFungibleItem, nil)
	cdc.RegisterConcrete(MsgTransferNonFungibleTokenOwnership{}, "nonFungible/"+MsgTypeTransferNonFungibleTokenOwnership, nil)
	cdc.RegisterConcrete(MsgAcceptNonFungibleTokenOwnership{}, "nonFungible/"+MsgTypeAcceptNonFungibleTokenOwnership, nil)
	cdc.RegisterConcrete(MsgSetNonFungibleItemStatus{}, "nonFungible/"+MsgTypeSetNonFungibleItemStatus, nil)
	cdc.RegisterConcrete(MsgEndorsement{}, "nonFungible/"+MsgTypeEndorsement, nil)
	cdc.RegisterConcrete(MsgUpdateItemMetadata{}, "nonFungible/"+MsgTypeUpdateItemMetadata, nil)
	cdc.RegisterConcrete(MsgUpdateNFTMetadata{}, "nonFungible/"+MsgTypeUpdateNFTMetadata, nil)
	cdc.RegisterConcrete(MsgUpdateEndorserList{}, "nonFungible/"+MsgTypeUpdateEndorserList, nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
	codec.RegisterCrypto(msgCdc)
}
