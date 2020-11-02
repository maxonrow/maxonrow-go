package cli

import (
	"bufio"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	token "github.com/maxonrow/maxonrow-go/x/token/nonfungible"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateNonFungibleTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-create [token-symbol]",
		Short: "creates new non fungible token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			properties := ""
			tokenSymbol := args[0]
			owner := cliCtx.GetFromAddress()
			metadata := viper.GetString("metadata")
			name := viper.GetString("token-name")
			payFeeTo := viper.GetString("pay-fee-to")
			feeValue := viper.GetString("fee-value")
			payFeeToAddr, payFeeToAddrErr := sdkTypes.AccAddressFromBech32(payFeeTo)
			if payFeeToAddrErr != nil {
				return payFeeToAddrErr
			}

			tokenFee := token.NewFee(payFeeToAddr, feeValue)

			msg := token.NewMsgCreateNonFungibleToken(tokenSymbol, owner, name, properties, metadata, tokenFee)
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
	cmd.Flags().String("pay-fee-to", "mxw1p8qrka5ua840quqa3a3yzae5k25wpssq9n7890", "Wallet address")
	cmd.Flags().String("fee-value", "1000000000cin", "Fee amount")

	return cmd
}

func TransferNonFungibleItem(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-transferitem [symbol] [itemid]",
		Short: "transfer nonfungible item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			from := cliCtx.GetFromAddress()
			toString := viper.GetString("to")
			to, err := sdkTypes.AccAddressFromBech32(toString)
			if err != nil {
				return err
			}
			itemID := args[1]

			msg := token.NewMsgTransferNonFungibleItem(args[0], from, to, itemID)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}

	cmd.Flags().String("to", "", "Address to which to transfer the non fungible tokens to")

	return cmd
}

func BurnNonFungibleItem(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-burnitem [symbol] [itemID]",
		Short: "Request for burning the preowned non fungible item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			tokenSymbol := args[0]
			itemID := args[1]
			owner := cliCtx.GetFromAddress()
			msg := token.NewMsgBurnNonFungibleItem(tokenSymbol, owner, itemID)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}
	return cmd
}

func TransferNonFungibleTokenOwnership(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-changeownership [symbol]",
		Short: "change the non fungible token ownership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			tokenSymbol := args[0]
			owner := cliCtx.GetFromAddress()
			toString := viper.GetString("to")
			to, err := sdkTypes.AccAddressFromBech32(toString)
			if err != nil {
				return err
			}

			msg := token.NewMsgTransferNonFungibleTokenOwnership(tokenSymbol, owner, to)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}
	cmd.Flags().String("to", "", "Address to which to transfer the ownership to")
	return cmd
}

func MintNonFungibleItem(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-mintitem [symbol] [itemId] [properties]",
		Short: "change the token ownership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			symbol := args[0]
			itemID := args[1]
			properties := args[2]

			metadata := viper.GetString("metadata")
			owner := cliCtx.GetFromAddress()

			toString := viper.GetString("to")
			to, err := sdkTypes.AccAddressFromBech32(toString)
			if err != nil {
				return err
			}

			msg := token.NewMsgMintNonFungibleItem(owner, symbol, to, itemID, properties, metadata)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}
	cmd.Flags().String("to", "", "Address to which to transfer the ownership to")
	cmd.Flags().String("metadata", "", "IPFS hash link to attach to this process")
	return cmd

}

func Endorsement(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft-endose [symbol] [itemId] [metadata]",
		Short: "endose the nft item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			symbol := args[0]
			itemID := args[1]
			owner := cliCtx.GetFromAddress()
			metadata := args[2]
			msg := token.NewMsgEndorsement(symbol, owner, itemID, metadata)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}
	return cmd

}
