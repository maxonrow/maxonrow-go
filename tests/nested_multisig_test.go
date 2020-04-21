package tests

import "github.com/maxonrow/maxonrow-go/utils"

func makeNestedMultisigTxs() []*testCase {

	//
	// `mxwcli keys multisig-address "mxw1zvgat76jxp3fasgqk7dcwzkyj0f0wz5edxg7ke" 1`
	tKeys["grp-addr-a"] = &keyInfo{
		utils.MustGetAccAddressFromBech32("mxw16zcuc9kxv9tgk0l5345uz6vllstxc224ax60r3"), nil, nil, "mxw16zcuc9kxv9tgk0l5345uz6vllstxc224ax60r3",
	}

	// Nested group address
	// `mxwcli keys multisig-address "mxw16zcuc9kxv9tgk0l5345uz6vllstxc224ax60r3" 1`
	tKeys["grp-addr-n"] = &keyInfo{
		utils.MustGetAccAddressFromBech32("mxw1mq57v4rs6ryfkhnya3smeug73m75fm8kyryt8t"), nil, nil, "mxw1mq57v4rs6ryfkhnya3smeug73m75fm8kyryt8t",
	}

	tcs := []*testCase{

		{"multiSig", false, false, "Create grp-addr-a - Happy Path", "nested-acc-4", "800400000cin", 0, MultisigInfo{"create", "nested-acc-4", "", 2, []string{"multisig-acc-1", "multisig-acc-2", "nested-acc-4"}, "", 0, nil}, "MEMO : Create grp-addr-a", nil},

		{"bank", false, false, "top-up grp-addr", "multisig-acc-1", "800400000cin", 0, bankInfo{"multisig-acc-1", "grp-addr-a", "1000000000cin"}, "MEMO : top-up account", nil},
		{"bank", false, false, "top-up grp-addr", "multisig-acc-2", "800400000cin", 0, bankInfo{"multisig-acc-2", "grp-addr-a", "1000000000cin"}, "MEMO : top-up account", nil},
		{"bank", false, false, "top-up grp-addr", "multisig-acc-3", "800400000cin", 0, bankInfo{"multisig-acc-3", "grp-addr-a", "1000000000cin"}, "MEMO : top-up account", nil},

		{"multiSig", false, false, "Create grp-addr-n", "grp-addr-a:multisig-acc-1,multisig-acc-2", "800400000cin", 0, MultisigInfo{"create", "grp-addr-a", "", 2, []string{"grp-addr-a", "multisig-acc-3", "nested-acc-4"}, "", 0, nil}, "MEMO : Create grp-addr-n", nil},
		{"bank", false, false, "top-up grp-addr-n", "multisig-acc-2", "800400000cin", 0, bankInfo{"multisig-acc-2", "grp-addr-n", "10000000000cin"}, "MEMO : top-up account", nil},

		{"bank", false, false, "grp-addr-n to mostafa1", "grp-addr-n:multisig-acc-1,multisig-acc-2,multisig-acc-3", "800400000cin", 0, bankInfo{"grp-addr-n", "mostafa", "1cin"}, "Nested multisig bank", nil},
		{"bank", false, false, "grp-addr-n to mostafa2", "grp-addr-n:multisig-acc-1,multisig-acc-2,multisig-acc-3", "800400000cin", 0, bankInfo{"grp-addr-n", "mostafa", "2cin"}, "Nested multisig bank", nil},
		{"bank", false, false, "grp-addr-n to mostafa3", "grp-addr-n:nested-acc-4,multisig-acc-3", "800400000cin", 0, bankInfo{"grp-addr-n", "mostafa", "2cin"}, "Nested multisig bank", nil},
		{"bank", true, true, "grp-addr-n to mostafa4  ", "grp-addr-n:multisig-acc-1,multisig-acc-3", "800400000cin", 0, bankInfo{"grp-addr-n", "mostafa", "1cin"}, "Nested multisig bank", nil},
		{"bank", true, true, "grp-addr-n to mostafa5  ", "grp-addr-n:multisig-acc-1,multisig-acc-2,multisig-acc-3,mostafa", "800400000cin", 0, bankInfo{"grp-addr-n", "mostafa", "1cin"}, "Nested multisig bank", nil},
		{"bank", true, true, "grp-addr-n to mostafa6  ", "grp-addr-n:multisig-acc-1,multisig-acc-2,mostafa", "800400000cin", 0, bankInfo{"grp-addr-n", "mostafa", "1cin"}, "Nested multisig bank", nil},
		{"bank", true, true, "grp-addr-n to mostafa6  ", "grp-addr-a:multisig-acc-1,multisig-acc-2,multisig-acc-3", "800400000cin", 0, bankInfo{"grp-addr-n", "mostafa", "1cin"}, "Nested multisig bank", nil},
	}

	return tcs
}
