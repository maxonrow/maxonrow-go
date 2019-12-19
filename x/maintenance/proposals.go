package maintenance

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/maxonrow/maxonrow-go/types"
)

// Proposal is a struct used by gov module internally
// embedds ProposalContent with additional fields to record the status of the proposal process
type Proposal struct {
	ProposalContent `json:"proposal_content"` // Proposal content interface
	ProposalID      uint64                    `json:"proposalId"`     //  ID of the proposal
	ProposalData    MsgProposalData           `json:"proposalData"`   // Proposal content
	Status          ProposalStatus            `json:"proposalStatus"` //  Status of the Proposal {Active, Completed}
	Approvers       types.AddressHolder       `json:"approvers"`      // Approvers of the proposal
	Rejecters       types.AddressHolder       `json:"rejecters"`
	SubmitTime      time.Time                 `json:"submitTime"` //  Time of the block where TxGovSubmitProposal was included
}

// nolint
func (p Proposal) String() string {
	return fmt.Sprintf(`Proposal %d:
  Title:              %s
  Type:               %s
  Status:             %s
  Submit Time:        %s
  Description:        %s`,
		p.ProposalID, p.GetTitle(), p.ProposalType(),
		p.Status, p.SubmitTime, p.GetDescription(),
	)
}

// ProposalContent is an interface that has title, description, and proposaltype
// that the governance module can use to identify them and generate human readable messages
// ProposalContent can have additional fields, which will handled by ProposalHandlers
// via type assertion, e.g. parameter change amount in ParameterChangeProposal
type ProposalContent interface {
	GetTitle() string
	GetDescription() string
	ProposalType() ProposalKind
}

// Text Proposals
type TextProposal struct {
	Title       string       `json:"title"`       //  Title of the proposal
	Description string       `json:"description"` //  Description of the proposal
	Type        ProposalKind `json:"proposalType"`
}

func NewTextProposal(title, description string, proposalType ProposalKind) TextProposal {
	return TextProposal{
		Title:       title,
		Description: description,
		Type:        proposalType,
	}
}

// Implements Proposal Interface
var _ ProposalContent = TextProposal{}

// nolint
func (tp TextProposal) GetTitle() string           { return tp.Title }
func (tp TextProposal) GetDescription() string     { return tp.Description }
func (tp TextProposal) ProposalType() ProposalKind { return tp.Type }

// String to proposalType byte. Returns 0xff if invalid.
func ProposalTypeFromString(str string) (ProposalKind, error) {
	switch str {
	case "token":
		return ProposalTypeModifyToken, nil
	case "kyc":
		return ProposalTypeModifyKyc, nil
	case "nameservice":
		return ProposalTypeModifyNameservice, nil
	case "fee":
		return ProposalTypeModifyFee, nil
	case "validator":
		return ProposalTypesModifyValidatorSet, nil
	case "nonFungible":
		return ProposalTypeModifyNonFungible, nil
	default:
		return ProposalKind(0xff), fmt.Errorf("'%s' is not a valid proposal type", str)
	}
}

// is defined ProposalType?
func validProposalType(pt ProposalKind) bool {
	if pt == ProposalTypeModifyFee ||
		pt == ProposalTypeModifyKyc ||
		pt == ProposalTypeModifyNameservice ||
		pt == ProposalTypeModifyToken ||
		pt == ProposalTypesModifyValidatorSet ||
		pt == ProposalTypeModifyNonFungible {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (pt ProposalKind) Marshal() ([]byte, error) {
	return []byte{byte(pt)}, nil
}

// Unmarshal needed for protobuf compatibility
func (pt *ProposalKind) Unmarshal(data []byte) error {
	*pt = ProposalKind(data[0])
	return nil
}

// Marshals to JSON using string
func (pt ProposalKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(pt.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (pt *ProposalKind) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := ProposalTypeFromString(s)
	if err != nil {
		return err
	}
	*pt = bz2
	return nil
}

// Turns VoteOption byte to String
func (pt ProposalKind) String() string {
	switch pt {
	case ProposalTypeModifyFee:
		return "ModifyFeeMaintainer"
	case ProposalTypeModifyKyc:
		return "ModifyKycMaintainer"
	case ProposalTypeModifyNameservice:
		return "ModifyNameserviceMaintainer"
	case ProposalTypeModifyToken:
		return "ModifyTokenMaintainer"
	case ProposalTypeModifyNonFungible:
		return "ModifyNonFungibleMaintainer"
	default:
		return ""
	}

}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (pt ProposalKind) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(pt.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(pt))))
	}
}

// ProposalStatus

// Type that represents Proposal Status as a byte
type ProposalStatus byte

//nolint
const (
	StatusActive    ProposalStatus = 0x00
	StatusCompleted ProposalStatus = 0x01
	StatusRejected  ProposalStatus = 0x02
)

// ProposalStatusToString turns a string into a ProposalStatus
func ProposalStatusFromString(str string) (ProposalStatus, error) {
	switch str {
	case "Active":
		return StatusActive, nil
	case "Completed":
		return StatusCompleted, nil
	case "Rejeced":
		return StatusRejected, nil
	default:
		return ProposalStatus(0xff), fmt.Errorf("'%s' is not a valid proposal status", str)
	}
}

// is defined ProposalType?
func validProposalStatus(status ProposalStatus) bool {
	if status == StatusActive ||
		status == StatusCompleted {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (status ProposalStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *ProposalStatus) Unmarshal(data []byte) error {
	*status = ProposalStatus(data[0])
	return nil
}

// Marshals to JSON using string
func (status ProposalStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (status *ProposalStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := ProposalStatusFromString(s)
	if err != nil {
		return err
	}
	*status = bz2
	return nil
}

// Turns VoteStatus byte to String
func (status ProposalStatus) String() string {
	switch status {
	case StatusActive:
		return "Active"
	case StatusCompleted:
		return "Completed"
	case StatusRejected:
		return "Rejected"
	default:
		return ""
	}
}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (status ProposalStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}
