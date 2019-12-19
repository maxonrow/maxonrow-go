package nameservice

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	AuthorisedAddresses []sdkTypes.AccAddress `json:"authorised_addresses"`
	IssuerAddresses     []sdkTypes.AccAddress `json:"issuer_addresses"`
	ProviderAddresses   []sdkTypes.AccAddress `json:"provider_addresses"`
	GenesisAliasOwners  []genesisAliasOwner   `json:"genesis_alias_owners"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}
func InitGenesis(ctx sdkTypes.Context, keeper Keeper, genesisState GenesisState) {

	var validAuthorisedAddresses []sdkTypes.AccAddress
	for _, authorisedAddressesString := range genesisState.AuthorisedAddresses {
		authorisedAddress, err := sdkTypes.AccAddressFromBech32(authorisedAddressesString.String())
		if err != nil {
			panic("Invalid authorised address")
		}

		validAuthorisedAddresses = append(validAuthorisedAddresses, authorisedAddress)
	}
	keeper.SetAuthorisedAddresses(ctx, validAuthorisedAddresses)

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

	for _, genesisAliasOwner := range genesisState.GenesisAliasOwners {
		alias := &Alias{
			Name:     genesisAliasOwner.alias,
			Owner:    genesisAliasOwner.owner,
			Metadata: "genesis alias owner",
			Approved: false,
			Fee:      sdkTypes.NewUintFromString("0"),
		}

		aliasOwner := &AliasOwner{
			Name:     genesisAliasOwner.alias,
			Approved: false,
		}
		keeper.storeAlias(ctx, genesisAliasOwner.alias, alias, aliasOwner)

		alias.Approved = true
		aliasOwner.Approved = true
		keeper.setAlias(ctx, alias.Name, alias, aliasOwner)
	}
}

func ExportGenesis(keeper *Keeper) GenesisState {
	return GenesisState{
		AuthorisedAddresses: keeper.authorizedAddresses,
	}
}

type genesisAliasOwner struct {
	alias string              `json:"alias"`
	owner sdkTypes.AccAddress `json:"owner"`
}
