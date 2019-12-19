package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/maxonrow/maxonrow-go/types"
)

type MsgCreateAlias struct {
	Name  string         `json:"name"`
	Owner sdk.AccAddress `json:"owner"`
	Fee   Fee            `json:"fee"`
}

type Fee struct {
	To    sdkTypes.AccAddress `json:"to"`
	Value string              `json:"value"`
}

func NewMsgCreateAlias(name string, owner sdk.AccAddress, fee Fee) MsgCreateAlias {
	return MsgCreateAlias{
		Name:  name,
		Owner: owner,
		Fee:   fee,
	}
}

func (msg MsgCreateAlias) Route() string {
	return "nameservice"
}

func (msg MsgCreateAlias) Type() string {
	return "createAlias"
}

func (msg MsgCreateAlias) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 {
		return sdkTypes.ErrUnknownRequest("Alias cannot be empty.")
	}

	if msg.Name == msg.Owner.String() {
		return types.ErrAliasIsInUsed()
	}

	if len(msg.Fee.Value) == 0 {
		return sdkTypes.ErrUnknownRequest("Fee cannot be empty.")
	}

	return nil
}

func (msg MsgCreateAlias) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateAlias) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgSetAliasStatus
type MsgSetAliasStatus struct {
	Owner      sdkTypes.AccAddress `json:"owner"`
	Payload    Payload             `json:"payload"`
	Signatures []Signature         `json:"signatures"`
}

type Payload struct {
	Alias         AliasData `json:"alias"`
	crypto.PubKey `json:"pub_key"`
	Signature     []byte `json:"signature"`
}

type AliasData struct {
	From   sdkTypes.AccAddress `json:"from"`
	Nonce  string              `json:"nonce"`
	Status string              `json:"status"`
	Name   string              `json:"name"`
}

func NewPayload(alisdata AliasData, pubkey crypto.PubKey, Signature []byte) *Payload {
	return &Payload{
		Alias:     alisdata,
		PubKey:    pubkey,
		Signature: Signature,
	}

}

func NewAlias(from sdkTypes.AccAddress, nonce, status, name string) *AliasData {
	return &AliasData{
		From:   from,
		Nonce:  nonce,
		Status: status,
		Name:   name,
	}
}

func NewMsgSetAliasStatus(owner sdkTypes.AccAddress, payload Payload, signatures []Signature) *MsgSetAliasStatus {
	return &MsgSetAliasStatus{
		Owner:      owner,
		Payload:    payload,
		Signatures: signatures,
	}
}

func (msg MsgSetAliasStatus) Route() string {
	return "nameservice"
}

func (msg MsgSetAliasStatus) Type() string {
	return "setAliasStatus"
}

func (msg MsgSetAliasStatus) ValidateBasic() sdkTypes.Error {
	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}

	if len(msg.Payload.Alias.Name) < 1 {
		return sdkTypes.ErrUnknownRequest("Alias name cannot be empty.")
	}

	if len(msg.Signatures) < 1 {
		return sdkTypes.ErrInvalidAddress("Insufficient issuer signature.")
	}

	if msg.Payload.Alias.From.Empty() {
		return sdkTypes.ErrInvalidAddress("From cannot be empty.")
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

func (alias AliasData) GetFromSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(alias)
	if err != nil {
		panic(err)
	}

	return sdkTypes.MustSortJSON(b)
}

func (msg MsgSetAliasStatus) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgSetAliasStatus) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
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
