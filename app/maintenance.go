package app

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/x/maintenance"
)

func (app *mxwApp) executeProposal(ctx sdkTypes.Context, proposal maintenance.Proposal) sdkTypes.Error {

	proposalID := proposal.ProposalID
	proposal, ok := app.maintenanceKeeper.GetProposal(ctx, proposalID)

	if !ok {
		return maintenance.ErrUnknownProposal(maintenance.DefaultCodespace, proposalID)
	}
	if proposal.Status != maintenance.StatusActive {
		return maintenance.ErrInactiveProposal(maintenance.DefaultCodespace, proposalID)
	}

	// Newly set addresses must be whitelisted.
	switch proposal.ProposalType() {
	case maintenance.ProposalTypeModifyFee:

		feeMaintainer, ok := proposal.ProposalData.(maintenance.FeeMaintainer)
		if !ok {
			return sdkTypes.ErrInternal("Converting to fee maintainer failed.")
		}
		executeErr := app.executeFeeMaintainerProposal(ctx, &feeMaintainer)
		if executeErr != nil {
			return executeErr
		}

	case maintenance.ProposalTypeModifyKyc:
		kycMaintainer, ok := proposal.ProposalData.(maintenance.KycMaintainer)
		if !ok {
			return sdkTypes.ErrInternal("Converting to kyc maintainer failed.")
		}
		executeErr := app.executeKycMaintainerProposal(ctx, kycMaintainer)
		if executeErr != nil {
			return executeErr
		}
	case maintenance.ProposalTypeModifyNameservice:
		nameserviceMaintainer, ok := proposal.ProposalData.(maintenance.NameserviceMaintainer)
		if !ok {
			return sdkTypes.ErrInternal("Converting to nameservice maintainer failed.")
		}
		executeErr := app.executeNameserviceProposal(ctx, nameserviceMaintainer)
		if executeErr != nil {
			return executeErr
		}
	case maintenance.ProposalTypeModifyToken:
		tokenMaintainer, ok := proposal.ProposalData.(maintenance.TokenMaintainer)
		if !ok {
			return sdkTypes.ErrInternal("Converting to token maintainer failed.")
		}
		executeErr := app.executeTokenProposal(ctx, tokenMaintainer)
		if executeErr != nil {
			return executeErr
		}
	// case maintenance.ProposalTypeModifyNonFungible:
	// 	nonFungibleMaintainer, ok := proposal.ProposalData.(maintenance.NonFungibleMaintainer)
	// 	if !ok {
	// 		return sdkTypes.ErrInternal("Converting to non fungible maintainer failed.")
	// 	}
	// 	executeErr := app.executeNonFungibleProposal(ctx, nonFungibleMaintainer)
	// 	if executeErr != nil {
	// 		return executeErr
	// 	}
	case maintenance.ProposalTypesModifyValidatorSet:
		whitelistValidator, ok := proposal.ProposalData.(maintenance.WhitelistValidator)
		if !ok {
			return sdkTypes.ErrInternal("Converting to whitelist validator failed.")
		}
		executeErr := app.executeWhitelistValidator(ctx, whitelistValidator)
		if executeErr != nil {
			return executeErr
		}

	default:
		return maintenance.ErrInvalidProposalType(maintenance.DefaultCodespace, proposal.ProposalType())
	}

	return nil
}

// Handle fee maintainer proposal
func (app *mxwApp) executeFeeMaintainerProposal(ctx sdkTypes.Context, feeMaintainer *maintenance.FeeMaintainer) sdkTypes.Error {
	switch feeMaintainer.Action {
	// TO-DO: Proper checking of empty slice.
	case maintenance.ADD:
		if feeMaintainer.FeeCollectors[0].Address != nil {
			for _, feeCollectorValue := range feeMaintainer.FeeCollectors {
				if !app.kycKeeper.IsWhitelisted(ctx, feeCollectorValue.Address) {
					return sdkTypes.ErrInternal("Address has to be whitelisted.")
				}
				app.feeKeeper.SetFeeCollectorAddresses(ctx, feeCollectorValue.Module, []sdkTypes.AccAddress{feeCollectorValue.Address})
			}
		}
		if feeMaintainer.AuthorisedAddresses[0] != nil {
			for _, authorisedAddress := range feeMaintainer.AuthorisedAddresses {
				if !app.kycKeeper.IsWhitelisted(ctx, authorisedAddress) {
					return sdkTypes.ErrInternal("Address has to be whitelisted.")
				}
			}
			app.feeKeeper.SetAuthorisedAddresses(ctx, feeMaintainer.AuthorisedAddresses)
		}
	case maintenance.REMOVE:
		if feeMaintainer.FeeCollectors[0].Address != nil {
			for _, feeCollectorValue := range feeMaintainer.FeeCollectors {
				if app.feeKeeper.IsFeeCollector(ctx, feeCollectorValue.Module, feeCollectorValue.Address) {
					app.feeKeeper.RemoveFeeCollectorAddress(ctx, feeCollectorValue.Module, feeCollectorValue.Address)
				}
			}
		}
		if feeMaintainer.AuthorisedAddresses[0] != nil {
			for _, authorisedAddress := range feeMaintainer.AuthorisedAddresses {
				if !app.feeKeeper.IsAuthorised(ctx, authorisedAddress) {
					return sdkTypes.ErrInternal("Address is not an authorisedAddress.")
				}
			}

			app.feeKeeper.RemoveAuthorisedAddresses(ctx, feeMaintainer.AuthorisedAddresses)
		}

	default:
		return sdkTypes.ErrInternal("Not recognised action.")
	}
	return nil
}

// Handle kyc proposal
func (app *mxwApp) executeKycMaintainerProposal(ctx sdkTypes.Context, kycMaintainer maintenance.KycMaintainer) sdkTypes.Error {
	switch kycMaintainer.Action {
	case maintenance.ADD:
		if kycMaintainer.AuthorisedAddresses[0] != nil {
			for _, authorisedAddress := range kycMaintainer.AuthorisedAddresses {
				if !app.kycKeeper.IsWhitelisted(ctx, authorisedAddress) {
					return sdkTypes.ErrInternal("Address has to be whitelisted.")
				}
			}
			app.kycKeeper.SetAuthorisedAddresses(ctx, kycMaintainer.AuthorisedAddresses)
		}
		if kycMaintainer.IssuerAddresses[0] != nil {
			for _, issuerAddress := range kycMaintainer.IssuerAddresses {
				if !app.kycKeeper.IsWhitelisted(ctx, issuerAddress) {
					return sdkTypes.ErrInternal("Address has to be whitelisted.")
				}
			}
			app.kycKeeper.SetIssuerAddresses(ctx, kycMaintainer.IssuerAddresses)
		}
		if kycMaintainer.ProviderAddresses[0] != nil {
			for _, providerAddress := range kycMaintainer.ProviderAddresses {
				if !app.kycKeeper.IsWhitelisted(ctx, providerAddress) {
					return sdkTypes.ErrInternal("Address has to be whitelisted.")
				}
			}
			app.kycKeeper.SetProviderAddresses(ctx, kycMaintainer.ProviderAddresses)
		}
	case maintenance.REMOVE:
		if kycMaintainer.AuthorisedAddresses[0] != nil {
			for _, authorisedAddress := range kycMaintainer.AuthorisedAddresses {
				if !app.kycKeeper.IsAuthorised(ctx, authorisedAddress) {
					return sdkTypes.ErrInternal("Address is not an authorised address.")
				}
			}
			app.kycKeeper.RemoveAuthorisedAddresses(ctx, kycMaintainer.AuthorisedAddresses)
		}
		if kycMaintainer.IssuerAddresses[0] != nil {
			for _, issuerAddress := range kycMaintainer.IssuerAddresses {
				if !app.kycKeeper.IsIssuer(ctx, issuerAddress) {
					return sdkTypes.ErrInternal("Address is not an issuer.")
				}
			}
			app.kycKeeper.RemoveIssuerAddresses(ctx, kycMaintainer.IssuerAddresses)

		}
		if kycMaintainer.ProviderAddresses[0] != nil {
			for _, providerAddress := range kycMaintainer.ProviderAddresses {
				if !app.kycKeeper.IsProvider(ctx, providerAddress) {
					return sdkTypes.ErrInternal("Address is not a provider.")
				}
			}
			app.kycKeeper.RemoveProviderAddresses(ctx, kycMaintainer.ProviderAddresses)
		}
	default:
		return sdkTypes.ErrInternal("Not recognised action.")
	}
	return nil
}

// Handle nameservice proposal
func (app *mxwApp) executeNameserviceProposal(ctx sdkTypes.Context, nameserviceMaintainer maintenance.NameserviceMaintainer) sdkTypes.Error {
	switch nameserviceMaintainer.Action {
	case maintenance.ADD:
		if nameserviceMaintainer.AuthorisedAddresses[0] != nil {
			for _, authorisedAddress := range nameserviceMaintainer.AuthorisedAddresses {
				if !app.kycKeeper.IsWhitelisted(ctx, authorisedAddress) {
					return sdkTypes.ErrInternal("Address has to be whitelisted.")
				}
			}
			app.nsKeeper.SetAuthorisedAddresses(ctx, nameserviceMaintainer.AuthorisedAddresses)
		}
		if nameserviceMaintainer.IssuerAddresses[0] != nil {
			for _, issuerAddress := range nameserviceMaintainer.IssuerAddresses {
				if !app.kycKeeper.IsWhitelisted(ctx, issuerAddress) {
					return sdkTypes.ErrInternal("Address has to be whitelisted.")
				}
			}
			app.nsKeeper.SetIssuerAddresses(ctx, nameserviceMaintainer.IssuerAddresses)

		}
		if nameserviceMaintainer.ProviderAddresses[0] != nil {
			for _, providerAddress := range nameserviceMaintainer.ProviderAddresses {
				if !app.kycKeeper.IsWhitelisted(ctx, providerAddress) {
					return sdkTypes.ErrInternal("Address has to be whitelisted.")
				}
			}
			app.nsKeeper.SetProviderAddresses(ctx, nameserviceMaintainer.ProviderAddresses)

		}
	case maintenance.REMOVE:
		if nameserviceMaintainer.AuthorisedAddresses[0] != nil {
			for _, authorisedAddress := range nameserviceMaintainer.AuthorisedAddresses {
				if !app.nsKeeper.IsAuthorised(ctx, authorisedAddress) {
					return sdkTypes.ErrInternal("Address is not an authorised address.")
				}
			}
			app.nsKeeper.RemoveAuthorisedAddresses(ctx, nameserviceMaintainer.AuthorisedAddresses)
		}
		if nameserviceMaintainer.IssuerAddresses[0] != nil {
			for _, issuerAddress := range nameserviceMaintainer.IssuerAddresses {
				if !app.nsKeeper.IsIssuer(ctx, issuerAddress) {
					return sdkTypes.ErrInternal("Address is not an issuer.")
				}
			}
			app.nsKeeper.RemoveIssuerAddresses(ctx, nameserviceMaintainer.IssuerAddresses)

		}
		if nameserviceMaintainer.ProviderAddresses[0] != nil {
			for _, providerAddress := range nameserviceMaintainer.ProviderAddresses {
				if !app.nsKeeper.IsProvider(ctx, providerAddress) {
					return sdkTypes.ErrInternal("Address is not a provider.")
				}
			}
			app.nsKeeper.RemoveProviderAddresses(ctx, nameserviceMaintainer.ProviderAddresses)

		}
	default:
		return sdkTypes.ErrInternal("Not recognised action.")
	}
	return nil
}

// Handle token proposal
func (app *mxwApp) executeTokenProposal(ctx sdkTypes.Context, tokenMaintainer maintenance.TokenMaintainer) sdkTypes.Error {
	switch tokenMaintainer.Action {
	case maintenance.ADD:
		if tokenMaintainer.AuthorisedAddresses[0] != nil {
			for _, authorisedAddress := range tokenMaintainer.AuthorisedAddresses {
				if !app.kycKeeper.IsWhitelisted(ctx, authorisedAddress) {
					return sdkTypes.ErrInternal("Address has to be whitelisted.")
				}
			}
			app.tokenKeeper.SetAuthorisedAddresses(ctx, tokenMaintainer.AuthorisedAddresses)
		}
		if tokenMaintainer.IssuerAddresses[0] != nil {
			for _, issuerAddress := range tokenMaintainer.IssuerAddresses {
				if !app.kycKeeper.IsWhitelisted(ctx, issuerAddress) {
					return sdkTypes.ErrInternal("Address has to be whitelisted.")
				}
			}
			app.tokenKeeper.SetIssuerAddresses(ctx, tokenMaintainer.IssuerAddresses)

		}
		if tokenMaintainer.ProviderAddresses[0] != nil {
			for _, providerAddress := range tokenMaintainer.ProviderAddresses {
				if !app.kycKeeper.IsWhitelisted(ctx, providerAddress) {
					return sdkTypes.ErrInternal("Address has to be whitelisted.")
				}
			}
			app.tokenKeeper.SetProviderAddresses(ctx, tokenMaintainer.ProviderAddresses)

		}
	case maintenance.REMOVE:

		if tokenMaintainer.AuthorisedAddresses[0] != nil {
			for _, authorisedAddress := range tokenMaintainer.AuthorisedAddresses {
				if !app.tokenKeeper.IsAuthorised(ctx, authorisedAddress) {
					return sdkTypes.ErrInternal("Address is not an authorised address.")
				}
			}
			app.tokenKeeper.RemoveAuthorisedAddresses(ctx, tokenMaintainer.AuthorisedAddresses)
		}
		if tokenMaintainer.IssuerAddresses[0] != nil {
			for _, issuerAddress := range tokenMaintainer.IssuerAddresses {
				if !app.tokenKeeper.IsIssuer(ctx, issuerAddress) {
					return sdkTypes.ErrInternal("Address is not an issuer.")
				}
			}
			app.tokenKeeper.RemoveIssuerAddresses(ctx, tokenMaintainer.IssuerAddresses)

		}
		if tokenMaintainer.ProviderAddresses[0] != nil {
			for _, providerAddress := range tokenMaintainer.ProviderAddresses {
				if !app.tokenKeeper.IsProvider(ctx, providerAddress) {
					return sdkTypes.ErrInternal("Address is not a provider.")
				}
				app.tokenKeeper.RemoveProviderAddresses(ctx, tokenMaintainer.ProviderAddresses)
			}

		}
	default:
		return sdkTypes.ErrInternal("Not recognised action.")
	}
	return nil
}

// Handle nonFungible proposal
// func (app *mxwApp) executeNonFungibleProposal(ctx sdkTypes.Context, nonFungibleMaintainer maintenance.NonFungibleMaintainer) sdkTypes.Error {
// 	switch nonFungibleMaintainer.Action {
// 	case maintenance.ADD:
// 		if nonFungibleMaintainer.AuthorisedAddresses[0] != nil {
// 			for _, authorisedAddress := range nonFungibleMaintainer.AuthorisedAddresses {
// 				if !app.kycKeeper.IsWhitelisted(ctx, authorisedAddress) {
// 					return sdkTypes.ErrInternal("Address has to be whitelisted.")
// 				}
// 			}
// 			app.nonFungibleTokenKeeper.SetAuthorisedAddresses(ctx, nonFungibleMaintainer.AuthorisedAddresses)
// 		}
// 		if nonFungibleMaintainer.IssuerAddresses[0] != nil {
// 			for _, issuerAddress := range nonFungibleMaintainer.IssuerAddresses {
// 				if !app.kycKeeper.IsWhitelisted(ctx, issuerAddress) {
// 					return sdkTypes.ErrInternal("Address has to be whitelisted.")
// 				}
// 			}
// 			app.nonFungibleTokenKeeper.SetIssuerAddresses(ctx, nonFungibleMaintainer.IssuerAddresses)

// 		}
// 		if nonFungibleMaintainer.ProviderAddresses[0] != nil {
// 			for _, providerAddress := range nonFungibleMaintainer.ProviderAddresses {
// 				if !app.kycKeeper.IsWhitelisted(ctx, providerAddress) {
// 					return sdkTypes.ErrInternal("Address has to be whitelisted.")
// 				}
// 			}
// 			app.nonFungibleTokenKeeper.SetProviderAddresses(ctx, nonFungibleMaintainer.ProviderAddresses)

// 		}
// 	case maintenance.REMOVE:

// 		if nonFungibleMaintainer.AuthorisedAddresses[0] != nil {
// 			for _, authorisedAddress := range nonFungibleMaintainer.AuthorisedAddresses {
// 				if !app.nonFungibleTokenKeeper.IsAuthorised(ctx, authorisedAddress) {
// 					return sdkTypes.ErrInternal("Address is not an authorised address.")
// 				}
// 			}
// 			app.nonFungibleTokenKeeper.RemoveAuthorisedAddresses(ctx, nonFungibleMaintainer.AuthorisedAddresses)
// 		}
// 		if nonFungibleMaintainer.IssuerAddresses[0] != nil {
// 			for _, issuerAddress := range nonFungibleMaintainer.IssuerAddresses {
// 				if !app.nonFungibleTokenKeeper.IsIssuer(ctx, issuerAddress) {
// 					return sdkTypes.ErrInternal("Address is not an issuer.")
// 				}
// 			}
// 			app.nonFungibleTokenKeeper.RemoveIssuerAddresses(ctx, nonFungibleMaintainer.IssuerAddresses)

// 		}
// 		if nonFungibleMaintainer.ProviderAddresses[0] != nil {
// 			for _, providerAddress := range nonFungibleMaintainer.ProviderAddresses {
// 				if !app.nonFungibleTokenKeeper.IsProvider(ctx, providerAddress) {
// 					return sdkTypes.ErrInternal("Address is not a provider.")
// 				}
// 				app.nonFungibleTokenKeeper.RemoveProviderAddresses(ctx, nonFungibleMaintainer.ProviderAddresses)
// 			}

// 		}
// 	default:
// 		return sdkTypes.ErrInternal("Not recognised action.")
// 	}
// 	return nil
// }

func (app *mxwApp) executeWhitelistValidator(ctx sdkTypes.Context, whitelistValidator maintenance.WhitelistValidator) sdkTypes.Error {
	switch whitelistValidator.Action {
	case maintenance.ADD:
		if whitelistValidator.ValidatorPubKeys[0] != nil {
			app.maintenanceKeeper.WhitelistValidator(ctx, whitelistValidator.ValidatorPubKeys)
		}
	case maintenance.REMOVE:
		if whitelistValidator.ValidatorPubKeys[0] != nil {

			app.maintenanceKeeper.RevokeValidator(ctx, whitelistValidator.ValidatorPubKeys)
		}
	default:
		return sdkTypes.ErrInternal("Not recognised action.")
	}
	return nil
}
