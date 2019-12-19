package cli

import (
	"fmt"
	"github.com/tendermint/tendermint/libs/bech32"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetCmdResolveName(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "resolve [alias]",
		Short: "resolve alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			alias := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/resolve/%s", queryRoute, alias), nil)
			if err != nil {
				fmt.Printf("could not resolve alias - %s \n", string(alias))
				return nil
			}

			addrString, _ := bech32.ConvertAndEncode("mxw", res)

			fmt.Println(addrString)

			return nil
		},
	}
}

func GetCmdWhois(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "whois [address]",
		Short: "Query whois info of address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			address := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/whois/%s", queryRoute, address), nil)
			if err != nil {
				fmt.Printf("could not resolve whois - %s \n", string(address))
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
