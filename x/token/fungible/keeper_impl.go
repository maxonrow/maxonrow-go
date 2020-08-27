package fungible

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/types"
)

func (k *Keeper) MintFungibleToken(ctx sdkTypes.Context, symbol string, from sdkTypes.AccAddress, to sdkTypes.AccAddress, value sdkTypes.Uint) sdkTypes.Result {

	var token = new(Token)
	if exists := k.GetFungibleTokenDataInfo(ctx, symbol, token); !exists {
		return types.ErrInvalidTokenSymbol(symbol).Result()
	}

	if !token.Flags.HasFlag(MintFlag) {
		return types.ErrInvalidTokenAction().Result()
	}

	// Wallet account
	minterAccount := k.accountKeeper.GetAccount(ctx, from)
	if minterAccount == nil {
		return types.ErrInvalidTokenMinter().Result()
	}

	// minter can only be the owner of the token
	tokenOwnerAccount := k.GetFungibleAccount(ctx, symbol, from)

	// token account
	if tokenOwnerAccount == nil {
		return types.ErrInvalidTokenMinter().Result()
	}

	if tokenOwnerAccount.Frozen {
		return types.ErrTokenAccountFrozen().Result()
	}

	if !token.Owner.Equals(from) {
		return types.ErrInvalidTokenMinter().Result()
	}

	if !token.Flags.HasFlag(ApprovedFlag) {
		return types.ErrTokenInvalid().Result()
	}

	if token.Flags.HasFlag(FrozenFlag) {
		return types.ErrTokenFrozen().Result()
	}

	account := k.GetFungibleAccount(ctx, symbol, to)
	if account == nil {
		account = k.createFungibleAccount(ctx, symbol, to)
	}

	if account.Frozen {
		return types.ErrTokenAccountFrozen().Result()
	}

	token.TotalSupply = token.TotalSupply.Add(value)

	// max supply 0 means is dynamic supply
	if !token.MaxSupply.IsZero() {
		if token.TotalSupply.GT(token.MaxSupply) {
			return types.ErrInvalidTokenSupply().Result()
		}
	}

	addFungibleTokenErr := k.addFungibleToken(ctx, symbol, to, value)
	if addFungibleTokenErr != nil {
		return addFungibleTokenErr.Result()
	}

	k.storeToken(ctx, symbol, token)

	eventParam := []string{symbol, "mxw1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqgcpfl3", to.String(), value.String()}
	eventSignature := "MintedFungibleToken(string,string,string,bignumber)"

	accountSequence := minterAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}
}

//* TransferFungibleToken
func (k *Keeper) TransferFungibleToken(ctx sdkTypes.Context, symbol string, from, to sdkTypes.AccAddress, value sdkTypes.Uint) sdkTypes.Result {
	var token = new(Token)
	if exists := k.GetFungibleTokenDataInfo(ctx, symbol, token); !exists {
		return types.ErrTokenInvalid().Result()
	}

	fromAccount := k.accountKeeper.GetAccount(ctx, from)
	if fromAccount == nil {
		return types.ErrInvalidTokenAccount().Result()
	}

	if token.Flags.HasFlag(FrozenFlag) {
		return types.ErrTokenFrozen().Result()
	}

	ownerAccount := k.GetFungibleAccount(ctx, symbol, from)
	if ownerAccount == nil {
		return sdkTypes.ErrUnknownRequest("Owner doesn't have such token.").Result()
	}

	if ownerAccount.Frozen {
		return types.ErrTokenAccountFrozen().Result()
	}

	if ownerAccount.Balance.LT(value) {
		return types.ErrInvalidTokenAccountBalance(fmt.Sprintf("Not enough tokens. Have only %v.", ownerAccount.Balance.String())).Result()
	}

	newOwnerAccount := k.GetFungibleAccount(ctx, symbol, to)
	if newOwnerAccount == nil {
		newOwnerAccount = k.createFungibleAccount(ctx, symbol, to)
	}

	if newOwnerAccount.Frozen {
		return types.ErrTokenAccountFrozen().Result()
	}

	subFungibleTokenErr := k.subFungibleToken(ctx, symbol, from, value)
	if subFungibleTokenErr != nil {
		return subFungibleTokenErr.Result()
	}

	addFungibleTokenErr := k.addFungibleToken(ctx, symbol, to, value)
	if addFungibleTokenErr != nil {
		return addFungibleTokenErr.Result()
	}

	eventParam := []string{symbol, from.String(), to.String(), value.String()}
	eventSignature := "TransferredFungibleToken(string,string,string,bignumber)"

	accountSequence := fromAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}

}

// BurnFungibleToken
func (k *Keeper) BurnFungibleToken(ctx sdkTypes.Context, symbol string, owner sdkTypes.AccAddress, value sdkTypes.Uint) sdkTypes.Result {
	var token = new(Token)
	if exists := k.GetFungibleTokenDataInfo(ctx, symbol, token); !exists {
		return types.ErrInvalidTokenSymbol(symbol).Result()
	}

	if !token.Flags.HasFlag(BurnFlag) {
		return types.ErrInvalidTokenAction().Result()
	}

	ownerAccount := k.accountKeeper.GetAccount(ctx, owner)
	if ownerAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid account to burn from.").Result()
	}

	if !token.Flags.HasFlag(ApprovedFlag) {
		return types.ErrTokenInvalid().Result()
	}

	if token.Flags.HasFlag(FrozenFlag) {
		return types.ErrTokenFrozen().Result()
	}

	account := k.GetFungibleAccount(ctx, symbol, owner)
	if account == nil {
		return types.ErrInvalidTokenAccount().Result()
	}

	if account.Frozen {
		return types.ErrTokenAccountFrozen().Result()
	}

	if account.Balance.LT(value) {
		return types.ErrInvalidTokenAccountBalance(fmt.Sprintf("Not enough tokens. Have only %v.", account.Balance.String())).Result()
	}

	token.TotalSupply = token.TotalSupply.Sub(value)
	k.storeToken(ctx, symbol, token)

	subFungibleTokenErr := k.subFungibleToken(ctx, symbol, owner, value)
	if subFungibleTokenErr != nil {
		return subFungibleTokenErr.Result()
	}

	eventParam := []string{symbol, owner.String(), "mxw1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqgcpfl3", value.String()}
	eventSignature := "BurnedFungibleToken(string,string,string,bignumber)"

	accountSequence := ownerAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, owner.String(), eventParam),
		Log:    resultLog.String(),
	}

}

// TO-DO: proper implementation to cater nonfungibletoken transfer ownership
func (k *Keeper) TransferTokenOwnership(ctx sdkTypes.Context, symbol string, from sdkTypes.AccAddress, to sdkTypes.AccAddress, metadata string) sdkTypes.Result {

	var token = new(Token)

	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return err.Result()
	}

	return (k.transferFungibleTokenOwnership(ctx, from, to, token, metadata))

}

func (k *Keeper) transferFungibleTokenOwnership(ctx sdkTypes.Context, from sdkTypes.AccAddress, to sdkTypes.AccAddress, token *Token, metadata string) sdkTypes.Result {

	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, token.Owner)
	if ownerWalletAccount == nil {
		return types.ErrInvalidTokenOwner().Result()
	}

	if ownerWalletAccount != nil && !token.Owner.Equals(from) {
		return types.ErrInvalidTokenOwner().Result()
	}

	if !token.IsApproved() {
		// TODO: Please define an error code
		return sdkTypes.ErrUnknownRequest("Token is not approved.").Result()
	}

	if token.IsFrozen() {
		return types.ErrTokenFrozen().Result()
	}

	newOwnerAccount := k.GetFungibleAccount(ctx, token.Symbol, to)
	if newOwnerAccount == nil {
		newOwnerAccount = k.createFungibleAccount(ctx, token.Symbol, to)
	}

	if newOwnerAccount.Frozen {
		return sdkTypes.ErrUnknownRequest("New owner is frozen.").Result()
	}

	// set token newowner to new owner, pending for accepting by new owner
	token.NewOwner = to
	token.Metadata = metadata
	token.Flags.AddFlag(TransferTokenOwnershipFlag)

	k.storeToken(ctx, token.Symbol, token)

	eventParam := []string{token.Symbol, from.String(), to.String()}
	eventSignature := "TransferredFungibleTokenOwnership(string,string,string)"

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}
}

// TO-DO: proper implementation to cater nonfungibletoken accept ownership
func (k *Keeper) AcceptTokenOwnership(ctx sdkTypes.Context, symbol string, from sdkTypes.AccAddress, metadata string) sdkTypes.Result {

	var token = new(Token)

	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return err.Result()
	}

	return (k.acceptFungibleTokenOwnership(ctx, from, token, metadata))
}

func (k *Keeper) acceptFungibleTokenOwnership(ctx sdkTypes.Context, from sdkTypes.AccAddress, token *Token, metadata string) sdkTypes.Result {

	if !token.Flags.HasFlag(AcceptTokenOwnershipFlag) && !token.Flags.HasFlag(ApproveTransferTokenOwnershipFlag) && !token.Flags.HasFlag(TransferTokenOwnershipFlag) {
		return types.ErrInvalidTokenAction().Result()
	}

	// validation of exisisting owner account
	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, token.Owner)
	if ownerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid token owner.").Result()
	}

	// validation of new owner account
	newOwnerWalletAccount := k.accountKeeper.GetAccount(ctx, token.Owner)
	if newOwnerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid token owner.").Result()
	}

	newOwnerAccount := k.GetFungibleAccount(ctx, token.Symbol, from)
	if newOwnerAccount == nil {
		return sdkTypes.ErrUnknownRequest("New owner account is not found.").Result()
	}

	if newOwnerAccount.Frozen {
		return sdkTypes.ErrUnknownRequest("New owner account is frozen.").Result()
	}

	if newOwnerWalletAccount != nil && token.NewOwner.String() != from.String() {
		return types.ErrInvalidTokenNewOwner().Result()
	}

	if !token.Flags.HasFlag(ApprovedFlag) {
		return sdkTypes.ErrUnknownRequest("Fungible token is not approved.").Result()
	}

	if token.Flags.HasFlag(FrozenFlag) {
		return types.ErrTokenFrozen().Result()
	}

	//TO-DO: if there is need to set token.NewOwner to empty
	// accepting token ownership, remove newowner move the newowner into owner.
	var emptyAccAddr sdkTypes.AccAddress
	token.Owner = from
	token.NewOwner = emptyAccAddr
	token.Metadata = metadata

	token.Flags.RemoveFlag(ApproveTransferTokenOwnershipFlag)
	token.Flags.RemoveFlag(AcceptTokenOwnershipFlag)
	token.Flags.RemoveFlag(TransferTokenOwnershipFlag)
	k.storeToken(ctx, token.Symbol, token)

	eventParam := []string{token.Symbol, from.String()}
	eventSignature := "AcceptedFungibleTokenOwnership(string,string)"

	accountSequence := newOwnerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}

}

func (k *Keeper) IsTokenOwnershipAcceptable(ctx sdkTypes.Context, symbol string) bool {

	var token = new(Token)

	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return false
	}

	if token.Flags.HasFlag(TransferTokenOwnershipFlag) &&
		token.Flags.HasFlag(AcceptTokenOwnershipFlag) &&
		token.Flags.HasFlag(ApproveTransferTokenOwnershipFlag) {
		return true
	}

	return false
}

func (k *Keeper) IsTokenOwnershipTransferrable(ctx sdkTypes.Context, symbol string) bool {

	var token = new(Token)

	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return false
	}

	if token.Flags.HasFlag(TransferTokenOwnershipFlag) || token.Flags.HasFlag(AcceptTokenOwnershipFlag) || token.Flags.HasFlag(ApproveTransferTokenOwnershipFlag) {
		return false
	}

	return true
}

func (k *Keeper) IsTokenOwner(ctx sdkTypes.Context, symbol string, owner sdkTypes.AccAddress) bool {
	var token = new(Token)

	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return false
	}

	if token.Owner.Equals(owner) {
		return true
	}

	return false
}

func (k *Keeper) IsTokenNewOwner(ctx sdkTypes.Context, symbol string, newOwner sdkTypes.AccAddress) bool {

	var token = new(Token)

	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return false
	}

	if token.NewOwner.Equals(newOwner) {
		return true
	}

	return false
}
