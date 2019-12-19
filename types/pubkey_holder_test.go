package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestPubKeyHolder(t *testing.T) {
	var a PubKeyHolder

	pubkey1 := ed25519.PubKeyEd25519{1}
	pubkey2 := ed25519.PubKeyEd25519{2}
	pubkey3 := ed25519.PubKeyEd25519{3}
	pubkey4 := ed25519.PubKeyEd25519{4}
	pubkey5 := ed25519.PubKeyEd25519{5}

	s := []crypto.PubKey{pubkey3, pubkey4, pubkey5}

	assert.True(t, a.Append(pubkey1))
	assert.True(t, a.Append(pubkey2))
	assert.True(t, a.Append(pubkey3))
	assert.False(t, a.Append(pubkey2))
	assert.Equal(t, a.Size(), 3)
	assert.True(t, a.Remove(pubkey1))
	assert.Equal(t, a.Size(), 2)
	i, ok := a.Contains(pubkey2)
	assert.Equal(t, i, 0)
	assert.Equal(t, ok, true)
	i, ok = a.Contains(pubkey4)
	assert.Equal(t, i, -1)
	assert.Equal(t, ok, false)

	a.AppendPubKeys(s)
	assert.Equal(t, a.Size(), 4)

}
