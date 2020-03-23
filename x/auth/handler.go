package auth

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	sdkAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/utils"
	"github.com/maxonrow/maxonrow-go/x/kyc"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/common"
	rpc "github.com/tendermint/tendermint/rpc/core"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	"golang.org/x/crypto/ripemd160"
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

	for _, signer := range msg.Signers {
		if !kycKeeper.IsWhitelisted(ctx, signer) {
			return sdkTypes.ErrUnknownRequest("Signer is not whitelisted.").Result()
		}
	}

	acc := accountKeeper.NewAccountWithAddress(ctx, addr)
	multisig := sdkAuthTypes.NewMultiSig(msg.Owner, msg.Threshold, msg.Signers)
	acc.SetMultiSig(multisig)
	accountKeeper.SetAccount(ctx, acc)

	// Whitelisted this address in kyc keeper.
	kycKeeper.Whitelist(ctx, addr, "msig:"+addr.String())

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

	addrBz := addr.Bytes()
	sequenceBz := sdkTypes.Uint64ToBigEndian(sequence)
	sequenceBz = bytes.TrimLeft(sequenceBz, "\x00")
	temp := append(addrBz[:], sequenceBz[:]...)

	hasherSHA256 := sha256.New()
	hasherSHA256.Write(temp[:]) // does not error
	sha := hasherSHA256.Sum(nil)

	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha) // does not error

	return sdkTypes.AccAddress(hasherRIPEMD160.Sum(nil))
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

	if !multiSig.GetOwner().Equals(msg.Owner) {
		return sdkTypes.ErrUnknownRequest("Owner address invalid.").Result()
	}

	err := multiSig.UpdateSigners(msg.NewSigners, msg.NewThreshold)
	if err != nil {
		return sdkTypes.ResultFromError(err)
	}
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
	multiSig := groupAcc.GetMultiSig()
	if !multiSig.IsOwner(msg.Owner) {
		return sdkTypes.ErrUnknownRequest("Owner of group address invalid.").Result()
	}

	multiSig.UpdateOwner(msg.NewOwner)
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

	if !groupAcc.IsMultiSig() {
		return sdkTypes.ErrUnknownRequest("Sender is not a multisig account.").Result()
	}

	if !groupAcc.IsSigner(msg.Sender) {
		return sdkTypes.ErrUnknownRequest("Sender is not signer of group address.").Result()
	}

	multiSig := groupAcc.GetMultiSig()
	ptx, err := multiSig.AddPendingTx(msg.StdTx, msg.Sender)
	if err != nil {
		return sdkTypes.ResultFromError(err)
	}

	// We need to first update keeper, then check signatures
	groupAcc.SetMultiSig(multiSig)
	accountKeeper.SetAccount(ctx, groupAcc)

	stdTx, ok := ptx.GetTx().(sdkAuth.StdTx)
	if !ok {
		return sdkTypes.ErrInternal("Pending tx must be StdTx.").Result()
	}
	_, sigErr := utils.CheckTxSig(ctx, stdTx, accountKeeper, kycKeeper)
	if sigErr != nil {
		return sdkTypes.ResultFromError(sigErr)
	}

	internalHash, broadcastedEvents, metricErr := checkIsMetric(ctx, ptx.GetID(), groupAcc, txEncoder)
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

	if !groupAcc.IsMultiSig() {
		return sdkTypes.ErrUnknownRequest("Sender is not a multisig account.").Result()
	}

	if !groupAcc.IsSigner(msg.Sender) {
		return sdkTypes.ErrUnknownRequest("Sender is not signer of group address.").Result()
	}

	multiSig := groupAcc.GetMultiSig()

	ptx, err := multiSig.AddSignature(msg.TxID, msg.Sender, msg.Signature.Signature)
	if err != nil {
		return sdkTypes.ResultFromError(err)
	}

	// We need to first update keeper, then check signatures
	groupAcc.SetMultiSig(multiSig)
	accountKeeper.SetAccount(ctx, groupAcc)

	stdTx, ok := ptx.GetTx().(sdkAuth.StdTx)
	if !ok {
		return sdkTypes.ErrInternal("Pending tx must be StdTx.").Result()
	}
	_, sigErr := utils.CheckTxSig(ctx, stdTx, accountKeeper, kycKeeper)
	if sigErr != nil {
		return sdkTypes.ResultFromError(sigErr)
	}

	accountSequence := senderAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	internalHash, broadcastedEvents, metricErr := checkIsMetric(ctx, msg.TxID, groupAcc, txEncoder)
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

	if !groupAcc.IsMultiSig() {
		return sdkTypes.ErrUnknownRequest("Sender is not a multisig account.").Result()
	}

	multiSig := groupAcc.GetMultiSig()
	pendingTx := multiSig.GetPendingTx(msg.TxID)
	if pendingTx == nil {
		return sdkTypes.ErrUnknownRequest("Pending tx is not found.").Result()
	}

	if !multiSig.IsOwner(msg.Sender) {
		return sdkTypes.ErrUnknownRequest("Only group account owner can remove pending tx.").Result()
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

func checkIsMetric(ctx sdkTypes.Context, txID uint64, groupAcc exported.Account, txEncoder sdkTypes.TxEncoder) (common.HexBytes, sdkTypes.Events, sdkTypes.Error) {

	multiSig := groupAcc.GetMultiSig()
	isMetric := multiSig.IsMetric(txID)
	broadcastedEvents := sdkTypes.EmptyEvents()
	var internalHash common.HexBytes
	if isMetric {
		ptx := multiSig.GetPendingTx(txID)
		if ptx == nil {
			return nil, nil, sdkTypes.ErrInternal("There is no pending tx.")
		}

		bz, err := txEncoder(ptx.GetTx())
		if err != nil {
			return nil, nil, sdkTypes.ErrInternal("Error encoding pending tx.")
		}

		internalHash = tmhash.Sum(bz)

		var rpcCtx *rpctypes.Context
		go func() {
			// sleep 3 seconds to make sure what ever necessary is completed in previous block.
			time.Sleep(3 * time.Second)
			res, err := rpc.BroadcastTxSync(rpcCtx, bz)
			if err != nil {
				ctx.Logger().Error("Panic on broadcasting internal transaction", "Error", err)
				panic(err)
			}

			if res.Code != 0 {
				ctx.Logger().Error("Broadcasting internal transaction failed", "Result", res)

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
