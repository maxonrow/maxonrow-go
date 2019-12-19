package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	token "github.com/maxonrow/maxonrow-go/x/token/fungible"
)

func CreateFungibleTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-token [token-symbol]",
		Short: "creates new fungible token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			tokenSymbol := args[0]
			owner := cliCtx.GetFromAddress()

			metadata := viper.GetString("metadata")
			fixedSupply := viper.GetBool("fixed-supply")
			totalSupplyStr := viper.GetString("total-supply")
			decimalsStr := viper.GetString("decimals")
			name := viper.GetString("token-name")
			payFeeTo := viper.GetString("pay-fee-to")
			feeValue := viper.GetString("fee-value")
			totalSupply := sdkTypes.NewUintFromString(totalSupplyStr)
			payFeeToAddr, payFeeToAddrErr := sdkTypes.AccAddressFromBech32(payFeeTo)
			if payFeeToAddrErr != nil {
				return payFeeToAddrErr
			}

			tokenFee := token.NewFee(payFeeToAddr, feeValue)

			decimals, decErr := strconv.Atoi(decimalsStr)
			if decErr != nil {
				return decErr
			}

			msg := token.NewMsgCreateFungibleToken(tokenSymbol, decimals, owner, name, fixedSupply, totalSupply, metadata, tokenFee)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}

	cmd.Flags().String("token-name", "", "Desired token name")
	cmd.Flags().String("metadata", "", "IPFS hash link to attach to this process")
	cmd.Flags().Bool("fixed-supply", false, "To set the token fixed supply")
	cmd.Flags().String("total-supply", "0", "Total supply in case the supply is fixed")
	cmd.Flags().String("decimals", "8", "Decimals places")
	cmd.Flags().String("pay-fee-to", "mxw1p8qrka5ua840quqa3a3yzae5k25wpssq9n7890", "Wallet address")
	cmd.Flags().String("fee-value", "1000000000cin", "Fee amount")

	return cmd
}

func TransferFungibleTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-fungible [symbol]",
		Short: "transfer fungible token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()

			toString := viper.GetString("to")
			to, err := sdkTypes.AccAddressFromBech32(toString)
			if err != nil {
				return err
			}

			amountStr := viper.GetString("amount")
			amount := sdkTypes.NewUintFromString(amountStr)

			msg := token.NewMsgTransferFungibleToken(args[0], amount, from, to)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}

	cmd.Flags().String("to", "", "Address to which to transfer the fungible tokens to")
	cmd.Flags().String("amount", "1", "Amount of fungible token to transfer")

	return cmd
}

func BurnFungibleTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-burn [symbol] [value]",
		Short: "Request for burning the preowned token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			tokenSymbol := args[0]
			totalToken := args[2]
			totalval := sdkTypes.NewUintFromString(totalToken)
			owner := cliCtx.GetFromAddress()
			msg := token.NewMsgBurnFungibleToken(tokenSymbol, totalval, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}
	cmd.Flags().String("value", "", "amount of token wish to burn")
	return cmd
}

func TransferFungibleTokenOwnership(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-changeownership [symbol]",
		Short: "change the token ownership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			tokenSymbol := args[0]

			owner := cliCtx.GetFromAddress()

			toString := viper.GetString("to")
			to, err := sdkTypes.AccAddressFromBech32(toString)
			if err != nil {
				return err
			}

			msg := token.NewMsgTransferFungibleTokenOwnership(tokenSymbol, owner, to)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}
	cmd.Flags().String("to", "", "Address to which to transfer the ownership to")
	return cmd
}

func MintFungibleToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-minting [symbol]",
		Short: "change the token ownership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			tokenSymbol := args[0]
			owner := cliCtx.GetFromAddress()
			toString := viper.GetString("to")
			to, err := sdkTypes.AccAddressFromBech32(toString)
			if err != nil {
				return err
			}

			totalmint := viper.GetString("value")
			val := sdkTypes.NewUintFromString(totalmint)

			msg := token.NewMsgIssueFungibleAsset(owner, tokenSymbol, to, val)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}
	cmd.Flags().String("to", "", "Address to which to transfer the ownership to")
	cmd.Flags().String("value", "", "amount of token wish to mint")
	return cmd

}

// func ApproveTokenCmd(cdc *codec.Codec) *cobra.Command {
// 	cmd := &cobra.Command{
// 	// 	Use:   "token-approve [symbol]",
// 	// 	Short: "change the token ownership",
// 	// 	Args:  cobra.ExactArgs(1),
// 	// 	RunE: func(cmd *cobra.Command, args []string) error {
// 	// 		cliCtx := context.NewCLIContext().WithCodec(cdc)
// 	// 		txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

// 	// 		tokenSymbol := args[0]
// 	// 		owner := cliCtx.GetFromAddress()

// 	// 	   var pl token.Payload

// 	// 		// toString := viper.GetString("to")
// 	// 		// to, err := sdkTypes.AccAddressFromBech32(toString)
// 	// 		// if err != nil {
// 	// 		// 	return err
// 	// 		// }

// 	// 		totalmint := viper.GetString("value")
// 	// 		val := sdkTypes.NewUintFromString(totalmint)

// 	// 		msg := token.NewMsgSetFungibleTokenStatus(owner, tokenSymbol, to, val)
// 	// 		if err := msg.ValidateBasic(); err != nil {
// 	// 			return err
// 	// 		}
// 	// 		return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
// 	// 	},
// 	}
// 	// cmd.Flags().String("to", "", "Address to which to transfer the ownership to")
// 	// cmd.Flags().String("value", "", "amount of token wish to mint")
// 	return cmd

// }
