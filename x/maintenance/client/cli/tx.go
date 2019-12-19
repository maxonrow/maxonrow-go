package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto"
	"github.com/maxonrow/maxonrow-go/x/maintenance"
)

// GetCmdSubmitProposal is the CLI command for sending a BuyName transaction
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-proposal [proposal name]",
		Short: "submit proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			proposer := cliCtx.GetFromAddress()

			proposalType := viper.GetString("proposal-type")
			title := viper.GetString("title")
			description := viper.GetString("description")
			authorisedAddrStr := viper.GetString("authorised-address")
			providerAddrStr := viper.GetString("provider-address")
			validatorPubKeyStr := viper.GetString("validator-pubkey")
			issuerAddrStr := viper.GetString("issuer-address")
			feeCollectorAddrStr := viper.GetString("fee-collector")
			feeCollectorModule := viper.GetString("fee-collector-module")
			action := viper.GetString("action")

			proposalKind, proposalKindErr := maintenance.ProposalTypeFromString(proposalType)
			if proposalKindErr != nil {
				return proposalKindErr
			}

			fmt.Println(proposalKind)

			var msg maintenance.MsgProposal
			// TO-DO: better implementation
			switch proposalKind {
			case maintenance.ProposalTypeModifyFee:

				if authorisedAddrStr == "" && feeCollectorAddrStr == "" {
					return sdkTypes.ErrInternal(fmt.Sprintf("Proposal type error, please check: %s, --proposalType %s", proposalKind.String(), proposalType))
				}

				var authorisedAddress sdkTypes.AccAddress
				if authorisedAddrStr != "" {
					authorisedAddr, authorisedAddrErr := sdkTypes.AccAddressFromBech32(authorisedAddrStr)
					if authorisedAddrErr != nil {
						return authorisedAddrErr
					}
					authorisedAddress = authorisedAddr
				} else {
					authorisedAddress = nil
				}
				feeCollectorAddr, feeCollectorAddrErr := sdkTypes.AccAddressFromBech32(feeCollectorAddrStr)
				if feeCollectorAddrErr != nil {
					return feeCollectorAddrErr
				}

				feeCollector := maintenance.FeeCollector{
					Module:  feeCollectorModule,
					Address: feeCollectorAddr,
				}
				feeMaintainer := maintenance.NewFeeeMaintainer(action, []sdkTypes.AccAddress{authorisedAddress}, []maintenance.FeeCollector{feeCollector})

				msg = maintenance.NewMsgSubmitProposal(title, description, proposalKind, &feeMaintainer, proposer)

			case maintenance.ProposalTypeModifyKyc:

				if authorisedAddrStr == "" && issuerAddrStr == "" && providerAddrStr == "" {
					return sdkTypes.ErrInternal(fmt.Sprintf("Proposal type error, please check: %s, --proposalType %s", proposalKind.String(), proposalType))
				}

				authorisedAddr, authorisedAddrErr := sdkTypes.AccAddressFromBech32(authorisedAddrStr)
				if authorisedAddrErr != nil {
					return authorisedAddrErr
				}
				providerAddr, providerAddrErr := sdkTypes.AccAddressFromBech32(providerAddrStr)
				if providerAddrErr != nil {
					return providerAddrErr
				}
				issuerAddr, issuerAddrErr := sdkTypes.AccAddressFromBech32(issuerAddrStr)
				if issuerAddrErr != nil {
					return issuerAddrErr
				}
				kycMaintainer := maintenance.NewKycMaintainer(action, []sdkTypes.AccAddress{authorisedAddr}, []sdkTypes.AccAddress{issuerAddr}, []sdkTypes.AccAddress{providerAddr})

				msg = maintenance.NewMsgSubmitProposal(title, description, proposalKind, &kycMaintainer, proposer)
			case maintenance.ProposalTypeModifyNameservice:

				if authorisedAddrStr == "" && issuerAddrStr == "" && providerAddrStr == "" {
					return sdkTypes.ErrInternal(fmt.Sprintf("Proposal type error, please check: %s, --proposalType %s", proposalKind.String(), proposalType))
				}

				authorisedAddr, authorisedAddrErr := sdkTypes.AccAddressFromBech32(authorisedAddrStr)
				if authorisedAddrErr != nil {
					return authorisedAddrErr
				}
				providerAddr, providerAddrErr := sdkTypes.AccAddressFromBech32(providerAddrStr)
				if providerAddrErr != nil {
					return providerAddrErr
				}
				issuerAddr, issuerAddrErr := sdkTypes.AccAddressFromBech32(issuerAddrStr)
				if issuerAddrErr != nil {
					return issuerAddrErr
				}
				nameserviceMaintainer := maintenance.NewNamerserviceMaintainer(action, []sdkTypes.AccAddress{authorisedAddr}, []sdkTypes.AccAddress{issuerAddr}, []sdkTypes.AccAddress{providerAddr})

				msg = maintenance.NewMsgSubmitProposal(title, description, proposalKind, &nameserviceMaintainer, proposer)
			case maintenance.ProposalTypeModifyToken:

				if authorisedAddrStr == "" && issuerAddrStr == "" && providerAddrStr == "" {
					return sdkTypes.ErrInternal(fmt.Sprintf("Proposal type error, please check: %s, --proposalType %s", proposalKind.String(), proposalType))
				}

				authorisedAddr, authorisedAddrErr := sdkTypes.AccAddressFromBech32(authorisedAddrStr)
				if authorisedAddrErr != nil {
					return authorisedAddrErr
				}
				providerAddr, providerAddrErr := sdkTypes.AccAddressFromBech32(providerAddrStr)
				if providerAddrErr != nil {
					return providerAddrErr
				}
				issuerAddr, issuerAddrErr := sdkTypes.AccAddressFromBech32(issuerAddrStr)
				if issuerAddrErr != nil {
					return issuerAddrErr
				}
				tokenMaintainer := maintenance.NewTokenMaintainer(action, []sdkTypes.AccAddress{authorisedAddr}, []sdkTypes.AccAddress{issuerAddr}, []sdkTypes.AccAddress{providerAddr})

				msg = maintenance.NewMsgSubmitProposal(title, description, proposalKind, &tokenMaintainer, proposer)
			case maintenance.ProposalTypesModifyValidatorSet:

				if validatorPubKeyStr == "" {
					return sdkTypes.ErrInternal(fmt.Sprintf("Proposal type error, please check: %s, --proposalType %s", proposalKind.String(), proposalType))
				}

				pubKeyAddr, pubKeyErr := sdkTypes.GetConsPubKeyBech32(validatorPubKeyStr)
				if pubKeyErr != nil {
					return pubKeyErr
				}
				whitelistValidator := maintenance.NewWhitelistValidator(action, []crypto.PubKey{pubKeyAddr})
				msg = maintenance.NewMsgSubmitProposal(title, description, proposalKind, &whitelistValidator, proposer)
			default:
				return sdkTypes.ErrInternal("Unregonised proposal type.")
			}

			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}

	cmd.Flags().String("proposal-type", "token", "Module name")
	cmd.Flags().String("title", "Modify kyc maintainer", "Proposal title")
	cmd.Flags().String("description", "adding new kyc maintainer", "Proposal description")
	cmd.Flags().String("authorised-address", "", "Address to be added/removed as authorised address.")
	cmd.Flags().String("issuer-address", "", "Address to be added/removed as issuer address.")
	cmd.Flags().String("provider-address", "", "Address to be added/removed as provider address.")
	cmd.Flags().String("fee-collector", "", "Address to be added/removed as fee collector.")
	cmd.Flags().String("fee-collector-module", "", "Fee collector has to assign to collect/removed fees for a module.")
	cmd.Flags().String("action", "add", "Action can be remove or add.")
	cmd.Flags().String("validator-address", "", "Validator address to be whitelisted or revoke.")
	return cmd
}

func GetCmdCastAction(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cast-action [proposalID]",
		Short: "Approve/ reject proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			caster := cliCtx.GetFromAddress()
			proposalID := args[0]
			proposalIDInt, err := strconv.ParseUint(proposalID, 10, 64)
			if err != nil {
				return err
			}

			action := viper.GetString("action")
			msg := maintenance.NewMsgCastAction(caster, proposalIDInt, action)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}
	cmd.Flags().String("action", "", "Action can be approve or reject.")
	return cmd
}
