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

// Initialize the anchor.
// Read the config.yaml file and generate the app.yaml file
// the app file is used to variaty functions
func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "initialize the anchor",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s init 
$ %s i
$ %s i --home [home_dir] --config [config_file_path]		
		`, defaultAppName, defaultAppName, defaultAppName)),
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
				if _, err := os.Stat(appFileDir); os.IsNotExist(err) {
					if _, err := os.Stat(home); os.IsNotExist(err) {
						if err = os.Mkdir(home, os.ModePerm); err != nil {
							return util.LogErr(types.ErrParseConfig, err)
						}
					}
					if err = os.Mkdir(appFileDir, os.ModePerm); err != nil {
						return util.LogErr(types.ErrParseConfig, err)
					}
				}

				yamlFile, err := os.ReadFile(configFile)
				if err != nil {
					return util.LogErr(types.ErrParseConfig, err)
				}

				var configType app.ConfigType
				err = yaml.Unmarshal(yamlFile, &configType)
				if err != nil {
					return util.LogErr(types.ErrParseConfig, err)
				}

				// Generate default app file in the home directory.
				err = app.GenDefaultApp(home, appFilePath, configType)
				if err != nil {
					return util.LogErr(types.ErrParseApp, err)
				}

				util.LogInfo(util.BB("success initialization"))
				return nil
			}

			return util.LogErr(types.ErrParseApp, "app file already exists")
		},
	}
	cmd.Flags().String(flagConfigFile, defaultConfigPath, "configuration file path")

	return cmd
}
