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

func makeMaintenaceTxs() []*testCase {

	var proposalTitleOutOfLength string
	proposalTitleOutOfLength = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333"

	var proposalDescOutOfLength string
	proposalDescOutOfLength = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333----aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzzaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccc--aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333----aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzzaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccc--aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---" +
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---" +
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaa-eeeeee"

	tcs := []*testCase{

		//=============================================start :tx modules
		// Add Proposal [ challenge ]
		{"maintenance", true, true, "Is invalid maintainer: Add fee provider address.", "nago", "0cin", 0, MaintenanceInfo{"add", "Add fee-provider address", "Add a party as fee-provider address", "fee", "nago", "", "", FeeCollector{}, "nago", ""}, "", nil},
		{"maintenance", true, true, "Zero-length of Proposal Title: Add fee provider address.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "", "Add a party as fee-provider address", "fee", "nago", "", "", FeeCollector{}, "maintainer-1", ""}, "", nil},
		{"maintenance", true, true, "Proposal Title was out-of-length: Add fee provider address.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", proposalTitleOutOfLength, "Add a party as fee-provider address", "fee", "maintainer-1", "", "", FeeCollector{}, "nago", ""}, "", nil},
		{"maintenance", true, true, "Zero-length of Proposal Description: Add fee provider address.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add fee-provider address", "", "fee", "nago", "", "", FeeCollector{}, "maintainer-1", ""}, "", nil},
		{"maintenance", true, true, "Proposal Description was out-of-length: Add fee provider address.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add fee-provider address", proposalDescOutOfLength, "fee", "nago", "", "", FeeCollector{}, "nago", ""}, "", nil},
		{"maintenance", true, true, "Signer and proposer is different", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add fee-provider address", "Add a party as fee-provider address", "fee", "nago", "", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", true, true, "Invalid proposal type: Add nameservice provider address.", "nago", "0cin", 0, MaintenanceInfo{"add", "Add nameservice-provider address", "Add a party as nameservice-provider address", "nameservice-999", "nago", "", "maintainer-3", FeeCollector{}, "maintainer-3", ""}, "", nil},

		//------------------------------
		// Maintenance - Fee
		// Maintenance - Fee - Add(Proposal) [proposal-1, proposal-2, proposal-3, proposal-4]
		{"maintenance", false, false, "1. Proposal, add fee authorised address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add authorised address", "Add cmo as fee authorised address", "fee", "maintainer-1", "", "", FeeCollector{}, "maintainer-1", ""}, "", nil},
		{"maintenance", false, false, "2. Proposal, add fee issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add isser address", "Add fee issuer address", "fee", "", "maintainer-1", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "3. Proposal, add fee provider address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"add", "Add isser address", "Add fee issuer address", "fee", "", "", "maintainer-2", FeeCollector{}, "maintainer-3", ""}, "", nil},
		{"maintenance", false, false, "4. Proposal, add fee collector address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add token fee collector", "Add maintainer-1 as token fee collector", "fee", "", "", "", FeeCollector{Module: "token", Address: "maintainer-1"}, "maintainer-2", ""}, "", nil},

		// Maintenance - Fee - Cast Action [proposal-1, proposal-2, proposal-3, proposal-4]
		// Cast Action - Proposal 1
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 1, signer and caster different person.", "maintainer-1", "0cin", 0, CastAction{"maintainer-2", "approve", 1}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 1, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 1}, "", nil},
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 1, caster and signer not maintainer.", "yk", "0cin", 0, CastAction{"yk", "approve", 1}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 1, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 1}, "", nil},
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to inactive proposal 1.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 1}, "", nil},
		// Cast Action - Proposal 2
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 2, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 2}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 2, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 2}, "", nil},
		// Cast Action - Proposal 3
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 3, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 3}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 3, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 3}, "", nil},
		// Cast Action - Proposal 4
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 3, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 4}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 3, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 4}, "", nil},

		// Maintenance - Fee - Remove(Proposal) [proposal-5, proposal-6, proposal-7, proposal-8]
		{"maintenance", false, false, "5. Proposal, remove fee authorised address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"remove", "Remove authorised address", "Remove cmo as fee authorised address", "fee", "maintainer-1", "", "", FeeCollector{}, "maintainer-1", ""}, "", nil},
		{"maintenance", false, false, "6. Proposal, remove fee issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"remove", "Remove isser address", "Remove fee issuer address", "fee", "", "maintainer-1", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "7. Proposal, remove fee provider address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"remove", "Remove isser address", "Remove fee issuer address", "fee", "", "", "maintainer-2", FeeCollector{}, "maintainer-3", ""}, "", nil},
		{"maintenance", false, false, "8. Proposal, remove fee collector address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"remove", "Remove token fee collector", "Remove maintainer-1 as token fee collector", "fee", "", "", "", FeeCollector{Module: "token", Address: "maintainer-1"}, "maintainer-2", ""}, "", nil},

		// Maintenance - Fee - Cast Action [proposal-5, proposal-6, proposal-7, proposal-8]
		// Cast Action - Proposal 5
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 5, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 5}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 5, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 5}, "", nil},
		// Cast Action - Proposal 6
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 6, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 6}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 6, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 6}, "", nil},
		// Cast Action - Proposal 7
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 7, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 7}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 7, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 7}, "", nil},
		// Cast Action - Proposal 8
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 8, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 8}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 8, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 8}, "", nil},

		//---------------------------------
		// Maintenance - KYC
		// Maintenance - KYC [proposal-9, proposal-10, proposal-11]
		{"maintenance", false, false, "9. Proposal, add kyc authorised address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add KYC-authorised address", "Add a party as KYC-authorised address", "kyc", "maintainer-1", "", "", FeeCollector{}, "maintainer-1", ""}, "", nil},
		{"maintenance", false, false, "10. Proposal, add kyc issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add KYC-issuer address", "Add a party as KYC-issuer address", "kyc", "", "maintainer-2", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "11. Proposal, add kyc provider address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"add", "Add KYC-provider address", "Add a party as KYC-provider address", "kyc", "", "", "maintainer-3", FeeCollector{}, "maintainer-3", ""}, "", nil},

		// Cast Action - Proposal 9
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 9, signer and caster different person.", "maintainer-1", "0cin", 0, CastAction{"maintainer-2", "approve", 13}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 9, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 9}, "", nil},
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 9, caster and signer not maintainer.", "yk", "0cin", 0, CastAction{"yk", "approve", 9}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 9, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 9}, "", nil},
		//{"maintenance-cast-action", true, true, "(Approve)-Cast action to inactive proposal 9.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "disable", 9}, "", nil},//kiv
		// Cast Action - Proposal 10
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 10, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 10}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 10, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 10}, "", nil},
		// Cast Action - Proposal 11
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 11, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 11}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 11, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 11}, "", nil},

		// Maintenance - KYC - Remove(Proposal) [proposal-12, proposal-13, proposal-14]
		{"maintenance", false, false, "12. Proposal, remove kyc authorised address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"remove", "Remove KYC-authorised address", "Remove a party as KYC-authorised address", "kyc", "maintainer-1", "", "", FeeCollector{}, "maintainer-1", ""}, "", nil},
		{"maintenance", false, false, "13. Proposal, remove kyc issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"remove", "Remove KYC-issuer address", "Remove a party as KYC-issuer address", "kyc", "", "maintainer-2", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "14. Proposal, remove kyc provider address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"remove", "Remove KYC-provider address", "Remove a party as KYC-provider address", "kyc", "", "", "maintainer-3", FeeCollector{}, "maintainer-3", ""}, "", nil},

		// Cast Action for KYC - Proposal 12
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 12, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 12}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 12, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 12}, "", nil},
		// Cast Action for KYC - Proposal 13
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 13, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 13}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 13, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 13}, "", nil},
		// Cast Action for KYC - Proposal 14
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 14, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 14}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 14, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 14}, "", nil},
		{"maintenance-cast-action", false, false, "(Reject)-Cast action to proposal 14, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "reject", 14}, "", nil},

		//---------------------------------
		// Maintenance - Token
		// Maintenance - Token [proposal-15, proposal-16, proposal-17]
		{"maintenance", false, false, "15. Proposal, add token authorised address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"add", "Add token-authorised address", "Add a party as token-authorised address", "token", "maintainer-3", "", "", FeeCollector{}, "maintainer-3", ""}, "", nil},
		{"maintenance", false, false, "16. Proposal, add token issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add token-issuer address", "Add a party as token-issuer address", "token", "", "maintainer-2", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "17. Proposal, add token provider address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add token-provider address", "Add a party as token-provider address", "token", "", "", "maintainer-1", FeeCollector{}, "maintainer-1", ""}, "", nil},

		// Cast Action - Proposal 15
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 15, signer and caster different person.", "maintainer-1", "0cin", 0, CastAction{"maintainer-3", "approve", 19}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 15, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 15}, "", nil},
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 15, caster and signer not maintainer.", "yk", "0cin", 0, CastAction{"yk", "approve", 15}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 15, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 15}, "", nil},
		//{"maintenance-cast-action", true, true, "(Approve)-Cast action to inactive proposal 15.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "disable", 15}, "", nil},//kiv
		// Cast Action - Proposal 16
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 16, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 16}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 16, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 16}, "", nil},
		// Cast Action - Proposal 17
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 17, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 17}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 17, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 17}, "", nil},

		// Maintenance - Token - Remove(Proposal) [proposal-18, proposal-19, proposal-20]
		{"maintenance", false, false, "18. Proposal, remove token authorised address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"remove", "Remove token-authorised address", "Remove a party as token-authorised address", "token", "maintainer-1", "", "", FeeCollector{}, "maintainer-1", ""}, "", nil},
		{"maintenance", false, false, "19. Proposal, remove token issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"remove", "Remove token-issuer address", "Remove a party as token-issuer address", "token", "", "maintainer-2", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "20. Proposal, remove token provider address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"remove", "Remove token-provider address", "Remove a party as token-provider address", "token", "", "", "maintainer-3", FeeCollector{}, "maintainer-3", ""}, "", nil},

		// Cast Action - Proposal 18
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 18, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 18}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 18, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 18}, "", nil},
		{"maintenance-cast-action", false, false, "(Reject)-Cast action to proposal 18, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "reject", 18}, "", nil},
		// Cast Action - Proposal 19
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 19, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 19}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 19, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 19}, "", nil},
		{"maintenance-cast-action", false, false, "(Reject)-Cast action to proposal 19, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "reject", 19}, "", nil},
		// Cast Action - Proposal 20
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 20, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 20}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 20, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 20}, "", nil},

		//---------------------------------
		// Maintenance - Nameservice
		// Maintenance - Nameservice [proposal-21, proposal-22, proposal-23]
		{"maintenance", false, false, "21. Proposal, add nameservice authorised address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"add", "Add nameservice-authorised address", "Add a party as nameservice-authorised address", "nameservice", "maintainer-3", "", "", FeeCollector{}, "maintainer-3", ""}, "", nil},
		{"maintenance", false, false, "22. Proposal, add nameservice issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add nameservice-issuer address", "Add a party as nameservice-issuer address", "nameservice", "", "maintainer-2", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "23. Proposal, add nameservice provider address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add nameservice-provider address", "Add a party as nameservice-provider address", "nameservice", "", "", "maintainer-1", FeeCollector{}, "maintainer-1", ""}, "", nil},

		// Cast Action - Proposal 21
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 21, signer and caster different person.", "maintainer-1", "0cin", 0, CastAction{"maintainer-2", "approve", 21}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 21, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 21}, "", nil},
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 21, caster and signer not maintainer.", "yk", "0cin", 0, CastAction{"yk", "approve", 21}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 21, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 21}, "", nil},
		//{"maintenance-cast-action", true, true, "(Approve)-Cast action to inactive proposal 21.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "disable", 21}, "", nil},//kiv
		// Cast Action - Proposal 22
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 22, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 22}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 22, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 22}, "", nil},
		// Cast Action - Proposal 23
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 23, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 23}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 23, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 23}, "", nil},

		// Maintenance - Nameservice - Remove(Proposal) [proposal-24, proposal-25, proposal-26]
		{"maintenance", false, false, "24. Proposal, remove nameservice authorised address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"remove", "Remove nameservice-authorised address", "Remove a party as nameservice-authorised address", "nameservice", "maintainer-3", "", "", FeeCollector{}, "maintainer-3", ""}, "", nil},
		{"maintenance", false, false, "25. Proposal, remove nameservice issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"remove", "Remove nameservice-issuer address", "Remove a party as nameservice-issuer address", "nameservice", "", "maintainer-2", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "26. Proposal, remove nameservice provider address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"remove", "Remove nameservice-provider address", "Remove a party as nameservice-provider address", "nameservice", "", "", "maintainer-1", FeeCollector{}, "maintainer-1", ""}, "", nil},

		// Cast Action - Proposal 24
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 24, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 24}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 24, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 24}, "", nil},
		{"maintenance-cast-action", false, false, "(Reject)-Cast action to proposal 24, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "reject", 24}, "", nil},
		// Cast Action - Proposal 25
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 25, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 25}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 25, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 25}, "", nil},
		{"maintenance-cast-action", false, false, "(Reject)-Cast action to proposal 25, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "reject", 25}, "", nil},
		// Cast Action - Proposal 26
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 26, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 26}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 26, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 26}, "", nil},

		//---------------------------------
		// Maintenance - Validator-set
		// Maintenance - Validator-set [proposal-27, proposal-28, proposal-29]
		{"maintenance", false, false, "27. Proposal, add validator set-1, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add validator-set", "Add a party for validator-set", "validator", "", "", "", FeeCollector{}, "maintainer-1", "mxwvalconspub1zcjduepq2vxnxwuzvf82w9mxhjuwm35q7e84pfglsexh5l0ffqz0ddfxjp5q8wjkgw"}, "", nil},
		{"maintenance", false, false, "28. Proposal, add validator set-2, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add validator-set", "Add a party for validator-set", "validator", "", "", "", FeeCollector{}, "maintainer-2", "mxwvalconspub1zcjduepqczwdy9dlmvazg3u3nml743xgprr2e82n2lt6wue5ycsga2nudvxq0avuc6"}, "", nil},
		{"maintenance", false, false, "29. Proposal, add validator set-3, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"add", "Add validator-set", "Add a party for validator-set", "validator", "", "", "", FeeCollector{}, "maintainer-3", "mxwvalconspub1zcjduepqvf9vf3cdxwtk65ya83q8uz36c8vqn5gylp3dmkghxjs253thve4qqzm5ca"}, "", nil},

		// Cast Action - Proposal 27
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 27, signer and caster different person.", "maintainer-1", "0cin", 0, CastAction{"maintainer-2", "approve", 27}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 27, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 27}, "", nil},
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 27, caster and signer not maintainer.", "gohck", "0cin", 0, CastAction{"gohck", "approve", 27}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 27, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 27}, "", nil},
		// Cast Action - Proposal 28
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 28, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 28}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 28, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 28}, "", nil},
		// Cast Action - Proposal 29
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 29, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 29}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 29, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 29}, "", nil},

		// Maintenance - Validator-set - Remove(Proposal) [proposal-30, proposal-31, proposal-32]
		{"maintenance", false, false, "30. Proposal, remove validator set-1, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"remove", "Remove validator-set", "Remove a party for validator-set", "validator", "", "", "", FeeCollector{}, "maintainer-1", "mxwvalconspub1zcjduepq2vxnxwuzvf82w9mxhjuwm35q7e84pfglsexh5l0ffqz0ddfxjp5q8wjkgw"}, "", nil},
		{"maintenance", false, false, "31. Proposal, remove validator set-2, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"remove", "Remove validator-set", "Remove a party for validator-set", "validator", "", "", "", FeeCollector{}, "maintainer-2", "mxwvalconspub1zcjduepqczwdy9dlmvazg3u3nml743xgprr2e82n2lt6wue5ycsga2nudvxq0avuc6"}, "", nil},
		{"maintenance", false, false, "32. Proposal, remove validator set-3, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"remove", "Remove validator-set", "Remove a party for validator-set", "validator", "", "", "", FeeCollector{}, "maintainer-3", "mxwvalconspub1zcjduepqvf9vf3cdxwtk65ya83q8uz36c8vqn5gylp3dmkghxjs253thve4qqzm5ca"}, "", nil},

		// Cast Action - Proposal 30
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 30, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 30}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 30, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 30}, "", nil},
		{"maintenance-cast-action", false, false, "(Reject)-Cast action to proposal 30, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "reject", 30}, "", nil},
		// Cast Action - Proposal 31
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 31, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 31}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 31, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 31}, "", nil},
		{"maintenance-cast-action", false, false, "(Reject)-Cast action to proposal 31, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "reject", 31}, "", nil},
		// Cast Action - Proposal 32
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 32, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 32}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 32, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 32}, "", nil},

		// Cast Action - Proposal 33
		// add nameservice fee collector with maintenance.
		{"maintenance", false, false, "33. Proposal, add nameservice fee collector address, Happy path. commit", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add nameservice fee collector", "Add ns-feecollector as nameservice fee collector", "fee", "", "", "", FeeCollector{Module: "nameservice", Address: "ns-feecollector"}, "maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve ns-feecollector as nameservice fee collector, Happy path. commit", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 33}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve ns-feecollector as nameservice fee collector, Happy path. commit", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 33}, "", nil},

		// Cast Action - Proposal 34
		// add token fee collector with maintenance. (mostafa is whitelisted.)
		{"maintenance", false, false, "34. Proposal, add token fee collector address, Happy path. commit", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add token fee collector", "Add mostafa as token fee collector", "fee", "", "", "", FeeCollector{Module: "token", Address: "mostafa"}, "maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as ft fee collector, Happy path. commit", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 34}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as ft fee collector, Happy path. commit", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 34}, "", nil},

		//=============================================start : used by nft modules
		//add nft-mostafa as nonfungible authorised address
		{"maintenance", false, false, "35. Proposal, add nonfungible authorised address, Happy path. commit", "nft-maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add non fungible authorised address", "Add mostafa as non fungible authorised address", "nonFungible", "nft-mostafa", "", "", FeeCollector{}, "nft-maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as non fungible token authorised address, Happy path. commit", "nft-maintainer-1", "0cin", 0, CastAction{"nft-maintainer-1", "approve", 35}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as non fungible token authorised address, Happy path. commit", "nft-maintainer-3", "0cin", 0, CastAction{"nft-maintainer-3", "approve", 35}, "", nil},

		//add nft-carlo as nonfungible issuer address
		{"maintenance", false, false, "36. Proposal, add nonfungible issuer address, Happy path. commit", "nft-maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add non fungible fee issuer address", "Add carlo as non fungible issuer address", "nonFungible", "", "nft-carlo", "", FeeCollector{}, "nft-maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve carlo as non fungible token issuer address, Happy path. commit", "nft-maintainer-1", "0cin", 0, CastAction{"nft-maintainer-1", "approve", 36}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve carlo as non fungible token issuer address, Happy path. commit", "nft-maintainer-3", "0cin", 0, CastAction{"nft-maintainer-3", "approve", 36}, "", nil},

		//add nft-jeansoon as nonfungible provider address
		{"maintenance", false, false, "37. Proposal, add nonfungible provider address, Happy path. commit", "nft-maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add non fungible fee provider address", "Add jeansoon as non fungible provider address", "nonFungible", "", "", "nft-jeansoon", FeeCollector{}, "nft-maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "Cast action to approve jeansoon as non fungible token provider address, Happy path. commit", "nft-maintainer-1", "0cin", 0, CastAction{"nft-maintainer-1", "approve", 37}, "", nil},
		{"maintenance-cast-action", false, false, "Cast action to approve jeansoon as non fungible token provider address, Happy path. commit", "nft-maintainer-3", "0cin", 0, CastAction{"nft-maintainer-3", "approve", 37}, "", nil},

		//add nameservice fee collector with maintenance. (nft-mostafa is whitelisted.)
		{"maintenance", false, false, "38. Proposal, add token fee collector address, Happy path. commit", "nft-maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add token fee collector", "Add mostafa as nameservice fee collector", "fee", "", "", "", FeeCollector{Module: "nonFungible", Address: "nft-mostafa"}, "nft-maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as nameservice fee collector, Happy path. commit", "nft-maintainer-1", "0cin", 0, CastAction{"nft-maintainer-1", "approve", 38}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as nameservice fee collector, Happy path. commit", "nft-maintainer-3", "0cin", 0, CastAction{"nft-maintainer-3", "approve", 38}, "", nil},

		// Cast Action - Proposal 39
		// add token fee collector with maintenance. (to pass sdk ft test - add maintainer-1 as fungible token fee collector)
		{"maintenance", false, false, "38. Proposal, add token fee collector address, Happy path. commit", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add token fee collector", "Add maintainer-1 as token fee collector", "fee", "", "", "", FeeCollector{Module: "token", Address: "maintainer-1"}, "maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve maintainer-1 as token fee collector, Happy path. commit", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 39}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve maintainer-1 as token fee collector, Happy path. commit", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 39}, "", nil},
	}

	return tcs
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
