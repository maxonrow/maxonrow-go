package fungible

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

var prefixAuthorised = []byte("token/authorised")
var prefixProvider = []byte("token/provider")
var prefixIssuer = []byte("token/issuer")

// keys
func getAuthorisedKey() []byte {
	return prefixAuthorised
}

func getProviderKey() []byte {
	return prefixProvider
}

func getIssuerKey() []byte {
	return prefixIssuer
}

func getTokenKey(symbol string) []byte {
	return []byte(fmt.Sprintf("symbol:%s", symbol))
}

func getFungibleAccountKey(symbol string, owner sdkTypes.AccAddress) []byte {
	key := make([]byte, 0, len(symbol)+1+len(owner))
	key = append(key, []byte(symbol)...)
	key = append(key, ':')
	key = append(key, owner...)
	return key
}
