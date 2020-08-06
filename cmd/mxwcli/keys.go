/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"bufio"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/fatih/color"
	"github.com/pkg/errors"

	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/bgentry/speakeasy"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
)

const (
	flagOverwrite  = "overwrite"
	DefaultChainID = "maxonrow-chain"
)

var (
	DefaultNodeHome = os.ExpandEnv("$HOME/.mxw/")
)

// keysCmd represents the keys command
func createKeyPairCommand() *cobra.Command {
	keysCmd := &cobra.Command{
		Use:   "create-keypair",
		Short: "create the account with mnemonic, private key, public key and address",
		Long:  `create the account with mnemonic, private key, public key and address`,
		RunE: func(_ *cobra.Command, _ []string) error {

			color.HiBlue(`
| $$$    /$$$
| $$$$  /$$$$  /$$$$$$  /$$   /$$  /$$$$$$  /$$$$$$$   /$$$$$$   /$$$$$$  /$$  /$$  /$$
| $$ $$/$$ $$ |____  $$|  $$ /$$/ /$$__  $$| $$__  $$ /$$__  $$ /$$__  $$| $$ | $$ | $$
| $$  $$$| $$  /$$$$$$$ \  $$$$/ | $$  \ $$| $$  \ $$| $$  \__/| $$  \ $$| $$ | $$ | $$
| $$\  $ | $$ /$$__  $$  >$$  $$ | $$  | $$| $$  | $$| $$      | $$  | $$| $$ | $$ | $$
| $$ \/  | $$|  $$$$$$$ /$$/\  $$|  $$$$$$/| $$  | $$| $$      |  $$$$$$/|  $$$$$/$$$$/
|__/     |__/ \_______/|__/  \__/ \______/ |__/  |__/|__/       \______/  \_____/\___/
                                                                                      `)

			cobra.EnableCommandSorting = false
			ctx := server.NewDefaultContext()
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			_ = viper.GetString(client.FlagChainID)

			var isnameTrue = false
			var accname string
			var nmsg string
			for isnameTrue == false {
				if nmsg == "" {
					nmsg = "enter the name for your account"
				}
				name, err := nameprompt(nmsg)
				if err != nil {
					fmt.Println(color.HiYellowString("something went wrong while entering name", err))
					nmsg = "re-enter the name for your account"
					isnameTrue = false
				}
				if len(name) > 0 {
					accname = name
					isnameTrue = true
				} else {
					nmsg = "re-enter the name for your account"
					isnameTrue = false
				}
			}
			fmt.Println(" ")

			//password for account
			message := `**Important** please keep the password safe.
password must contain one uppercase letter
password must have length of seven character
password must contain one number
password must contain lowercase
password must contain special characters`

			color.HiRed(message)

			var ispasswordTrue = false
			var password string
			var pssmsg string
			for ispasswordTrue == false {
				if pssmsg == "" {
					pssmsg = color.CyanString("please enter password for your account:")
				}
				passphrase, err := prompt(pssmsg)
				if err != nil {
					fmt.Println(color.HiYellowString("something went wrong while entering password", err))
					pssmsg = "re-enter  password :"
					ispasswordTrue = false
				} else {
					valid := isValid(passphrase)
					if valid == true {
						pssmsg = color.CyanString("please confirm password for your account :")
						confirmPassword, err := prompt(pssmsg)
						if err != nil {
							fmt.Println(color.HiYellowString("password mismath please try again"))
							pssmsg = "re-enter  password :"
							ispasswordTrue = false
						}
						if passphrase != confirmPassword {
							fmt.Println(color.HiYellowString("password and confirmPassword mismath please try again"))
							pssmsg = "re-enter  password :"
							ispasswordTrue = false
						} else {
							password = passphrase
							ispasswordTrue = true
						}
					} else {
						fmt.Println(color.HiYellowString("Invalid password, password must match the condition"))
						pssmsg = "re-enter  password :"
						ispasswordTrue = false
					}
				}
			}

			fmt.Println(" ")

			keyPath := filepath.Join(config.RootDir, "keys")
			keybase := keys.New("keys", keyPath)
			info, mnemonic, err := keybase.CreateMnemonic(accname, keys.English, password, keys.Secp256k1)
			imp := "\n**Important** write this mnemonic phrase in a safe place.It is the only way to recover your account if you ever forget your password."
			color.HiRed(imp)
			fmt.Println(" ")
			fmt.Println(mnemonic)
			seed := strings.Split(mnemonic, " ")
			mn := make([]string, len(seed))
			//new slice vaiable
			copy(mn, seed)
			clearscreen()
			for i := range seed {
				seed[i] = strconv.Itoa(i+1) + "." + seed[i]
			}
			IterSplits(SliceSplit(seed, 6), 6)
			fmt.Println(" ")
			color.HiYellow("Please write down mnemonic, and press 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')

			fmt.Println(" ")
			var isTrue = false
			for isTrue == false {
				sd, err := readseed()
				if err != nil {
					return fmt.Errorf("Unable read the seed %v", err)
				}
				rs := strings.Split(sd, " ")
				isTrue = reflect.DeepEqual(rs, mn)
				if isTrue == true {
					isTrue = true

				}
			}

			//encrypt the priv_key file with password
			encryptmessage := `**Important** encrypt password use to encrypt the private_key,please keep the password safe.
passsword must contain one uppercase letter
passsword must contain special characters
passsword must be length of seven character
passsword must contain one number
passsword must contain lowercase`
			color.HiRed(encryptmessage)
			var isEncryptTrue = false
			var encryptPhrase string
			var encryptmsg string
			for isEncryptTrue == false {
				if encryptmsg == "" {
					encryptmsg = color.CyanString("please enter password to encrypt private_key :")
				}
				encryptpassphrase, err := prompt(encryptmsg)
				if err != nil {
					fmt.Println(color.HiYellowString("something went wrong while entering encrypt_key", err))
					encryptmsg = "re-enter  password to encrypt private_key :"
					isEncryptTrue = false
				} else {
					encryptvalid := isValid(encryptpassphrase)
					if encryptvalid == true {
						encryptmsg = color.CyanString("please enter confirm password to encrypt private_key :")
						encryptconfirmPassword, err := prompt(encryptmsg)
						if err != nil {
							fmt.Println(color.HiYellowString("password mismath please try again"))
							encryptmsg = "re-enter  password to encrypt private_key :"
							isEncryptTrue = false
						}
						if encryptpassphrase != encryptconfirmPassword {
							fmt.Println(color.HiYellowString("password and confirmPassword mismath please try again"))
							encryptmsg = "re-enter  password to encrypt private_key :"
							isEncryptTrue = false
						} else {
							encryptPhrase = encryptpassphrase
							isEncryptTrue = true
						}
					} else {
						fmt.Println(color.HiYellowString("invalid password please try again"))
						encryptmsg = "re-enter  password to encrypt private_key :"
						isEncryptTrue = false
					}
				}
			}

			//export private key
			privkey, err := keybase.ExportPrivKey(accname, password, encryptPhrase)
			if err != nil {
				return fmt.Errorf("Unable export key %v", err)
			}
			err = common.EnsureDir(config.RootDir+"/private_key", 0777)
			if err != nil {
				return errors.Wrap(err, "Failed to create developer folder")
			}

			priv_key_file := config.RootDir + "/private_key/priv_key"
			dataBytes := []byte(privkey)
			err = writefile(priv_key_file, dataBytes)
			if err != nil {
				fmt.Println("writng the private key ", err)
			}

			fmt.Println(" ")
			fmt.Println(color.CyanString("your wallet address:"), info.GetAddress().String())
			fmt.Println(" ")
			fmt.Println(color.CyanString("your account Name:"), info.GetName())
			fmt.Println(" ")
			fmt.Println(color.CyanString("your encrypted private key has been saved in:"), priv_key_file)
			config.ProfListenAddress = ""
			return nil
		},
	}

	keysCmd.Flags().String(cli.HomeFlag, DefaultNodeHome, "node's home directory")
	keysCmd.Flags().String(client.FlagChainID, DefaultChainID, "genesis file chain-id")

	return keysCmd
}

//prompt for password
func prompt(msg string) (pass string, err error) {
	password, err := speakeasy.Ask(msg)
	if err != nil {
		fmt.Println(err)

	}
	return password, nil
}

//prompt for enter the seeds
func readseed() (seed string, err error) {
	CallClear()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(color.CyanString("Please verify the seeds dispayed above: "))
	scanner.Scan()
	text := scanner.Text()
	return text, nil
}

/* prompt name*/
func nameprompt(nmsg string) (seed string, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	msg := color.CyanString(nmsg)
	fmt.Println(msg)
	scanner.Scan()
	text := scanner.Text()
	return text, nil
}

/* dealy func before ask  user input seed*/
func DelaySecond(n time.Duration) {
	time.Sleep(n * time.Second)
}

/*check the  valid password
    `password must contain one uppercase letter
	 must have length of seven
	 must contain one number
	 must contain lowercase` */
func isValid(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

/* arrange the seed in list with row and column*/
func SliceSplit(aSlice []string, splitSize int) [][]string {

	var splits [][]string

	for i := 0; i < len(aSlice); i += splitSize {
		end := i + splitSize

		if end > len(aSlice) {
			end = len(aSlice)
		}

		splits = append(splits, aSlice[i:end])
	}

	return splits
}

func IterSplits(slices [][]string, splitSize int) {
	for i := 0; i < splitSize; i++ {
		for _, s := range slices {
			if len(s) > i {
				fmt.Printf("%-15v", s[i])
			}
		}
		println("")
	}
}

//clear the terminal

var clear map[string]func() //create a map for storing clear funcs

func clearscreen() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

//Take slice bytes and create the output file to disk
func writefile(outputDocument string, value []byte) error {
	outputFile, err := os.OpenFile(outputDocument, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("error cause due to open file from disk", err)
	}
	defer outputFile.Close()
	_, err = fmt.Fprintf(outputFile, "%s", value)
	if err != nil {
		fmt.Println("error cause while writing the private key to  disk", err)
	}
	return nil
}
