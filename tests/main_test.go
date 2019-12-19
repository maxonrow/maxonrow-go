package tests

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	cp "github.com/otiai10/copy"
	tmCrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
	"github.com/maxonrow/maxonrow-go/app"
)

type info struct {
	addr    sdkTypes.AccAddress
	priv    tmCrypto.PrivKey
	pub     tmCrypto.PubKey
	addrStr string
	seq     uint64
}

var tKeys map[string]*info
var tCdc *codec.Codec
var tClient *rpcclient.JSONRPCClient
var tPort = "26657"
var tWorkingDir string
var tValidator = "mxwvaloper1rjgjjkkjqtd676ydahysmnfsg0v4yvwfp2n965"

func startServer(done chan struct{}) *Process {
	configFldr := path.Join(tWorkingDir, "config")
	dataFldr := path.Join(tWorkingDir, "data")
	cp.Copy("./config", configFldr)
	os.Mkdir(dataFldr, 0755)
	ioutil.WriteFile(path.Join(dataFldr, "priv_validator_state.json"), []byte("{}"), 0755)

	tCdc = app.MakeDefaultCodec()
	proc, err := CreateProcess("", "mxwd", []string{"start", "--home", tWorkingDir})
	if err != nil {
		panic(err)
	}
	time.Sleep(1000)

	go func() {
		err := proc.Cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		out, err1, _ := proc.ReadAll()
		fmt.Printf("%s\n%s\n", string(out), string(err1))

		done <- struct{}{}
	}()

	tClient = WaitForRPC(tPort)
	WaitForNextHeightTM(tPort)

	return proc
}

func TestMain(m *testing.M) {
	dir, err := ioutil.TempDir("", "mxw")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting test node at " + dir)
	tWorkingDir = dir

	type key struct {
		Name        string
		MasterPriv  string
		DerivedPriv string
		Address     string
		Mnemonic    string
	}
	var keys []key
	content, _ := ioutil.ReadFile("./config/keys.json")
	json.Unmarshal(content, &keys)
	tKeys = make(map[string]*info)

	for _, k := range keys {
		bz, _ := hex.DecodeString(k.DerivedPriv)
		var priv [32]byte
		copy(priv[:], bz)
		addr, _ := sdkTypes.AccAddressFromBech32(k.Address)

		tKeys[k.Name] = &info{
			addr,
			secp256k1.PrivKeySecp256k1(priv),
			secp256k1.PrivKeySecp256k1(priv).PubKey(),
			k.Address,
			0,
		}

		proc, err := CreateProcess("", "mxwcli", []string{"keys", "import-mnemonic", k.Name, k.Mnemonic, "--encryption_passphrase", "12345678", "--home", tWorkingDir})
		if err != nil {
			panic(err)
		}

		err = proc.Cmd.Start()
		if err != nil {
			panic(err)
		}
		out, err1, _ := proc.ReadAll()
		fmt.Printf("%s%s", string(out), string(err1))
		proc.Cmd.Wait()
	}

	tKeys["nope"] = &info{
		sdkTypes.AccAddress{}, nil, nil,
		"", 0,
	}

	done := make(chan struct{})
	proc := startServer(done)

	exitCode := m.Run()

	proc.Stop(true)

	// waiting for gallactic to exit
	proc.Wait()
	<-done

	os.Exit(exitCode)
}
