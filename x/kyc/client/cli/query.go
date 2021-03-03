package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/x/kyc"
	"github.com/spf13/cobra"
)

func GetCmdIsWhitelisted(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "check [address]",
		Short: "check if address (in hex) is whitelisted",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			addressStr := args[0]

			// To prevent invalid address going to the server
			if _, err := sdkTypes.AccAddressFromBech32(addressStr); err != nil {
				return sdkTypes.ErrInvalidAddress(err.Error())
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, kyc.QueryIsWhitelisted, addressStr), nil)
			if err != nil {
				fmt.Printf("Could not check %s: %s\n", addressStr, err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
