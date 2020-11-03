package fungible

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	MsgRoute                              = "token"
	MsgTypeCreateFungibleToken            = "createFungibleToken"
	MsgTypeSetFungibleTokenStatus         = "setFungibleTokenStatus"
	MsgTypeTransferFungibleToken          = "transferFungibleToken"
	MsgTypeMintFungibleToken              = "mintFungibleToken"
	MsgTypeBurnFungibleToken              = "burnFungibleToken"
	MsgTypeTransferFungibleTokenOwnership = "transferFungibleTokenOwnership"
	MsgTypeAcceptFungibleTokenOwnership   = "acceptFungibleTokenOwnership"
	MsgTypeSetFungibleTokenAccountStatus  = "setFungibleTokenAccountStatus"
)

const (
	TransferFungibleToken  = "transferFungibleToken"
	MintFungibleToken      = "mintFungibleToken"
	BurnFungibleToken      = "burnFungibleToken"
	TransferTokenOwnership = "transferFungibleTokenOwnership"
	AcceptTokenOwnership   = "acceptFungibleTokenOwnership"
)

// MsgCreateFungibleToken
type MsgCreateFungibleToken struct {
	Name        string              `json:"name"`
	Symbol      string              `json:"symbol"`
	Decimals    int                 `json:"decimals"`
	Metadata    string              `json:"metadata"`
	FixedSupply bool                `json:"fixedSupply"`
	Owner       sdkTypes.AccAddress `json:"owner"`
	MaxSupply   sdkTypes.Uint       `json:"maxSupply"`
	Fee         Fee                 `json:"fee"`
}

type Fee struct {
	To    sdkTypes.AccAddress `json:"to"`
	Value string              `json:"value"`
}

func NewMsgCreateFungibleToken(symbol string, decimals int, owner sdkTypes.AccAddress, name string, fixedSupply bool, maxSupply sdkTypes.Uint, metadata string, fee Fee) *MsgCreateFungibleToken {
	return &MsgCreateFungibleToken{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		Metadata:    metadata,
		FixedSupply: fixedSupply,
		Owner:       owner,
		MaxSupply:   maxSupply,
		Fee:         fee,
	}
}

func NewFee(to sdkTypes.AccAddress, value string) Fee {
	return Fee{
		To:    to,
		Value: value,
	}
}

func (msg MsgCreateFungibleToken) Route() string {
	return MsgRoute
}

func (msg MsgCreateFungibleToken) Type() string {
	return MsgTypeCreateFungibleToken
}

func (msg MsgCreateFungibleToken) ValidateBasic() sdkTypes.Error {
	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}

	if err := validateTokenName(msg.Name); err != nil {
		return err
	}

	if err := validateMetadata(msg.Metadata); err != nil {
		return err
	}

	if err := validateSymbol(msg.Symbol); err != nil {
		return err
	}

	zero := sdkTypes.NewUintFromString("0")
	if msg.FixedSupply {
		if msg.MaxSupply.IsZero() || msg.MaxSupply.LTE(zero) {
			return sdkTypes.ErrUnknownRequest("Cannot have max supply 0 or less than 0 for fixed supply token")
		}
	} else {
		if msg.MaxSupply.LT(zero) {
			return sdkTypes.ErrUnknownRequest("Cannot have max supply less than 0 for dynamic supply token")
		}
	}

	if err := validateDecimal(msg.Decimals); err != nil {
		return err
	}

	return nil
}

func (msg MsgCreateFungibleToken) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateFungibleToken) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
}

// MsgSetFungibleTokenStatus
type MsgSetFungibleTokenStatus struct {
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
	From      sdkTypes.AccAddress `json:"from"`
	Nonce     string              `json:"nonce"`
	Status    string              `json:"status"`
	Symbol    string              `json:"symbol"`
	Burnable  bool                `json:"burnable"`
	TokenFees []TokenFee          `json:"tokenFees,omitempty"`
}

type TokenFee struct {
	Action  string `json:"action"`
	FeeName string `json:"feeName"`
}

func NewMsgSetFungibleTokenStatus(owner sdkTypes.AccAddress, payload Payload, signatures []Signature) *MsgSetFungibleTokenStatus {
	return &MsgSetFungibleTokenStatus{
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

func NewToken(from sdkTypes.AccAddress, nonce, status, symbol string, burnable bool, tokenFees []TokenFee) *TokenData {
	return &TokenData{
		From:      from,
		Nonce:     nonce,
		Status:    status,
		Symbol:    symbol,
		Burnable:  burnable,
		TokenFees: tokenFees,
	}
}

func (msg MsgSetFungibleTokenStatus) Route() string {
	return MsgRoute
}

func (msg MsgSetFungibleTokenStatus) Type() string {
	return MsgTypeSetFungibleTokenStatus
}

func (msg MsgSetFungibleTokenStatus) ValidateBasic() sdkTypes.Error {
	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}

	if err := validateSymbol(msg.Payload.Token.Symbol); err != nil {
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

func (msg MsgSetFungibleTokenStatus) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgSetFungibleTokenStatus) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
}

// MintFungibleToken - only for token without fixed supply
type MsgMintFungibleToken struct {
	Symbol string              `json:"symbol"`
	Value  sdkTypes.Uint       `json:"value"`
	Owner  sdkTypes.AccAddress `json:"owner"`
	To     sdkTypes.AccAddress `json:"to"`
}

func NewMsgIssueFungibleAsset(owner sdkTypes.AccAddress, symbol string, to sdkTypes.AccAddress, value sdkTypes.Uint) *MsgMintFungibleToken {
	return &MsgMintFungibleToken{
		Symbol: symbol,
		Value:  value,
		To:     to,
		Owner:  owner,
	}
}

func (msg MsgMintFungibleToken) Route() string {
	return MsgRoute
}

func (msg MsgMintFungibleToken) Type() string {
	return MsgTypeMintFungibleToken
}

func (msg MsgMintFungibleToken) ValidateBasic() sdkTypes.Error {
	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}

	if msg.To.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.To.String())
	}

	if err := validateSymbol(msg.Symbol); err != nil {
		return err
	}

	zero := sdkTypes.NewUintFromString("0")

	if msg.Value.LT(zero) {
		return sdkTypes.ErrUnknownRequest("Value cannot be empty.")
	}

	return nil
}

func (msg MsgMintFungibleToken) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgMintFungibleToken) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
}

// TransferFungibleToken
type MsgTransferFungibleToken struct {
	Symbol string              `json:"symbol"`
	Value  sdkTypes.Uint       `json:"value"`
	From   sdkTypes.AccAddress `json:"from"`
	To     sdkTypes.AccAddress `json:"to"`
}

func NewMsgTransferFungibleToken(symbol string, value sdkTypes.Uint, from, to sdkTypes.AccAddress) *MsgTransferFungibleToken {
	return &MsgTransferFungibleToken{
		Symbol: symbol,
		Value:  value,
		From:   from,
		To:     to,
	}
}

func (msg MsgTransferFungibleToken) Route() string {
	return MsgRoute
}

func (msg MsgTransferFungibleToken) Type() string {
	return MsgTypeTransferFungibleToken
}

func (msg MsgTransferFungibleToken) ValidateBasic() sdkTypes.Error {
	if msg.From.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.From.String())
	}

	if msg.To.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.To.String())
	}

	if err := validateSymbol(msg.Symbol); err != nil {
		return err
	}

	zero := sdkTypes.NewUintFromString("0")

	if msg.Value.LT(zero) {
		return sdkTypes.ErrUnknownRequest("Value cannot be empty.")
	}

	return nil
}

func (msg MsgTransferFungibleToken) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferFungibleToken) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.From}
}

// BurnFungibleToken - only for fixed supply fungible token
type MsgBurnFungibleToken struct {
	Symbol string              `json:"symbol"`
	Value  sdkTypes.Uint       `json:"value"`
	From   sdkTypes.AccAddress `json:"from"`
}

func NewMsgBurnFungibleToken(symbol string, value sdkTypes.Uint, from sdkTypes.AccAddress) *MsgBurnFungibleToken {
	return &MsgBurnFungibleToken{
		Symbol: symbol,
		Value:  value,
		From:   from,
	}
}

func (msg MsgBurnFungibleToken) Route() string {
	return MsgRoute
}

func (msg MsgBurnFungibleToken) Type() string {
	return MsgTypeBurnFungibleToken
}

func (msg MsgBurnFungibleToken) ValidateBasic() sdkTypes.Error {
	if msg.From.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.From.String())
	}

	if err := validateSymbol(msg.Symbol); err != nil {
		return err
	}

	zero := sdkTypes.NewUintFromString("0")

	if msg.Value.LT(zero) {
		return sdkTypes.ErrUnknownRequest("Value cannot be empty.")
	}

	return nil
}

func (msg MsgBurnFungibleToken) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnFungibleToken) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.From}
}

// MsgSetFungibleTokenAccountStatus
type MsgSetFungibleTokenAccountStatus struct {
	Owner               sdkTypes.AccAddress `json:"owner"`
	TokenAccountPayload TokenAccountPayload `json:"payload"`
	Signatures          []Signature         `json:"signatures"`
}

type TokenAccountPayload struct {
	TokenAccount  TokenAccount `json:"tokenAccount"`
	crypto.PubKey `json:"pub_key"`
	Signature     []byte `json:"signature"`
}

type TokenAccount struct {
	From    sdkTypes.AccAddress `json:"from"`
	Nonce   string              `json:"nonce"`
	Status  string              `json:"status"`
	Symbol  string              `json:"symbol"`
	Account sdkTypes.AccAddress `json:"to"`
}

func NewMsgSetFungibleTokenAccountStatus(owner sdkTypes.AccAddress, tokenAccountPayload TokenAccountPayload, signatures []Signature) *MsgSetFungibleTokenAccountStatus {
	return &MsgSetFungibleTokenAccountStatus{
		Owner:               owner,
		TokenAccountPayload: tokenAccountPayload,
		Signatures:          signatures,
	}
}

func (msg MsgSetFungibleTokenAccountStatus) Route() string {
	return MsgRoute
}

func (msg MsgSetFungibleTokenAccountStatus) Type() string {
	return MsgTypeSetFungibleTokenAccountStatus
}

func (msg MsgSetFungibleTokenAccountStatus) ValidateBasic() sdkTypes.Error {
	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}

	if msg.TokenAccountPayload.TokenAccount.From.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.TokenAccountPayload.TokenAccount.From.String())
	}

	if err := validateSymbol(msg.TokenAccountPayload.TokenAccount.Symbol); err != nil {
		return err
	}

	return nil
}

func (tokenAccountPayload TokenAccountPayload) GetAccountStatusSettingSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(tokenAccountPayload)
	if err != nil {
		panic(err)
	}

	return sdkTypes.MustSortJSON(b)
}

func (tokenAccount TokenAccount) GetAccountStatusSettingFromSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(tokenAccount)
	if err != nil {
		panic(err)
	}

	return sdkTypes.MustSortJSON(b)
}

func (msg MsgSetFungibleTokenAccountStatus) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgSetFungibleTokenAccountStatus) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
}

// TransferFungibleTokenOwnership
type MsgTransferFungibleTokenOwnership struct {
	Symbol string              `json:"symbol"`
	From   sdkTypes.AccAddress `json:"from"`
	To     sdkTypes.AccAddress `json:"to"`
}

func NewMsgTransferFungibleTokenOwnership(symbol string, from sdkTypes.AccAddress, to sdkTypes.AccAddress) *MsgTransferFungibleTokenOwnership {
	return &MsgTransferFungibleTokenOwnership{
		Symbol: symbol,
		From:   from,
		To:     to,
	}
}

func (msg MsgTransferFungibleTokenOwnership) Route() string {
	return MsgRoute
}

func (msg MsgTransferFungibleTokenOwnership) Type() string {
	return MsgTypeTransferFungibleTokenOwnership
}

func (msg MsgTransferFungibleTokenOwnership) ValidateBasic() sdkTypes.Error {
	if msg.From.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.From.String())
	}

	if msg.To.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.To.String())
	}

	if err := validateSymbol(msg.Symbol); err != nil {
		return err
	}

	return nil
}

func (msg MsgTransferFungibleTokenOwnership) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferFungibleTokenOwnership) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.From}
}

// MsgAcceptFungibleTokenOwnership
type MsgAcceptFungibleTokenOwnership struct {
	Symbol string              `json:"symbol"`
	From   sdkTypes.AccAddress `json:"from"`
}

func NewMsgAcceptFungibleTokenOwnership(symbol string, from sdkTypes.AccAddress) *MsgAcceptFungibleTokenOwnership {
	return &MsgAcceptFungibleTokenOwnership{
		Symbol: symbol,
		From:   from,
	}
}

func (msg MsgAcceptFungibleTokenOwnership) Route() string {
	return MsgRoute
}

func (msg MsgAcceptFungibleTokenOwnership) Type() string {
	return MsgTypeAcceptFungibleTokenOwnership
}

func (msg MsgAcceptFungibleTokenOwnership) ValidateBasic() sdkTypes.Error {
	if msg.From.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.From.String())
	}

	if err := validateSymbol(msg.Symbol); err != nil {
		return err
	}

	return nil
}

func (msg MsgAcceptFungibleTokenOwnership) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgAcceptFungibleTokenOwnership) GetSigners() []sdkTypes.AccAddress {
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
