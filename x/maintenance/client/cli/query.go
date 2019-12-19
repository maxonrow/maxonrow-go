package cli

import (
	"fmt"
	"strconv"

	"github.com/maxonrow/maxonrow-go/x/maintenance"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetCmdGetProposal(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proposal [proposalID]",
		Short: "get proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposalID := args[0]

			proposalIDInt, err := strconv.ParseUint(proposalID, 10, 64)

			params := maintenance.NewQueryProposalParams(proposalIDInt)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/proposal", queryRoute), bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
