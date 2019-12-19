package maintenance

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdkTypes.CodespaceType = "Maintenance"

	CodeUnknownProposal sdkTypes.CodeType = iota
	CodeInactiveProposal
	CodeInvalidTitle
	CodeInvalidDescription
	CodeInvalidProposalType
	CodeInvalidProposalStatus
	CodeInvalidGenesis
)

// Error constructors

func ErrUnknownProposal(codespace sdkTypes.CodespaceType, proposalID uint64) sdkTypes.Error {
	return sdkTypes.NewError(codespace, CodeUnknownProposal, fmt.Sprintf("Unknown proposal with id %d", proposalID))
}

func ErrInactiveProposal(codespace sdkTypes.CodespaceType, proposalID uint64) sdkTypes.Error {
	return sdkTypes.NewError(codespace, CodeInactiveProposal, fmt.Sprintf("Inactive proposal with id %d", proposalID))
}

func ErrInvalidTitle(codespace sdkTypes.CodespaceType, errorMsg string) sdkTypes.Error {
	return sdkTypes.NewError(codespace, CodeInvalidTitle, errorMsg)
}

func ErrInvalidDescription(codespace sdkTypes.CodespaceType, errorMsg string) sdkTypes.Error {
	return sdkTypes.NewError(codespace, CodeInvalidDescription, errorMsg)
}

func ErrInvalidProposalType(codespace sdkTypes.CodespaceType, proposalType ProposalKind) sdkTypes.Error {
	return sdkTypes.NewError(codespace, CodeInvalidProposalType, fmt.Sprintf("Proposal Type '%s' is not valid", proposalType))
}

func ErrInvalidProposalData(codespace sdkTypes.CodespaceType, proposalData  MsgProposalData) sdkTypes.Error {
	return sdkTypes.NewError(codespace, CodeInvalidProposalType, fmt.Sprintf("Proposal data '%s' is not valid", proposalData))
}

func ErrInvalidGenesis(codespace sdkTypes.CodespaceType, msg string) sdkTypes.Error {
	return sdkTypes.NewError(codespace, CodeInvalidGenesis, msg)
}
