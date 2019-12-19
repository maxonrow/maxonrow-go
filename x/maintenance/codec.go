package maintenance

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgProposal{}, "maintenance/msgProposal", nil)
	cdc.RegisterConcrete(MsgCastAction{}, "maintenance/msgCastAction", nil)

	cdc.RegisterInterface((*ProposalContent)(nil), nil)
	cdc.RegisterConcrete(TextProposal{}, "maintenance/proposal", nil)

	cdc.RegisterInterface((*MsgProposalData)(nil), nil)
	cdc.RegisterConcrete(FeeMaintainer{}, "maintenance/data/feeMaintainer", nil)
	cdc.RegisterConcrete(KycMaintainer{}, "maintenance/data/kycMaintainer", nil)
	cdc.RegisterConcrete(NameserviceMaintainer{}, "maintenance/data/nameserviceMaintainer", nil)
	cdc.RegisterConcrete(TokenMaintainer{}, "maintenance/data/tokenMaintainer", nil)
	cdc.RegisterConcrete(WhitelistValidator{}, "maintenance/data/whitelistValidator", nil)
	cdc.RegisterConcrete(NonFungibleMaintainer{}, "maintenance/data/nonFungibleMaintainer", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
	codec.RegisterCrypto(msgCdc)
}
