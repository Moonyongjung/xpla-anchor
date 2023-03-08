package gw

import (
	"sync"

	"github.com/Moonyongjung/xpla-anchor/types"
	"github.com/Moonyongjung/xpla-anchor/util"
)

var blockListInstance *BlockList
var blockListOnce sync.Once

// Recording block structure includes end block number.
// The end block indicates latest block that retrieved by anchor.
type BlockList struct {
	EndBlock string
}

func BlockListMng() *BlockList {
	blockListOnce.Do(func() {
		blockListInstance = &BlockList{}
	})
	return blockListInstance
}

func (b *BlockList) NewLatestBlockHeight(EndBlock string) {
	b.EndBlock = EndBlock
}

func (b *BlockList) NowLatestBlockHeight() string {
	return b.EndBlock
}

func (b *BlockList) IncreaseLatestBlockHeight() {
	temp, err := util.ToInt(b.EndBlock)
	if err != nil {
		util.LogErr(types.ErrBlockMng, err)
		panic(err)
	}
	increasedBlockHeight := util.ToString(temp+1, "")
	b.EndBlock = increasedBlockHeight
}

func (b *BlockList) DecreaseLatestBlockHeight() {
	temp, err := util.ToInt(b.EndBlock)
	if err != nil {
		util.LogErr(types.ErrBlockMng, err)
		panic(err)
	}
	increasedBlockHeight := util.ToString(temp-1, "")
	b.EndBlock = increasedBlockHeight
}
