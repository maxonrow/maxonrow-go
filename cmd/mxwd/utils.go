package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	tm "github.com/tendermint/tendermint/types"
)

//Take the input from the terminal for creating the node
func createNode() (int, error) {
	return 1, nil
	var n int
	fmt.Println("Enter an integer value for node creation  : ")
	_, err := fmt.Scanf("%d", &n)
	if err != nil {
		fmt.Println(err)
	}
	if n > 30 {
		fmt.Println("please enter below 15")
		os.Exit(1)
	}

	fmt.Println("You have entered : ", n)
	return n, nil
}

//check the config and data file exist and remove
func fileexists(configpath string, datapath string) error {
	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		return err
	}

	err := os.RemoveAll(configpath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(datapath); os.IsNotExist(err) {
		return err
	}

	err = os.RemoveAll(datapath)
	if err != nil {
		return err
	}

	return nil
}

func folderexists(configpath string, datapath string) error {
	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		return err
	}

	err := os.Remove(configpath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(datapath); os.IsNotExist(err) {
		return err
	}
	err = os.Remove(datapath)
	if err != nil {
		return err
	}

	return nil
}

func copyKeyNode(keyPath string, nodepath []string) error {
	for _, path := range nodepath {
		n1 := filepath.Join(path, "keys")
		err := CopyDir(keyPath, n1)
		if err != nil {
			return fmt.Errorf("Unable to create new account: %v", err)
		}
	}
	return nil
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {

	if fileExists(dst) {
		return nil
	} else {

		in, err := os.Open(src)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer func() {
			if e := out.Close(); e != nil {
				err = e
			}
		}()

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}

		err = out.Sync()
		if err != nil {
			return err
		}

		si, err := os.Stat(src)
		if err != nil {
			return err
		}
		err = os.Chmod(dst, si.Mode())
		if err != nil {
			return err
		}
	}
	return err
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func exportGenesisFile(
	genFile, chainID string, validators []tm.GenesisValidator, appState json.RawMessage,
) error {
	genDoc := tm.GenesisDoc{
		ChainID:    chainID,
		Validators: validators,
		AppState:   appState,
	}
	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}
	return genDoc.SaveAs(genFile)
}

func accountnameprompt(nmsg string) (name string, err error) {
	fmt.Println(" ")
	scanner := bufio.NewScanner(os.Stdin)
	msg := color.CyanString(nmsg)
	fmt.Println(msg)
	scanner.Scan()
	text := scanner.Text()
	return text, nil
}

func stakeamount(smg string) (amount string, err error) {
	fmt.Println(" ")
	scanner := bufio.NewScanner(os.Stdin)
	msg := color.CyanString(smg)
	fmt.Println(msg)
	scanner.Scan()
	text := scanner.Text()

	viper.SetDefault("amount", text)
	return text, nil
}

func monikerprompt(mmsg string) (moniker string, err error) {
	fmt.Println(" ")
	scanner := bufio.NewScanner(os.Stdin)
	msg := color.CyanString(mmsg)
	fmt.Println(msg)
	scanner.Scan()
	monikervalue := scanner.Text()
	return monikervalue, nil
}

func chainIdprompt(cmsg string) (chainID string, err error) {
	fmt.Println(" ")
	scanner := bufio.NewScanner(os.Stdin)
	msg := color.CyanString(cmsg)
	fmt.Println(msg)
	scanner.Scan()
	chainId := scanner.Text()
	return chainId, nil
}

func DefaultCommissionRatePrompt(crmsg string) (amount string, err error) {
	fmt.Println(" ")
	scanner := bufio.NewScanner(os.Stdin)
	msg := color.CyanString(crmsg)
	fmt.Println(msg)
	scanner.Scan()
	DefaultCommissionRate := scanner.Text()

	return DefaultCommissionRate, nil
}

func DefaultCommissionMaxRatePrompt(cmrmsg string) (amount string, err error) {
	fmt.Println(" ")
	scanner := bufio.NewScanner(os.Stdin)
	msg := color.CyanString(cmrmsg)
	fmt.Println(msg)
	scanner.Scan()
	DefaultCommissionMaxRate := scanner.Text()

	return DefaultCommissionMaxRate, nil
}

func DefaultCommissionMaxChangeRatePrompt(cmcrmsg string) (amount string, err error) {
	fmt.Println(" ")
	scanner := bufio.NewScanner(os.Stdin)
	msg := color.CyanString(cmcrmsg)
	fmt.Println(msg)
	scanner.Scan()
	DefaultCommissionMaxChangeRate := scanner.Text()

	return DefaultCommissionMaxChangeRate, nil
}

func DefaultMinSelfDelegationPrompt(msdp string) (amount string, err error) {
	fmt.Println(" ")
	scanner := bufio.NewScanner(os.Stdin)
	msg := color.CyanString(msdp)
	fmt.Println(msg)
	scanner.Scan()
	DefaultMinSelfDelegation := scanner.Text()
	return DefaultMinSelfDelegation, nil
}

func prompt(msg string) (amount string, err error) {
	fmt.Println(" ")
	scanner := bufio.NewScanner(os.Stdin)
	out := color.CyanString(msg)
	fmt.Println(out)
	scanner.Scan()
	value := scanner.Text()
	return value, nil
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

func EnsuredeleteDir(dir string, mode os.FileMode) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("Could not create directory %v. %v", dir, err)
	}
	return nil
}

func writefile(outputDocument string, value Createvalidator) error {
	outputFile, err := os.OpenFile(outputDocument, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	json, err := amino.MarshalJSON(value)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(outputFile, "%s\n", json)
	return err
}

func writeNodes(outputDocument string, value NodeDetails) error {
	//write the validator key info
	dst := outputDocument + "/Node_Pub_Id.json"
	file, err := json.MarshalIndent(value, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dst, file, 0644)
	if err != nil {
		return err
	}
	return nil
}
