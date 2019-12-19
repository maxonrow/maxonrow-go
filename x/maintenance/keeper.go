package maintenance

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/types"
)

const (
	APPROVERS_THRESHOLD = 3
	REJECTERS_THRESHOLD = 3
)

type ExecuteMaintenaceProposalFn func(ctx sdkTypes.Context, proposal Proposal) sdkTypes.Error

type Keeper struct {
	maintenanceKey  sdkTypes.StoreKey
	validatorSetKey sdkTypes.StoreKey
	cdc             *codec.Codec
	execProposalFn  ExecuteMaintenaceProposalFn
}

func NewKeeper(cdc *codec.Codec, maintenanceKey, validatorSetKey sdkTypes.StoreKey, execProposalFn ExecuteMaintenaceProposalFn) Keeper {
	return Keeper{
		cdc:             cdc,
		maintenanceKey:  maintenanceKey,
		validatorSetKey: validatorSetKey,
		execProposalFn:  execProposalFn,
	}
}

func (k *Keeper) setMaintainersAddresses(ctx sdkTypes.Context, addresses []sdkTypes.AccAddress) {

	ah := k.getMaintainersAddress(ctx)
	ah.AppendAccAddrs(addresses)

	store := ctx.KVStore(k.maintenanceKey)
	key := getMaintenanceAddressKey()
	val, err := k.cdc.MarshalBinaryLengthPrefixed(&ah)
	if err != nil {
		panic("Setting maintainers address failed.")
	}
	store.Set([]byte(key), val)

}

func (k *Keeper) getMaintainersAddress(ctx sdkTypes.Context) types.AddressHolder {
	var ah types.AddressHolder
	store := ctx.KVStore(k.maintenanceKey)
	key := getMaintenanceAddressKey()
	val := store.Get([]byte(key))
	if val == nil {
		return ah
	}
	err := k.cdc.UnmarshalBinaryLengthPrefixed(val, &ah)
	if err != nil {
		panic("Get maintainers address failed")
	}

	return ah
}

func (k *Keeper) IsMaintainers(ctx sdkTypes.Context, address sdkTypes.AccAddress) bool {
	existingAddrs := k.getMaintainersAddress(ctx)
	_, isExist := existingAddrs.Contains(address)
	return isExist
}

// Proposals
func (keeper Keeper) SubmitProposal(ctx sdkTypes.Context, content ProposalContent, data MsgProposalData) (proposal *Proposal, err sdkTypes.Error) {
	proposalID, err := keeper.getNewProposalID(ctx)
	if err != nil {
		return nil, err
	}

	submitTime := ctx.BlockHeader().Time

	proposal = &Proposal{
		ProposalContent: content,
		ProposalID:      proposalID,
		ProposalData:    data,
		Status:          StatusActive,
		SubmitTime:      submitTime,
	}

	keeper.SetProposal(ctx, *proposal)
	return
}

// Gets the next available ProposalID and increments it
func (keeper Keeper) getNewProposalID(ctx sdkTypes.Context) (proposalID uint64, err sdkTypes.Error) {
	store := ctx.KVStore(keeper.maintenanceKey)
	bz := store.Get(KeyNextProposalID)
	if bz == nil {
		return 0, sdkTypes.ErrInternal("Proposal initialID never set")
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposalID)
	bz = keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID + 1)
	store.Set(KeyNextProposalID, bz)
	return proposalID, nil
}

// Set proposal
func (keeper Keeper) SetProposal(ctx sdkTypes.Context, proposal Proposal) {
	store := ctx.KVStore(keeper.maintenanceKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposal)
	store.Set(getProposalKey(proposal.ProposalID), bz)
}

// Delete proposal
func (keeper Keeper) DeleteProposal(ctx sdkTypes.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.maintenanceKey)
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		panic("DeleteProposal cannot fail to GetProposal.")
	}
	store.Delete(getProposalKey(proposal.ProposalID))
}

// Get Proposal from store by ProposalID
func (keeper Keeper) GetProposal(ctx sdkTypes.Context, proposalID uint64) (proposal Proposal, ok bool) {
	store := ctx.KVStore(keeper.maintenanceKey)
	bz := store.Get(getProposalKey(proposalID))
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposal)
	return proposal, true
}

func (keeper Keeper) IsProposalActive(ctx sdkTypes.Context, proposalId uint64) (active bool) {
	proposal, ok := keeper.GetProposal(ctx, proposalId)
	if !ok {
		return false
	}
	if proposal.Status == StatusActive {
		return true
	}
	return false
}

// Set the initial proposal ID
func (keeper Keeper) setInitialProposalID(ctx sdkTypes.Context, proposalID uint64) sdkTypes.Error {
	store := ctx.KVStore(keeper.maintenanceKey)
	bz := store.Get(KeyNextProposalID)
	if bz != nil {
		return ErrInvalidGenesis(DefaultCodespace, "Initial ProposalID already set")
	}
	bz = keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(KeyNextProposalID, bz)
	return nil
}

// Peeks the next available ProposalID without incrementing it
func (keeper Keeper) peekCurrentProposalID(ctx sdkTypes.Context) (proposalID uint64, err sdkTypes.Error) {
	store := ctx.KVStore(keeper.maintenanceKey)
	bz := store.Get(KeyNextProposalID)
	if bz == nil {
		return 0, ErrInvalidGenesis(DefaultCodespace, "InitialProposalID never set")
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposalID)
	return proposalID, nil
}

func (keeper Keeper) ApproveProposal(ctx sdkTypes.Context, proposalID uint64, approver sdkTypes.AccAddress) (sdkTypes.Error, sdkTypes.Events) {
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return ErrUnknownProposal(DefaultCodespace, proposalID), nil
	}
	if proposal.Status != StatusActive {
		return ErrInactiveProposal(DefaultCodespace, proposalID), nil
	}

	proposal.Approvers.Append(approver)
	keeper.SetProposal(ctx, proposal)

	var executedEvents sdkTypes.Events
	if len(proposal.Approvers) == APPROVERS_THRESHOLD {
		executeProposalErr := keeper.execProposalFn(ctx, proposal)
		if executeProposalErr != nil {
			return executeProposalErr, nil
		}
		proposal.Status = StatusCompleted
		keeper.SetProposal(ctx, proposal)

		// Event: executed proposal
		executedEventParam := []string{EXECUTED, "mxw000000000000000000000000000000000000000", string(proposal.ProposalID)}
		executedEventSignature := "ExecutedProposal(string,string,string)"
		executedEvents = types.MakeMxwEvents(executedEventSignature, "mxw000000000000000000000000000000000000000", executedEventParam)

	}

	return nil, executedEvents

}

func (keeper Keeper) RejectProposal(ctx sdkTypes.Context, proposalID uint64, rejecter sdkTypes.AccAddress) (sdkTypes.Error, sdkTypes.Events) {
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return ErrUnknownProposal(DefaultCodespace, proposalID), nil
	}
	if proposal.Status != StatusActive {
		return ErrInactiveProposal(DefaultCodespace, proposalID), nil
	}

	proposal.Rejecters.Append(rejecter)

	keeper.SetProposal(ctx, proposal)

	var executedEvents sdkTypes.Events
	if len(proposal.Rejecters) == REJECTERS_THRESHOLD {
		executeProposalErr := keeper.execProposalFn(ctx, proposal)
		if executeProposalErr != nil {
			return executeProposalErr, nil
		}
		proposal.Status = StatusRejected
		keeper.SetProposal(ctx, proposal)

		// Event: executed proposal
		executedEventParam := []string{EXECUTED, "mxw000000000000000000000000000000000000000", string(proposal.ProposalID)}
		executedEventSignature := "ExecutedProposal(string,string,string)"
		executedEvents = types.MakeMxwEvents(executedEventSignature, "mxw000000000000000000000000000000000000000", executedEventParam)
	}

	return nil, executedEvents

}
