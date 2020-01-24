package tests

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/fee"
	"github.com/maxonrow/maxonrow-go/x/nameservice"

	"github.com/maxonrow/maxonrow-go/x/bank"
	"github.com/maxonrow/maxonrow-go/x/kyc"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type bankInfo struct {
	from   string
	to     string
	amount string
}

type kycInfo struct {
	authorised string
	issuer     string
	provider   string
	action     string
	from       string
	signer     string
	data       string
	nonce      string
}

type feeInfo struct {
	function   string
	name       string
	assignee   string
	multiplier string
	min        string
	max        string
	percentage string
	issuer     string
}

type testCase struct {
	msgType       string
	checkFailed   bool
	deliverFailed bool
	desc          string
	signer        string
	fees          string
	gas           uint64
	msgInfo       interface{}
	memo          string
	hash          []byte
}

type NameServiceInfo struct {
	Action         string
	Name           string
	From           string
	ApplicationFee string
	FeeCollector   string
	Provider       string
	ProviderNonce  string
	Issuer         string
	approved       string
}

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

func TestTxs(t *testing.T) {

	//acc1 := Account(tKeys["alice"].addrStr)
	val1 := Validator(tValidator)
	fmt.Println(val1)
	//bal1 := acc1.GetCoins()[0]

	_, err := time.ParseDuration("60s")
	if err != nil {
		fmt.Println("timeout", err)
	}

	var proposalTitleOutOfLength string
	proposalTitleOutOfLength = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333"

	var proposalDescOutOfLength string
	proposalDescOutOfLength = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333----aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzzaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccc--aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333----aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzzaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccc--aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---" +
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---" +
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccc---xxxxxxxxxxxxxxxxxxxxxxxyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyzzzzzzzzzzzzzzzzzzz---1111112222222233333---aaaaaaaaaaaaaaaa-eeeeee"

	tcs := []*testCase{

		//---------------------------------------------------------------------------------------------------------------------
		// bank ---------------------------------------------------------------------------------------------------------------
		//---------------------------------------------------------------------------------------------------------------------
		{"bank", false, false, "sending 1 cin", "alice", "200000000cin", 0, bankInfo{"alice", "bob", "1cin"}, "", nil},
		{"bank", true, true, "sending 0 cin", "alice", "200000000cin", 0, bankInfo{"alice", "bob", "0cin"}, "", nil},
		{"bank", true, true, "insufficient amount", "carlo", "1000000000cin", 0, bankInfo{"carlo", "eve", "999999999999000000001cin"}, "", nil},
		{"bank", false, false, "transffer all coins", "gohck", "1000000000cin", 0, bankInfo{"gohck", "eve", "999999999999000000000cin"}, "", nil},
		{"bank", true, true, "sending 1 abc", "alice", "200000000cin", 0, bankInfo{"alice", "bob", "1abc"}, "", nil},
		{"bank", true, true, "sending 1 cin & 1 abc", "alice", "200000000cin", 0, bankInfo{"alice", "bob", "1cin, 1abc"}, "", nil},
		{"bank", false, false, "sending 1 mxw", "alice", "1000000000cin", 0, bankInfo{"alice", "bob", "1000000000000000000cin"}, "", nil},
		{"bank", false, false, "more fee", "alice", "100000000000cin", 0, bankInfo{"alice", "bob", "1cin"}, "", nil},
		{"bank", false, false, "more fee", "alice", "1000000000000000000cin", 0, bankInfo{"alice", "bob", "1cin"}, "", nil},
		{"bank", false, false, "with memo", "alice", "200000000cin", 0, bankInfo{"alice", "bob", "1cin"}, "alice to bob", nil},
		{"bank", true, true, "no fee", "alice", "", 0, bankInfo{"alice", "bob", "1cin"}, "", nil},
		{"bank", true, true, "invalid denom for fee", "alice", "200000000abc", 0, bankInfo{"alice", "bob", "1cin"}, "", nil},
		{"bank", true, true, "less fee", "alice", "1cin", 0, bankInfo{"alice", "bob", "1cin"}, "", nil},
		{"bank", true, true, "zero fee", "alice", "0cin", 0, bankInfo{"alice", "bob", "1cin"}, "", nil},
		{"bank", true, true, "wrong fee", "alice", "200000000abc", 0, bankInfo{"alice", "bob", "1cin"}, "", nil},
		{"bank", true, true, "zero amount", "alice", "1cin", 0, bankInfo{"alice", "bob", "0cin"}, "", nil},
		{"bank", true, true, "wrong signer", "eve", "200000000cin", 0, bankInfo{"alice", "bob", "1cin"}, "", nil},
		{"bank", true, true, "no sender", "alice", "200000000cin", 0, bankInfo{"nope", "bob", "1cin"}, "", nil},
		{"bank", true, true, "no receiver", "alice", "200000000cin", 0, bankInfo{"alice", "nope", "1cin"}, "", nil},
		{"bank", true, true, "no amount", "alice", "", 0, bankInfo{"alice", "bob", "1cin"}, "", nil},
		{"bank", true, true, "wrong gas", "alice", "200000000cin", 1, bankInfo{"alice", "bob", "1cin"}, "", nil},
		{"bank", true, true, "long memo", "alice", "200000000cin", 0, bankInfo{"alice", "bob", "1cin"}, "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890", nil},
		{"bank", false, false, "to non-kyc account-commit", "alice", "200000000cin", 0, bankInfo{"alice", "josephin", "2000000000cin"}, "", nil},
		{"bank", true, true, "from non-kyc account", "josephin", "200000000cin", 0, bankInfo{"josephin", "bob", "1cin"}, "", nil},

		//---------------------------------------------------------------------------------------------------------------------
		// kyc ----------------------------------------------------------------------------------------------------------------
		//---------------------------------------------------------------------------------------------------------------------
		//Whitelist
		{"kyc", true, true, "Doing kyc - INVALID SIGNER-1", "bob", "0cin", 0, kycInfo{"bob", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - INVALID SIGNER-2", "kyc-issuer-1", "0cin", 0, kycInfo{"kyc-issuer-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - MISMATCH SIGNER", "kyc-auth-2", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1234", "0"}, "", nil},
		{"kyc", false, false, "Doing kyc - HAPPY PATH-commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - EMPTY DATA", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "acc-30", "acc-30", "", "0"}, "", nil},

		// goh - last-time
		{"kyc", false, false, "Doing kyc - HAPPY PATH for dont-use-this-1", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "dont-use-this-1", "dont-use-this-1", "testKyc1251", "0"}, "", nil},

		{"kyc", true, true, "Doing kyc - KYC AGAIN!", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1234", "0"}, "", nil}, //code 1
		{"kyc", false, false, "Doing kyc - KYC AGAIN, DIFFERENT KYC DATA-commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc12399", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - DUPLICATE KYC DATA", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "acc-23", "acc-23", "testKyc12399", "0"}, "", nil},

		{"kyc", true, true, "Doing kyc - NOT ISSUER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "jeansoon", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc12345", "0"}, "", nil},                        //code 4
		{"kyc", true, true, "Doing kyc - NO ISSUER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "nope", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc12345", "0"}, "", nil},                             //code 4
		{"kyc", true, true, "Doing kyc - NOT AUTHORISED", "jeansoon", "0cin", 0, kycInfo{"jeansoon", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc12346", "0"}, "", nil},                    //code 1
		{"kyc", true, true, "Doing kyc - NOT AUTHORISED, BUT WHITELISTED ADDRESS", "alice", "0cin", 0, kycInfo{"alice", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc12346", "0"}, "", nil}, //code 1
		{"kyc", true, true, "Doing kyc - NOT PROVIDER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "jeansoon", "whitelist", "josephin", "josephin", "testKyc12347", "0"}, "", nil},                    //code 4
		{"kyc", true, true, "Doing kyc - NO PROVIDER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "nope", "whitelist", "josephin", "josephin", "testKyc12347", "0"}, "", nil},                         //code 4
		{"kyc", true, true, "Doing kyc - WRONG NONCE", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1238", "1"}, "", nil},
		{"kyc", true, true, "Doing kyc - NOT PROVIDER, NOT ISSUER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "jeansoon", "acc-19", "whitelist", "josephin", "josephin", "testKyc1239", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - TWO PROVIDER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-prov-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1239", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - TWO ISSUER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-issuer-1", "whitelist", "josephin", "josephin", "testKyc1239", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - EMPTY DATA", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "jeansoon", "acc-19", "whitelist", "acc-24", "acc-24", "", "0"}, "", nil},
		{"kyc", false, false, "Doing kyc - PAY FEE-commit", "kyc-auth-1", "1000cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc123456", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - WRONG GAS", "kyc-auth-1", "0cin", 1, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "josephin", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - NOPE ADDR", "kyc-auth-1", "0cin", 1, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "nope", "acc-23", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - INVALID FROM SIGNATURE", "kyc-auth-1", "0cin", 1, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "acc-23", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - INVALID AUTHORISED SIGNATURE-1", "kyc-auth-2", "0cin", 1, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "acc-23", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - INVALID AUTHORISED SIGNATURE-2", "kyc-issuer-1", "0cin", 1, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "acc-23", "testKyc1234", "0"}, "", nil},
		{"kyc", true, true, "Doing kyc - INVALID AUTHORISED SIGNATURE", "kyc-prov-1", "0cin", 1, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "josephin", "acc-23", "testKyc1234", "0"}, "", nil},

		// It should works now
		{"bank", false, false, "sending after whitelisting an account", "josephin", "200000000cin", 0, bankInfo{"josephin", "bob", "2cin"}, "", nil},

		//RevokeWhitelist
		// revoke whitelist needs provider signature(RevokeKycData)
		// revoke whitelist needs issuer signature(RevokeKycPayload)
		// revoke whitelist needs authorised signature(tx)
		{"kyc", true, true, "Undoing kyc - INVALID SIGNER", "bob", "0cin", 0, kycInfo{"bob", "kyc-issuer-1", "jeansoon", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - MISMATCH SIGNER", "kyc-auth-2", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "jeansoon", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - NOT PROVIDER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "jeansoon", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - NO PROVIDER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "nope", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - NOT ISSUER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "jeansoon", "kyc-prov-1", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - NO ISSUER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "nope", "kyc-prov-1", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - TWO ISSUER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-issuer-1", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - TWO PROVIDER", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-prov-1", "kyc-prov-1", "revokeWhitelist", "josephin", "", "", ""}, "", nil},
		{"kyc", true, true, "Undoing kyc - REVOKE SOMEONE IS NOT WHITELISTED", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "revokeWhitelist", "acc-23", "", "", ""}, "", nil},
		{"kyc", false, false, "Undoing kyc - HAPPY PATH-commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "revokeWhitelist", "josephin", "", "", ""}, "", nil},

		// It should fail now
		{"bank", true, true, "sending after revoking an account", "josephin", "200000000cin", 0, bankInfo{"josephin", "bob", "1cin"}, "", nil},

		{"bank", false, false, "receiving after revoking an account", "alice", "200000000cin", 0, bankInfo{"alice", "josephin", "1cin"}, "", nil},

		//---------------------------------------------------------------------------------------------------------------------
		// fee ----------------------------------------------------------------------------------------------------------------
		//---------------------------------------------------------------------------------------------------------------------
		{"fee", true, true, "Updating default fee settings-invalid signer", "alice", "0cin", 0, feeInfo{"sys-fee", "default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", true, true, "Updating default fee settings-wrong issuer", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "default", "", "", "400000000cin", "2000000000cin", "0.001", "alice"}, "", nil},
		{"fee", false, false, "Updating default fee settings-commit", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "default", "", "", "400000000cin", "2000000000cin", "0.001", "fee-auth"}, "", nil},
		{"fee", true, true, "Updating default fee settings-wrong fee", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "wrong", "", "", "888888888cin", "111111111cin", "0.001", "fee-auth"}, "", nil},
		{"fee", true, true, "Add double fee settings-wrong signer", "alice", "0cin", 0, feeInfo{"sys-fee", "double", "", "", "800000000cin", "4000000000cin", "0.002", "fee-auth"}, "", nil},
		{"fee", false, false, "Add double fee settings-commit", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "double", "", "", "800000000cin", "4000000000cin", "0.002", "fee-auth"}, "", nil},
		{"fee", false, false, "Add double fee settings-commit - extra fee", "fee-auth", "100000cin", 0, feeInfo{"sys-fee", "double1", "", "", "800000000cin", "4000000000cin", "0.002", "fee-auth"}, "", nil},
		{"fee", true, true, "assign double to bank txs-wrong issuer", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "double", "bank-send", "", "", "", "", "bob"}, "", nil},
		{"fee", true, true, "assign double to bank txs-wrong signer", "bob", "0cin", 0, feeInfo{"assign-msg", "double", "bank-send", "", "", "", "", "fee-auth"}, "", nil},
		{"fee", false, true, "assign double to bank txs-wrong name", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "wrong", "bank-send", "", "", "", "", "fee-auth"}, "", nil},
		{"fee", true, true, "assign double to bank txs-empty name", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "", "bank-send", "", "", "", "", "fee-auth"}, "", nil},
		//??? No way to check msg-type is valid?
		{"fee", false, false, "assign double to bank txs-wrong msg", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "double", "wrong-send", "", "", "", "", "fee-auth"}, "", nil},
		{"fee", true, true, "assign double to bank txs-empty msg", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "double", "", "", "", "", "", "fee-auth"}, "", nil},
		{"fee", false, false, "assign double to bank txs-commit", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "double", "bank-send", "", "", "", "", "fee-auth"}, "", nil},
		// should pay double
		{"kyc", false, false, "Doing kyc - yk - commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "yk", "yk", "testKyc12345678", "0"}, "", nil},
		{"bank", true, true, "sending after updating fee- wrong fee", "yk", "500000000cin", 0, bankInfo{"yk", "gohck", "1cin"}, "", nil},

		{"bank", false, false, "sending after updating fee", "yk", "800000000cin", 0, bankInfo{"yk", "mostafa", "1000cin"}, "", nil},

		// assign zero fee to an account
		{"kyc", false, false, "Doing kyc - mostafa - commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "mostafa", "mostafa", "testKyc123456789", "0"}, "", nil},
		{"fee", true, true, "assign zero-fee to nil address", "fee-auth", "0cin", 0, feeInfo{"assign-acc", "zero", "nope", "", "", "", "", "fee-auth"}, "", nil},
		{"fee", false, false, "assign zero-fee to mostafa-commit", "fee-auth", "0cin", 0, feeInfo{"assign-acc", "zero", "mostafa", "", "", "", "", "fee-auth"}, "", nil},

		// zero fee
		{"bank", false, false, "sending after updating acc-fee", "mostafa", "0cin", 0, bankInfo{"mostafa", "bob", "1cin"}, "", nil},

		// updating multiplier
		{"fee", true, true, "invalid signer", "eve", "0cin", 0, feeInfo{"multiplier", "", "", "2", "", "", "", "fee-auth"}, "", nil},
		{"fee", true, true, "mismatch signer", "fee-auth", "0cin", 0, feeInfo{"multiplier", "", "", "2", "", "", "", "eve"}, "", nil},
		{"fee", true, true, "invalid multiplier", "fee-auth", "0cin", 0, feeInfo{"multiplier", "", "", "1.*", "", "", "", "fee-auth"}, "", nil},
		{"fee", false, false, "updating multiplier", "fee-auth", "0cin", 0, feeInfo{"multiplier", "", "", "2", "", "", "", "fee-auth"}, "", nil},
		{"fee", false, false, "updating multiplier", "fee-auth", "0cin", 0, feeInfo{"multiplier", "", "", "1.0005", "", "", "", "fee-auth"}, "", nil},
		// Reset default fee setting to genesis default
		{"fee", false, false, "Updating default fee settings-commit", "fee-auth", "0cin", 0, feeInfo{"sys-fee", "default", "", "", "100000000cin", "1000000000cin", "0.001", "fee-auth"}, "", nil},

		//---------------------------------------------------------------------------------------------------------------------
		// maintenance --------------------------------------------------------------------------------------------------------
		//---------------------------------------------------------------------------------------------------------------------
		// Add Proposal [ challenge ]
		{"maintenance", true, true, "Is invalid maintainer: Add fee provider address.", "nago", "0cin", 0, MaintenanceInfo{"add", "Add fee-provider address", "Add a party as fee-provider address", "fee", "nago", "", "", FeeCollector{}, "nago", ""}, "", nil},
		{"maintenance", true, true, "Zero-length of Proposal Title: Add fee provider address.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "", "Add a party as fee-provider address", "fee", "nago", "", "", FeeCollector{}, "maintainer-1", ""}, "", nil},
		{"maintenance", true, true, "Proposal Title was out-of-length: Add fee provider address.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", proposalTitleOutOfLength, "Add a party as fee-provider address", "fee", "maintainer-1", "", "", FeeCollector{}, "nago", ""}, "", nil},
		{"maintenance", true, true, "Zero-length of Proposal Description: Add fee provider address.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add fee-provider address", "", "fee", "nago", "", "", FeeCollector{}, "maintainer-1", ""}, "", nil},
		{"maintenance", true, true, "Proposal Description was out-of-length: Add fee provider address.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add fee-provider address", proposalDescOutOfLength, "fee", "nago", "", "", FeeCollector{}, "nago", ""}, "", nil},
		{"maintenance", true, true, "Signer and proposer is different", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add fee-provider address", "Add a party as fee-provider address", "fee", "nago", "", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", true, true, "Invalid proposal type: Add nameservice provider address.", "nago", "0cin", 0, MaintenanceInfo{"add", "Add nameservice-provider address", "Add a party as nameservice-provider address", "nameservice-999", "nago", "", "maintainer-3", FeeCollector{}, "maintainer-3", ""}, "", nil},

		//------------------------------
		// Maintenance - Fee
		// Maintenance - Fee - Add(Proposal) [proposal-1, proposal-2, proposal-3, proposal-4]
		{"maintenance", false, false, "1. Proposal, add fee authorised address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add authorised address", "Add cmo as fee authorised address", "fee", "maintainer-1", "", "", FeeCollector{}, "maintainer-1", ""}, "", nil},
		{"maintenance", false, false, "2. Proposal, add fee issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add isser address", "Add fee issuer address", "fee", "", "maintainer-1", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "3. Proposal, add fee provider address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"add", "Add isser address", "Add fee issuer address", "fee", "", "", "maintainer-2", FeeCollector{}, "maintainer-3", ""}, "", nil},
		{"maintenance", false, false, "4. Proposal, add fee collector address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add token fee collector", "Add maintainer-1 as token fee collector", "fee", "", "", "", FeeCollector{Module: "token", Address: "maintainer-1"}, "maintainer-2", ""}, "", nil},

		// Maintenance - Fee - Cast Action [proposal-1, proposal-2, proposal-3, proposal-4]
		// Cast Action - Proposal 1
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 1, signer and caster different person.", "maintainer-1", "0cin", 0, CastAction{"maintainer-2", "approve", 1}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 1, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 1}, "", nil},
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 1, caster and signer not maintainer.", "yk", "0cin", 0, CastAction{"yk", "approve", 1}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 1, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 1}, "", nil},
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to inactive proposal 1.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 1}, "", nil},
		// Cast Action - Proposal 2
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 2, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 2}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 2, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 2}, "", nil},
		// Cast Action - Proposal 3
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 3, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 3}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 3, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 3}, "", nil},
		// Cast Action - Proposal 4
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 3, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 4}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 3, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 4}, "", nil},

		// Maintenance - Fee - Remove(Proposal) [proposal-5, proposal-6, proposal-7, proposal-8]
		{"maintenance", false, false, "5. Proposal, remove fee authorised address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"remove", "Remove authorised address", "Remove cmo as fee authorised address", "fee", "maintainer-1", "", "", FeeCollector{}, "maintainer-1", ""}, "", nil},
		{"maintenance", false, false, "6. Proposal, remove fee issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"remove", "Remove isser address", "Remove fee issuer address", "fee", "", "maintainer-1", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "7. Proposal, remove fee provider address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"remove", "Remove isser address", "Remove fee issuer address", "fee", "", "", "maintainer-2", FeeCollector{}, "maintainer-3", ""}, "", nil},
		{"maintenance", false, false, "8. Proposal, remove fee collector address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"remove", "Remove token fee collector", "Remove maintainer-1 as token fee collector", "fee", "", "", "", FeeCollector{Module: "token", Address: "maintainer-1"}, "maintainer-2", ""}, "", nil},

		// Maintenance - Fee - Cast Action [proposal-5, proposal-6, proposal-7, proposal-8]
		// Cast Action - Proposal 5
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 5, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 5}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 5, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 5}, "", nil},
		// Cast Action - Proposal 6
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 6, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 6}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 6, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 6}, "", nil},
		// Cast Action - Proposal 7
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 7, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 7}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 7, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 7}, "", nil},
		// Cast Action - Proposal 8
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 8, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 8}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 8, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 8}, "", nil},

		//---------------------------------
		// Maintenance - KYC
		// Maintenance - KYC [proposal-9, proposal-10, proposal-11]
		{"maintenance", false, false, "9. Proposal, add kyc authorised address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add KYC-authorised address", "Add a party as KYC-authorised address", "kyc", "maintainer-1", "", "", FeeCollector{}, "maintainer-1", ""}, "", nil},
		{"maintenance", false, false, "10. Proposal, add kyc issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add KYC-issuer address", "Add a party as KYC-issuer address", "kyc", "", "maintainer-2", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "11. Proposal, add kyc provider address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"add", "Add KYC-provider address", "Add a party as KYC-provider address", "kyc", "", "", "maintainer-3", FeeCollector{}, "maintainer-3", ""}, "", nil},

		// Cast Action - Proposal 9
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 9, signer and caster different person.", "maintainer-1", "0cin", 0, CastAction{"maintainer-2", "approve", 9}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 9, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 9}, "", nil},
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 9, caster and signer not maintainer.", "yk", "0cin", 0, CastAction{"yk", "approve", 9}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 9, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 9}, "", nil},
		//{"maintenance-cast-action", true, true, "(Approve)-Cast action to inactive proposal 9.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "disable", 9}, "", nil},//kiv
		// Cast Action - Proposal 10
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 10, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 10}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 10, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 10}, "", nil},
		// Cast Action - Proposal 11
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 11, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 11}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 11, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 11}, "", nil},

		// Maintenance - KYC - Remove(Proposal) [proposal-12, proposal-13, proposal-14]
		{"maintenance", false, false, "12. Proposal, remove kyc authorised address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"remove", "Remove KYC-authorised address", "Remove a party as KYC-authorised address", "kyc", "maintainer-1", "", "", FeeCollector{}, "maintainer-1", ""}, "", nil},
		{"maintenance", false, false, "13. Proposal, remove kyc issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"remove", "Remove KYC-issuer address", "Remove a party as KYC-issuer address", "kyc", "", "maintainer-2", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "14. Proposal, remove kyc provider address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"remove", "Remove KYC-provider address", "Remove a party as KYC-provider address", "kyc", "", "", "maintainer-3", FeeCollector{}, "maintainer-3", ""}, "", nil},

		// Cast Action for KYC - Proposal 12
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 12, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 12}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 12, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 12}, "", nil},
		// Cast Action for KYC - Proposal 13
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 13, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 13}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 13, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 13}, "", nil},
		// Cast Action for KYC - Proposal 14
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 14, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 14}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 14, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 14}, "", nil},
		{"maintenance-cast-action", false, false, "(Reject)-Cast action to proposal 14, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "reject", 14}, "", nil},

		//---------------------------------
		// Maintenance - Token
		// Maintenance - Token [proposal-15, proposal-16, proposal-17]
		{"maintenance", false, false, "15. Proposal, add token authorised address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"add", "Add token-authorised address", "Add a party as token-authorised address", "token", "maintainer-3", "", "", FeeCollector{}, "maintainer-3", ""}, "", nil},
		{"maintenance", false, false, "16. Proposal, add token issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add token-issuer address", "Add a party as token-issuer address", "token", "", "maintainer-2", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "17. Proposal, add token provider address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add token-provider address", "Add a party as token-provider address", "token", "", "", "maintainer-1", FeeCollector{}, "maintainer-1", ""}, "", nil},

		// Cast Action - Proposal 15
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 15, signer and caster different person.", "maintainer-1", "0cin", 0, CastAction{"maintainer-3", "approve", 15}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 15, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 15}, "", nil},
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 15, caster and signer not maintainer.", "yk", "0cin", 0, CastAction{"yk", "approve", 15}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 15, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 15}, "", nil},
		//{"maintenance-cast-action", true, true, "(Approve)-Cast action to inactive proposal 15.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "disable", 15}, "", nil},//kiv
		// Cast Action - Proposal 16
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 16, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 16}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 16, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 16}, "", nil},
		// Cast Action - Proposal 17
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 17, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 17}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 17, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 17}, "", nil},

		// Maintenance - Token - Remove(Proposal) [proposal-18, proposal-19, proposal-20]
		{"maintenance", false, false, "18. Proposal, remove token authorised address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"remove", "Remove token-authorised address", "Remove a party as token-authorised address", "token", "maintainer-1", "", "", FeeCollector{}, "maintainer-1", ""}, "", nil},
		{"maintenance", false, false, "19. Proposal, remove token issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"remove", "Remove token-issuer address", "Remove a party as token-issuer address", "token", "", "maintainer-2", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "20. Proposal, remove token provider address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"remove", "Remove token-provider address", "Remove a party as token-provider address", "token", "", "", "maintainer-3", FeeCollector{}, "maintainer-3", ""}, "", nil},

		// Cast Action - Proposal 18
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 18, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 18}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 18, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 18}, "", nil},
		{"maintenance-cast-action", false, false, "(Reject)-Cast action to proposal 18, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "reject", 18}, "", nil},
		// Cast Action - Proposal 19
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 19, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 19}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 19, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 19}, "", nil},
		{"maintenance-cast-action", false, false, "(Reject)-Cast action to proposal 19, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "reject", 19}, "", nil},
		// Cast Action - Proposal 20
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 20, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 20}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 20, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 20}, "", nil},

		//---------------------------------
		// Maintenance - Nameservice
		// Maintenance - Nameservice [proposal-21, proposal-22, proposal-23]
		{"maintenance", false, false, "21. Proposal, add nameservice authorised address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"add", "Add nameservice-authorised address", "Add a party as nameservice-authorised address", "nameservice", "maintainer-3", "", "", FeeCollector{}, "maintainer-3", ""}, "", nil},
		{"maintenance", false, false, "22. Proposal, add nameservice issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add nameservice-issuer address", "Add a party as nameservice-issuer address", "nameservice", "", "maintainer-2", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "23. Proposal, add nameservice provider address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add nameservice-provider address", "Add a party as nameservice-provider address", "nameservice", "", "", "maintainer-1", FeeCollector{}, "maintainer-1", ""}, "", nil},

		// Cast Action - Proposal 21
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 21, signer and caster different person.", "maintainer-1", "0cin", 0, CastAction{"maintainer-2", "approve", 21}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 21, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 21}, "", nil},
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 21, caster and signer not maintainer.", "yk", "0cin", 0, CastAction{"yk", "approve", 21}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 21, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 21}, "", nil},
		//{"maintenance-cast-action", true, true, "(Approve)-Cast action to inactive proposal 21.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "disable", 21}, "", nil},//kiv
		// Cast Action - Proposal 22
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 22, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 22}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 22, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 22}, "", nil},
		// Cast Action - Proposal 23
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 23, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 23}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 23, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 23}, "", nil},

		// Maintenance - Nameservice - Remove(Proposal) [proposal-24, proposal-25, proposal-26]
		{"maintenance", false, false, "24. Proposal, remove nameservice authorised address, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"remove", "Remove nameservice-authorised address", "Remove a party as nameservice-authorised address", "nameservice", "maintainer-3", "", "", FeeCollector{}, "maintainer-3", ""}, "", nil},
		{"maintenance", false, false, "25. Proposal, remove nameservice issuer address, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"remove", "Remove nameservice-issuer address", "Remove a party as nameservice-issuer address", "nameservice", "", "maintainer-2", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance", false, false, "26. Proposal, remove nameservice provider address, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"remove", "Remove nameservice-provider address", "Remove a party as nameservice-provider address", "nameservice", "", "", "maintainer-1", FeeCollector{}, "maintainer-1", ""}, "", nil},

		// Cast Action - Proposal 24
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 24, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 24}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 24, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 24}, "", nil},
		{"maintenance-cast-action", false, false, "(Reject)-Cast action to proposal 24, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "reject", 24}, "", nil},
		// Cast Action - Proposal 25
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 25, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 25}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 25, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 25}, "", nil},
		{"maintenance-cast-action", false, false, "(Reject)-Cast action to proposal 25, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "reject", 25}, "", nil},
		// Cast Action - Proposal 26
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 26, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 26}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 26, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 26}, "", nil},

		//---------------------------------
		// Maintenance - Validator-set
		// Maintenance - Validator-set [proposal-27, proposal-28, proposal-29]
		{"maintenance", false, false, "27. Proposal, add validator set-1, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"add", "Add validator-set", "Add a party for validator-set", "validator", "", "", "", FeeCollector{}, "maintainer-1", "mxwvalconspub1zcjduepq2vxnxwuzvf82w9mxhjuwm35q7e84pfglsexh5l0ffqz0ddfxjp5q8wjkgw"}, "", nil},
		{"maintenance", false, false, "28. Proposal, add validator set-2, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add validator-set", "Add a party for validator-set", "validator", "", "", "", FeeCollector{}, "maintainer-2", "mxwvalconspub1zcjduepqczwdy9dlmvazg3u3nml743xgprr2e82n2lt6wue5ycsga2nudvxq0avuc6"}, "", nil},
		{"maintenance", false, false, "29. Proposal, add validator set-3, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"add", "Add validator-set", "Add a party for validator-set", "validator", "", "", "", FeeCollector{}, "maintainer-3", "mxwvalconspub1zcjduepqvf9vf3cdxwtk65ya83q8uz36c8vqn5gylp3dmkghxjs253thve4qqzm5ca"}, "", nil},

		// Cast Action - Proposal 27
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 27, signer and caster different person.", "maintainer-1", "0cin", 0, CastAction{"maintainer-2", "approve", 27}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 27, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 27}, "", nil},
		{"maintenance-cast-action", true, true, "(Approve)-Cast action to proposal 27, caster and signer not maintainer.", "gohck", "0cin", 0, CastAction{"gohck", "approve", 27}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 27, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 27}, "", nil},
		// Cast Action - Proposal 28
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 28, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 28}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 28, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 28}, "", nil},
		// Cast Action - Proposal 29
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 29, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 29}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 29, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 29}, "", nil},

		// Maintenance - Validator-set - Remove(Proposal) [proposal-30, proposal-31, proposal-32]
		{"maintenance", false, false, "30. Proposal, remove validator set-1, Happy path.", "maintainer-1", "0cin", 0, MaintenanceInfo{"remove", "Remove validator-set", "Remove a party for validator-set", "validator", "", "", "", FeeCollector{}, "maintainer-1", "mxwvalconspub1zcjduepq2vxnxwuzvf82w9mxhjuwm35q7e84pfglsexh5l0ffqz0ddfxjp5q8wjkgw"}, "", nil},
		{"maintenance", false, false, "31. Proposal, remove validator set-2, Happy path.", "maintainer-2", "0cin", 0, MaintenanceInfo{"remove", "Remove validator-set", "Remove a party for validator-set", "validator", "", "", "", FeeCollector{}, "maintainer-2", "mxwvalconspub1zcjduepqczwdy9dlmvazg3u3nml743xgprr2e82n2lt6wue5ycsga2nudvxq0avuc6"}, "", nil},
		{"maintenance", false, false, "32. Proposal, remove validator set-3, Happy path.", "maintainer-3", "0cin", 0, MaintenanceInfo{"remove", "Remove validator-set", "Remove a party for validator-set", "validator", "", "", "", FeeCollector{}, "maintainer-3", "mxwvalconspub1zcjduepqvf9vf3cdxwtk65ya83q8uz36c8vqn5gylp3dmkghxjs253thve4qqzm5ca"}, "", nil},

		// Cast Action - Proposal 30
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 30, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 30}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 30, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 30}, "", nil},
		{"maintenance-cast-action", false, false, "(Reject)-Cast action to proposal 30, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "reject", 30}, "", nil},
		// Cast Action - Proposal 31
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 31, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 31}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 31, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 31}, "", nil},
		{"maintenance-cast-action", false, false, "(Reject)-Cast action to proposal 31, Happy path.", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "reject", 31}, "", nil},
		// Cast Action - Proposal 32
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 32, Happy path.", "maintainer-2", "0cin", 0, CastAction{"maintainer-2", "approve", 32}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to proposal 32, Happy path.", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 32}, "", nil},

		//---------------------------------------------------------------------------------------------------------------------
		// nameservice --------------------------------------------------------------------------------------------------------
		//---------------------------------------------------------------------------------------------------------------------
		// whitelist the fee collector first
		{"kyc", false, false, "Doing kyc - ns-feecollector - HAPPY PATH-commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "ns-feecollector", "ns-feecollector", "iamNS-FeeCollector", "0"}, "", nil},

		// add nameservice fee collector with maintenance.
		{"maintenance", false, false, "33. Proposal, add nameservice fee collector address, Happy path. commit", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add nameservice fee collector", "Add ns-feecollector as nameservice fee collector", "fee", "", "", "", FeeCollector{Module: "nameservice", Address: "ns-feecollector"}, "maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve ns-feecollector as nameservice fee collector, Happy path. commit", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 33}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve ns-feecollector as nameservice fee collector, Happy path. commit", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 33}, "", nil},

		// Create NameService
		{"kyc", false, false, "Doing kyc - nago - commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "nago", "nago", "nago-data", "0"}, "", nil},
		{"nameservice", true, true, "creating alias with name mxw-alias-invalid signer", "yk", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias", "nago", "10000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", false, false, "creating alias with name mxw-alias", "nago", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias", "nago", "10000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", false, true, "creating alias with name mxw-alias-again", "nago", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias2", "nago", "10000", "ns-feecollector", "", "", "", ""}, "", nil},

		{"nameservice", false, false, "creating alias with name mxw-alias-1", "acc-40", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias-1", "acc-40", "1000000000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", false, false, "creating alias with name mxw-alias-2", "carlo", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias-2", "carlo", "1000000000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", false, false, "creating alias with name MXW-ALIAS", "yk", "100000000cin", 0, NameServiceInfo{"create", "MXW-ALIAS", "yk", "10000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", false, false, "creating alias with name MXW-@LI@S", "dont-use-this-1", "100000000cin", 0, NameServiceInfo{"create", "MXW-@LI@S", "dont-use-this-1", "10000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", true, true, "creating alias with wrong free collector", "jeansoon", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias", "jeansoon", "1000000000", "acc-40", "", "", "", ""}, "", nil},
		{"nameservice", false, true, "Already existed  name mxw-alias", "nago", "100000000cin", 0, NameServiceInfo{"create", "new-alias", "nago", "100000000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", true, true, "create alias with zero system fee", "nago", "0cin", 0, NameServiceInfo{"create", "mxw-alias", "nago", "100000000", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", false, true, "create alias with zero application-fee, have pending", "nago", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias-nago", "nago", "0", "ns-feecollector", "", "", "", ""}, "", nil},
		{"nameservice", false, false, "create alias with zero application-fee", "mostafa", "100000000cin", 0, NameServiceInfo{"create", "mxw-alias-mostafa", "mostafa", "0", "ns-feecollector", "", "", "", ""}, "", nil},
		///{"nameservice", true, true, "create alias with non-kyc account", "gohck", "100000000cin", 0, NameServiceInfo{"create", "mxwone", "gohck", "100000000", "ns-feecollector", "", "", "", ""}, "", nil}, //kiv
		{"nameservice", true, true, "creating alias with signer and owner are different", "nago", "100000000cin", 0, NameServiceInfo{"create", "mxwone", "gohck", "100000000", "ns-feecollector", "", "", "", ""}, "", nil},

		//NameServie
		//Approve
		{"fee", false, false, "assign zero-fee to bank alias msg-commit", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "zero", "nameservice-setAliasStatus", "", "", "", "", "fee-auth"}, "", nil},
		{"nameservice", true, true, "Non nameservice authorizer account try to authorize", "gohck", "0cin", 0, NameServiceInfo{"approve", "mxw-alias", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", true, true, "all account different than nameservice provider,issuer,auth account try to sign", "mostafa", "0cin", 0, NameServiceInfo{"approve", "mxw-alias", "", "", " ", "nago", "0", "acc-40", "true"}, "", nil},

		{"nameservice", true, true, "Approve name mxw-alias without authorizer", "nago", "0cin", 0, NameServiceInfo{"approve", "mxw-alias-2", "", "", "", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", false, false, "Approve name mxw-alias", "ns-auth", "0cin", 0, NameServiceInfo{"approve", "mxw-alias", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", false, true, "Approve name mxw-alias-repeated", "ns-auth", "0cin", 0, NameServiceInfo{"approve", "mxw-alias", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},

		//NameService
		//Revoke
		{"nameservice", true, true, "Non Authorizer revoke name mxw-alias", "acc-40", "0cin", 0, NameServiceInfo{"revoke", "mxw-alias", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", true, true, "All  different accounts provider,issuer,auth in nameservice try to sign and revoke", "acc-40", "0cin", 0, NameServiceInfo{"revoke", "mxw-alias", "", "", " ", "nago", "0", "acc-40", "true"}, "", nil},
		{"nameservice", true, true, "Authorizer revoke name mxw-alias without issuer", "ns-auth", "0cin", 0, NameServiceInfo{"revoke", "mxw-alias-1", "", "0", "", "ns-provider", "0", "acc-40", "true"}, "", nil},
		{"nameservice", false, false, "Authorizer revoke  name mxw-alias", "ns-auth", "0cin", 0, NameServiceInfo{"revoke", "mxw-alias", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", false, true, "Authorizer revoke  name mxw-alias-again", "ns-auth", "0cin", 0, NameServiceInfo{"revoke", "mxw-alias", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},

		//NameService
		//Reject
		{"nameservice", true, true, "Non Authorizer reject name mxw-alias", "acc-40", "0cin", 0, NameServiceInfo{"reject", "mxw-alias", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", true, true, "all account different than nameservice provider,issuer,auth account try to sign", "acc-40", "0cin", 0, NameServiceInfo{"reject", "mxw-alias", "", "", " ", "nago", "0", "acc-40", "true"}, "", nil},
		{"nameservice", false, false, "Authorizer reject name mxw-alias-1", "ns-auth", "0cin", 0, NameServiceInfo{"reject", "mxw-alias-1", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", false, false, "Authorizer reject name mxw-alias-2", "ns-auth", "0cin", 0, NameServiceInfo{"reject", "mxw-alias-2", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},
		{"nameservice", false, true, "Authorizer reject name mxw-alias-2-again", "ns-auth", "0cin", 0, NameServiceInfo{"reject", "mxw-alias-2", "", "", " ", "ns-provider", "0", "ns-issuer", "true"}, "", nil},

		//---------------------------------------------------------------------------------------------------------------------
		// fungible token -----------------------------------------------------------------------------------------------------
		//---------------------------------------------------------------------------------------------------------------------
		//set fee fore msgTokenMultiplier to fee 0cin
		{"fee", false, false, "assign msgTokenMultiplier to fee 0cin. commit", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "zero", "fee-updateTokenMultiplier", "", "", "", "", "fee-auth"}, "", nil},

		//add token fee multiplier
		{"fee", false, false, "create token fee multiplier. commit", "fee-auth", "0cin", 0, feeInfo{"token-fee-multiplier", "", "", "1", "", "", "", "fee-auth"}, "", nil},

		// add nameservice fee collector with maintenance. (mostafa is whitelisted.)
		{"maintenance", false, false, "34. Proposal, add token fee collector address, Happy path. commit", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add token fee collector", "Add mostafa as nameservice fee collector", "fee", "", "", "", FeeCollector{Module: "token", Address: "mostafa"}, "maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as nameservice fee collector, Happy path. commit", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 34}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as nameservice fee collector, Happy path. commit", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 34}, "", nil},

		// Create Fungible Token
		{"token", false, false, "Create token - Happy Path", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken", "TT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Not fee collector", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "acc-20", "TestToken", "TT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Creating symbol that has been created", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken", "TT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Create token - Generous person pay for system fees", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken", "ToT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Poor people trying to create enormous application fee token", "acc-40", "100000000cin", 0, TokenInfo{"create", "1000000000000000000000000000", "mostafa", "TestToken", "ToT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		//{"token", true, true, "Create token - Some stranger(not whitelisted) trying to create token", "josephin", "0cin", 0, TokenInfo{"create", "1000000", "mostafa", "TestToken", "TooT", 8, "josephin","", "", "", true, "100000", false, false, false, "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Some smart people steal people signature, but owner is not signer", "acc-40", "100000000cin", 0, TokenInfo{"create", "1000000000000000000000000000", "mostafa", "TestToken", "ToT", 8, "josephin", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Stingy people trying to escape application fee", "acc-40", "100000000cin", 0, TokenInfo{"create", "0", "mostafa", "TestToken", "ToT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Wrong gas", "acc-40", "100000000cin", 1, TokenInfo{"create", "10000000", "mostafa", "TestToken", "TTT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Some smart people trying negative amount hoping negative*negative=positive", "acc-40", "100000000cin", 0, TokenInfo{"create", "-10000000", "mostafa", "TestToken", "TTT", 8, "acc-40", "", "", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Create token - Some naggy people have the world longest metadata", "acc-40", "100000000cin", 0, TokenInfo{"create", "10000000", "mostafa", "TestToken", "TTT", 8, "acc-40", "", "abcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcbacbabcabc", true, "100000", false, false, false, "", "", "", "", "", ""}, "", nil},
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

		//set MsgVerifyTokenTransferOwnership fee to 0cin
		{"fee", false, false, "assign msgVerifyTokenTransferOwnership to fee 0cin. commit", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "zero", "token-verifyTransferTokenOwnership", "", "", "", "", "fee-auth"}, "", nil},

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
		{"token", false, true, "Mint token - Not allow to mint if was Fixed Supply", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "", "TT-2", 8, "acc-40", "carlo", "", true, "100", true, false, true, "100000", "", "", "", "", ""}, "", nil},

		{"token", false, false, "Transfer fixed supply token - Happy path", "acc-40", "100000000cin", 0, TokenInfo{"transfer", "", "", "", "TT-2", 8, "acc-40", "eve", "", false, "", true, false, false, "0", "", "", "", "", ""}, "", nil},

		{"token", true, true, "Burn token - Token not existed", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-3", 0, "acc-40", "carlo", "", true, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Burn token - Invalid owner", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2", 0, "nago", "", "", true, "", true, false, false, "1", "", "", "", "", ""}, "", nil},
		{"token", false, true, "Burn token - Not enough balance for TT-2", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2", 0, "acc-40", "", "", true, "", true, false, true, "100000000", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Burn token - Fixed Supply Happy path", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2", 0, "acc-40", "", "", true, "", true, false, false, "200", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Burn token - Allow to burn same TT-2 again", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2", 0, "acc-40", "", "", true, "", true, false, true, "100", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Burn token - Not allow to burn if was approved and burnable", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2", 0, "nago", "", "", true, "", true, false, true, "1", "", "", "", "", ""}, "", nil},
		{"token", false, false, "Freeze token - Happy path for TT-2", "token-auth-1", "0cin", 0, TokenInfo{"freeze", "", "", "", "TT-2", 0, "", "", "", true, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, false, "Burn token - Not allow to burn if was approved and frozen", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2", 0, "nago", "", "", true, "", true, true, false, "1", "", "", "", "", ""}, "", nil},

		// Create-Approval-Burn : using Dynamic Supply
		{"token", true, true, "Create token - Use dynamic supply with long memo", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-2b", "TT-2b", 8, "acc-40", "", "", false, "100000", false, false, true, "", "", "", "", "", ""}, "long memo:123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890", nil},
		{"token", false, false, "Create token - Use dynamic supply with memo", "acc-40", "100000000cin", 0, TokenInfo{"create", "100000", "mostafa", "TestToken-2b", "TT-2b", 8, "acc-40", "", "", false, "100000", false, false, true, "", "", "", "", "", ""}, "1234567890", nil},
		{"token", false, false, "Approve token - Happy path for TT-2b", "token-auth-1", "0cin", 0, TokenInfo{"approve", "", "", "", "TT-2b", 0, "", "", "", false, "", false, false, true, "100000", "token-prov-1", "0", "token-issuer-1", "zero", ""}, "", nil},
		{"token", false, false, "Mint token - Allow to mint if was Dynamic Supply", "acc-40", "100000000cin", 0, TokenInfo{"mint", "", "", "", "TT-2b", 8, "acc-40", "carlo", "", false, "0", true, false, true, "100000", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Burn token - Token not existed", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-3b", 0, "acc-40", "carlo", "", false, "", true, false, true, "200", "", "", "", "", ""}, "", nil},
		{"token", true, true, "Burn token - Invalid owner", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2b", 0, "nago", "", "", false, "", true, false, true, "200", "", "", "", "", ""}, "", nil},
		{"token", false, true, "Burn token - Not enough balance for TT-2b", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2b", 0, "acc-40", "", "", false, "", true, false, true, "100000000", "", "", "", "", ""}, "", nil},

		{"token", false, true, "Burn token - Dynamic Supply Happy path", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2b", 0, "acc-40", "", "", false, "", true, false, true, "200", "", "", "", "", ""}, "", nil},      //KIV
		{"token", false, true, "Burn token - Allow to burn same TT-2b again", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2b", 0, "acc-40", "", "", false, "", true, false, true, "300", "", "", "", "", ""}, "", nil}, //KIV
		{"token", false, false, "Freeze token - Happy path for TT-2b", "token-auth-1", "0cin", 0, TokenInfo{"freeze", "", "", "", "TT-2b", 0, "", "", "", false, "", true, false, true, "1", "token-prov-1", "0", "token-issuer-1", "", ""}, "", nil},
		{"token", true, false, "Burn token - Not allow to burn if was approved and frozen", "acc-40", "100000000cin", 0, TokenInfo{"burn", "", "", "", "TT-2b", 0, "nago", "", "", false, "", true, false, true, "1", "", "", "", "", ""}, "", nil},
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

		//---------------------------------------------------------------------------------------------------------------------
		// non fungible token -------------------------------------------------------------------------------------------------
		//---------------------------------------------------------------------------------------------------------------------
		// add mostafa as nonfungible authorised address
		{"maintenance", false, false, "35. Proposal, add nonfungible authorised address, Happy path. commit", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add non fungible authorised address", "Add mostafa as non fungible authorised address", "nonFungible", "mostafa", "", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as non fungible token authorised address, Happy path. commit", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 35}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as non fungible token authorised address, Happy path. commit", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 35}, "", nil},

		// add carlo as nonfungible issuer address
		{"maintenance", false, false, "36. Proposal, add nonfungible issuer address, Happy path. commit", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add non fungible fee issuer address", "Add carlo as non fungible issuer address", "nonFungible", "", "carlo", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve carlo as non fungible token issuer address, Happy path. commit", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 36}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve carlo as non fungible token issuer address, Happy path. commit", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 36}, "", nil},

		// add jeansoon as nonfungible provider address
		{"maintenance", false, false, "37. Proposal, add nonfungible provider address, Happy path. commit", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add non fungible fee provider address", "Add jeansoon as non fungible provider address", "nonFungible", "", "", "jeansoon", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "Cast action to approve jeansoon as non fungible token provider address, Happy path. commit", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 37}, "", nil},
		{"maintenance-cast-action", false, false, "Cast action to approve jeansoon as non fungible token provider address, Happy path. commit", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 37}, "", nil},

		//set MsgTransferNonFungibleTokenOwnership fee to 0cin
		{"fee", false, false, "assign msgVerifyTokenTransferOwnership to fee 0cin. commit", "fee-auth", "0cin", 0, feeInfo{"assign-msg", "zero", "nonFungible-transferNonFungibleTokenOwnership", "", "", "", "", "fee-auth"}, "", nil},

		//====================================== start : nft modules
		// 1. create non fungible - without ItemID
		{"nonFungibleToken", false, false, "Create non fungible token - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT", "acc-40", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Happy Path.", nil},
		{"nonFungibleToken", true, true, "Create non fungible token - Token already exists (TNFT). commit", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT", "acc-29", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Token already exists (TFT).", nil},
		{"nonFungibleToken", true, true, "Create non fungible token - Insufficient fee amount. commit", "acc-29", "0cin", 0, NonFungibleTokenInfo{"create", "0", "mostafa", "TestNonFungibleToken", "TNFT", "acc-29", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Insufficient fee amount.", nil},
		{"nonFungibleToken", true, true, "Create non fungible token - Very long metadata!", "acc-29", "0cin", 0, NonFungibleTokenInfo{"create", "0", "mostafa", "TestNonFungibleToken", "TNFT", "acc-29", "", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Very long metadata!", nil},
		{"nonFungibleToken", true, true, "goh-Create non fungible token - Fee collector invalid", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "yk", "TestNonFungibleToken-191", "TNFT-191", "acc-29", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "goh-Create non fungible token - Fee collector invalid", nil},                                                                         // ok
		{"nonFungibleToken", true, true, "goh-Create non fungible token - Invalid fee amount", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"create", "abcXXX", "mostafa", "TestNonFungibleToken-191", "TNFT-191", "acc-29", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "goh-Create non fungible token - Invalid fee amount", nil},                                                                            // ok
		{"nonFungibleToken", true, true, "goh-Create non fungible token - Insufficient balance to pay for application fee", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"create", "77999999999999900000000", "mostafa", "TestNonFungibleToken-191", "TNFT-191", "acc-29", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "goh-Create non fungible token - Insufficient balance to pay for application fee", nil}, //ok

		// 2. approve non fungible - without ItemID
		{"nonFungibleToken", false, false, "Approve non fungible token(TNFT) : TransferLimit(2) Mintlimit(2) Endorser(jeanson,carlo) - Happy path", "mostafa", "0cin", 0, NonFungibleTokenInfo{"approve", "", "", "TestNonFungibleToken", "TNFT", "", "", "", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "0", "carlo", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil},
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT) : The Signer Not authorised to approve", "yk", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "mostafa", "TestNonFungibleToken", "TNFT", "", "", "", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "0", "carlo", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil},                    //ok
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT-191) : The Token symbol does not exist", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "mostafa", "TestNonFungibleToken-191", "TNFT-191", "", "", "", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "0", "carlo", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil},        //ok
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT) : Unauthorized signature - yk", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "mostafa", "TestNonFungibleToken", "TNFT", "", "", "", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "0", "yk", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil},                           //ok
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT) : Fee setting is not valid - fee-setting-191", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "mostafa", "TestNonFungibleToken", "TNFT", "", "", "", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "0", "carlo", "fee-setting-191", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT) : Token already approved - TNFT", "mostafa", "0cin", 0, NonFungibleTokenInfo{"approve", "", "", "TestNonFungibleToken", "TNFT", "", "", "", []byte(""), []string{"properties"}, []string{"metadata"}, true, false, false, "jeansoon", "0", "carlo", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil},                                              // ok

		// 3. mint non fungible - with ItemID
		{"nonFungibleToken", false, false, "Mint non fungible token - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "acc-40", "mostafa", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                    //ok
		{"nonFungibleToken", false, false, "Mint non fungible token - (mint for burn)Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "acc-40", "yk", "", []byte("223344"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},          //ok
		{"nonFungibleToken", false, false, "Mint non fungible token - (mint for endorsement)Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "acc-40", "nago", "", []byte("334455"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Mint non fungible token - Invalid Token Symbol", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT-191", "acc-40", "mostafa", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},        //ok
		{"nonFungibleToken", false, true, "Mint non fungible token - Token item id is in used.", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "acc-40", "mostafa", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},      //ok
		{"nonFungibleToken", false, true, "Mint non fungible token - Invalid token minter.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "yk", "mostafa", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                  //ok

		//====================== with ItemID :
		// --------------------OK [tc.signer, i.Owner, i.Symbol, i.ItemID]
		// 4. make endorsement - with ItemID
		{"nonFungibleToken", false, false, "endorse a nonfungible item - Happy path", "carlo", "100000000cin", 0, NonFungibleTokenInfo{"endorsement", "", "", "", "TNFT", "carlo", "", "", []byte("334455"), []string{""}, []string{""}, true, false, false, "", "", "", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},
		// start : goh-nft
		{"nonFungibleToken", false, true, "endorse a nonfungible item - Invalid endorser", "yk", "100000000cin", 0, NonFungibleTokenInfo{"endorsement", "", "", "", "TNFT", "yk", "", "", []byte("778899"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},               // ok
		{"nonFungibleToken", false, true, "endorse a nonfungible item - Invalid Token Symbol", "carlo", "100000000cin", 0, NonFungibleTokenInfo{"endorsement", "", "", "", "TNFT-111", "carlo", "", "", []byte("334455"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil}, //ok
		{"nonFungibleToken", false, true, "endorse a nonfungible item - Invalid Item-ID", "carlo", "100000000cin", 0, NonFungibleTokenInfo{"endorsement", "", "", "", "TNFT", "carlo", "", "", []byte("999111"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},          // ok

		// --------------------OK [i.Owner, i.NewOwner, i.Symbol, i.ItemID]
		// 5. transfer non fungible item - with ItemID
		{"nonFungibleToken", false, false, "Transfer non fungible token item - Happy path", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"transfer", "", "", "", "TNFT", "mostafa", "yk", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		// start : goh-nftmostafa
		{"nonFungibleToken", true, true, "Transfer non fungible token item - Invalid Token Symbol", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"transfer", "", "", "", "TNFT-111", "mostafa", "yk", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},         // ok
		{"nonFungibleToken", true, true, "Transfer non fungible token item - Invalid Account to transfer from", "carlo", "100000000cin", 0, NonFungibleTokenInfo{"transfer", "", "", "", "TNFT-111", "carlo", "yk", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, // ok
		{"nonFungibleToken", false, true, "Transfer non fungible token item - Invalid Item-ID", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"transfer", "", "", "", "TNFT", "mostafa", "yk", "", []byte("999111"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                 // ok

		// --------------------OK [i.Owner, i.Symbol, i.ItemID]
		// 6. burn non fungible item - with ItemID
		{"nonFungibleToken", false, true, "Burn non fungible token item - Invalid token owner", "carlo", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT", "carlo", "", "", []byte("223344"), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", false, false, "Burn non fungible token item -  Happy path", "yk", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT", "yk", "", "", []byte("223344"), []string{""}, []string{""}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		// start : goh-nft
		{"nonFungibleToken", true, true, "Burn non fungible token item - Invalid Token Symbol", "yk", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT-111", "yk", "", "", []byte("223344"), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                             //ok
		{"nonFungibleToken", false, true, "Burn non fungible token item - Invalid Item-ID", "yk", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT", "yk", "", "", []byte("999111"), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                                     //ok
		{"nonFungibleToken", true, true, "Burn non fungible token item - Invalid account to burn from due to yet pass KYC", "yk", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT", "acc-19", "", "", []byte("223344"), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok

		// --------------------OK
		//=============================== start goh : base on 'TNFT-191'
		// create non fungible :
		{"nonFungibleToken", false, false, "Create non fungible token - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken-191", "TNFT-191", "acc-40", "", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Happy Path.", nil}, // ok
		// mint non fungible :
		{"nonFungibleToken", true, true, "Mint non fungible token item - Invalid token as yet to approved", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT-191", "acc-40", "mostafa", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		// burn non fungible :
		{"nonFungibleToken", true, true, "Burn non fungible token item - Invalid token", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "'TNFT-191'", "acc-40", "mostafa", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		//=============================== end goh : base on 'TNFT-191'

		//====================== without ItemID :
		// --------------------OK
		// 7. transfer ownership - without ItemID
		// 1. validation.go : [case nonFungible.MsgTransferNonFungibleTokenOwnership:]
		{"nonFungibleToken", false, false, "Transfer non fungible token ownership - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "", "", "", "TNFT", "acc-40", "yk", "", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		// start : goh-nft
		// 7.1. must 'CheckApprovedToken'
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-T1 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-T1", "acc-40", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Transfer non fungible token ownership - Invalid token as yet to approved", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-T1", "acc-40", "yk", "", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},            //ok
		// 7.2. must 'IsTokenOwner'
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-T2 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-T2", "acc-40", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Transfer non fungible token ownership - Invalid token owner", "yk", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-T2", "yk", "acc-40", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                     //ok

		// --------------------OK
		// 8. verify transfer token ownership - without ItemID
		{"nonFungibleToken", false, false, "Approve non fungible token transfer ownership - Happy path for TNFT", "mostafa", "0cin", 0, NonFungibleTokenInfo{"verify-transfer-tokenOwnership", "", "", "", "TNFT", "mostafa", "yk", "", []byte(""), []string{""}, []string{""}, true, false, false, "jeansoon", "0", "carlo", "", "APPROVE_TRANFER_TOKEN_OWNERSHIP", "", "", []string{""}}, "", nil},

		// --------------------OK
		// 9. accept token ownership - without ItemID
		// 1. validation.go : [case nonFungible.MsgAcceptNonFungibleTokenOwnership:]
		{"nonFungibleToken", false, false, "Accept non fungible token ownership - Happy path. commit", "yk", "100000000cin", 0, NonFungibleTokenInfo{"accept-ownership", "", "", "", "TNFT", "", "yk", "", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		// start : goh-nft
		// 9.1.1 must 'CheckApprovedToken'
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-Q1 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-Q1", "acc-40", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Accept non fungible token ownership - Invalid token as yet to approved", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-Q1", "acc-40", "yk", "", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},              //ok
		// 9.2. must 'IsTokenNewOwner'
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-Q2 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-Q2", "acc-40", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Accept non fungible token ownership - Invalid token new-owner", "yk", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-Q2", "yk", "acc-40", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                   //ok
		// 9.3. must 'IsTokenOwnershipTransferrable'
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-Q3 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-Q3", "acc-40", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                       //ok
		{"nonFungibleToken", true, true, "Accept non fungible token ownership - Invalid token due to IsTokenOwnershipTransferrable == FALSE", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-Q3", "acc-40", "yk", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok

		// --------------------OK
		// 10. freeze non fungible token - without ItemID
		{"nonFungibleToken", false, false, "Freeze non fungible token - Happy path. commit", "mostafa", "0cin", 0, NonFungibleTokenInfo{"freeze", "", "", "", "TNFT", "", "", "", []byte(""), []string{""}, []string{""}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", true, true, "Transfer non fungible token ownership - Invalid token action (due to Token not approved) ", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "", "", "", "TNFT", "acc-40", "yk", "", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		// start : goh-nft
		{"nonFungibleToken", false, false, "Create non fungible token - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken-192", "TNFT-192", "acc-40", "", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Happy Path.", nil}, // ok
		{"nonFungibleToken", true, true, "Freeze non fungible token - Not authorised to approve due to Invalid Fee collector", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"freeze", "10000000", "yk", "TestNonFungibleToken-192", "TNFT-192", "", "", "", []byte(""), []string{""}, []string{""}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{""}}, "", nil},             //ok
		{"nonFungibleToken", true, true, "Freeze non fungible token - Invalid Token symbol - TNFT-111", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"freeze", "10000000", "mostafa", "TestNonFungibleToken-111", "TNFT-111", "", "", "", []byte(""), []string{""}, []string{""}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{""}}, "", nil},                              //ok

		// --------------------OK
		// 11. unfreeze non fungible token - without ItemID
		{"nonFungibleToken", false, false, "Unfreeze non fungible token - Happy path", "mostafa", "0cin", 0, NonFungibleTokenInfo{"unfreeze", "", "", "", "TNFT", "", "", "", []byte(""), []string{""}, []string{""}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{""}}, "", nil},
		// start : goh-nft
		{"nonFungibleToken", true, true, "Unfreeze non fungible token - Not authorised to approve due to Invalid Fee collector", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze", "10000000", "yk", "TestNonFungibleToken-192", "TNFT-192", "", "", "", []byte(""), []string{""}, []string{""}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{""}}, "", nil}, // ok
		{"nonFungibleToken", true, true, "Unfreeze non fungible token - Invalid Token symbol - TNFT-111", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze", "10000000", "mostafa", "TestNonFungibleToken-111", "TNFT-111", "", "", "", []byte(""), []string{""}, []string{""}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{""}}, "", nil},                  // ok

		//====================== without ItemID :
		// freeze and THEN unfreeze
		{"nonFungibleToken", false, false, "CREATE non fungible token [TNFT-B2] - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-B2", "acc-40", "", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},
		{"nonFungibleToken", false, false, "APPROVE non fungible token [TNFT-B2] - Happy path.  commit", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-B2", "acc-40", "", "metadata", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "0", "carlo", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil},
		{"nonFungibleToken", false, false, "MINT non fungible item [TNFT-B2] - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-B2", "acc-40", "mostafa", "metadata", []byte("001177"), []string{"properties"}, []string{"metadata"}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil}, //ok
		{"nonFungibleToken", false, false, "FREEZE non fungible item [TNFT-B2] - Happy path.", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-B2", "mostafa", "", "metadata", []byte("001177"), []string{"properties"}, []string{"metadata"}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},
		{"nonFungibleToken", false, false, "UNFREEZE non fungible item [TNFT-B2] - Happy path.", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "mostafa", "", "TNFT-B2", "mostafa", "", "metadata", []byte("001177"), []string{"properties"}, []string{"metadata"}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},

		// --------------------OK [tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.ItemID]
		// 6.1. freeze non fungible item - with ItemID
		{"nonFungibleToken", false, false, "CREATE non fungible token [TNFT-D2] - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-D2", "acc-40", "", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},
		{"nonFungibleToken", false, false, "APPROVE non fungible token [TNFT-D2] - Happy path.  commit", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-D2", "acc-40", "", "metadata", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "0", "carlo", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil},
		{"nonFungibleToken", false, false, "MINT non fungible item [TNFT-D2] - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-D2", "acc-40", "mostafa", "metadata", []byte("880099"), []string{"properties"}, []string{"metadata"}, true, true, false, "jeansoon", "0", "carlo", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil}, //ok
		{"nonFungibleToken", false, false, "Freeze non fungible item [TNFT-D2] - Happy path.", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "mostafa", "", "", []byte("880099"), []string{""}, []string{""}, true, true, false, "jeansoon", "0", "carlo", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", false, true, "Freeze non fungible item [TNFT-D2] - Invalid signer.", "jeansoon", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "jeansoon", "", "", []byte("880099"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                //ok
		{"nonFungibleToken", false, true, "Freeze non fungible item [TNFT-D2] - No such non fungible token.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-9988", "yk", "", "", []byte("880099"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},              //ok
		{"nonFungibleToken", false, true, "Freeze non fungible item [TNFT-D2] - No such item to freeze.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "yk", "", "", []byte("991111"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                    //ok
		{"nonFungibleToken", false, true, "Freeze non fungible item [TNFT-D2] - Not authorised to freeze non token item.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "yk", "", "", []byte("880099"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},   //ok
		{"nonFungibleToken", false, true, "Freeze non fungible item [TNFT-D2] - Non Fungible item already frozen.", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "mostafa", "", "", []byte("880099"), []string{""}, []string{""}, true, true, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil}, //??? --- why need nonce==1
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-D2] - Invalid nonce.", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-D2", "mostafa", "", "", []byte("880099"), []string{""}, []string{""}, true, true, false, "jeansoon", "2", "carlo", "", "", "", "", []string{""}}, "", nil},                //???--- why need nonce==1
		{"nonFungibleToken", false, false, "Unfreeze non fungible item [TNFT-D2] - Happy path.", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-D2", "mostafa", "", "", []byte("880099"), []string{""}, []string{""}, true, true, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                  //???--- why need nonce==1

		// --------------------OK [tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.ItemID]
		// 6.2. unfreeze non fungible item - with ItemID
		{"nonFungibleToken", false, false, "CREATE non fungible token [TNFT-E2] - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-E2", "acc-40", "", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},
		{"nonFungibleToken", false, false, "APPROVE non fungible token [TNFT-E2] - Happy path.  commit", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-E2", "acc-40", "", "metadata", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "1", "carlo", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil},
		{"nonFungibleToken", false, false, "MINT non fungible item [TNFT-E2] - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-E2", "acc-40", "mostafa", "metadata", []byte("770099"), []string{"properties"}, []string{"metadata"}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil}, //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - Invalid signer.", "jeansoon", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-E2", "jeansoon", "", "", []byte("770099"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                                                         //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - No such non fungible token.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-9988", "yk", "", "", []byte("770099"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                                                       //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - No such  non fungible item to unfreeze.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-E2", "yk", "", "", []byte("991111"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                                             //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - Not authorised to unfreeze token account.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-E2", "yk", "", "", []byte("770099"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                                           //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - Non fungible item not frozen.", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-E2", "mostafa", "", "", []byte("770099"), []string{""}, []string{""}, true, true, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                                              //???--- why need nonce==1

	}
	var totalFee = sdkTypes.NewInt64Coin("cin", 0)
	var totalAmt = sdkTypes.NewInt64Coin("cin", 0)

	for n, tc := range tcs {

		fees, err := types.ParseCoins(tc.fees)
		assert.NoError(t, err, tc.fees)

		var msg sdkTypes.Msg
		switch tc.msgType {
		case "bank":
			{
				i := tc.msgInfo.(bankInfo)
				amt, err := types.ParseCoins(i.amount)
				assert.NoError(t, err, i.amount)

				msg = bank.NewMsgSend(tKeys[i.from].addr, tKeys[i.to].addr, amt)

				if !tc.deliverFailed {
					if i.from == "alice" {
						totalFee = totalFee.Add(fees[0])
						totalAmt = totalAmt.Add(amt[0])
					}
				}
			}
		case "kyc":
			{
				i := tc.msgInfo.(kycInfo)

				switch i.action {
				case "whitelist":
					msg = makeKycWhitelistMsg(t, i.authorised, i.issuer, i.provider, i.from, i.signer, i.data, i.nonce)

				case "revokeWhitelist":
					msg = makeKycRevokeWhitelistMsg(t, i.authorised, i.issuer, i.provider, i.from)
				}
			}

		case "fee":
			{
				i := tc.msgInfo.(feeInfo)

				msg = makeFeeMsg(t, i.function, i.name, i.assignee, i.multiplier, i.min, i.max, i.percentage, i.issuer)
			}

		case "nameservice":
			{
				i := tc.msgInfo.(NameServiceInfo)

				switch i.Action {
				case "create":
					msg = makeCreateNameServiceMsg(t, i.Name, i.From, i.ApplicationFee, i.FeeCollector)
				case "approve":
					msg = setStatusAliasMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Name, "APPROVE")
				case "reject":
					msg = setStatusAliasMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Name, "REJECT")
				case "revoke":
					msg = setStatusAliasMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Name, "REVOKE")
				}

			}
		case "token":
			{
				i := tc.msgInfo.(TokenInfo)

				switch i.Action {
				case "create":
					msg = makeCreateFungibleTokenMsg(t, i.Name, i.Symbol, i.Metadata, i.Owner, i.MaxSupply, i.ApplicationFee, i.FeeCollector, i.Decimals, i.FixedSupply)
				case "approve":
					msg = makeApproveFungibleTokenMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, "APPROVE", i.Burnable, i.FeeSettingName)
				case "transfer":
					msg = makeTransferFungibleTokenMsg(t, i.Owner, i.NewOwner, i.Symbol, i.AmountOfToken)
				case "mint":
					msg = makeMintFungibleTokenMsg(t, i.Owner, i.NewOwner, i.Symbol, i.AmountOfToken)
				case "burn":
					msg = makeBurnFungibleTokenMsg(t, i.Owner, i.Symbol, i.AmountOfToken)
				case "transfer-ownership":
					msg = makeTransferFungibleTokenOwnershipMsg(t, i.Owner, i.NewOwner, i.Symbol)
				case "accept-ownership":
					msg = makeAcceptFungibleTokenOwnershipMsg(t, i.NewOwner, i.Symbol)
				case "freeze":
					msg = makeFreezeFungibleTokenMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.Burnable)
				case "unfreeze":
					msg = makeUnfreezeFungibleTokenMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.Burnable)
				case "verify-transfer-tokenOwnership":
					msg = makeVerifyTransferTokenOwnership(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.VerifyTransferTokenOwnership)
				}
			}
		case "nonFungibleToken":
			{
				i := tc.msgInfo.(NonFungibleTokenInfo)

				switch i.Action {
				case "create":
					msg = makeCreateNonFungibleTokenMsg(t, i.Name, i.Symbol, i.TokenMetadata, i.Owner, i.ApplicationFee, i.FeeCollector)
				case "approve":
					msg = makeApproveNonFungibleTokenMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, "APPROVE", i.FeeSettingName, i.MintLimit, i.TransferLimit, i.EndorserList)
				case "transfer":
					msg = makeTransferNonFungibleTokenMsg(t, i.Owner, i.NewOwner, i.Symbol, i.ItemID)
				case "mint":
					msg = makeMintNonFungibleTokenMsg(t, i.Owner, i.NewOwner, i.Symbol, i.ItemID, i.Properties, i.Metadata)
				case "burn":
					msg = makeBurnNonFungibleTokenMsg(t, i.Owner, i.Symbol, i.ItemID)
				case "transfer-ownership":
					msg = makeTransferNonFungibleTokenOwnershipMsg(t, i.Owner, i.NewOwner, i.Symbol)
				case "accept-ownership":
					msg = makeAcceptNonFungibleTokenOwnershipMsg(t, i.NewOwner, i.Symbol)
				case "freeze-item":
					msg = makeFreezeNonFungibleItemMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.ItemID)
				case "unfreeze-item":
					msg = makeUnfreezeNonFungibleItemMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.ItemID)
				case "freeze":
					msg = makeFreezeNonFungibleTokenMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.Burnable)
				case "unfreeze":
					msg = makeUnfreezeNonFungibleTokenMsg(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.Burnable)
				case "verify-transfer-tokenOwnership":
					msg = makeVerifyTransferNonFungibleTokenOwnership(t, tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, i.VerifyTransferTokenOwnership)
				case "endorsement":
					msg = makeEndorsement(t, tc.signer, i.Owner, i.Symbol, i.ItemID)
				}
			}
		case "maintenance":
			{
				i := tc.msgInfo.(MaintenanceInfo)
				msg = makeMaintenanceMsg(t, i.Action, i.Title, i.Description, i.ProposalType, i.AuthorisedAddress, i.IssuerAddress, i.ProviderAddress, i.Proposer, i.ValidatorPubKey, i.FeeCollector)
			}

		case "maintenance-cast-action":
			{
				i := tc.msgInfo.(CastAction)
				msg = makeCastActionMsg(t, i.Action, i.Caster, i.ProposalId)
			}
		}

		tx, bz := MakeSignedTx(t, tc.signer, tc.gas, fees, tc.memo, msg)
		txb, err := tCdc.MarshalJSON(tx)
		assert.NoError(t, err)
		fmt.Printf("============\nBroadcasting %d tx (%s): %s\n", n+1, tc.desc, string(txb))

		res := BroadcastTxCommit(bz)
		tc.hash = res.Hash.Bytes()

		if !tc.checkFailed {
			tKeys[tc.signer].seq++
			require.Zero(t, res.CheckTx.Code, "test case %v(%v) check failed: %v", n+1, tc.desc, res.CheckTx.Log)

			if !tc.deliverFailed {
				require.Zero(t, res.DeliverTx.Code, "test case %v(%v) deliver failed: %v", n+1, tc.desc, res.DeliverTx.Log)
			} else {
				require.NotZero(t, res.DeliverTx.Code, "test case %v(%v) deliver failed: %v", n+1, tc.desc, res.DeliverTx.Log)
			}

		} else {
			require.NotZero(t, res.CheckTx.Code, "test case %v(%v) check failed: %v", n+1, tc.desc, res.CheckTx.Log)
		}

		if strings.Contains(tc.desc, "commit") {
			WaitForNextHeightTM(tPort)
		}
	}

	WaitForNextHeightTM(tPort)

	/// check results
	for i, tc := range tcs {
		res := Tx(tc.hash)
		if tc.checkFailed {
			assert.Nil(t, res, "test case %v(%v) failed", i, tc.desc)
		} else {
			if tc.deliverFailed {
				assert.NotZero(t, res.TxResult.Code, "test case %v(%v) failed", i, tc.desc)
			} else {
				assert.Zero(t, res.TxResult.Code, "test case %v(%v) failed: %v", i, tc.desc, res.TxResult.Log)
			}
		}
	}

	//acc2 := Account(tKeys["alice"].addrStr)
	//bal2 := acc2.GetCoins()[0]
	//diff := bal1.Sub(bal2)

	//total := totalAmt.Add(totalFee)
	//require.Equal(t, diff, total)
	val2 := Validator(tValidator)
	fmt.Println(val2)

	// accGohck := Account(tKeys["gohck"].addrStr)
	// require.Empty(t, accGohck.GetCoins())
}

func MakeSignedTx(t *testing.T, name string, gas uint64, fees sdkTypes.Coins, memo string, msg sdkTypes.Msg) (sdkAuth.StdTx, []byte) {
	acc := Account(tKeys[name].addrStr)
	require.NotNil(t, acc)

	signMsg := authTypes.StdSignMsg{
		AccountNumber: acc.GetAccountNumber(),
		ChainID:       "maxonrow-chain",
		Fee:           authTypes.NewStdFee(gas, fees),
		Memo:          memo,
		Msgs:          []sdkTypes.Msg{msg},
		Sequence:      tKeys[name].seq,
	}

	signBz, signBzErr := tCdc.MarshalJSON(signMsg)
	if signBzErr != nil {
		panic(signBzErr)
	}

	sig, err := tKeys[name].priv.Sign(sdkTypes.MustSortJSON(signBz))
	if err != nil {
		panic(err)
	}

	pub := tKeys[name].priv.PubKey()
	stdSig := sdkAuth.StdSignature{
		PubKey:    pub,
		Signature: sig,
	}

	sdtTx := authTypes.NewStdTx(signMsg.Msgs, signMsg.Fee, []authTypes.StdSignature{stdSig}, signMsg.Memo)
	bz, err := tCdc.MarshalBinaryLengthPrefixed(sdtTx)
	if err != nil {
		panic(err)
	}
	return sdtTx, bz
}

func makeKycWhitelistMsg(t *testing.T, authorised, issuer, provider, from, signer, data, nonce string) sdkTypes.Msg {
	// create new kyc data to be whitelisted
	kycData := kyc.NewKyc(tKeys[from].addr, nonce, data)

	// kyc signed by the address which want to be whitelisted
	kycDataBz, err := tCdc.MarshalJSON(kycData)
	require.NoError(t, err)
	signedKycDataBz, err := tKeys[signer].priv.Sign(sdkTypes.MustSortJSON(kycDataBz))
	require.NoError(t, err)

	// creating the kyc payload
	kycPayload := kyc.NewPayload(kycData, tKeys[from].pub, signedKycDataBz)

	// kycPayload to be signed by issuer and provider
	kycPayloadBz, err := tCdc.MarshalJSON(kycPayload)
	require.NoError(t, err)
	var signatures []kyc.Signature

	if issuer != "nope" {
		issuerSignedBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(kycPayloadBz))
		require.NoError(t, err)
		issuerSignature := kyc.NewSignature(tKeys[issuer].pub, issuerSignedBz)
		signatures = append(signatures, issuerSignature)
	}

	if provider != "nope" {
		providerSignedBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(kycPayloadBz))
		require.NoError(t, err)
		providerSignature := kyc.NewSignature(tKeys[provider].pub, providerSignedBz)
		signatures = append(signatures, providerSignature)
	}
	toWhitelistData := kyc.NewKycData(kycPayload, signatures)

	return kyc.NewMsgWhitelist(tKeys[authorised].addr, toWhitelistData)
}

func makeKycRevokeWhitelistMsg(t *testing.T, authorised, issuer, provider, to string) sdkTypes.Msg {
	// convert uint64 to string
	providerNonceStr := strconv.FormatUint(tKeys[provider].seq, 10)

	// create new kyc data to be revoked
	revokeKycData := kyc.NewRevokeKycData(tKeys[provider].addr, providerNonceStr, tKeys[to].addr)

	var signedRevokeKycDataBz []byte

	if provider != "nope" {
		// revokeKycData signed by the provider
		revokeKycDataBz, err := tCdc.MarshalJSON(revokeKycData)
		require.NoError(t, err)
		signedRevokeKycDataBz, err = tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(revokeKycDataBz))
		require.NoError(t, err)
	}
	// creating the revoke kyc payload
	revokeKycPayload := kyc.NewRevokePayload(revokeKycData, tKeys[provider].pub, signedRevokeKycDataBz)

	var signatures []kyc.Signature

	if issuer != "nope" {
		// revokeKycPayload to be signed by issuer
		revokeKycPayloadBz, err := tCdc.MarshalJSON(revokeKycPayload)
		require.NoError(t, err)
		issuerSignedBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(revokeKycPayloadBz))
		require.NoError(t, err)
		issuerSignature := kyc.NewSignature(tKeys[issuer].pub, issuerSignedBz)
		signatures = append(signatures, issuerSignature)
	}
	return kyc.NewMsgRevokeWhitelist(sdkTypes.AccAddress(tKeys[authorised].addr), revokeKycPayload, signatures)
}

func makeFeeMsg(t *testing.T, function, name, assignee, multiplier, _min, _max, percentage, issuer string) sdkTypes.Msg {
	var msg sdkTypes.Msg
	switch function {
	case "sys-fee":
		min, err := sdkTypes.ParseCoins(_min)
		require.NoError(t, err)
		max, err := sdkTypes.ParseCoins(_max)
		require.NoError(t, err)

		msg = fee.NewMsgSysFeeSetting(name, min, max, percentage, tKeys[issuer].addr)
	case "assign-msg":
		msg = fee.NewMsgAssignFeeToMsg(name, assignee, tKeys[issuer].addr)
	case "assign-acc":
		msg = fee.NewMsgAssignFeeToAcc(name, tKeys[assignee].addr, tKeys[issuer].addr)
	case "multiplier":
		msg = fee.NewMsgMultiplier(multiplier, tKeys[issuer].addr)
	case "token-fee-multiplier":
		msg = fee.NewMsgTokenMultiplier(multiplier, tKeys[issuer].addr)
	}
	return msg
}

func makeCreateNameServiceMsg(t *testing.T, name, owner, applicationFee, aliasFeeCollector string) sdkTypes.Msg {

	// create new alias
	ownerAddr := tKeys[owner].addr
	fee := nameservice.Fee{
		To:    tKeys[aliasFeeCollector].addr,
		Value: applicationFee,
	}
	msgCreateAlias := nameservice.NewMsgCreateAlias(name, ownerAddr, fee)

	return msgCreateAlias
}

func setStatusAliasMsg(t *testing.T, signer, provider, providerNonce, issuer, name, status string) sdkTypes.Msg {

	providerAddr := tKeys[provider].addr

	nameserviceDoc := nameservice.NewAlias(providerAddr, providerNonce, status, name)

	// provider sign the nameservice
	nsProvider, err := tCdc.MarshalJSON(nameserviceDoc)
	require.NoError(t, err)
	signedAlias, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(nsProvider))
	require.NoError(t, err)

	nameservicePayload := nameservice.NewPayload(*nameserviceDoc, tKeys[provider].pub, signedAlias)

	// issuer sign the nameservice
	aliasPayload, err := tCdc.MarshalJSON(nameservicePayload)
	require.NoError(t, err)
	signedAliasPayload, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(aliasPayload))
	require.NoError(t, err)

	var signatures []nameservice.Signature
	signature := nameservice.NewSignature(tKeys[issuer].pub, signedAliasPayload)
	signatures = append(signatures, signature)

	msgSetFungibleTokenStatus := nameservice.NewMsgSetAliasStatus(tKeys[signer].addr, *nameservicePayload, signatures)

	return msgSetFungibleTokenStatus
}
