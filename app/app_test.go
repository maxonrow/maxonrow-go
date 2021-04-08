package app

import (
	"os"
	"strconv"
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	sdkAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/maxonrow/maxonrow-go/utils"
	"github.com/maxonrow/maxonrow-go/x/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

/*
How to sign messages offline:

1- Import mnemonics:
mxwcli keys import-mnemonic acc-1 "since height latin shiver gallery cage sure face twelve already leisure shop super garden maid else summer search half robot bicycle game life disease"
mxwcli keys import-mnemonic acc-2 "youth grief wagon trim fix next hammer differ minimum grit stuff actress swap episode outdoor trophy seat hero floor word wink comfort outer nasty"
mxwcli keys import-mnemonic acc-3 "shove when pass black expose blouse dial glue original wonder move glad rice guide trophy dish beach legal animal kitchen maze concert ahead keep"

2- create a tx file for internal transacion: internal_tx.json

3- Sign internal transaction:
mxwcli tx sign internal_tx.json --from acc-1 --offline --account-number 1 --sequence 0 --chain-id mxw
mxwcli tx sign internal_tx.json --from acc-2 --offline --account-number 2 --sequence 0 --chain-id mxw

4- create file for tx: tx.json

5- Sign transaction:
mxwcli tx sign tx.json --from acc-1 --offline --account-number 1 --sequence 0 --chain-id mxw

6- Change the internal tx signature and sign with account 2
mxwcli tx sign tx.json --from acc-2 --offline --account-number 2 --sequence 0 --chain-id mxw
*/

var addr1 sdkTypes.AccAddress
var addr2 sdkTypes.AccAddress
var addr3 sdkTypes.AccAddress
var addr4 sdkTypes.AccAddress
var gaddr sdkTypes.AccAddress

func defaultLogger() log.Logger {
	return log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
}

func newMxwApp(t *testing.T) (*mxwApp, sdkTypes.Context) {
	logger := defaultLogger()
	db := dbm.NewMemDB()
	app := NewMXWApp(logger, db)

	ctx := app.NewContext(true, abci.Header{})
	ctx = ctx.WithChainID("mxw")
	ctx = ctx.WithBlockHeight(1)

	params := sdkAuthTypes.DefaultParams()
	app.accountKeeper.SetParams(ctx, params)

	coins, _ := sdkTypes.ParseCoins("1000000cin")
	addr1 = utils.MustGetAccAddressFromBech32("mxw1ld3stcsk5l8xjngw2ucuazux895rk2hxve69gr")
	addr2 = utils.MustGetAccAddressFromBech32("mxw1cntgdz7z8lm26xmzxkn2jcr9djx8rnhdf3dl7c")
	addr3 = utils.MustGetAccAddressFromBech32("mxw1gh9hlcdevh94rrxj7gd3wtngjrhj50gkt9txvk")
	addr4 = utils.MustGetAccAddressFromBech32("mxw1858220s2tacfhzgjrxrgyyn5huxsf4r9s9tq5j")
	gaddr = utils.MustGetAccAddressFromBech32("mxw1q6nmfejarl5e4xzceqxcygner7a6llgwnrdtl6")

	acc1 := sdkAuth.BaseAccount{
		Address:       addr1,
		Coins:         coins,
		AccountNumber: 1,
		Sequence:      0,
	}
	acc2 := sdkAuth.BaseAccount{
		Address:       addr2,
		Coins:         coins,
		AccountNumber: 2,
		Sequence:      0,
	}
	acc3 := sdkAuth.BaseAccount{
		Address:       addr3,
		Coins:         coins,
		AccountNumber: 3,
		Sequence:      0,
	}
	acc4 := sdkAuth.BaseAccount{
		Address:       addr4,
		Coins:         coins,
		AccountNumber: 4,
		Sequence:      0,
	}
	gaAcc := sdkAuth.BaseAccount{
		Address:       gaddr,
		Coins:         coins,
		AccountNumber: 5,
		Sequence:      0,
		MultiSig:      sdkAuthTypes.NewMultiSig(addr1, 2, []sdkTypes.AccAddress{addr1, addr2, addr3}),
	}
	app.accountKeeper.SetAccount(ctx, &acc1)
	app.accountKeeper.SetAccount(ctx, &acc2)
	app.accountKeeper.SetAccount(ctx, &acc3)
	app.accountKeeper.SetAccount(ctx, &acc4)
	app.accountKeeper.SetAccount(ctx, &gaAcc)
	app.kycKeeper.Whitelist(ctx, addr1, "acc-addr-1")
	app.kycKeeper.Whitelist(ctx, addr2, "acc-addr-2")
	app.kycKeeper.Whitelist(ctx, addr3, "acc-addr-3")
	app.kycKeeper.Whitelist(ctx, gaddr, "group-addr-1")

	mnemonic1 := "since height latin shiver gallery cage sure face twelve already leisure shop super garden maid else summer search half robot bicycle game life disease"
	mnemonic2 := "youth grief wagon trim fix next hammer differ minimum grit stuff actress swap episode outdoor trophy seat hero floor word wink comfort outer nasty"
	mnemonic3 := "shove when pass black expose blouse dial glue original wonder move glad rice guide trophy dish beach legal animal kitchen maze concert ahead keep"
	mnemonic4 := "fire milk legal six result shoulder cake globe quote absorb beauty glass ski crash tilt suspect paddle speed gather tunnel project wife fatal abstract"

	proc, err := utils.CreateProcess("", "mxwcli", []string{"keys", "import-mnemonic", "acc-1", mnemonic1, "--keyring-backend", "os"})
	err = proc.Cmd.Run()
	require.NoError(t, err)

	proc, err = utils.CreateProcess("", "mxwcli", []string{"keys", "import-mnemonic", "acc-2", mnemonic2, "--keyring-backend", "os"})
	err = proc.Cmd.Run()
	require.NoError(t, err)

	proc, err = utils.CreateProcess("", "mxwcli", []string{"keys", "import-mnemonic", "acc-3", mnemonic3, "--keyring-backend", "os"})
	err = proc.Cmd.Run()
	require.NoError(t, err)

	proc, err = utils.CreateProcess("", "mxwcli", []string{"keys", "import-mnemonic", "acc-4", mnemonic4, "--keyring-backend", "os"})
	err = proc.Cmd.Run()
	require.NoError(t, err)

	return app, ctx
}

// accNum is group account account number which is not same as signer account number
func signInternalTx(t *testing.T, app *mxwApp, from string, accNum, txID int) sdkAuth.StdTx {
	bz := `{
		"type": "cosmos-sdk/StdTx",
		"value": {
			"msg": [
				{
					"type": "mxw/msgSend",
					"value": {
						"from_address": "mxw1q6nmfejarl5e4xzceqxcygner7a6llgwnrdtl6",
						"to_address": "mxw193cra94tl3uj9zcmx2x2j7vrc3y38z39469dqy",
						"amount": [
							{
								"denom": "cin",
								"amount": "1"
							}
						]
					}
				}
			],
			"fee": {
				"amount": [
					{
						"denom": "cin",
						"amount": "800400000"
					}
				],
				"gas": "0"
			},
			"memo": ""
		}
	}`

	file, err := os.Create("/tmp/internal_tx.json")
	assert.NoError(t, err)
	defer file.Close()

	_, err = file.WriteString(bz)
	assert.NoError(t, err)

	proc, err := utils.CreateProcess("", "mxwcli", []string{"tx", "sign", "/tmp/internal_tx.json", "--from", from, "--offline", "--account-number", strconv.Itoa(accNum), "--sequence", strconv.Itoa(txID), "--chain-id", "mxw", "--keyring-backend", "os"})
	err = proc.Cmd.Start()
	require.NoError(t, err)

	out, _, _ := proc.ReadAll()
	//fmt.Printf("%s%s", string(out))
	proc.Cmd.Wait()

	var tx sdkAuth.StdTx
	err = app.cdc.UnmarshalJSON(out, &tx)
	assert.NoError(t, err)

	return tx
}

func signMultisigTx(t *testing.T, app *mxwApp, gaddr, sender sdkTypes.AccAddress, itx sdkAuth.StdTx, from string, seq int) sdkAuth.StdTx {

	var msg auth.MsgCreateMultiSigTx
	msg.Sender = sender
	msg.GroupAddress = gaddr
	msg.StdTx = itx

	tx1 := sdkAuth.NewStdTx([]sdkTypes.Msg{msg}, itx.Fee, nil, "")

	bz, err := app.cdc.MarshalJSON(&tx1)
	assert.NoError(t, err)

	file, err := os.Create("/tmp/tx.json")
	assert.NoError(t, err)
	defer file.Close()

	_, err = file.Write(bz)
	assert.NoError(t, err)

	proc, err := utils.CreateProcess("", "mxwcli", []string{"tx", "sign", "/tmp/tx.json", "--from", from, "--offline", "--account-number", from[4:], "--sequence", strconv.Itoa(seq), "--chain-id", "mxw", "--keyring-backend", "os"})
	err = proc.Cmd.Start()
	require.NoError(t, err)

	out, _, _ := proc.ReadAll()
	//fmt.Printf("%s%s", string(out))
	proc.Cmd.Wait()

	var tx2 sdkAuth.StdTx
	err = app.cdc.UnmarshalJSON(out, &tx2)
	assert.NoError(t, err)

	return tx2
}

func TestMultisig_1(t *testing.T) {
	app, ctx := newMxwApp(t)

	// Signer is acc-1, sender is acc-1
	itx := signInternalTx(t, app, "acc-1", 5, 0)
	tx := signMultisigTx(t, app, gaddr, addr1, itx, "acc-1", 0)
	//bz, _ := app.cdc.MarshalJSON(&tx)
	//fmt.Println(string(bz))

	_, err := utils.CheckTxSig(ctx, tx, app.accountKeeper, app.kycKeeper)
	assert.NoError(t, err)
}

func TestMultisig_2(t *testing.T) {
	app, ctx := newMxwApp(t)

	// Signer is acc-1, sender is acc-2
	itx := signInternalTx(t, app, "acc-1", 5, 0)
	tx := signMultisigTx(t, app, gaddr, addr2, itx, "acc-1", 0)
	_, err := utils.CheckTxSig(ctx, tx, app.accountKeeper, app.kycKeeper)
	assert.Error(t, err) // sender is different than signer

	// Signer is acc-4, sender is acc-4, no kyc
	itx = signInternalTx(t, app, "acc-1", 5, 0)
	tx = signMultisigTx(t, app, gaddr, addr4, itx, "acc-4", 0)

	_, err = utils.CheckTxSig(ctx, tx, app.accountKeeper, app.kycKeeper)
	assert.Error(t, err)

	// Signer is acc-2, sender is acc-2
	itx = signInternalTx(t, app, "acc-1", 5, 0)
	tx = signMultisigTx(t, app, gaddr, addr2, itx, "acc-2", 0)

	_, err = utils.CheckTxSig(ctx, tx, app.accountKeeper, app.kycKeeper)
	assert.NoError(t, err)
}

func TestMultisig_3(t *testing.T) {
	app, ctx := newMxwApp(t)

	pub1 := sdkTypes.MustGetAccPubKeyBech32("mxwpub1addwnpepq24jdanvwecsda8l7ldr9xdgr374yms25g9ztrsxmwn4y0pvszxxcw83mrv")
	pub2 := sdkTypes.MustGetAccPubKeyBech32("mxwpub1addwnpepqwnf4wxmkgmzn8m8ntqf6d7tvhasf6svtmhrfwf4vh64za08wv8a5nknyus")
	pub3 := sdkTypes.MustGetAccPubKeyBech32("mxwpub1addwnpepqgqzznmhvjykz0kyt3aeds6xmpsyu4qeznueu8y2xw3hw0ykluepvr0e52u")
	pub4 := sdkTypes.MustGetAccPubKeyBech32("mxwpub1addwnpepqdpmalfk8jmg3ryptfdacpvekzf8rsqddvpmyfwtnanr93ehav7tczp9r0n")

	acc1 := app.accountKeeper.GetAccount(ctx, addr1)
	acc1.SetPubKey(pub1)
	app.accountKeeper.SetAccount(ctx, acc1)

	acc2 := app.accountKeeper.GetAccount(ctx, addr2)
	acc2.SetPubKey(pub2)
	app.accountKeeper.SetAccount(ctx, acc2)

	acc3 := app.accountKeeper.GetAccount(ctx, addr3)
	acc3.SetPubKey(pub3)
	app.accountKeeper.SetAccount(ctx, acc3)

	acc4 := app.accountKeeper.GetAccount(ctx, addr4)
	acc4.SetPubKey(pub4)
	app.accountKeeper.SetAccount(ctx, acc4)

	{
		// itx signed by acc-4
		itx := signInternalTx(t, app, "acc-4", 5, 0)
		_, err := utils.CheckTxSig(ctx, itx, app.accountKeeper, app.kycKeeper)
		assert.Error(t, err)
	}

	{
		itx1 := signInternalTx(t, app, "acc-1", 5, 0)
		_, err := utils.CheckTxSig(ctx, itx1, app.accountKeeper, app.kycKeeper)
		assert.NoError(t, err)

		// Remove Public Key
		itx1.Signatures[0].PubKey = nil
		_, err = utils.CheckTxSig(ctx, itx1, app.accountKeeper, app.kycKeeper)
		assert.NoError(t, err)

		itx2 := signInternalTx(t, app, "acc-2", 5, 0)
		_, err = utils.CheckTxSig(ctx, itx2, app.accountKeeper, app.kycKeeper)
		assert.NoError(t, err)

		// Remove Public Keys
		itx2.Signatures[0].PubKey = nil
		_, err = utils.CheckTxSig(ctx, itx2, app.accountKeeper, app.kycKeeper)
		assert.NoError(t, err)

		itx2.Signatures = append(itx2.Signatures, itx1.Signatures...)

		_, err = utils.CheckTxSig(ctx, itx2, app.accountKeeper, app.kycKeeper)
		assert.NoError(t, err)
	}
}
