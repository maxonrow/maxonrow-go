package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	feeCmd "github.com/maxonrow/maxonrow-go/x/fee/client/cli"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "fee",
		Short: "Querying commands for the fee module",
	}

	queryCmd.AddCommand(client.GetCommands(
		feeCmd.GetSysFeeSetting(mc.cdc),
		feeCmd.GetMsgFeeSetting(mc.cdc),
		feeCmd.GetFungibleTokenFeeSetting(mc.cdc),
		feeCmd.GetNonFungibleTokenFeeSetting(mc.cdc),
		feeCmd.GetFeeMultiplier(mc.cdc),
		feeCmd.GetFungibleTokenFeeMultiplier(mc.cdc),
		feeCmd.GetNonFungibleTokenFeeMultiplier(mc.cdc),
		feeCmd.GetAccFeeSetting(mc.cdc),
		feeCmd.GetNonFungibleTokenFeeCollector(mc.cdc),
	)...)

	return queryCmd
}

func (mc ModuleClient) GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "fee",
		Short: "fee transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		feeCmd.CreateMsgSysFeeSetting(mc.cdc),
		feeCmd.EditMsgSysFeeSetting(mc.cdc),
		feeCmd.DeleteMsgSysFeeSetting(mc.cdc),
		feeCmd.CreateMsgAssignFeeToMsg(mc.cdc),
		feeCmd.CreateMsgAssignFeeToAcc(mc.cdc),
		feeCmd.CreateFeeMultiplier(mc.cdc),
		feeCmd.CreateFungibleTokenFeeMultiplier(mc.cdc),
		feeCmd.CreateNonFungibleTokenFeeMultiplier(mc.cdc),
		feeCmd.SetFungibleTokenFeeSetting(mc.cdc),
		feeCmd.SetNonFungibleTokenFeeSetting(mc.cdc),
		feeCmd.CreateMsgDeleteAccountFeeSetting(mc.cdc),
		//feeCmd.AddSysFeeSetting(mc.cdc),
	)...)

	return txCmd
}
