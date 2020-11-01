package cli

import (
	"fmt"

	"github.com/maxonrow/maxonrow-go/x/fee"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetSysFeeSetting(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "name [fee-setting-name]",
		Short: "get system fee setting",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			name := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", "fee", fee.QuerySysFeeSetting, name), nil)
			if err != nil {
				fmt.Printf("Could not get fee setting: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}

	return cmd
}

func GetMsgFeeSetting(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "msgType [message-type]",
		Short: "get message fee setting by message type",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			msgType := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", "fee", fee.QueryMsgFeeSetting, msgType), nil)
			if err != nil {
				fmt.Printf("Could not get fee setting: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}

	return cmd

}

func GetFungibleTokenFeeSetting(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fungible-token [token-symbol] [action]",
		Short: "get fee setting by fungible-token symbol and action",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			symbol := args[0]
			action := args[1]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s", "fee", fee.QueryFungibleTokenFeeSetting, symbol, action), nil)
			if err != nil {
				fmt.Printf("Could not get token fee setting: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}

	return cmd
}

func GetNonFungibleTokenFeeSetting(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nonFungible-token [token-symbol] [action]",
		Short: "get fee setting by nonFungible-token symbol and action",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			symbol := args[0]
			action := args[1]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s", "fee", fee.QueryNonFungibleTokenFeeSetting, symbol, action), nil)
			if err != nil {
				fmt.Printf("Could not get token fee setting: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}

	return cmd
}

func GetAccFeeSetting(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [account-address]",
		Short: "get fee setting by account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			accStr := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", "fee", fee.QueryAccFeeSetting, accStr), nil)
			if err != nil {
				fmt.Printf("Could not get fee setting: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}

	return cmd
}

func GetFeeMultiplier(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multiplier",
		Short: "get multiplier",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", "fee", fee.QueryFeeMultiplier), nil)
			if err != nil {
				fmt.Printf("Could not get fee multiplier: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}

	return cmd
}

func GetFungibleTokenFeeMultiplier(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fungible-token-fee-multiplier",
		Short: "get fungible token fee multiplier",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", "fee", fee.QueryFungibleTokenFeeMultiplier), nil)
			if err != nil {
				fmt.Printf("Could not get fungible-token fee multiplier: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}

	return cmd
}

func GetNonFungibleTokenFeeMultiplier(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nonFungible-token-fee-multiplier",
		Short: "get nonFungible token fee multiplier",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", "fee", fee.QueryNonFungibleTokenFeeMultiplier), nil)
			if err != nil {
				fmt.Printf("Could not get nonFungible-token fee multiplier: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}

	return cmd
}

func GetNonFungibleTokenFeeCollector(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nonFungible-token-fee-collector",
		Short: "get nonFungible token fee multiplier",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", "fee", fee.QueryNonFungibleTokenFeeCollector), nil)
			if err != nil {
				fmt.Printf("Could not get nonFungible-token fee collector: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}

	return cmd
}
