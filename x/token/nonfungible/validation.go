package nonfungible

import (
	"fmt"
	"strconv"
	"strings"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	MaxLength            = 256
	TokenNameMaxLength   = 100
	TokenSymbolMaxLength = 40
	TokenItemIDMaxLength = 128
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

func ValidateSymbol(symbol string) sdkTypes.Error {
	if len(symbol) == 0 || len(symbol) > TokenSymbolMaxLength {
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

func ValidateItemID(itemID string) sdkTypes.Error {
	if len(itemID) > TokenItemIDMaxLength {
		return sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid itemID field length: %d", len(itemID)))
	}
	return nil
}
