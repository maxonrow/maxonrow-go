package cli

import (
	"fmt"
	"strconv"

	"github.com/maxonrow/maxonrow-go/x/token/fungible"
	"github.com/maxonrow/maxonrow-go/x/token/nonfungible"

	"github.com/maxonrow/maxonrow-go/x/fee"
	"github.com/maxonrow/maxonrow-go/x/kyc"
	"github.com/maxonrow/maxonrow-go/x/maintenance"
	"github.com/maxonrow/maxonrow-go/x/nameservice"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetCmdGetProposal(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proposal [proposalID]",
		Short: "get proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposalID := args[0]

			proposalIDInt, err := strconv.ParseUint(proposalID, 10, 64)

			params := maintenance.NewQueryProposalParams(proposalIDInt)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/proposal", queryRoute), bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func GetCmdGetKycMaintainerAddresses(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "kyc",
		Short: "query kyc maintenance parties address for Issuer, Provider, Middleware",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, kyc.QueryGetKycMaintainerAddresses), nil)
			if err != nil {
				fmt.Printf("Could not get kyc maintenance parties addresses: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func GetCmdGetFungibleTokenMaintainerAddresses(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "token",
		Short: "query fungible token maintenance parties address for Issuer, Provider, Middleware",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, fungible.QueryGetFungibleTokenMaintainerAddresses), nil)
			if err != nil {
				fmt.Printf("Could not get fungible token maintenance parties addresses: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func GetCmdGetNonfungibleTokenMaintainerAddresses(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "nonfungible-token",
		Short: "query nonfungible-token maintenance parties address for Issuer, Provider, Middleware",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, nonfungible.QueryGetNonfungibleTokenMaintainerAddresses), nil)
			if err != nil {
				fmt.Printf("Could not get nonfungible token maintenance parties addresses: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func GetCmdGetNameserviceMaintainerAddresses(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "nameservice",
		Short: "query nameservice maintenance parties address for Issuer, Provider, Middleware",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, nameservice.QueryGetNameserviceMaintainerAddresses), nil)
			if err != nil {
				fmt.Printf("Could not get nameservice maintenance parties addresses: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func GetCmdGetFeeMaintainerAddresses(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "fee",
		Short: "query fee maintenance parties address for Middleware, Fee-collector",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, fee.QueryGetFeeMaintainerAddresses), nil)
			if err != nil {
				fmt.Printf("Could not get fee maintenance parties addresses: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
