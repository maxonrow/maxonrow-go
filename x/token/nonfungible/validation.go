package nonfungible

import (
	"fmt"
	"strconv"
	"strings"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	// TODO TODO - calculate exactly
	MaxLength            = 256
	TokenNameMaxLength   = 100
	TokenSymbolMaxLength = 100
)

func validateTokenName(tokenName string) sdkTypes.Error {
	if len(tokenName) == 0 || len(tokenName) > TokenNameMaxLength {
		// TODO - return appropriate error - need to have something like ErrInvalidRequest
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
		// TODO - return appropriate error - need to have something like ErrInvalidRequest
		return sdkTypes.ErrUnknownRequest(
			fmt.Sprintf("Invalid token symbol field length: %d", len(symbol)))
	}

	if strings.ContainsAny(symbol, ";:") {
		return sdkTypes.ErrUnknownRequest("Token symbol cannot contain following characters: ;:")
	}

	return nil
}

func validateMetadata(metadata string) sdkTypes.Error {
	if len(metadata) > MaxLength {
		return sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid metadata field length: %d", len(metadata)))
	}

	return nil
}

func validateProperties(properties string) sdkTypes.Error {
	if len(properties) > MaxLength {
		return sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid properties field length: %d", len(properties)))
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
