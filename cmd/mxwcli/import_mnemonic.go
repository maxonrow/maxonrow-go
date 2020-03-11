package main

import (
	"fmt"

	clientKeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func addImportMnemonicCommand(keyCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-mnemonic <name> <mnemonic>",
		Short: "Import key from mnemonic phrase",
		Args:  cobra.ExactArgs(2),
		RunE:  runImportMnemonicCmd,
	}

	keyCmd.AddCommand(cmd)

	cmd.Flags().String("encryption_passphrase", "testtest", "Passphrase that will be used to encrypt private key")

	return cmd
}

func runImportMnemonicCmd(cmd *cobra.Command, args []string) error {

	var kb keys.Keybase

	kb, err := clientKeys.NewKeyringFromHomeFlag(cmd.InOrStdin())
	if err != nil {
		return err
	}

	name := args[0]
	mnemonic := args[1]

	encryptPwd, err := cmd.Flags().GetString("encryption_passphrase")
	if err != nil {
		return err
	}
	hdParams, err := hd.NewParamsFromPath(sdkTypes.GetConfig().GetFullFundraiserPath())
	if err != nil {
		return err
	}

	// TODO - bip39 password
	_, err = kb.Derive(name, mnemonic, keys.DefaultBIP39Passphrase, encryptPwd, *hdParams)
	if err != nil {
		return err
	}

	fmt.Println("Key Imported!")

	return nil

}
