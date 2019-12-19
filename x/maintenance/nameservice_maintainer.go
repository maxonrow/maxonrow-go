package maintenance

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

type NameserviceMaintainer struct {
	Action              string                `json:"action"`
	AuthorisedAddresses []sdkTypes.AccAddress `json:"authorisedAddresses"`
	IssuerAddresses     []sdkTypes.AccAddress `json:"issuerAddresses"`
	ProviderAddresses   []sdkTypes.AccAddress `json:"providerAddresses"`
}

func NewNamerserviceMaintainer(action string, authorisedAddresses, issuerAddresses, providerAddresses []sdkTypes.AccAddress) NameserviceMaintainer {
	return NameserviceMaintainer{
		Action:              action,
		AuthorisedAddresses: authorisedAddresses,
		IssuerAddresses:     issuerAddresses,
		ProviderAddresses:   providerAddresses,
	}
}

var _ MsgProposalData = &NameserviceMaintainer{}

func (nameserviceMaintainer NameserviceMaintainer) GetType() ProposalKind {
	return ProposalTypeModifyNameservice
}

func (nameserviceMaintainer *NameserviceMaintainer) Unmarshal(data []byte) error {
	err := msgCdc.UnmarshalBinaryLengthPrefixed(data, nameserviceMaintainer)
	if err != nil {
		return err
	}
	return nil
}

func (nameserviceMaintainer NameserviceMaintainer) Marshal() ([]byte, error) {
	bz, err := msgCdc.MarshalBinaryLengthPrefixed(nameserviceMaintainer)
	if err != nil {
		return nil, err
	}
	return bz, nil
}
