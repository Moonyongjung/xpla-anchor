package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Moonyongjung/xpla-anchor/app"
	"github.com/Moonyongjung/xpla-anchor/gw"
	"github.com/Moonyongjung/xpla-anchor/gw/db"
	"github.com/Moonyongjung/xpla-anchor/types"
	"github.com/Moonyongjung/xpla-anchor/util"
	"github.com/spf13/cobra"
)

const (
	defaultLog = "none"
	dbLog      = "db"
)

// Running the anchor gateway.
// The gateway aggregates info of block in the private chain,
// and records info to the anchor contract.
func StartCmd(a *types.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"s"},
		Short:   "start anchor gateway",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s start
$ %s s
$ %s s --address [contract_address]
$ %s s --priv-block-api [blockinfo_api_of_private_chain]
$ %s s --log [db|file] 
		`, defaultAppName, defaultAppName, defaultAppName, defaultAppName, defaultAppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := cmd.Flags().GetString(flagContractAddr)
			if err != nil {
				return util.LogErr(types.ErrGw, err)
			}

			if addr == "" {
				addr = app.AppFile().Get().Contract.Address
			}

			log, err := cmd.Flags().GetString(flagLog)
			if err != nil {
				return util.LogErr(types.ErrGw, err)
			}

			if !(log == defaultLog || log == dbLog) {
				return util.LogErr(types.ErrGw, "invalid log type")
			}

			if log == dbLog {
				err = db.DbInit()
				if err != nil {
					return util.LogErr(types.ErrGw, err)
				}
			}

			// Set the API which can check the block info.
			// If the private chain has not the default block info API, the anchor need the new API.
			// e.g. default block info API such as XPLA
			//      https://LCD_URL/blocks
			//		other API such as EVMOS
			//      https://LCD_URL/cosmos/base/tendermint/v1beta1/blocks
			blockApi, err := cmd.Flags().GetString(flagPrivBlockApi)
			if err != nil {
				return util.LogErr(types.ErrGw, err)
			}

			// Thread gateway.
			go gw.SendAnchoringTx(a, log)
			go gw.StartGW(a, addr, blockApi, log)

			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
			<-stop

			util.LogInfo("shutting down the gateway...")
			util.LogInfo("gateway gracefully stopped")

			return nil
		},
	}
	cmd.Flags().String(flagContractAddr, "", "address of the anchor contract")
	cmd.Flags().String(flagPrivBlockApi, defaultPrivBlockApi, "block query API of the private chain(except height)")
	cmd.Flags().String(flagLog, defaultLog, "select log type")

	return cmd
}
