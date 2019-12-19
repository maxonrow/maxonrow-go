package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	nameservicecmd "github.com/maxonrow/maxonrow-go/x/nameservice/client/cli"
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
		Use:   "nameservice",
		Short: "Querying commands for the nameservice module",
	}

	queryCmd.AddCommand(client.GetCommands(
		nameservicecmd.GetCmdResolveName(mc.storeKey, mc.cdc),
		nameservicecmd.GetCmdWhois(mc.storeKey, mc.cdc),
	)...)

	return queryCmd
}

func (mc ModuleClient) GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "nameservice",
		Short: "Nameservice transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		// nameservicecmd.GetCmdSetAlias(mc.cdc),
		// nameservicecmd.GetCmdRemoveAlias(mc.cdc),
	)...)

	return txCmd
}
