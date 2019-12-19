package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	cliKeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tendermintConfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/maxonrow/maxonrow-go/genesis"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/fee"
)

func InitCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize genesis configuration, priv-validator file and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			chainID := viper.GetString(client.FlagChainID)

			nodeID, pk, err := genutil.InitializeNodeValidatorFiles(config)
			if err != nil {
				return errors.Wrap(err, "Failed to initialize validator files")
			}

			fmt.Printf("Node ID: %s. Pub key: %s\n", nodeID, pk)

			var appStateJSON json.RawMessage
			genesisFilePath := config.GenesisFile()

			if !viper.GetBool(flagOverwrite) && common.FileExists(genesisFilePath) {
				return fmt.Errorf("genesis.json file already exists at path: %v", genesisFilePath)
			}

			gen := genesis.NewDefaultGenesisState()
			min, _ := sdkTypes.NewIntFromString("0")
			max, _ := sdkTypes.NewIntFromString("0")
			feezero := []fee.GenesisFeeSetting{fee.GenesisFeeSetting{
				Name: "zero",
				Min: []sdkTypes.Coin{
					sdkTypes.Coin{
						Denom:  types.CIN,
						Amount: min,
					},
				},
				Max: []sdkTypes.Coin{
					sdkTypes.Coin{
						Denom:  types.CIN,
						Amount: max,
					},
				},
				Percentage: "0",
			}}

			for i := range feezero {
				gen.FeeState.FeeSettings = append(gen.FeeState.FeeSettings, feezero[i])
			}

			feestate := []fee.AssignMsgFeeSetting{fee.AssignMsgFeeSetting{Name: "zero", MsgType: "kyc-whitelist"},
				{Name: "zero", MsgType: "kyc-whitelist"},
				{Name: "zero", MsgType: "kyc-revokeWhitelist"},
				{Name: "zero", MsgType: "fee-createFeeSetting"},
				{Name: "zero", MsgType: "fee-createTxFeeSetting"},
				{Name: "zero", MsgType: "fee-createMultiplier"},
				{Name: "zero", MsgType: "fee-createAccFeeSetting"},
			}

			for i := range feestate {
				gen.FeeState.AssignedMsgFeeSettings = append(gen.FeeState.AssignedMsgFeeSettings, feestate[i])
			}

			// keyPath := filepath.Join(config.RootDir, "keys")
			// keybase := keys.New("keys", keyPath)

			keybase, kbErr := cliKeys.NewKeyringFromHomeFlag(cmd.InOrStdin())
			if kbErr != nil {
				return err
			}

			fmt.Println("Creating account. (All Passwords set to `12345678`)")
			for i := 0; i < 30; i++ {
				name := fmt.Sprintf("acc-%v", i+1)
				info, mnemonic, err := keybase.CreateMnemonic(name, keys.English, "12345678", keys.Secp256k1)

				if err != nil {
					return fmt.Errorf("Unable to create new account: %v", err)
				}
				fmt.Printf("Create new account. name: %v, address: %s, mnemonic:%s\n", name, info.GetAddress(), mnemonic)

				addr := info.GetAddress()
				acc := auth.NewBaseAccountWithAddress(addr)
				bal, _ := sdkTypes.NewIntFromString("1000000000000000000000")
				acc.Coins = []sdkTypes.Coin{sdkTypes.Coin{
					Denom:  types.CIN,
					Amount: bal,
				}}

				gen.Accounts = append(gen.Accounts, &acc)
				gen.KycState.AuthorizedAddresses = append(gen.KycState.AuthorizedAddresses, addr)
				gen.KycState.WhitelistedAddresses = append(gen.KycState.WhitelistedAddresses, addr)

			}

			//Add account for AccountKyc
			genState, accIndex, err := addAccountKyc(ctx, cdc, gen)
			if err != nil {
				return fmt.Errorf("Unable to add account: %v", err)
			}
			//Add account for NameService
			genState, accIndex, err = addAccounNameService(ctx, cdc, genState, accIndex)
			if err != nil {
				return fmt.Errorf("Unable to add account: %v", err)
			}

			//Add account for token
			genState, accIndex, err = addAccounttoken(ctx, cdc, genState, accIndex)
			if err != nil {
				return fmt.Errorf("Unable to add account: %v", err)
			}

			//Add account for fee
			genState, accIndex, err = addAccountfee(ctx, cdc, genState, accIndex)
			if err != nil {
				return fmt.Errorf("Unable to add account: %v", err)
			}

			//Add maintainers account
			var accIndex1 = 8
			genState, accIndex, err = addMainter(ctx, cdc, genState, accIndex1)
			if err != nil {
				return fmt.Errorf("Unable to add maintainers account: %v", err)
			}

			appStateJSON, err = codec.MarshalJSONIndent(cdc, genState)
			if err != nil {
				return err
			}

			if err = exportGenesisFile(genesisFilePath, chainID, nil, appStateJSON); err != nil {
				return errors.Wrap(err, "Failed to populate genesis.json")
			}

			config.ProfListenAddress = ""

			// This generates a tendermint configuration. Not sure where it should go
			tendermintConfig.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			fmt.Printf("Initialized mxwd configuration and bootstrapping files in %s...\n", viper.GetString(cli.HomeFlag))

			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, DefaultNodeHome, "node's home directory")
	cmd.Flags().String(client.FlagChainID, DefaultChainID, "genesis file chain-id")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")

	return cmd
}
