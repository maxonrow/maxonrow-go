package auth

import (
	"encoding/binary"
	"fmt"
	"time"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/kyc"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/common"
	rpc "github.com/tendermint/tendermint/rpc/core"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	"golang.org/x/crypto/sha3"
)

func NewHandler(accountKeeper sdkAuth.AccountKeeper, kycKeeper kyc.Keeper, txEncoder sdkTypes.TxEncoder) sdkTypes.Handler {
	return func(ctx sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		switch msg := msg.(type) {
		case MsgCreateMultiSigAccount:
			return handleMsgCreateMultiSigAccount(ctx, msg, accountKeeper, kycKeeper)
		case MsgUpdateMultiSigAccount:
			return handleMsgUpdateMultiSigAccount(ctx, msg, accountKeeper, kycKeeper)
		case MsgTransferMultiSigOwner:
			return handleMsgTransferMultiSigOwner(ctx, msg, accountKeeper, kycKeeper)
		case MsgCreateMultiSigTx:
			return handleMsgCreateMultiSigTx(ctx, msg, accountKeeper, kycKeeper, txEncoder)
		case MsgSignMultiSigTx:
			return handleMsgSignMultiSigTx(ctx, msg, accountKeeper, kycKeeper, txEncoder)
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
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	eventParam := []string{msg.Owner.String(), acc.GetAddress().String()}
	eventSignature := "CreatedMultiSigAccount(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msg.Owner.String(), eventParam),
		Log:    resultLog.String(),
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
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	eventParam := []string{msg.Owner.String(), msg.GroupAddress.String()}
	eventSignature := "UpdatedMultiSigAccount(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msg.Owner.String(), eventParam),
		Log:    resultLog.String(),
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
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	eventParam := []string{msg.Owner.String(), msg.GroupAddress.String()}
	eventSignature := "TransferredMultiSigOwnership(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msg.Owner.String(), eventParam),
		Log:    resultLog.String(),
	}

}

func handleMsgCreateMultiSigTx(ctx sdkTypes.Context, msg MsgCreateMultiSigTx, accountKeeper auth.AccountKeeper, kycKeeper kyc.Keeper, txEncoder sdkTypes.TxEncoder) sdkTypes.Result {

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
	txID := multiSig.GetCounter()

	pendingTx := sdkTypes.NewPendingTx(txID, msg.StdTx, msg.Sender, []sdkTypes.AccAddress{msg.Sender})

	err := multiSig.AddTx(pendingTx)
	if err != nil {
		return err.Result()
	}

	multiSig.IncCounter()
	groupAcc.SetMultiSig(multiSig)
	accountKeeper.SetAccount(ctx, groupAcc)

	internalHash, broadcastedEvents, metricErr := checkIsMetric(txID, groupAcc, txEncoder)
	if metricErr != nil {
		return metricErr.Result()
	}

	// TO-DO: event
	accountSequence := senderAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	if internalHash != nil {
		resultLog = resultLog.WithInternalHash(internalHash)
	}

	eventParam := []string{msg.Sender.String(), msg.GroupAddress.String()}
	eventSignature := "CreatedMultiSigTx(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msg.Sender.String(), eventParam).AppendEvents(broadcastedEvents),
		Log:    resultLog.String(),
	}

}

func handleMsgSignMultiSigTx(ctx sdkTypes.Context, msg MsgSignMultiSigTx, accountKeeper auth.AccountKeeper, kycKeeper kyc.Keeper, txEncoder sdkTypes.TxEncoder) sdkTypes.Result {

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

	//TO-DO: append real signatures also
	multiSig.SignTx(msg.Sender, msg.TxID)
	pendingTx := multiSig.GetTx(msg.TxID)
	if pendingTx == nil {
		return sdkTypes.ErrInternal("Pending tx not found.").Result()
	}

	stdTx, ok := pendingTx.(sdkAuth.StdTx)
	if !ok {
		return sdkTypes.ErrInternal("Pending tx must be StdTx.").Result()
	}

	stdTx.Signatures = append(stdTx.Signatures, msg.Signature)
	multiSig.UpdatePendingTx(msg.TxID, stdTx)

	accountKeeper.SetAccount(ctx, groupAcc)

	accountSequence := senderAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	internalHash, broadcastedEvents, metricErr := checkIsMetric(msg.TxID, groupAcc, txEncoder)
	if metricErr != nil {
		return metricErr.Result()
	}

	if internalHash != nil {
		resultLog = resultLog.WithInternalHash(internalHash)
	}

	eventParam := []string{msg.Sender.String(), msg.GroupAddress.String(), string(msg.TxID)}
	eventSignature := "SignedMultiSigTx(string,string,string)"
	events := types.MakeMxwEvents(eventSignature, msg.Sender.String(), eventParam)

	return sdkTypes.Result{
		Events: events.AppendEvents(broadcastedEvents),
		Log:    resultLog.String(),
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
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	eventParam := []string{msg.Sender.String(), msg.GroupAddress.String(), string(msg.TxID)}
	eventSignature := "DeletedMultiSigTx(string,string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, msg.Sender.String(), eventParam),
		Log:    resultLog.String(),
	}

}

func checkIsMetric(txID uint64, groupAcc exported.Account, txEncoder sdkTypes.TxEncoder) (common.HexBytes, sdkTypes.Events, sdkTypes.Error) {

	multiSig := groupAcc.GetMultiSig()
	isMetric := multiSig.IsMetric(txID)
	broadcastedEvents := sdkTypes.EmptyEvents()
	var internalHash common.HexBytes
	if isMetric {
		tx := multiSig.GetTx(txID)
		if tx == nil {
			return nil, nil, sdkTypes.ErrInternal("There is no pending tx.")
		}

		bz, err := txEncoder(tx)
		if err != nil {
			return nil, nil, sdkTypes.ErrInternal("Error encoding pending tx.")
		}

		internalHash = tmhash.Sum(bz)

		var rpcCtx *rpctypes.Context
		go func() {
			// sleep 3 seconds to make sure what ever necessary is completed in preivous block.
			time.Sleep(3 * time.Second)
			_, err := rpc.BroadcastTxSync(rpcCtx, bz)
			if err != nil {
				panic(err)
			}
		}()

		// Event: broadcast tx
		broadcastedEventParam := []string{groupAcc.GetAddress().String(), string(txID)}
		broadcastedEventSignature := "BroadcastedTx(string,string)"
		broadcastedEvents = types.MakeMxwEvents(broadcastedEventSignature, groupAcc.GetAddress().String(), broadcastedEventParam)

		return internalHash, broadcastedEvents, nil
	}
	return nil, nil, nil
}
