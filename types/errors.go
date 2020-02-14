package types

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

const (
	// Token
	CodeTokenDuplicated                     sdkTypes.CodeType = 2001
	CodeTokenInvalidSymbol                  sdkTypes.CodeType = 2002
	CodeTokenApproved                       sdkTypes.CodeType = 2003
	CodeTokenFrozen                         sdkTypes.CodeType = 2004
	CodeTokenUnfrozen                       sdkTypes.CodeType = 2005
	CodeTokenInvalid                        sdkTypes.CodeType = 2006
	CodeTokenAccountFrozen                  sdkTypes.CodeType = 2007
	CodeTokenAccountUnfrozen                sdkTypes.CodeType = 2008
	CodeTokenInvalidMinter                  sdkTypes.CodeType = 2009
	CodeTokenInvalidAccount                 sdkTypes.CodeType = 2010
	CodeTokenInvalidSupply                  sdkTypes.CodeType = 2099
	CodeTokenInvalidAccountBalance          sdkTypes.CodeType = 2100
	CodeTokenInvalidAction                  sdkTypes.CodeType = 2101
	CodeTokenInvalidNewOwner                sdkTypes.CodeType = 2102
	CodeTokenInvalidOwner                   sdkTypes.CodeType = 2103
	CodeTokenTransferTokenOwnershipApproved sdkTypes.CodeType = 2104
	CodeTokenItemIDInUsed                   sdkTypes.CodeType = 2105
	CodeTokenInvalidEndorser                sdkTypes.CodeType = 2106
	CodeTokenItemFrozen                     sdkTypes.CodeType = 2107

	CodeFeeNotFound             sdkTypes.CodeType = 3001
	CodeTokenFeeSettingNotFound sdkTypes.CodeType = 3002

	// Alias
	CodeAliasInUsed                 sdkTypes.CodeType = 4001
	CodeAliasNoSuchPendingAlias     sdkTypes.CodeType = 4002
	CodeAliasNotAllowedToCreate     sdkTypes.CodeType = 4003
	CodeAliasNotFound               sdkTypes.CodeType = 4004
	CodeAliasCouldNotResolveAddress sdkTypes.CodeType = 4005

	CodespaceMXW sdkTypes.CodespaceType = "mxw"
)

func newErrorWithMXWCodespace(code sdkTypes.CodeType, format string, args ...interface{}) sdkTypes.Error {
	return sdkTypes.NewError(CodespaceMXW, code, format, args...)
}

/// --- Fee errors
func ErrFeeSettingNotExists(feeName string) sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeFeeNotFound, "Fee setting in not valid: %s", feeName)
}
func ErrTokenFeeSettingNotExists(symbol string) sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenFeeSettingNotFound, "Token fee setting not found, token symbol: %s", symbol)
}

/// --- Token errors
func ErrTokenExists(symbol string) sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenDuplicated, "Token already exists: %s", symbol)
}
func ErrInvalidTokenSymbol(symbol string) sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenInvalidSymbol, "Token does not exist: %s", symbol)
}
func ErrTokenAlreadyApproved(symbol string) sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenApproved, "Token already approved: %s", symbol)
}
func ErrTokenFrozen() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenFrozen, "Token is frozen.")
}
func ErrTokenUnFrozen() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenUnfrozen, "Token already unfrozen.")
}
func ErrTokenInvalid() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenInvalid, "Invalid token.")
}
func ErrTokenAccountFrozen() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenAccountFrozen, "Token account frozen.")
}
func ErrTokenAccountUnFrozen() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenAccountUnfrozen, "Token account already unfrozen.")
}
func ErrInvalidTokenMinter() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenInvalidMinter, "Invalid token minter.")
}
func ErrInvalidTokenAccount() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenInvalidAccount, "Invalid token account.")
}
func ErrInvalidTokenSupply() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenInvalidSupply, "Invalid token supply.")
}
func ErrInvalidTokenAccountBalance(msg string) sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenInvalidAccountBalance, msg)
}

func ErrInvalidTokenAction() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenInvalidAction, "Invalid token action.")
}

func ErrInvalidTokenNewOwner() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenInvalidNewOwner, "Invalid token new owner.")
}

func ErrInvalidTokenOwner() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenInvalidOwner, "Invalid token owner.")
}

func ErrTokenTransferTokenOwnershipApproved() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenTransferTokenOwnershipApproved, "Transfer token ownership has been approved.")
}

// Alias
func ErrAliasIsInUsed() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeAliasInUsed, "Alias in used.")
}

func ErrAliasNoSuchPendingAlias() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeAliasNoSuchPendingAlias, "No such pending alias.")
}

func ErrAliasNotAllowedToCreate() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeAliasNotAllowedToCreate, "Not allowed to create new alias, you have pending alias approval.")
}

func ErrAliasNotFound() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeAliasNotFound, "Alias not found.")
}

func ErrAliasCouldNotResolveAddress() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeAliasCouldNotResolveAddress, "Could not resolve address.")
}

func ErrTokenItemIDInUsed() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenItemIDInUsed, "Token item id is in used.")
}

func ErrInvalidEndorser() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenInvalidEndorser, "Item endorser invalid.")
}

func ErrTokenItemFronzen() sdkTypes.Error {
	return newErrorWithMXWCodespace(CodeTokenItemFrozen, "Token item frozen.")
}
