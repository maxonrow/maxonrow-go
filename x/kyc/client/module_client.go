package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	kyccli "github.com/maxonrow/maxonrow-go/x/kyc/client/cli"
)

type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "kyc",
		Short: "Kyc querying subcommands",
	}

	queryCmd.AddCommand(client.GetCommands(
		kyccli.GetCmdIsWhitelisted(mc.storeKey, mc.cdc),
	)...)

	return queryCmd
}

func (mc ModuleClient) GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "kyc",
		Short: "Kyc transaction subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		kyccli.GetCmdWhitelist(mc.cdc),
	)...)

	return txCmd
}
