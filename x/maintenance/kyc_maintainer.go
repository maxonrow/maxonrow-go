package maintenance

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

type KycMaintainer struct {
	Action              string                `json:"action"`
	AuthorisedAddresses []sdkTypes.AccAddress `json:"authorisedAddresses"`
	IssuerAddresses     []sdkTypes.AccAddress `json:"issuerAddresses"`
	ProviderAddresses   []sdkTypes.AccAddress `json:"providerAddresses"`
}

func NewKycMaintainer(action string, authorisedAddresses, issuerAddresses, providerAddresses []sdkTypes.AccAddress) KycMaintainer {
	return KycMaintainer{
		Action:              action,
		AuthorisedAddresses: authorisedAddresses,
		IssuerAddresses:     issuerAddresses,
		ProviderAddresses:   providerAddresses,
	}
}

var _ MsgProposalData = &KycMaintainer{}

func (kycMaintainer KycMaintainer) GetType() ProposalKind {
	return ProposalTypeModifyKyc
}

func (kycMaintainer *KycMaintainer) Unmarshal(data []byte) error {
	err := msgCdc.UnmarshalBinaryLengthPrefixed(data, kycMaintainer)
	if err != nil {
		return err
	}
	return nil
}

func (kycMaintainer KycMaintainer) Marshal() ([]byte, error) {
	bz, err := msgCdc.MarshalBinaryLengthPrefixed(kycMaintainer)
	if err != nil {
		return nil, err
	}
	return bz, nil
}
