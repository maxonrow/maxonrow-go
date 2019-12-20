package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/fee"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type FeeTransaction struct {
	FeeName       string
	Max           string
	Min           string
	IssuerAddress sdkTypes.AccAddress
	percentage    string
	Txfee         string
	GasPrice      string
}

func CreateMsgSysFeeSetting(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-sysfee",
		Short: "create new fee setting",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			issuer := cliCtx.GetFromAddress()

			name := viper.GetString("name")

			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/fee/is_fee_setting_exist/%s", name), nil)
			if err != nil {
				return err
			}

			if string(bz) == "true" {
				return fmt.Errorf("Fee setting name is taken.")
			}

			minStr := viper.GetString("min")
			maxStr := viper.GetString("max")
			percentageStr := viper.GetString("percentage")

			amtMin, ok := sdkTypes.NewIntFromString(minStr)
			if !ok {
				return fmt.Errorf("Invalid min amount.")
			}
			mincoin := sdkTypes.NewCoin(types.CIN, amtMin)

			amtMax, ok := sdkTypes.NewIntFromString(maxStr)
			if !ok {
				return fmt.Errorf("Invalid max amount.")
			}
			maxcoin := sdkTypes.NewCoin(types.CIN, amtMax)

			msg := fee.NewMsgSysFeeSetting(name, sdkTypes.Coins{mincoin}, sdkTypes.Coins{maxcoin}, percentageStr, issuer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}

	cmd.Flags().String("name", "default", "Fee setting name")
	cmd.Flags().String("min", "200000000cin", "minimum fee")
	cmd.Flags().String("max", "1000000000cin", "maximum fee")
	cmd.Flags().String("percentage", "0.05", "percentage example: 10% = 10, 0.1% = 0.1")

	return cmd
}

func EditMsgSysFeeSetting(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-sysfee",
		Short: "update fee setting",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			issuer := cliCtx.GetFromAddress()

			name := viper.GetString("name")

			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/fee/is_fee_setting_exist/%s", name), nil)
			if err != nil {
				return err
			}

			if string(bz) == "false" {
				return fmt.Errorf("Fee setting name is not exist.")
			}

			minStr := viper.GetString("min")
			maxStr := viper.GetString("max")
			percentageStr := viper.GetString("percentage")

			amtMin, ok := sdkTypes.NewIntFromString(minStr)
			if !ok {
				return fmt.Errorf("Invalid min amount.")
			}
			mincoin := sdkTypes.NewCoin(types.CIN, amtMin)

			amtMax, ok := sdkTypes.NewIntFromString(maxStr)
			if !ok {
				return fmt.Errorf("Invalid max amount.")
			}
			maxcoin := sdkTypes.NewCoin(types.CIN, amtMax)

			msg := fee.NewMsgSysFeeSetting(name, sdkTypes.Coins{mincoin}, sdkTypes.Coins{maxcoin}, percentageStr, issuer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}

	cmd.Flags().String("name", "default", "Fee setting name")
	cmd.Flags().String("min", "200000000cin", "minimum fee")
	cmd.Flags().String("max", "1000000000cin", "maximum fee")
	cmd.Flags().String("percentage", "0.05", "percentage example: 10% = 10, 0.1% = 0.1")

	return cmd
}

func DeleteMsgSysFeeSetting(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-sysfee [fee name]",
		Short: "delete fee setting",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			issuer := cliCtx.GetFromAddress()

			name := args[0]

			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/fee/is_fee_setting_in_used/%s", name), nil)
			if err != nil {
				return err
			}
			if string(bz) == "true" {
				return fmt.Errorf("Unable to delete, fee setting is in used.")
			}

			msg := fee.NewMsgDeleteSysFeeSetting(name, issuer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}
	return cmd
}

func CreateMsgAssignFeeToMsg(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "msg [msg-type]",
		Short: "Assign a fee setting to a message type",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msgType := args[0]
			issuer := cliCtx.GetFromAddress()

			name := viper.GetString("name")
			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/fee/is_fee_setting_exist/%s", name), nil)
			if err != nil {
				return err
			}

			if string(bz) == "false" {
				return fmt.Errorf("Fee setting name is not exist.")
			}

			msg := fee.NewMsgAssignFeeToMsg(name, msgType, issuer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}

	cmd.Flags().String("name", "default", "Fee setting name")

	return cmd
}

func SetTokenFeeSetting(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token [token-symbol] [action]",
		Short: "Assign a fee setting to a token(fungible/ nonfungible)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			tokenSymbol := args[0]
			tokenAction := args[1]

			issuer := cliCtx.GetFromAddress()
			name := viper.GetString("name")

			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/fee/is_fee_setting_exist/%s", name), nil)
			if err != nil {
				return err
			}
			if string(bz) == "false" {
				return fmt.Errorf("Fee setting name is not exist.")
			}

			tokenData, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/token/token_data/%s", tokenSymbol), nil)
			if tokenData == nil {
				return fmt.Errorf("No such token symbol.")
			}

			isValid, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/fee/is_token_action_valid/%s", tokenAction), nil)
			if err != nil {
				return err
			}

			if string(isValid) == "false" {
				return fmt.Errorf("Token action is not valid.")
			}

			msg := fee.NewMsgAssignFeeToToken(name, tokenSymbol, tokenAction, issuer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}

	cmd.Flags().String("name", "default", "Fee setting name")

	return cmd
}

func CreateMsgAssignFeeToAcc(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [account address]",
		Short: "Assign a fee setting to an account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			accStr := args[0]
			acc, accErr := sdkTypes.AccAddressFromBech32(accStr)
			if accErr != nil {
				return accErr
			}

			issuer := cliCtx.GetFromAddress()
			name := viper.GetString("name")

			msg := fee.NewMsgAssignFeeToAcc(name, acc, issuer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}

	cmd.Flags().String("name", "default", "Fee setting name")

	return cmd
}

func CreateMsgDeleteAccountFeeSetting(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-acc-fee [account address]",
		Short: "Delete account fee setting",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			accStr := args[0]
			acc, accErr := sdkTypes.AccAddressFromBech32(accStr)
			if accErr != nil {
				return accErr
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", "fee", fee.QueryAccFeeSetting, accStr), nil)
			if err != nil {
				return fmt.Errorf("Could not get fee setting: %s\n", err)
			}

			if string(res) == "null" {
				return fmt.Errorf("Account does not have any fee setting.")
			}

			issuer := cliCtx.GetFromAddress()

			msg := fee.NewMsgDeleteAccFeeSetting(acc, issuer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}

	return cmd
}

func CreateFeeMultiplier(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fee-multiplier [multiplier]",
		Short: "Set/update fee multiplier",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			multiplier := args[0]

			issuer := cliCtx.GetFromAddress()

			msg := fee.NewMsgMultiplier(multiplier, issuer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}

	return cmd
}

func CreateTokenFeeMultiplier(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-fee-multiplier [multiplier]",
		Short: "Set/update token fee multiplier",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			multiplier := args[0]

			issuer := cliCtx.GetFromAddress()

			msg := fee.NewMsgTokenMultiplier(multiplier, issuer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdkTypes.Msg{msg})
		},
	}

	return cmd
}

//Create the fee setting
func AddSysFeeSetting(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-fee",
		Short: "create new system fee setting,use add-fee with --chain-id",
		Long:  "add-fee cmd which use to create new system fee setting,use --chain-id",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			file, err := ioutil.ReadFile(dir + "/fee.json")
			if err != nil {
				log.Fatal(err)
			}

			data := FeeTransaction{}
			unmarshalErr := json.Unmarshal([]byte(file), &data)
			if unmarshalErr != nil {
				return unmarshalErr
			}
			viper.SetDefault("fees", data.Txfee)
			flags.GasFlagVar.Gas = 0
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			issuer := data.IssuerAddress
			feename := data.FeeName
			minStr := data.Min
			maxStr := data.Max
			percentageStr := data.percentage

			min, minErr := sdkTypes.ParseCoins(minStr)
			if minErr != nil {
				return minErr
			}
			max, maxErr := sdkTypes.ParseCoins(maxStr)
			if maxErr != nil {
				return maxErr
			}

			msg := fee.NewMsgSysFeeSetting(feename, min, max, percentageStr, issuer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			viper.SetDefault("from", issuer)
			viper.SetDefault("broadcast-mode", "sync")
			viper.SetDefault("output", "json")

			cli := context.NewCLIContext().WithCodec(cdc)
			return utils.CompleteAndBroadcastTxCLI(txBldr, cli, []sdkTypes.Msg{msg})
		},
	}
	return cmd
}
