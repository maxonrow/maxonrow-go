package cli

import (
	"fmt"

	"github.com/maxonrow/maxonrow-go/x/auth"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetMultiSigAcc(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acc [group address]",
		Short: "get multisig account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			groupAddr := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", "auth", auth.QueryMultiSigAcc, groupAddr), nil)
			if err != nil {
				fmt.Printf("Could not get  multisig account: %s\n", err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}

	return cmd
}
