package cmd

import (
	"fmt"
	"strings"

	"github.com/Moonyongjung/xpla-anchor/app"
	"github.com/Moonyongjung/xpla-anchor/types"
	"github.com/Moonyongjung/xpla-anchor/util"
	"github.com/Moonyongjung/xpla.go/key"
	xtypes "github.com/Moonyongjung/xpla.go/types"
	"github.com/spf13/cobra"
)

const (
	// Generate the contract on the MacOS.
	defaultContractDirPath = "./contract"
	defaultContractPath    = "/artifacts/xpla_anchor_contract-aarch64.wasm"
)

// Contract command to execute.
func EContractCmd(a *types.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "contract",
		Aliases: []string{"ctrt"},
		Short:   "execute to the anchor contract",
	}

	cmd.AddCommand(
		storeContract(a),
		instantiateContract(a),
	)
	return cmd
}

// Contract command to query.
func QContractCmd(a *types.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "contract",
		Aliases: []string{"ctrt"},
		Short:   "query to the anchor contract",
	}

	cmd.AddCommand(
		queryContractLatestBlock(a),
		queryContractBlockData(a),
	)
	return cmd
}

// Store the anchor contract to the main chain.
func storeContract(a *types.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store",
		Short: "store anchor contract",
		Args:  withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s execute contract store
$ %s e ctrt store
$ %s e ctrt store --path [contract_file_path]
		`, defaultAppName, defaultAppName, defaultAppName)),
		RunE: func(cmd *cobra.Command, args []string) error {

			path, err := cmd.Flags().GetString(flagContractFilePath)
			if err != nil {
				return util.LogErr(types.ErrContract, err)
			}
			msg := xtypes.StoreMsg{
				FilePath: path,
			}
			// Create and sign the transaction to execute.
			txByte, err := a.PubClient.StoreCode(msg).CreateAndSignTx()
			if err != nil {
				return util.LogErr(types.ErrContract, err)
			}
			util.LogWait("send tx to store contract...")

			// Broadcast
			res, err := a.PubClient.Broadcast(txByte)
			if err != nil {
				return util.LogErr(types.ErrContract, err)
			}

			// Save the code ID of the anchor contract in the app.yaml
			err = app.SetCodeId(a.AppFilePath, path, res)
			if err != nil {
				return util.LogErr(types.ErrContract, err)
			}

			util.LogInfo(res.Response)
			util.LogInfo(util.BB("store contract successfully"))
			return nil

		},
	}
	cmd.Flags().String(flagContractFilePath, defaultContractDirPath+defaultContractPath, "file path of the anchor contract")

	return cmd
}

// Instantiate the stored anchor contract.
func instantiateContract(a *types.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "instantiate",
		Aliases: []string{"inst"},
		Short:   "instantiate anchor contract",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s excute contract instantiate 		
$ %s e ctrt inst
$ %s e ctrt inst --code-id [code_id_of_contract]
		`, defaultAppName, defaultAppName, defaultAppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			codeId, err := cmd.Flags().GetString(flagContractCodeId)
			if err != nil {
				return util.LogErr(types.ErrContract, err)
			}

			// If the flag is not exist, use the code ID is saved in the app.yaml.
			if codeId == "" {
				codeId = app.AppFile().Get().Contract.CodeID
			}

			// Set the admin of the contract.
			admin, err := key.Bech32AddrString(a.PubClient.GetPrivateKey())
			if err != nil {
				return util.LogErr(types.ErrContract, err)
			}

			// Generate the instantiage message.
			msg := xtypes.InstantiateMsg{
				CodeId:  codeId,
				Label:   "Anchor contract",
				InitMsg: "{}",
				Amount:  "0",
				Admin:   admin,
			}
			txByte, err := a.PubClient.InstantiateContract(msg).CreateAndSignTx()
			if err != nil {
				return util.LogErr(types.ErrContract, err)
			}

			util.LogWait("send tx to instantiate contract...")
			// Broadcast.
			res, err := a.PubClient.Broadcast(txByte)
			if err != nil {
				return util.LogErr(types.ErrContract, err)
			}

			// Save the contract address of the anchor contract in the app.yaml
			err = app.SetContractAddress(a.AppFilePath, res)
			if err != nil {
				return util.LogErr(types.ErrContract, err)
			}

			util.LogInfo(res.Response)
			util.LogInfo(util.BB("instantiate contract successfully"))

			return nil
		},
	}
	cmd.Flags().String(flagContractCodeId, "", "code ID of the anchor contract")

	return cmd
}

// Query the latest block height is recorded in the anchor contract.
func queryContractLatestBlock(a *types.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest",
		Short: "query the latest block which is recorded in the anchor contract",
		Args:  withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s query contract latest
$ %s q ctrt latest
$ %s q ctrt latest --address [contract_address]
		`, defaultAppName, defaultAppName, defaultAppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := cmd.Flags().GetString(flagContractAddr)
			if err != nil {
				return util.LogErr(types.ErrContract, err)
			}

			// If the address flag is not exist, use the contract address is saved in the app.yaml.
			if addr == "" {
				addr = app.AppFile().Get().Contract.Address
			}

			// Generate the query message.
			queryMsg := xtypes.QueryMsg{
				ContractAddress: addr,
				QueryMsg:        types.QueryLatestBlockMsg,
			}

			// Request query by using XPLA client.
			res, err := a.PubClient.QueryContract(queryMsg).Query()
			if err != nil {
				return util.LogErr(types.ErrContract, err)
			}

			util.LogInfo(res)
			util.LogInfo(util.BB("response data successfully"))
			return nil

		},
	}
	cmd.Flags().String(flagContractAddr, "", "address of the anchor contract")

	return cmd
}

// Query the block info is recorded in the anchor contract.
func queryContractBlockData(a *types.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block [height]",
		Short: "query the block data which is recorded in the anchor contract",
		Args:  withUsage(cobra.ExactArgs(1)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s query contract block [height]
$ %s q ctrt latest [height]
$ %s q ctrt latest [height] --address [contract_address]
		`, defaultAppName, defaultAppName, defaultAppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := cmd.Flags().GetString(flagContractAddr)
			if err != nil {
				return util.LogErr(types.ErrContract, err)
			}

			// If the address flag is not exist, use the contract address is saved in the app.yaml.
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
			res, err := a.PubClient.QueryContract(queryMsg).Query()
			if err != nil {
				return util.LogErr(types.ErrContract, err)
			}

			util.LogInfo(res)
			util.LogInfo(util.BB("response data successfully"))
			return nil

		},
	}
	cmd.Flags().String(flagContractAddr, "", "address of the anchor contract")

	return cmd
}
