package main

import (
	"fmt"
	"strconv"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/x/auth"
	"github.com/spf13/cobra"
)

func multisigAddressCommand(keyCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multisig-address <address> <sequence>",
		Short: "Generate multisig address from owner address and sequence (group address)",
		Args:  cobra.ExactArgs(2),
		RunE:  runMultisigAddressCmd,
	}

	keyCmd.AddCommand(cmd)

	return cmd
}

func runMultisigAddressCmd(cmd *cobra.Command, args []string) error {

	addr, err := sdkTypes.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	seq, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return err
	}

	groupAddr := auth.DeriveMultiSigAddress(addr, seq)
	fmt.Println(groupAddr.String())
	return nil

}
