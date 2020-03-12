package app

import (
	"fmt"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
)

// const addr string = "http://localhost"

func (app *mxwApp) deliverMultiSigTxAsync(tx []byte) *ctypes.ResultBroadcastTx {

	fmt.Printf("\n============START : deliverMultiSigTxAsync()\n")
	return BroadcastTxAsync(tx)
}

func BroadcastTxAsync(tx []byte) *ctypes.ResultBroadcastTx {
	result := new(ctypes.ResultBroadcastTx)
	fmt.Printf("\n============BroadcastTxAsync() - tx : %v\n", string(tx))

	client := rpcclient.NewJSONRPCClient("tcp://localhost:26657") // tcp://0.0.0.0:32603
	_, err := client.Call("broadcast_tx_async", map[string]interface{}{"tx": tx}, result)

	if err != nil {
		panic(err)
	}

	return result
}
