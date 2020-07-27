package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tendermint/tendermint/crypto"

	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	cliKeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	app "github.com/maxonrow/maxonrow-go/app"
	"github.com/maxonrow/maxonrow-go/genesis"
	"github.com/maxonrow/maxonrow-go/types"
	"github.com/maxonrow/maxonrow-go/x/fee"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	tendermintConfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tm "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	flagOverwrite  = "overwrite"
	DefaultChainID = "maxonrow-chain"
)

var (
	DefaultNodeHome = os.ExpandEnv("$HOME/.mxw")
)

type NodeDetails struct {
	NodeId     []string
	ConsPublic []string
}

func main() {

	cobra.EnableCommandSorting = false

	cdc := app.MakeDefaultCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "mxw",
		Short:             "maxonrow App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(InitCmd(ctx, cdc))
	rootCmd.AddCommand(InitAuto(ctx, cdc))
	rootCmd.AddCommand(InitNode(ctx, cdc))
	rootCmd.AddCommand(GenTxCmd(ctx, cdc))
	rootCmd.AddCommand(AddGenesisAccountCmd(ctx, cdc))
	rootCmd.AddCommand(CreateStaking(ctx, cdc))
	rootCmd.AddCommand(Version)
	server.AddCommands(ctx, cdc, rootCmd, appCreator(ctx), appExporter(ctx))

	executor := cli.PrepareBaseCmd(rootCmd, "MXW", DefaultNodeHome)
	if err := executor.Execute(); err != nil {
		panic(err)
	}
}

func appCreator(ctx *server.Context) server.AppCreator {
	return func(logger log.Logger, database dbm.DB, traceStore io.Writer) abci.Application {
		return app.NewMXWApp(logger, database)
	}
}

func appExporter(ctx *server.Context) server.AppExporter {
	return func(logger log.Logger, db dbm.DB, _ io.Writer, _ int64, _ bool, _ []string) (
		json.RawMessage, []tm.GenesisValidator, error) {
		mxwapp := app.NewMXWApp(logger, db)
		return mxwapp.ExportStateAndValidators()
	}
}

// AddGenesisAccountCmd allows users to add accounts to the genesis file
func AddGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address] [coins[,coins]]",
		Short: "Adds an account to the genesis file",
		Args:  cobra.ExactArgs(2),
		Long: strings.TrimSpace(`
Adds accounts to the genesis file so that you can start a chain with coins in the CLI:

$ nsd add-genesis-account cosmos1tse7r2fadvlrrgau3pa0ss7cqh55wrv6y9alwh 1000STAKE,1000mycoin
`),
		RunE: func(_ *cobra.Command, args []string) error {
			addr, err := sdkTypes.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			coins, err := sdkTypes.ParseCoins(args[1])
			if err != nil {
				return err
			}
			coins.Sort()

			var genDoc tm.GenesisDoc
			config := ctx.Config
			genesisFilePath := config.GenesisFile()

			if !common.FileExists(genesisFilePath) {
				return fmt.Errorf("%s does not exist, run `mxwd init` first", genesisFilePath)
			}
			genContents, err := ioutil.ReadFile(genesisFilePath)
			if err != nil {
			}

			if err = cdc.UnmarshalJSON(genContents, &genDoc); err != nil {
				return err
			}

			var appState genesis.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			for _, stateAcc := range appState.Accounts {
				if stateAcc.Address.Equals(addr) {
					return fmt.Errorf("the application state already contains account %v", addr)
				}
			}

			acc := auth.NewBaseAccountWithAddress(addr)
			acc.Coins = coins
			appState.Accounts = append(appState.Accounts, &acc)
			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			return exportGenesisFile(genesisFilePath, genDoc.ChainID, genDoc.Validators, appStateJSON)
		},
	}

	return cmd
}

// SimpleAppGenTx returns a simple GenTx command that makes the node a valdiator from the start
func SimpleAppGenTx(cdc *codec.Codec, pk crypto.PubKey) (
	appGenTx, cliPrint json.RawMessage, validator tm.GenesisValidator, err error) {

	addr, secret, err := server.GenerateCoinKey()
	if err != nil {
		return
	}

	bz, err := cdc.MarshalJSON(struct {
		Addr sdkTypes.AccAddress `json:"addr"`
	}{addr})
	if err != nil {
		return
	}

	appGenTx = json.RawMessage(bz)

	bz, err = cdc.MarshalJSON(map[string]string{"secret": secret})
	if err != nil {
		return
	}

	cliPrint = json.RawMessage(bz)

	validator = tm.GenesisValidator{
		PubKey: pk,
		Power:  10,
	}

	return
}

// This will set up everything needed and create a genesis file with one validator
func InitAuto(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init-auto",
		Short: "Auto initialize node with genesis configuration, gentx, priv-validator file and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			//use to take input from user and create the node dir
			node, err := createNode()
			//store the list of node dir path
			var nodepath []string
			config := ctx.Config
			chainID := viper.GetString(client.FlagChainID)

			//get the root directory and delete the common config file
			if config.RootDir == viper.GetString(cli.HomeFlag) {
				configpath := config.RootDir + "/config"
				data := config.RootDir + "/data"
				fileexists(configpath, data)
			}

			var node_ids []string
			var pub_keys []crypto.PubKey
			var validator []string
			//generate the nodeId and private validator
			for j := 0; j < node; j++ {
				if err != nil {
					fmt.Println("node creating", err)
				}
				val := strconv.Itoa(j)
				config.SetRoot(viper.GetString(cli.HomeFlag) + "/Node" + val)
				nodepath = append(nodepath, config.RootDir)
				err = common.EnsureDir(config.RootDir, 0777)
				if err != nil {
					return errors.Wrap(err, "Failed to initialize validator files")
				}

				err = common.EnsureDir(config.RootDir+"/config", 0777)
				if err != nil {
					return errors.Wrap(err, "Failed to initialize validator files")
				}

				//Initialize the node key and private validator_key
				nodeID, pk, err := genutil.InitializeNodeValidatorFiles(config)
				if err != nil {
					return errors.Wrap(err, "Failed to initialize validator files")
				}
				fmt.Printf("Node ID: %s. Pub key: %s\n", nodeID, pk)

				node_ids = append(node_ids, nodeID)
				pub_keys = append(pub_keys, pk)
				va, _ := sdkTypes.Bech32ifyConsPub(pub_keys[j])
				validator = append(validator, va)

			}

			var appStateJSON json.RawMessage
			//generate the genesis
			gen := genesis.NewDefaultGenesisState()

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

			//create the dir for keys
			keyPath := filepath.Join(viper.GetString(cli.HomeFlag), "keys")
			// // keybase := keys.New("keys", keyPath)
			keybase, kbErr := cliKeys.NewKeyringFromHomeFlag(cmd.InOrStdin())
			if kbErr != nil {
				return kbErr
			}

			fmt.Println("Creating account. (All Passwords set to `12345678`)")
			//create the  list of accounts and append to genesis
			for i := 0; i < 30; i++ {
				name := fmt.Sprintf("acc-%v", i+1)
				info, mnemonic, err := keybase.CreateMnemonic(name, keys.English, "12345678", keys.Secp256k1)

				if err != nil {
					return fmt.Errorf("Unable to create new account: %v", err)
				}
				fmt.Printf("Create new account. name: %v, address: %s, mnemonic:%s\n", name, info.GetAddress(), mnemonic)
				addr := info.GetAddress()
				acc := auth.NewBaseAccountWithAddress(addr)
				bal, _ := sdkTypes.NewIntFromString("1000000000000000000000000000000000")
				acc.Coins = []sdkTypes.Coin{sdkTypes.Coin{
					Denom:  types.CIN,
					Amount: bal,
				}}
				gen.Accounts = append(gen.Accounts, &acc)
				gen.KycState.AuthorizedAddresses = append(gen.KycState.AuthorizedAddresses, addr)
				gen.KycState.WhitelistedAddresses = append(gen.KycState.WhitelistedAddresses, addr)

			}

			config.SetRoot(nodepath[0])
			//get the genesis file path
			genesisFilePath := config.GenesisFile()
			if !viper.GetBool(flagOverwrite) && common.FileExists(genesisFilePath) {
				return fmt.Errorf("genesis.json file already exists at path: %v", genesisFilePath)
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

			//Add validator account
			genState, err = addValidatorSet(ctx, cdc, genState, accIndex1, validator)
			if err != nil {
				return fmt.Errorf("Unable to add validator set account: %v", err)
			}

			appStateJSON, err = codec.MarshalJSONIndent(cdc, genState)
			if err != nil {
				return err
			}
			if err = exportGenesisFile(genesisFilePath, chainID, nil, appStateJSON); err != nil {
				return errors.Wrap(err, "Failed to populate genesis.json")
			}
			//copy the keys dir to all node dir
			copyKeyNode(keyPath, nodepath)
			for j := 0; j < node; j++ {
				valAddress, _ := sdkTypes.Bech32ifyConsPub(pub_keys[j])
				cmd := exec.Command("mxwd", "gentx", "--name", fmt.Sprintf("acc-%v", j+1), "--home", nodepath[0], "--pubkey", valAddress, "--node-id", node_ids[j])
				stdin, err := cmd.StdinPipe()
				if err != nil {
					fmt.Println(fmt.Sprint(err))
				}
				defer stdin.Close()
				err = cmd.Start()
				if err != nil {
					fmt.Println(fmt.Sprint(err))
				}
				io.WriteString(stdin, "12345678\n")
				cmd.Wait()
			}
			config.ProfListenAddress = ""

			var np = NodeDetails{NodeId: node_ids, ConsPublic: validator}
			writeNodes(viper.GetString(cli.HomeFlag), np)

			// This generates a tendermint configuration. Not sure where it should go
			tendermintConfig.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			fmt.Printf("Initialized mxwd configuration and bootstrapping files in %s...\n", viper.GetString(cli.HomeFlag))
			fileexists(keyPath, "")
			if len(nodepath) > 1 {
				for _, path := range nodepath {
					//copy the genesis file
					err := CopyFile(config.GenesisFile(), path+"/config/genesis.json")
					if err != nil {
						return err
					}
					//copy the config file
					err = CopyFile(config.RootDir+"/config"+"/config.toml", path+"/config/config.toml")
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
	cmd.Flags().String(cli.HomeFlag, DefaultNodeHome, "node's home directory")
	cmd.Flags().String(client.FlagChainID, DefaultChainID, "genesis file chain-id")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(client.FlagKeyringBackend, client.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	return cmd
}

func InitNode(ctx *server.Context, cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "init-node",
		Short: "Auto initialize node configuration, priv-validator file and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {

			//store the list of node dir path
			var nodepath []string
			config := ctx.Config

			//get the root directory and delete the common config file
			if config.RootDir == viper.GetString(cli.HomeFlag) {
				configpath := config.RootDir + "/config"
				data := config.RootDir + "/data"
				fileexists(configpath, data)
			}

			config.SetRoot(viper.GetString(cli.HomeFlag))
			nodepath = append(nodepath, config.RootDir)
			err := common.EnsureDir(config.RootDir, 0777)
			if err != nil {
				return errors.Wrap(err, "Failed to initialize validator files")
			}

			err = common.EnsureDir(config.RootDir+"/config", 0777)
			if err != nil {
				return errors.Wrap(err, "Failed to initialize validator files")
			}

			//Initialize the node key and private validator_key
			nodeID, pk, err := genutil.InitializeNodeValidatorFiles(config)
			if err != nil {
				return errors.Wrap(err, "Failed to initialize validator files")
			}
			fmt.Printf("Node ID: %s. Pub key: %s\n", nodeID, pk)

			return nil
		},
	}
	cmd.Flags().String(cli.HomeFlag, DefaultNodeHome, "node's home directory")
	cmd.Flags().String(client.FlagChainID, DefaultChainID, "genesis file chain-id")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	return cmd
}

//Add accounts to kyc
func addAccountKyc(ctx *server.Context, cdc *codec.Codec, genState genesis.GenesisState) (genesis.GenesisState, int, error) {
	var startIndex = 0
	var endIndex = 0
	acc := genState.Accounts
	j := startIndex
	for ; j < 4; j++ {
		addr := acc[j].Address
		genState.KycState.IssuerAddresses = append(genState.KycState.IssuerAddresses, addr)
	}

	for ; j < 4+3; j++ {
		addr := acc[j].Address
		genState.KycState.ProviderAddresses = append(genState.KycState.ProviderAddresses, addr)
		endIndex = j
	}

	return genState, endIndex, nil
}

//Add accounts to NameService
func addAccounNameService(ctx *server.Context, cdc *codec.Codec, genState genesis.GenesisState, startIndex int) (genesis.GenesisState, int, error) {
	var endIndex = 0
	acc := genState.Accounts
	j := startIndex
	for ; j < 10; j++ {
		addr := acc[j].Address
		genState.NameServiceState.AuthorisedAddresses = append(genState.NameServiceState.AuthorisedAddresses, addr)

	}
	for ; j < 10+2; j++ {
		addr := acc[j].Address
		genState.NameServiceState.IssuerAddresses = append(genState.NameServiceState.IssuerAddresses, addr)
	}
	for ; j < 10+2+2; j++ {
		addr := acc[j].Address
		genState.NameServiceState.ProviderAddresses = append(genState.NameServiceState.ProviderAddresses, addr)
		endIndex = j
	}
	return genState, endIndex, nil
}

//Add accounts to Token
func addAccounttoken(ctx *server.Context, cdc *codec.Codec, genState genesis.GenesisState, startIndex int) (genesis.GenesisState, int, error) {
	var endIndex = 0
	acc := genState.Accounts
	j := startIndex
	for ; j < 16; j++ {
		addr := acc[j].Address
		genState.TokenState.AuthorizedAddresses = append(genState.TokenState.AuthorizedAddresses, addr)
	}
	for ; j < 16+3; j++ {
		addr := acc[j].Address
		genState.TokenState.IssuerAddresses = append(genState.TokenState.IssuerAddresses, addr)
	}
	for ; j < 16+3+3; j++ {
		addr := acc[j].Address
		genState.TokenState.ProviderAddresses = append(genState.TokenState.ProviderAddresses, addr)
		endIndex = j
	}
	return genState, endIndex, nil
}

//Add accounts to fee
func addAccountfee(ctx *server.Context, cdc *codec.Codec, genState genesis.GenesisState, startIndex int) (genesis.GenesisState, int, error) {
	var endIndex = 0
	acc := genState.Accounts
	j := startIndex
	for ; j < 26; j++ {
		addr := acc[j].Address
		genState.FeeState.AuthorisedAddresses = append(genState.FeeState.AuthorisedAddresses, addr)
		endIndex = j
	}
	return genState, endIndex, nil
}

//Add accounts to Maintainer
func addMainter(ctx *server.Context, cdc *codec.Codec, genState genesis.GenesisState, startIndex int) (genesis.GenesisState, int, error) {
	var endIndex = 0
	acc := genState.Accounts
	j := startIndex
	for ; j < 15; j++ {
		addr := acc[j].Address
		genState.MaintenanceState.Maintainers = append(genState.MaintenanceState.Maintainers, addr)
		endIndex = j
	}
	return genState, endIndex, nil
}

//Add accounts to ValidatorSet
func addValidatorSet(ctx *server.Context, cdc *codec.Codec, genState genesis.GenesisState, startIndex int, validator []string) (genesis.GenesisState, error) {
	for _, val := range validator {
		genState.MaintenanceState.ValidatorSet = append(genState.MaintenanceState.ValidatorSet, val)
	}
	return genState, nil
}
