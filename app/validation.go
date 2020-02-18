package app

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/bank"
	"github.com/maxonrow/maxonrow-go/x/fee"
	"github.com/maxonrow/maxonrow-go/x/kyc"
	"github.com/maxonrow/maxonrow-go/x/maintenance"
	"github.com/maxonrow/maxonrow-go/x/nameservice"
	fungible "github.com/maxonrow/maxonrow-go/x/token/fungible"
	nonFungible "github.com/maxonrow/maxonrow-go/x/token/nonfungible"
)

func (app *mxwApp) validateMsg(ctx sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Error {

	switch msg := msg.(type) {
	case nameservice.MsgCreateAlias:
		if !app.feeKeeper.IsFeeCollector(ctx, "nameservice", msg.Fee.To) {
			return sdkTypes.ErrInvalidAddress("Fee collector invalid.")
		}

		ownerAcc := app.accountKeeper.GetAccount(ctx, msg.Owner)
		appFeeAmt, err := sdkTypes.ParseCoins(msg.Fee.Value + types.CIN)
		if err != nil {
			return sdkTypes.ErrInternal("Invalid fee amount.")
		}
		if ownerAcc.GetCoins().IsAllLT(appFeeAmt) {
			return sdkTypes.ErrInternal("Insufficient balance to pay for application fee.")
		}

		if app.nsKeeper.IsAliasExists(ctx, msg.Name) {
			return types.ErrAliasIsInUsed()
		}
	case nameservice.MsgSetAliasStatus:
		if !app.nsKeeper.IsAuthorised(ctx, msg.GetSigners()[0]) {
			return sdkTypes.ErrInvalidAddress("Not authorised to set alias status.")
		}
		err := app.nsKeeper.ValidateSignatures(ctx, msg)
		if err != nil {
			return err
		}
	case kyc.MsgWhitelist:
		if app.kycKeeper.IsKycAddressExist(ctx, msg.KycData.Payload.Kyc.KycAddress) {
			return sdkTypes.NewError("mxw", 1001, "Kyc Address duplicated.")
		}

		whitelistSignErr := app.kycKeeper.ValidateSignatures(ctx, msg)
		if whitelistSignErr != nil {
			return whitelistSignErr
		}

		if !app.kycKeeper.IsAuthorised(ctx, msg.Owner) {
			return sdkTypes.ErrUnauthorized("Not authorized to whitelist.")
		}

	case kyc.MsgRevokeWhitelist:
		revokeWhitelistSignErr := app.kycKeeper.ValidateRevokeWhitelistSignatures(ctx, msg)
		if revokeWhitelistSignErr != nil {
			return revokeWhitelistSignErr
		}
		addr := msg.RevokePayload.RevokeKycData.To
		if !app.kycKeeper.IsWhitelisted(ctx, addr) {
			return sdkTypes.ErrInternal("Target address is not whitelisted.")
		}

		if !app.kycKeeper.IsAuthorised(ctx, msg.Owner) {
			return sdkTypes.ErrUnauthorized("Not authorized to whitelist.")
		}
	case fungible.MsgCreateFungibleToken:
		ownerAcc := app.accountKeeper.GetAccount(ctx, msg.Owner)
		appFeeAmt, err := sdkTypes.ParseCoins(msg.Fee.Value + types.CIN)
		if err != nil {
			return sdkTypes.ErrInternal("Invalid fee amount.")
		}
		if ownerAcc.GetCoins().IsAllLT(appFeeAmt) {
			return sdkTypes.ErrInternal("Insufficient balance to pay for application fee.")
		}

		if !app.feeKeeper.IsFeeCollector(ctx, "token", msg.Fee.To) {
			return sdkTypes.ErrInvalidAddress("Fee collector invalid.")
		}

		if app.tokenKeeper.TokenExists(ctx, msg.Symbol) {
			return types.ErrTokenExists(msg.Symbol)
		}
	case fungible.MsgSetFungibleTokenStatus:
		// TO-DO: revisit
		if !app.tokenKeeper.IsAuthorised(ctx, msg.GetSigners()[0]) {
			return sdkTypes.ErrUnauthorized("Not authorised to approve.")
		}
		if !app.tokenKeeper.TokenExists(ctx, msg.Payload.Token.Symbol) {
			return types.ErrInvalidTokenSymbol(msg.Payload.Token.Symbol)
		}
		err := app.tokenKeeper.ValidateSignatures(ctx, msg)
		if err != nil {
			return err
		}
		if msg.Payload.Token.Status == fungible.ApproveToken {

			for _, val := range msg.Payload.Token.TokenFees {
				if !app.feeKeeper.FeeSettingExists(ctx, val.FeeName) {
					return types.ErrFeeSettingNotExists(val.FeeName)
				}

				if !fee.ContainAction(val.Action) {
					return types.ErrInvalidTokenAction()
				}
			}

			if app.tokenKeeper.CheckApprovedToken(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenAlreadyApproved(msg.Payload.Token.Symbol)
			}
		}
		if msg.Payload.Token.Status == fungible.FreezeToken {
			if app.tokenKeeper.IsTokenFrozen(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenFrozen()
			}
			if !app.tokenKeeper.CheckApprovedToken(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenInvalid()
			}
		}
		if msg.Payload.Token.Status == fungible.UnfreezeToken {
			if !app.tokenKeeper.IsTokenFrozen(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenUnFrozen()
			}
		}

		if msg.Payload.Token.Status == fungible.ApproveTransferTokenOwnership || msg.Payload.Token.Status == fungible.RejectTransferTokenOwnership {
			if !app.tokenKeeper.IsVerifyableTransferTokenOwnership(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenTransferTokenOwnershipApproved()
			}
		}
	case fungible.MsgSetFungibleTokenAccountStatus:
		if msg.TokenAccountPayload.TokenAccount.Status == fungible.FreezeTokenAccount {
			if app.tokenKeeper.IsFungibleTokenAccountFrozen(ctx, msg.TokenAccountPayload.TokenAccount.Account, msg.TokenAccountPayload.TokenAccount.Symbol) {
				return types.ErrTokenAccountFrozen()
			}
		}
		if msg.TokenAccountPayload.TokenAccount.Status == fungible.UnfreezeTokenAccount {
			if !app.tokenKeeper.IsFungibleTokenAccountFrozen(ctx, msg.TokenAccountPayload.TokenAccount.Account, msg.TokenAccountPayload.TokenAccount.Symbol) {
				return types.ErrTokenAccountUnFrozen()
			}
		}
	case fungible.MsgTransferFungibleToken:
		if !app.tokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
		if app.tokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}
		if app.tokenKeeper.IsFungibleTokenAccountFrozen(ctx, msg.From, msg.Symbol) {
			return types.ErrTokenAccountFrozen()
		}
		//check receiver
		if app.tokenKeeper.IsFungibleTokenAccountFrozen(ctx, msg.To, msg.Symbol) {
			return types.ErrTokenAccountFrozen()
		}
	case fungible.MsgMintFungibleToken:
		if !app.tokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}

		if app.tokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}

		// check sender
		if app.tokenKeeper.IsFungibleTokenAccountFrozen(ctx, msg.Owner, msg.Symbol) {
			return types.ErrTokenAccountFrozen()
		}

		//check receiver
		if app.tokenKeeper.IsFungibleTokenAccountFrozen(ctx, msg.To, msg.Symbol) {
			return types.ErrTokenAccountFrozen()
		}
	case fungible.MsgBurnFungibleToken:
		if !app.tokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
		if app.tokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}

		var token = new(nonFungible.Token)
		isTokenExisted := app.nonFungibleTokenKeeper.GetTokenDataInfo(ctx, msg.Symbol, token)
		if isTokenExisted == true {
			if !token.Owner.Equals(msg.From) {
				return sdkTypes.ErrUnknownRequest("Invalid token owner.")
			}
		}
	case fungible.MsgTransferFungibleTokenOwnership:
		if !app.tokenKeeper.IsTokenOwnershipTransferrable(ctx, msg.Symbol) {
			return types.ErrInvalidTokenAction()
		}
		if !app.tokenKeeper.IsTokenOwner(ctx, msg.Symbol, msg.From) {
			return types.ErrInvalidTokenOwner()
		}
		if !app.tokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
	case fungible.MsgAcceptFungibleTokenOwnership:
		if !app.tokenKeeper.IsTokenOwnershipAcceptable(ctx, msg.Symbol) {
			return types.ErrInvalidTokenAction()
		}

		if !app.tokenKeeper.IsTokenNewOwner(ctx, msg.Symbol, msg.From) {
			return types.ErrInvalidTokenNewOwner()
		}

		if !app.tokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
	case nonFungible.MsgCreateNonFungibleToken:
		ownerAcc := app.accountKeeper.GetAccount(ctx, msg.Owner)
		appFeeAmt, err := sdkTypes.ParseCoins(msg.Fee.Value + types.CIN)
		if err != nil {
			return sdkTypes.ErrInternal("Invalid fee amount.")
		}
		if ownerAcc.GetCoins().IsAllLT(appFeeAmt) {
			return sdkTypes.ErrInternal("Insufficient balance to pay for application fee.")
		}

		if !app.feeKeeper.IsFeeCollector(ctx, "nonFungible", msg.Fee.To) {
			return sdkTypes.ErrInvalidAddress("Fee collector invalid.")
		}

		if app.nonFungibleTokenKeeper.TokenExists(ctx, msg.Symbol) {
			return types.ErrTokenExists(msg.Symbol)
		}
	case nonFungible.MsgSetNonFungibleTokenStatus:
		// TO-DO: revisit
		if !app.nonFungibleTokenKeeper.IsAuthorised(ctx, msg.GetSigners()[0]) {
			return sdkTypes.ErrUnauthorized("Not authorised to approve.")
		}
		if !app.nonFungibleTokenKeeper.TokenExists(ctx, msg.Payload.Token.Symbol) {
			return types.ErrInvalidTokenSymbol(msg.Payload.Token.Symbol)
		}
		err := app.nonFungibleTokenKeeper.ValidateSignatures(ctx, msg)
		if err != nil {
			return err
		}
		if msg.Payload.Token.Status == fungible.ApproveToken {

			for _, val := range msg.Payload.Token.TokenFees {
				if !app.feeKeeper.FeeSettingExists(ctx, val.FeeName) {
					return types.ErrFeeSettingNotExists(val.FeeName)
				}

				if !fee.ContainAction(val.Action) {
					return types.ErrInvalidTokenAction()
				}
			}

			if app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenAlreadyApproved(msg.Payload.Token.Symbol)
			}
		}
		if msg.Payload.Token.Status == fungible.FreezeToken {
			if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenFrozen()
			}
			if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenInvalid()
			}
		}
		if msg.Payload.Token.Status == fungible.UnfreezeToken {
			if !app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenUnFrozen()
			}
		}

		if msg.Payload.Token.Status == fungible.ApproveTransferTokenOwnership || msg.Payload.Token.Status == fungible.RejectTransferTokenOwnership {
			if !app.nonFungibleTokenKeeper.IsVerifyableTransferTokenOwnership(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenTransferTokenOwnershipApproved()
			}
		}
	case nonFungible.MsgSetNonFungibleItemStatus:
		var token = new(nonFungible.Token)
		if exists := app.nonFungibleTokenKeeper.GetTokenDataInfo(ctx, msg.ItemPayload.Item.Symbol, token); !exists {
			return sdkTypes.ErrUnknownRequest("No such non fungible token.")
		}

		nonFungibleItem := app.nonFungibleTokenKeeper.GetNonFungibleItem(ctx, msg.ItemPayload.Item.Symbol, msg.ItemPayload.Item.ItemID)
		if nonFungibleItem == nil {
			return sdkTypes.ErrUnknownRequest("No such item to freeze.")
		}

		if !app.nonFungibleTokenKeeper.IsAuthorised(ctx, msg.Owner) {
			return sdkTypes.ErrUnauthorized("Not authorised to unfreeze token account.")
		}

		signatureErr := app.nonFungibleTokenKeeper.ValidateSignatures(ctx, msg)
		if signatureErr != nil {
			return signatureErr
		}

		if msg.ItemPayload.Item.Status == nonFungible.FreezeItem {
			if app.nonFungibleTokenKeeper.IsNonFungibleItemFrozen(ctx, msg.ItemPayload.Item.Symbol, msg.ItemPayload.Item.ItemID) {
				return types.ErrTokenAccountFrozen()
			}
		}
		if msg.ItemPayload.Item.Status == nonFungible.UnfreezeItem {
			if !app.nonFungibleTokenKeeper.IsNonFungibleItemFrozen(ctx, msg.ItemPayload.Item.Symbol, msg.ItemPayload.Item.ItemID) {
				return types.ErrTokenAccountUnFrozen()
			}
		}

	case nonFungible.MsgTransferNonFungibleItem:
		if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
		if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}
		if app.nonFungibleTokenKeeper.IsNonFungibleItemFrozen(ctx, msg.Symbol, msg.ItemID) {
			return types.ErrTokenAccountFrozen()
		}

		// 1. [Transfer non fungible token item - Invalid Item-ID]
		nonFungibleItem := app.nonFungibleTokenKeeper.GetNonFungibleItem(ctx, msg.Symbol, msg.ItemID)
		if nonFungibleItem == nil {
			return sdkTypes.ErrUnknownRequest("Invalid Item ID.")
		}

	case nonFungible.MsgMintNonFungibleItem:
		if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
		if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}

		//1. [Mint (by Public==TRUE) non fungible token(TNFT-public-01) - Error, Public token can only be minted to itself.]
		var token = new(nonFungible.Token)
		app.nonFungibleTokenKeeper.GetTokenDataInfo(ctx, msg.Symbol, token)
		if token.Flags.HasFlag(0x0080) {
			ownerAcc := msg.Owner
			newOwnerAcc := msg.To
			if !ownerAcc.Equals(newOwnerAcc) {
				return sdkTypes.ErrInternal("Public token can only be minted to oneself.")
			}
		}

	case nonFungible.MsgBurnNonFungibleItem:
		if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
		if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}
		// 1. [Burn non fungible token item - Invalid Item-owner]
		item := app.nonFungibleTokenKeeper.GetNonFungibleItem(ctx, msg.Symbol, msg.ItemID)
		itemOwner := app.nonFungibleTokenKeeper.GetNonFungibleItemOwnerInfo(ctx, msg.Symbol, msg.ItemID)
		if item == nil {
			return types.ErrTokenItemIDInUsed()
		}
		if !itemOwner.Equals(msg.From) {
			return types.ErrInvalidItemOwner()
		}

	case nonFungible.MsgTransferNonFungibleTokenOwnership:
		if !app.nonFungibleTokenKeeper.IsTokenOwnershipTransferrable(ctx, msg.Symbol) {
			return types.ErrInvalidTokenAction()
		}
		if !app.nonFungibleTokenKeeper.IsTokenOwner(ctx, msg.Symbol, msg.From) {
			return types.ErrInvalidTokenOwner()
		}
		if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
	case nonFungible.MsgAcceptNonFungibleTokenOwnership:
		if !app.nonFungibleTokenKeeper.IsTokenOwnershipAcceptable(ctx, msg.Symbol) {
			return types.ErrInvalidTokenAction()
		}

		if !app.nonFungibleTokenKeeper.IsTokenNewOwner(ctx, msg.Symbol, msg.From) {
			return types.ErrInvalidTokenNewOwner()
		}

		if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
	case nonFungible.MsgEndorsement:
		if !app.nonFungibleTokenKeeper.IsTokenEndorser(ctx, msg.Symbol, msg.From) {
			return types.ErrInvalidEndorser()
		}

		// 1. [endorse a nonfungible item - Invalid Item-ID]
		item := app.nonFungibleTokenKeeper.GetNonFungibleItem(ctx, msg.Symbol, msg.ItemID)
		if item == nil {
			return types.ErrTokenInvalid()
		}
		// 2. [endorse a nonfungible item - Invalid Token Symbol]
		if err := nonFungible.ValidateSymbol(msg.Symbol); err != nil {
			return err
		}

	case nonFungible.MsgUpdateNFTMetadata:
		if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
		if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}
	case nonFungible.MsgUpdateItemMetadata:
		if app.nonFungibleTokenKeeper.IsNonFungibleItemFrozen(ctx, msg.Symbol, msg.ItemID) {
			return types.ErrTokenAccountFrozen()
		}
		if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}
		// 1. [Update Item Metadata non fungible token - Invalid Item Id.]
		item := app.nonFungibleTokenKeeper.GetNonFungibleItem(ctx, msg.Symbol, msg.ItemID)
		if item == nil {
			return types.ErrTokenInvalid()
		}

		// 2. [Update Item Metadata non fungible token - Item owner not match.]
		itemOwner := app.nonFungibleTokenKeeper.GetNonFungibleItemOwnerInfo(ctx, msg.Symbol, msg.ItemID)
		if itemOwner == nil {
			return sdkTypes.ErrUnknownRequest("Item owner not match.")
		} else {
			if !itemOwner.Equals(msg.From) {
				return sdkTypes.ErrUnknownRequest("Item owner not match.")
			}
		}
		// 3. [Update Item Metadata non fungible token - Invalid token symbol.]
		if err := nonFungible.ValidateSymbol(msg.Symbol); err != nil {
			return err
		}
	case maintenance.MsgProposal:
		if !app.maintenanceKeeper.IsMaintainers(ctx, msg.Proposer) {
			return sdkTypes.ErrUnauthorized("Not authorised to submit proposal.")
		}
	case maintenance.MsgCastAction:
		if !app.maintenanceKeeper.IsMaintainers(ctx, msg.Owner) {
			return sdkTypes.ErrUnauthorized("Not authorised to cast vote.")
		}
		if !app.maintenanceKeeper.IsProposalActive(ctx, msg.ProposalID) {
			return maintenance.ErrUnknownProposal(maintenance.DefaultCodespace, msg.ProposalID)
		}
	case staking.MsgCreateValidator:
		if !app.maintenanceKeeper.IsValidator(ctx, msg.PubKey) {
			return sdkTypes.ErrUnauthorized("Not authorised validator address.")
		}

	case fee.MsgSysFeeSetting:
		if !app.feeKeeper.IsAuthorised(ctx, msg.Issuer) {
			return sdkTypes.ErrUnauthorized("Not authorised to create fee setting.")
		}
	case bank.MsgMxwSend:
		if !app.bankKeeper.HasCoins(ctx, msg.FromAddress, msg.Amount) {
			return sdkTypes.ErrInsufficientCoins("Insufficient balance to do transaction.")
		}

	default:
		return nil
	}
	return nil
}
