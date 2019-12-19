package kyc

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

type GenesisState struct {
	AuthorizedAddresses  []sdkTypes.AccAddress `json:"authorised_addresses"`
	IssuerAddresses      []sdkTypes.AccAddress `json:"issuer_addresses"`
	ProviderAddresses    []sdkTypes.AccAddress `json:"provider_addresses"`
	WhitelistedAddresses []sdkTypes.AccAddress `json:"whitelisted_addresses"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}
