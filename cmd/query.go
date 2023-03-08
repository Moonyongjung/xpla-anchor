package cmd

import (
	"fmt"
	"strings"

	"github.com/Moonyongjung/xpla-anchor/app"
	"github.com/Moonyongjung/xpla-anchor/gw"
	"github.com/Moonyongjung/xpla-anchor/types"
	"github.com/Moonyongjung/xpla-anchor/util"
	xtypes "github.com/Moonyongjung/xpla.go/types"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

const (
	verified    = "-verifed-"
	notVerified = "NOT CONSISTENT, check state of the private chain"
)

// Query the anchor.
// The query command is able to get responses of the contract.
func QueryCmd(a *types.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "query anchor",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return util.LogErr(types.ErrParseConfig, err)
			}

			appFilePath, err := getAppFile(home)
			if err != nil {
				return util.LogErr(types.ErrParseApp, err)
			}

			pubClient, privClient, err := initXplaClient(home, false)
			if err != nil {
				return err
			}

			a.PubClient = pubClient
			a.PrivClient = privClient
			a.AppFilePath = appFilePath

			return nil
		},
	}

	cmd.AddCommand(
		QContractCmd(a),
		AccountCmd(a),
		verify(a),
	)
	return cmd
}

// verify by comparing recorded block info in the contract with response of the private chain.
func verify(a *types.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "verify [height]",
		Aliases: []string{"v"},
		Short:   "verify the block consistency of the private chain",
		Args:    withUsage(cobra.ExactArgs(1)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s query verify [height]
$ %s q v [height]
$ %s q v [height] --address [contract_address] --priv-block-api [blockinfo_api_of_private_chain] 
		`, defaultAppName, defaultAppName, defaultAppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set the API which can check the block info.
			blockApi, err := cmd.Flags().GetString(flagPrivBlockApi)
			if err != nil {
				return util.LogErr(types.ErrQuery, err)
			}

			// Get the block info from the private chain.
			res, err := gw.DoRequest(a, blockApi, args[0], false)
			if err != nil {
				return util.LogErr(types.ErrQuery, err)
			}

			var block types.Block
			responseData := util.JsonUnmarshalData(&block, res)
			mapstructure.Decode(responseData, &block)

			addr, err := cmd.Flags().GetString(flagContractAddr)
			if err != nil {
				return util.LogErr(types.ErrQuery, err)
			}

			if addr == "" {
				addr = app.AppFile().Get().Contract.Address
			}

			height := args[0]

			// Generate the query message.
			queryMsg := xtypes.QueryMsg{
				ContractAddress: addr,
				QueryMsg:        `{"block_data":{"height":"` + height + `"}}`,
			}

			// Request query by using XPLA client.
			resContract, err := a.PubClient.QueryContract(queryMsg).Query()
			if err != nil {
				return util.LogErr(types.ErrContract, err)
			}

			var blockInfo types.QueryBlockInfoResponse
			resContractData := util.JsonUnmarshalData(&blockInfo, []byte(resContract))
			mapstructure.Decode(resContractData, &blockInfo)

			// Compare.
			util.LogInfo("[priv chain]", util.BB("height=")+block.Block.Header.Height)
			util.LogInfo("[contract]  ", util.BB("height=")+blockInfo.Data.Height)
			if block.Block.Header.Height == blockInfo.Data.Height {
				util.LogInfo(util.G("height " + verified))
			} else {
				util.LogWarning(util.R(notVerified))
			}

			util.LogInfo("[priv chain]", util.BB("block hash=")+block.BlockID.Hash)
			util.LogInfo("[contract]  ", util.BB("block hash=")+blockInfo.Data.BlockHash)
			if block.BlockID.Hash == blockInfo.Data.BlockHash {
				util.LogInfo(util.G("hash " + verified))
			} else {
				util.LogWarning(util.R(notVerified))
			}

			util.LogInfo("[priv chain]", util.BB("merkle root=")+block.Block.Header.DataHash)
			util.LogInfo("[contract]  ", util.BB("merkle root=")+blockInfo.Data.DataMerkle)
			if block.Block.Header.DataHash == blockInfo.Data.DataMerkle {
				util.LogInfo(util.G("merkle " + verified))
			} else {
				util.LogWarning(util.R(notVerified))
			}

			util.LogInfo("[priv chain]", util.BB("timestamp=")+block.Block.Header.Time)
			util.LogInfo("[contract]  ", util.BB("timestamp=")+blockInfo.Data.Timestamp)
			if block.Block.Header.Time == blockInfo.Data.Timestamp {
				util.LogInfo(util.G("timestamp " + verified))
			} else {
				util.LogWarning(util.R(notVerified))
			}

			return nil
		},
	}
	cmd.Flags().String(flagContractAddr, "", "address of the anchor contract")
	cmd.Flags().String(flagPrivBlockApi, defaultPrivBlockApi, "block query API of the private chain(except height)")

	return cmd
}
