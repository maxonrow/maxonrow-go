package auth

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = "auth"

type MsgCreateMultiSigAccount struct {
	Owner     sdkTypes.AccAddress
	Threshold int
	Signers   []sdkTypes.AccAddress
}

var _ sdkTypes.Msg = MsgCreateMultiSigAccount{}

func NewMsgCreateMultiSigAccount(master sdkTypes.AccAddress, threshold int, signers []sdkTypes.AccAddress) MsgCreateMultiSigAccount {
	return MsgCreateMultiSigAccount{master, threshold, signers}
}

func (msg MsgCreateMultiSigAccount) Route() string {
	return RouterKey
}

func (msg MsgCreateMultiSigAccount) Type() string {
	return "createMultiSigAccount"
}

func (msg MsgCreateMultiSigAccount) ValidateBasic() sdkTypes.Error {
	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}

	if len(msg.Signers) < msg.Threshold {
		return sdkTypes.ErrInternal("Invalid signers/ threshold signers: " + string(len(msg.Signers)) + " <  threshold: " + string(msg.Threshold))
	}

	return nil
}

func (msg MsgCreateMultiSigAccount) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateMultiSigAccount) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
}

type MsgUpdateMultiSigAccount struct {
	Owner        sdkTypes.AccAddress   `json:owner`
	GroupAddress sdkTypes.AccAddress   `json:groupAddress`
	NewThreshold int                   `json:threshold`
	NewSigners   []sdkTypes.AccAddress `json:signers`
}

func NewMsgUpdateMultiSigAccount(owner, groupAddress sdkTypes.AccAddress, threshold int, signers []sdkTypes.AccAddress) MsgUpdateMultiSigAccount {
	return MsgUpdateMultiSigAccount{owner, groupAddress, threshold, signers}
}

func (msg MsgUpdateMultiSigAccount) Route() string {
	return RouterKey
}

func (msg MsgUpdateMultiSigAccount) Type() string {
	return "updateMultiSigAccount"
}

func (msg MsgUpdateMultiSigAccount) ValidateBasic() sdkTypes.Error {
	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}

	if msg.GroupAddress.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.GroupAddress.String())
	}

	if len(msg.NewSigners) < msg.NewThreshold {
		return sdkTypes.ErrInternal("Invalid signers/ threshold signers: " + string(len(msg.NewSigners)) + " <  threshold: " + string(msg.NewThreshold))
	}
	return nil
}

func (msg MsgUpdateMultiSigAccount) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgUpdateMultiSigAccount) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
}

type MsgTransferMultiSigOwner struct {
	GroupAddress sdkTypes.AccAddress `json:groupAddress`
	Owner        sdkTypes.AccAddress `json:owner`
	NewOwner     sdkTypes.AccAddress `json:newOwner`
}

func NewMsgTransferMultiSigOwner(groupAddress sdkTypes.AccAddress, newOwner, owner sdkTypes.AccAddress) MsgTransferMultiSigOwner {
	return MsgTransferMultiSigOwner{groupAddress, owner, newOwner}
}

func (msg MsgTransferMultiSigOwner) Route() string {
	return RouterKey
}

func (msg MsgTransferMultiSigOwner) Type() string {
	return "transferMultiSigOwner"
}

func (msg MsgTransferMultiSigOwner) ValidateBasic() sdkTypes.Error {
	if msg.Owner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Owner.String())
	}

	if msg.GroupAddress.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.GroupAddress.String())
	}

	if msg.NewOwner.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.NewOwner.String())
	}
	return nil
}

func (msg MsgTransferMultiSigOwner) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferMultiSigOwner) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Owner}
}

type MsgCreateMultiSigTx struct {
	GroupAddress sdkTypes.AccAddress `json:groupAddress`
	TxID         uint64              `json:txID`
	Tx           sdkTypes.Tx         `json:tx`
	Sender       sdkTypes.AccAddress `json:sender`
}

func NewMsgCreateMultiSigTx(groupAddress sdkTypes.AccAddress, txID uint64, tx sdkTypes.Tx, sender sdkTypes.AccAddress) MsgCreateMultiSigTx {
	return MsgCreateMultiSigTx{groupAddress, txID, tx, sender}
}

func (msg MsgCreateMultiSigTx) Route() string {
	return RouterKey
}

func (msg MsgCreateMultiSigTx) Type() string {
	return "createMutiSigTx"
}

func (msg MsgCreateMultiSigTx) ValidateBasic() sdkTypes.Error {
	if msg.Sender.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Sender.String())
	}

	if msg.GroupAddress.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.GroupAddress.String())
	}

	if msg.TxID < 0 {
		return sdkTypes.ErrInternal("TxID not allowed to be less than 0.")
	}
	return nil
}

func (msg MsgCreateMultiSigTx) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateMultiSigTx) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Sender}
}

type MsgSignMultiSigTx struct {
	GroupAddress sdkTypes.AccAddress `json:groupAddress`
	TxID         uint64              `json:txId`
	Sender       sdkTypes.AccAddress `json:sender`
}

func NewMsgSignMultiSigTx(groupAddress sdkTypes.AccAddress, txID uint64, sender sdkTypes.AccAddress) MsgSignMultiSigTx {
	return MsgSignMultiSigTx{groupAddress, txID, sender}
}

func (msg MsgSignMultiSigTx) Route() string {
	return RouterKey
}

func (msg MsgSignMultiSigTx) Type() string {
	return "signMutiSigTx"
}

func (msg MsgSignMultiSigTx) ValidateBasic() sdkTypes.Error {
	if msg.Sender.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Sender.String())
	}

	if msg.GroupAddress.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.GroupAddress.String())
	}

	if msg.TxID < 0 {
		return sdkTypes.ErrInternal("TxID not allowed to be less than 0.")
	}
	return nil
}

func (msg MsgSignMultiSigTx) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgSignMultiSigTx) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Sender}
}

type MsgDeleteMultiSigTx struct {
	GroupAddress sdkTypes.AccAddress `json:groupAddress`
	TxID         uint64              `json:txId`
	Sender       sdkTypes.AccAddress `json:sender`
}

func NewMsgDeleteMultiSigTx(groupAddress sdkTypes.AccAddress, txID uint64, sender sdkTypes.AccAddress) MsgDeleteMultiSigTx {
	return MsgDeleteMultiSigTx{groupAddress, txID, sender}
}

func (msg MsgDeleteMultiSigTx) Route() string {
	return RouterKey
}

func (msg MsgDeleteMultiSigTx) Type() string {
	return "deleteMutiSigTx"
}

func (msg MsgDeleteMultiSigTx) ValidateBasic() sdkTypes.Error {
	if msg.Sender.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.Sender.String())
	}

	if msg.GroupAddress.Empty() {
		return sdkTypes.ErrInvalidAddress(msg.GroupAddress.String())
	}

	if msg.TxID < 0 {
		return sdkTypes.ErrInternal("TxID not allowed to be less than 0.")
	}
	return nil
}

func (msg MsgDeleteMultiSigTx) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgDeleteMultiSigTx) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{msg.Sender}
}
