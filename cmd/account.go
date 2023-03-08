package cmd

import (
	"fmt"
	"strings"

	"github.com/Moonyongjung/xpla-anchor/types"
	"github.com/Moonyongjung/xpla-anchor/util"
	xtypes "github.com/Moonyongjung/xpla.go/types"
	"github.com/spf13/cobra"
)

// Get the account info.
func AccountCmd(a *types.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account",
		Aliases: []string{"acc"},
		Short:   "get info of the anchor account",
	}

	cmd.AddCommand(
		balance(a),
		info(a),
	)
	return cmd
}

// Query balances of the anchor account.
func balance(a *types.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance",
		Short: "query balances of the account",
		Args:  withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s query account balance
$ %s q acc balance
		`, defaultAppName, defaultAppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return util.LogErr(types.ErrAccount, err)
			}

			util.LogInfo(util.Y("(need passphrase. not tx, only using gen address from the key)"))
			_, addr, err := extractKey(home)
			if err != nil {
				return util.LogErr(types.ErrAccount, err)
			}

			msg := xtypes.BankBalancesMsg{
				Address: addr,
			}

			res, err := a.PubClient.BankBalances(msg).Query()
			if err != nil {
				return util.LogErr(types.ErrAccount, err)
			}

			util.LogInfo(res)
			util.LogInfo(util.BB("response data successfully"))
			return nil

		},
	}

	return cmd
}

// Query the anchor account
func info(a *types.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "query info of the account",
		Args:  withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s query account info
$ %s q acc info
		`, defaultAppName, defaultAppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return util.LogErr(types.ErrAccount, err)
			}

			util.LogInfo(util.Y("(need passphrase. not tx, only using gen address from the key)"))
			_, addr, err := extractKey(home)
			if err != nil {
				return util.LogErr(types.ErrAccount, err)
			}

			msg := xtypes.QueryAccAddressMsg{
				Address: addr,
			}

			res, err := a.PubClient.AccAddress(msg).Query()
			if err != nil {
				return util.LogErr(types.ErrAccount, err)
			}

			util.LogInfo(res)
			util.LogInfo(util.BB("response data successfully"))
			return nil

		},
	}

	return cmd
}
