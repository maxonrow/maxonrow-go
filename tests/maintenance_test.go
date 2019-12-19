package tests

import (
	"testing"

	"github.com/maxonrow/maxonrow-go/x/maintenance"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

type MaintenanceInfo struct {
	Action            string
	Title             string
	Description       string
	ProposalType      string
	AuthorisedAddress string
	IssuerAddress     string
	ProviderAddress   string
	FeeCollector      FeeCollector
	Proposer          string
	ValidatorPubKey   string
}

type CastAction struct {
	Caster     string
	Action     string
	ProposalId uint64
}

type FeeCollector struct {
	Module  string
	Address string
}

func makeMaintenanceMsg(t *testing.T, action, title, description, proposalType, authorisedAddress, issuerAddress, providerAddress, proposer, validatorPubKey string, feeCollector FeeCollector) maintenance.MsgProposal {

	proposalKind, proposalKindErr := maintenance.ProposalTypeFromString(proposalType)
	if proposalKindErr != nil {
		var msg1 maintenance.MsgProposal
		var authorisedAddr1 sdkTypes.AccAddress
		var issuerAddr1 sdkTypes.AccAddress
		var providerAddr1 sdkTypes.AccAddress
		if authorisedAddress != "" {
			authorisedAddr1 = tKeys[authorisedAddress].addr

		} else {
			authorisedAddr1 = nil
		}

		if issuerAddress != "" {
			issuerAddr1 = tKeys[issuerAddress].addr

		} else {
			issuerAddr1 = nil
		}

		if providerAddress != "" {
			providerAddr1 = tKeys[providerAddress].addr

		} else {
			providerAddr1 = nil
		}

		_maintainer := maintenance.NewNamerserviceMaintainer(action, []sdkTypes.AccAddress{authorisedAddr1}, []sdkTypes.AccAddress{issuerAddr1}, []sdkTypes.AccAddress{providerAddr1})
		msg1 = maintenance.NewMsgSubmitProposal(title, description, proposalKind, _maintainer, tKeys[proposer].addr)
		return msg1
	}
	require.NoError(t, proposalKindErr)

	proposerAddr := tKeys[proposer].addr

	var msg maintenance.MsgProposal
	// TO-DO: better implementation
	switch proposalKind {
	case maintenance.ProposalTypeModifyFee:
		var authorisedAddr sdkTypes.AccAddress
		var feeCollectorAddr sdkTypes.AccAddress
		if authorisedAddress != "" {
			authorisedAddr = tKeys[authorisedAddress].addr

		} else {
			authorisedAddr = nil
		}

		if feeCollector.Address != "" {
			feeCollectorAddr = tKeys[feeCollector.Address].addr
		}

		feeCollector := maintenance.FeeCollector{
			Module:  feeCollector.Module,
			Address: feeCollectorAddr,
		}
		feeMaintainer := maintenance.NewFeeeMaintainer(action, []sdkTypes.AccAddress{authorisedAddr}, []maintenance.FeeCollector{feeCollector})

		msg = maintenance.NewMsgSubmitProposal(title, description, proposalKind, feeMaintainer, proposerAddr)

	case maintenance.ProposalTypeModifyKyc:
		var authorisedAddr sdkTypes.AccAddress
		var issuerAddr sdkTypes.AccAddress
		var providerAddr sdkTypes.AccAddress
		if authorisedAddress != "" {
			authorisedAddr = tKeys[authorisedAddress].addr
		}

		if providerAddress != "" {
			providerAddr = tKeys[providerAddress].addr
		}

		if issuerAddress != "" {
			issuerAddr = tKeys[issuerAddress].addr
		}
		kycMaintainer := maintenance.NewKycMaintainer(action, []sdkTypes.AccAddress{authorisedAddr}, []sdkTypes.AccAddress{issuerAddr}, []sdkTypes.AccAddress{providerAddr})
		msg = maintenance.NewMsgSubmitProposal(title, description, proposalKind, kycMaintainer, proposerAddr)

	case maintenance.ProposalTypeModifyNameservice:
		var authorisedAddr sdkTypes.AccAddress
		var issuerAddr sdkTypes.AccAddress
		var providerAddr sdkTypes.AccAddress
		if authorisedAddress != "" {
			authorisedAddr = tKeys[authorisedAddress].addr
		}

		if providerAddress != "" {
			providerAddr = tKeys[providerAddress].addr
		}

		if issuerAddress != "" {
			issuerAddr = tKeys[issuerAddress].addr
		}
		nameserviceMaintainer := maintenance.NewNamerserviceMaintainer(action, []sdkTypes.AccAddress{authorisedAddr}, []sdkTypes.AccAddress{issuerAddr}, []sdkTypes.AccAddress{providerAddr})
		msg = maintenance.NewMsgSubmitProposal(title, description, proposalKind, nameserviceMaintainer, proposerAddr)

	case maintenance.ProposalTypeModifyToken:
		var authorisedAddr sdkTypes.AccAddress
		var issuerAddr sdkTypes.AccAddress
		var providerAddr sdkTypes.AccAddress
		if authorisedAddress != "" {
			authorisedAddr = tKeys[authorisedAddress].addr
		}

		if providerAddress != "" {
			providerAddr = tKeys[providerAddress].addr
		}

		if issuerAddress != "" {
			issuerAddr = tKeys[issuerAddress].addr
		}
		tokenMaintainer := maintenance.NewTokenMaintainer(action, []sdkTypes.AccAddress{authorisedAddr}, []sdkTypes.AccAddress{issuerAddr}, []sdkTypes.AccAddress{providerAddr})
		msg = maintenance.NewMsgSubmitProposal(title, description, proposalKind, tokenMaintainer, proposerAddr)
	case maintenance.ProposalTypeModifyNonFungible:
		var authorisedAddr sdkTypes.AccAddress
		var issuerAddr sdkTypes.AccAddress
		var providerAddr sdkTypes.AccAddress
		if authorisedAddress != "" {
			authorisedAddr = tKeys[authorisedAddress].addr
		}

		if providerAddress != "" {
			providerAddr = tKeys[providerAddress].addr
		}

		if issuerAddress != "" {
			issuerAddr = tKeys[issuerAddress].addr
		}
		nonFungibleMaintainer := maintenance.NewNonFungibleMaintainer(action, []sdkTypes.AccAddress{authorisedAddr}, []sdkTypes.AccAddress{issuerAddr}, []sdkTypes.AccAddress{providerAddr})
		msg = maintenance.NewMsgSubmitProposal(title, description, proposalKind, nonFungibleMaintainer, proposerAddr)

	case maintenance.ProposalTypesModifyValidatorSet:
		pubKeyAddr, pubKeyErr := sdkTypes.GetConsPubKeyBech32(validatorPubKey)
		require.NoError(t, pubKeyErr)
		whitelistValidator := maintenance.NewWhitelistValidator(action, []crypto.PubKey{pubKeyAddr})
		msg = maintenance.NewMsgSubmitProposal(title, description, proposalKind, whitelistValidator, proposerAddr)

	}
	return msg

}

func makeCastActionMsg(t *testing.T, action, caster string, proposalId uint64) maintenance.MsgCastAction {

	casterAddr := tKeys[caster].addr
	msg := maintenance.NewMsgCastAction(casterAddr, proposalId, action)

	return msg

}
