package cli

import (
	"bufio"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/bank"
)

const (
	flagTo     = "to"
	flagAmount = "amount"
)

type Resolve struct {
	Alias   string `json:"alias"`
	Address string `json:"address"`
}

func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-alias [alias] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Create and sign a send-alias tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc)

			toStr := args[0]
			to, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/nameservice/resolve/%s", toStr), nil)
			if err != nil {
				return err
			}

			amount := args[1]
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}

			if len(coins) > 1 {
				return errors.New("Only 1 token is supported")
			}

			coin := coins[0]
			if coin.Denom == "mxw" {
				coin = types.MXWtoCIN(coin)
			}

			// ensure account has enough coins
			coins = sdk.Coins{coin}
			from := cliCtx.GetFromAddress()
			msg := bank.NewMsgSend(from, to, coins)
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg})
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}
