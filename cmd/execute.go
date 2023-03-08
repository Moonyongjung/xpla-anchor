package cmd

import (
	"github.com/Moonyongjung/xpla-anchor/types"
	"github.com/Moonyongjung/xpla-anchor/util"
	"github.com/spf13/cobra"
)

// Execute the anchor.
// The execute command is able to send execute transaction of the contract or running the anchor gateway.
func ExecuteCmd(a *types.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "execute",
		Aliases: []string{"e"},
		Short:   "execute anchor",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return util.LogErr(types.ErrParseConfig, err)
			}

			appFilePath, err := getAppFile(home)
			if err != nil {
				return util.LogErr(types.ErrParseApp, err)
			}

			var channels types.Channels
			channels.AnchringTx = make(chan types.Anchoring)
			channels.HttpClientStartSignal = make(chan bool)

			pubClient, privClient, err := initXplaClient(home, true)
			if err != nil {
				return err
			}

			a.PubClient = pubClient
			a.PrivClient = privClient
			a.Channels = channels
			a.AppFilePath = appFilePath

			return nil
		},
	}

	cmd.AddCommand(
		EContractCmd(a),
		StartCmd(a),
	)
	return cmd
}
