package tests

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"testing"

	"github.com/maxonrow/maxonrow-go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendCli(t *testing.T) {
	proc, err := utils.CreateProcess("", "mxwcli", []string{"tx", "send", tKeys["eve"].addrStr, tKeys["acc-20"].addrStr, "1cin", "--gas", "0", "--fees", "200000000cin", "--chain-id", "maxonrow-chain", "--home", tWorkingDir, "--broadcast-mode", "block", "-y", "--keyring-backend", "test"})

	// key password
	_, err = proc.StdinPipe.Write([]byte("12345678\n"))
	require.NoError(t, err)

	err = proc.Cmd.Start()
	require.NoError(t, err)

	out, err1, err := proc.ReadAll()
	require.NoError(t, err)
	fmt.Printf("%s\n%s\n", string(out), string(err1))
	res := string(out)

	re := regexp.MustCompile(`(txhash:.*)`)
	matched := re.FindString(res)
	require.NotEmpty(t, matched)
	hash, err := hex.DecodeString(matched[8:])
	require.NoError(t, err)

	WaitForNextHeightTM(tPort)
	res1 := Tx(hash)
	assert.Zero(t, res1.TxResult.Code)
}
