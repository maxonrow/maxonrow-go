package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/badoux/checkmail"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/manifoldco/promptui"
	"github.com/maxonrow/maxonrow-go/app"
	mxwAuthCmd "github.com/maxonrow/maxonrow-go/x/auth/client/cli"
	bankcmd "github.com/maxonrow/maxonrow-go/x/bank/client/cli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	//kycClient "github.com/maxonrow/maxonrow-go/x/kyc/client"
	ver "github.com/maxonrow/maxonrow-go/version"
	authClient "github.com/maxonrow/maxonrow-go/x/auth/client"
	feeClient "github.com/maxonrow/maxonrow-go/x/fee/client"
	maintenanceClient "github.com/maxonrow/maxonrow-go/x/maintenance/client"
	nsClient "github.com/maxonrow/maxonrow-go/x/nameservice/client"
	tokenClient "github.com/maxonrow/maxonrow-go/x/token/fungible/client"
	nftClient "github.com/maxonrow/maxonrow-go/x/token/nonfungible/client"
	"gopkg.in/cheggaaa/pb.v1"
)

const (
	storeAcc         = "acc"
	storeStake       = "staking"
	storeDistr       = "distr"
	storeNS          = "nameservice"
	storeKyc         = "kyc"
	storeToken       = "token"
	storeFee         = "fee"
	storeMaintenance = "maintenance"
	storeAuth        = "auth"
)

var (
	defaultCLIHome = os.ExpandEnv("$HOME/.mxw")
)

func main() {

	cobra.EnableCommandSorting = false
	cdc := app.MakeDefaultCodec()

	rootCmd := &cobra.Command{
		Use:   "mxwcli",
		Short: "maxonrow client",
	}

	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	keyComd := keys.Commands()
	addImportMnemonicCommand(keyComd)
	multisigAddressCommand(keyComd)

	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(defaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		client.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		keyComd,
		client.LineBreak,
		client.LineBreak,
		bechCommand(),
		kycCommand(),
		createKeyPairCommand(),
		Version,
	)

	executor := cli.PrepareMainCmd(rootCmd, "MXW", defaultCLIHome)
	if err := executor.Execute(); err != nil {
		panic(err)
	}
}

func queryCmd(cdc *amino.Codec) *cobra.Command {

	maintenanceModuleClient := maintenanceClient.NewModuleClient(storeMaintenance, cdc)
	nsModuleClient := nsClient.NewModuleClient(storeNS, cdc)
	//kycModuleClient := kycClient.NewModuleClient(storeKyc, cdc)
	tokenModuleClient := tokenClient.NewModuleClient(storeToken, cdc)
	feeModuleClient := feeClient.NewModuleClient(storeFee, cdc)
	authModuleClient := authClient.NewModuleClient(storeAuth, cdc)

	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(cdc),
		client.LineBreak,
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(cdc),
		client.LineBreak,

		// TO-DO: implement appmodulebasic interface in every module.
		// fee cli-tx
		nsModuleClient.GetQueryCmd(),
		feeModuleClient.GetQueryCmd(),
		tokenModuleClient.GetQueryCmd(),
		maintenanceModuleClient.GetQueryCmd(),
		authModuleClient.GetQueryCmd(),
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {

	maintenanceModuleClient := maintenanceClient.NewModuleClient(storeMaintenance, cdc)
	nsModuleClient := nsClient.NewModuleClient(storeNS, cdc)
	//kycModuleClient := kycClient.NewModuleClient(storeKyc, cdc)
	tokenModuleClient := tokenClient.NewModuleClient(storeToken, cdc)
	nftmoduleClient := nftClient.NewModuleClient(storeToken, cdc)
	feeModuleClient := feeClient.NewModuleClient(storeFee, cdc)
	authModuleClient := authClient.NewModuleClient(storeAuth, cdc)

	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		mxwAuthCmd.CreateMultiSigAccountCmd(cdc),
		bankcmd.SendTxCmd(cdc),
		client.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		client.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		client.LineBreak,

		// TO-DO: implement appmodulebasic interface in every module.
		// fee cli-tx
		nsModuleClient.GetTxCmd(),
		feeModuleClient.GetTxCmd(),
		tokenModuleClient.GetTxCmd(),
		nftmoduleClient.GetTxCmd(),
		maintenanceModuleClient.GetTxCmd(),
		authModuleClient.GetTxCmd(),
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	// remove auth and bank commands as they're mounted under the root tx command
	var cmdsToRemove []*cobra.Command

	for _, cmd := range txCmd.Commands() {
		if cmd.Use == auth.ModuleName || cmd.Use == bank.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	txCmd.RemoveCommand(cmdsToRemove...)

	return txCmd
}

// registerRoutes registers the routes from the different modules for the LCD.
// NOTE: details on the routes added for each module are in the module documentation
// NOTE: If making updates here you also need to update the test helper in client/lcd/test_helper.go
func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}

func kycCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kyc",
		Aliases: []string{"k"},
		Short:   "kyc register",
	}

	registerCmd := &cobra.Command{
		Use:   "register <address>",
		Short: "request KYC whitelisting from KYC authorised",
		Args:  cobra.ExactArgs(1),
		RunE:  runKYCRegisterCmd,
	}

	cmd.AddCommand(registerCmd)

	return cmd
}

func runKYCRegisterCmd(_ *cobra.Command, args []string) error {
	address := args[0]

	namePrompt := promptui.Prompt{
		Label: "Full name",
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("Invalid name!")
			}

			return nil
		},
	}

	name, err := namePrompt.Run()

	if err != nil {
		fmt.Printf("Name input failed, %v\n", err)
		return err
	}

	emailPrompt := promptui.Prompt{
		Label: "Email address",
		Validate: func(input string) error {
			err := checkmail.ValidateFormat(input)

			if err != nil {
				return errors.New("Invalid email address!")
			}

			return nil
		},
	}

	email, err := emailPrompt.Run()

	if err != nil {
		fmt.Printf("Email input failed, %v\n", err)
		return err
	}

	bar := pb.StartNew(2)
	bar.ShowTimeLeft = false
	bar.ShowSpeed = false

	url := fmt.Sprintf("http://68.183.213.3:5555/kyc/%s", address)

	jsonStr, err := json.Marshal(map[string]string{"name": name, "email": email})

	if err != nil {
		fmt.Printf("Error while converting to JSON, %v\n", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	bar.Increment()

	if err != nil {
		fmt.Printf("Error while communicating with KYC middleware, %v\n", err)
		return err
	}

	defer resp.Body.Close()

	if resp.Status == "422 Unprocessable Entity" || resp.Status == "500 Internal Server Error" {
		fmt.Printf("Invalid data provided.")
		return nil
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	data := result["data"].(map[string]interface{})

	bar.Increment()

	if data["kyc"].(bool) {
		bar.FinishPrint("Approval complete!")
	} else {
		bar.FinishPrint("An error ocurred. Please try again.")
	}

	return nil
}

func bechCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bech",
		Aliases: []string{"b"},
		Short:   "bech encoding/decoding",
	}

	encodeCmd := &cobra.Command{
		Use:   "encode <prefix> <value>",
		Short: "encodes as bech32 value with given prefix",
		Args:  cobra.ExactArgs(2),
		RunE:  runEncodeCmd,
	}
	encodeCmd.Flags().Int("base", 16, "Input encoding")

	decodeCmd := &cobra.Command{
		Use:   "decode <value>",
		Short: "decodes bech32-encoded data",
		Args:  cobra.ExactArgs(1),
		RunE:  runDecodeCmd,
	}

	cmd.AddCommand(encodeCmd, decodeCmd)

	return cmd
}

func runEncodeCmd(_ *cobra.Command, args []string) error {
	prefix := args[0]
	value := args[1]

	var binaryValue []byte
	base := viper.GetInt("base")
	switch base {
	case 16:
		binaryValue = make([]byte, hex.DecodedLen(len(value)))
		_, err := hex.Decode(binaryValue, []byte(value))
		if err != nil {
			return errors.Wrap(err, "Invalid hex value")
		}
	case 64:
		binaryValue = make([]byte, base64.StdEncoding.DecodedLen(len(value)))
		l, err := base64.StdEncoding.Decode(binaryValue, []byte(value))
		if err != nil {
			return errors.Wrap(err, "Invalid base64 value")
		}
		binaryValue = binaryValue[:l]

		fmt.Printf("%d bytes. First: %d. Last: %d",
			len(binaryValue), binaryValue[0], binaryValue[len(binaryValue)-1])
	default:
		return errors.New("Invalid base")
	}

	b32Value, err := bech32.ConvertBits(binaryValue, 8, 5, true)
	if err != nil {
		return errors.Wrap(err, "Failed to convert bits")
	}

	encoded, err := bech32.Encode(prefix, b32Value)
	if err != nil {
		return errors.Wrap(err, "Failed to encode")
	}

	fmt.Println(string(encoded))

	return nil
}

func runDecodeCmd(_ *cobra.Command, args []string) error {
	value := args[0]

	readable, bech32Data, err := bech32.Decode(value)
	if err != nil {
		return errors.Wrap(err, "Failed to decode")
	}

	data, err := bech32.ConvertBits(bech32Data, 5, 8, false)
	if err != nil {
		return errors.Wrap(err, "Failed to convert bits")
	}

	fmt.Println(readable)

	var hexData = make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(hexData, data)

	fmt.Println(string(hexData))

	return nil
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}

var Version = &cobra.Command{
	Use:   "version",
	Short: "Print version info of maxonrow",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Maxonrow:", ver.Version)
	},
}
