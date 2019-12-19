package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	feeCmd "github.com/maxonrow/maxonrow-go/x/fee/client/cli"
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
		feeCmd.GetTokenFeeSetting(mc.cdc),
		feeCmd.GetFeeMultiplier(mc.cdc),
		feeCmd.GetTokenFeeMultiplier(mc.cdc),
		feeCmd.GetAccFeeSetting(mc.cdc),
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
		feeCmd.CreateTokenFeeMultiplier(mc.cdc),
		feeCmd.AddSysFeeSetting(mc.cdc),
	)...)

	return txCmd
}
