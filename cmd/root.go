package cmd

import (
	"context"
	"errors"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/Moonyongjung/xpla-anchor/app"
	"github.com/Moonyongjung/xpla-anchor/types"
	"github.com/Moonyongjung/xpla-anchor/util"
	"github.com/Moonyongjung/xpla.go/client"
	"github.com/Moonyongjung/xpla.go/key"
	xutil "github.com/Moonyongjung/xpla.go/util"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var (
	// default values of the anchor
	defaultAppName      = "anc"
	defaultHome         = filepath.Join(os.Getenv("HOME"), ".anchor")
	defaultAppPath      = "config"
	defaultKeyPath      = "keys"
	defaultAppFilePath  = "app.yaml"
	defaultConfigPath   = "./config.yaml"
	defaultKeyName      = "AnchorKey"
	defaultPrivBlockApi = "/blocks"
)

// Start the command for the anchor.
// The anchor provides CLI functions and running anchor module.
func Start() {
	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd()
	rootCmd.SilenceUsage = true

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

// New root command for the anchor.
// In order to run the anchor, implement commands
// such as intialization, execute and query to the anchor contract.
func NewRootCmd() *cobra.Command {
	a := &types.App{
		Viper: viper.New(),
	}

	var rootCmd = &cobra.Command{
		Use:   defaultAppName,
		Short: "The anchor provides secure to the private chain by sending tx to main chain as XPLA",
		Long: util.ToString(`the anchor provides functions as below

1. Initialize the anchor by reading config.yaml file which an user set
2. Recover the key to send transaction to the main chain as XPLA
3. Easily execute and query to the anchor contract such as store, instantiate and query
4. Run the anchor gateway
5. Verify the consistency between block info in the anchor contract and query response from the private chain
		`, ""),
	}

	rootCmd.PersistentFlags().StringVar(&a.HomePath, flagHome, defaultHome, "set home directory")
	if err := a.Viper.BindPFlag(flagHome, rootCmd.PersistentFlags().Lookup(flagHome)); err != nil {
		panic(err)
	}

	// Register subcommands
	rootCmd.AddCommand(
		InitCmd(),
		KeyCmd(),
		ConfigCmd(),
		ExecuteCmd(a),
		QueryCmd(a),
	)

	return rootCmd
}

// Initialize XPLA client in order to send a transaction or execute/query contract.
// Chain ID, LCD URL of public and private chain is mandatory, but parameters to generate transaction is optional.
// The data field is structed in the config.yaml file
func initXplaClient(home string, isExecute bool) (*client.XplaClient, *client.XplaClient, error) {
	conf := app.AppFile().Get().Config

	// mandatory
	chainId := conf.PublicChain.ChainID
	if chainId == "" {
		return nil, nil, util.LogErr(types.ErrGenXplaClient, "the config must include chain ID")
	}
	lcd := conf.PublicChain.LCD
	if lcd == "" {
		return nil, nil, util.LogErr(types.ErrGenXplaClient, "the config must include LCD URL")
	}

	var privKey cryptotypes.PrivKey

	// The private key is needed to only send a transaction
	// When the user queries to the anchor contract, the private key is not used
	if isExecute {
		extPrivKey, addr, err := extractKey(home)
		if err != nil {
			return nil, nil, util.LogErr(types.ErrGenXplaClient)
		}
		privKey = extPrivKey

		util.LogInfo("public chain account to send tx=" + addr)
	}

	// optional params(mode, gas adjustment, gas limit)
	broadcastMode := ""
	if conf.PublicChain.BroadcastMode != "" {
		broadcastMode = conf.PublicChain.BroadcastMode
	}

	gasAdj := ""
	if conf.PublicChain.GasAdj != "" {
		gasAdj = conf.PublicChain.GasAdj
	}

	gasLimit := ""
	if conf.PublicChain.GasLimit != "" {
		gasLimit = conf.PublicChain.GasLimit
	}

	pubXplac := client.NewXplaClient(chainId).WithOptions(
		client.Options{
			LcdURL:        lcd,
			BroadcastMode: broadcastMode,
			GasAdjustment: gasAdj,
			GasLimit:      gasLimit,
		},
	)

	// include the private key in the XPLA client when send a transaction.
	if isExecute {
		pubXplac.WithPrivateKey(privKey)
	}

	// mandatory
	privChainId := conf.PrivateChain.ChainID
	if privChainId == "" {
		return nil, nil, util.LogErr(types.ErrGenXplaClient, "the config must include chain ID")
	}
	privLcd := conf.PrivateChain.LCD
	if privLcd == "" {
		return nil, nil, util.LogErr(types.ErrGenXplaClient, "the config must include LCD URL")
	}

	priXplac := client.NewXplaClient(privChainId).WithURL(privLcd)
	util.LogInfo("generate XPLA client successfully")

	return pubXplac, priXplac, nil
}

// Extract the private key
func extractKey(home string) (cryptotypes.PrivKey, string, error) {
	keyFileDir := path.Join(home, defaultKeyPath)
	keyFilePath := path.Join(keyFileDir, defaultKeyName)
	if _, err := os.Stat(keyFilePath); os.IsNotExist(err) {
		return nil, "", util.LogErr(types.ErrGenXplaClient, "run recovering key before execute")
	}

	keyFile, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, "", util.LogErr(types.ErrKey, err)
	}

	xutil.MakeEncodingConfig()
	util.LogInfo("input passphrase")
	passphraseByte, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, "", util.LogErr(types.ErrGenXplaClient, err)
	}

	// The private key is saved by armored type
	privKey, _, err := key.UnarmorDecryptPrivKey(string(keyFile), string(passphraseByte))
	if err != nil {
		return nil, "", util.LogErr(types.ErrGenXplaClient, err)
	}

	addr, err := key.Bech32AddrString(privKey)
	if err != nil {
		return nil, "", util.LogErr(types.ErrGenXplaClient)
	}

	return privKey, addr, nil
}

// Wrapping the args using CLI.
func withUsage(inner cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := inner(cmd, args); err != nil {
			cmd.Root().SilenceUsage = false
			cmd.SilenceUsage = false
			return err
		}

		return nil
	}
}

// Get instance of the app.yaml file
func getAppFile(home string) (string, error) {
	appFileDir := path.Join(home, defaultAppPath)
	appFilePath := path.Join(appFileDir, defaultAppFilePath)
	if _, err := os.Stat(appFilePath); os.IsNotExist(err) {
		return "", errors.New("invalid request, run init first")
	}

	err := app.AppFile().Read(appFilePath)
	if err != nil {
		return "", err
	}

	return appFilePath, nil
}
