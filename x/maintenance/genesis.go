package maintenance

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type GenesisState struct {
	Maintainers        []sdkTypes.AccAddress `json:"maintainers"`
	StartingProposalID uint64                `json:"starting_proposal_id"`
	ValidatorSet       []string              `json:"validator_set"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		StartingProposalID: 1,
	}
}

func InitGenesis(ctx sdkTypes.Context, keeper *Keeper, genesisState GenesisState) {
	var pks []crypto.PubKey
	for _, s := range genesisState.ValidatorSet {
		pk, err := sdkTypes.GetConsPubKeyBech32(s)
		if err != nil {
			panic(err)
		}
		pks = append(pks, pk)
	}
	keeper.setMaintainersAddresses(ctx, genesisState.Maintainers)
	keeper.WhitelistValidator(ctx, pks)
	err := keeper.setInitialProposalID(ctx, genesisState.StartingProposalID)
	if err != nil {
		panic(err)
	}
}

func ExportGenesis(ctx sdkTypes.Context, keeper *Keeper) GenesisState {
	startingProposalID, _ := keeper.peekCurrentProposalID(ctx)
	return GenesisState{
		StartingProposalID: startingProposalID,
	}
}
