package cli

// GetCmdBuyName is the CLI command for sending a BuyName transaction
// func GetCmdSetAlias(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "set-alias [name]",
// 		Short: "set alias",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)

// 			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

// 			if err := cliCtx.EnsureAccountExists(); err != nil {
// 				return err
// 			}

// 			account := cliCtx.GetFromAddress()

// 			msg := nameservice.NewMsgSetName(args[0], account)
// 			err := msg.ValidateBasic()
// 			if err != nil {
// 				return err
// 			}

//

// 			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
// 		},
// 	}
// }

// func GetCmdRemoveAlias(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "remove-alias",
// 		Short: "remove alias",
// 		Args:  cobra.ExactArgs(0),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)

// 			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

// 			if err := cliCtx.EnsureAccountExists(); err != nil {
// 				return err
// 			}

// 			account := cliCtx.GetFromAddress()

// 			msg := nameservice.NewMsgRemoveAlias(account)
// 			err := msg.ValidateBasic()
// 			if err != nil {
// 				return err
// 			}

//

// 			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
// 		},
// 	}
// }
