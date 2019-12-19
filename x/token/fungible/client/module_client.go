package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	tokenCmd "github.com/maxonrow/maxonrow-go/x/token/fungible/client/cli"
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
		Use:   "token",
		Short: "Querying commands for the token module",
	}

	queryCmd.AddCommand(client.GetCommands(
		tokenCmd.ListTokenSymbolCmd(mc.storeKey, mc.cdc),
		tokenCmd.GetTokenDataCmd(mc.storeKey, mc.cdc),
		tokenCmd.GetAccountCmd(mc.storeKey, mc.cdc),
		//assetcmd.GetNonfungibleAssetCmd(mc.storeKey, mc.cdc),
	)...)

	return queryCmd
}

func (mc ModuleClient) GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "token",
		Short: "token transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		tokenCmd.CreateFungibleTokenCmd(mc.cdc),
		tokenCmd.MintFungibleToken(mc.cdc),
		tokenCmd.TransferFungibleTokenCmd(mc.cdc),
		tokenCmd.TransferFungibleTokenOwnership(mc.cdc),
		tokenCmd.BurnFungibleTokenCmd(mc.cdc),
		// tokenCmd.ApproveTokenCmd(mc.cdc),
		// tokenCmd.RejectAssetClassCmd(mc.cdc),
		// tokenCmd.FreezeAssetClassCmd(mc.cdc),
		// tokenCmd.UnfreezeAssetClassCmd(mc.cdc),
		// tokenCmd.IssueFungibleAssetCmd(mc.cdc),
		// tokenCmd.FreezeFungibleAccountCmd(mc.cdc),
		// tokenCmd.UnfreezeFungibleAccountCmd(mc.cdc),

	)...)

	return txCmd
}
