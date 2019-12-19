package types

import (
	"encoding/json"

	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/common"
)

type ResultLog struct {
	Hash  common.HexBytes `json:"hash"`
	Nonce uint64          `json:"nonce"`
}

func MakeResultLog(nonce uint64, txBytes []byte) string {
	respData, _ := json.Marshal(&ResultLog{
		Hash:  tmhash.Sum(txBytes),
		Nonce: nonce,
	})
	return string(respData)
}
