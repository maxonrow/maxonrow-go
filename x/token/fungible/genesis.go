package fungible

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	AuthorizedAddresses []sdkTypes.AccAddress `json:"authorised_addresses"`
	IssuerAddresses     []sdkTypes.AccAddress `json:"issuer_addresses"`
	ProviderAddresses   []sdkTypes.AccAddress `json:"provider_addresses"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func InitGenesis(ctx sdkTypes.Context, keeper *Keeper, genesisState GenesisState) {
	var validAuthorizedAddresses []sdkTypes.AccAddress

	for _, AuthorizedAddressesString := range genesisState.AuthorizedAddresses {
		authorisedAddress, err := sdkTypes.AccAddressFromBech32(AuthorizedAddressesString.String())

		if err != nil {
			panic("Invalid authorised address")
		}
		validAuthorizedAddresses = append(validAuthorizedAddresses, authorisedAddress)
	}
	keeper.SetAuthorisedAddresses(ctx, validAuthorizedAddresses)

	var validIssuerAddresses []sdkTypes.AccAddress
	for _, issuerAddress := range genesisState.IssuerAddresses {
		issuerAdd, err := sdkTypes.AccAddressFromBech32(issuerAddress.String())
		if err != nil {
			panic("Invalid issuer address")
		}
		validIssuerAddresses = append(validIssuerAddresses, issuerAdd)
	}
	keeper.SetIssuerAddresses(ctx, validIssuerAddresses)

	var validProviderAddresses []sdkTypes.AccAddress
	for _, providerAddress := range genesisState.ProviderAddresses {
		providerAdd, err := sdkTypes.AccAddressFromBech32(providerAddress.String())
		if err != nil {
			panic("Invalid provider address")
		}
		validProviderAddresses = append(validProviderAddresses, providerAdd)
	}
	keeper.SetProviderAddresses(ctx, validProviderAddresses)
}

func ExportGenesis(keeper *Keeper) GenesisState {
	return GenesisState{}
}
