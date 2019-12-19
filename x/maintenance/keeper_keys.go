package maintenance

import (
	"fmt"
)

const (
	MaintainerAdd = "maintainer"
)

var (
	KeyNextProposalID = []byte("newProposalID")

	PrefixValidatorSet = []byte("0x033")
)

// Key for getting a specific proposal from the store
func getProposalKey(proposalID uint64) []byte {
	return []byte(fmt.Sprintf("maintenance/proposal:%d", proposalID))
}

//* key
func getMaintenanceAddressKey() string {
	return fmt.Sprintf("maintenace/sys_addr")
}

func getValidatorSetKey() []byte {
	return PrefixValidatorSet
}
