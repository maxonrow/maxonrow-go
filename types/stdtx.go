package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func GetSignBytes(ctx sdk.Context, tx sdkAuth.StdTx, acc exported.Account) []byte {
	genesis := ctx.BlockHeight() == 0
	chainID := ctx.ChainID()
	var accNum uint64
	if !genesis {
		accNum = acc.GetAccountNumber()
	}

	if acc.IsMultiSig() {
		// 1- check if tx exists in pending tx, then tx_id is same as pending.id
		// 2. if not, seq = acc.GetCounter()
		multisig := acc.GetMultiSig()
		txID, exist := multisig.ValidateMultiSigTx(tx)
		if exist {
			return sdkAuth.StdSignBytes(
				chainID, accNum, txID, tx.Fee, tx.Msgs, tx.Memo,
			)
		}
		return sdkAuth.StdSignBytes(
			chainID, accNum, multisig.GetCounter(), tx.Fee, tx.Msgs, tx.Memo,
		)
	} else {
		seq := acc.GetSequence()
		return sdkAuth.StdSignBytes(
			chainID, accNum, seq, tx.Fee, tx.Msgs, tx.Memo,
		)
	}
}
