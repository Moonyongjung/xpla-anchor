package gw

import (
	"sync"

	"github.com/Moonyongjung/xpla-anchor/util"
)

var SequenceInstance *SequenceStruct
var SequenceOnce sync.Once

// Manage account sequence
type SequenceStruct struct {
	Sequence string
}

func SequenceMng() *SequenceStruct {
	SequenceOnce.Do(func() {
		SequenceInstance = &SequenceStruct{}
	})
	return SequenceInstance
}

func (n *SequenceStruct) NewSequence(sequence string) {
	n.Sequence = sequence
}

func (n *SequenceStruct) NowSequence() string {
	return n.Sequence
}

func (n *SequenceStruct) AddSequence() {
	Sequence := n.Sequence
	SequenceNum := util.FromStringToUint64(Sequence)
	SequenceNum = SequenceNum + 1
	Sequence = util.FromUint64ToString(SequenceNum)
	n.Sequence = Sequence
}
