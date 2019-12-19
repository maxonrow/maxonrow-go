package fee

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/types"
)

type GenesisState struct {
	FeeSettings            []GenesisFeeSetting   `json:"fee_settings"`
	AuthorisedAddresses    []sdkTypes.AccAddress `json:"authorised_addresses"`
	Multiplier             string                `json:"multiplier"`
	AssignedMsgFeeSettings []AssignMsgFeeSetting `json:"assigned_fee"`
}

type AssignMsgFeeSetting struct {
	Name    string `json:"name"`
	MsgType string `json:"msg_type"`
}

type GenesisFeeSetting struct {
	Name       string         `json:"name"`
	Min        sdkTypes.Coins `json:"min"`
	Max        sdkTypes.Coins `json:"max"`
	Percentage string         `json:"percentage"`
}

func DefaultGenesisState() GenesisState {
	min, _ := sdkTypes.NewIntFromString("10000000000000000")
	max, _ := sdkTypes.NewIntFromString("1000000000000000000000000")
	defaultGenesisFeeSetting := GenesisFeeSetting{
		Name: "default",
		Min: []sdkTypes.Coin{
			sdkTypes.Coin{
				Denom:  types.CIN,
				Amount: min,
			},
		},
		Max: []sdkTypes.Coin{
			sdkTypes.Coin{
				Denom:  types.CIN,
				Amount: max,
			},
		},
		Percentage: "0.05",
	}
	return GenesisState{
		FeeSettings: []GenesisFeeSetting{defaultGenesisFeeSetting},
		Multiplier:  "1",
	}
}

func InitGenesis(ctx sdkTypes.Context, keeper *Keeper, genesisState GenesisState) {

	for _, feeSetting := range genesisState.FeeSettings {
		sysFee := NewMsgSysFeeSetting(feeSetting.Name, feeSetting.Min, feeSetting.Max, feeSetting.Percentage, genesisState.AuthorisedAddresses[0])
		keeper.storeFeeSetting(ctx, sysFee)
	}

	for _, assignMsgFeeSetting := range genesisState.AssignedMsgFeeSettings {
		msgAssignFeeToMsg := NewMsgAssignFeeToMsg(assignMsgFeeSetting.Name, assignMsgFeeSetting.MsgType, genesisState.AuthorisedAddresses[0])
		keeper.assignFeeToMsg(ctx, msgAssignFeeToMsg)
	}

	keeper.storeFeeMultiplier(ctx, genesisState.Multiplier)
	keeper.SetAuthorisedAddresses(ctx, genesisState.AuthorisedAddresses)

}

func ExportGenesis(keeper *Keeper) GenesisState {
	return GenesisState{}
}
