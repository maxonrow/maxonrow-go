package auth

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryMultiSigAcc = "get_multisig_acc"
)

func NewQuerier(cdc *codec.Codec, accountKeeper sdkAuth.AccountKeeper) sdkTypes.Querier {
	return func(ctx sdkTypes.Context, path []string, req abci.RequestQuery) ([]byte, sdkTypes.Error) {
		switch path[0] {
		case QueryMultiSigAcc:
			return queryMultiSigAcc(cdc, ctx, path[1:], req, accountKeeper)
		default:
			return nil, sdkTypes.ErrUnknownRequest("unknown mxw/Auth query endpoint")
		}
	}
}

func queryMultiSigAcc(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, accountKeeper sdkAuth.AccountKeeper) ([]byte, sdkTypes.Error) {

	if len(path) != 1 {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid path %s", strings.Join(path, "/")))
	}

	groupAddr, err := sdkTypes.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkTypes.ErrUnknownAddress(fmt.Sprintf("Invliad group address %s", path[0]))
	}

	groupAcc := accountKeeper.GetAccount(ctx, groupAddr)

	respData := cdc.MustMarshalJSON(groupAcc)

	return respData, nil
}

type GroupAccount struct {
	
}
