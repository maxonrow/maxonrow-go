package auth

import (
	"encoding/binary"
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/kyc"
	"golang.org/x/crypto/sha3"
)

// Internal broadcast tx
type BroadcastTx func(ctx sdkTypes.Context, tx sdkTypes.Tx, simulate bool) (sdkTypes.Context, error)

func NewHandler(accountKeeper sdkAuth.AccountKeeper, kycKeeper kyc.Keeper, broadcastTx BroadcastTx) sdkTypes.Handler {
	return func(ctx sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		switch msg := msg.(type) {
		case MsgCreateMultiSigAccount:
			return handleMsgCreateMultiSigAccount(ctx, msg, accountKeeper, kycKeeper)
		case MsgUpdateMultiSigAccount:
			return handleMsgUpdateMultiSigAccount(ctx, msg, accountKeeper, kycKeeper)
		case MsgTransferMultiSigOwner:
			return handleMsgTransferMultiSigOwner(ctx, msg, accountKeeper, kycKeeper)
		case MsgCreateMultiSigTx:
			return handleMsgCreateMultiSigTx(ctx, msg, accountKeeper, kycKeeper)
		case MsgSignMultiSigTx:
			return handleMsgSignMultiSigTx(ctx, msg, accountKeeper, kycKeeper, broadcastTx)
		case MsgDeleteMultiSigTx:
			return handleMsgDeleteMultiSigTx(ctx, msg, accountKeeper, kycKeeper)
		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdkTypes.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateMultiSigAccount(ctx sdkTypes.Context, msg MsgCreateMultiSigAccount, accountKeeper auth.AccountKeeper, kycKeeper kyc.Keeper) sdkTypes.Result {
	OwnerAcc := accountKeeper.GetAccount(ctx, msg.Owner)
	if OwnerAcc == nil {
		return sdkTypes.ErrInvalidAddress(fmt.Sprintf("Invalid account address: %s", msg.Owner)).Result()
	}
	addr := DeriveMultiSigAddress(msg.Owner, OwnerAcc.GetSequence())

	// TODO: Check signers has passed KYC
	for _, signer := range msg.Signers {
		if !kycKeeper.IsWhitelisted(ctx, signer) {
			return sdkTypes.ErrUnknownRequest("Signer is not whitelisted.").Result()
		}
	}

	acc := accountKeeper.NewAccountWithAddress(ctx, addr)
	multisig := new(sdkTypes.MultiSig)
	multisig.Owner = msg.Owner
	multisig.Threshold = msg.Threshold
	multisig.Signers = msg.Signers
	acc.SetMultiSig(multisig)
	accountKeeper.SetAccount(ctx, acc)

	// TODO
	// Whitelisted this address in kyc keeper.
	kycKeeper.Whitelist(ctx, addr, "msig:"+addr.String())

	// TODO: Events
	accountSequence := OwnerAcc.GetSequence()
	var log string
	if accountSequence == 0 {
		log = types.MakeResultLog(accountSequence, ctx.TxBytes())
	} else {
		log = types.MakeResultLog(accountSequence-1, ctx.TxBytes())
	}

	eventParam := []string{msg.Owner.String(), acc.GetAddress().String()}
	eventSignature := "CreatedMultiSigAccount(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msg.Owner.String(), eventParam),
		Log:    log,
	}

}

func DeriveMultiSigAddress(addr sdkTypes.AccAddress, sequence uint64) sdkTypes.AccAddress {
	temp := make([]byte, 20+8)
	copy(temp, addr.Bytes())
	binary.BigEndian.PutUint64(temp[20:], sequence)
	hasher := sha3.New256()
	hasher.Write(temp) // does not error
	hash := hasher.Sum(nil)

	return sdkTypes.AccAddress(hash[12:])
}

func handleMsgUpdateMultiSigAccount(ctx sdkTypes.Context, msg MsgUpdateMultiSigAccount, accountKeeper auth.AccountKeeper, kycKeeper kyc.Keeper) sdkTypes.Result {

	groupAcc := accountKeeper.GetAccount(ctx, msg.GroupAddress)
	if groupAcc == nil {
		return sdkTypes.ErrUnknownRequest("Group address invalid.").Result()
	}

	ownerAccount := accountKeeper.GetAccount(ctx, msg.Owner)
	if ownerAccount == nil {
		return sdkTypes.ErrUnknownRequest("Owner address invalid.").Result()
	}

	// TO-DO: check add or remove signers.
	multiSig := groupAcc.GetMultiSig()

	if !multiSig.Owner.Equals(msg.Owner) {
		return sdkTypes.ErrUnknownRequest("Owner address invalid.").Result()
	}

	if len(multiSig.PendingTxs) > 0 {
		return sdkTypes.ErrUnknownRequest("Please clear the pending tx before editting signers.").Result()
	}

	multiSig.Signers = msg.NewSigners
	multiSig.Threshold = msg.NewThreshold
	groupAcc.SetMultiSig(multiSig)
	accountKeeper.SetAccount(ctx, groupAcc)

	// TO-DO: event
	accountSequence := ownerAccount.GetSequence()
	var log string
	if accountSequence == 0 {
		log = types.MakeResultLog(accountSequence, ctx.TxBytes())
	} else {
		log = types.MakeResultLog(accountSequence-1, ctx.TxBytes())
	}

	eventParam := []string{msg.Owner.String(), msg.GroupAddress.String()}
	eventSignature := "UpdatedMultiSigAccount(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msg.Owner.String(), eventParam),
		Log:    log,
	}

}

func handleMsgTransferMultiSigOwner(ctx sdkTypes.Context, msg MsgTransferMultiSigOwner, accountKeeper auth.AccountKeeper, kycKeeper kyc.Keeper) sdkTypes.Result {

	groupAcc := accountKeeper.GetAccount(ctx, msg.GroupAddress)
	if groupAcc == nil {
		return sdkTypes.ErrUnknownRequest("Group address invalid.").Result()
	}

	ownerAccount := accountKeeper.GetAccount(ctx, msg.Owner)
	if ownerAccount == nil {
		return sdkTypes.ErrUnknownRequest("Owner address invalid.").Result()
	}

	if !groupAcc.GetMultiSig().Owner.Equals(msg.Owner) {
		return sdkTypes.ErrUnknownRequest("Owner of group address invalid.").Result()
	}

	multiSig := groupAcc.GetMultiSig()
	multiSig.Owner = msg.NewOwner
	groupAcc.SetMultiSig(multiSig)
	accountKeeper.SetAccount(ctx, groupAcc)

	// TO-DO: event
	accountSequence := ownerAccount.GetSequence()
	var log string
	if accountSequence == 0 {
		log = types.MakeResultLog(accountSequence, ctx.TxBytes())
	} else {
		log = types.MakeResultLog(accountSequence-1, ctx.TxBytes())
	}

	eventParam := []string{msg.Owner.String(), msg.GroupAddress.String()}
	eventSignature := "TransferredMultiSigOwnership(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msg.Owner.String(), eventParam),
		Log:    log,
	}

}

func handleMsgCreateMultiSigTx(ctx sdkTypes.Context, msg MsgCreateMultiSigTx, accountKeeper auth.AccountKeeper, kycKeeper kyc.Keeper) sdkTypes.Result {

	groupAcc := accountKeeper.GetAccount(ctx, msg.GroupAddress)
	if groupAcc == nil {
		return sdkTypes.ErrUnknownRequest("Group address invalid.").Result()
	}

	senderAccount := accountKeeper.GetAccount(ctx, msg.Sender)
	if senderAccount == nil {
		return sdkTypes.ErrUnknownRequest("Sender address invalid.").Result()
	}

	if !groupAcc.GetMultiSig().IsSigner(msg.Sender) {
		return sdkTypes.ErrUnknownRequest("Sender is not signer of group address.").Result()
	}

	multiSig := groupAcc.GetMultiSig()
	txID := multiSig.GetNewTxID()

	pendingTx := sdkTypes.NewPendingTx(txID, msg.StdTx, msg.Sender, []sdkTypes.AccAddress{msg.Sender})

	err := multiSig.AddTx(pendingTx)
	if err != nil {
		return err.Result()
	}

	multiSig.IncCounter()
	groupAcc.SetMultiSig(multiSig)
	accountKeeper.SetAccount(ctx, groupAcc)

	// TO-DO: event
	accountSequence := senderAccount.GetSequence()
	var log string
	if accountSequence == 0 {
		log = types.MakeResultLog(accountSequence, ctx.TxBytes())
	} else {
		log = types.MakeResultLog(accountSequence-1, ctx.TxBytes())
	}

	eventParam := []string{msg.Sender.String(), msg.GroupAddress.String()}
	eventSignature := "CreatedMultiSigTx(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msg.Sender.String(), eventParam),
		Log:    log,
	}

}

func handleMsgSignMultiSigTx(ctx sdkTypes.Context, msg MsgSignMultiSigTx, accountKeeper auth.AccountKeeper, kycKeeper kyc.Keeper, broadcastTx BroadcastTx) sdkTypes.Result {

	groupAcc := accountKeeper.GetAccount(ctx, msg.GroupAddress)
	if groupAcc == nil {
		return sdkTypes.ErrUnknownRequest("Group address invalid.").Result()
	}

	senderAccount := accountKeeper.GetAccount(ctx, msg.Sender)
	if senderAccount == nil {
		return sdkTypes.ErrUnknownRequest("Sender address invalid.").Result()
	}

	if !groupAcc.GetMultiSig().IsSigner(msg.Sender) {
		return sdkTypes.ErrUnknownRequest("Sender is not signer of group address.").Result()
	}

	multiSig := groupAcc.GetMultiSig()
	multiSig.SignTx(msg.Sender, msg.TxID)

	groupAcc.SetMultiSig(multiSig)
	accountKeeper.SetAccount(ctx, groupAcc)

	isMetric := multiSig.IsMetric(msg.TxID)
	broadcastedEvents := sdkTypes.EmptyEvents()
	if isMetric {
		// TO-DO: broadcast tx
		tx := multiSig.GetTx(msg.TxID)
		if tx != nil {
			_, err := broadcastTx(ctx, tx, false)
			if err != nil {
				return sdkTypes.ErrInternal("Multisig broadcast tx failed.").Result()
			}

			// Event: broadcast tx
			broadcastedEventParam := []string{groupAcc.GetAddress().String(), string(msg.TxID)}
			broadcastedEventSignature := "BroadcastedTx(string,string)"
			broadcastedEvents = types.MakeMxwEvents(broadcastedEventSignature, groupAcc.GetAddress().String(), broadcastedEventParam)
		}

		isDeleted := multiSig.RemoveTx(msg.TxID)
		if !isDeleted {
			return sdkTypes.ErrUnknownRequest("Delete failed.").Result()
		}

		groupAcc.SetMultiSig(multiSig)
		accountKeeper.SetAccount(ctx, groupAcc)
	}

	// TO-DO: event
	accountSequence := senderAccount.GetSequence()
	var log string
	if accountSequence == 0 {
		log = types.MakeResultLog(accountSequence, ctx.TxBytes())
	} else {
		log = types.MakeResultLog(accountSequence-1, ctx.TxBytes())
	}

	eventParam := []string{msg.Sender.String(), msg.GroupAddress.String(), string(msg.TxID)}
	eventSignature := "SignedMultiSigTx(string,string,string)"
	events := types.MakeMxwEvents(eventSignature, msg.Sender.String(), eventParam)

	return sdkTypes.Result{
		Events: events.AppendEvents(broadcastedEvents),
		Log:    log,
	}

}

func handleMsgDeleteMultiSigTx(ctx sdkTypes.Context, msg MsgDeleteMultiSigTx, accountKeeper auth.AccountKeeper, kycKeeper kyc.Keeper) sdkTypes.Result {

	groupAcc := accountKeeper.GetAccount(ctx, msg.GroupAddress)
	if groupAcc == nil {
		return sdkTypes.ErrUnknownRequest("Group address invalid.").Result()
	}

	senderAccount := accountKeeper.GetAccount(ctx, msg.Sender)
	if senderAccount == nil {
		return sdkTypes.ErrUnknownRequest("Sender address invalid.").Result()
	}

	if !groupAcc.GetMultiSig().IsSigner(msg.Sender) {
		return sdkTypes.ErrUnknownRequest("Sender is not signer of group address.").Result()
	}

	multiSig := groupAcc.GetMultiSig()
	ok, pendingTx := multiSig.GetPendingTx(msg.TxID)
	if !ok {
		return sdkTypes.ErrUnknownRequest("Pending tx is not found.").Result()
	}

	if !pendingTx.Sender.Equals(msg.Sender) && !multiSig.Owner.Equals(msg.Sender) {
		return sdkTypes.ErrUnknownRequest("Sender is invalid.").Result()
	}

	isDeleted := multiSig.RemoveTx(msg.TxID)
	if !isDeleted {
		return sdkTypes.ErrUnknownRequest("Delete failed.").Result()
	}

	groupAcc.SetMultiSig(multiSig)
	accountKeeper.SetAccount(ctx, groupAcc)

	// TO-DO: event
	accountSequence := senderAccount.GetSequence()
	var log string
	if accountSequence == 0 {
		log = types.MakeResultLog(accountSequence, ctx.TxBytes())
	} else {
		log = types.MakeResultLog(accountSequence-1, ctx.TxBytes())
	}

	eventParam := []string{msg.Sender.String(), msg.GroupAddress.String(), string(msg.TxID)}
	eventSignature := "DeletedMultiSigTx(string,string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msg.Sender.String(), eventParam),
		Log:    log,
	}

}
