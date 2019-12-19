package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type Value struct {
	Hash   string   `json:"hash"`
	Params []string `json:"params"`
}

// MakeMxwEvents create mxw events for scanner.
func MakeMxwEvents(eventSignature string, eventOwner string, eventParams []string) sdkTypes.Events {

	key := eventOwner
	var value = new(Value)

	hasher := sha256.New()
	hasher.Write([]byte(eventSignature))
	md := hasher.Sum(nil)

	buff := bytes.NewBuffer(md)
	buff.Truncate(20)
	md = buff.Bytes()

	hexEncoded := hex.EncodeToString(md)

	value.Hash = hexEncoded
	value.Params = eventParams

	out, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	attribute := sdkTypes.NewAttribute(key, string(out))
	event := sdkTypes.NewEvent(SYSTEM, attribute)
	return sdkTypes.Events{event}

}
