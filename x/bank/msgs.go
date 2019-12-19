package bank

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkBank "github.com/cosmos/cosmos-sdk/x/bank"
)

// RouterKey is they name of the bank module
const RouterKey = "bank"

// MsgSend - high level transaction of the coin module
type MsgMxwSend struct {
	sdkBank.MsgSend
}

var _ sdk.Msg = MsgMxwSend{}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins) MsgMxwSend {
	return MsgMxwSend{
		MsgSend: sdkBank.MsgSend{
			FromAddress: fromAddr,
			ToAddress:   toAddr,
			Amount:      amount,
		},
	}
}

// MarshalJSON marshals to JSON using Bech32.
func (msg MsgMxwSend) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.MsgSend)
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (msg *MsgMxwSend) UnmarshalJSON(data []byte) error {
	var m sdkBank.MsgSend
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	*msg = MsgMxwSend{MsgSend: m}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMxwSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}
