package types

import (
	"encoding/json"

	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/common"
)

type ResultLog struct {
	Hash         common.HexBytes `json:"hash"`
	InternalHash common.HexBytes `json:"internalHash"`
	Nonce        uint64          `json:"nonce"`
}

func NewResultLog(nonce uint64, txBytes []byte) *ResultLog {

	if nonce != 0 {
		nonce = nonce - 1
	}
	return &ResultLog{
		Hash:  tmhash.Sum(txBytes),
		Nonce: nonce,
	}
}

func (r *ResultLog) WithInternalHash(internalHash common.HexBytes) *ResultLog {
	r.InternalHash = internalHash
	return r
}

func (r *ResultLog) String() string {
	respData, _ := json.Marshal(r)
	return string(respData)
}
