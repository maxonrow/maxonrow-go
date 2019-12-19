package maintenance

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the governance Querier
const (
	QueryProposal = "proposal"
)

func NewQuerier(keeper *Keeper) sdkTypes.Querier {
	return func(ctx sdkTypes.Context, path []string, req abci.RequestQuery) ([]byte, sdkTypes.Error) {
		switch path[0] {
		case QueryProposal:
			return queryProposal(ctx, path[1:], req, keeper)
		default:
			return nil, sdkTypes.ErrUnknownRequest("unknown gov query endpoint")
		}
	}
}

type QueryProposalParams struct {
	ProposalID uint64
}

// creates a new instance of QueryProposalParams
func NewQueryProposalParams(proposalID uint64) QueryProposalParams {
	return QueryProposalParams{
		ProposalID: proposalID,
	}
}

// nolint: unparam
func queryProposal(ctx sdkTypes.Context, path []string, req abci.RequestQuery, keeper *Keeper) ([]byte, sdkTypes.Error) {

	var params QueryProposalParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkTypes.ErrUnknownRequest(sdkTypes.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	proposal, ok := keeper.GetProposal(ctx, params.ProposalID)
	if !ok {
		return nil, ErrUnknownProposal(DefaultCodespace, params.ProposalID)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, proposal)
	if err != nil {
		return nil, sdkTypes.ErrInternal(sdkTypes.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
