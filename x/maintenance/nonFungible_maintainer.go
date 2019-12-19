package maintenance

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

type NonFungibleMaintainer struct {
	Action              string                `json:"action"`
	AuthorisedAddresses []sdkTypes.AccAddress `json:"authorisedAddresses"`
	IssuerAddresses     []sdkTypes.AccAddress `json:"issuerAddresses"`
	ProviderAddresses   []sdkTypes.AccAddress `json:"providerAddresses"`
}

func NewNonFungibleMaintainer(action string, authorisedAddresses, issuerAddresses, providerAddresses []sdkTypes.AccAddress) NonFungibleMaintainer {
	return NonFungibleMaintainer{
		Action:              action,
		AuthorisedAddresses: authorisedAddresses,
		IssuerAddresses:     issuerAddresses,
		ProviderAddresses:   providerAddresses,
	}
}

var _ MsgProposalData = &NonFungibleMaintainer{}

func (nonFungibleMaintainer NonFungibleMaintainer) GetType() ProposalKind {
	return ProposalTypeModifyNonFungible
}

func (nonFungibleMaintainer *NonFungibleMaintainer) Unmarshal(data []byte) error {
	err := msgCdc.UnmarshalBinaryLengthPrefixed(data, nonFungibleMaintainer)
	if err != nil {
		return err
	}
	return nil
}

func (nonFungibleMaintainer NonFungibleMaintainer) Marshal() ([]byte, error) {
	bz, err := msgCdc.MarshalBinaryLengthPrefixed(nonFungibleMaintainer)
	if err != nil {
		return nil, err
	}
	return bz, nil
}
