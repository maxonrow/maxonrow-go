package fee

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/types"
)

type Keeper struct {
	key sdkTypes.StoreKey
	cdc *codec.Codec
}

type FeeSetting struct {
	Name       string              `json:"name"`
	Min        sdkTypes.Coins      `json:"min"`
	Max        sdkTypes.Coins      `json:"max"`
	Percentage string              `json:"percentage"`
	Issuer     sdkTypes.AccAddress `json:"issuer"`
}

const (
	TransferFungibleToken  = "transfer"
	MintFungibleToken      = "mint"
	BurnFungibleToken      = "burn"
	TransferTokenOwnership = "transferOwnership"
	AcceptTokenOwnership   = "acceptOwnership"
)

var prefixAuthorised = []byte("0x01")
var prefixFeeCollector = []byte("0x02")
var prefixSysFeeSetting = []byte("0x03")
var prefixMsgFeeSetting = []byte("0x04")
var prefixAccFeeSetting = []byte("0x05")
var prefixTokenFeeSetting = []byte("0x06")
var prefixMultiplier = []byte("0x50")
var prefixTokenMultiplier = []byte("0x51")

// Token Actions
var tokenActions = []string{TransferFungibleToken, MintFungibleToken, BurnFungibleToken, TransferTokenOwnership, AcceptTokenOwnership}

// keys
func getAuthorisedKey() []byte {
	return prefixAuthorised
}

func getFeeCollectorKey(msgType string) []byte {
	return append(prefixFeeCollector, []byte(msgType)...)
}

func getSysFeeSettingKey(name string) []byte {
	return append(prefixSysFeeSetting, []byte(name)...)
}

func getMsgFeeSettingKey(msgType string) []byte {
	return append(prefixMsgFeeSetting, []byte(msgType)...)
}

func getAccFeeSettingKey(acc sdkTypes.AccAddress) []byte {
	return append(prefixAccFeeSetting, acc.Bytes()...)
}

func getTokenFeeSettingKey(symbol, action string) []byte {
	tokenAction := symbol + ":" + action
	return append(prefixTokenFeeSetting, []byte(tokenAction)...)
}

func getFeeMultiplierKey() []byte {
	return prefixMultiplier
}

func getTokenFeeMultiplierKey() []byte {
	return prefixTokenMultiplier
}

func NewKeeper(cdc *codec.Codec, key sdkTypes.StoreKey) Keeper {
	return Keeper{
		cdc: cdc,
		key: key,
	}
}

func (k *Keeper) SetAuthorisedAddresses(ctx sdkTypes.Context, addresses []sdkTypes.AccAddress) {
	feeStore := ctx.KVStore(k.key)
	key := getAuthorisedKey()

	ah := k.GetAuthorisedAddresses(ctx)
	ah.AppendAccAddrs(addresses)
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(ah)
	if err != nil {
		panic(err)
	}

	feeStore.Set(key, bz)
}

func (k *Keeper) GetAuthorisedAddresses(ctx sdkTypes.Context) types.AddressHolder {
	var ah types.AddressHolder
	feeStore := ctx.KVStore(k.key)
	key := getAuthorisedKey()

	bz := feeStore.Get(key)
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

	feeStore := ctx.KVStore(k.key)

	key := getAuthorisedKey()
	authorisedAddresses := k.GetAuthorisedAddresses(ctx)

	for _, authorisedAddress := range addresses {
		authorisedAddresses.Remove(authorisedAddress)
	}

	bz, err := k.cdc.MarshalBinaryLengthPrefixed(authorisedAddresses)
	if err != nil {
		panic(err)
	}

	feeStore.Set(key, bz)
}

//IsAuthorised Check if is authorised
func (k *Keeper) IsAuthorised(ctx sdkTypes.Context, address sdkTypes.AccAddress) bool {
	ah := k.GetAuthorisedAddresses(ctx)

	_, ok := ah.Contains(address)
	return ok
}

func (k *Keeper) SetFeeCollectorAddresses(ctx sdkTypes.Context, msgType string, addresses []sdkTypes.AccAddress) {
	feeStore := ctx.KVStore(k.key)
	key := getFeeCollectorKey(msgType)

	ah := k.GetFeeCollectorAddresses(ctx, msgType)
	ah.AppendAccAddrs(addresses)
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(ah)
	if err != nil {
		panic(err)
	}

	feeStore.Set(key, bz)
}

func (k *Keeper) GetFeeCollectorAddresses(ctx sdkTypes.Context, msgType string) types.AddressHolder {
	var ah types.AddressHolder
	feeStore := ctx.KVStore(k.key)
	key := getFeeCollectorKey(msgType)

	bz := feeStore.Get(key)
	if bz == nil {
		return ah
	}

	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &ah)
	if err != nil {
		panic(err)
	}
	return ah
}

func (k *Keeper) IsFeeCollector(ctx sdkTypes.Context, msgType string, address sdkTypes.AccAddress) bool {
	ah := k.GetFeeCollectorAddresses(ctx, msgType)

	_, ok := ah.Contains(address)
	return ok
}

func (k *Keeper) RemoveFeeCollectorAddress(ctx sdkTypes.Context, msgType string, feeCollectorAddress sdkTypes.AccAddress) {

	feeStore := ctx.KVStore(k.key)

	key := getFeeCollectorKey(msgType)
	feeCollectorAddresses := k.GetFeeCollectorAddresses(ctx, msgType)

	feeCollectorAddresses.Remove(feeCollectorAddress)

	bz, err := k.cdc.MarshalBinaryLengthPrefixed(feeCollectorAddresses)
	if err != nil {
		panic(err)
	}

	feeStore.Set(key, bz)
}

// FeeSettingExists check if the fee setting already exists in KVStore
func (k *Keeper) FeeSettingExists(ctx sdkTypes.Context, feeSettingType string) bool {
	store := ctx.KVStore(k.key)
	key := getSysFeeSettingKey(feeSettingType)
	return store.Has(key)
}

// CreateFeeSetting check if fee setting existed, else store.
func (k *Keeper) CreateFeeSetting(
	ctx sdkTypes.Context,
	feeSettingMsg MsgSysFeeSetting,
) sdkTypes.Result {
	if !k.IsAuthorised(ctx, feeSettingMsg.Issuer) {
		return sdkTypes.ErrUnknownRequest("Not authorised to create fee setting.").Result()
	}

	return k.storeFeeSetting(ctx, feeSettingMsg)
}

func (k *Keeper) GetFeeSettingByName(ctx sdkTypes.Context, name string) (*FeeSetting, sdkTypes.Error) {
	var msgFeeSetting = new(FeeSetting)

	store := ctx.KVStore(k.key)
	key := getSysFeeSettingKey(name)
	feeSettingData := store.Get(key)

	if feeSettingData == nil {
		return nil, types.ErrFeeSettingNotExists(name)
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(feeSettingData, &msgFeeSetting)

	return msgFeeSetting, nil
}

func (k *Keeper) storeFeeSetting(ctx sdkTypes.Context, msgFeeSetting MsgSysFeeSetting) sdkTypes.Result {
	var feeSetting = new(FeeSetting)
	feeSetting.Name = msgFeeSetting.Name
	feeSetting.Min = msgFeeSetting.Min
	feeSetting.Max = msgFeeSetting.Max
	feeSetting.Issuer = msgFeeSetting.Issuer
	feeSetting.Percentage = msgFeeSetting.Percentage

	store := ctx.KVStore(k.key)
	keyFeeSettingType := getSysFeeSettingKey(feeSetting.Name)
	feeSettingData := k.cdc.MustMarshalBinaryLengthPrefixed(feeSetting)

	store.Set(keyFeeSettingType, feeSettingData)

	eventParam := []string{msgFeeSetting.GetSigners()[0].String(), msgFeeSetting.Name}
	eventSignature := "CreatedFeeSetting(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msgFeeSetting.GetSigners()[0].String(), eventParam),
	}
}

func (k *Keeper) AssignFeeToMsg(
	ctx sdkTypes.Context,
	msgAssignFeeToMsg MsgAssignFeeToMsg,
) sdkTypes.Result {

	if !k.IsAuthorised(ctx, msgAssignFeeToMsg.Issuer) {
		return sdkTypes.ErrUnknownRequest("Not authorised to create msg fee setting.").Result()
	}

	if !k.FeeSettingExists(ctx, msgAssignFeeToMsg.FeeName) {
		return types.ErrFeeSettingNotExists(msgAssignFeeToMsg.FeeName).Result()
	}

	k.assignFeeToMsg(ctx, msgAssignFeeToMsg)

	eventParam := []string{msgAssignFeeToMsg.GetSigners()[0].String(), msgAssignFeeToMsg.MsgType}
	eventSignature := "CreatedTxFeeSetting(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msgAssignFeeToMsg.GetSigners()[0].String(), eventParam),
	}
}

func (k *Keeper) AssignFeeToAcc(
	ctx sdkTypes.Context,
	msgAssignFeeToAcc MsgAssignFeeToAcc,
) sdkTypes.Result {

	if !k.IsAuthorised(ctx, msgAssignFeeToAcc.Issuer) {
		return sdkTypes.ErrUnknownRequest("Not authorised to create msg fee setting.").Result()
	}

	if !k.FeeSettingExists(ctx, msgAssignFeeToAcc.FeeName) {
		return types.ErrFeeSettingNotExists(msgAssignFeeToAcc.FeeName).Result()
	}

	k.assignFeeToAcc(ctx, msgAssignFeeToAcc)

	eventParam := []string{msgAssignFeeToAcc.GetSigners()[0].String(), msgAssignFeeToAcc.Account.String()}
	eventSignature := "CreatedAccountFeeSetting(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msgAssignFeeToAcc.GetSigners()[0].String(), eventParam),
	}
}

func (k *Keeper) AssignFeeToToken(
	ctx sdkTypes.Context,
	msgAssignFeeToToken MsgAssignFeeToToken,
) sdkTypes.Result {

	if !k.IsAuthorised(ctx, msgAssignFeeToToken.Issuer) {
		return sdkTypes.ErrUnknownRequest("Not authorised to create msg fee setting.").Result()
	}

	err := k.AssignFeeToTokenAction(ctx, msgAssignFeeToToken.FeeName, msgAssignFeeToToken.Symbol, msgAssignFeeToToken.Action)
	if err != nil {
		return err.Result()
	}

	eventParam := []string{msgAssignFeeToToken.GetSigners()[0].String(), msgAssignFeeToToken.Symbol}
	eventSignature := "CreatedTokenFeeSetting(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msgAssignFeeToToken.GetSigners()[0].String(), eventParam),
	}
}

func (k *Keeper) CreateMultiplier(
	ctx sdkTypes.Context,
	msgMultiplier MsgMultiplier,
) sdkTypes.Result {

	if !k.IsAuthorised(ctx, msgMultiplier.Issuer) {
		return sdkTypes.ErrUnknownRequest("Not authorised to create msg fee setting.").Result()
	}

	k.storeFeeMultiplier(ctx, msgMultiplier.Multiplier)

	eventParam := []string{msgMultiplier.GetSigners()[0].String()}
	eventSignature := "CreatedFeeMultiplier(string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msgMultiplier.GetSigners()[0].String(), eventParam),
	}
}

func (k *Keeper) CreateTokenMultiplier(
	ctx sdkTypes.Context,
	msgTokenMultiplier MsgTokenMultiplier,
) sdkTypes.Result {

	if !k.IsAuthorised(ctx, msgTokenMultiplier.Issuer) {
		return sdkTypes.ErrUnknownRequest("Not authorised to create msg fee setting.").Result()
	}

	k.storeTokenFeeMultiplier(ctx, msgTokenMultiplier.Multiplier)

	eventParam := []string{msgTokenMultiplier.GetSigners()[0].String()}
	eventSignature := "CreatedTokenFeeMultiplier(string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msgTokenMultiplier.GetSigners()[0].String(), eventParam),
	}
}

func (k *Keeper) assignFeeToMsg(ctx sdkTypes.Context, msg MsgAssignFeeToMsg) {
	store := ctx.KVStore(k.key)
	key := getMsgFeeSettingKey(msg.MsgType)

	store.Set(key, []byte(msg.FeeName))
}

func (k *Keeper) assignFeeToAcc(ctx sdkTypes.Context, msg MsgAssignFeeToAcc) {
	store := ctx.KVStore(k.key)
	key := getAccFeeSettingKey(msg.Account)

	store.Set(key, []byte(msg.FeeName))
}

func (k *Keeper) AssignFeeToTokenAction(ctx sdkTypes.Context, feeName, symbol, action string) sdkTypes.Error {
	if !k.FeeSettingExists(ctx, feeName) {
		return types.ErrFeeSettingNotExists(feeName)
	}

	ok := ContainAction(action)
	if !ok {
		return sdkTypes.ErrUnknownRequest("Token action is not recognize.")
	}

	store := ctx.KVStore(k.key)
	key := getTokenFeeSettingKey(symbol, action)

	store.Set(key, []byte(feeName))

	return nil
}

// Fee Multiplier
func (k *Keeper) storeFeeMultiplier(ctx sdkTypes.Context, multiplier string) {
	store := ctx.KVStore(k.key)
	keyMsgMultiplier := getFeeMultiplierKey()

	store.Set(keyMsgMultiplier, []byte(multiplier))
}

// Token Fee multiplier
func (k *Keeper) storeTokenFeeMultiplier(ctx sdkTypes.Context, multiplier string) {
	store := ctx.KVStore(k.key)
	keyMsgTokenMultiplier := getTokenFeeMultiplierKey()

	store.Set(keyMsgTokenMultiplier, []byte(multiplier))
}

func (k *Keeper) GetMsgFeeSetting(ctx sdkTypes.Context, msgType string) (*FeeSetting, sdkTypes.Error) {
	store := ctx.KVStore(k.key)

	keyMsgType := getMsgFeeSettingKey(msgType)
	feeName := store.Get(keyMsgType)
	if feeName == nil {
		ctx.Logger().Debug("No such tx fee setting. Try to get default fee setting.", "MsgType", msgType)
		feeName = []byte("default")
	}

	feeSetting, err := k.GetFeeSettingByName(ctx, string(feeName))
	if err != nil {
		return nil, err
	}

	return feeSetting, nil
}

func (k *Keeper) GetAccFeeSetting(ctx sdkTypes.Context, acc sdkTypes.AccAddress) (*FeeSetting, sdkTypes.Error) {
	store := ctx.KVStore(k.key)

	keyAccFeeSetting := getAccFeeSettingKey(acc)
	feeName := store.Get(keyAccFeeSetting)
	if feeName != nil {
		feeSetting, err := k.GetFeeSettingByName(ctx, string(feeName))
		if err != nil {
			return nil, err
		}
		return feeSetting, nil
	}

	return nil, nil
}

func (k *Keeper) GetTokenFeeSetting(ctx sdkTypes.Context, tokenSymbol, action string) (*FeeSetting, sdkTypes.Error) {
	store := ctx.KVStore(k.key)

	keyTokenFeeSetting := getTokenFeeSettingKey(tokenSymbol, action)
	feeName := store.Get(keyTokenFeeSetting)
	if feeName == nil {
		return nil, types.ErrTokenFeeSettingNotExists(tokenSymbol)
	}

	feeSetting, err := k.GetFeeSettingByName(ctx, string(feeName))
	if err != nil {
		return nil, err
	}

	return feeSetting, nil
}

func (k *Keeper) GetFeeMultiplier(ctx sdkTypes.Context) (string, sdkTypes.Error) {

	store := ctx.KVStore(k.key)
	key := getFeeMultiplierKey()
	multiplier := store.Get(key)

	if multiplier == nil {
		return string(multiplier), sdkTypes.ErrUnknownRequest(fmt.Sprintf("No multiplier found"))
	}

	return string(multiplier), nil
}

func (k *Keeper) GetTokenFeeMultiplier(ctx sdkTypes.Context) (string, sdkTypes.Error) {

	store := ctx.KVStore(k.key)
	key := getTokenFeeMultiplierKey()
	multiplier := store.Get(key)

	if multiplier == nil {
		return string(multiplier), sdkTypes.ErrUnknownRequest(fmt.Sprintf("No token fee multiplier found"))
	}

	return string(multiplier), nil
}

func (k *Keeper) DeleteFeeSetting(ctx sdkTypes.Context, msgDeleteSysFeeSetting MsgDeleteSysFeeSetting) sdkTypes.Result {

	if k.IsFeeSettingUsed(ctx, msgDeleteSysFeeSetting.Name) {
		return sdkTypes.ErrInternal("Fee setting is in used, delete failed.").Result()
	}

	store := ctx.KVStore(k.key)
	key := getSysFeeSettingKey(msgDeleteSysFeeSetting.Name)

	store.Delete(key)

	eventParam := []string{msgDeleteSysFeeSetting.GetSigners()[0].String(), msgDeleteSysFeeSetting.Name}
	eventSignature := "DeletedFeeSetting(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msgDeleteSysFeeSetting.GetSigners()[0].String(), eventParam),
	}
}

func (k *Keeper) DeleteAccFeeSetting(ctx sdkTypes.Context, msgDeleteAccFeeSetting MsgDeleteAccFeeSetting) sdkTypes.Result {

	accFeeSetting, err := k.GetAccFeeSetting(ctx, msgDeleteAccFeeSetting.Account)
	if err != nil {
		return err.Result()
	}

	if accFeeSetting == nil {
		return sdkTypes.ErrInternal("Account Fee setting is not set, delete failed.").Result()
	}

	store := ctx.KVStore(k.key)
	key := getAccFeeSettingKey(msgDeleteAccFeeSetting.Account)

	store.Delete(key)

	eventParam := []string{msgDeleteAccFeeSetting.GetSigners()[0].String(), msgDeleteAccFeeSetting.Account.String()}
	eventSignature := "DeletedAccountFeeSetting(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msgDeleteAccFeeSetting.GetSigners()[0].String(), eventParam),
	}
}

// get list of sysfeesetting
func (k *Keeper) ListAllSysFeeSetting(ctx sdkTypes.Context) []FeeSetting {

	store := ctx.KVStore(k.key)
	start := append(prefixSysFeeSetting, 0x00)
	end := append(prefixSysFeeSetting, 0xFF)
	iter := store.Iterator(start, end)
	defer iter.Close()

	var lst = make([]FeeSetting, 0)

	for {
		if !iter.Valid() {
			break
		}
		var feeSetting = new(FeeSetting)
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &feeSetting)
		lst = append(lst, *feeSetting)

		iter.Next()
	}
	return lst
}

// get list of tokenfeesetting
func (k *Keeper) ListTokenFeeSetting(ctx sdkTypes.Context) []FeeSetting {

	store := ctx.KVStore(k.key)
	start := append(prefixTokenFeeSetting, 0x06)
	end := append(prefixTokenFeeSetting, 0xFF)
	iter := store.Iterator(start, end)
	defer iter.Close()

	var lst = make([]FeeSetting, 0)

	for {
		if !iter.Valid() {
			break
		}
		var TokenfeeSetting = new(FeeSetting)
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &TokenfeeSetting)
		lst = append(lst, *TokenfeeSetting)

		iter.Next()
	}
	return lst
}

// get list of accfeesetting
func (k *Keeper) ListAccFeeSetting(ctx sdkTypes.Context) []FeeSetting {

	store := ctx.KVStore(k.key)
	start := append(prefixAccFeeSetting, 0x05)
	end := append(prefixAccFeeSetting, 0xFF)
	iter := store.Iterator(start, end)
	defer iter.Close()

	var lst = make([]FeeSetting, 0)

	for {
		if !iter.Valid() {
			break
		}
		var AccfeeSetting = new(FeeSetting)
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &AccfeeSetting)
		lst = append(lst, *AccfeeSetting)

		iter.Next()
	}
	return lst
}

func (k *Keeper) IsFeeSettingUsed(ctx sdkTypes.Context, feeName string) bool {

	// default created upon genesis.
	if feeName == "default" {
		return true
	}

	store := ctx.KVStore(k.key)
	start := []byte("0x04")
	end := []byte("0x07")
	iter := store.Iterator(start, end)
	defer iter.Close()

	for {
		if !iter.Valid() {
			break
		}

		if string(iter.Value()) == feeName {
			return true
		}

		iter.Next()
	}
	return false
}

func ContainAction(tokenAction string) bool {
	for _, action := range tokenActions {
		if tokenAction == action {
			return true
		}
	}
	return false
}
