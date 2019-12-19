package types

import (
	"github.com/tendermint/go-amino"
	tmCrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmTypes "github.com/tendermint/tendermint/types"
)

type PVKeyFile struct {
	Address tmTypes.Address  `json:"address"`
	PubKey  tmCrypto.PubKey  `json:"pub_key"`
	PrivKey tmCrypto.PrivKey `json:"priv_key"`
}

func GenPrvKeyFile() PVKeyFile {
	privKey := ed25519.GenPrivKey()
	return PVKeyFile{
		Address: privKey.PubKey().Address(),
		PubKey:  privKey.PubKey(),
		PrivKey: privKey,
	}
}

func (pvKey PVKeyFile) Save(filePath string, cdc *amino.Codec) {
	if filePath == "" {
		panic("cannot save PrivValidator key: filePath not set")
	}

	jsonBytes, err := cdc.MarshalJSONIndent(pvKey, "", "  ")
	if err != nil {
		panic(err)
	}
	err = cmn.WriteFileAtomic(filePath, jsonBytes, 0600)
	if err != nil {
		panic(err)
	}

}
