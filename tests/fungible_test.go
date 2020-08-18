package tests

import (
	"testing"

	token "github.com/maxonrow/maxonrow-go/x/token/fungible"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

type TokenInfo struct {
	Action                       string
	ApplicationFee               string
	FeeCollector                 string
	Name                         string
	Symbol                       string
	Decimals                     int
	Owner                        string
	NewOwner                     string
	Metadata                     string
	FixedSupply                  bool
	MaxSupply                    string
	Approved                     bool
	Frozen                       bool
	Burnable                     bool
	AmountOfToken                string // to be : transfer/mint/burn
	Provider                     string
	ProviderNonce                string
	Issuer                       string
	FeeSettingName               string
	VerifyTransferTokenOwnership string
}

func makeFungibleTokenTxs() []*testCase {

	tcs := []*testCase{

		//-D
		///2. 20200818-MsgMintFungibleToken
		{"token", false, false, "Create token TFT-772 - Dynamic-supply with mint unlimited", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-772", "TFT-772", 8, "acc-40", "", "", false, "0", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token TFT-772 - for Dynamic-supply with mint unlimited purpose.", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-772", "TFT-772", 0, "", "", "", false, "0", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", true, true, "Mint token TFT-772 - Error, due to Invalid token symbol field length which maximum is 40.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-772-name", "TFT-772-aaaaabbbbbcccccdddddeeeeefffffggggghhhhhiiiiijjjjjkkkkklllllmmmmmnnnnnooooopppppqqqqqrrrrrssssstttttuuuuuvvvvvwwwwwxxxxxyyyyyzzzzzz11111222223333344444555556666677777888889999900000", 8, "acc-40", "carlo", "", false, "0", false, false, true, "0", "", "", "", "", ""}, "", nil},

		//-D
		///3. 20200818-MsgTransferFungibleToken
		{"token", false, false, "Create token TFT-773 - Dynamic-supply with mint unlimited", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-773", "TFT-773", 8, "acc-40", "", "", false, "0", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token TFT-773 - for Dynamic-supply with mint unlimited purpose.", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-773", "TFT-773", 0, "", "", "", false, "0", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", false, false, "Mint token TFT-773 - Dynamic-supply which allowed to mint-token with 0 amount.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-773-name", "TFT-773", 8, "acc-40", "carlo", "", false, "0", false, false, true, "0", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Transfer token TFT-773 -Error due to Invalid token symbol field length which maximum is 40.", "acc-40", "100000000cin", 0, TokenInfo{"transfer", "", "", "", "TFT-773-aaaaabbbbbcccccdddddeeeeefffffggggghhhhhiiiiijjjjjkkkkklllllmmmmmnnnnnooooopppppqqqqqrrrrrssssstttttuuuuuvvvvvwwwwwxxxxxyyyyyzzzzzz11111222223333344444555556666677777888889999900000", 8, "acc-40", "eve", "", false, "", true, false, false, "0", "", "", "", "", ""}, "", nil},

		//-D
		///4. 20200818-MsgBurnFungibleToken
		{"token", false, false, "Create token TFT-774 - Dynamic-supply with mint unlimited", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-774", "TFT-774", 8, "acc-40", "", "", false, "0", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token TFT-774 - for Dynamic-supply with mint unlimited purpose.", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-774", "TFT-774", 0, "", "", "", false, "0", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", false, false, "Mint token TFT-774 - Dynamic-supply which allowed to mint-token with 0 amount.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-774-name", "TFT-774", 8, "acc-40", "carlo", "", false, "0", false, false, true, "0", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Burn token TFT-774 - Error due to Invalid token symbol field length which maximum is 40.", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TFT-774-aaaaabbbbbcccccdddddeeeeefffffggggghhhhhiiiiijjjjjkkkkklllllmmmmmnnnnnooooopppppqqqqqrrrrrssssstttttuuuuuvvvvvwwwwwxxxxxyyyyyzzzzzz11111222223333344444555556666677777888889999900000", 0, "acc-40", "", "", false, "", false, false, true, "100", "", "", "", "", ""}, "", nil},

		//-D
		///1. 20200818-MsgTransferFungibleTokenOwnership
		{"token", false, false, "Create token TFT-776 - Dynamic-supply with mint unlimited", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-776", "TFT-776", 8, "acc-40", "", "", false, "0", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token TFT-776 - for Dynamic-supply with mint unlimited purpose.", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-776", "TFT-776", 0, "", "", "", false, "0", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", true, true, "Transfer token ownership TFT-776 - Error due to Invalid token symbol field length which maximum is 40.", "acc-40", "100000000cin", 0, TokenInfo{"transfer-ownership", "", "", "", "TFT-776-aaaaabbbbbcccccdddddeeeeefffffggggghhhhhiiiiijjjjjkkkkklllllmmmmmnnnnnooooopppppqqqqqrrrrrssssstttttuuuuuvvvvvwwwwwxxxxxyyyyyzzzzzz11111222223333344444555556666677777888889999900000", 8, "acc-40", "carlo", "", false, "", true, false, false, "1", "", "", "", "", ""}, "", nil},

		//-D
		///1. 20200818-MsgAcceptFungibleTokenOwnership
		{"token", false, false, "Create token TFT-775 - Dynamic-supply with mint unlimited", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-775", "TFT-775", 8, "acc-40", "", "", false, "0", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token TFT-775 - for Dynamic-supply with mint unlimited purpose.", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-775", "TFT-775", 0, "", "", "", false, "0", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", false, false, "Transfer token ownership TFT-775 - Happy path", "acc-40", "100000000cin", 0, TokenInfo{"transfer-ownership", "", "", "", "TFT-775", 8, "acc-40", "carlo", "", false, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token transfer ownership TFT-775 - Happy path", "token-auth-1", "0cin", 0, TokenInfo{"verify-transfer-tokenOwnership", "", "", "", "TFT-775", 0, "", "", "", true, "", true, false, true, "", "token-prov-1", "0", "token-issuer-1", "", "APPROVE_TRANFER_TOKEN_OWNERSHIP"}, "", nil},
		{"token", true, true, "Accept token ownership TFT-775 - Error due to Invalid token symbol field length which maximum is 40.", "carlo", "100000000cin", 0, TokenInfo{"accept-ownership", "", "", "", "TFT-775--aaaaabbbbbcccccdddddeeeeefffffggggghhhhhiiiiijjjjjkkkkklllllmmmmmnnnnnooooopppppqqqqqrrrrrssssstttttuuuuuvvvvvwwwwwxxxxxyyyyyzzzzzz11111222223333344444555556666677777888889999900000", 8, "acc-40", "carlo", "", false, "", true, false, false, "1", "", "", "", "", ""}, "", nil},

		// Create Fungible Token
		{"token", false, false, "Create token - Happy Path", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken", "TT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Empty Symbol", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken", "", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Empty Name", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "", "TT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Negative application fee", "acc-40", "100000000cin", 0, TokenInfo{"create", "-0.000000000008", "mostafa", "TestToken", "TTTTTT1", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Negative application fee", "acc-40", "100000000cin", 0, TokenInfo{"create", "-123.4567", "mostafa", "TestToken", "TTTTTT1", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Creat e token - Existed Token Name", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken", "TTV2", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - 19 decimal places", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken", "TTV3", 19, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Create token - 0 decimal places", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken", "TTV3", 0, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Not fee collector", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "acc-20", "TestToken", "TT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Creating symbol that has been created", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken", "TT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Create token - Generous person pay for system fees", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken", "ToT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Poor people trying to create enormous application fee token", "acc-40", "100000000cin", 0, TokenInfo{"create", "1000000000000000000000000000", "mostafa", "TestToken", "ToT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		//{"token", true, true, "Create token - Some stranger(not whitelisted) trying to create token", "josephin", "0cin", 0, TokenInfo{"create", "1000000", "mostafa", "TestToken", "TooT", 8, "josephin","", "", "", true, "100000", false, false, false, "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Some smart people steal people signature, but owner is not signer", "acc-40", "100000000cin", 0, TokenInfo{"create", "1000000000000000000000000000", "mostafa", "TestToken", "ToT", 8, "josephin", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Stingy people trying to escape application fee", "acc-40", "100000000cin", 0, TokenInfo{"create", "0", "mostafa", "TestToken", "ToT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Wrong gas", "acc-40", "100000000cin", 1, TokenInfo{"create", "10000000", "mostafa", "TestToken", "TTT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Some smart people trying negative amount hoping negative*negative=positive", "acc-40", "100000000cin", 0, TokenInfo{"create", "-10000000", "mostafa", "TestToken", "TTT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Some naggy people have the world longest metadata", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken", "TTT", 8, "acc-40", "", "abcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcbacbabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcbacbabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcbacbabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcbacbabcabc", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Some naggy people have the world longest symbol", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken", "abcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcbacbabcabcabcabcabcabcab", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Some naggy people have the world longest name", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "abcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcbacbabcabcabcabcabcabcab", "TTOTT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		//{"token", true, true, "Create token - Some forgetful people, forgot about token owner", "acc-40", "0cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken", "TTOTT", 8, "", "","", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Create token - Happy Path which is dynamic supply", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-1", "TT-1", 8, "acc-40", "", "", false, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Already existed", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-1", "TT-1", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "Token already existed", nil},
		{"token", false, false, "Create token - Happy Path which for freeze purpose", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-4", "TT-4", 8, "acc-40", "", "", false, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Create token - Happy Path which for unfreeze purpose", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-5", "TT-5", 8, "acc-40", "", "", false, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},

		// Approve Fungible Token
		// Only can approve that token created at Token - Create Fungible Token
		// If need more token please create them at Token - Create Fungible Token
		{"token", false, false, "Approve token - Happy path", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken", "TT", 8, "", "", "", true, "", false, false, true, "1", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", true, true, "Approve token - again", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken", "TT", 0, "", "", "", true, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", true, true, "Approve token - Token not existed", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-3", "TT-3", 0, "", "", "", true, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "zero", ""}, "", nil},
		{"token", true, true, "Approve token - Invalid signer", "mostafa", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-1", "TT-1", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "zero", ""}, "", nil},
		{"token", true, true, "Approve token - Invalid provider", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-1", "TT-1", 0, "", "", "", false, "", true, false, true, "1", "mostafa", "0", "token-issuer-1", "zero", ""}, "", nil},
		{"token", true, true, "Approve token - Invalid issuer", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-1", "TT-1", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "mostafa", "zero", ""}, "", nil},
		{"token", true, true, "Approve token - Invalid fee", "token-auth-1", "10000cin", 0, TokenInfo{"approve", "", "", "TestToken-1", "TT-1", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", false, false, "Approve token - Happy path for TT-1", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-1", "TT-1", 0, "", "", "", false, "", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", true, true, "Approve token - Not allow for TT-1 again", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-1", "TT-1", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", false, false, "Approve token - Happy path for TT-4", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-4", "TT-4", 0, "", "", "", false, "", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", false, false, "Approve token - Happy path for TT-5", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-5", "TT-5", 0, "", "", "", false, "", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},

		// Freeze Fungible Token - only if after the approved
		{"token", true, true, "Freeze token - Token not existed", "token-auth-1", "0cin", 0, TokenInfo{"freeze", "", "", "", "TT-3", 0, "", "", "", true, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, true, "Freeze token - Invalid signer", "mostafa", "0cin", 0, TokenInfo{"freeze", "", "", "", "TT-4", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, true, "Freeze token - Invalid provider", "token-auth-1", "0cin", 0, TokenInfo{"freeze", "", "", "", "TT-4", 0, "", "", "", false, "", true, false, true, "1", "mostafa", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, true, "Freeze token - Invalid issuer", "token-auth-1", "0cin", 0, TokenInfo{"freeze", "", "", "", "TT-4", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "mostafa", "", ""}, "", nil},
		{"token", true, true, "Freeze token - Invalid fee", "token-auth-1", "10000cin", 0, TokenInfo{"freeze", "", "", "", "TT-4", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, true, "Freeze token - Not Allow to freeze without the approval", "token-auth-1", "0cin", 0, TokenInfo{"freeze", "", "", "", "ToT", 0, "acc-40", "eve", "", false, "", false, false, false, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil}, //KIV
		{"token", false, false, "Freeze token - Happy path for TT-4-commit", "token-auth-1", "0cin", 0, TokenInfo{"freeze", "", "", "", "TT-4", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, true, "Freeze token - Not allow for TT-4 again", "token-auth-1", "0cin", 0, TokenInfo{"freeze", "", "", "", "TT-4", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},

		//[Once Frozen, Token Not allowed to : mint, burn, transfer-ownership, verify-ownership, accept-ownership] : OK
		{"token", false, false, "Create token (FT-NEW01) - Happy Path", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "FT-NEW01", "FT-NEW01", 8, "acc-40", "", "", false, "10000", false, false, true, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token (FT-NEW01) - for Dynamic-supply purpose.", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "FT-NEW01", "FT-NEW01", 0, "", "", "", false, "10000", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", false, false, "Freeze token (FT-NEW01) - Happy path. commit.", "token-auth-1", "0cin", 0, TokenInfo{"freeze", "", "", "", "FT-NEW01", 0, "", "", "", false, "", false, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, true, "Mint token (FT-NEW01) - Error due to Token is frozen.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "FT-NEW01", "FT-NEW01", 8, "acc-40", "carlo", "", false, "10000", false, false, true, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Burn token (FT-NEW01) - Error due to Token is frozen.", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "FT-NEW01", 0, "acc-40", "", "", false, "", false, false, true, "200", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Transfer token ownership (FT-NEW01) - Error due to Token is frozen.", "acc-40", "100000000cin", 0, TokenInfo{"transfer-ownership", "", "", "", "FT-NEW01", 8, "acc-40", "carlo", "", false, "", false, false, true, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Approve transfer token ownership (FT-NEW01) - Error due to Invalid token action because of Token is frozen.", "token-auth-2", "0cin", 0, TokenInfo{"verify-transfer-tokenOwnership", "", "", "", "FT-NEW01", 0, "", "", "", false, "", false, false, true, "", "token-prov-2", "0", "token-issuer-2", "", "APPROVE_TRANFER_TOKEN_OWNERSHIP"}, "", nil},
		{"token", true, true, "Accept token ownership (FT-NEW01) - Error due to Invalid token action because of Token is frozen.", "carlo", "100000000cin", 0, TokenInfo{"accept-ownership", "", "", "", "FT-NEW01", 8, "acc-40", "carlo", "", false, "", false, false, true, "1", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Unfreeze token (FT-NEW01) - Happy path", "token-auth-1", "0cin", 0, TokenInfo{"unfreeze", "", "", "", "FT-NEW01", 0, "", "", "", false, "", false, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},

		{"token", true, true, "Mint token - Not allow if was approved and frozen", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "", "TT-4", 8, "acc-40", "carlo", "", false, "100000", true, true, false, "0", "", "", "", "", ""}, "", nil},

		// Unfreeze Fungible Token - only if after the approved and already frozen
		{"token", true, true, "Unfreeze token - Token not existed", "token-auth-1", "0cin", 0, TokenInfo{"unfreeze", "", "", "", "TT-3", 0, "", "", "", true, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, true, "Unfreeze token - Invalid signer", "mostafa", "0cin", 0, TokenInfo{"unfreeze", "", "", "", "TT-4", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, true, "Unfreeze token - Invalid provider", "token-auth-1", "0cin", 0, TokenInfo{"unfreeze", "", "", "", "TT-4", 0, "", "", "", false, "", true, false, true, "1", "mostafa", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, true, "Unfreeze token - Invalid issuer", "token-auth-1", "0cin", 0, TokenInfo{"unfreeze", "", "", "", "TT-4", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "mostafa", "", ""}, "", nil},
		{"token", true, true, "Unfreeze token - Invalid fee", "token-auth-1", "10000cin", 0, TokenInfo{"unfreeze", "", "", "", "TT-4", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, true, "Unfreeze token - Approved but not yet freeze", "token-auth-1", "0cin", 0, TokenInfo{"unfreeze", "", "", "", "TT-5", 0, "acc-40", "eve", "", false, "", true, false, false, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil}, //KIV
		{"token", false, false, "Unfreeze token - Happy path for TT-4", "token-auth-1", "0cin", 0, TokenInfo{"unfreeze", "", "", "", "TT-4", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, true, "Unfreeze token - Not allow for TT-4 again", "token-auth-1", "0cin", 0, TokenInfo{"unfreeze", "", "", "", "TT-4", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},

		// mint
		{"token", true, true, "Mint token - Token not existed", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "", "TT-3", 8, "acc-40", "carlo", "", true, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Mint token - Invalid owner", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "", "TT-1", 8, "nago", "carlo", "", false, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Mint token - Happy path for TT-1 which is dynamic supply", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "", "TT-1", 8, "acc-40", "carlo", "", false, "100000", true, false, false, "100", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Mint token - Allow to issue another TWO amount of TT-1 again", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "", "TT-1", 8, "acc-40", "carlo", "", false, "100000", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Mint token - Allow to issue ZERO amount of TT-1", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "", "TT-1", 8, "acc-40", "carlo", "", false, "100000", true, false, false, "0", "", "", "", "", ""}, "", nil},

		// case-1.1 : Dynamic-supply where MaxSupply="1000" and FixedSupply = FALSE and allowed to mint-token with any amount [AmountOfToken] u like, but not allowed > MaxSupply
		{"token", false, false, "Create token TFT-1101 - Happy Path", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-1101", "TFT-1101", 8, "acc-40", "", "", false, "10000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token TFT-1101 - for Dynamic-supply purpose.", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-1101", "TFT-1101", 0, "", "", "", false, "10000", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", false, false, "Mint token TFT-1101 - Dynamic-supply which allowed to mint-token with any amount.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-1101", "TFT-1101", 8, "acc-40", "carlo", "", false, "10000", false, false, true, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Mint token TFT-1101 - Error of Invalid token supply, due to current TotalSupply > MaxSupply which is only 1000.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-1101", "TFT-1101", 8, "acc-40", "carlo", "", false, "10000", false, false, true, "99999", "", "", "", "", ""}, "", nil},

		// case-1.2 : Dynamic-supply where MaxSupply="1" and FixedSupply = FALSE and allowed to mint-token with any amount [AmountOfToken] u like, but not allowed > MaxSupply
		{"token", false, false, "Create token TFT-1201 - Happy Path", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-1201", "TFT-1201", 8, "acc-40", "", "", false, "1", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token TFT-1201 - for Dynamic-supply purpose.", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-1201", "TFT-1201", 0, "", "", "", false, "1", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", true, true, "Mint token TFT-1201 - Error of Invalid token supply, due to current TotalSupply > MaxSupply which is only 1.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-1201", "TFT-1201", 8, "acc-40", "carlo", "", false, "1", false, false, true, "2", "", "", "", "", ""}, "", nil},

		// case-1.3 : Dynamic-supply where MaxSupply="-1" and FixedSupply = FALSE
		//{"token", true, true, "Create token TFT-1301 - Dynamic-supply with Error panic while MaxSupply '-1'", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-1301", "TFT-1301", 8, "acc-40", "", "", false, "-1", false, false, false, "", "", "", "", "", ""}, "", nil},  // Error of - panic

		// case-1.4 : Dynamic-supply where MaxSupply="0" and FixedSupply = FALSE
		{"token", false, false, "Create token TFT-1401 - Dynamic-supply with mint unlimited", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-1401", "TFT-1401", 8, "acc-40", "", "", false, "0", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token TFT-1401 - for Dynamic-supply with mint unlimited purpose.", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-1401", "TFT-1401", 0, "", "", "", false, "0", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", false, false, "Mint token TFT-1401 - Dynamic-supply which allowed to mint-token with 0 amount.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-1401", "TFT-1401", 8, "acc-40", "carlo", "", false, "0", false, false, true, "0", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Mint token TFT-1401 - Dynamic-supply which allowed to mint-token with 2 amount.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-1401", "TFT-1401", 8, "acc-40", "carlo", "", false, "0", false, false, true, "2", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Mint token TFT-1401 - Dynamic-supply which allowed to mint-token with 100 amount.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-1401", "TFT-1401", 8, "acc-40", "carlo", "", false, "0", false, false, true, "100", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Mint token TFT-1401 - Dynamic-supply which allowed to mint-token with 300 amount.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-1401", "TFT-1401", 8, "acc-40", "carlo", "", false, "0", false, false, true, "300", "", "", "", "", ""}, "", nil},
		//{"token", true, true, "Mint token TFT-1401 - Dynamic-supply that allowed to mint-token with AmountOfToken=ZERO-with-78.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-1401", "TFT-1401", 8, "acc-40", "carlo", "", false, "0", false, false, true, "1000000000000000000000000000000000000000000000000000000000000000000000000000000", "", "", "", "", ""}, "", nil}, //AmountOfToken=ZERO-with-78			// Error of - 'panic: bit length 260 greater than 256'
		//{"token", true, true, "Mint token TFT-1401 - Dynamic-supply that allowed to mint-token with AmountOfToken=ZERO-with-77.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-1401", "TFT-1401", 8, "acc-40", "carlo", "", false, "0", false, false, true, "100000000000000000000000000000000000000000000000000000000000000000000000000000", "", "", "", "", ""}, "", nil},  //AmountOfToken=ZERO-with-77			// Error of - 'panic: bit length 257 greater than 256'
		{"token", true, true, "Mint token TFT-1401 - Dynamic-supply that allowed to mint-token with AmountOfToken=ZERO-with-76.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-1401", "TFT-1401", 8, "acc-40", "carlo", "", false, "0", false, false, true, "10000000000000000000000000000000000000000000000000000000000000000000000000000", "", "", "", "", ""}, "", nil}, //AmountOfToken=ZERO-with-76 			// Insufficient fee amount

		// case-2 : FixedSupply where MaxSupply!="0" & FixedSupply = TRUE, burnable=TRUE, is not allowed to do mint-token
		{"token", false, false, "Create token TFT-1112 MaxSupply != 0 - Happy Path", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-1112", "TFT-1112", 8, "acc-40", "", "", true, "100", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token TFT-1112 MaxSupply != 0 - for Fixed-supply purpose.", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-1112", "TFT-1112", 0, "", "", "", true, "2", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", true, true, "Mint token TFT-1112 MaxSupply != 0 - Error, due to Invalid token action where FixedSupply not allowed do mint-token.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-1112", "TFT-1112", 8, "acc-40", "carlo", "", true, "2", false, false, true, "222", "", "", "", "", ""}, "", nil},

		// case-2.2 : FixedSupply where MaxSupply!="0" & FixedSupply = TRUE, burnable=FALSE, is not allowed to do mint-token
		{"token", false, false, "Create token TFT-1441 MaxSupply != 0 - Happy Path", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-1441", "TFT-1441", 8, "acc-40", "", "", true, "100", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token TFT-1441 MaxSupply != 0 - for Fixed-supply purpose.", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-1441", "TFT-1441", 0, "", "", "", true, "2", false, false, false, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", true, true, "Mint token TFT-1441 MaxSupply != 0 - Error, due to Invalid token action where FixedSupply not allowed do mint-token.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-1441", "TFT-1441", 8, "acc-40", "carlo", "", true, "2", false, false, false, "222", "", "", "", "", ""}, "", nil},

		// case-3 : [fixedSupply cannot have MaxSupply = 0] and FixedSupply = TRUE, is not allowed to do mint-token
		{"token", false, false, "Create token TFT-1103 MaxSupply = 0 - Happy Path", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-1103", "TFT-1103", 8, "acc-40", "", "", true, "100", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token TFT-1103 MaxSupply = 0 - for Fixed-supply purpose.", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-1103", "TFT-1103", 0, "", "", "", true, "0", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", true, true, "Mint token TFT-1103 MaxSupply = 0 - Error, due to Invalid token action where FixedSupply not allowed do mint-token.", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-1103", "TFT-1103", 8, "acc-40", "carlo", "", true, "2", false, false, true, "222", "", "", "", "", ""}, "", nil},

		// transfer
		{"token", true, true, "Transfer token - Not allow to transfer if was approved and frozen", "acc-40", "100000000cin", 0, TokenInfo{"transfer", "", "", "", "TT-1", 8, "acc-40", "eve", "", false, "", true, true, false, "1000000000000000", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Transfer token - Token not existed", "acc-40", "100000000cin", 0, TokenInfo{"transfer", "", "", "", "TT-3", 8, "acc-40", "eve", "", false, "", false, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Transfer token - Invalid owner", "acc-40", "100000000cin", 0, TokenInfo{"transfer", "", "", "", "TT-1", 8, "nago", "eve", "", false, "", false, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Transfer token - Not approved", "acc-40", "100000000cin", 0, TokenInfo{"transfer", "", "", "", "ToT", 8, "acc-40", "eve", "", false, "", false, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Transfer token - Happy path", "acc-40", "100000000cin", 0, TokenInfo{"transfer", "", "", "", "TT-1", 8, "acc-40", "eve", "", false, "", true, false, false, "0", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Transfer token - Invalid fee", "acc-40", "1cin", 0, TokenInfo{"transfer", "", "", "", "TT-1", 8, "acc-40", "eve", "", false, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Transfer token - The remaining balance which allow to transfer same TT-1 again", "acc-40", "100000000cin", 0, TokenInfo{"transfer", "", "", "", "TT-1", 8, "acc-40", "eve", "", false, "", true, false, false, "0", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Transfer token - Not enough balance to transfer same TT-1 again", "acc-40", "100000000cin", 0, TokenInfo{"transfer", "", "", "", "TT-1", 8, "acc-40", "eve", "", false, "", true, false, false, "1000000000000000", "", "", "", "", ""}, "", nil},

		// transfer ownership
		{"token", true, true, "Transfer token ownership - Token not existed", "acc-40", "100000000cin", 0, TokenInfo{"transfer-ownership", "", "", "", "TT-3", 8, "acc-40", "carlo", "", true, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Transfer token ownership - Invalid owner", "acc-40", "100000000cin", 0, TokenInfo{"transfer-ownership", "", "", "", "TT-1", 8, "nago", "carlo", "", false, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Transfer token ownership - Happy path", "acc-40", "100000000cin", 0, TokenInfo{"transfer-ownership", "", "", "", "TT-1", 8, "acc-40", "carlo", "", false, "", true, false, false, "1", "", "", "", "", ""}, "", nil},

		// verify transfer token ownership
		{"token", false, false, "Approve token transfer ownership - Happy path for TT-1", "token-auth-1", "0cin", 0, TokenInfo{"verify-transfer-tokenOwnership", "", "", "", "TT-1", 0, "", "", "", true, "", true, false, true, "", "token-prov-1", "0", "token-issuer-1", "", "APPROVE_TRANFER_TOKEN_OWNERSHIP"}, "", nil},

		// accept ownership
		{"token", true, true, "Accept token ownership - Token not existed", "carlo", "100000000cin", 0, TokenInfo{"accept-ownership", "", "", "", "TT-3", 8, "acc-40", "carlo", "", true, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Accept token ownership - Invalid new owner", "carlo", "100000000cin", 0, TokenInfo{"accept-ownership", "", "", "", "TT-1", 8, "acc-40", "nago", "", false, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Accept token ownership - Happy path. commit", "carlo", "100000000cin", 0, TokenInfo{"accept-ownership", "", "", "", "TT-1", 8, "acc-40", "carlo", "", false, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Accept token ownership - Not Allow for TT-1 again", "carlo", "100000000cin", 0, TokenInfo{"accept-ownership", "", "", "", "TT-1", 8, "acc-40", "carlo", "", false, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Transfer token ownership - Allow for TT-1 transfer ownership again as the ownership already accepted by new party", "acc-40", "100000000cin", 0, TokenInfo{"transfer-ownership", "", "", "", "TT-1", 8, "acc-40", "carlo", "", false, "", true, false, false, "1", "", "", "", "", ""}, "", nil},

		// Create-Approval-Burn : using Fixed Supply
		{"token", true, true, "Create token - Use fixed supply with long memo", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-2", "TT-2", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "long memo:123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890", nil},
		{"token", false, false, "Create token - Use fixed supply with memo", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-2", "TT-2", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "1234567890", nil},
		{"token", false, false, "Approve token - Happy path for TT-2", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "", "TT-2", 0, "", "", "", true, "", true, false, true, "100000", "token-prov-1", "0", "token-issuer-1", "zero", ""}, "", nil},
		{"token", true, true, "Mint token - Not allow to mint if was Fixed Supply", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "", "TT-2", 8, "acc-40", "carlo", "", true, "100", true, false, true, "100000", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Transfer fixed supply token - Happy path", "acc-40", "100000000cin", 0, TokenInfo{"transfer", "", "", "", "TT-2", 8, "acc-40", "eve", "", false, "", true, false, false, "0", "", "", "", "", ""}, "", nil},

		{"token", true, true, "Burn token - Token not existed", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-3", 0, "acc-40", "carlo", "", true, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Burn token - Invalid owner", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2", 0, "nago", "", "", true, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Burn token - Not enough balance for TT-2", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2", 0, "acc-40", "", "", true, "", true, false, true, "100000000", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Burn token - Fixed Supply Happy path", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2", 0, "acc-40", "", "", true, "", true, false, false, "200", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Burn token - Allow to burn same TT-2 again", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2", 0, "acc-40", "", "", true, "", true, false, true, "100", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Burn token - Not allow to burn if was approved and burnable", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2", 0, "nago", "", "", true, "", true, false, true, "1", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Freeze token - Happy path for TT-2", "token-auth-1", "0cin", 0, TokenInfo{"freeze", "", "", "", "TT-2", 0, "", "", "", true, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, true, "Burn token - Not allow to burn if was approved and frozen", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2", 0, "nago", "", "", true, "", true, true, false, "1", "", "", "", "", ""}, "", nil},

		// Create-Approval-Burn : using Dynamic Supply
		{"token", true, true, "Create token - Use dynamic supply with long memo", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-2b", "TT-2b", 8, "acc-40", "", "", false, "100000", true, true, true, "", "", "", "", "", ""}, "long memo:123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890", nil},
		{"token", false, false, "Create token - Use dynamic supply with memo", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-2b", "TT-2b", 8, "acc-40", "", "", false, "200000", true, true, true, "", "", "", "", "", ""}, "1234567890", nil},
		{"token", false, false, "Approve token - Happy path for TT-2b", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "", "TT-2b", 0, "", "", "", false, "", true, true, true, "200000", "token-prov-1", "0", "token-issuer-1", "zero", ""}, "", nil},
		{"token", false, false, "Mint token - Allow to mint if was Dynamic Supply", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "", "TT-2b", 8, "acc-40", "carlo", "", false, "0", true, true, true, "100000", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Burn token - Token not existed", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-3b", 0, "acc-40", "carlo", "", false, "", true, true, true, "200", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Burn token - Invalid owner", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2b", 0, "nago", "", "", false, "", true, true, true, "200", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Burn token - Dynamic Supply of due to not enough balance for TT-2b", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2b", 0, "acc-40", "", "", false, "", true, true, true, "101", "", "", "", "", ""}, "", nil},

		{"token", true, true, "Burn token - Dynamic Supply of due to not enough balance for TT-2b", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2b", 0, "acc-40", "", "", false, "", true, false, true, "201", "", "", "", "", ""}, "", nil},       //KIV
		{"token", true, true, "Burn token - Dynamic Supply of due to not enough balance for TT-2b again", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2b", 0, "acc-40", "", "", false, "", true, false, true, "301", "", "", "", "", ""}, "", nil}, //KIV
		{"token", false, false, "Freeze token - Happy path for TT-2b", "token-auth-1", "0cin", 0, TokenInfo{"freeze", "", "", "", "TT-2b", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, true, "Burn token - Not allow to burn if was approved and frozen", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2b", 0, "nago", "", "", false, "", true, false, true, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Mint token - Not allow to mint if was approved and frozen", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "", "TT-2b", 8, "acc-40", "carlo", "", false, "1000", true, false, true, "0", "", "", "", "", ""}, "", nil},

		// Create - Approval - Approve transfer token ownership - Accept token ownership: Approve transfer token ownership
		{"token", false, false, "Create token - Happy Path TT-6 which for Approve transfer token ownership purpose", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-6", "TT-6", 8, "acc-40", "", "", false, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token - Happy path for TT-6", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-6", "TT-6", 0, "", "", "", false, "", false, false, true, "0", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", false, false, "Mint token - Happy path for TT-6", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "", "TT-6", 8, "acc-40", "carlo", "", false, "100000", true, false, false, "100", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Transfer token - Happy path for TT-6", "acc-40", "100000000cin", 0, TokenInfo{"transfer", "", "", "", "TT-6", 8, "acc-40", "eve", "", false, "", true, false, false, "0", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Transfer token ownership - Happy path", "acc-40", "100000000cin", 0, TokenInfo{"transfer-ownership", "", "", "", "TT-6", 8, "acc-40", "carlo", "", false, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Approve transfer token ownership - Error for TT-6 due to Invalid Authorised Signer", "eve", "0cin", 0, TokenInfo{"verify-transfer-tokenOwnership", "", "", "", "TT-6", 0, "", "", "", false, "", true, false, false, "", "token-prov-2", "0", "token-issuer-1", "", "APPROVE_TRANFER_TOKEN_OWNERSHIP"}, "", nil},
		{"token", true, true, "Approve transfer token ownership - Error for TT-6 due to Invalid Token Provider", "token-auth-1", "0cin", 0, TokenInfo{"verify-transfer-tokenOwnership", "", "", "", "TT-6", 0, "", "", "", false, "", true, false, false, "", "eve", "0", "token-issuer-1", "", "APPROVE_TRANFER_TOKEN_OWNERSHIP"}, "", nil},
		{"token", true, true, "Approve transfer token ownership - Error for TT-6 due to Invalid Token Issuer", "token-auth-1", "0cin", 0, TokenInfo{"verify-transfer-tokenOwnership", "", "", "", "TT-6", 0, "", "", "", false, "", true, false, false, "", "token-prov-2", "0", "eve", "", "APPROVE_TRANFER_TOKEN_OWNERSHIP"}, "", nil},
		// try again with the valid scenarios
		{"token", false, false, "Approve transfer token ownership - Happy-path for TT-6", "token-auth-2", "0cin", 0, TokenInfo{"verify-transfer-tokenOwnership", "", "", "", "TT-6", 0, "", "", "", false, "", true, false, false, "", "token-prov-2", "0", "token-issuer-2", "", "APPROVE_TRANFER_TOKEN_OWNERSHIP"}, "", nil},
		{"token", false, false, "Accept token ownership - Happy-path for TT-6", "carlo", "100000000cin", 0, TokenInfo{"accept-ownership", "", "", "", "TT-6", 8, "acc-40", "carlo", "", false, "", true, false, false, "1", "", "", "", "", ""}, "", nil},

		// 1. burnable flag
		{"token", false, false, "Create token - Happy Path (Burnable false)", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken123", "TT123", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Approve token - Happy path (Burnable false)", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken123", "TT123", 8, "", "", "", false, "", false, false, false, "1", "token-prov-1", "0", "token-issuer-1", "default", ""}, "", nil},
		{"token", true, true, "Burn token (TestToken123)", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT123", 0, "acc-40", "", "", false, "", false, false, true, "1", "", "", "", "", ""}, "", nil},

		// Base on : validation.go
		{"token", true, true, "Create token - Error due to Invalid fee amount", "acc-40", "1cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken-validate01", "TFT-validate01", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "Create token - Error due to Invalid fee amount.", nil},
		{"token", true, true, "Create token - Error due to Insufficient balance to pay for application fee", "acc-40", "10000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken-validate01", "TFT-validate01", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "Create token - Error due to Insufficient balance to pay for application fee.", nil},
		{"token", true, true, "Create token - Error due to Fee collector invalid", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "bob", "TestToken-validate01", "TFT-validate01", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "Create token - Error due to Fee collector invalid.", nil},
		{"token", false, false, "Create token TFT-validate01 - Happy Path, commit.", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken-validate01", "TFT-validate01", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "Create token TFT-validate01 - Happy Path.", nil},
		{"token", true, true, "Create token - Error due to Already existed", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken-validate01", "TFT-validate01", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "Create token - Error due to Already existed.", nil},

		{"token", true, true, "Approve token - Error due to Not authorised to approve", "acc-40", "100000000cin", 0, TokenInfo{"approve", "10000000", "mostafa", "TestToken-validate01", "TFT-validate01", 8, "", "", "", true, "", true, false, true, "2", "token-prov-1", "0", "token-issuer-1", "default", ""}, "MEMO: XXXX", nil},
		{"token", true, true, "Approve token - Error due to Token does not exist", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "mostafa", "TestToken-validate002", "TFT-validate002", 8, "", "", "", true, "", true, false, true, "2", "token-prov-1", "0", "token-issuer-1", "default", ""}, "MEMO: XXXX", nil},
		{"token", true, true, "Approve token - Error due to Fee setting in not valid 'default-validate01' ", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-validate01", "TFT-validate01", 8, "", "", "", false, "", true, false, true, "2", "token-prov-1", "0", "token-issuer-1", "default-validate01", ""}, "MEMO: XXXX", nil},

		{"token", false, false, "Approve token - Happy path for TFT-validate01", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-validate01", "TFT-validate01", 8, "", "", "", false, "", true, false, true, "2", "token-prov-1", "0", "token-issuer-1", "default", ""}, "MEMO: XXXX", nil},
		{"token", true, true, "Unfreeze token - Approved but not yet freeze", "token-auth-1", "0cin", 0, TokenInfo{"unfreeze", "", "", "TestToken-validate01", "TFT-validate01", 8, "acc-40", "eve", "", false, "", true, false, false, "2", "token-prov-1", "0", "token-issuer-1", "", ""}, "MEMO: XXXX", nil},
		{"token", false, false, "Freeze token - Happy path for TFT-validate01", "token-auth-1", "0cin", 0, TokenInfo{"freeze", "", "", "TestToken-validate01", "TFT-validate01", 8, "", "", "", false, "", true, false, true, "2", "token-prov-1", "0", "token-issuer-1", "", ""}, "MEMO: XXXX", nil},
		{"token", true, true, "Approve token - Error due to Token is frozen", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-validate01", "TFT-validate01", 8, "", "", "", false, "", true, false, true, "2", "token-prov-1", "0", "token-issuer-1", "default", ""}, "MEMO: XXXX", nil},
		{"token", true, true, "Approve token - Error due to Token already approved", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "TestToken-validate01", "TFT-validate01", 8, "", "", "", false, "", true, false, true, "2", "token-prov-1", "0", "token-issuer-1", "default", ""}, "MEMO: XXXX", nil},

		{"token", true, true, "Transfer token - Error due to Token is frozen", "acc-40", "100000000cin", 0, TokenInfo{"transfer", "", "", "TestToken-validate01", "TFT-validate01", 8, "acc-40", "eve", "", false, "", true, true, false, "100000", "", "", "", "", ""}, "MEMO: XXXX", nil},

		{"token", true, true, "Mint token - Error due to Token is frozen", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "TestToken-validate01", "TFT-validate01", 8, "acc-40", "carlo", "", false, "100000", true, true, false, "2", "", "", "", "default", ""}, "MEMO: XXXX", nil},

		{"token", true, true, "Burn token - Error due to Token is frozen", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "TestToken-validate01", "TFT-validate01", 8, "acc-40", "", "", true, "", true, true, false, "2", "", "", "", "", ""}, "MEMO: XXXX", nil},

		{"token", true, true, "Transfer token ownership - Invalid token owner", "acc-29", "100000000cin", 0, TokenInfo{"transfer-ownership", "", "", "TestToken-validate01", "TFT-validate01", 8, "acc-29", "carlo", "", false, "", true, true, false, "2", "", "", "", "", ""}, "MEMO: XXXX", nil},

		{"token", true, true, "Accept token ownership - Invalid new token owner", "acc-29", "100000000cin", 0, TokenInfo{"accept-ownership", "", "", "TestToken-validate01", "TFT-validate01", 8, "acc-29", "carlo", "", false, "", true, true, false, "2", "", "", "", "", ""}, "MEMO: XXXX", nil},
	}

	return tcs
}

func makeCreateFungibleTokenMsg(t *testing.T, name, symbol, metadata, owner, maxSupply, applicationFee, tokenFeeCollector string, decimals int, fixedSupply bool) sdkTypes.Msg {

	// create new token
	ownerAddr := tKeys[owner].addr
	maxSupplyUint := sdkTypes.NewUintFromString(maxSupply)
	fee := token.Fee{
		To:    tKeys[tokenFeeCollector].addr,
		Value: applicationFee,
	}
	msgCreateFungibleToken := token.NewMsgCreateFungibleToken(symbol, decimals, ownerAddr, name, fixedSupply, maxSupplyUint, metadata, fee)

	return msgCreateFungibleToken
}

func makeApproveFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, status string, burnable bool, feeSettingName string) sdkTypes.Msg {

	providerAddr := tKeys[provider].addr

	var tokenFee = []token.TokenFee{
		{
			Action:  "transfer",
			FeeName: feeSettingName,
		},
		{
			Action:  "mint",
			FeeName: feeSettingName,
		},
		{
			Action:  "burn",
			FeeName: feeSettingName,
		},
		{
			Action:  "transferOwnership",
			FeeName: feeSettingName,
		},
		{
			Action:  "acceptOwnership",
			FeeName: feeSettingName,
		},
	}
	tokenDoc := token.NewToken(providerAddr, providerNonce, status, symbol, burnable, tokenFee)

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := token.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []token.Signature
	signature := token.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetFungibleTokenStatus := token.NewMsgSetFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetFungibleTokenStatus
}

//module of transfer
func makeTransferFungibleTokenMsg(t *testing.T, owner string, newOwner string, symbol string, transferAmountOfToken string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr
	transferAmountOfTokenUint := sdkTypes.NewUintFromString(transferAmountOfToken)

	msgTransferPayload := token.NewMsgTransferFungibleToken(symbol, transferAmountOfTokenUint, ownerAddr, newOwnerAddr)
	return msgTransferPayload

}

//module of mint
func makeMintFungibleTokenMsg(t *testing.T, owner string, newOwner string, symbol string, mintAmountOfToken string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr
	mintAmountOfTokenUint := sdkTypes.NewUintFromString(mintAmountOfToken)

	msgMintPayload := token.NewMsgIssueFungibleAsset(ownerAddr, symbol, newOwnerAddr, mintAmountOfTokenUint)
	return msgMintPayload

}

//module of burn
func makeBurnFungibleTokenMsg(t *testing.T, owner string, symbol string, burnAmountOfToken string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	burnAmountOfTokenUint := sdkTypes.NewUintFromString(burnAmountOfToken)

	msgTransferPayload := token.NewMsgBurnFungibleToken(symbol, burnAmountOfTokenUint, ownerAddr)
	return msgTransferPayload

}

//moduel of transferOwnership
func makeTransferFungibleTokenOwnershipMsg(t *testing.T, owner string, newOwner string, symbol string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr

	msgTransferOwnershipPayload := token.NewMsgTransferFungibleTokenOwnership(symbol, ownerAddr, newOwnerAddr)
	return msgTransferOwnershipPayload

}

//module of acceptOwnership
func makeAcceptFungibleTokenOwnershipMsg(t *testing.T, newOwner string, symbol string) sdkTypes.Msg {

	fromAddr := tKeys[newOwner].addr

	msgAcceptOwnershipPayload := token.NewMsgAcceptFungibleTokenOwnership(symbol, fromAddr)
	return msgAcceptOwnershipPayload

}

func makeFreezeFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, burnable bool) sdkTypes.Msg {

	status := "FREEZE"
	providerAddr := tKeys[provider].addr

	tokenDoc := token.NewToken(providerAddr, providerNonce, status, symbol, burnable, nil) // status : FREEZE / UNFREEZE

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := token.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []token.Signature
	signature := token.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetFungibleTokenStatus := token.NewMsgSetFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetFungibleTokenStatus
}

func makeUnfreezeFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, burnable bool) sdkTypes.Msg {

	status := "UNFREEZE"
	providerAddr := tKeys[provider].addr

	tokenDoc := token.NewToken(providerAddr, providerNonce, status, symbol, burnable, nil) // status : FREEZE / UNFREEZE

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := token.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []token.Signature
	signature := token.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetFungibleTokenStatus := token.NewMsgSetFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetFungibleTokenStatus
}

func makeVerifyTransferTokenOwnership(t *testing.T, signer, provider, providerNonce, issuer, symbol, action string) sdkTypes.Msg {

	providerAddr := tKeys[provider].addr

	// burnable and tokenfees is not in used for verifying transfer token status, we just set it to false and leave it empty.
	verifyTransferTokenOwnershipDoc := token.NewToken(providerAddr, providerNonce, action, symbol, false, []token.TokenFee{})

	// provider sign
	verifyTransferTokenOwnershipDocBz, err := tCdc.MarshalJSON(verifyTransferTokenOwnershipDoc)
	require.NoError(t, err)
	signedVerifyTransferTokenOwnershipDoc, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(verifyTransferTokenOwnershipDocBz))
	require.NoError(t, err)

	verifyTransferTokenOwnershipPayload := token.NewPayload(*verifyTransferTokenOwnershipDoc, tKeys[provider].pub, signedVerifyTransferTokenOwnershipDoc)

	// issuer sign
	verifyTransferPayloadBz, err := tCdc.MarshalJSON(verifyTransferTokenOwnershipPayload)
	require.NoError(t, err)
	signedVerifyTransferPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(verifyTransferPayloadBz))
	require.NoError(t, err)

	var signatures []token.Signature
	signature := token.NewSignature(tKeys[issuer].pub, signedVerifyTransferPayloadBz)
	signatures = append(signatures, signature)

	msgVerifyTransferTokenOwnership := token.NewMsgSetFungibleTokenStatus(tKeys[signer].addr, *verifyTransferTokenOwnershipPayload, signatures)

	return msgVerifyTransferTokenOwnership
}
