package tests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/maxonrow/maxonrow-go/types"

	"github.com/maxonrow/maxonrow-go/x/bank"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	nonFungible "github.com/maxonrow/maxonrow-go/x/token/nonfungible"
)

type NonFungibleTokenInfo struct {
	Action                       string
	ApplicationFee               string
	FeeCollector                 string
	Name                         string
	Symbol                       string
	Owner                        string
	NewOwner                     string
	TokenMetadata                string
	ItemID                       []byte
	Properties                   []string
	Metadata                     []string
	Approved                     bool
	Frozen                       bool
	Burnable                     bool
	Provider                     string
	ProviderNonce                string
	Issuer                       string
	FeeSettingName               string
	VerifyTransferTokenOwnership string
	TransferLimit                string
	MintLimit                    string
	EndorserList                 []string
}

func TestNFT(t *testing.T) {

	val1 := Validator(tValidator)
	fmt.Println(val1)

	_, err := time.ParseDuration("60s")
	if err != nil {
		fmt.Println("timeout", err)
	}

	tcs := []*testCase{

		//YK
		{"kyc", false, false, "Doing kyc - yk - commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "yk", "yk", "testKyc12345678", "0"}, "", nil},
		{"bank", false, false, "sending after updating fee", "yk", "800000000cin", 0, bankInfo{"yk", "mostafa", "1000cin"}, "", nil},

		//mostafa (role As Fee-collector)
		{"bank", false, false, "sending after updating fee", "yk", "800000000cin", 0, bankInfo{"yk", "mostafa", "1000cin"}, "", nil},
		{"kyc", false, false, "Doing kyc - mostafa - commit", "kyc-auth-1", "0cin", 0, kycInfo{"kyc-auth-1", "kyc-issuer-1", "kyc-prov-1", "whitelist", "mostafa", "mostafa", "testKyc123456789", "0"}, "", nil},
		{"fee", false, false, "assign zero-fee to mostafa-commit", "fee-auth", "0cin", 0, feeInfo{"assign-acc", "zero", "mostafa", "", "", "", "", "fee-auth"}, "", nil},
		{"bank", false, false, "sending after updating acc-fee", "mostafa", "0cin", 0, bankInfo{"mostafa", "bob", "1cin"}, "", nil},

		//add mostafa as nonfungible authorised address
		{"maintenance", false, false, "1. Proposal, add nonfungible authorised address, Happy path. commit", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add non fungible authorised address", "Add mostafa as non fungible authorised address", "nonFungible", "mostafa", "", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as non fungible token authorised address, Happy path. commit", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 1}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as non fungible token authorised address, Happy path. commit", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 1}, "", nil},

		//add carlo as nonfungible issuer address
		{"maintenance", false, false, "2. Proposal, add nonfungible issuer address, Happy path. commit", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add non fungible fee issuer address", "Add carlo as non fungible issuer address", "nonFungible", "", "carlo", "", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve carlo as non fungible token issuer address, Happy path. commit", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 2}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve carlo as non fungible token issuer address, Happy path. commit", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 2}, "", nil},

		//add jeansoon as nonfungible provider address
		{"maintenance", false, false, "3. Proposal, add nonfungible provider address, Happy path. commit", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add non fungible fee provider address", "Add jeansoon as non fungible provider address", "nonFungible", "", "", "jeansoon", FeeCollector{}, "maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "Cast action to approve jeansoon as non fungible token provider address, Happy path. commit", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 3}, "", nil},
		{"maintenance-cast-action", false, false, "Cast action to approve jeansoon as non fungible token provider address, Happy path. commit", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 3}, "", nil},

		//add nameservice fee collector with maintenance. (mostafa is whitelisted.)
		{"maintenance", false, false, "4. Proposal, add token fee collector address, Happy path. commit", "maintainer-2", "0cin", 0, MaintenanceInfo{"add", "Add token fee collector", "Add mostafa as nameservice fee collector", "fee", "", "", "", FeeCollector{Module: "token", Address: "mostafa"}, "maintainer-2", ""}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as nameservice fee collector, Happy path. commit", "maintainer-1", "0cin", 0, CastAction{"maintainer-1", "approve", 4}, "", nil},
		{"maintenance-cast-action", false, false, "(Approve)-Cast action to approve mostafa as nameservice fee collector, Happy path. commit", "maintainer-3", "0cin", 0, CastAction{"maintainer-3", "approve", 4}, "", nil},

		// 1. create non fungible - without ItemID
		{"nonFungibleToken", false, false, "Create non fungible token - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT", "acc-40", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Happy Path.", nil},
		{"nonFungibleToken", true, true, "Create non fungible token - Token already exists (TNFT). commit", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT", "acc-29", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Token already exists (TFT).", nil},
		{"nonFungibleToken", true, true, "Create non fungible token - Insufficient fee amount. commit", "acc-29", "0cin", 0, NonFungibleTokenInfo{"create", "0", "mostafa", "TestNonFungibleToken", "TNFT", "acc-29", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Insufficient fee amount.", nil},
		{"nonFungibleToken", true, true, "Create non fungible token - Very long metadata!", "acc-29", "0cin", 0, NonFungibleTokenInfo{"create", "0", "mostafa", "TestNonFungibleToken", "TNFT", "acc-29", "", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Very long metadata!", nil},
		{"nonFungibleToken", true, true, "goh-Create non fungible token - Fee collector invalid", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "yk", "TestNonFungibleToken-191", "TNFT-191", "acc-29", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "goh-Create non fungible token - Fee collector invalid", nil},                                                                         // ok
		{"nonFungibleToken", true, true, "goh-Create non fungible token - Invalid fee amount", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"create", "abcXXX", "mostafa", "TestNonFungibleToken-191", "TNFT-191", "acc-29", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "goh-Create non fungible token - Invalid fee amount", nil},                                                                            // ok
		{"nonFungibleToken", true, true, "goh-Create non fungible token - Insufficient balance to pay for application fee", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"create", "77999999999999900000000", "mostafa", "TestNonFungibleToken-191", "TNFT-191", "acc-29", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "goh-Create non fungible token - Insufficient balance to pay for application fee", nil}, //ok

		// 2. approve non fungible - without ItemID [tc.signer, i.Provider, i.ProviderNonce, i.Issuer, i.Symbol, "APPROVE", i.FeeSettingName, i.MintLimit, i.TransferLimit, i.EndorserList]
		{"nonFungibleToken", false, false, "Approve non fungible token(TNFT) : TransferLimit(2) Mintlimit(2) Endorser(jeanson,carlo) - Happy path", "mostafa", "0cin", 0, NonFungibleTokenInfo{"approve", "", "", "TestNonFungibleToken", "TNFT", "", "", "", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "0", "carlo", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil},
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT) : The Signer Not authorised to approve", "yk", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "mostafa", "TestNonFungibleToken", "TNFT", "", "", "", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "0", "carlo", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil},                    //ok
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT-191) : The Token symbol does not exist", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "mostafa", "TestNonFungibleToken-191", "TNFT-191", "", "", "", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "0", "carlo", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil},        //ok
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT) : Unauthorized signature - yk", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "mostafa", "TestNonFungibleToken", "TNFT", "", "", "", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "0", "yk", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil},                           //ok
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT) : Fee setting is not valid - fee-setting-191", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "mostafa", "TestNonFungibleToken", "TNFT", "", "", "", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "0", "carlo", "fee-setting-191", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT) : Token already approved - TNFT", "mostafa", "0cin", 0, NonFungibleTokenInfo{"approve", "", "", "TestNonFungibleToken", "TNFT", "", "", "", []byte(""), []string{"properties"}, []string{"metadata"}, true, false, false, "jeansoon", "0", "carlo", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil},                                              // ok

		// 3. mint non fungible - with ItemID [i.Owner, i.NewOwner, i.Symbol, i.ItemID, i.Properties, i.Metadata]
		{"nonFungibleToken", false, false, "Mint non fungible token - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "acc-40", "mostafa", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                    //ok
		{"nonFungibleToken", false, false, "Mint non fungible token - (mint for burn)Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "acc-40", "yk", "", []byte("223344"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},          //ok
		{"nonFungibleToken", false, false, "Mint non fungible token - (mint for endorsement)Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "acc-40", "nago", "", []byte("334455"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Mint non fungible token - Invalid Token Symbol", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT-191", "acc-40", "mostafa", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},        //ok
		{"nonFungibleToken", false, true, "Mint non fungible token - Token item id is in used.", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "acc-40", "mostafa", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},      //ok
		{"nonFungibleToken", false, true, "Mint non fungible token - Invalid token minter.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "yk", "mostafa", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                  //ok

		//====================== with ItemID :
		// 4. make endorsement - with ItemID
		{"nonFungibleToken", false, false, "endorse a nonfungible item - Happy path", "carlo", "100000000cin", 0, NonFungibleTokenInfo{"endorsement", "", "", "", "TNFT", "carlo", "", "", []byte("334455"), []string{""}, []string{""}, true, false, false, "", "", "", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},
		{"nonFungibleToken", false, true, "endorse a nonfungible item - Invalid endorser", "yk", "100000000cin", 0, NonFungibleTokenInfo{"endorsement", "", "", "", "TNFT", "yk", "", "", []byte("778899"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},               // ok
		{"nonFungibleToken", false, true, "endorse a nonfungible item - Invalid Token Symbol", "carlo", "100000000cin", 0, NonFungibleTokenInfo{"endorsement", "", "", "", "TNFT-111", "carlo", "", "", []byte("334455"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil}, //ok
		{"nonFungibleToken", false, true, "endorse a nonfungible item - Invalid Item-ID", "carlo", "100000000cin", 0, NonFungibleTokenInfo{"endorsement", "", "", "", "TNFT", "carlo", "", "", []byte("999111"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},          // ok

		// 5. transfer non fungible item - with ItemID
		{"nonFungibleToken", false, false, "Transfer non fungible token item - Happy path", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"transfer", "", "", "", "TNFT", "mostafa", "yk", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", true, true, "Transfer non fungible token item - Invalid Token Symbol", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"transfer", "", "", "", "TNFT-111", "mostafa", "yk", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},         // ok
		{"nonFungibleToken", true, true, "Transfer non fungible token item - Invalid Account to transfer from", "carlo", "100000000cin", 0, NonFungibleTokenInfo{"transfer", "", "", "", "TNFT-111", "carlo", "yk", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, // ok
		{"nonFungibleToken", false, true, "Transfer non fungible token item - Invalid Item-ID", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"transfer", "", "", "", "TNFT", "mostafa", "yk", "", []byte("999111"), []string{"properties"}, []string{"metadata"}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                 // ok

		// 6. burn non fungible item - with ItemID
		{"nonFungibleToken", false, true, "Burn non fungible token item - Invalid token owner", "carlo", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT", "carlo", "", "", []byte("223344"), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", false, false, "Burn non fungible token item -  Happy path", "yk", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT", "yk", "", "", []byte("223344"), []string{""}, []string{""}, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", true, true, "Burn non fungible token item - Invalid Token Symbol", "yk", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT-111", "yk", "", "", []byte("223344"), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                             //ok
		{"nonFungibleToken", false, true, "Burn non fungible token item - Invalid Item-ID", "yk", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT", "yk", "", "", []byte("999111"), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                                     //ok
		{"nonFungibleToken", true, true, "Burn non fungible token item - Invalid account to burn from due to yet pass KYC", "yk", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT", "acc-19", "", "", []byte("223344"), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok

		//=============================== start goh : base on 'TNFT-191'
		// create non fungible :
		{"nonFungibleToken", false, false, "Create non fungible token - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken-191", "TNFT-191", "acc-40", "", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Happy Path.", nil}, // ok
		// mint non fungible :
		{"nonFungibleToken", true, true, "Mint non fungible token item - Invalid token as yet to approved", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT-191", "acc-40", "mostafa", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		// burn non fungible :
		{"nonFungibleToken", true, true, "Burn non fungible token item - Invalid token", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "'TNFT-191'", "acc-40", "mostafa", "", []byte("112233"), []string{"properties"}, []string{"metadata"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		//=============================== end goh : base on 'TNFT-191'

		//====================== without ItemID :
		// 7. transfer ownership - without ItemID
		{"nonFungibleToken", false, false, "Transfer non fungible token ownership - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "", "", "", "TNFT", "acc-40", "yk", "", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-T1 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-T1", "acc-40", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Transfer non fungible token ownership - Invalid token as yet to approved", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-T1", "acc-40", "yk", "", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},            //ok
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-T2 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-T2", "acc-40", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Transfer non fungible token ownership - Invalid token owner", "yk", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-T2", "yk", "acc-40", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                     //ok

		// 8. verify transfer token ownership - without ItemID
		{"nonFungibleToken", false, false, "Approve non fungible token transfer ownership - Happy path for TNFT", "mostafa", "0cin", 0, NonFungibleTokenInfo{"verify-transfer-tokenOwnership", "", "", "", "TNFT", "mostafa", "yk", "", []byte(""), []string{""}, []string{""}, true, false, false, "jeansoon", "0", "carlo", "", "APPROVE_TRANFER_TOKEN_OWNERSHIP", "", "", []string{""}}, "", nil},

		// 9. accept token ownership - without ItemID
		{"nonFungibleToken", false, false, "Accept non fungible token ownership - Happy path. commit", "yk", "100000000cin", 0, NonFungibleTokenInfo{"accept-ownership", "", "", "", "TNFT", "", "yk", "", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-Q1 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-Q1", "acc-40", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                       //ok
		{"nonFungibleToken", true, true, "Accept non fungible token ownership - Invalid token as yet to approved", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-Q1", "acc-40", "yk", "", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                                    //ok
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-Q2 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-Q2", "acc-40", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                       //ok
		{"nonFungibleToken", true, true, "Accept non fungible token ownership - Invalid token new-owner", "yk", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-Q2", "yk", "acc-40", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                                         //ok
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-Q3 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-Q3", "acc-40", "", "metadata", []byte("test test test"), []string{"hi", "bye"}, []string{"hi", "bye"}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                       //ok
		{"nonFungibleToken", true, true, "Accept non fungible token ownership - Invalid token due to IsTokenOwnershipTransferrable == FALSE", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-Q3", "acc-40", "yk", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok

		// 10. freeze non fungible token - without ItemID
		{"nonFungibleToken", false, false, "Freeze non fungible token - Happy path. commit", "mostafa", "0cin", 0, NonFungibleTokenInfo{"freeze", "", "", "", "TNFT", "", "", "", []byte(""), []string{""}, []string{""}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", true, true, "Transfer non fungible token ownership - Invalid token action (due to Token not approved) ", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "", "", "", "TNFT", "acc-40", "yk", "", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", false, false, "Create non fungible token - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken-192", "TNFT-192", "acc-40", "", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Happy Path.", nil}, // ok
		{"nonFungibleToken", true, true, "Freeze non fungible token - Not authorised to approve due to Invalid Fee collector", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"freeze", "10000000", "yk", "TestNonFungibleToken-192", "TNFT-192", "", "", "", []byte(""), []string{""}, []string{""}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{""}}, "", nil},             //ok
		{"nonFungibleToken", true, true, "Freeze non fungible token - Invalid Token symbol - TNFT-111", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"freeze", "10000000", "mostafa", "TestNonFungibleToken-111", "TNFT-111", "", "", "", []byte(""), []string{""}, []string{""}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{""}}, "", nil},                              //ok

		// 11. unfreeze non fungible token - without ItemID
		{"nonFungibleToken", false, false, "Unfreeze non fungible token - Happy path", "mostafa", "0cin", 0, NonFungibleTokenInfo{"unfreeze", "", "", "", "TNFT", "", "", "", []byte(""), []string{""}, []string{""}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", true, true, "Unfreeze non fungible token - Not authorised to approve due to Invalid Fee collector", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze", "10000000", "yk", "TestNonFungibleToken-192", "TNFT-192", "", "", "", []byte(""), []string{""}, []string{""}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{""}}, "", nil}, // ok
		{"nonFungibleToken", true, true, "Unfreeze non fungible token - Invalid Token symbol - TNFT-111", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze", "10000000", "mostafa", "TestNonFungibleToken-111", "TNFT-111", "", "", "", []byte(""), []string{""}, []string{""}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{""}}, "", nil},                  // ok

		//====================== without ItemID :
		// freeze and THEN unfreeze
		{"nonFungibleToken", false, false, "CREATE non fungible token [TNFT-B2] - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-B2", "acc-40", "", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},
		{"nonFungibleToken", false, false, "APPROVE non fungible token [TNFT-B2] - Happy path.  commit", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-B2", "acc-40", "", "metadata", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "0", "carlo", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil}, //ok
		{"nonFungibleToken", false, false, "MINT non fungible item [TNFT-B2] - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-B2", "acc-40", "mostafa", "metadata", []byte("001177"), []string{"properties"}, []string{"metadata"}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},               //ok
		{"nonFungibleToken", false, false, "FREEZE non fungible item [TNFT-B2] - Happy path.", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-B2", "mostafa", "", "metadata", []byte("001177"), []string{"properties"}, []string{"metadata"}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},                                     //ok
		{"nonFungibleToken", false, false, "UNFREEZE non fungible item [TNFT-B2] - Happy path.", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "mostafa", "", "TNFT-B2", "mostafa", "", "metadata", []byte("001177"), []string{"properties"}, []string{"metadata"}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},                          //ok

		// 12. freeze non fungible item - with ItemID
		{"nonFungibleToken", false, false, "CREATE non fungible token [TNFT-D2] - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-D2", "acc-40", "", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},                                             //ok
		{"nonFungibleToken", false, false, "APPROVE non fungible token [TNFT-D2] - Happy path.  commit", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-D2", "acc-40", "", "metadata", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "0", "carlo", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil}, //ok
		{"nonFungibleToken", false, false, "MINT non fungible item [TNFT-D2] - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-D2", "acc-40", "mostafa", "metadata", []byte("880099"), []string{"properties"}, []string{"metadata"}, true, true, false, "jeansoon", "0", "carlo", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},                //ok
		{"nonFungibleToken", false, false, "Freeze non fungible item [TNFT-D2] - Happy path.", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "mostafa", "", "", []byte("880099"), []string{""}, []string{""}, true, true, false, "jeansoon", "0", "carlo", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", false, true, "Freeze non fungible item [TNFT-D2] - Invalid signer.", "jeansoon", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "jeansoon", "", "", []byte("880099"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                //ok
		{"nonFungibleToken", false, true, "Freeze non fungible item [TNFT-D2] - No such non fungible token.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-9988", "yk", "", "", []byte("880099"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},              //ok
		{"nonFungibleToken", false, true, "Freeze non fungible item [TNFT-D2] - No such item to freeze.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "yk", "", "", []byte("991111"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                    //ok
		{"nonFungibleToken", false, true, "Freeze non fungible item [TNFT-D2] - Not authorised to freeze non token item.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "yk", "", "", []byte("880099"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},   //ok
		{"nonFungibleToken", false, true, "Freeze non fungible item [TNFT-D2] - Non Fungible item already frozen.", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "mostafa", "", "", []byte("880099"), []string{""}, []string{""}, true, true, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-D2] - Invalid nonce.", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-D2", "mostafa", "", "", []byte("880099"), []string{""}, []string{""}, true, true, false, "jeansoon", "2", "carlo", "", "", "", "", []string{""}}, "", nil},                //ok
		{"nonFungibleToken", false, false, "Unfreeze non fungible item [TNFT-D2] - Happy path.", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-D2", "mostafa", "", "", []byte("880099"), []string{""}, []string{""}, true, true, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                  //ok

		// 13. unfreeze non fungible item - with ItemID
		{"nonFungibleToken", false, false, "CREATE non fungible token [TNFT-E2] - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-E2", "acc-40", "", "metadata", []byte(""), []string{""}, []string{""}, false, false, false, "", "", "", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil},
		{"nonFungibleToken", false, false, "APPROVE non fungible token [TNFT-E2] - Happy path.  commit", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-E2", "acc-40", "", "metadata", []byte(""), []string{"properties"}, []string{"metadata"}, false, false, false, "jeansoon", "1", "carlo", "default", "", "2", "2", []string{"jeansoon", "carlo"}}, "", nil},
		{"nonFungibleToken", false, false, "MINT non fungible item [TNFT-E2] - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "10000000", "mostafa", "TestNonFungibleToken", "TNFT-E2", "acc-40", "mostafa", "metadata", []byte("770099"), []string{"properties"}, []string{"metadata"}, true, false, false, "jeansoon", "0", "carlo", "", "", "", "", []string{"jeansoon", "carlo"}}, "", nil}, //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - Invalid signer.", "jeansoon", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-E2", "jeansoon", "", "", []byte("770099"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                                                         //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - No such non fungible token.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-9988", "yk", "", "", []byte("770099"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                                                       //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - No such  non fungible item to unfreeze.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-E2", "yk", "", "", []byte("991111"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                                             //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - Not authorised to unfreeze token account.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-E2", "yk", "", "", []byte("770099"), []string{""}, []string{""}, true, false, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                                           //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - Non fungible item not frozen.", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-E2", "mostafa", "", "", []byte("770099"), []string{""}, []string{""}, true, true, false, "jeansoon", "1", "carlo", "", "", "", "", []string{""}}, "", nil},                                              //ok

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

	val2 := Validator(tValidator)
	fmt.Println(val2)

}

func makeCreateNonFungibleTokenMsg(t *testing.T, name, symbol, metadata, owner, applicationFee, tokenFeeCollector string) sdkTypes.Msg {

	// create new token
	ownerAddr := tKeys[owner].addr
	fee := nonFungible.Fee{
		To:    tKeys[tokenFeeCollector].addr,
		Value: applicationFee,
	}
	msgCreateNonFungibleToken := nonFungible.NewMsgCreateNonFungibleToken(symbol, ownerAddr, name, metadata, fee)

	return msgCreateNonFungibleToken
}

func makeApproveNonFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, status string, feeSettingName string, mintLimit, transferLimit string, endorserList []string) sdkTypes.Msg {

	providerAddr := tKeys[provider].addr

	var tokenFee = []nonFungible.TokenFee{
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

	mintL := sdkTypes.NewUintFromString(mintLimit)
	transferL := sdkTypes.NewUintFromString(transferLimit)

	var endorsers []sdkTypes.AccAddress

	for _, v := range endorserList {
		addr := tKeys[v].addr
		endorsers = append(endorsers, addr)
	}

	tokenDoc := nonFungible.NewToken(providerAddr, providerNonce, status, symbol, transferL, mintL, tokenFee, endorsers)

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := nonFungible.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetFungibleTokenStatus := nonFungible.NewMsgSetNonFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetFungibleTokenStatus
}

//module of transfer
func makeTransferNonFungibleTokenMsg(t *testing.T, owner string, newOwner string, symbol string, itemID []byte) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr

	msgTransferPayload := nonFungible.NewMsgTransferNonFungibleToken(symbol, ownerAddr, newOwnerAddr, itemID)
	return msgTransferPayload

}

//module of mint
func makeMintNonFungibleTokenMsg(t *testing.T, owner string, newOwner string, symbol string, itemID []byte, properties, metadata []string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr

	msgMintPayload := nonFungible.NewMsgMintNonFungibleToken(ownerAddr, symbol, newOwnerAddr, itemID, properties, metadata)
	return msgMintPayload

}

//module of burn
func makeBurnNonFungibleTokenMsg(t *testing.T, owner string, symbol string, itemID []byte) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr

	msgBurnNonFungible := nonFungible.NewMsgBurnNonFungibleToken(symbol, ownerAddr, itemID)
	return msgBurnNonFungible

}

//moduel of transferOwnership
func makeTransferNonFungibleTokenOwnershipMsg(t *testing.T, owner string, newOwner string, symbol string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr

	msgTransferOwnershipPayload := nonFungible.NewMsgTransferNonFungibleTokenOwnership(symbol, ownerAddr, newOwnerAddr)
	return msgTransferOwnershipPayload

}

//module of acceptOwnership
func makeAcceptNonFungibleTokenOwnershipMsg(t *testing.T, newOwner string, symbol string) sdkTypes.Msg {

	fromAddr := tKeys[newOwner].addr

	msgAcceptOwnershipPayload := nonFungible.NewMsgAcceptNonFungibleTokenOwnership(symbol, fromAddr)
	return msgAcceptOwnershipPayload

}

func makeFreezeNonFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, burnable bool) sdkTypes.Msg {

	status := "FREEZE"
	providerAddr := tKeys[provider].addr

	tokenDoc := nonFungible.NewToken(providerAddr, providerNonce, status, symbol, sdkTypes.NewUint(0), sdkTypes.NewUint(0), nil, nil) // status : FREEZE / UNFREEZE

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := nonFungible.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetNonFungibleTokenStatus := nonFungible.NewMsgSetNonFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetNonFungibleTokenStatus
}

func makeUnfreezeNonFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, burnable bool) sdkTypes.Msg {

	status := "UNFREEZE"
	providerAddr := tKeys[provider].addr

	tokenDoc := nonFungible.NewToken(providerAddr, providerNonce, status, symbol, sdkTypes.NewUint(0), sdkTypes.NewUint(0), nil, nil) // status : FREEZE / UNFREEZE

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := nonFungible.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetNonFungibleTokenStatus := nonFungible.NewMsgSetNonFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetNonFungibleTokenStatus
}

func makeVerifyTransferNonFungibleTokenOwnership(t *testing.T, signer, provider, providerNonce, issuer, symbol, action string) sdkTypes.Msg {

	providerAddr := tKeys[provider].addr

	// burnable and tokenfees is not in used for verifying transfer token status, we just set it to false and leave it empty.
	verifyTransferTokenOwnershipDoc := nonFungible.NewToken(providerAddr, providerNonce, action, symbol, sdkTypes.NewUint(0), sdkTypes.NewUint(0), []nonFungible.TokenFee{}, nil)

	// provider sign
	verifyTransferTokenOwnershipDocBz, err := tCdc.MarshalJSON(verifyTransferTokenOwnershipDoc)
	require.NoError(t, err)
	signedVerifyTransferTokenOwnershipDoc, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(verifyTransferTokenOwnershipDocBz))
	require.NoError(t, err)

	verifyTransferTokenOwnershipPayload := nonFungible.NewPayload(*verifyTransferTokenOwnershipDoc, tKeys[provider].pub, signedVerifyTransferTokenOwnershipDoc)

	// issuer sign
	verifyTransferPayloadBz, err := tCdc.MarshalJSON(verifyTransferTokenOwnershipPayload)
	require.NoError(t, err)
	signedVerifyTransferPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(verifyTransferPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedVerifyTransferPayloadBz)
	signatures = append(signatures, signature)

	msgVerifyTransferNonFungibleTokenOwnership := nonFungible.NewMsgSetNonFungibleTokenStatus(tKeys[signer].addr, *verifyTransferTokenOwnershipPayload, signatures)

	return msgVerifyTransferNonFungibleTokenOwnership
}

func makeEndorsement(t *testing.T, signer, to, symbol string, itemID []byte) sdkTypes.Msg {

	signerAddr := tKeys[signer].addr
	toAddr := tKeys[to].addr

	return nonFungible.NewMsgEndorsement(symbol, signerAddr, toAddr, itemID)
}

// Freeze Item
func makeFreezeNonFungibleItemMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, itemID []byte) sdkTypes.Msg {

	status := "FREEZE_ITEM"
	providerAddr := tKeys[provider].addr

	itemDetails := nonFungible.NewItemDetails(providerAddr, providerNonce, status, symbol, itemID) // status : FREEZE / UNFREEZE

	// provider sign the item
	itemBz, err := tCdc.MarshalJSON(itemDetails)
	require.NoError(t, err)
	signedItemBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(itemBz))
	require.NoError(t, err)

	itemPayload := nonFungible.NewItemPayload(*itemDetails, tKeys[provider].pub, signedItemBz)

	// issuer sign the itemPayload
	itemPayloadBz, err := tCdc.MarshalJSON(itemPayload)
	require.NoError(t, err)
	signedItemPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(itemPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedItemPayloadBz)
	signatures = append(signatures, signature)

	msgSetNonFungibleTokenStatus := nonFungible.NewMsgSetNonFungibleItemStatus(tKeys[signer].addr, *itemPayload, signatures)
	return msgSetNonFungibleTokenStatus
}

// UnFreeze Item
func makeUnfreezeNonFungibleItemMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, itemID []byte) sdkTypes.Msg {

	status := "UNFREEZE_ITEM"
	providerAddr := tKeys[provider].addr

	itemDetails := nonFungible.NewItemDetails(providerAddr, providerNonce, status, symbol, itemID) // status : FREEZE / UNFREEZE

	// provider sign the item
	itemBz, err := tCdc.MarshalJSON(itemDetails)
	require.NoError(t, err)
	signedItemBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(itemBz))
	require.NoError(t, err)

	itemPayload := nonFungible.NewItemPayload(*itemDetails, tKeys[provider].pub, signedItemBz)

	// issuer sign the itemPayload
	itemPayloadBz, err := tCdc.MarshalJSON(itemPayload)
	require.NoError(t, err)
	signedItemPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(itemPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedItemPayloadBz)
	signatures = append(signatures, signature)

	msgSetNonFungibleTokenStatus := nonFungible.NewMsgSetNonFungibleItemStatus(tKeys[signer].addr, *itemPayload, signatures)
	return msgSetNonFungibleTokenStatus
}
