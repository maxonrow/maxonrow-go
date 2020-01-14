package app

import (
	"fmt"
	"math/big"
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	multisig "github.com/maxonrow/maxonrow-go/x/auth"
	"github.com/maxonrow/maxonrow-go/x/bank"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func TestDecodeTx(t *testing.T) {
	app := NewMXWApp(log.NewNopLogger(), dbm.NewMemDB())
	_, _, addr1 := KeyTestPubAddr()
	_, _, addr2 := KeyTestPubAddr()

	amt := sdkTypes.Coins{
		{
			Denom:  "cin",
			Amount: sdkTypes.NewInt(1000000),
		},
	}

	fee := auth.NewStdFee(200000, amt)
	msgSend := bank.NewMsgSend(addr1, addr2, amt)
	tx1 := auth.NewStdTx([]sdkTypes.Msg{msgSend}, fee, nil, "")

	bz, _ := app.cdc.MarshalJSON(tx1)
	js1 := string(bz)

	bz, err := app.EncodeTx(nil, js1)
	assert.NoError(t, err)

	js2, err := app.DecodeTx(nil, bz)
	assert.NoError(t, err)
	assert.Equal(t, js1, js2)
	fmt.Println(js1)
}

func KeyTestPubAddr() (crypto.PrivKey, crypto.PubKey, sdkTypes.AccAddress) {

	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdkTypes.AccAddress(pub.Address())
	return key, pub, addr
}

func test() {
	verybig := big.NewInt(1)
	ten := big.NewInt(10)
	for i := 0; i < 100000; i++ {
		verybig.Mul(verybig, ten)
	}
	fmt.Println(verybig)
}

func TestDecodeMultiSig(t *testing.T) {
	jsonStr := "eyAidHlwZSI6ICJjb3Ntb3Mtc2RrL1N0ZFR4IiwgInZhbHVlIjogeyAibXNnIjogWyB7ICJ0eXBlIjogIm14dy9tc2dDcmVhdGVNdWx0aVNpZ0FjY291bnQiLCAidmFsdWUiOiB7ICJvd25lciI6ICJteHcxZHd3M253dHB2ZmNxMmg5NHJtbGZ0d3l3eTdza2M0OHlha3UyN3AiLCAidGhyZXNob2xkIjogIjIiLCAic2lnbmVycyI6IFsgIm14dzFkd3czbnd0cHZmY3EyaDk0cm1sZnR3eXd5N3NrYzQ4eWFrdTI3cCIsICJteHcxbmo1eGR6NnljaHZhMm1qcjdkbnpwMzZ0c2Z6ZWZwaGFkcTIzMG0iLCAibXh3MXFnd3pkeGY2NnRwNW1qcGtwZmU1OTNudnNzdDdxemZ4enFxNzNkIiBdIH0gfSBdLCAiZmVlIjogeyAiYW1vdW50IjogWyB7ICJhbW91bnQiOiAiMTAwMDAwMDAwIiwgImRlbm9tIjogImNpbiIgfSBdLCAiZ2FzIjogIjAiIH0sICJzaWduYXR1cmVzIjogWyB7ICJwdWJfa2V5IjogeyAidHlwZSI6ICJ0ZW5kZXJtaW50L1B1YktleVNlY3AyNTZrMSIsICJ2YWx1ZSI6ICJBdVZOY003dDNPcW8vMTFTcktvenZLL3ZUQm5YREFxL1prSVM0SWtITy8xcSIgfSwgInNpZ25hdHVyZSI6ICJqQVg2YnBla0R6U2JPdTgwMnJFd1dnVUlhLytPdFBydkxEVng5QUJITDhkdVNCcGlQVWZubHRaNnlnWFNuNmMzcUVTNFNxeXBKOGdBUEg2Sk5XeUZPZz09IiB9IF0sICJtZW1vIjogIiIgfSB9"

	cdc := MakeDefaultCodec()

	bz := parseJSON(jsonStr)
	var tx sdkAuth.StdTx
	err := cdc.UnmarshalJSON(bz, &tx)
	assert.NoError(t, err)
	msg := tx.Msgs[0].(multisig.MsgCreateMultiSigAccount)
	assert.Equal(t, msg.Threshold, 2)
	assert.NotEmpty(t, tx.GetMsgs())

}
