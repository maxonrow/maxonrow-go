package fungible

import (
	"fmt"
	"strconv"
	"strings"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	MetadataMaxLength    = 256
	TokenNameMaxLength   = 100
	TokenSymbolMaxLength = 40
)

func validateTokenName(tokenName string) sdkTypes.Error {
	if len(tokenName) == 0 || len(tokenName) > TokenNameMaxLength {
		return sdkTypes.ErrUnknownRequest(
			fmt.Sprintf("Invalid token name field length: %d", len(tokenName)))
	}

	if strings.ContainsAny(tokenName, ";:") {
		return sdkTypes.ErrUnknownRequest("Token name cannot contain following characters: ;:")
	}

	return nil
}

func validateSymbol(symbol string) sdkTypes.Error {
	if len(symbol) == 0 || len(symbol) > TokenSymbolMaxLength {
		return sdkTypes.ErrUnknownRequest(
			fmt.Sprintf("Invalid token symbol field length: %d", len(symbol)))
	}

	if strings.ContainsAny(symbol, ";:") {
		return sdkTypes.ErrUnknownRequest("Token symbol cannot contain following characters: ;:")
	}

	return nil
}

func validateMetadata(link string) sdkTypes.Error {
	if len(link) > MetadataMaxLength {
		return sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid metadata field length: %d", len(link)))
	}

	return nil
}

func validateAmount(amount string) sdkTypes.Error {
	_, err := strconv.Atoi(amount)
	if err != nil {
		return sdkTypes.ErrInvalidCoins(fmt.Sprintf("Invalid amount string: %s", err))
	}

	return nil
}

func validateDecimal(decimal int) sdkTypes.Error {
	if decimal < 0 || decimal > 18 {
		return sdkTypes.ErrInternal("Invalid decimal")
	}
	return nil
}
