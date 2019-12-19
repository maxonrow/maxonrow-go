package fee

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const routeName = "fee"

// Notes
// Many messages contain fields like Owner or Approver that are a bit redundant - they could
// be derived from the message signer itself, but the CosmosSDK structure prevents this

// MsgSysFeeSetting Create Fee Setting
type MsgSysFeeSetting struct {
	Name       string              `json:"name"`
	Min        sdkTypes.Coins      `json:"min"`
	Max        sdkTypes.Coins      `json:"max"`
	Percentage string              `json:"percentage"`
	Issuer     sdkTypes.AccAddress `json:"issuer"`
}

func NewMsgSysFeeSetting(name string, min sdkTypes.Coins, max sdkTypes.Coins, percentage string, issuer sdkTypes.AccAddress) MsgSysFeeSetting {
	return MsgSysFeeSetting{
		Min:        min,
		Max:        max,
		Percentage: percentage,
		Issuer:     issuer,
		Name:       name,
	}
}

func (msg MsgSysFeeSetting) Route() string {
	return routeName
}

func (msg MsgSysFeeSetting) Type() string {
	return "updateFeeSetting"
}

func (msg MsgSysFeeSetting) ValidateBasic() sdkTypes.Error {
	if msg.Issuer.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Issuer.String())
	}

	if msg.Max.IsAllLT(msg.Min) {
		return sdkTypes.ErrInvalidCoins("Max fee cannot lower than minimum fee.")
	}

	return nil
}

func (msg MsgSysFeeSetting) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgSysFeeSetting) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Issuer}
}

type MsgDeleteSysFeeSetting struct {
	Name   string              `json:"name"`
	Issuer sdkTypes.AccAddress `json:"issuer"`
}

func NewMsgDeleteSysFeeSetting(name string, issuer sdkTypes.AccAddress) MsgDeleteSysFeeSetting {
	return MsgDeleteSysFeeSetting{
		Issuer: issuer,
		Name:   name,
	}
}

func (msg MsgDeleteSysFeeSetting) Route() string {
	return routeName
}

func (msg MsgDeleteSysFeeSetting) Type() string {
	return "deleteFeeSetting"
}

func (msg MsgDeleteSysFeeSetting) ValidateBasic() sdkTypes.Error {
	if msg.Issuer.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Issuer.String())
	}

	return nil
}

func (msg MsgDeleteSysFeeSetting) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgDeleteSysFeeSetting) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Issuer}
}

// MsgAssignFeeToMsg create tx type fee
type MsgAssignFeeToMsg struct {
	FeeName string              `json:"fee_name"`
	MsgType string              `json:"msg_type"`
	Issuer  sdkTypes.AccAddress `json:"issuer"`
}

func NewMsgAssignFeeToMsg(name string, msgType string, issuer sdkTypes.AccAddress) MsgAssignFeeToMsg {
	return MsgAssignFeeToMsg{
		FeeName: name,
		MsgType: msgType,
		Issuer:  issuer,
	}
}

func (msg MsgAssignFeeToMsg) Route() string {
	return routeName
}

func (msg MsgAssignFeeToMsg) Type() string {
	return "assignFeeToMsg"
}

func (msg MsgAssignFeeToMsg) ValidateBasic() sdkTypes.Error {
	if msg.Issuer.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Issuer.String())
	}

	if len(msg.FeeName) <= 0 {
		return sdkTypes.ErrInvalidCoins("Fee name type cant be empty.")
	}

	if len(msg.MsgType) <= 0 {
		return sdkTypes.ErrInvalidCoins("Message Type cant be empty.")
	}

	return nil
}

func (msg MsgAssignFeeToMsg) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgAssignFeeToMsg) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Issuer}
}

// MsgAssignFeeToAcc create account fee type
type MsgAssignFeeToAcc struct {
	FeeName string              `json:"fee_name"`
	Account sdkTypes.AccAddress `json:"account"`
	Issuer  sdkTypes.AccAddress `json:"issuer"`
}

func NewMsgAssignFeeToAcc(name string, account, issuer sdkTypes.AccAddress) MsgAssignFeeToAcc {
	return MsgAssignFeeToAcc{
		FeeName: name,
		Account: account,
		Issuer:  issuer,
	}
}

func (msg MsgAssignFeeToAcc) Route() string {
	return routeName
}

func (msg MsgAssignFeeToAcc) Type() string {
	return "assignFeeToAcc"
}

func (msg MsgAssignFeeToAcc) ValidateBasic() sdkTypes.Error {
	if msg.Issuer.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Issuer.String())
	}

	if len(msg.FeeName) <= 0 {
		return sdkTypes.ErrInvalidCoins("Fee name type cant be empty.")
	}

	if msg.Account.Empty() {
		return sdkTypes.ErrInvalidCoins("Account cant be empty.")
	}

	return nil
}

func (msg MsgAssignFeeToAcc) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgAssignFeeToAcc) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Issuer}
}

// MsgMultiplier create/update fee multiplier
type MsgMultiplier struct {
	Multiplier string              `json:"multiplier"`
	Issuer     sdkTypes.AccAddress `json:"issuer"`
}

func NewMsgMultiplier(multiplier string, issuer sdkTypes.AccAddress) MsgMultiplier {
	return MsgMultiplier{
		Multiplier: multiplier,
		Issuer:     issuer,
	}
}

func (msg MsgMultiplier) Route() string {
	return routeName
}

func (msg MsgMultiplier) Type() string {
	return "updateMultiplier"
}

func (msg MsgMultiplier) ValidateBasic() sdkTypes.Error {

	minMultiplier, _ := sdkTypes.NewDecFromStr("0")
	multiplier, err := sdkTypes.NewDecFromStr(msg.Multiplier)
	if err != nil {
		return err
	}

	if msg.Issuer.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Issuer.String())
	}

	if len(msg.Multiplier) <= 0 {
		return sdkTypes.ErrInvalidCoins("Multiplier cant be empty.")
	}

	if !multiplier.GT(minMultiplier) {
		return sdkTypes.ErrInternal("Multiplier invalid.")
	}

	return nil
}

func (msg MsgMultiplier) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgMultiplier) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Issuer}
}

// MsgTokenMultiplier create/update token fee multiplier
type MsgTokenMultiplier struct {
	Multiplier string              `json:"multiplier"`
	Issuer     sdkTypes.AccAddress `json:"issuer"`
}

func NewMsgTokenMultiplier(multiplier string, issuer sdkTypes.AccAddress) MsgTokenMultiplier {
	return MsgTokenMultiplier{
		Multiplier: multiplier,
		Issuer:     issuer,
	}
}

func (msg MsgTokenMultiplier) Route() string {
	return routeName
}

func (msg MsgTokenMultiplier) Type() string {
	return "updateTokenMultiplier"
}

func (msg MsgTokenMultiplier) ValidateBasic() sdkTypes.Error {

	minMultiplier, _ := sdkTypes.NewDecFromStr("0")
	multiplier, err := sdkTypes.NewDecFromStr(msg.Multiplier)
	if err != nil {
		return err
	}

	if msg.Issuer.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Issuer.String())
	}

	if len(msg.Multiplier) <= 0 {
		return sdkTypes.ErrInvalidCoins("Multiplier cant be empty.")
	}

	if !multiplier.GT(minMultiplier) {
		return sdkTypes.ErrInternal("Multiplier invalid.")
	}

	return nil
}

func (msg MsgTokenMultiplier) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgTokenMultiplier) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Issuer}
}

type MsgAssignFeeToToken struct {
	FeeName string              `json:"fee_name"`
	Symbol  string              `json:"symbol"`
	Action  string              `json:"action"`
	Issuer  sdkTypes.AccAddress `json:"issuer"`
}

func NewMsgAssignFeeToToken(name, symbol, action string, issuer sdkTypes.AccAddress) MsgAssignFeeToToken {
	return MsgAssignFeeToToken{
		FeeName: name,
		Symbol:  symbol,
		Action:  action,
		Issuer:  issuer,
	}
}

func (msg MsgAssignFeeToToken) Route() string {
	return routeName
}

func (msg MsgAssignFeeToToken) Type() string {
	return "assignFeeToToken"
}

func (msg MsgAssignFeeToToken) ValidateBasic() sdkTypes.Error {
	if msg.Issuer.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Issuer.String())
	}

	if len(msg.FeeName) <= 0 {
		return sdkTypes.ErrInvalidCoins("Fee name type cant be empty.")
	}

	if len(msg.Symbol) <= 0 {
		return sdkTypes.ErrInvalidCoins("Symbol cant be empty.")
	}

	if len(msg.Action) <= 0 {
		return sdkTypes.ErrInvalidCoins("Action cant be empty.")
	}

	return nil
}

func (msg MsgAssignFeeToToken) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgAssignFeeToToken) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Issuer}
}
