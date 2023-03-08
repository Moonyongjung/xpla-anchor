package gw

import (
	"strings"
	"time"

	"github.com/Moonyongjung/xpla-anchor/app"
	"github.com/Moonyongjung/xpla-anchor/types"
	"github.com/Moonyongjung/xpla-anchor/util"
	"github.com/mitchellh/mapstructure"
)

const (
	requestBiggerHeightErr = "requested block height is bigger then the chain length"
	waitingBlockTime       = 6000
)

var dataAggregate []types.Data

// Aggregate info of blocks.
// Handle parameters which is some of the cosmos based block data such as height, hash and etc.
func aggregate(a *types.App, responseBody []byte) {
	if strings.Contains(string(responseBody), requestBiggerHeightErr) {
		util.LogWait("wating for creating new block...")
		BlockListMng().DecreaseLatestBlockHeight()
		time.Sleep(time.Millisecond * time.Duration(waitingBlockTime))

		a.Channels.HttpClientStartSignal <- true

	} else {
		count := app.AppFile().Get().Config.Anchor.CollectBlockCount

		var block types.Block
		responseData := util.JsonUnmarshalData(&block, responseBody)
		mapstructure.Decode(responseData, &block)

		util.LogInfo(util.BB("height=")+block.Block.Header.Height, util.BB("hash=")+block.BlockID.Hash)

		if block.BlockID.Hash == "" || block.Block.Header.Height == "" {
			util.LogWarning("empty response, check the LCD URL or block info API")
			BlockListMng().DecreaseLatestBlockHeight()
			time.Sleep(time.Millisecond * time.Duration(waitingBlockTime))

		} else {
			height := block.Block.Header.Height
			block_hash := block.BlockID.Hash
			data_merkle := block.Block.Header.DataHash
			timestamp := block.Block.Header.Time

			newData := types.NewData(height, block_hash, data_merkle, timestamp)

			// Listing aggreated info.
			dataAggregate = append(dataAggregate, newData)
		}

		if len(dataAggregate) == count {
			latest := dataAggregate[count-1].Height

			util.LogInfo(util.BB("fin aggregate"))
			util.LogInfo(util.BB("aggregated first block height=") + dataAggregate[0].Height)
			util.LogInfo(util.BB("aggregated latest block height=") + latest)

			newAnchoring := types.NewAncoring(dataAggregate, latest)

			a.Channels.AnchringTx <- newAnchoring

			dataAggregate = nil
		} else {
			a.Channels.HttpClientStartSignal <- true
		}
	}

}
