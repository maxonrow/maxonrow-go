package maintenance

import (
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/maxonrow/maxonrow-go/types"
)

const (
	ADD      = "add"
	REMOVE   = "remove"
	APPROVE  = "approve"
	REJECT   = "reject"
	EXECUTED = "executed"
)

func NewHandler(keeper *Keeper, accountKeeper *auth.AccountKeeper) sdkTypes.Handler {
	return func(ctx sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		switch msg := msg.(type) {

		case MsgProposal:
			return handleMsgSubmitProposal(ctx, keeper, msg, accountKeeper)
		case MsgCastAction:
			return handleMsgCastAction(ctx, keeper, msg, accountKeeper)
		default:
			errMsg := fmt.Sprintf("Unrecognized fee Msg type: %v", msg.Type())
			return sdkTypes.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitProposal(ctx sdkTypes.Context, keeper *Keeper, msg MsgProposal, accountKeeper *auth.AccountKeeper) sdkTypes.Result {

	proposerAccount := accountKeeper.GetAccount(ctx, msg.Proposer)
	if proposerAccount == nil {
		return sdkTypes.ErrInternal("Invalid propeser account.").Result()
	}

	// Check if maintainers
	if !keeper.IsMaintainers(ctx, proposerAccount.GetAddress()) {
		return sdkTypes.ErrUnauthorized("Not authorised to submit proposal.").Result()
	}

	content := NewTextProposal(msg.Title, msg.Description, msg.ProposalType)

	proposal, submitProposalErr := keeper.SubmitProposal(ctx, content, msg.ProposalData)
	if submitProposalErr != nil {
		return submitProposalErr.Result()
	}

	approveErr, executedEvents := keeper.ApproveProposal(ctx, proposal.ProposalID, msg.Proposer)
	if approveErr != nil {
		return approveErr.Result()
	}

	// When proposal just submitted, the executedTags should be empty.
	if executedEvents != nil {
		return ErrInactiveProposal(DefaultCodespace, proposal.ProposalID).Result()
	}

	// Event: Submit Proposal
	proposalEventParam := []string{string(proposal.ProposalID), msg.Proposer.String(), proposal.ProposalType().String()}
	proposalEventSignature := "SubmittedProposal(string,string,string)"
	proposalEvents := types.MakeMxwEvents(proposalEventSignature, msg.Proposer.String(), proposalEventParam)

	// Event: Add Vote
	approveEventParam := []string{APPROVE, string(proposal.ProposalID), msg.Proposer.String()}
	approveEventSignature := "ApprovedProposal(string,string,string)"
	approveEvents := types.MakeMxwEvents(approveEventSignature, msg.Proposer.String(), approveEventParam)

	accountSequence := proposerAccount.GetSequence()

	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: proposalEvents.AppendEvents(approveEvents),
		Log:    resultLog.String(),
	}

}

func handleMsgCastAction(ctx sdkTypes.Context, keeper *Keeper, msg MsgCastAction, accountKeeper *auth.AccountKeeper) sdkTypes.Result {

	account := accountKeeper.GetAccount(ctx, msg.Owner)
	if account == nil {
		return sdkTypes.ErrInternal("Invalid account.").Result()
	}

	// Check if maintainers
	if !keeper.IsMaintainers(ctx, account.GetAddress()) {
		return sdkTypes.ErrUnauthorized("Not authorised to submit proposal.").Result()
	}

	switch msg.Action {
	case APPROVE:
		err, executedEvents := keeper.ApproveProposal(ctx, msg.ProposalID, msg.Owner)
		if err != nil {
			return err.Result()
		}

		// Event: approve proposal
		approveEventParam := []string{APPROVE, string(msg.ProposalID), msg.Owner.String()}
		approveEventSignature := "ApprovedProposal(string,string,string)"
		approveEvents := types.MakeMxwEvents(approveEventSignature, msg.Owner.String(), approveEventParam)

		accountSequence := account.GetSequence()
		resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

		if executedEvents != nil {
			approveEvents = approveEvents.AppendEvents(executedEvents)
		}
		return sdkTypes.Result{
			Events: approveEvents,
			Log:    resultLog.String(),
		}

	case REJECT:
		err, executedEvents := keeper.RejectProposal(ctx, msg.ProposalID, msg.Owner)
		if err != nil {
			return err.Result()
		}

		// Event: reject proposal
		rejectEventParam := []string{REJECT, string(msg.ProposalID), msg.Owner.String()}
		rejectEventSignature := "RejectedProposal(string,string,string)"
		rejectEvents := types.MakeMxwEvents(rejectEventSignature, msg.Owner.String(), rejectEventParam)

		accountSequence := account.GetSequence()
		resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

		if executedEvents != nil {
			rejectEvents = rejectEvents.AppendEvents(executedEvents)
		}
		return sdkTypes.Result{
			Events: rejectEvents,
			Log:    resultLog.String(),
		}

	}
	return sdkTypes.Result{}
}
