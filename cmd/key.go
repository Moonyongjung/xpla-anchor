package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/Moonyongjung/xpla-anchor/types"
	"github.com/Moonyongjung/xpla-anchor/util"
	"github.com/Moonyongjung/xpla.go/key"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// Manage the private key for handling the transaction.
// The anchor does not create a new key, but it is able to recover by mnemonic.
// The mnemonic words is not saved in the local directory as home dir,
// and the private key is recorded in home directory by armored type.
// The key is encrypted by using passphrase when use key commands.
func KeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "key",
		Aliases: []string{"k"},
		Short:   "Manage the private key for the anchor",
	}

	cmd.AddCommand(
		recover(),
		change(),
	)
	return cmd
}

// Recovering the private key by using mnemonic words.
func recover() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recover",
		Short: `recover key by using mnemonic words. the anchor can own only one key`,
		Args:  withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s key recover
		`, defaultAppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				util.LogErr(types.ErrInit, err)
				return err
			}

			err = genKey(home, true)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

// Change the private key.
// Should implement the change command to switch the key because The anchor uses only the one private key.
func change() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change",
		Short: `change key by using mnemonic words. the anchor can own only one key`,
		Args:  withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s key change
		`, defaultAppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				util.LogErr(types.ErrInit, err)
				return err
			}

			err = genKey(home, false)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

// Generate the private key.
func genKey(home string, isRecover bool) error {
	appFileDir := path.Join(home, defaultAppPath)
	appFilePath := path.Join(appFileDir, defaultAppFilePath)
	if _, err := os.Stat(appFilePath); os.IsNotExist(err) {
		util.LogErr(types.ErrKey, "run init before gen key")
		return err
	}

	// Default key name.
	keyFileDir := path.Join(home, defaultKeyPath)
	keyFilePath := path.Join(keyFileDir, defaultKeyName)

	if isRecover {
		if err := os.Mkdir(keyFileDir, os.ModePerm); err != nil {
			util.LogErr(types.ErrKey, err)
			return err
		}
	}

	util.LogInfo("input BIP39 mnemonic\n")
	mnemonicByte, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		util.LogErr(types.ErrKey, err)
		return err
	}

	mnemonic := string(mnemonicByte)

	if !bip39.IsMnemonicValid(mnemonic) {
		err := "invalid mnemonic"
		util.LogErr(types.ErrKey, err)
		return errors.New(err)
	}

	util.LogInfo("input passphrase\n")
	passphraseByte, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		util.LogErr(types.ErrKey, err)
		return err
	}

	passphrase := string(passphraseByte)

	privKey, err := key.NewPrivKey(mnemonic)
	if err != nil {
		util.LogErr(types.ErrKey, err)
		return err
	}

	addr, err := key.Bech32AddrString(privKey)
	if err != nil {
		util.LogErr(types.ErrKey, err)
		return err
	}
	util.LogInfo("address=" + addr)

	// generate the armored private key which is encrypted by using passphrase.
	armoredKey := key.EncryptArmorPrivKey(privKey, passphrase)

	f, err := os.Create(keyFilePath)
	if err != nil {
		util.LogErr(types.ErrKey, err)
		return err
	}
	defer f.Close()

	if _, err = f.Write([]byte(armoredKey)); err != nil {
		util.LogErr(types.ErrKey, err)
		return err
	}

	util.LogInfo(util.BB("success recover key"))

	return nil
}
