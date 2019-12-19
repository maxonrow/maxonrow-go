package kyc

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdkTypes.Context, keeper *Keeper, genesisState GenesisState) {
	var validAuthorizedAddresses []sdkTypes.AccAddress
	for _, AuthorizedAddressesString := range genesisState.AuthorizedAddresses {

		authorisedAddress, err := sdkTypes.AccAddressFromBech32(AuthorizedAddressesString.String())

		if err != nil {
			panic("Invalid authorised address")
		}

		keeper.Whitelist(ctx, authorisedAddress, authorisedAddress.String())
		validAuthorizedAddresses = append(validAuthorizedAddresses, authorisedAddress)
	}
	keeper.SetAuthorisedAddresses(ctx, validAuthorizedAddresses)

	for _, addressString := range genesisState.WhitelistedAddresses {
		address, err := sdkTypes.AccAddressFromBech32(addressString.String())
		if err != nil {
			panic("Invalid whitelisted address")
		}

		keeper.Whitelist(ctx, address, address.String())
	}

	var validIssuerAddresses []sdkTypes.AccAddress
	for _, issuerAddress := range genesisState.IssuerAddresses {
		issuerAdd, err := sdkTypes.AccAddressFromBech32(issuerAddress.String())

		if err != nil {
			panic("Invalid issuer address")
		}

		keeper.Whitelist(ctx, issuerAdd, issuerAdd.String())
		validIssuerAddresses = append(validIssuerAddresses, issuerAdd)
	}
	keeper.SetIssuerAddresses(ctx, validIssuerAddresses)

	var validProviderAddresses []sdkTypes.AccAddress
	for _, providerAddress := range genesisState.ProviderAddresses {
		providerAdd, err := sdkTypes.AccAddressFromBech32(providerAddress.String())

		if err != nil {
			panic("Invalid provider address")
		}

		keeper.Whitelist(ctx, providerAdd, providerAdd.String())
		validProviderAddresses = append(validProviderAddresses, providerAdd)
	}
	keeper.SetProviderAddresses(ctx, validProviderAddresses)
}

func ExportGenesis(keeper *Keeper) GenesisState {
	return GenesisState{}
}
