package types

import (
	"github.com/tendermint/tendermint/crypto"
)

type PubKeyHolder []crypto.PubKey

func (a *PubKeyHolder) AppendPubKeys(pubKeys []crypto.PubKey) {
	for _, pubKey := range pubKeys {
		a.Append(pubKey)
	}
}

func (a *PubKeyHolder) Append(pubKey crypto.PubKey) bool {
	if _, exist := a.Contains(pubKey); exist {
		return false
	}
	*a = append(*a, pubKey)
	return true
}

func (a *PubKeyHolder) Contains(pubKey crypto.PubKey) (int, bool) {
	for i, e := range *a {
		if e.Equals(pubKey) {
			return i, true
		}
	}
	return -1, false
}

func (a *PubKeyHolder) Remove(pubKey crypto.PubKey) bool {
	index, exists := a.Contains(pubKey)
	if exists {
		s := *a
		*a = append(s[:index], s[index+1:]...)
		return true
	}
	return false
}

func (a *PubKeyHolder) Size() int {
	return len(*a)
}
