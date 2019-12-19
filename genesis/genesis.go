package genesis

import (
	"encoding/json"
	"fmt"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	sdkBank "github.com/cosmos/cosmos-sdk/x/bank"
	sdkDist "github.com/cosmos/cosmos-sdk/x/distribution"
	sdkStaking "github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/fee"
	"github.com/maxonrow/maxonrow-go/x/kyc"
	"github.com/maxonrow/maxonrow-go/x/maintenance"
	"github.com/maxonrow/maxonrow-go/x/nameservice"
	token "github.com/maxonrow/maxonrow-go/x/token/fungible"
)

// Represents the state at the start. Should store initial accounts etc here
type GenesisState struct {
	AuthState        sdkAuth.GenesisState     `json:"auth"`
	Accounts         []*sdkAuth.BaseAccount   `json:"accounts"`
	BankState        sdkBank.GenesisState     `json:"bank"`
	StakingState     sdkStaking.GenesisState  `json:"staking"`
	DistrState       sdkDist.GenesisState     `json:"distribution"`
	KycState         kyc.GenesisState         `json:"kyc"`
	TokenState       token.GenesisState       `json:"token"`
	NameServiceState nameservice.GenesisState `json:"nameservice"`
	FeeState         fee.GenesisState         `json:"fee"`
	MaintenanceState maintenance.GenesisState `json:"maintenance"`
	GenTxs           []json.RawMessage        `json:"gentxs"`
}

// NewDefaultGenesisState generates the default state.
func NewDefaultGenesisState() GenesisState {
	gen := GenesisState{
		AuthState:        sdkAuth.DefaultGenesisState(),
		Accounts:         nil,
		BankState:        sdkBank.DefaultGenesisState(),
		StakingState:     sdkStaking.DefaultGenesisState(),
		DistrState:       sdkDist.DefaultGenesisState(),
		KycState:         kyc.DefaultGenesisState(),
		TokenState:       token.DefaultGenesisState(),
		NameServiceState: nameservice.DefaultGenesisState(),
		FeeState:         fee.DefaultGenesisState(),
		MaintenanceState: maintenance.DefaultGenesisState(),
		GenTxs:           nil,
	}

	gen.AuthState.Params.TxSizeCostPerByte = uint64(0)
	gen.AuthState.Params.SigVerifyCostED25519 = uint64(0)
	gen.AuthState.Params.SigVerifyCostSecp256k1 = uint64(0)
	gen.StakingState.Params.BondDenom = types.CIN
	gen.DistrState.CommunityTax = sdkTypes.NewDec(0)
	gen.DistrState.BaseProposerReward = sdkTypes.NewDec(0)
	gen.DistrState.BonusProposerReward = sdkTypes.NewDec(0)

	return gen
}

func (s *GenesisState) Validate() error {
	// https://github.com/maxonrow/maxonrow-go/issues/44
	//if err := sdkAuth.ValidateGenesis(s.AuthState); err != nil {
	//	panic(fmt.Sprintf("Invalid genesis auth state: %s", err))
	//}

	if err := sdkBank.ValidateGenesis(s.BankState); err != nil {
		return fmt.Errorf("Invalid genesis bank state: %s", err)
	}
	if err := sdkStaking.ValidateGenesis(s.StakingState); err != nil {
		return fmt.Errorf("Invalid genesis stake state: %s", err)
	}
	if err := sdkDist.ValidateGenesis(s.DistrState); err != nil {
		return fmt.Errorf("Invalid genesis distribution state: %s", err)
	}

	if s.AuthState.Params.TxSizeCostPerByte != uint64(0) {
		return fmt.Errorf("TxSizeCostPerByte should be zero: %v", s.AuthState.Params.TxSizeCostPerByte)
	}
	if s.AuthState.Params.SigVerifyCostED25519 != uint64(0) {
		return fmt.Errorf("SigVerifyCostED25519 should be zero: %v", s.AuthState.Params.SigVerifyCostED25519)
	}
	if s.AuthState.Params.SigVerifyCostSecp256k1 != uint64(0) {
		return fmt.Errorf("SigVerifyCostSecp256k1 should be zero: %v", s.AuthState.Params.SigVerifyCostSecp256k1)
	}
	if s.StakingState.Params.BondDenom != types.CIN {
		return fmt.Errorf("BondDenom should be cin: %v", s.StakingState.Params.BondDenom)
	}
	if !s.DistrState.CommunityTax.Equal(sdkTypes.NewDec(0)) {
		return fmt.Errorf("CommunityTax should be zero: %v", s.DistrState.CommunityTax)
	}
	if !s.DistrState.BaseProposerReward.Equal(sdkTypes.NewDec(0)) {
		return fmt.Errorf("BaseProposerReward should be zero: %v", s.DistrState.BaseProposerReward)
	}
	if !s.DistrState.BonusProposerReward.Equal(sdkTypes.NewDec(0)) {
		return fmt.Errorf("BonusProposerReward should be zero. %v", s.DistrState.BonusProposerReward)
	}

	/*

	   TODO::::
	   		kycValidateGenesis(s.DistrState);
	   		tokenValidateGenesis(s.DistrState);
	   		nameserviceValidateGenesis(s.DistrState);
	   		feeValidateGenesis(s.DistrState);
	   		maintenanceValidateGenesis(s.DistrState);

	*/

	// check for duplicated accounts in genesis
	accounts := make(map[string]bool)

	for _, initialAccount := range s.Accounts {
		addrStr := initialAccount.Address.String()
		if _, exists := accounts[addrStr]; exists {
			return fmt.Errorf("Duplicate account %s in genesis", addrStr)
		}

		accounts[addrStr] = true
	}

	return nil
}
