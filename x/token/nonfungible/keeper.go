package nonfungible

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
	NonFungibleFlag  types.Bitmask = 0x0001
	MintFlag         types.Bitmask = 0x0002
	BurnFlag         types.Bitmask = 0x0004
	FrozenFlag       types.Bitmask = 0x0008
	ApprovedFlag     types.Bitmask = 0x0010
	TransferableFlag types.Bitmask = 0x0020
	ModifiableFlag   types.Bitmask = 0x0040
	PubFlag          types.Bitmask = 0x0080

	TransferTokenOwnershipFlag        types.Bitmask = 0x0100
	ApproveTransferTokenOwnershipFlag types.Bitmask = 0x0200
	AcceptTokenOwnershipFlag          types.Bitmask = 0x0400

	NonFungibleTokenMask = NonFungibleFlag + MintFlag
)

type Token struct {
	Flags             types.Bitmask
	Name              string
	Symbol            string
	Owner             sdkTypes.AccAddress
	NewOwner          sdkTypes.AccAddress
	Properties        string
	Metadata          string
	TotalSupply       sdkTypes.Uint
	TransferLimit     sdkTypes.Uint
	MintLimit         sdkTypes.Uint
	EndorserList      []sdkTypes.AccAddress
	EndorserListLimit sdkTypes.Uint
}

type Item struct {
	ID            string
	Properties    string
	Metadata      string
	TransferLimit sdkTypes.Uint
	Frozen        bool
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
	return store.Has(key)
}

func (k *Keeper) CreateNonFungibleToken(
	ctx sdkTypes.Context,
	name string,
	symbol string,
	owner sdkTypes.AccAddress,
	properties string,
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
		Flags:       NonFungibleTokenMask,
		Symbol:      symbol,
		Owner:       owner,
		Properties:  properties,
		Metadata:    metadata,
		TotalSupply: zero,
	}

	k.storeToken(ctx, symbol, token)

	eventParam := []string{symbol, owner.String(), fee.To.String(), fee.Value}
	eventSignature := "CreatedNonFungibleToken(string,string,string,bignumber)"
	event := types.MakeMxwEvents(eventSignature, owner.String(), eventParam)

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())
	return sdkTypes.Result{
		Events: applicationFeeResult.Events.AppendEvents(event),
		Log:    resultLog.String(),
	}
}

// ApproveToken
func (k *Keeper) ApproveToken(ctx sdkTypes.Context, symbol string, tokenFees []TokenFee, mintLimit, transferLimit sdkTypes.Uint, endorserList []sdkTypes.AccAddress, signer sdkTypes.AccAddress, burnable, transferable, modifiable, public bool, endorserListLimit sdkTypes.Uint) sdkTypes.Result {
	if !k.IsAuthorised(ctx, signer) {
		return sdkTypes.ErrUnauthorized("Not authorised to approve.").Result()
	}

	return k.approveNonFungibleToken(ctx, symbol, tokenFees, mintLimit, transferLimit, signer, endorserList, burnable, transferable, modifiable, public, endorserListLimit)
}

func (k *Keeper) approveNonFungibleToken(ctx sdkTypes.Context, symbol string, tokenFees []TokenFee, mintLimit, transferLimit sdkTypes.Uint, signer sdkTypes.AccAddress, endorserList []sdkTypes.AccAddress, burnable, transferable, modifiable, public bool, endorserListLimit sdkTypes.Uint) sdkTypes.Result {
	var token = new(Token)
	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return err.Result()
	}

	if endorserListLimit.LTE(sdkTypes.NewUintFromString("0")) {
		return sdkTypes.ErrInternal("Endorserlist limit cannot less than or equal zero.").Result()
	}

	if sdkTypes.NewUint(uint64(len(endorserList))).GT(endorserListLimit) {
		return sdkTypes.ErrInternal("Endorserlist limit exceeded.").Result()

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
		err := k.feeKeeper.AssignFeeToNonFungibleTokenAction(ctx, tokenFee.FeeName, token.Symbol, tokenFee.Action)
		if err != nil {
			return err.Result()
		}
	}

	token.Flags.AddFlag(ApprovedFlag)
	if burnable {
		token.Flags.AddFlag(BurnFlag)
	}
	if transferable {
		token.Flags.AddFlag(TransferableFlag)
	}
	if modifiable {
		token.Flags.AddFlag(ModifiableFlag)
	}
	if public {
		token.Flags.AddFlag(PubFlag)
	}

	token.TransferLimit = transferLimit
	token.MintLimit = mintLimit
	token.EndorserList = endorserList
	token.EndorserListLimit = endorserListLimit

	k.storeToken(ctx, symbol, token)

	// Event: Approved fungible token
	eventParam := []string{symbol, token.Owner.String()}
	eventSignature := "ApprovedNonFungibleToken(string,string)"
	events := types.MakeMxwEvents(eventSignature, signer.String(), eventParam)

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: events,
		Log:    resultLog.String(),
	}

}

func (k *Keeper) RejectToken(ctx sdkTypes.Context, symbol string, signer sdkTypes.AccAddress) sdkTypes.Result {

	var token = new(Token)

	if !k.IsAuthorised(ctx, signer) {
		return sdkTypes.ErrUnauthorized("Not authorised to reject.").Result()
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
	store.Delete(tokenTypeKey)

	eventParam := []string{symbol, token.Owner.String()}
	eventSignature := "RejectedNonFungibleToken(string,string)"

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, signer.String(), eventParam),
		Log:    resultLog.String(),
	}

}

func (k *Keeper) FreezeToken(ctx sdkTypes.Context, symbol string, signer sdkTypes.AccAddress) sdkTypes.Result {
	if !k.IsAuthorised(ctx, signer) {
		return sdkTypes.ErrUnauthorized("Not authorised to freeze.").Result()
	}

	var token = new(Token)
	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return err.Result()
	}

	return k.freezeNonFungibleToken(ctx, symbol, signer)

}

func (k *Keeper) freezeNonFungibleToken(ctx sdkTypes.Context, symbol string, signer sdkTypes.AccAddress) sdkTypes.Result {
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
	k.storeToken(ctx, symbol, token)

	eventParam := []string{symbol, token.Owner.String()}
	eventSignature := "FrozenNonFungibleToken(string,string)"

	accountSequence := signerAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, signer.String(), eventParam),
		Log:    resultLog.String(),
	}

}

func (k *Keeper) UnfreezeToken(ctx sdkTypes.Context, symbol string, signer sdkTypes.AccAddress) sdkTypes.Result {
	if !k.IsAuthorised(ctx, signer) {
		return sdkTypes.ErrUnauthorized("Not authorised to unfreeze.").Result()
	}

	var token = new(Token)
	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return err.Result()
	}
	return k.unfreezeNonFungibleToken(ctx, symbol, signer)
}

func (k *Keeper) unfreezeNonFungibleToken(ctx sdkTypes.Context, symbol string, signer sdkTypes.AccAddress) sdkTypes.Result {

	var token = new(Token)
	k.mustGetTokenData(ctx, symbol, token)

	signerAccount := k.accountKeeper.GetAccount(ctx, signer)
	if signerAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid signer.").Result()
	}

	if !token.Flags.HasFlag(FrozenFlag) {
		return sdkTypes.ErrUnknownRequest("Non-fungible token is not frozen.").Result()
	}

	token.Flags.RemoveFlag(FrozenFlag)

	k.storeToken(ctx, symbol, token)

	eventParam := []string{symbol, token.Owner.String()}
	eventSignature := "UnfreezeNonFungibleToken(string,string)"

	accountSequence := signerAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, signer.String(), eventParam),
		Log:    resultLog.String(),
	}

}

// FreezeNonFungibleItem
func (k *Keeper) FreezeNonFungibleItem(ctx sdkTypes.Context, symbol string, owner, itemOwner sdkTypes.AccAddress, itemID string, metadata string) sdkTypes.Result {
	var token = new(Token)
	if exists := k.GetNonfungibleTokenDataInfo(ctx, symbol, token); !exists {
		return sdkTypes.ErrUnknownRequest("No such non-fungible token.").Result()
	}

	if !k.IsAuthorised(ctx, owner) {
		return sdkTypes.ErrUnauthorized("Not authorised to freeze non-fungible item.").Result()
	}

	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, owner)
	if ownerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid signer.").Result()
	}

	nonFungibleItem := k.GetNonFungibleItem(ctx, symbol, itemID)
	if nonFungibleItem == nil {
		return sdkTypes.ErrUnknownRequest("No such non-fungible item to freeze.").Result()
	}

	itemOwner = k.GetNonFungibleItemOwnerInfo(ctx, symbol, itemID)
	if itemOwner == nil {
		return sdkTypes.ErrUnknownRequest("Invalid item owner.").Result()
	}

	if nonFungibleItem.Frozen {
		return sdkTypes.ErrUnknownRequest("Non-fungible item already frozen.").Result()
	}

	nonFungibleItem.Frozen = true

	k.storeNonFungibleItem(ctx, symbol, itemOwner, nonFungibleItem)

	eventParam := []string{symbol, itemID, owner.String()}
	eventSignature := "FrozenNonFungibleItem(string,string,string)"

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, owner.String(), eventParam),
		Log:    resultLog.String(),
	}
}

func (k *Keeper) UnfreezeNonFungibleItem(ctx sdkTypes.Context, symbol string, owner sdkTypes.AccAddress, itemID string, metadata string) sdkTypes.Result {
	if !k.IsAuthorised(ctx, owner) {
		return sdkTypes.ErrUnauthorized("Not authorised to unfreeze non-fungible token account.").Result()
	}

	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, owner)
	if ownerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid signer.").Result()
	}

	var token = new(Token)
	if exists := k.GetNonfungibleTokenDataInfo(ctx, symbol, token); !exists {
		return sdkTypes.ErrUnknownRequest("No such non-fungible token.").Result()
	}

	nonFungibleItem := k.GetNonFungibleItem(ctx, symbol, itemID)
	if nonFungibleItem == nil {
		return sdkTypes.ErrUnknownRequest("No such non-fungible item to unfreeze.").Result()
	}

	itemOwner := k.GetNonFungibleItemOwnerInfo(ctx, symbol, itemID)
	if itemOwner == nil {
		return sdkTypes.ErrUnknownRequest("Invalid item owner.").Result()
	}

	if !nonFungibleItem.Frozen {
		return sdkTypes.ErrUnknownRequest("Non-fungible item not frozen.").Result()
	}

	nonFungibleItem.Frozen = false

	k.storeNonFungibleItem(ctx, symbol, itemOwner, nonFungibleItem)

	eventParam := []string{symbol, itemID, owner.String()}
	eventSignature := "UnfreezeNonFungibleItem(string,string,string)"

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, owner.String(), eventParam),
		Log:    resultLog.String(),
	}

}

func (k *Keeper) ApproveTransferTokenOwnership(ctx sdkTypes.Context, symbol string, from sdkTypes.AccAddress) sdkTypes.Result {
	if !k.IsAuthorised(ctx, from) {
		return sdkTypes.ErrUnauthorized("Not authorised to approve transfer token ownership.").Result()
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
	eventSignature := "ApprovedTransferNonFungibleTokenOwnership(string,string,string)"

	accountSequence := fromWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}
}

func (k *Keeper) RejectTransferTokenOwnership(ctx sdkTypes.Context, symbol string, from sdkTypes.AccAddress) sdkTypes.Result {
	if !k.IsAuthorised(ctx, from) {
		return sdkTypes.ErrUnauthorized("Not authorised to reject transfer token ownership.").Result()
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
	eventSignature := "RejectedTransferNonFungibleTokenOwnership(string,string,string)"

	accountSequence := fromWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}
}

func (k *Keeper) GetNonfungibleTokenDataInfo(ctx sdkTypes.Context, symbol string, target interface{}) bool {
	store := ctx.KVStore(k.key)
	key := getTokenKey(symbol)

	tokenData := store.Get(key)
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

	store.Set(key, tokenData)
}

func (k *Keeper) increaseMintItemLimit(ctx sdkTypes.Context, symbol string, owner sdkTypes.AccAddress) {
	store := ctx.KVStore(k.key)
	key := getMintItemLimitKey(symbol, owner)

	v := store.Get(key)
	if v != nil {
		counter := sdkTypes.NewUintFromString(string(v))
		counter = counter.Add(sdkTypes.NewUintFromString("1"))
		store.Set(key, []byte(counter.String()))
	} else {
		counter := sdkTypes.NewUintFromString("1")
		store.Set(key, []byte(counter.String()))
	}

}

// Item
func (k *Keeper) GetNonFungibleItem(ctx sdkTypes.Context, symbol string, itemID string) *Item {
	itemKey := getNonFungibleItemKey(symbol, []byte(itemID))
	store := ctx.KVStore(k.key)

	itemValue := store.Get(itemKey)
	if len(itemValue) == 0 {
		return nil
	}

	var item = new(Item)
	k.cdc.MustUnmarshalBinaryLengthPrefixed(itemValue, item)

	return item
}

func (k *Keeper) GetNonFungibleItemOwnerInfo(ctx sdkTypes.Context, symbol string, itemID string) sdkTypes.AccAddress {
	ownerKey := getNonFungibleItemOwnerKey(symbol, []byte(itemID))
	store := ctx.KVStore(k.key)

	return store.Get(ownerKey)
}

func (k *Keeper) createNonFungibleItem(ctx sdkTypes.Context, symbol string, owner sdkTypes.AccAddress, itemID, properties, metadata string) *Item {
	item := &Item{
		ID:         itemID,
		Properties: properties,
		Metadata:   metadata,
		Frozen:     false,
	}

	k.storeNonFungibleItem(ctx, symbol, owner, item)

	return item
}

func (k *Keeper) storeNonFungibleItem(ctx sdkTypes.Context, symbol string, owner sdkTypes.AccAddress, item *Item) {
	store := ctx.KVStore(k.key)
	itemKey := getNonFungibleItemKey(symbol, []byte(item.ID))
	ownerKey := getNonFungibleItemOwnerKey(symbol, []byte(item.ID))

	itemData := k.cdc.MustMarshalBinaryLengthPrefixed(item)

	store.Set(itemKey, itemData)
	store.Set(ownerKey, owner.Bytes())
}

func (k *Keeper) getAnyItem(ctx sdkTypes.Context, symbol string, itemID string) interface{} {
	return k.GetNonFungibleItem(ctx, symbol, itemID)
}

func (k *Keeper) mustGetTokenData(ctx sdkTypes.Context, symbol string, target interface{}) sdkTypes.Error {
	if exists := k.GetNonfungibleTokenDataInfo(ctx, symbol, target); !exists {
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
func (k *Keeper) ListTokenData(ctx sdkTypes.Context) ([]string, []string) {
	/*store := ctx.KVStore(k.key)
	start := "symbol:"
	end := "symbol;"
	iter := store.Iterator([]byte(start), []byte(end))
	defer iter.Close()

	var Token = make([]string, 0)
	var nonfungibleToken = make([]string, 0)

	for {
		if !iter.Valid() {
			break
		}

		key := string(iter.Key())
		value := string(iter.Value())

		keysplit := strings.Split(string(key), ":")
		if len(keysplit) != 2 {
			panic(fmt.Sprintf("Invalid key: %s", key))
		}

		tokenSymbol := keysplit[1]

		if tokenSymbol == "" {
			panic(fmt.Sprintf("Invalid token symbol"))
		}

		if value == FungibleTokenType {
			Token = append(Token, tokenSymbol)
		} else if value == NonFungibleTokenType {
			nonfungibleToken = append(nonfungibleToken, tokenSymbol)
		} else {
			panic(fmt.Sprintf("Invalid value name %s for key %s", value, key))
		}

		iter.Next()
	}

	return Token, nonfungibleToken
	*/
	// check ListAllSysFeeSetting as an example
	return nil, nil
}

func (k *Keeper) GetTokenData(ctx sdkTypes.Context, symbol string) (interface{}, sdkTypes.Error) {
	res, err := k.mustGetAnyTokenData(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (k *Keeper) GetItem(ctx sdkTypes.Context, symbol string, itemID string) (interface{}, sdkTypes.Error) {
	_, err := k.mustGetAnyTokenData(ctx, symbol)
	if err != nil {
		return nil, err
	}

	return k.getAnyItem(ctx, symbol, itemID), nil
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

func (k *Keeper) IsNonFungibleItemFrozen(ctx sdkTypes.Context, symbol string, itemID string) bool {

	item := k.GetNonFungibleItem(ctx, symbol, itemID)

	if item != nil {
		if item.Frozen {
			return true
		}
	}

	return false
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

func (k *Keeper) IsItemIDUnique(ctx sdkTypes.Context, symbol string, itemID string) bool {

	item := k.GetNonFungibleItem(ctx, symbol, itemID)

	if item != nil {
		return false
	}

	return true
}

func (k *Keeper) IsTokenEndorser(ctx sdkTypes.Context, symbol string, endorser sdkTypes.AccAddress) bool {
	token := new(Token)
	k.GetNonfungibleTokenDataInfo(ctx, symbol, token)

	var endorsers types.AddressHolder

	if token != nil {

		endorsers = token.EndorserList
		if endorsers != nil {
			_, contain := endorsers.Contains(endorser)
			return contain
		}
		return true
	}
	return false
}

func (k *Keeper) GetEndorserList(ctx sdkTypes.Context, symbol string) []sdkTypes.AccAddress {
	token := new(Token)
	k.GetNonfungibleTokenDataInfo(ctx, symbol, token)

	if token != nil {

		return token.EndorserList
	}
	return nil
}

// Querying
func (k *Keeper) ListTokens(ctx sdkTypes.Context) []Token {
	store := ctx.KVStore(k.key)
	start := getTokenKey(string(0x00))
	end := getTokenKey(string(0xFF))
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
