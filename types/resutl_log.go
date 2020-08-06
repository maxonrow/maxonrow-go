package types

import (
	"encoding/json"

	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/common"
)

type ResultLog struct {
	Hash         common.HexBytes `json:"hash"`
	InternalHash common.HexBytes `json:"internalHash,omitempty"`
	Nonce        uint64          `json:"nonce"`
}

func ResultLogFromTMLog(log string) *ResultLog {
	type _tmLog struct {
		Log string `json:"log"`
	}
	var tmLog []_tmLog
	json.Unmarshal([]byte(log), &tmLog)

	resultLog := new(ResultLog)
	json.Unmarshal([]byte(tmLog[0].Log), resultLog)

	return resultLog
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
