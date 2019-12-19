package nameservice

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	sdkBank "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/tendermint/tendermint/crypto"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/bank"
)

type Keeper struct {
	authorizedAddresses []sdkTypes.AccAddress
	namesStoreKey       sdkTypes.StoreKey // map(alias => aliasData)
	ownersStoreKey      sdkTypes.StoreKey // map(address => alias)
	accountKeeper       *sdkAuth.AccountKeeper
	bankKeeper          sdkBank.Keeper
	cdc                 *codec.Codec
}

var prefixAuthorised = []byte("ns/authorised")
var prefixProvider = []byte("ns/provider")
var prefixIssuer = []byte("ns/issuer")

// keys
func getAuthorisedKey() []byte {
	return prefixAuthorised
}

func getProviderKey() []byte {
	return prefixProvider
}

func getIssuerKey() []byte {
	return prefixIssuer
}

func NewKeeper(namesStoreKey sdkTypes.StoreKey, ownersStoreKey sdkTypes.StoreKey, accountKeeper *sdkAuth.AccountKeeper, bankKeeper sdkBank.Keeper, cdc *codec.Codec) Keeper {
	return Keeper{
		namesStoreKey:  namesStoreKey,
		ownersStoreKey: ownersStoreKey,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
		cdc:            cdc,
	}
}

// store the alias before approval, after approval will update the status
func (k Keeper) IsAliasExists(ctx sdkTypes.Context, alias string) bool {
	pendingStore := ctx.KVStore(k.namesStoreKey)
	approvedStore := ctx.KVStore(k.ownersStoreKey)

	aliasKey := getAliasKey(alias)

	if pendingStore.Has([]byte(aliasKey)) {
		return true
	}

	if approvedStore.Has([]byte(aliasKey)) {
		return true
	}

	return false
}

func (k *Keeper) SetAuthorisedAddresses(ctx sdkTypes.Context, addresses []sdkTypes.AccAddress) {
	authorisedStore := ctx.KVStore(k.ownersStoreKey)
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
	authorisedStore := ctx.KVStore(k.ownersStoreKey)
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
	authorisedStore := ctx.KVStore(k.ownersStoreKey)
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
	authorisedStore := ctx.KVStore(k.ownersStoreKey)
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

	authorisedStore := ctx.KVStore(k.ownersStoreKey)
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
	authorisedStore := ctx.KVStore(k.ownersStoreKey)
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
	authorisedStore := ctx.KVStore(k.ownersStoreKey)
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

	authorisedStore := ctx.KVStore(k.ownersStoreKey)
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
	authorisedStore := ctx.KVStore(k.ownersStoreKey)
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

// CreateAlias
func (k *Keeper) CreateAlias(ctx sdkTypes.Context, from sdkTypes.AccAddress, alias string, fee Fee) sdkTypes.Result {
	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, from)

	if ownerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Not authorised to apply for alias creation.").Result()
	}

	if k.IsAliasExists(ctx, alias) {
		return types.ErrAliasIsInUsed().Result()
	}

	// If have any pending alias, not allowed to create
	if k.isHavingAnyPendingAlias(ctx, from) {

		return types.ErrAliasNotAllowedToCreate().Result()
	}

	amt, parseErr := sdkTypes.ParseCoins(fee.Value + types.CIN)
	if parseErr != nil {
		return sdkTypes.ErrInvalidCoins("Parse value to coins failed.").Result()
	}
	sendCoinsErr := k.bankKeeper.SendCoins(ctx, from, fee.To, amt)
	if sendCoinsErr != nil {
		return sendCoinsErr.Result()
	}

	// Overwrite the cosmos sdk events.
	applicationFeeResult := bank.MakeBankSendEvent(ctx, from, fee.To, amt, *k.accountKeeper)

	aliasData := &Alias{
		Name:     alias,
		Owner:    from,
		Approved: false,
		Fee:      sdkTypes.NewUintFromString(fee.Value),
	}

	aliasOwner := &AliasOwner{
		Name:     alias,
		Approved: false,
	}

	k.storeAlias(ctx, alias, aliasData, aliasOwner)

	eventParam := []string{alias, from.String(), fee.To.String(), fee.Value}
	eventSignature := "CreatedAlias(string,string,string,bignumber)"

	accountSequence := ownerWalletAccount.GetSequence()
	var log string
	if accountSequence == 0 {
		log = types.MakeResultLog(accountSequence, ctx.TxBytes())
	} else {
		log = types.MakeResultLog(accountSequence-1, ctx.TxBytes())
	}
	applicationFeeResult.Events = applicationFeeResult.Events.AppendEvents(types.MakeMxwEvents(eventSignature, from.String(), eventParam))

	return sdkTypes.Result{
		Events: applicationFeeResult.Events,
		Log:    log,
	}
}

// ApproveAlias
func (k *Keeper) ApproveAlias(ctx sdkTypes.Context, alias string, signer sdkTypes.AccAddress, metadata string) sdkTypes.Result {
	if !k.IsAuthorised(ctx, signer) {
		return sdkTypes.ErrUnauthorized("Not authorised to approve.").Result()
	}

	pendingAliasData, pendingAliasDataErr := k.getPendingAlias(ctx, alias)
	if pendingAliasDataErr != nil {
		return pendingAliasDataErr.Result()
	}

	pendingOwnerAliasData, pendingAliasOwnerDataErr := k.getPendingAliasOwnerData(ctx, pendingAliasData.Owner)
	if pendingAliasOwnerDataErr != nil {
		return pendingAliasOwnerDataErr.Result()
	}

	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, pendingAliasData.Owner)
	if ownerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid alias owner.").Result()
	}

	pendingAliasData.Approved = true
	pendingAliasData.Metadata = metadata

	pendingOwnerAliasData.Approved = true

	// set alias data
	k.setAlias(ctx, alias, pendingAliasData, pendingOwnerAliasData)

	eventParam := []string{alias, pendingAliasData.Owner.String()}
	eventSignature := "ApprovedAlias(string,string)"

	accountSequence := ownerWalletAccount.GetSequence()
	var log string
	if accountSequence == 0 {
		log = types.MakeResultLog(accountSequence, ctx.TxBytes())
	} else {
		log = types.MakeResultLog(accountSequence-1, ctx.TxBytes())
	}

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, signer.String(), eventParam),
		Log:    log,
	}

}

func (k *Keeper) RejectAlias(ctx sdkTypes.Context, alias string, signer sdkTypes.AccAddress) sdkTypes.Result {

	if !k.IsAuthorised(ctx, signer) {
		return sdkTypes.ErrUnauthorized("Not authorised to reject").Result()
	}

	pendingAliasData, pendingAliasDataErr := k.getPendingAlias(ctx, alias)
	if pendingAliasDataErr != nil {
		return pendingAliasDataErr.Result()
	}

	if pendingAliasData.Name != alias {
		return sdkTypes.ErrUnknownRequest("Alias is not valid.").Result()
	}

	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, pendingAliasData.Owner)
	if ownerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid alias owner.").Result()
	}

	aliasDataStore := ctx.KVStore(k.namesStoreKey)
	aliasKey := getAliasKey(alias)

	aliasDataStore.Delete([]byte(pendingAliasData.Owner.String()))
	aliasDataStore.Delete([]byte(aliasKey))

	eventParam := []string{alias, pendingAliasData.Owner.String()}
	eventSignature := "RejectedAlias(string,string)"

	accountSequence := ownerWalletAccount.GetSequence()
	var log string
	if accountSequence == 0 {
		log = types.MakeResultLog(accountSequence, ctx.TxBytes())
	} else {
		log = types.MakeResultLog(accountSequence-1, ctx.TxBytes())
	}

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, signer.String(), eventParam),
		Log:    log,
	}

}

func (k Keeper) ResolveAlias(ctx sdkTypes.Context, alias string) (string, sdkTypes.Error) {

	aliasData, error := k.getAliasData(ctx, alias)

	return aliasData.Owner.String(), error
}

func (k Keeper) Whois(ctx sdkTypes.Context, address sdkTypes.AccAddress) string {

	aliasOwnerData, _ := k.getAliasOwnerData(ctx, address)

	return aliasOwnerData.Name
}

// List all alias
func (k *Keeper) ListUsedAlias(ctx sdkTypes.Context) []string {
	store := ctx.KVStore(k.namesStoreKey)
	start := "alias:"
	end := "alias;"
	iter := store.Iterator([]byte(start), []byte(end))
	defer iter.Close()

	var alias = make([]string, 0)

	for {
		if !iter.Valid() {
			break
		}

		key := string(iter.Key())

		keysplit := strings.Split(string(key), ":")
		if len(keysplit) != 2 {
			panic(fmt.Sprintf("Invalid key: %s", key))
		}

		aliasName := keysplit[1]

		if aliasName == "" {
			panic(fmt.Sprintf("Invalid alias"))
		}

		alias = append(alias, aliasName)

		iter.Next()
	}

	return alias
}

func (k Keeper) RevokeAlias(ctx sdkTypes.Context, alias string, signer sdkTypes.AccAddress) sdkTypes.Result {

	if !k.IsAuthorised(ctx, signer) {
		return sdkTypes.ErrUnauthorized("Not authorised to remove alias.").Result()
	}

	aliasData, aliasErr := k.getAliasData(ctx, alias)
	if aliasErr != nil {
		return aliasErr.Result()
	}

	signerWalletAccount := k.accountKeeper.GetAccount(ctx, signer)
	if signerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid signer.").Result()
	}

	ownerStore := ctx.KVStore(k.ownersStoreKey)
	aliasKey := getAliasKey(alias)

	ownerStore.Delete([]byte(aliasData.Owner.String()))
	ownerStore.Delete([]byte(aliasKey))

	eventParam := []string{alias, aliasData.Owner.String()}
	eventSignature := "RevokedAlias(string,string)"

	accountSequence := signerWalletAccount.GetSequence()
	var log string
	if accountSequence == 0 {
		log = types.MakeResultLog(accountSequence, ctx.TxBytes())
	} else {
		log = types.MakeResultLog(accountSequence-1, ctx.TxBytes())
	}

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, signer.String(), eventParam),
		Log:    log,
	}

}

func (k Keeper) ValidateSignatures(ctx sdkTypes.Context, msg MsgSetAliasStatus) sdkTypes.Error {

	fromSignature := NewSignature(msg.Payload.PubKey, msg.Payload.Signature)
	// * from sign bytes
	fromSignBytes := msg.Payload.Alias.GetFromSignBytes()
	fromAddr := msg.Payload.Alias.From
	fromAccountNonce := msg.Payload.Alias.Nonce

	//* issuer sign bytes
	issuerSignBytes := msg.Payload.GetIssuerSignBytes()
	issuerSignatures := msg.Signatures

	//* check account sequence with passed in nonce
	acc := k.accountKeeper.GetAccount(ctx, fromAddr)
	if acc == nil {
		nonce, nonceErr := strconv.ParseUint(fromAccountNonce, 10, 64)
		if nonceErr != nil {
			return sdkTypes.ErrInvalidSequence("Wallet signature is invalid.")
		}

		if nonce != 0 {
			return sdkTypes.ErrInvalidSequence("Wallet signature is invalid.")
		}

	} else {
		nonce, nonceErr := strconv.ParseUint(fromAccountNonce, 10, 64)
		sequence := acc.GetSequence()
		if nonceErr != nil {
			return sdkTypes.ErrInvalidSequence("Wallet signature is invalid.")
		}

		if nonce != sequence {
			return sdkTypes.ErrInvalidSequence("Wallet signature is invalid.")
		}
		acc.SetSequence(sequence + 1)
	}

	if !k.IsProvider(ctx, fromAddr) {
		return sdkTypes.ErrUnauthorized("Insufficient provider signature.")
	}

	//* verify from sign
	if !(processSig(acc, fromSignature, fromSignBytes)) {

		return sdkTypes.ErrUnauthorized("From signature verification failed.")
	}

	//* at least one issuer
	issuerCounter := 0

	//* verify issuer sign
	for i := 0; i < len(issuerSignatures); i++ {
		issuerAddr := sdkTypes.AccAddress(issuerSignatures[i].PubKey.Address())
		issuerAcc := k.accountKeeper.GetAccount(ctx, issuerAddr)

		if k.IsIssuer(ctx, issuerAddr) {
			issuerCounter++
		} else {
			return sdkTypes.ErrUnauthorized("Unauthorized signature.")
		}

		if !(processSig(issuerAcc, issuerSignatures[i], issuerSignBytes)) {

			return sdkTypes.ErrUnauthorized("Signature verification failed.")
		}

		issuerAcc.SetSequence(issuerAcc.GetSequence() + 1)
	}

	if issuerCounter < 1 {
		return sdkTypes.ErrUnauthorized("Insufficient issuer signature.")
	}

	return nil
}

func processSig(
	signerAcc exported.Account, signature Signature, signBytes []byte) bool {

	pubKey, res := ProcessPubKey(signerAcc, signature)
	if !res.IsOK() {

		return false
	}

	if signerAcc != nil {
		err := signerAcc.SetPubKey(pubKey)
		if err != nil {

			return false
		}
	}

	return pubKey.VerifyBytes(signBytes, signature.Signature)
}

// ProcessPubKey verifies that the given account address matches that of the
// StdSignature. In addition, it will set the public key of the account if it
// has not been set.
func ProcessPubKey(acc exported.Account, sig Signature) (crypto.PubKey, sdkTypes.Result) {
	var pubKey crypto.PubKey
	if acc != nil {
		pubKey = acc.GetPubKey()
	}

	if pubKey == nil {

		return sig.PubKey, sdkTypes.Result{}
	}
	if acc != nil {
		cryptoPubKey := pubKey
		if cryptoPubKey == nil {
			return nil, sdkTypes.ErrInvalidPubKey("PubKey not found").Result()
		}

		if !bytes.Equal(cryptoPubKey.Address(), acc.GetAddress()) {
			return nil, sdkTypes.ErrInvalidPubKey(
				fmt.Sprintf("PubKey does not match Signer address %s", acc.GetAddress())).Result()
		}

		return cryptoPubKey, sdkTypes.Result{}
	}
	return nil, sdkTypes.Result{}
}

// Pending alias set here
func (k *Keeper) storeAlias(ctx sdkTypes.Context, alias string, aliasData *Alias, aliasOwnerData *AliasOwner) {
	nameStore := ctx.KVStore(k.namesStoreKey)
	aliasKey := getAliasKey(alias)
	aliasInfo := k.cdc.MustMarshalBinaryLengthPrefixed(aliasData)

	ownerKey := aliasData.Owner.String()
	aliasOwnerInfo := k.cdc.MustMarshalBinaryLengthPrefixed(aliasOwnerData)

	nameStore.Set([]byte(ownerKey), aliasOwnerInfo)
	nameStore.Set([]byte(aliasKey), aliasInfo)
}

// Approved alias set here
func (k *Keeper) setAlias(ctx sdkTypes.Context, alias string, aliasData *Alias, aliasOwnerData *AliasOwner) {
	ownerStore := ctx.KVStore(k.ownersStoreKey)
	aliasKey := getAliasKey(alias)
	aliasInfo := k.cdc.MustMarshalBinaryLengthPrefixed(aliasData)

	ownerKey := aliasData.Owner.String()
	aliasOwnerInfo := k.cdc.MustMarshalBinaryLengthPrefixed(aliasOwnerData)

	ownerStore.Set([]byte(ownerKey), aliasOwnerInfo)
	ownerStore.Set([]byte(aliasKey), aliasInfo)

	// Remove from namestore after approved and set into ownerstore
	nameStore := ctx.KVStore(k.namesStoreKey)
	nameStore.Delete([]byte(ownerKey))
	nameStore.Delete([]byte(aliasKey))
}

func (k *Keeper) getAliasData(ctx sdkTypes.Context, alias string) (*Alias, sdkTypes.Error) {

	var aliasData = new(Alias)

	store := ctx.KVStore(k.ownersStoreKey)
	key := getAliasKey(alias)

	aliasInfo := store.Get([]byte(key))

	if aliasInfo == nil {
		return aliasData, types.ErrAliasNotFound()
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(aliasInfo, &aliasData)

	return aliasData, nil
}

func (k *Keeper) getAliasOwnerData(ctx sdkTypes.Context, owner sdkTypes.AccAddress) (*AliasOwner, sdkTypes.Error) {

	var aliasOwner = new(AliasOwner)

	store := ctx.KVStore(k.ownersStoreKey)

	aliasOwnerInfo := store.Get([]byte(owner.String()))

	if aliasOwnerInfo == nil {
		return aliasOwner, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Address does not have any alias %s", owner))
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(aliasOwnerInfo, &aliasOwner)

	return aliasOwner, nil
}

func (k *Keeper) isHavingAnyPendingAlias(ctx sdkTypes.Context, owner sdkTypes.AccAddress) bool {
	// Pending approval alias will be stored in namestore
	store := ctx.KVStore(k.namesStoreKey)
	pendingAliasOwnerInfo := store.Get([]byte(owner.String()))

	if pendingAliasOwnerInfo == nil {
		return false
	}

	return true
}

func (k *Keeper) getPendingAlias(ctx sdkTypes.Context, alias string) (*Alias, sdkTypes.Error) {

	var aliasData = new(Alias)

	store := ctx.KVStore(k.namesStoreKey)
	key := getAliasKey(alias)

	aliasInfo := store.Get([]byte(key))

	if aliasInfo == nil {
		return aliasData, types.ErrAliasNoSuchPendingAlias()
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(aliasInfo, &aliasData)

	return aliasData, nil
}

func (k *Keeper) getPendingAliasOwnerData(ctx sdkTypes.Context, owner sdkTypes.AccAddress) (*AliasOwner, sdkTypes.Error) {

	var aliasOwner = new(AliasOwner)

	store := ctx.KVStore(k.namesStoreKey)

	aliasOwnerInfo := store.Get([]byte(owner.String()))

	if aliasOwnerInfo == nil {
		return aliasOwner, types.ErrAliasNoSuchPendingAlias()
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(aliasOwnerInfo, &aliasOwner)

	return aliasOwner, nil
}

func (k *Keeper) isOwnAnyAlias(ctx sdkTypes.Context, owner sdkTypes.AccAddress) bool {
	// Approved alias will be stored in ownerstore
	store := ctx.KVStore(k.ownersStoreKey)
	approvedAliasOwnerInfo := store.Get([]byte(owner.String()))

	if approvedAliasOwnerInfo == nil {
		return false
	}

	return true
}

func getAliasKey(alias string) string {
	return fmt.Sprintf("alias:%s", alias)
}

type Alias struct {
	Name     string
	Owner    sdkTypes.AccAddress
	Metadata string
	Approved bool
	Fee      sdkTypes.Uint
}

type AliasOwner struct {
	Name     string
	Approved bool
}
