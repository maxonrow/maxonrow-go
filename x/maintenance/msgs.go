package maintenance

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

// Governance message types and routes
const (
	TypeMsgCastAction     = "castAction"
	TypeMsgSubmitProposal = "submitProposal"
	RouterKey             = "maintenance"

	MaxDescriptionLength int = 5000
	MaxTitleLength       int = 140
)

//nolint
type ProposalKind byte

//nolint
const (
	ProposalTypeModifyFee           ProposalKind = 0x01
	ProposalTypeModifyToken         ProposalKind = 0x02
	ProposalTypeModifyNameservice   ProposalKind = 0x03
	ProposalTypeModifyKyc           ProposalKind = 0x04
	ProposalTypesModifyValidatorSet ProposalKind = 0x05
	ProposalTypeModifyNonFungible   ProposalKind = 0x06
)

// MsgSubmitProposal
type MsgProposal struct {
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	ProposalType ProposalKind        `json:"proposalType"`
	ProposalData MsgProposalData     `json:"proposalData"`
	Proposer     sdkTypes.AccAddress `json:"proposer"`
}

type MsgProposalData interface {
	GetType() ProposalKind
}

func NewMsgSubmitProposal(title, description string, proposalType ProposalKind, proposalData MsgProposalData, proposer sdkTypes.AccAddress) MsgProposal {
	return MsgProposal{
		Title:        title,
		Description:  description,
		ProposalType: proposalType,
		ProposalData: proposalData,
		Proposer:     proposer,
	}
}

//nolint
func (msg MsgProposal) Route() string { return RouterKey }
func (msg MsgProposal) Type() string  { return TypeMsgSubmitProposal }

// Implements Msg.
func (msg MsgProposal) ValidateBasic() sdkTypes.Error {
	if len(msg.Title) == 0 {
		return ErrInvalidTitle(DefaultCodespace, "No title present in proposal")
	}
	if len(msg.Title) > MaxTitleLength {
		return ErrInvalidTitle(DefaultCodespace, fmt.Sprintf("Proposal title is longer than max length of %d", MaxTitleLength))
	}
	if len(msg.Description) == 0 {
		return ErrInvalidDescription(DefaultCodespace, "No description present in proposal")
	}
	if len(msg.Description) > MaxDescriptionLength {
		return ErrInvalidDescription(DefaultCodespace, fmt.Sprintf("Proposal description is longer than max length of %d", MaxDescriptionLength))
	}
	if !validProposalType(msg.ProposalType) {
		return ErrInvalidProposalType(DefaultCodespace, msg.ProposalType)
	}
	if msg.ProposalData == nil {
		return ErrInvalidProposalData(DefaultCodespace, msg.ProposalData)
	}
	if msg.Proposer.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Proposer.String())
	}
	return nil
}

func (msg MsgProposal) String() string {
	return fmt.Sprintf("MsgSubmitProposal{%s, %s, %s, %v}", msg.Title, msg.Description, msg.ProposalType, msg.ProposalData)
}

// Implements Msg.
func (msg MsgProposal) GetSignBytes() []byte {
	bz := msgCdc.MustMarshalJSON(msg)
	return sdkTypes.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgProposal) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Proposer}
}

// MsgCastAction
type MsgCastAction struct {
	Action     string              `json:"action"`
	ProposalID uint64              `json:"proposalId"` // ID of the proposal
	Owner      sdkTypes.AccAddress `json:"owner"`      //  address of the voter
}

func NewMsgCastAction(owner sdkTypes.AccAddress, proposalID uint64, action string) MsgCastAction {
	return MsgCastAction{
		ProposalID: proposalID,
		Owner:      owner,
		Action:     action,
	}
}

// Implements Msg.
// nolint
func (msg MsgCastAction) Route() string { return RouterKey }
func (msg MsgCastAction) Type() string  { return TypeMsgCastAction }

// Implements Msg.
func (msg MsgCastAction) ValidateBasic() sdkTypes.Error {
	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}
	if msg.ProposalID < 0 {
		return ErrUnknownProposal(DefaultCodespace, msg.ProposalID)
	}

	if msg.Action == "" {
		return sdkTypes.ErrInternal("Invalid action.")
	}
	return nil
}

// Implements Msg.
func (msg MsgCastAction) GetSignBytes() []byte {
	bz := msgCdc.MustMarshalJSON(msg)
	return sdkTypes.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgCastAction) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
}
