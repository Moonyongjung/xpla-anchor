package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Moonyongjung/xpla-anchor/app"
	"github.com/Moonyongjung/xpla-anchor/types"
	"github.com/Moonyongjung/xpla-anchor/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Handle the config.yaml file.
// If the owner of the anchor need to change the config params in the app.yaml,
// run the config command.
func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"conf"},
		Short:   "handle the config file",
	}

	cmd.AddCommand(
		set(),
	)
	return cmd
}

// Change the config part of app file.
func set() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set",
		Aliases: []string{"s"},
		Short:   "modify the app file by changing the config file",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s config set 
$ %s conf s		
		`, defaultAppName, defaultAppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return util.LogErr(types.ErrInit, err)
			}

			configFile, err := cmd.Flags().GetString(flagConfigFile)
			if err != nil {
				return util.LogErr(types.ErrInit, err)
			}

			appFileDir := path.Join(home, defaultAppPath)
			appFilePath := path.Join(appFileDir, defaultAppFilePath)
			if _, err := os.Stat(appFilePath); os.IsNotExist(err) {
				return util.LogErr(types.ErrParseApp, "run init before set config params")
			}

			configYamlFile, err := os.ReadFile(configFile)
			if err != nil {
				return util.LogErr(types.ErrParseConfig, err)
			}

			appYamlFile, err := os.ReadFile(appFilePath)
			if err != nil {
				return util.LogErr(types.ErrParseConfig, err)
			}

			var configType app.ConfigType
			err = yaml.Unmarshal(configYamlFile, &configType)
			if err != nil {
				return util.LogErr(types.ErrParseConfig, err)
			}

			var appType app.AppType
			err = yaml.Unmarshal(appYamlFile, &appType)
			if err != nil {
				return util.LogErr(types.ErrParseConfig, err)
			}

			err = appType.SetConfig(configType, appFilePath)
			if err != nil {
				return util.LogErr(types.ErrParseConfig, err)
			}

			util.LogInfo(util.BB("set configuration to app file successfully"))

			return nil
		},
	}
	cmd.Flags().String(flagConfigFile, defaultConfigPath, "configuration file path")

	return cmd
}
