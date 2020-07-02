package fungible

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/bank"
	"github.com/maxonrow/maxonrow-go/x/fee"
)

type Keeper struct {
	accountKeeper *sdkAuth.AccountKeeper
	feeKeeper     *fee.Keeper
	key           sdkTypes.StoreKey
	cdc           *codec.Codec
}

const (
	FungibleFlag types.Bitmask = 0x0001
	MintFlag     types.Bitmask = 0x0002
	BurnFlag     types.Bitmask = 0x0004
	FrozenFlag   types.Bitmask = 0x0008
	ApprovedFlag types.Bitmask = 0x0010

	TransferTokenOwnershipFlag        types.Bitmask = 0x0100
	ApproveTransferTokenOwnershipFlag types.Bitmask = 0x0200
	AcceptTokenOwnershipFlag          types.Bitmask = 0x0400

	DynamicFungibleTokenMask                = FungibleFlag + MintFlag + BurnFlag
	FixedSupplyBurnableFungibleTokenMask    = FungibleFlag + BurnFlag
	FixedSupplyNotBurnableFungibleTokenMask = FungibleFlag
)

type Token struct {
	Flags       types.Bitmask
	Name        string
	Symbol      string
	Decimals    int
	Owner       sdkTypes.AccAddress
	NewOwner    sdkTypes.AccAddress
	Metadata    string
	TotalSupply sdkTypes.Uint
	MaxSupply   sdkTypes.Uint
}

type FungibleTokenAccount struct {
	Owner    sdkTypes.AccAddress
	Frozen   bool
	Metadata string
	Balance  sdkTypes.Uint
}

func (t *Token) IsApproved() bool {
	return t.Flags.HasFlag(ApprovedFlag)
}

func (t *Token) IsFrozen() bool {
	return t.Flags.HasFlag(FrozenFlag)
}

func NewKeeper(cdc *codec.Codec, accountKeeper *auth.AccountKeeper, feeKeeper *fee.Keeper, key sdkTypes.StoreKey) Keeper {
	return Keeper{
		cdc:           cdc,
		key:           key,
		accountKeeper: accountKeeper,
		feeKeeper:     feeKeeper,
	}
}

func (k *Keeper) SetAuthorisedAddresses(ctx sdkTypes.Context, addresses []sdkTypes.AccAddress) {
	authorisedStore := ctx.KVStore(k.key)
	key := getAuthorisedKey()

	ah := k.GetAuthorisedAddresses(ctx)
	ah.AppendAccAddrs(addresses)
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(ah)
	if err != nil {
		panic(err)
	}

	authorisedStore.Set(key, bz)
}

func (k *Keeper) GetAuthorisedAddresses(ctx sdkTypes.Context) types.AddressHolder {
	var ah types.AddressHolder
	authorisedStore := ctx.KVStore(k.key)
	key := getAuthorisedKey()

	bz := authorisedStore.Get(key)
	if bz == nil {
		return ah
	}
	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &ah)
	if err != nil {
		panic(err)
	}
	return ah
}

func (k *Keeper) RemoveAuthorisedAddresses(ctx sdkTypes.Context, addresses []sdkTypes.AccAddress) {
	authorisedStore := ctx.KVStore(k.key)
	key := getAuthorisedKey()
	authorisedAddresses := k.GetAuthorisedAddresses(ctx)

	for _, authorisedAddress := range addresses {
		authorisedAddresses.Remove(authorisedAddress)
	}

	bz, err := k.cdc.MarshalBinaryLengthPrefixed(authorisedAddresses)
	if err != nil {
		panic(err)
	}

	authorisedStore.Set(key, bz)
}

//IsAuthorised Check if is authorised
func (k Keeper) IsAuthorised(ctx sdkTypes.Context, address sdkTypes.AccAddress) bool {

	ah := k.GetAuthorisedAddresses(ctx)
	_, ok := ah.Contains(address)

	return ok
}

func (k *Keeper) GetIssuerAddresses(ctx sdkTypes.Context) types.AddressHolder {
	var ah types.AddressHolder
	authorisedStore := ctx.KVStore(k.key)
	key := getIssuerKey()

	bz := authorisedStore.Get(key)
	if bz == nil {
		return ah
	}
	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &ah)
	if err != nil {
		panic(err)
	}
	return ah
}

func (k Keeper) SetIssuerAddresses(ctx sdkTypes.Context, addresses []sdkTypes.AccAddress) {

	authorisedStore := ctx.KVStore(k.key)
	key := getIssuerKey()

	ah := k.GetIssuerAddresses(ctx)
	ah.AppendAccAddrs(addresses)
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(ah)
	if err != nil {
		panic(err)
	}

	authorisedStore.Set(key, bz)
}

func (k *Keeper) RemoveIssuerAddresses(ctx sdkTypes.Context, addresses []sdkTypes.AccAddress) {
	authorisedStore := ctx.KVStore(k.key)
	key := getIssuerKey()
	issuerAddresses := k.GetIssuerAddresses(ctx)

	for _, issuerAddress := range addresses {
		issuerAddresses.Remove(issuerAddress)
	}

	bz, err := k.cdc.MarshalBinaryLengthPrefixed(issuerAddresses)
	if err != nil {
		panic(err)
	}

	authorisedStore.Set(key, bz)

}

func (k Keeper) IsIssuer(ctx sdkTypes.Context, address sdkTypes.AccAddress) bool {
	ah := k.GetIssuerAddresses(ctx)
	_, ok := ah.Contains(address)

	return ok
}

func (k *Keeper) GetProviderAddresses(ctx sdkTypes.Context) types.AddressHolder {
	var ah types.AddressHolder
	authorisedStore := ctx.KVStore(k.key)
	key := getProviderKey()

	bz := authorisedStore.Get(key)
	if bz == nil {
		return ah
	}
	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &ah)
	if err != nil {
		panic(err)
	}
	return ah
}

func (k Keeper) SetProviderAddresses(ctx sdkTypes.Context, addresses []sdkTypes.AccAddress) {

	authorisedStore := ctx.KVStore(k.key)
	key := getProviderKey()

	ah := k.GetProviderAddresses(ctx)
	ah.AppendAccAddrs(addresses)
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(ah)
	if err != nil {
		panic(err)
	}

	authorisedStore.Set(key, bz)
}

func (k *Keeper) RemoveProviderAddresses(ctx sdkTypes.Context, addresses []sdkTypes.AccAddress) {
	authorisedStore := ctx.KVStore(k.key)
	key := getProviderKey()
	providerAddresses := k.GetProviderAddresses(ctx)

	for _, providerAddress := range addresses {
		providerAddresses.Remove(providerAddress)
	}

	bz, err := k.cdc.MarshalBinaryLengthPrefixed(providerAddresses)
	if err != nil {
		panic(err)
	}

	authorisedStore.Set(key, bz)

}

func (k Keeper) IsProvider(ctx sdkTypes.Context, address sdkTypes.AccAddress) bool {
	ah := k.GetProviderAddresses(ctx)
	_, ok := ah.Contains(address)

	return ok
}

func (k *Keeper) TokenExists(ctx sdkTypes.Context, symbol string) bool {
	store := ctx.KVStore(k.key)
	key := getTokenKey(symbol)
	return store.Has([]byte(key))
}

func (k *Keeper) CreateFungibleToken(
	ctx sdkTypes.Context,
	name string,
	symbol string,
	decimals int,
	owner sdkTypes.AccAddress,
	fixedSupply bool,
	maxSupply sdkTypes.Uint,
	metadata string,
	fee Fee,
) sdkTypes.Result {

	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, owner)

	if ownerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Not authorised to apply for token creation.").Result()
	}

	if k.TokenExists(ctx, symbol) {
		return types.ErrTokenExists(symbol).Result()
	}

	// Overwrite the cosmos sdk tags.
	// Application fee paid at ante.go, event created here.
	amt, parseErr := sdkTypes.ParseCoins(fee.Value + types.CIN)
	if parseErr != nil {
		return sdkTypes.ErrInvalidCoins("Parse value to coins failed.").Result()
	}
	applicationFeeResult := bank.MakeBankSendEvent(ctx, owner, fee.To, amt, *k.accountKeeper)

	zero := sdkTypes.NewUintFromString("0")

	token := &Token{
		Name:        name,
		Flags:       FungibleFlag,
		Symbol:      symbol,
		Decimals:    decimals,
		Owner:       owner,
		Metadata:    metadata,
		MaxSupply:   maxSupply,
		TotalSupply: zero,
	}

	if !fixedSupply {
		token.Flags.AddFlag(MintFlag)
	}

	k.storeToken(ctx, symbol, token)

	eventParam := []string{symbol, owner.String(), fee.To.String(), fee.Value}
	eventSignature := "CreatedFungibleToken(string,string,string,bignumber)"
	event := types.MakeMxwEvents(eventSignature, owner.String(), eventParam)

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: applicationFeeResult.Events.AppendEvents(event),
		Log:    resultLog.String(),
	}

}

// ApproveToken
func (k *Keeper) ApproveToken(ctx sdkTypes.Context, symbol string, tokenFees []TokenFee, burnable bool, signer sdkTypes.AccAddress, metadata string) sdkTypes.Result {
	if !k.IsAuthorised(ctx, signer) {
		return sdkTypes.ErrUnauthorized("Not authorised to approve.").Result()
	}

	return k.approveFungibleToken(ctx, symbol, tokenFees, burnable, metadata, signer)
}

func (k *Keeper) approveFungibleToken(ctx sdkTypes.Context, symbol string, tokenFees []TokenFee, burnable bool, metadata string, signer sdkTypes.AccAddress) sdkTypes.Result {
	var token = new(Token)
	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return err.Result()
	}

	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, token.Owner)
	if ownerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid token owner.").Result()
	}

	if token.Flags.HasFlag(ApprovedFlag) {
		return types.ErrTokenAlreadyApproved(symbol).Result()
	}

	// Assign fee to token
	for _, tokenFee := range tokenFees {
		if !k.feeKeeper.FeeSettingExists(ctx, tokenFee.FeeName) {
			return types.ErrFeeSettingNotExists(tokenFee.FeeName).Result()
		}
		err := k.feeKeeper.AssignFeeToTokenAction(ctx, tokenFee.FeeName, token.Symbol, tokenFee.Action)
		if err != nil {
			return err.Result()
		}
	}

	var flags types.Bitmask

	// if fixed supply: fixed supply doesn't allow to mint.
	if !token.Flags.HasFlag(MintFlag) {
		if burnable {
			flags = FixedSupplyBurnableFungibleTokenMask
		} else {
			flags = FixedSupplyNotBurnableFungibleTokenMask
		}
	} else {
		flags = DynamicFungibleTokenMask
	}

	token.Flags = flags + ApprovedFlag
	token.Metadata = metadata

	account := k.getFungibleAccount(ctx, symbol, token.Owner)
	if account == nil {
		account = k.createFungibleAccount(ctx, symbol, token.Owner)
	}

	if account.Frozen {
		return sdkTypes.ErrInternal("Fungible token account is frozen.").Result()
	}

	// if fixed supply: fixed supply doesn't allow to mint.
	if !token.Flags.HasFlag(MintFlag) {
		addFungibleTokenErr := k.addFungibleToken(ctx, symbol, token.Owner, token.MaxSupply)
		if addFungibleTokenErr != nil {
			return addFungibleTokenErr.Result()
		}

		token.TotalSupply = token.TotalSupply.Add(token.MaxSupply)
	}

	k.storeToken(ctx, symbol, token)

	// Get the token account again.
	account = k.getFungibleAccount(ctx, symbol, token.Owner)

	var transferEvents sdkTypes.Events
	if !account.Balance.IsZero() {
		// Event: After approve, added total supply into token owner account.
		transferEventParam := []string{symbol, "mxw000000000000000000000000000000000000000", ownerWalletAccount.GetAddress().String(), token.TotalSupply.String()}
		transferEventSignature := "TransferredFungibleToken(string,string,string,bignumber)"
		transferEvents = types.MakeMxwEvents(transferEventSignature, "mxw000000000000000000000000000000000000000", transferEventParam)
	}
	// Event: Approved fungible token
	eventParam := []string{symbol, token.Owner.String()}
	eventSignature := "ApprovedFungibleToken(string,string)"
	events := types.MakeMxwEvents(eventSignature, signer.String(), eventParam)

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: events.AppendEvents(transferEvents),
		Log:    resultLog.String(),
	}

}

func (k *Keeper) RejectToken(ctx sdkTypes.Context, symbol string, signer sdkTypes.AccAddress) sdkTypes.Result {

	var token = new(Token)

	if !k.IsAuthorised(ctx, signer) {
		return sdkTypes.ErrUnauthorized("Not authorised to reject").Result()
	}

	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return err.Result()
	}

	var isApproved = token.Flags.HasFlag(ApprovedFlag)
	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, token.Owner)
	if ownerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid token owner.").Result()
	}

	if isApproved {
		return types.ErrTokenAlreadyApproved(symbol).Result()
	}

	store := ctx.KVStore(k.key)

	tokenTypeKey := getTokenKey(symbol)
	store.Delete([]byte(tokenTypeKey))

	eventParam := []string{symbol, token.Owner.String()}
	eventSignature := "RejectedFungibleToken(string,string)"

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, signer.String(), eventParam),
		Log:    resultLog.String(),
	}

}

func (k *Keeper) FreezeToken(ctx sdkTypes.Context, symbol string, signer sdkTypes.AccAddress, metadata string) sdkTypes.Result {
	if !k.IsAuthorised(ctx, signer) {
		return sdkTypes.ErrUnauthorized("Not authorised to freeze.").Result()
	}

	var token = new(Token)
	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return err.Result()
	}

	return k.freezeFungibleToken(ctx, symbol, signer, metadata)

}

func (k *Keeper) freezeFungibleToken(ctx sdkTypes.Context, symbol string, signer sdkTypes.AccAddress, metadata string) sdkTypes.Result {
	var token = new(Token)
	k.mustGetTokenData(ctx, symbol, token)

	signerAccount := k.accountKeeper.GetAccount(ctx, signer)
	if signerAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid signer.").Result()
	}

	if token.Flags.HasFlag(FrozenFlag) {
		return types.ErrTokenFrozen().Result()
	}

	token.Flags.AddFlag(FrozenFlag)
	token.Metadata = metadata
	k.storeToken(ctx, symbol, token)

	eventParam := []string{symbol, token.Owner.String()}
	eventSignature := "FrozenFungibleToken(string,string)"

	accountSequence := signerAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, signer.String(), eventParam),
		Log:    resultLog.String(),
	}

}

func (k *Keeper) UnfreezeToken(ctx sdkTypes.Context, symbol string, signer sdkTypes.AccAddress, metadata string) sdkTypes.Result {
	if !k.IsAuthorised(ctx, signer) {
		return sdkTypes.ErrUnauthorized("Not authorised to unfreeze.").Result()
	}

	var token = new(Token)
	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return err.Result()
	}
	return k.unfreezeFungibleToken(ctx, symbol, signer, metadata)
}

func (k *Keeper) unfreezeFungibleToken(ctx sdkTypes.Context, symbol string, signer sdkTypes.AccAddress, metadata string) sdkTypes.Result {

	var token = new(Token)
	k.mustGetTokenData(ctx, symbol, token)

	signerAccount := k.accountKeeper.GetAccount(ctx, signer)
	if signerAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid signer.").Result()
	}

	if !token.Flags.HasFlag(FrozenFlag) {
		return sdkTypes.ErrUnknownRequest("Fungible token is not frozen.").Result()
	}

	token.Flags.RemoveFlag(FrozenFlag)
	token.Metadata = metadata

	k.storeToken(ctx, symbol, token)

	eventParam := []string{symbol, token.Owner.String()}
	eventSignature := "UnfreezeFungibleToken(string,string)"

	accountSequence := signerAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, signer.String(), eventParam),
		Log:    resultLog.String(),
	}

}

// FreezeFungibleTokenAccount
func (k *Keeper) FreezeFungibleTokenAccount(ctx sdkTypes.Context, symbol string, owner sdkTypes.AccAddress, tokenAccount sdkTypes.AccAddress, metadata string) sdkTypes.Result {
	var token = new(Token)
	if exists := k.getTokenData(ctx, symbol, token); !exists {
		return sdkTypes.ErrUnknownRequest("No such fungible token.").Result()
	}

	if !k.IsAuthorised(ctx, owner) {
		return sdkTypes.ErrUnauthorized("Not authorised to freeze token account.").Result()
	}

	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, owner)
	if ownerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid signer.").Result()
	}

	fungibleAccount := k.getFungibleAccount(ctx, symbol, tokenAccount)
	if fungibleAccount == nil {
		return sdkTypes.ErrUnknownRequest("No such token account to freeze.").Result()
	}

	if fungibleAccount.Frozen {
		return sdkTypes.ErrUnknownRequest("Fungible token account already frozen.").Result()
	}

	fungibleAccount.Frozen = true
	fungibleAccount.Metadata = metadata

	k.storeFungibleAccount(ctx, symbol, fungibleAccount)

	eventParam := []string{symbol, tokenAccount.String()}
	eventSignature := "FrozenFungibleTokenAccount(string,string)"

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, owner.String(), eventParam),
		Log:    resultLog.String(),
	}
}

func (k *Keeper) UnfreezeFungibleTokenAccount(ctx sdkTypes.Context, symbol string, owner sdkTypes.AccAddress, tokenAccount sdkTypes.AccAddress, metadata string) sdkTypes.Result {
	if !k.IsAuthorised(ctx, owner) {
		return sdkTypes.ErrUnauthorized("Not authorised to unfreeze token account.").Result()
	}

	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, owner)
	if ownerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid signer.").Result()
	}

	var token = new(Token)
	if exists := k.getTokenData(ctx, symbol, token); !exists {
		return sdkTypes.ErrUnknownRequest("No such fungible token.").Result()
	}

	fungibleAccount := k.getFungibleAccount(ctx, symbol, tokenAccount)
	if fungibleAccount == nil {
		return sdkTypes.ErrUnknownRequest("No such fungible token account to unfreeze.").Result()
	}

	if !fungibleAccount.Frozen {
		return sdkTypes.ErrUnknownRequest("Fungible token account not frozen.").Result()
	}

	fungibleAccount.Frozen = false
	fungibleAccount.Metadata = metadata

	k.storeFungibleAccount(ctx, symbol, fungibleAccount)

	eventParam := []string{symbol, tokenAccount.String()}
	eventSignature := "UnfreezeFungibleTokenAccount(string,string)"

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, owner.String(), eventParam),
		Log:    resultLog.String(),
	}

}

func (k *Keeper) ApproveTransferTokenOwnership(ctx sdkTypes.Context, symbol string, from sdkTypes.AccAddress) sdkTypes.Result {
	if !k.IsAuthorised(ctx, from) {
		return sdkTypes.ErrUnauthorized("Not authorised to accept transfer token ownership.").Result()
	}

	fromWalletAccount := k.accountKeeper.GetAccount(ctx, from)
	if fromWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid signer.").Result()
	}

	var token = new(Token)
	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return err.Result()
	}

	if !token.Flags.HasFlag(TransferTokenOwnershipFlag) {
		return types.ErrInvalidTokenAction().Result()
	}

	token.Flags.AddFlag(AcceptTokenOwnershipFlag)
	token.Flags.AddFlag(ApproveTransferTokenOwnershipFlag)

	k.storeToken(ctx, symbol, token)

	eventParam := []string{symbol, token.Owner.String(), token.NewOwner.String()}
	eventSignature := "ApprovedTransferTokenOwnership(string,string,string)"

	accountSequence := fromWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}
}

func (k *Keeper) RejectTransferTokenOwnership(ctx sdkTypes.Context, symbol string, from sdkTypes.AccAddress) sdkTypes.Result {
	if !k.IsAuthorised(ctx, from) {
		return sdkTypes.ErrUnauthorized("Not authorised to accept transfer token ownership.").Result()
	}

	fromWalletAccount := k.accountKeeper.GetAccount(ctx, from)
	if fromWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid signer.").Result()
	}

	var token = new(Token)
	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return err.Result()
	}

	if !token.Flags.HasFlag(TransferTokenOwnershipFlag) {
		return types.ErrInvalidTokenAction().Result()
	}

	token.Flags.RemoveFlag(TransferTokenOwnershipFlag)

	var emptyAccAddr sdkTypes.AccAddress
	token.NewOwner = emptyAccAddr

	k.storeToken(ctx, symbol, token)

	eventParam := []string{symbol, token.Owner.String(), token.NewOwner.String()}
	eventSignature := "RejectedTransferTokenOwnership(string,string,string)"

	accountSequence := fromWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}
}

func (k *Keeper) getTokenData(ctx sdkTypes.Context, symbol string, target interface{}) bool {
	store := ctx.KVStore(k.key)
	key := getTokenKey(symbol)

	tokenData := store.Get([]byte(key))
	if tokenData == nil {
		return false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(tokenData, target)

	return true
}

func (k *Keeper) storeToken(ctx sdkTypes.Context, symbol string, token interface{}) {
	store := ctx.KVStore(k.key)
	key := getTokenKey(symbol)
	tokenData := k.cdc.MustMarshalBinaryLengthPrefixed(token)

	store.Set([]byte(key), tokenData)
}

// Accounts
func (k *Keeper) getFungibleAccount(ctx sdkTypes.Context, symbol string, owner sdkTypes.AccAddress) *FungibleTokenAccount {
	key := getFungibleAccountKey(symbol, owner)

	store := ctx.KVStore(k.key)
	value := store.Get(key)
	if len(value) == 0 {
		return nil
	}

	var account = new(FungibleTokenAccount)
	k.cdc.MustUnmarshalBinaryLengthPrefixed(value, account)

	return account
}

func (k *Keeper) createFungibleAccount(ctx sdkTypes.Context, symbol string, owner sdkTypes.AccAddress) *FungibleTokenAccount {
	account := &FungibleTokenAccount{
		Owner:   owner,
		Frozen:  false,
		Balance: sdkTypes.NewUint(0),
	}

	k.storeFungibleAccount(ctx, symbol, account)

	return account
}

func (k *Keeper) storeFungibleAccount(ctx sdkTypes.Context, symbol string, account *FungibleTokenAccount) {
	store := ctx.KVStore(k.key)
	key := getFungibleAccountKey(symbol, account.Owner)
	accountData := k.cdc.MustMarshalBinaryLengthPrefixed(account)

	store.Set(key, accountData)
}

func (k *Keeper) getAnyAccount(ctx sdkTypes.Context, symbol string, owner sdkTypes.AccAddress) interface{} {
	return k.getFungibleAccount(ctx, symbol, owner)
}

func (k *Keeper) mustGetTokenData(ctx sdkTypes.Context, symbol string, target interface{}) sdkTypes.Error {
	if exists := k.getTokenData(ctx, symbol, target); !exists {
		return types.ErrInvalidTokenSymbol(symbol)
	}
	return nil
}

func (k *Keeper) mustGetAnyTokenData(ctx sdkTypes.Context, symbol string) (interface{}, sdkTypes.Error) {
	var target interface{}

	target = new(Token)

	err := k.mustGetTokenData(ctx, symbol, target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

// Querying
func (k *Keeper) ListTokens(ctx sdkTypes.Context) []Token {
	store := ctx.KVStore(k.key)
	start := []byte("symbol:")
	end := []byte("symbom:")
	iter := store.Iterator(start, end)
	defer iter.Close()

	var lst = make([]Token, 0)

	for {
		if !iter.Valid() {
			break
		}
		var t = new(Token)
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &t)
		lst = append(lst, *t)

		iter.Next()
	}
	return lst
}

func (k *Keeper) GetTokenData(ctx sdkTypes.Context, symbol string) (interface{}, sdkTypes.Error) {
	res, err := k.mustGetAnyTokenData(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (k *Keeper) GetAccount(ctx sdkTypes.Context, symbol string, account sdkTypes.AccAddress) (interface{}, sdkTypes.Error) {
	_, err := k.mustGetAnyTokenData(ctx, symbol)
	if err != nil {
		return nil, err
	}

	return k.getAnyAccount(ctx, symbol, account), nil
}

func (k *Keeper) CheckApprovedToken(ctx sdkTypes.Context, symbol string) bool {
	var token = new(Token)
	k.mustGetTokenData(ctx, symbol, token)
	if token != nil {
		if token.Flags.HasFlag(ApprovedFlag) {
			return true
		}
	}

	return false
}

func (k *Keeper) IsTokenFrozen(ctx sdkTypes.Context, symbol string) bool {
	var token = new(Token)
	k.mustGetTokenData(ctx, symbol, token)

	if token != nil {
		if token.Flags.HasFlag(FrozenFlag) {
			return true
		}
	}

	return false
}

func (k *Keeper) IsFungibleTokenAccountFrozen(ctx sdkTypes.Context, account sdkTypes.AccAddress, symbol string) bool {

	tokenAccount := k.getFungibleAccount(ctx, symbol, account)

	if tokenAccount != nil {
		if tokenAccount.Frozen {
			return true
		}
	}

	return false
}

func (k *Keeper) subFungibleToken(ctx sdkTypes.Context, symbol string, address sdkTypes.AccAddress, value sdkTypes.Uint) sdkTypes.Error {
	account := k.getFungibleAccount(ctx, symbol, address)
	if account != nil {
		account.Balance = account.Balance.Sub(value)
		k.storeFungibleAccount(ctx, symbol, account)
		return nil
	}

	return types.ErrInvalidTokenAccount()
}

func (k *Keeper) addFungibleToken(ctx sdkTypes.Context, symbol string, address sdkTypes.AccAddress, value sdkTypes.Uint) sdkTypes.Error {
	account := k.getFungibleAccount(ctx, symbol, address)
	if account != nil {
		account.Balance = account.Balance.Add(value)
		k.storeFungibleAccount(ctx, symbol, account)
		return nil
	}

	return types.ErrInvalidTokenAccount()
}

func (k *Keeper) IsVerifyableTransferTokenOwnership(ctx sdkTypes.Context, symbol string) bool {
	var token = new(Token)
	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return false
	}

	if token.Flags.HasFlag(TransferTokenOwnershipFlag) && !token.Flags.HasFlag(AcceptTokenOwnershipFlag) {
		return true
	}

	return false
}
