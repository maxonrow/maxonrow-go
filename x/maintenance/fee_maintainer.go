package maintenance

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type FeeMaintainer struct {
	Action              string                `json:"action"`
	FeeCollectors       []FeeCollector        `json:"feeCollectors"`
	AuthorisedAddresses []sdkTypes.AccAddress `json:"authorisedAddresses"`
}

func NewFeeeMaintainer(action string, authorisedAddresses []sdkTypes.AccAddress, feeCollectors []FeeCollector) FeeMaintainer {
	return FeeMaintainer{
		Action:              action,
		FeeCollectors:       feeCollectors,
		AuthorisedAddresses: authorisedAddresses,
	}
}

type FeeCollector struct {
	Module  string              `json:"module"`
	Address sdkTypes.AccAddress `json:"address"`
}

var _ MsgProposalData = &FeeMaintainer{}

func (feeMaintainer FeeMaintainer) GetType() ProposalKind {
	return ProposalTypeModifyFee
}

func (feeMaintainer *FeeMaintainer) Unmarshal(data []byte) error {
	err := msgCdc.UnmarshalBinaryLengthPrefixed(data, feeMaintainer)
	if err != nil {
		return err
	}
	return nil
}

func (feeMaintainer FeeMaintainer) Marshal() ([]byte, error) {
	bz, err := msgCdc.MarshalBinaryLengthPrefixed(feeMaintainer)
	if err != nil {
		return nil, err
	}
	return bz, nil
}
