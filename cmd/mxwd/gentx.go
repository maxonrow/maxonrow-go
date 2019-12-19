package main

// DONTCOVER

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	kbkeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/maxonrow/maxonrow-go/genesis"
	"github.com/maxonrow/maxonrow-go/types"
)

var (
	defaultTokens                  = sdk.TokensFromConsensusPower(100)
	defaultAmount                  = defaultTokens.String() + types.CIN
	defaultCommissionRate          = "0.1"
	defaultCommissionMaxRate       = "0.2"
	defaultCommissionMaxChangeRate = "0.01"
	defaultMinSelfDelegation       = "1"
)

// GenTxCmd builds the gentx command.
// nolint: errcheck
func GenTxCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gentx",
		Short: "Generate a genesis tx carrying a self delegation",
		Args:  cobra.NoArgs,
		Long: fmt.Sprintf(`This command is an alias of the 'mxwd tx create-validator' command'.

It creates a genesis piece carrying a self delegation with the
following delegation and commission default parameters:

	delegation amount:           %s
	commission rate:             %s
	commission max rate:         %s
	commission max change rate:  %s
	minimum self delegation:     %s
`, defaultAmount, defaultCommissionRate, defaultCommissionMaxRate, defaultCommissionMaxChangeRate, defaultMinSelfDelegation),
		RunE: func(cmd *cobra.Command, args []string) error {

			config := ctx.Config
			config.SetRoot(viper.GetString(tmcli.HomeFlag))
			nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(ctx.Config)
			if err != nil {
				return err
			}

			// Read --nodeID, if empty take it from priv_validator.json
			if nodeIDString := viper.GetString(cli.FlagNodeID); nodeIDString != "" {
				nodeID = nodeIDString
			}

			ip := viper.GetString(cli.FlagIP)
			if ip == "" {
				fmt.Fprintf(os.Stderr, "couldn't retrieve an external IP; "+
					"the tx's memo field will be unset")
			}

			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return err
			}

			genesisState := genesis.GenesisState{}
			if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
				return err
			}

			if err = genesisState.Validate(); err != nil {
				return err
			}

			keyPath := filepath.Join(config.RootDir, "keys")
			kb := kbkeys.New("keys", keyPath)
			name := viper.GetString(client.FlagName)
			key, err := kb.Get(name)
			if err != nil {
				return err
			}

			// Read --pubkey, if empty take it from priv_validator.json
			if valPubKeyString := viper.GetString(cli.FlagPubKey); valPubKeyString != "" {
				valPubKey, err = sdk.GetConsPubKeyBech32(valPubKeyString)
				if err != nil {
					return err
				}
			}

			va, err := sdk.Bech32ifyConsPub(valPubKey)
			if err != nil {
				return err
			}
			genesisState.MaintenanceState.ValidatorSet = append(genesisState.MaintenanceState.ValidatorSet, va)

			website := viper.GetString(cli.FlagWebsite)
			details := viper.GetString(cli.FlagDetails)
			identity := viper.GetString(cli.FlagIdentity)

			// Set flags for creating gentx
			prepareFlagsForTxCreateValidator(config, nodeID, ip, genDoc.ChainID, valPubKey, website, details, identity, config.Moniker)

			// Fetch the amount of coins staked
			amount := viper.GetString(cli.FlagAmount)
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}

			err = accountInGenesis(genesisState, key.GetAddress(), coins)
			if err != nil {
				return err
			}

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc)).WithKeybase(kb)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// XXX: Set the generate-only flag here after the CLI context has
			// been created. This allows the from name/key to be correctly populated.
			//
			// TODO: Consider removing the manual setting of generate-only in
			// favor of a 'gentx' flag in the create-validator command.
			viper.Set(client.FlagGenerateOnly, true)

			// create a 'create-validator' message
			txBldr, msg, err := cli.BuildCreateValidatorMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}

			info, err := txBldr.Keybase().Get(name)
			if err != nil {
				return err
			}

			if info.GetType() == kbkeys.TypeOffline || info.GetType() == kbkeys.TypeMulti {
				fmt.Println("Offline key passed in. Use `mxwcli tx sign` command to sign:")
				return utils.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg})
			}

			// write the unsigned transaction to the buffer
			w := bytes.NewBuffer([]byte{})
			cliCtx = cliCtx.WithOutput(w)

			if err = utils.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg}); err != nil {
				return err
			}

			// read the transaction
			stdTx, err := readUnsignedGenTxFile(cdc, w)
			if err != nil {
				return err
			}

			// mxw msgs are gas free
			stdTx.Fee.Gas = 0
			stdTx.Fee.Amount, _ = types.ParseCoins("0cin")

			// sign the transaction and write it to the output file
			signedTx, err := utils.SignStdTx(txBldr, cliCtx, name, stdTx, false, true)
			if err != nil {
				return err
			}

			outputDocument := viper.GetString(client.FlagOutputDocument)

			if outputDocument == "" {
				outputDocument, err = makeOutputFilepath(config.RootDir, nodeID)
				if err != nil {
					return err
				}
			}
			if err := writeSignedGenTx(cdc, outputDocument, signedTx); err != nil {
				return err
			}

			if err := updateGenesis(cdc, ctx, genDoc, genesisState, signedTx); err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "Genesis transaction written to %q\n", outputDocument)
			return nil

		},
	}

	ip, _ := server.ExternalIP()

	cmd.Flags().String(tmcli.HomeFlag, DefaultNodeHome, "node's home directory")
	cmd.Flags().String(client.FlagName, "", "name of private key with which to sign the gentx")
	cmd.Flags().String(client.FlagOutputDocument, "",
		"write the genesis transaction JSON document to the given file instead of the default location")
	cmd.Flags().String(cli.FlagIP, ip, "The node's public IP")
	cmd.Flags().String(cli.FlagNodeID, "", "The node's NodeID")
	cmd.Flags().String(cli.FlagWebsite, "", "The validator's (optional) website")
	cmd.Flags().String(cli.FlagDetails, "", "The validator's (optional) details")
	cmd.Flags().String(cli.FlagIdentity, "", "The (optional) identity signature (ex. UPort or Keybase)")
	cmd.Flags().AddFlagSet(cli.FsCommissionCreate)
	cmd.Flags().AddFlagSet(cli.FsMinSelfDelegation)
	cmd.Flags().AddFlagSet(cli.FsAmount)
	cmd.Flags().AddFlagSet(cli.FsPk)
	cmd.Flags().String(cli.FlagMoniker, "", "The validator's node name")
	cmd.MarkFlagRequired(client.FlagName)
	return cmd
}

func accountInGenesis(genesisState genesis.GenesisState, key sdk.AccAddress, coins sdk.Coins) error {
	accountIsInGenesis := false
	bondDenom := genesisState.StakingState.Params.BondDenom

	// Check if the account is in genesis
	for _, acc := range genesisState.Accounts {
		// Ensure that account is in genesis
		if acc.Address.Equals(key) {

			// Ensure account contains enough funds of default bond denom
			if coins.AmountOf(bondDenom).GT(acc.Coins.AmountOf(bondDenom)) {
				return fmt.Errorf(
					"account %v is in genesis, but it only has %v%v available to stake, not %v%v",
					key.String(), acc.Coins.AmountOf(bondDenom), bondDenom, coins.AmountOf(bondDenom), bondDenom,
				)
			}
			accountIsInGenesis = true
			break
		}
	}

	if accountIsInGenesis {
		return nil
	}

	return fmt.Errorf("account %s in not in the app_state.accounts array of genesis.json", key)
}

func prepareFlagsForTxCreateValidator(
	config *cfg.Config, nodeID, ip, chainID string, valPubKey crypto.PubKey, website, details, identity, moniker string,
) {
	viper.Set(client.FlagChainID, chainID)
	viper.Set(client.FlagFrom, viper.GetString(client.FlagName))
	viper.Set(cli.FlagNodeID, nodeID)
	viper.Set(cli.FlagIP, ip)
	viper.Set(cli.FlagPubKey, sdk.MustBech32ifyConsPub(valPubKey))

	viper.Set(cli.FlagWebsite, website)
	viper.Set(cli.FlagDetails, details)
	viper.Set(cli.FlagIdentity, identity)
	if moniker == "" {
		viper.Set(cli.FlagMoniker, config.Moniker)
	} else {
		viper.Set(cli.FlagMoniker, moniker)
	}
	if config.Moniker == "" {
		viper.Set(cli.FlagMoniker, viper.GetString(client.FlagName))
	}
	if viper.GetString(cli.FlagAmount) == "" {
		viper.Set(cli.FlagAmount, defaultAmount)
	}
	if viper.GetString(cli.FlagCommissionRate) == "" {
		viper.Set(cli.FlagCommissionRate, defaultCommissionRate)
	}
	if viper.GetString(cli.FlagCommissionMaxRate) == "" {
		viper.Set(cli.FlagCommissionMaxRate, defaultCommissionMaxRate)
	}
	if viper.GetString(cli.FlagCommissionMaxChangeRate) == "" {
		viper.Set(cli.FlagCommissionMaxChangeRate, defaultCommissionMaxChangeRate)
	}
	if viper.GetString(cli.FlagMinSelfDelegation) == "" {
		viper.Set(cli.FlagMinSelfDelegation, defaultMinSelfDelegation)
	}
}

func makeOutputFilepath(rootDir, nodeID string) (string, error) {
	writePath := filepath.Join(rootDir, "config", "gentx")
	if err := common.EnsureDir(writePath, 0700); err != nil {
		return "", err
	}
	return filepath.Join(writePath, fmt.Sprintf("gentx-%v.json", nodeID)), nil
}

func readUnsignedGenTxFile(cdc *codec.Codec, r io.Reader) (auth.StdTx, error) {
	var stdTx auth.StdTx
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return stdTx, err
	}
	err = cdc.UnmarshalJSON(bytes, &stdTx)
	return stdTx, err
}

// nolint: errcheck
func writeSignedGenTx(cdc *codec.Codec, outputDocument string, tx auth.StdTx) error {
	outputFile, err := os.OpenFile(outputDocument, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	json, err := cdc.MarshalJSON(tx)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(outputFile, "%s\n", json)
	return err
}

// Update the genesis file when the gentx executed
func updateGenesis(cdc *codec.Codec, ctx *server.Context, genDoc *tmtypes.GenesisDoc, genContent genesis.GenesisState, tx auth.StdTx) error {
	config := ctx.Config
	signedtx, err := cdc.MarshalJSON(tx)
	if err != nil {
		return err
	}
	genContent.GenTxs = append(genContent.GenTxs, signedtx)
	appstate, err := cdc.MarshalJSON(genContent)
	if err != nil {
		return err
	}
	genDoc.AppState = appstate
	genutil.ExportGenesisFile(genDoc, config.GenesisFile())
	return nil
}

//create the staking transaction without the signature and crearte the file in gentx folder
type Createvalidator struct {
	ValConPubkey string
	Moniker      string
	IpAddress    string
	NodeID       string
	ChainId      string
}

func CreateStaking(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-staking",
		Short: "Generate a validator transaction  carrying a self delegation",
		Args:  cobra.NoArgs,
		Long: fmt.Sprintf(`This command is an alias of the 'mxwd tx create-validator' command'.

It creates a validator transaction piece carrying a self delegation with the
following delegation and commission default parameters:
	delegation amount:           %s
	commission rate:             %s
	commission max rate:         %s
	commission max change rate:  %s
	minimum self delegation:     %s
`, defaultAmount, defaultCommissionRate, defaultCommissionMaxRate, defaultCommissionMaxChangeRate, defaultMinSelfDelegation),
		RunE: func(cmd *cobra.Command, args []string) error {

			ip, _ := server.ExternalIP()
			config := ctx.Config
			config.SetRoot(viper.GetString(tmcli.HomeFlag))
			nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(ctx.Config)
			if err != nil {
				return err
			}

			// Read --nodeID, if empty take it from priv_validator.json
			if nodeIDString := viper.GetString(cli.FlagNodeID); nodeIDString != "" {
				nodeID = nodeIDString
			}

			viper.SetDefault("ip", ip)
			if ip == "" {
				fmt.Fprintf(os.Stderr, "couldn't retrieve an external IP; "+
					"the tx's memo field will be unset")
			}

			var isAccNameTrue = false
			var nmsg string
			for isAccNameTrue == false {
				if nmsg == "" {
					nmsg = "enter the name of your account"
				}
				name, err := accountnameprompt(nmsg)
				if err != nil {
					fmt.Println(color.HiYellowString("something went wrong while entering name", err))
					nmsg = "re-enter the name of your account"
					isAccNameTrue = false
				}
				if len(name) > 0 {
					viper.SetDefault("name", name)
					isAccNameTrue = true
				} else {
					nmsg = "re-enter the name of your account"
					isAccNameTrue = false
				}
			}

			keyPath := filepath.Join(config.RootDir, "keys")
			kb := kbkeys.New("keys", keyPath)
			name := viper.GetString(client.FlagName)
			key, err := kb.Get(name)
			if err != nil {
				return err
			}
			key.GetAddress()

			var ischainIdTrue = false
			var cimsg string
			var chain string

			for ischainIdTrue == false {
				if cimsg == "" {
					cimsg = "enter the name for chainId"
				}
				chainvalue, err := chainIdprompt(cimsg)
				if err != nil {
					fmt.Println(color.HiYellowString("something went wrong while entering chainid", err))
					cimsg = "re-enter the name of chainid"
					ischainIdTrue = false
				}
				if len(chainvalue) > 0 {
					chain = chainvalue
					ischainIdTrue = true
				} else {
					cimsg = "re-enter the name of you are chainid"
					ischainIdTrue = false
				}
			}

			var ismonikerTrue = false
			var mmsg string
			var moniker string
			for ismonikerTrue == false {
				if mmsg == "" {
					mmsg = "enter the name for moniker"
				}
				monikervalue, err := monikerprompt(mmsg)
				if err != nil {
					fmt.Println(color.HiYellowString("something went wrong while moniker", err))
					mmsg = "re-enter the name of moniker"
					ismonikerTrue = false
				}
				if len(monikervalue) > 0 {
					moniker = monikervalue
					ismonikerTrue = true
				} else {
					mmsg = "re-enter the name of you are moniker"
					ismonikerTrue = false
				}
			}

			// Read --pubkey, if empty take it from priv_validator.json
			if valPubKeyString := viper.GetString(cli.FlagPubKey); valPubKeyString != "" {
				valPubKey, err = sdk.GetConsPubKeyBech32(valPubKeyString)
				if err != nil {
					return err
				}
			}

			// prompt for staking cin
			var isStakingTrue = false
			var smsg string
			for isStakingTrue == false {
				if smsg == "" {
					smsg = "Enter amount want to stake,please make sure you enter proper amount to stake {Eg: 22cin}: "
				}
				stakeamount, err := stakeamount(smsg)
				if err != nil {
					fmt.Println("something went wrong while staking amount", err)
					isStakingTrue = false
				}

				if len(stakeamount) <= 0 {
					smsg = "please re-enter staking amount {Eg: 22cin}:"
					isStakingTrue = false

				} else {

					_, err := sdk.ParseCoins(stakeamount)
					if err != nil {
						smsg = ("you didn't specify coin type cin please re-enter stake amount {Eg: 22cin}")
						isStakingTrue = false
					} else {
						viper.SetDefault("amount", stakeamount)
						color.HiGreen("you enter to stake amount:%s", stakeamount)
						isStakingTrue = true
					}

				}
			}

			// set for setting DefaultCommissionRate
			var iscommissionRateTrue = false
			var crmsg string
			for iscommissionRateTrue == false {
				if crmsg == "" {
					crmsg = "Enter default commission rate {Eg: 0.1} : "
				}
				defaultcommisionrate, err := DefaultCommissionRatePrompt(crmsg)
				if err != nil {
					crmsg = ("something went wrong with default commision rate please re-enter")
					iscommissionRateTrue = false
				}
				if len(defaultcommisionrate) <= 0 {
					crmsg = "please re-enter default commision rate  {Eg: 0.1}"
					iscommissionRateTrue = false

				} else {
					_, err := strconv.ParseFloat(defaultcommisionrate, 64)
					if err != nil {
						crmsg = "please re-enter default commision rate, {Eg: 0.1}"
						iscommissionRateTrue = false

					} else {
						viper.SetDefault("commission-rate", defaultcommisionrate)
						color.HiGreen("you enter default commission rate:%s", defaultcommisionrate)
						iscommissionRateTrue = true
					}

				}
			}

			// set for setting DefaultCommissionMaxRate
			var iscommissionMaxRateTrue = false
			var cmrmsg string
			for iscommissionMaxRateTrue == false {
				if cmrmsg == "" {
					cmrmsg = "Enter default commission max rate {Eg: 0.2}: "
				}
				defaultcommissionMaxRate, err := DefaultCommissionMaxRatePrompt(cmrmsg)
				if err != nil {
					cmrmsg = ("something went wrong with default commision rate please re-enter")
					iscommissionMaxRateTrue = false
				}
				if len(defaultcommissionMaxRate) <= 0 {
					cmrmsg = "please re-enter default commission max rate  {Eg: 0.2}"
					iscommissionMaxRateTrue = false

				} else {
					_, err := strconv.ParseFloat(defaultcommissionMaxRate, 64)
					if err != nil {
						cmrmsg = "please re-enter default commission max rate  {Eg: 0.2}"
						iscommissionRateTrue = false

					} else {
						viper.SetDefault("commission-max-rate", defaultcommissionMaxRate)
						color.HiGreen("you enter default commission max rate:  %s", defaultcommissionMaxRate)
						iscommissionMaxRateTrue = true
					}

				}
			}

			// set for setting DefaultCommissionMaxRate
			var isCommissionMaxChangeRateTrue = false
			var cmcrmsg string
			for isCommissionMaxChangeRateTrue == false {
				if cmcrmsg == "" {
					cmcrmsg = "Enter commission max change rate {Eg: 0.01}:  "
				}
				defaultCommissionMaxChangeRate, err := DefaultCommissionMaxChangeRatePrompt(cmcrmsg)
				if err != nil {
					cmcrmsg = ("something went wrong with  commission max change rate rate please re-enter")
					isCommissionMaxChangeRateTrue = false
				}
				if len(defaultCommissionMaxChangeRate) <= 0 {
					cmcrmsg = "please re-enter commission max change rate {Eg: 0.01}:"
					isCommissionMaxChangeRateTrue = false

				} else {
					_, err := strconv.ParseFloat(defaultCommissionMaxChangeRate, 64)
					if err != nil {
						cmcrmsg = "please re-enter commission max change rate {Eg: 0.01}:"
						iscommissionRateTrue = false
					} else {
						viper.SetDefault("commission-max-change-rate", defaultCommissionMaxChangeRate)
						color.HiGreen("you enter to default commission max change rate:  %s", defaultCommissionMaxChangeRate)
						isCommissionMaxChangeRateTrue = true
					}

				}
			}

			// set for setting DefaultMinSelfDelegation
			var isMinSelfDelegationTrue = false
			var msdmsg string
			for isMinSelfDelegationTrue == false {
				if msdmsg == "" {
					msdmsg = "Enter default min self delegation {min : 1}:"
				}
				defaultMinSelfDelegation, err := DefaultMinSelfDelegationPrompt(msdmsg)
				if err != nil {
					msdmsg = ("something went wrong with min self delegation, please re-enter {min : 1}")
					isMinSelfDelegationTrue = false
				}
				if len(defaultMinSelfDelegation) <= 0 {
					msdmsg = "please re-enter minimum self delegation {min : 1}"
					isMinSelfDelegationTrue = false

				} else {
					viper.SetDefault("min-self-delegation", defaultMinSelfDelegation)
					color.HiGreen("you enter  default min self delegation:  %s", defaultMinSelfDelegation)
					isMinSelfDelegationTrue = true
				}
			}

			web, err := prompt("enter the website url {optional}")
			if err != nil {
				fmt.Println()

			}
			if web != "" {
				viper.SetDefault("website", web)
			}

			Description, err := prompt("enter the description {optional}")
			if err != nil {
				fmt.Println()

			}
			if Description != "" {
				viper.SetDefault("details", Description)
			}

			Identity, err := prompt("enter the identity {optional}")
			if err != nil {
				fmt.Println()

			}
			if Identity != "" {
				viper.SetDefault("identity", Identity)
			}

			website := viper.GetString(cli.FlagWebsite)
			details := viper.GetString(cli.FlagDetails)
			identity := viper.GetString(cli.FlagIdentity)

			fmt.Println(website)
			fmt.Println(details)
			fmt.Println(identity)

			// Set flags for creating gentx
			prepareFlagsForTxCreateValidator(config, nodeID, ip, chain, valPubKey, website, details, identity, moniker)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc)).WithKeybase(kb)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			viper.Set(client.FlagGenerateOnly, true)

			// create a 'create-validator' message
			txBldr, msg, err := cli.BuildCreateValidatorMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}

			info, err := txBldr.Keybase().Get(name)
			if err != nil {
				return err
			}

			if info.GetType() == kbkeys.TypeOffline || info.GetType() == kbkeys.TypeMulti {
				fmt.Println("Offline key passed in. Use `mxwcli tx sign` command to sign:")
				return utils.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg})
			}

			// write the unsigned transaction to the buffer
			w := bytes.NewBuffer([]byte{})
			cliCtx = cliCtx.WithOutput(w)

			if err = utils.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg}); err != nil {
				return err
			}

			// read the transaction
			stdTx, err := readUnsignedGenTxFile(cdc, w)
			if err != nil {
				return err
			}

			// mxw msgs are gas free
			stdTx.Fee.Gas = 0
			stdTx.Fee.Amount, _ = types.ParseCoins("0cin")

			// sign the transaction and write it to the output file
			signedTx, err := utils.SignStdTx(txBldr, cliCtx, name, stdTx, false, true)
			if err != nil {
				return err
			}

			outputDocument := viper.GetString(client.FlagOutputDocument)

			if outputDocument == "" {
				outputDocument, err = makeOutputFilepath(config.RootDir, nodeID)
				if err != nil {
					return err
				}
			}
			if err := writeSignedGenTx(cdc, outputDocument, signedTx); err != nil {
				return err
			}

			//get the validator public key, Ip ,node id and save it to disk
			valpub, err := sdk.Bech32ifyConsPub(valPubKey)
			if err != nil {
				return err
			}

			var val = Createvalidator{NodeID: nodeID, IpAddress: ip, ValConPubkey: valpub, ChainId: chain, Moniker: moniker}

			err = common.EnsureDir(config.RootDir+"/developer", 0777)
			if err != nil {
				return errors.Wrap(err, "Failed to create developer folder")
			}

			src := config.RootDir + "/config/gentx/gentx-" + nodeID + ".json"
			dst := config.RootDir + "/developer/gentx-" + nodeID + ".json"
			err = CopyFile(src, dst)
			if err != nil {
				return err
			}

			srcdev := config.RootDir + "/developer/val_info.json"

			//Remove the file
			if config.RootDir == viper.GetString(tmcli.HomeFlag) {
				configpath := config.RootDir + "/config/config.toml"
				apppath := config.RootDir + "/config/app.toml"
				// data := config.RootDir + "/data"
				fileexists(configpath, apppath)

			}

			//copy the dir
			src = config.RootDir + "/config/"
			dst = config.RootDir + "/validator/config"
			err = CopyDir(src, dst)
			if err != nil {
				return err
			}

			//copy the dir
			src = config.RootDir + "/private_key"
			dst = config.RootDir + "/validator/private_key/"
			err = CopyDir(src, dst)
			if err != nil {
				return err
			}

			//write the validator key info
			err = writefile(srcdev, val)
			if err != nil {
				fmt.Println("writng the node file ", err)
			}

			//delete the file
			if config.RootDir == viper.GetString(tmcli.HomeFlag) {
				configpath := config.RootDir + "/config"
				keypath := config.RootDir + "/keys"
				data := config.RootDir + "/data"
				privatekey := config.RootDir + "/private_key"
				EnsuredeleteDir(configpath, 0777)
				EnsuredeleteDir(keypath, 0777)
				EnsuredeleteDir(data, 0777)
				EnsuredeleteDir(privatekey, 0777)
			}

			fmt.Println("Successfully transcation are created")
			return nil

		},
	}
	ip, _ := server.ExternalIP()
	viper.SetDefault("ip", ip)
	cmd.Flags().String(cli.FlagIP, ip, "The node's public IP")
	return cmd
}
