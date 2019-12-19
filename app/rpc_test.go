package app

import (
	"fmt"
	"math/big"
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tm-db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/maxonrow/maxonrow-go/x/bank"
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
