package kyc

import (
	"encoding/json"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type MsgWhitelist struct {
	Owner   sdkTypes.AccAddress `json:"owner"`
	KycData KycData             `json:"kycData"`
}

type KycData struct {
	Payload    Payload     `json:"payload"`
	Signatures []Signature `json:"signatures"`
}

type Payload struct {
	Kyc           Kyc `json:"kyc"`
	crypto.PubKey `json:"pub_key"`
	Signature     []byte `json:"signature"`
}

type Kyc struct {
	From       sdkTypes.AccAddress `json:"from"`
	Nonce      string              `json:"nonce"`
	KycAddress string              `json:"kycAddress"` /// It's a reference to the kyc data
}

type Signature struct {
	crypto.PubKey `json:"pub_key"`
	Signature     []byte `json:"signature"`
}

func NewMsgWhitelist(owner sdkTypes.AccAddress, kycData KycData) MsgWhitelist {
	return MsgWhitelist{
		Owner:   owner,
		KycData: kycData,
	}
}

func NewPayload(kyc Kyc, pubKey crypto.PubKey, signature []byte) Payload {
	return Payload{
		Kyc:       kyc,
		PubKey:    pubKey,
		Signature: signature,
	}
}

func NewKycData(payload Payload, signatures []Signature) KycData {
	return KycData{
		Payload:    payload,
		Signatures: signatures,
	}
}

func NewKyc(from sdkTypes.AccAddress, nonce string, kycAddress string) Kyc {
	return Kyc{
		From:       from,
		Nonce:      nonce,
		KycAddress: kycAddress,
	}
}

func NewSignature(pubKey crypto.PubKey, signature []byte) Signature {
	return Signature{
		PubKey:    pubKey,
		Signature: signature,
	}
}

func (msg MsgWhitelist) Route() string {
	return "kyc"
}

func (msg MsgWhitelist) Type() string {
	return "whitelist"
}

func (msg MsgWhitelist) ValidateBasic() sdkTypes.Error {

	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}

	if len(msg.KycData.Signatures) < 2 {
		return sdkTypes.ErrInvalidAddress("Insufficient issuer signature.")
	}

	if msg.KycData.Payload.Kyc.From.Empty() {
		return sdkTypes.ErrInvalidAddress("From cannot be empty.")
	}

	if len(msg.KycData.Payload.Kyc.KycAddress) <= 0 {
		return sdkTypes.ErrInvalidAddress(msg.KycData.Payload.Kyc.KycAddress)
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

func (kyc Kyc) GetFromSignBytes() []byte {
	b, err := json.Marshal(kyc)
	if err != nil {
		panic(err)
	}

	return sdkTypes.MustSortJSON(b)
}

func (msg MsgWhitelist) GetSignBytes() []byte {

	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners get signers
func (msg MsgWhitelist) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
}

type MsgRevokeWhitelist struct {
	Owner         sdkTypes.AccAddress `json:"owner"`
	RevokePayload RevokePayload       `json:"payload"`
	Signatures    []Signature         `json:"signatures"`
}

type RevokePayload struct {
	RevokeKycData RevokeKycData `json:"kyc"`
	crypto.PubKey `json:"pub_key"`
	Signature     []byte `json:"signature"`
}

type RevokeKycData struct {
	From  sdkTypes.AccAddress `json:"from"`
	Nonce string              `json:"nonce"`
	To    sdkTypes.AccAddress `json:"to"`
}

func NewMsgRevokeWhitelist(owner sdkTypes.AccAddress, revokePayload RevokePayload, signatures []Signature) MsgRevokeWhitelist {
	return MsgRevokeWhitelist{
		Owner:         owner,
		RevokePayload: revokePayload,
		Signatures:    signatures,
	}
}

func NewRevokeKycData(from sdkTypes.AccAddress, nonce string, to sdkTypes.AccAddress) RevokeKycData {
	return RevokeKycData{
		From:  from,
		Nonce: nonce,
		To:    to,
	}
}

func NewRevokePayload(revokeKycData RevokeKycData, pub_key crypto.PubKey, signature []byte) RevokePayload {
	return RevokePayload{
		RevokeKycData: revokeKycData,
		PubKey:        pub_key,
		Signature:     signature,
	}
}

func (msg MsgRevokeWhitelist) Route() string {
	return "kyc"
}

func (msg MsgRevokeWhitelist) Type() string {
	return "revokeWhitelist"
}

func (msg MsgRevokeWhitelist) ValidateBasic() sdkTypes.Error {

	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}

	if len(msg.Signatures) < 1 {
		return sdkTypes.ErrInvalidAddress("Insufficient issuer signature.")
	}

	if msg.RevokePayload.RevokeKycData.From.Empty() {
		return sdkTypes.ErrInvalidAddress("From cannot be empty.")
	}

	return nil
}

func (revokePayload RevokePayload) GetRevokeIssuerSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(revokePayload)
	if err != nil {
		panic(err)
	}

	return sdkTypes.MustSortJSON(b)
}

func (revokeKycData RevokeKycData) GetRevokeFromSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(revokeKycData)
	if err != nil {
		panic(err)
	}

	return sdkTypes.MustSortJSON(b)
}

func (msg MsgRevokeWhitelist) GetSignBytes() []byte {

	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners get signers
func (msg MsgRevokeWhitelist) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
}
