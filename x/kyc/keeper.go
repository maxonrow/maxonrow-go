package kyc

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	KycAddressMaxLength = 64
)

type Keeper struct {
	whitelistedStoreKey sdkTypes.StoreKey
	kycDataStoreKey     sdkTypes.StoreKey
	accountKeeper       *sdkAuth.AccountKeeper

	cdc *codec.Codec
}

var prefixWhitelisted = []byte("0x01")
var prefixAuthorised = []byte("0x02")
var prefixProvider = []byte("0x03")
var prefixIssuer = []byte("0x04")
var prefixKycData = []byte("0x05")

// keys
func getWhitelistedKey(addr sdkTypes.AccAddress) []byte {
	return append(prefixWhitelisted, addr.Bytes()...)
}
func getKycDataKey(data []byte) []byte {
	return append(prefixKycData, data...)
}
func getAuthorisedKey() []byte {
	return prefixAuthorised
}

func getProviderKey() []byte {
	return prefixProvider
}

func getIssuerKey() []byte {
	return prefixIssuer
}

func NewKeeper(cdc *codec.Codec, accountKeeper *sdkAuth.AccountKeeper, whitelistedStoreKey, kycDataStoreKey sdkTypes.StoreKey) Keeper {
	return Keeper{
		cdc:                 cdc,
		accountKeeper:       accountKeeper,
		whitelistedStoreKey: whitelistedStoreKey,
		kycDataStoreKey:     kycDataStoreKey,
	}
}

func (k *Keeper) SetAuthorisedAddresses(ctx sdkTypes.Context, addresses []sdkTypes.AccAddress) {
	authorisedStore := ctx.KVStore(k.whitelistedStoreKey)
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
	authorisedStore := ctx.KVStore(k.whitelistedStoreKey)
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
	authorisedStore := ctx.KVStore(k.whitelistedStoreKey)
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
	authorisedStore := ctx.KVStore(k.whitelistedStoreKey)
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

	authorisedStore := ctx.KVStore(k.whitelistedStoreKey)
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
	authorisedStore := ctx.KVStore(k.whitelistedStoreKey)
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
	authorisedStore := ctx.KVStore(k.whitelistedStoreKey)
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

	authorisedStore := ctx.KVStore(k.whitelistedStoreKey)
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
	authorisedStore := ctx.KVStore(k.whitelistedStoreKey)
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

// * TO DO: VALIDATION
func (k Keeper) Whitelist(ctx sdkTypes.Context, targetAddress sdkTypes.AccAddress, kycAddress string) {

	acc := k.accountKeeper.GetAccount(ctx, targetAddress)
	if acc == nil {
		acc = k.accountKeeper.NewAccountWithAddress(ctx, targetAddress)
		acc.SetAccountNumber(k.accountKeeper.GetNextAccountNumber(ctx))
		k.accountKeeper.SetAccount(ctx, acc)
	}

	whitelistStore := ctx.KVStore(k.whitelistedStoreKey)
	kycDataStore := ctx.KVStore(k.kycDataStoreKey)
	whitelistedKey := getWhitelistedKey(targetAddress)
	kycDataKey := getKycDataKey([]byte(kycAddress))
	whitelistStore.Set(whitelistedKey, []byte(kycAddress))
	kycDataStore.Set(kycDataKey, targetAddress.Bytes())

}

func (k Keeper) RevokeWhitelist(ctx sdkTypes.Context, targetAddress sdkTypes.AccAddress, owner sdkTypes.AccAddress) sdkTypes.Result {

	kycDataByte := k.GetKycAddress(ctx, targetAddress)

	whitelistStore := ctx.KVStore(k.whitelistedStoreKey)
	kycDataStore := ctx.KVStore(k.kycDataStoreKey)
	whitelistedKey := getWhitelistedKey(targetAddress)
	kycDataKey := getKycDataKey(kycDataByte)

	kycDataStore.Delete(kycDataKey)
	whitelistStore.Delete(whitelistedKey)

	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, owner)
	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	eventParam := []string{targetAddress.String(), string(kycDataByte)}
	eventSignature := "RevokedWhitelist(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, owner.String(), eventParam),
		Log:    resultLog.String()}
}

func (k Keeper) IsKycAddressExist(ctx sdkTypes.Context, kycAddress string) bool {
	kycDataStore := ctx.KVStore(k.kycDataStoreKey)
	key := getKycDataKey([]byte(kycAddress))
	return kycDataStore.Has(key)
}

func (k Keeper) IsWhitelisted(ctx sdkTypes.Context, address sdkTypes.AccAddress) bool {
	whitelistStore := ctx.KVStore(k.whitelistedStoreKey)
	key := getWhitelistedKey(address)
	return whitelistStore.Has(key)
}

func (k Keeper) GetKycAddress(ctx sdkTypes.Context, address sdkTypes.AccAddress) []byte {
	whitelistStore := ctx.KVStore(k.whitelistedStoreKey)
	key := getWhitelistedKey(address)
	return whitelistStore.Get(key)
}

func (k Keeper) CheckTx(ctx sdkTypes.Context, tx sdkAuth.StdTx) bool {
	allSigners := tx.GetSigners()
	for _, signer := range allSigners {
		if !k.IsWhitelisted(ctx, signer) {
			return false
		}
	}

	return true
}

func (k Keeper) ValidateKycAddress(kycAdd string) sdkTypes.Error {
	if len(kycAdd) == 0 || len(kycAdd) > KycAddressMaxLength {
		return sdkTypes.ErrUnknownRequest(
			fmt.Sprintf("Invalid KycAddress field length: %d", len(kycAdd)))
	}

	if strings.ContainsAny(kycAdd, ";:") {
		return sdkTypes.ErrUnknownRequest("KycAddress cannot contain following characters: ;:")
	}

	return nil
}

func (k Keeper) ValidateSignatures(ctx sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Error {

	var fromSignature Signature
	var issuerSignatures []Signature
	var fromSignBytes, issuerSignBytes []byte
	var fromAddr sdkTypes.AccAddress
	var fromAccountNonce string

	switch msg := msg.(type) {
	case MsgWhitelist:
		fromSignature = NewSignature(msg.KycData.Payload.PubKey, msg.KycData.Payload.Signature)
		// * from sign bytes
		fromSignBytes = msg.KycData.Payload.Kyc.GetFromSignBytes()
		fromAddr = msg.KycData.Payload.Kyc.From
		fromAccountNonce = msg.KycData.Payload.Kyc.Nonce

		//* issuer sign bytes
		issuerSignBytes = msg.KycData.Payload.GetIssuerSignBytes()
		issuerSignatures = msg.KycData.Signatures

	default:
		errMsg := fmt.Sprintf("Unrecognized kyc Msg type: %v", msg.Type())
		return sdkTypes.ErrUnknownRequest(errMsg)
	}

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

	//* verify from sign
	if !(processSig(acc, fromSignature, fromSignBytes)) {

		return sdkTypes.ErrUnauthorized("From signature verification failed.")
	}

	//* at least one issuer and one provider
	providerCounter := 0
	issuerCounter := 0

	//* verify issuer sign
	for i := 0; i < len(issuerSignatures); i++ {
		issuerAddr := sdkTypes.AccAddress(issuerSignatures[i].PubKey.Address())
		issuerAcc := k.accountKeeper.GetAccount(ctx, issuerAddr)

		if k.IsProvider(ctx, issuerAddr) {
			providerCounter++
		} else if k.IsIssuer(ctx, issuerAddr) {
			issuerCounter++
		} else {
			return sdkTypes.ErrUnauthorized("Unauthorized signature.")
		}

		if !(processSig(issuerAcc, issuerSignatures[i], issuerSignBytes)) {
			return sdkTypes.ErrUnauthorized("Signature verification failed.")
		}

		issuerAcc.SetSequence(issuerAcc.GetSequence() + 1)
	}

	if providerCounter < 1 || issuerCounter < 1 {
		return sdkTypes.ErrUnauthorized("Required provider and issuer signatures.")
	}

	return nil
}

func processSig(signerAcc exported.Account, signature Signature, signBytes []byte) bool {

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

func (k Keeper) ValidateRevokeWhitelistSignatures(ctx sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Error {

	var fromSignature Signature
	var issuerSignatures []Signature
	var fromSignBytes, issuerSignBytes []byte
	var fromAddr sdkTypes.AccAddress
	var fromAccountNonce string

	switch msg := msg.(type) {
	case MsgRevokeWhitelist:
		fromSignature = NewSignature(msg.RevokePayload.PubKey, msg.RevokePayload.Signature)
		// * from sign bytes
		fromSignBytes = msg.RevokePayload.RevokeKycData.GetRevokeFromSignBytes()
		fromAddr = msg.RevokePayload.RevokeKycData.From
		fromAccountNonce = msg.RevokePayload.RevokeKycData.Nonce

		//* issuer sign bytes
		issuerSignBytes = msg.RevokePayload.GetRevokeIssuerSignBytes()
		issuerSignatures = msg.Signatures

	default:
		errMsg := fmt.Sprintf("Unrecognized Revoke Whitelist Msg type: %v", msg.Type())
		return sdkTypes.ErrUnknownRequest(errMsg)
	}

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

// ProcessPubKey verifies that the given account address matches that of the
// StdSignature. In addition, it will set the public key of the account if it
// has not been smiet.
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

func (k *Keeper) ListAllWhitelistedAccounts(ctx sdkTypes.Context) []sdkTypes.AccAddress {

	store := ctx.KVStore(k.whitelistedStoreKey)
	start := append(prefixWhitelisted, 0x00)
	end := append(prefixWhitelisted, 0xFF)
	iter := store.Iterator(start, end)
	defer iter.Close()

	var lst = make([]sdkTypes.AccAddress, 0)

	for {
		if !iter.Valid() {
			break
		}

		addr := sdkTypes.AccAddress(iter.Key()[len(prefixWhitelisted):])
		lst = append(lst, addr)

		iter.Next()
	}

	return lst
}

/// TODO: check if we can improve the performance
func (k *Keeper) NumOfWhitelisted(ctx sdkTypes.Context) int {
	return len(k.ListAllWhitelistedAccounts(ctx))
}
