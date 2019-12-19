package maintenance

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/maxonrow/maxonrow-go/types"
)

type WhitelistValidator struct {
	Action           string          `json:"action"`
	ValidatorPubKeys []crypto.PubKey `json:"validator_pubkeys"`
}

func NewWhitelistValidator(action string, validatorPubKeys []crypto.PubKey) WhitelistValidator {
	return WhitelistValidator{
		Action:           action,
		ValidatorPubKeys: validatorPubKeys,
	}
}

var _ MsgProposalData = &WhitelistValidator{}

func (whitelistValidator WhitelistValidator) GetType() ProposalKind {
	return ProposalTypesModifyValidatorSet
}

func (whitelistValidator *WhitelistValidator) Unmarshal(data []byte) error {
	err := msgCdc.UnmarshalBinaryLengthPrefixed(data, whitelistValidator)
	if err != nil {
		return err
	}
	return nil
}

func (whitelistValidator WhitelistValidator) Marshal() ([]byte, error) {
	bz, err := msgCdc.MarshalBinaryLengthPrefixed(whitelistValidator)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

// WhitelistValidator Sets or Removes Validator into validator set.
func (k *Keeper) WhitelistValidator(ctx sdkTypes.Context, pubkeys []crypto.PubKey) {
	validatorStore := ctx.KVStore(k.validatorSetKey)
	key := getValidatorSetKey()

	pubKeyHolder := k.GetValidatorSet(ctx)
	pubKeyHolder.AppendPubKeys(pubkeys)
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(pubKeyHolder)
	if err != nil {
		panic(err)
	}

	validatorStore.Set(key, bz)
}

func (k *Keeper) GetValidatorSet(ctx sdkTypes.Context) types.PubKeyHolder {
	var pubKeyHolder types.PubKeyHolder
	valSetStore := ctx.KVStore(k.validatorSetKey)
	key := getValidatorSetKey()

	bz := valSetStore.Get(key)
	if bz == nil {
		return pubKeyHolder
	}
	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &pubKeyHolder)
	if err != nil {
		panic(err)
	}
	return pubKeyHolder
}

func (k *Keeper) RevokeValidator(ctx sdkTypes.Context, pubKeyrs []crypto.PubKey) {

	valSetStore := ctx.KVStore(k.validatorSetKey)

	key := getValidatorSetKey()
	valSet := k.GetValidatorSet(ctx)

	for _, pubKeyr := range pubKeyrs {
		valSet.Remove(pubKeyr)
	}

	bz, err := k.cdc.MarshalBinaryLengthPrefixed(valSet)
	if err != nil {
		panic(err)
	}

	valSetStore.Set(key, bz)
}

//IsValidatorSet check validator is in authorised set
func (k *Keeper) IsValidator(ctx sdkTypes.Context, pubKeyr crypto.PubKey) bool {
	valSet := k.GetValidatorSet(ctx)

	_, ok := valSet.Contains(pubKeyr)
	return ok
}
