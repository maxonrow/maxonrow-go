package tests

type bankInfo struct {
	from   string
	to     string
	amount string
}

func makeBankTxs() []*testCase {

	tcs := []*testCase{

		{"bank", false, false, "sending 1 cin", "alice", "200000000cin", 0, bankInfo{"alice", "bob", "1cin"}, "", nil},
		{"bank", true, true, "sending 0 cin", "alice", "200000000cin", 0, bankInfo{"alice", "bob", "0cin"}, "", nil},
		{"bank", true, true, "insufficient amount", "carlo", "1000000000cin", 0, bankInfo{"carlo", "eve", "99999999999999999999999000000001cin"}, "", nil},
		                                                                                                                                       
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
	}

	return tcs
}
