package cli

import (
	"bufio"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/maxonrow/maxonrow-go/x/auth"
)

const (
	flagMultisig          = "multisig"
	flagMultiSigThreshold = "multisig-threshold"
)

func CreateMultiSigAccountCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-multisig-account",
		Short: "Create multi signature account",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := sdkAuth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			master := cliCtx.GetFromAddress()
			multisigKeys := viper.GetStringSlice(flagMultisig)
			multisigThreshold := viper.GetInt(flagMultiSigThreshold)

			kb, err := keys.NewKeyringFromHomeFlag(cmd.InOrStdin())
			if err != nil {
				return err
			}

			var signers []sdkTypes.AccAddress
			if len(multisigKeys) != 0 {
				for _, keyname := range multisigKeys {
					k, err := kb.Get(keyname)
					if err != nil {
						return err
					}
					signers = append(signers, k.GetAddress())
				}
			}

			msg := auth.NewMsgCreateMultiSigAccount(master, multisigThreshold, signers)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdkTypes.Msg{msg})
		},
	}

	cmd.Flags().StringSlice(flagMultisig, nil, "list of signers eg. "+strconv.Quote("acc1,acc2,acc3")+" by local wallet names.")
	cmd.Flags().Uint(flagMultiSigThreshold, 1, "K out of N required signatures. For use in conjunction with --multisig")

	//cmd = client.PostCommands(cmd)[0]

	return cmd
}

func UpdateMultiSigAccountCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-multisig-account [group address]",
		Short: "Update multi signature account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := sdkAuth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			owner := cliCtx.GetFromAddress()
			multisigKeys := viper.GetStringSlice(flagMultisig)
			multisigThreshold := viper.GetInt(flagMultiSigThreshold)

			groupAddress, err := sdkTypes.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			kb, err := keys.NewKeyringFromHomeFlag(cmd.InOrStdin())
			if err != nil {
				return err
			}

			var signers []sdkTypes.AccAddress
			if len(multisigKeys) != 0 {
				for _, keyname := range multisigKeys {
					k, err := kb.Get(keyname)
					if err != nil {
						return err
					}
					signers = append(signers, k.GetAddress())
				}
			}

			msg := auth.NewMsgUpdateMultiSigAccount(owner, groupAddress, multisigThreshold, signers)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdkTypes.Msg{msg})
		},
	}

	cmd.Flags().StringSlice(flagMultisig, nil, "list of signers eg. "+strconv.Quote("acc1,acc2,acc3")+" by local wallet names.")
	cmd.Flags().Uint(flagMultiSigThreshold, 1, "K out of N required signatures. For use in conjunction with --multisig")

	//cmd = client.PostCommands(cmd)[0]

	return cmd
}

func TransferMultiSigAccountOwnershipCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-multisig-account-ownership [from owner] [to new owner] [group address]",
		Short: "Transfer multi signature account ownership",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := sdkAuth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to, err := sdkTypes.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			groupAddress, err := sdkTypes.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			// TO-DO: Checking for group address, from address and to address

			msg := auth.NewMsgTransferMultiSigOwner(groupAddress, to, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdkTypes.Msg{msg})
		},
	}

	return cmd
}
