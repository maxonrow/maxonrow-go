package nonfungible

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	MsgRoute                                 = "nonFungible"
	MsgTypeCreateNonFungibleToken            = "createNonFungibleToken"
	MsgTypeSetNonFungibleTokenStatus         = "setNonFungibleTokenStatus"
	MsgTypeTransferNonFungibleItem           = "transferNonFungibleItem"
	MsgTypeMintNonFungibleItem               = "mintNonFungibleItem"
	MsgTypeBurnNonFungibleItem               = "burnNonFungibleItem"
	MsgTypeTransferNonFungibleTokenOwnership = "transferNonFungibleTokenOwnership"
	MsgTypeAcceptNonFungibleTokenOwnership   = "acceptNonFungibleTokenOwnership"
	MsgTypeSetNonFungibleItemStatus          = "setNonFungibleItemStatus"
	MsgTypeEndorsement                       = "endorsement"
	MsgTypeUpdateItemMetadata                = "updateItemMetadata"
	MsgTypeUpdateNFTMetadata                 = "updateNFTMetadata"
	MsgTypeUpdateEndorserList                = "updateNFTEndorserList"
)

// MsgCreateNonFungibleToken
type MsgCreateNonFungibleToken struct {
	Name       string              `json:"name"`
	Symbol     string              `json:"symbol"`
	Properties string              `json:"properties"`
	Metadata   string              `json:"metadata"`
	Owner      sdkTypes.AccAddress `json:"owner"`
	Fee        Fee                 `json:"fee"`
}

type Fee struct {
	To    sdkTypes.AccAddress `json:"to"`
	Value string              `json:"value"`
}

func NewMsgCreateNonFungibleToken(symbol string, owner sdkTypes.AccAddress, name string, properties string, metadata string, fee Fee) *MsgCreateNonFungibleToken {
	return &MsgCreateNonFungibleToken{
		Name:       name,
		Symbol:     symbol,
		Properties: properties,
		Metadata:   metadata,
		Owner:      owner,
		Fee:        fee,
	}
}

func NewFee(to sdkTypes.AccAddress, value string) Fee {
	return Fee{
		To:    to,
		Value: value,
	}
}

func (msg MsgCreateNonFungibleToken) Route() string {
	return MsgRoute
}

func (msg MsgCreateNonFungibleToken) Type() string {
	return MsgTypeCreateNonFungibleToken
}

func (msg MsgCreateNonFungibleToken) ValidateBasic() sdkTypes.Error {
	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}

	if err := validateTokenName(msg.Name); err != nil {
		return err
	}

	if err := validateMetadata(msg.Metadata); err != nil {
		return err
	}

	if err := validateProperties(msg.Properties); err != nil {
		return err
	}

	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	return nil
}

func (msg MsgCreateNonFungibleToken) GetSignBytes() []byte {

	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateNonFungibleToken) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
}

// MsgSetFungibleTokenStatus
type MsgSetNonFungibleTokenStatus struct {
	Owner      sdkTypes.AccAddress `json:"owner"`
	Payload    Payload             `json:"payload"`
	Signatures []Signature         `json:"signatures"`
}

type Payload struct {
	Token         TokenData `json:"token"`
	crypto.PubKey `json:"pub_key"`
	Signature     []byte `json:"signature"`
}

type TokenData struct {
	From          sdkTypes.AccAddress   `json:"from"`
	Nonce         string                `json:"nonce"`
	Status        string                `json:"status"`
	Symbol        string                `json:"symbol"`
	TransferLimit sdkTypes.Uint         `json:"transferLimit"`
	MintLimit     sdkTypes.Uint         `json:"mintLimit"`
	Burnable      bool                  `json:"burnable"`
	Transferable  bool                  `json:"transferable"`
	Modifiable    bool                  `json:"modifiable"`
	Public        bool                  `json:"pub"`
	TokenFees     []TokenFee            `json:"tokenFees,omitempty"`
	EndorserList  []sdkTypes.AccAddress `json:"endorserList"`
}

type TokenFee struct {
	Action  string `json:"action"`
	FeeName string `json:"feeName"`
}

func NewMsgSetNonFungibleTokenStatus(owner sdkTypes.AccAddress, payload Payload, signatures []Signature) *MsgSetNonFungibleTokenStatus {
	return &MsgSetNonFungibleTokenStatus{
		Owner:      owner,
		Payload:    payload,
		Signatures: signatures,
	}
}

func NewPayload(token TokenData, pubKey crypto.PubKey, signature []byte) *Payload {
	return &Payload{
		Token:     token,
		PubKey:    pubKey,
		Signature: signature,
	}
}

func NewToken(from sdkTypes.AccAddress, nonce, status, symbol string, transferLimit, mintLimit sdkTypes.Uint, tokenFees []TokenFee, endorserList []sdkTypes.AccAddress, burnable, transferable, modifiable, public bool) *TokenData {
	return &TokenData{
		From:          from,
		Nonce:         nonce,
		Status:        status,
		Symbol:        symbol,
		TransferLimit: transferLimit,
		MintLimit:     mintLimit,
		Burnable:      burnable,
		Transferable:  transferable,
		Modifiable:    modifiable,
		Public:        public,
		TokenFees:     tokenFees,
		EndorserList:  endorserList,
	}
}

func (msg MsgSetNonFungibleTokenStatus) Route() string {
	return MsgRoute
}

func (msg MsgSetNonFungibleTokenStatus) Type() string {
	return MsgTypeSetNonFungibleTokenStatus
}

func (msg MsgSetNonFungibleTokenStatus) ValidateBasic() sdkTypes.Error {
	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}

	if err := ValidateSymbol(msg.Payload.Token.Symbol); err != nil {
		return err
	}

	if len(msg.Signatures) < 1 {
		return sdkTypes.ErrInvalidAddress("Insufficient issuer signature.")
	}

	if msg.Payload.Token.From.Empty() {
		return sdkTypes.ErrInvalidAddress("From cannot be empty.")
	}

	if msg.Payload.Token.Status == ApproveToken && msg.Payload.Token.TokenFees == nil {
		return sdkTypes.ErrUnknownRequest("Approve token, token fees cannot be empty.")
	}

	return nil
}

func (payload Payload) GetIssuerSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(payload)
	if err != nil {
		panic(err)
	}

	return sdkTypes.MustSortJSON(b)
}

func (token TokenData) GetFromSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(token)
	if err != nil {
		panic(err)
	}

	return sdkTypes.MustSortJSON(b)
}

func (msg MsgSetNonFungibleTokenStatus) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgSetNonFungibleTokenStatus) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
}

// MintFungibleToken - only for token without fixed supply
type MsgMintNonFungibleItem struct {
	ItemID     string              `json:"itemID"`
	Symbol     string              `json:"symbol"`
	Owner      sdkTypes.AccAddress `json:"owner"`
	To         sdkTypes.AccAddress `json:"to"`
	Properties string              `json:"properties"`
	Metadata   string              `json:"metadata"`
}

func NewMsgMintNonFungibleItem(owner sdkTypes.AccAddress, symbol string, to sdkTypes.AccAddress, itemID, properties, metadata string) *MsgMintNonFungibleItem {
	return &MsgMintNonFungibleItem{
		ItemID:     itemID,
		Symbol:     symbol,
		To:         to,
		Owner:      owner,
		Properties: properties,
		Metadata:   metadata,
	}
}

func (msg MsgMintNonFungibleItem) Route() string {
	return MsgRoute
}

func (msg MsgMintNonFungibleItem) Type() string {
	return MsgTypeMintNonFungibleItem
}

func (msg MsgMintNonFungibleItem) ValidateBasic() sdkTypes.Error {
	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}

	if msg.To.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.To.String())
	}

	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	if len(msg.ItemID) < 1 {
		return sdkTypes.ErrInternal("Item id not allow to be empty.")
	}

	return nil
}

func (msg MsgMintNonFungibleItem) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgMintNonFungibleItem) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
}

// TransferFungibleToken
type MsgTransferNonFungibleItem struct {
	Symbol string              `json:"symbol"`
	From   sdkTypes.AccAddress `json:"from"`
	To     sdkTypes.AccAddress `json:"to"`
	ItemID string              `json:"itemID"`
}

func NewMsgTransferNonFungibleItem(symbol string, from, to sdkTypes.AccAddress, itemID string) *MsgTransferNonFungibleItem {
	return &MsgTransferNonFungibleItem{
		Symbol: symbol,
		From:   from,
		To:     to,
		ItemID: itemID,
	}
}

func (msg MsgTransferNonFungibleItem) Route() string {
	return MsgRoute
}

func (msg MsgTransferNonFungibleItem) Type() string {
	return MsgTypeTransferNonFungibleItem
}

func (msg MsgTransferNonFungibleItem) ValidateBasic() sdkTypes.Error {
	if msg.From.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.From.String())
	}

	if msg.To.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.To.String())
	}

	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	if err := ValidateItemID(msg.ItemID); err != nil {
		return err
	}

	return nil
}

func (msg MsgTransferNonFungibleItem) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferNonFungibleItem) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.From}
}

// BurnFungibleToken - only for fixed supply fungible token
type MsgBurnNonFungibleItem struct {
	Symbol string              `json:"symbol"`
	From   sdkTypes.AccAddress `json:"from"`
	ItemID string              `json:"itemID"`
}

func NewMsgBurnNonFungibleItem(symbol string, from sdkTypes.AccAddress, itemID string) *MsgBurnNonFungibleItem {
	return &MsgBurnNonFungibleItem{
		Symbol: symbol,
		From:   from,
		ItemID: itemID,
	}
}

func (msg MsgBurnNonFungibleItem) Route() string {
	return MsgRoute
}

func (msg MsgBurnNonFungibleItem) Type() string {
	return MsgTypeBurnNonFungibleItem
}

func (msg MsgBurnNonFungibleItem) ValidateBasic() sdkTypes.Error {
	if msg.From.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.From.String())
	}

	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	return nil
}

func (msg MsgBurnNonFungibleItem) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnNonFungibleItem) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.From}
}

// MsgSetFungibleTokenAccountStatus
type MsgSetNonFungibleItemStatus struct {
	Owner       sdkTypes.AccAddress `json:"owner"`
	ItemPayload ItemPayload         `json:"payload"`
	Signatures  []Signature         `json:"signatures"`
}

type ItemPayload struct {
	Item          ItemDetails `json:"item"`
	crypto.PubKey `json:"pub_key"`
	Signature     []byte `json:"signature"`
}

type ItemDetails struct {
	From   sdkTypes.AccAddress `json:"from"`
	Nonce  string              `json:"nonce"`
	Status string              `json:"status"`
	Symbol string              `json:"symbol"`
	ItemID string              `json:"itemID"`
}

func NewMsgSetNonFungibleItemStatus(owner sdkTypes.AccAddress, itemPayload ItemPayload, signatures []Signature) *MsgSetNonFungibleItemStatus {
	return &MsgSetNonFungibleItemStatus{
		Owner:       owner,
		ItemPayload: itemPayload,
		Signatures:  signatures,
	}
}

func (msg MsgSetNonFungibleItemStatus) Route() string {
	return MsgRoute
}

func (msg MsgSetNonFungibleItemStatus) Type() string {
	return MsgTypeSetNonFungibleItemStatus
}

func (msg MsgSetNonFungibleItemStatus) ValidateBasic() sdkTypes.Error {
	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}

	if msg.ItemPayload.Item.From.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.ItemPayload.Item.From.String())
	}

	if err := ValidateSymbol(msg.ItemPayload.Item.Symbol); err != nil {
		return err
	}

	return nil
}

func (itemPayload ItemPayload) GetAccountStatusSettingSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(itemPayload)
	if err != nil {
		panic(err)
	}

	return sdkTypes.MustSortJSON(b)
}

func (item ItemDetails) GetAccountStatusSettingFromSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(item)
	if err != nil {
		panic(err)
	}

	return sdkTypes.MustSortJSON(b)
}

func (msg MsgSetNonFungibleItemStatus) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgSetNonFungibleItemStatus) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
}

// TransferFungibleTokenOwnership
type MsgTransferNonFungibleTokenOwnership struct {
	Symbol string              `json:"symbol"`
	From   sdkTypes.AccAddress `json:"from"`
	To     sdkTypes.AccAddress `json:"to"`
}

func NewMsgTransferNonFungibleTokenOwnership(symbol string, from sdkTypes.AccAddress, to sdkTypes.AccAddress) *MsgTransferNonFungibleTokenOwnership {
	return &MsgTransferNonFungibleTokenOwnership{
		Symbol: symbol,
		From:   from,
		To:     to,
	}
}

func (msg MsgTransferNonFungibleTokenOwnership) Route() string {
	return MsgRoute
}

func (msg MsgTransferNonFungibleTokenOwnership) Type() string {
	return MsgTypeTransferNonFungibleTokenOwnership
}

func (msg MsgTransferNonFungibleTokenOwnership) ValidateBasic() sdkTypes.Error {
	if msg.From.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.From.String())
	}

	if msg.To.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.To.String())
	}

	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	return nil
}

func (msg MsgTransferNonFungibleTokenOwnership) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferNonFungibleTokenOwnership) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.From}
}

// MsgAcceptFungibleTokenOwnership
type MsgAcceptNonFungibleTokenOwnership struct {
	Symbol string              `json:"symbol"`
	From   sdkTypes.AccAddress `json:"from"`
}

func NewMsgAcceptNonFungibleTokenOwnership(symbol string, from sdkTypes.AccAddress) *MsgAcceptNonFungibleTokenOwnership {
	return &MsgAcceptNonFungibleTokenOwnership{
		Symbol: symbol,
		From:   from,
	}
}

func (msg MsgAcceptNonFungibleTokenOwnership) Route() string {
	return MsgRoute
}

func (msg MsgAcceptNonFungibleTokenOwnership) Type() string {
	return MsgTypeAcceptNonFungibleTokenOwnership
}

func (msg MsgAcceptNonFungibleTokenOwnership) ValidateBasic() sdkTypes.Error {
	if msg.From.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.From.String())
	}

	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	return nil
}

func (msg MsgAcceptNonFungibleTokenOwnership) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgAcceptNonFungibleTokenOwnership) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.From}
}

type MsgEndorsement struct {
	Symbol   string              `json:"symbol"`
	From     sdkTypes.AccAddress `json:"from"`
	ItemID   string              `json:"itemID"`
	Metadata string              `json:"metadata"`
}

func NewMsgEndorsement(symbol string, from sdkTypes.AccAddress, itemID string, metadata string) *MsgEndorsement {
	return &MsgEndorsement{
		Symbol:   symbol,
		From:     from,
		ItemID:   itemID,
		Metadata: metadata,
	}
}

func (msg MsgEndorsement) Route() string {
	return MsgRoute
}

func (msg MsgEndorsement) Type() string {
	return MsgTypeEndorsement
}

func (msg MsgEndorsement) ValidateBasic() sdkTypes.Error {
	if msg.From.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.From.String())
	}

	if len(msg.ItemID) < 1 {
		return sdkTypes.ErrInternal("Item id can not be empty.")
	}

	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	if err := validateMetadata(msg.Metadata); err != nil {
		return err
	}

	return nil
}

func (msg MsgEndorsement) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgEndorsement) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.From}
}

type MsgUpdateItemMetadata struct {
	Symbol   string              `json:"symbol"`
	From     sdkTypes.AccAddress `json:"from"`
	ItemID   string              `json:"itemID"`
	Metadata string              `json:"metadata"`
}

func NewMsgUpdateItemMetadata(symbol string, from sdkTypes.AccAddress, itemID string, metadata string) *MsgUpdateItemMetadata {
	return &MsgUpdateItemMetadata{
		Symbol:   symbol,
		From:     from,
		ItemID:   itemID,
		Metadata: metadata,
	}
}

func (msg MsgUpdateItemMetadata) Route() string {
	return MsgRoute
}

func (msg MsgUpdateItemMetadata) Type() string {
	return MsgTypeUpdateItemMetadata
}

func (msg MsgUpdateItemMetadata) ValidateBasic() sdkTypes.Error {

	if msg.From.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.From.String())
	}

	if len(msg.ItemID) < 1 {
		return sdkTypes.ErrInternal("Item id can not be empty.")
	}
	if err := validateMetadata(msg.Metadata); err != nil {
		return err
	}

	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	return nil
}

func (msg MsgUpdateItemMetadata) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgUpdateItemMetadata) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.From}
}

type MsgUpdateNFTMetadata struct {
	Symbol   string              `json:"symbol"`
	From     sdkTypes.AccAddress `json:"from"`
	Metadata string              `json:"metadata"`
}

func NewMsgUpdateNFTMetadata(symbol string, from sdkTypes.AccAddress, metadata string) *MsgUpdateNFTMetadata {
	return &MsgUpdateNFTMetadata{
		Symbol:   symbol,
		From:     from,
		Metadata: metadata,
	}
}

func (msg MsgUpdateNFTMetadata) Route() string {
	return MsgRoute
}

func (msg MsgUpdateNFTMetadata) Type() string {
	return MsgTypeUpdateNFTMetadata
}

func (msg MsgUpdateNFTMetadata) ValidateBasic() sdkTypes.Error {
	if msg.From.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.From.String())
	}

	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}
	if err := validateMetadata(msg.Metadata); err != nil {
		return err
	}

	return nil
}

func (msg MsgUpdateNFTMetadata) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgUpdateNFTMetadata) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.From}
}

type MsgUpdateEndorserList struct {
	Symbol    string                `json:"symbol"`
	From      sdkTypes.AccAddress   `json:"from"`
	Endorsers []sdkTypes.AccAddress `json:"endorsers"`
}

func NewMsgUpdateEndorserList(symbol string, from sdkTypes.AccAddress, endorsers []sdkTypes.AccAddress) *MsgUpdateEndorserList {
	return &MsgUpdateEndorserList{
		Symbol:    symbol,
		From:      from,
		Endorsers: endorsers,
	}
}

func (msg MsgUpdateEndorserList) Route() string {
	return MsgRoute
}

func (msg MsgUpdateEndorserList) Type() string {
	return MsgTypeUpdateEndorserList
}

func (msg MsgUpdateEndorserList) ValidateBasic() sdkTypes.Error {
	if msg.From.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.From.String())
	}

	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	return nil
}

func (msg MsgUpdateEndorserList) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgUpdateEndorserList) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.From}
}

// Signature structure
type Signature struct {
	crypto.PubKey `json:"pub_key"`
	Signature     []byte `json:"signature"`
}

func NewSignature(pubKey crypto.PubKey, signature []byte) Signature {
	return Signature{
		PubKey:    pubKey,
		Signature: signature,
	}
}

func NewItemPayload(itemDetails ItemDetails, pubKey crypto.PubKey, signature []byte) *ItemPayload {
	return &ItemPayload{
		Item:      itemDetails,
		PubKey:    pubKey,
		Signature: signature,
	}

}

func NewItemDetails(from sdkTypes.AccAddress, nonce string, status string, symbol string, itemID string) *ItemDetails {
	return &ItemDetails{
		From:   from,
		Nonce:  nonce,
		Status: status,
		Symbol: symbol,
		ItemID: itemID,
	}

}
