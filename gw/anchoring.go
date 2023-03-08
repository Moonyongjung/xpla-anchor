package gw

import (
	"github.com/Moonyongjung/xpla-anchor/app"
	"github.com/Moonyongjung/xpla-anchor/types"
	"github.com/Moonyongjung/xpla-anchor/util"
	xtypes "github.com/Moonyongjung/xpla.go/types"
)

// Send the transaction is anchoring message.
// The message has aggregated blocks info of the private chain.
func SendAnchoringTx(a *types.App, log string) {
	channel := a.Channels.AnchringTx
	for {
		select {
		case anchoringTx := <-channel:
			addr := app.AppFile().Get().Contract.Address

			bytes, err := util.JsonMarshalData(anchoringTx)
			if err != nil {
				util.LogErr(types.ErrGw, err)
				panic(err)
			}

			executeMsg := xtypes.ExecuteMsg{
				ContractAddress: addr,
				Amount:          "0",
				ExecMsg:         `{"anchoring":` + string(bytes) + `}`,
			}

			xplac := a.PubClient.WithSequence(SequenceMng().NowSequence())
			txbytes, err := xplac.ExecuteContract(executeMsg).CreateAndSignTx()
			if err != nil {
				util.LogErr(types.ErrGw, err)
				panic(err)
			}

			util.LogWait("send anchoring tx...")
			// The mode of broadcasting is "block" because of waiting until confirmed time.
			_, err = xplac.BroadcastBlock(txbytes)
			if err != nil {
				util.LogErr(types.ErrGw, err)
				panic(err)
			}

			util.LogInfo(util.BB("anchoring success"))

			SequenceMng().AddSequence()
			a.Channels.HttpClientStartSignal <- true
		}
	}
}
