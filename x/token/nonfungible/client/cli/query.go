package cli

// import (
// 	"fmt"

// 	"github.com/cosmos/cosmos-sdk/client/context"
// 	"github.com/cosmos/cosmos-sdk/codec"
// 	"github.com/spf13/cobra"
// 	"gitlab.com/mxw.old/maxonrow-go/x/token"
// )

// func ListTokenSymbolCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "symbols",
// 		Short: "list all token symbols, fungible and nonfungible",
// 		Args:  cobra.ExactArgs(0),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)

// 			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, token.QueryListTokenSymbol), nil)
// 			if err != nil {
// 				fmt.Printf("Could not list token symbols: %s\n", err)
// 				return nil
// 			}

// 			fmt.Println(string(res))

// 			return nil
// 		},
// 	}
// }

// func GetTokenDataCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "token-data [token-symbol]",
// 		Short: "get a single token data, fungible or nonfungible",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)

// 			tokenSymbol := args[0]

// 			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, token.QueryTokenData, tokenSymbol), nil)
// 			if err != nil {
// 				fmt.Printf("Could not get asset class: %s\n", err)
// 				return nil
// 			}

// 			fmt.Println(string(res))

// 			return nil
// 		},
// 	}
// }

// func GetAccountCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "token-account [token-symbol] [account]",
// 		Short: "get information about token belonging to a single account for the given token symbol",
// 		Args:  cobra.ExactArgs(2),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)

// 			symbol := args[0]
// 			account := args[1]

// 			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s", queryRoute, token.QueryAccount, symbol, account), nil)
// 			if err != nil {
// 				fmt.Printf("Could not get account: %s\n", err)
// 				return nil
// 			}

// 			fmt.Println(string(res))

// 			return nil
// 		},
// 	}
// }
