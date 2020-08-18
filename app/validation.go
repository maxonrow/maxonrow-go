package app

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/utils"
	"github.com/maxonrow/maxonrow-go/x/auth"
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

		ownerAcc := utils.GetAccount(ctx, app.accountKeeper, msg.Owner)
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
			return types.ErrKycDuplicated()
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
		ownerAcc := utils.GetAccount(ctx, app.accountKeeper, msg.Owner)
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

		if app.fungibleTokenKeeper.TokenExists(ctx, msg.Symbol) {
			return types.ErrTokenExists(msg.Symbol)
		}
	case fungible.MsgSetFungibleTokenStatus:
		// TO-DO: revisit
		if !app.fungibleTokenKeeper.IsAuthorised(ctx, msg.GetSigners()[0]) {
			return sdkTypes.ErrUnauthorized("Not authorised to approve.")
		}
		if !app.fungibleTokenKeeper.TokenExists(ctx, msg.Payload.Token.Symbol) {
			return types.ErrInvalidTokenSymbol(msg.Payload.Token.Symbol)
		}
		err := app.fungibleTokenKeeper.ValidateSignatures(ctx, msg)
		if err != nil {
			return err
		}
		if msg.Payload.Token.Status == fungible.ApproveToken {

			for _, val := range msg.Payload.Token.TokenFees {
				if !app.feeKeeper.FeeSettingExists(ctx, val.FeeName) {
					return types.ErrFeeSettingNotExists(val.FeeName)
				}

				if !fee.ContainFungibleAction(val.Action) {
					return types.ErrInvalidTokenAction()
				}
			}

			if app.fungibleTokenKeeper.CheckApprovedToken(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenAlreadyApproved(msg.Payload.Token.Symbol)
			}
		}
		if msg.Payload.Token.Status == fungible.FreezeToken {
			if app.fungibleTokenKeeper.IsTokenFrozen(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenFrozen()
			}
			if !app.fungibleTokenKeeper.CheckApprovedToken(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenInvalid()
			}
		}
		if msg.Payload.Token.Status == fungible.UnfreezeToken {
			if !app.fungibleTokenKeeper.IsTokenFrozen(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenUnFrozen()
			}
		}

		if msg.Payload.Token.Status == fungible.ApproveTransferTokenOwnership || msg.Payload.Token.Status == fungible.RejectTransferTokenOwnership {
			if app.fungibleTokenKeeper.IsTokenFrozen(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenFrozen()
			}

			if !app.fungibleTokenKeeper.IsVerifyableTransferTokenOwnership(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenTransferTokenOwnershipInvalid()
			}
		}
	case fungible.MsgSetFungibleTokenAccountStatus:
		if msg.TokenAccountPayload.TokenAccount.Status == fungible.FreezeTokenAccount {
			if app.fungibleTokenKeeper.IsFungibleTokenAccountFrozen(ctx, msg.TokenAccountPayload.TokenAccount.Account, msg.TokenAccountPayload.TokenAccount.Symbol) {
				return types.ErrTokenAccountFrozen()
			}
		}
		if msg.TokenAccountPayload.TokenAccount.Status == fungible.UnfreezeTokenAccount {
			if !app.fungibleTokenKeeper.IsFungibleTokenAccountFrozen(ctx, msg.TokenAccountPayload.TokenAccount.Account, msg.TokenAccountPayload.TokenAccount.Symbol) {
				return types.ErrTokenAccountUnFrozen()
			}
		}
	case fungible.MsgTransferFungibleToken:
		if !app.fungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
		if app.fungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}
		if app.fungibleTokenKeeper.IsFungibleTokenAccountFrozen(ctx, msg.From, msg.Symbol) {
			return types.ErrTokenAccountFrozen()
		}
		//check receiver
		if app.fungibleTokenKeeper.IsFungibleTokenAccountFrozen(ctx, msg.To, msg.Symbol) {
			return types.ErrTokenAccountFrozen()
		}
	case fungible.MsgMintFungibleToken:
		if !app.fungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}

		if app.fungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}

		// check sender
		if app.fungibleTokenKeeper.IsFungibleTokenAccountFrozen(ctx, msg.Owner, msg.Symbol) {
			return types.ErrTokenAccountFrozen()
		}

		//check receiver
		if app.fungibleTokenKeeper.IsFungibleTokenAccountFrozen(ctx, msg.To, msg.Symbol) {
			return types.ErrTokenAccountFrozen()
		}

		var fungibleToken = new(fungible.Token)
		app.fungibleTokenKeeper.GetFungibleTokenDataInfo(ctx, msg.Symbol, fungibleToken)
		fungibleToken.TotalSupply = fungibleToken.TotalSupply.Add(msg.Value)
		if !fungibleToken.MaxSupply.IsZero() {
			if fungibleToken.TotalSupply.GT(fungibleToken.MaxSupply) {
				return types.ErrInvalidTokenSupply()
			}
		}

		// (FixedSupply-FungibleToken) MintFlag - types.Bitmask = 0x0002
		if !fungibleToken.Flags.HasFlag(fungible.MintFlag) {
			return types.ErrInvalidTokenAction()
		}

	case fungible.MsgBurnFungibleToken:
		if !app.fungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
		if app.fungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}

		var fungibleToken = new(fungible.Token)
		app.fungibleTokenKeeper.GetFungibleTokenDataInfo(ctx, msg.Symbol, fungibleToken)
		fungibleToken.TotalSupply = fungibleToken.TotalSupply.Add(msg.Value)
		if !fungibleToken.Flags.HasFlag(fungible.BurnFlag) {
			return types.ErrInvalidTokenAction()
		}

		var account = new(fungible.FungibleTokenAccount)
		account = app.fungibleTokenKeeper.GetFungibleAccount(ctx, msg.Symbol, msg.From)

		if account.Frozen {
			return types.ErrTokenAccountFrozen()
		}

		if account.Balance.LT(msg.Value) {
			return types.ErrInvalidTokenAccountBalance(fmt.Sprintf("Not enough tokens. Have only %v", account.Balance.String()))
		}

	case fungible.MsgTransferFungibleTokenOwnership:
		if app.fungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}

		if !app.fungibleTokenKeeper.IsTokenOwnershipTransferrable(ctx, msg.Symbol) {
			return types.ErrInvalidTokenAction()
		}
		if !app.fungibleTokenKeeper.IsTokenOwner(ctx, msg.Symbol, msg.From) {
			return types.ErrInvalidTokenOwner()
		}
		if !app.fungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
	case fungible.MsgAcceptFungibleTokenOwnership:

		if app.fungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}

		if !app.fungibleTokenKeeper.IsTokenOwnershipAcceptable(ctx, msg.Symbol) {
			return types.ErrInvalidTokenAction()
		}

		if !app.fungibleTokenKeeper.IsTokenNewOwner(ctx, msg.Symbol, msg.From) {
			return types.ErrInvalidTokenNewOwner()
		}

		if !app.fungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
	case nonFungible.MsgCreateNonFungibleToken:
		ownerAcc := utils.GetAccount(ctx, app.accountKeeper, msg.Owner)
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
		if !app.nonFungibleTokenKeeper.TokenExists(ctx, msg.Payload.Token.Symbol) {
			return types.ErrInvalidTokenSymbol(msg.Payload.Token.Symbol)
		}
		if !app.nonFungibleTokenKeeper.IsAuthorised(ctx, msg.GetSigners()[0]) {
			return sdkTypes.ErrUnauthorized("Not authorised to approve.")
		}
		err := app.nonFungibleTokenKeeper.ValidateSignatures(ctx, msg)
		if err != nil {
			return err
		}
		if msg.Payload.Token.Status == nonFungible.ApproveToken {

			for _, val := range msg.Payload.Token.TokenFees {
				if !app.feeKeeper.FeeSettingExists(ctx, val.FeeName) {
					return types.ErrFeeSettingNotExists(val.FeeName)
				}

				if !fee.ContainNonFungibleAction(val.Action) {
					return types.ErrInvalidTokenAction()
				}
			}

			if app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenAlreadyApproved(msg.Payload.Token.Symbol)
			}
		}

		if msg.Payload.Token.Status == nonFungible.RejectToken {
			if app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenAlreadyApproved(msg.Payload.Token.Symbol)
			}
		}

		if msg.Payload.Token.Status == nonFungible.FreezeToken {
			if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenFrozen()
			}
			if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenInvalid()
			}
		}
		if msg.Payload.Token.Status == nonFungible.UnfreezeToken {
			if !app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenUnFrozen()
			}
		}

		if msg.Payload.Token.Status == nonFungible.ApproveTransferTokenOwnership || msg.Payload.Token.Status == nonFungible.RejectTransferTokenOwnership {
			if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenFrozen()
			}

			if !app.nonFungibleTokenKeeper.IsVerifyableTransferTokenOwnership(ctx, msg.Payload.Token.Symbol) {
				return types.ErrTokenTransferTokenOwnershipInvalid()
			}
		}
	case nonFungible.MsgSetNonFungibleItemStatus:
		if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.ItemPayload.Item.Symbol) {
			return types.ErrTokenInvalid()
		}

		if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.ItemPayload.Item.Symbol) {
			return types.ErrTokenFrozen()
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

		if !app.kycKeeper.IsWhitelisted(ctx, msg.To) {
			return types.ErrReceiverNotKyc()
		}

		if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
		if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}

		if !app.nonFungibleTokenKeeper.IsItemOwner(ctx, msg.Symbol, msg.ItemID, msg.From) {
			return types.ErrTokenItemNotFound()
		}
		if app.nonFungibleTokenKeeper.IsNonFungibleItemFrozen(ctx, msg.Symbol, msg.ItemID) {
			return types.ErrTokenItemFrozen()
		}

		// 1. [Transfer non fungible token item - Invalid Item-ID]
		nonFungibleItem := app.nonFungibleTokenKeeper.GetNonFungibleItem(ctx, msg.Symbol, msg.ItemID)
		if nonFungibleItem == nil {
			return sdkTypes.ErrUnknownRequest("Invalid Item ID.")
		}

		if app.nonFungibleTokenKeeper.IsItemTransferLimitExceeded(ctx, msg.Symbol, msg.ItemID) {
			return sdkTypes.ErrInternal("Transfer limit exceeded.")
		}

	case nonFungible.MsgMintNonFungibleItem:
		if !app.kycKeeper.IsWhitelisted(ctx, msg.To) {
			return types.ErrReceiverNotKyc()
		}

		if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
		if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}

		if app.nonFungibleTokenKeeper.IsMintLimitExceeded(ctx, msg.Symbol, msg.To) {
			return sdkTypes.ErrInternal("Mint limit exceeded.")
		}

		//1. checking: (flag of Public equals to TRUE)
		var token = new(nonFungible.Token)
		app.nonFungibleTokenKeeper.GetNonfungibleTokenDataInfo(ctx, msg.Symbol, token)
		if token.Flags.HasFlag(nonFungible.PubFlag) {
			ownerAcc := msg.Owner
			newOwnerAcc := msg.To
			if !ownerAcc.Equals(newOwnerAcc) {
				return sdkTypes.ErrInternal("Public token can only be minted to oneself.")
			}
		}

		//2. Check for Unique status
		if !app.nonFungibleTokenKeeper.IsItemIDUnique(ctx, msg.Symbol, msg.ItemID) {
			return types.ErrTokenItemIDInUsed()
		}

	case nonFungible.MsgBurnNonFungibleItem:
		if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
		if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}

		var token = new(nonFungible.Token)
		app.nonFungibleTokenKeeper.GetNonfungibleTokenDataInfo(ctx, msg.Symbol, token)
		if !token.Flags.HasFlag(nonFungible.BurnFlag) {
			return types.ErrInvalidTokenAction()
		}

		// 1. [Burn non fungible token item - Invalid Item-owner]
		item := app.nonFungibleTokenKeeper.GetNonFungibleItem(ctx, msg.Symbol, msg.ItemID)
		if item == nil {
			return types.ErrTokenItemNotFound()
		}

		if item.Frozen {
			return types.ErrTokenItemFrozen()
		}

		if !app.nonFungibleTokenKeeper.IsItemOwner(ctx, msg.Symbol, msg.ItemID, msg.From) {
			return types.ErrInvalidItemOwner()
		}

	case nonFungible.MsgTransferNonFungibleTokenOwnership:

		if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}

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

		if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}

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

	case nonFungible.MsgUpdateNFTMetadata:
		if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
		if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}
		if !app.nonFungibleTokenKeeper.IsTokenOwner(ctx, msg.Symbol, msg.From) {
			return types.ErrInvalidTokenOwner()
		}
	case nonFungible.MsgUpdateItemMetadata:
		if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
		if app.nonFungibleTokenKeeper.IsNonFungibleItemFrozen(ctx, msg.Symbol, msg.ItemID) {
			return types.ErrTokenAccountFrozen()
		}
		if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}
		if !app.nonFungibleTokenKeeper.IsItemMetadataModifiable(ctx, msg.Symbol, msg.From, msg.ItemID) {
			return sdkTypes.ErrInternal("Non fungible item metadata is not modifiable.")
		}
	case nonFungible.MsgUpdateEndorserList:
		if !app.nonFungibleTokenKeeper.CheckApprovedToken(ctx, msg.Symbol) {
			return types.ErrTokenInvalid()
		}
		if app.nonFungibleTokenKeeper.IsTokenFrozen(ctx, msg.Symbol) {
			return types.ErrTokenFrozen()
		}
		if !app.nonFungibleTokenKeeper.IsTokenOwner(ctx, msg.Symbol, msg.From) {
			return types.ErrInvalidTokenOwner()
		}
		for _, v := range msg.Endorsers {
			if !app.kycKeeper.IsWhitelisted(ctx, v) {
				return types.ErrUnauthorisedEndorser()
			}
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
	case auth.MsgCreateMultiSigAccount:
		for _, signer := range msg.Signers {
			if !app.kycKeeper.IsWhitelisted(ctx, signer) {
				return sdkTypes.ErrUnknownRequest("Signer is not whitelisted.")
			}
		}
	case auth.MsgCreateMultiSigTx:
		groupAcc := utils.GetAccount(ctx, app.accountKeeper, msg.GroupAddress)
		if groupAcc == nil {
			return sdkTypes.ErrInternal("Group account not found.")
		}
		if !groupAcc.IsMultiSig() {
			return sdkTypes.ErrInternal("MultiSig Tx create failed, group address is not a multisig account.")
		}
		if !groupAcc.IsSigner(msg.Sender) {
			return sdkTypes.ErrInternal("Sender is not group account's signer.")
		}
		_, exist := groupAcc.GetMultiSig().CheckTx(msg.StdTx)
		if exist {
			return sdkTypes.ErrInternal("Tx already existed in pending tx.")
		}
		internalMsgErr := msg.StdTx.GetMsgs()[0].ValidateBasic()
		if internalMsgErr != nil {
			return sdkTypes.ErrInternal("Internal transaction invalid.")
		}
		_, sigErr := utils.CheckTxSig(ctx, msg.StdTx, app.accountKeeper, app.kycKeeper)
		if sigErr != nil {
			return sdkTypes.ConvertError(sigErr)
		}
	case auth.MsgUpdateMultiSigAccount:
		groupAcc := utils.GetAccount(ctx, app.accountKeeper, msg.GroupAddress)
		if groupAcc == nil {
			return sdkTypes.ErrInternal("Group account not found.")
		}
		if !groupAcc.IsMultiSig() {
			return sdkTypes.ErrInternal("MultiSig Tx update failed, group address is not a multisig account.")
		}
		if !groupAcc.GetMultiSig().GetOwner().Equals(msg.Owner) {
			return sdkTypes.ErrUnknownRequest("Owner address invalid.")
		}
	case auth.MsgTransferMultiSigOwner:
		groupAcc := utils.GetAccount(ctx, app.accountKeeper, msg.GroupAddress)
		if groupAcc == nil {
			return sdkTypes.ErrInternal("Group account not found.")
		}
		if !groupAcc.IsMultiSig() {
			return sdkTypes.ErrInternal("MultiSig Tx update failed, group address is not a multisig account.")
		}
		if !groupAcc.GetMultiSig().IsOwner(msg.Owner) {
			return sdkTypes.ErrUnknownRequest("Owner of group address invalid.")
		}
	case auth.MsgDeleteMultiSigTx:
		groupAcc := utils.GetAccount(ctx, app.accountKeeper, msg.GroupAddress)
		if groupAcc == nil {
			return sdkTypes.ErrInternal("Group account not found.")
		}
		if !groupAcc.IsMultiSig() {
			return sdkTypes.ErrInternal("MultiSig Tx update failed, group address is not a multisig account.")
		}
		if !groupAcc.GetMultiSig().IsOwner(msg.Sender) {
			return sdkTypes.ErrUnknownRequest("Only group account owner can remove pending tx.")
		}
		pendingTx := groupAcc.GetMultiSig().GetPendingTx(msg.TxID)
		if pendingTx == nil {
			return sdkTypes.ErrInternal("MultiSig pending tx not found.")
		}
	case auth.MsgSignMultiSigTx:
		groupAcc := utils.GetAccount(ctx, app.accountKeeper, msg.GroupAddress)
		if groupAcc == nil {
			return sdkTypes.ErrInternal("Group account not found.")
		}
		if !groupAcc.IsMultiSig() {
			return sdkTypes.ErrInternal("MultiSig Tx sign failed, group address is not a multisig account.")
		}
		if !groupAcc.IsSigner(msg.Sender) {
			return sdkTypes.ErrInternal("Sender is not group account's signer.")
		}
		pendingTx := groupAcc.GetMultiSig().GetPendingTx(msg.TxID)
		if pendingTx == nil {
			return sdkTypes.ErrInternal("MultiSig pending tx not found.")
		}
		if pendingTx.IsSignedBy(msg.Sender) {
			return sdkTypes.ErrInternal("Signer already signed this tx.")
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
