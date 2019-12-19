package types

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

const (
	Bech32PrefixAccAddr  = "mxw"
	Bech32PrefixAccPub   = "mxwpub"
	Bech32PrefixValAddr  = "mxwvaloper"
	Bech32PrefixValPub   = "mxwvaloperpub"
	Bech32PrefixConsAddr = "mxwvalcons"
	Bech32PrefixConsPub  = "mxwvalconspub"

	SYSTEM = "system" // Event type
)

var ZeroCoins = []sdkTypes.Coin{
	sdkTypes.Coin{
		Denom:  CIN,
		Amount: sdkTypes.NewInt(0),
	},
}
