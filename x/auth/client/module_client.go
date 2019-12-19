package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	multiSigCmd "github.com/maxonrow/maxonrow-go/x/auth/client/cli"
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
		Use:   "multisig",
		Short: "Querying commands for the multisig module",
	}

	queryCmd.AddCommand(client.GetCommands(
		multiSigCmd.GetMultiSigAcc(mc.cdc),
	)...)

	return queryCmd
}

func (mc ModuleClient) GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "multisig",
		Short: "multisig transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		multiSigCmd.CreateMultiSigAccountCmd(mc.cdc),
		multiSigCmd.UpdateMultiSigAccountCmd(mc.cdc),
		multiSigCmd.TransferMultiSigAccountOwnershipCmd(mc.cdc),
	)...)

	return txCmd
}
