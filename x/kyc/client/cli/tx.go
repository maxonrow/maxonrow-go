package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetCmdWhitelist(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist [address] [kyc address]",
		Short: "whitelist address from authorised address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// cliCtx := context.NewCLIContext().WithCodec(cdc)

			// txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// signer := cliCtx.GetFromName()

			// recipient, err := sdkTypes.AccAddressFromBech32(args[0])
			// if err != nil {
			// 	return err
			// }
			// kycAddress := args[1]
			// walletSignature, walletErr := sdkTypes.AccAddressFromBech32(args[2])
			// if walletErr != nil {
			// 	return walletErr
			// }
			// kyc := kyc.NewKyc(recipient, "0", kycAddress)
			//* fix populating data
			// kycPayload := kyc.NewPayload(recipient, "", kycAddress)
			// kycData := kyc.NewKycData(kycPayload, []byte{walletSignature}, crypto.PubKey{nil})
			// msg := kyc.NewMsgWhitelist(sender, kycData, []kyc.Signature{nil})
			// if err := msg.ValidateBasic(); err != nil {
			// 	return err
			// }

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{})
			return nil
		},
	}

	cmd.Flags().String("issuer-address", "", "Issuer address")
	cmd.Flags().String("provider-address", "", "Provider Address")

	return cmd
}
