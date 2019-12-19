package maintenance

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

type TokenMaintainer struct {
	Action              string                `json:"action"`
	AuthorisedAddresses []sdkTypes.AccAddress `json:"authorisedAddresses"`
	IssuerAddresses     []sdkTypes.AccAddress `json:"issuerAddresses"`
	ProviderAddresses   []sdkTypes.AccAddress `json:"providerAddresses"`
}

func NewTokenMaintainer(action string, authorisedAddresses, issuerAddresses, providerAddresses []sdkTypes.AccAddress) TokenMaintainer {
	return TokenMaintainer{
		Action:              action,
		AuthorisedAddresses: authorisedAddresses,
		IssuerAddresses:     issuerAddresses,
		ProviderAddresses:   providerAddresses,
	}
}

var _ MsgProposalData = &TokenMaintainer{}

func (tokenMaintainer TokenMaintainer) GetType() ProposalKind {
	return ProposalTypeModifyToken
}

func (tokenMaintainer *TokenMaintainer) Unmarshal(data []byte) error {
	err := msgCdc.UnmarshalBinaryLengthPrefixed(data, tokenMaintainer)
	if err != nil {
		return err
	}
	return nil
}

func (tokenMaintainer TokenMaintainer) Marshal() ([]byte, error) {
	bz, err := msgCdc.MarshalBinaryLengthPrefixed(tokenMaintainer)
	if err != nil {
		return nil, err
	}
	return bz, nil
}
