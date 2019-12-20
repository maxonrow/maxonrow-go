package main

import (
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

// REST handler to get the latest block
func BlockResultsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		height, err := strconv.ParseInt(vars["height"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				"ERROR: Couldn't parse block height. Assumed format is '/block/{height}'.")
			return
		}
		chainHeight, err := rpc.GetChainHeight(cliCtx)
		if height > chainHeight {
			rest.WriteErrorResponse(w, http.StatusNotFound,
				"ERROR: Requested block height is bigger then the chain length.")
			return
		}

		// get the node
		node, err := cliCtx.GetNode()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound,
				"ERROR: Requested block height invalid.")
			return
		}

		// header -> BlockchainInfo
		// header, tx -> Block
		// results -> BlockResults
		res, err := node.BlockResults(&height)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound,
				"ERROR: Requested block height is bigger then the chain length.")
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
