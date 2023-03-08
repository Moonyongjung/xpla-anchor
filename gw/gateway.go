package gw

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"github.com/Moonyongjung/xpla-anchor/app"
	"github.com/Moonyongjung/xpla-anchor/types"
	"github.com/Moonyongjung/xpla-anchor/util"
	"github.com/Moonyongjung/xpla.go/client"
	"github.com/Moonyongjung/xpla.go/key"
	xtypes "github.com/Moonyongjung/xpla.go/types"
	"github.com/mitchellh/mapstructure"
)

// Start the gateway of the anchor.
// Request the block info to the private chain periodically,
// and send the transaction as anchoring message to the main chain.
func StartGW(a *types.App, contractAddr, blockApi, log string) {
	util.LogInfo(util.BB("target anchor contract=") + contractAddr)

	// Set requested period in the config.yaml (milliseconds)
	requestPeriod := app.AppFile().Get().Config.Anchor.RequestPeriod
	if requestPeriod < 0 {
		util.LogErr(types.ErrGw, "request period must be not negative")
	}

	go request(a, blockApi, requestPeriod)

	initGW(a, contractAddr)
}

// Set the periodic time.
func request(a *types.App, blockApi string, requestPeriod int) {
	for {
		if <-a.Channels.HttpClientStartSignal {
			_, err := DoRequest(a, blockApi, "", true)
			if err != nil {
				util.LogErr(types.ErrGw, err)
				panic(err)
			}
			BlockListMng().IncreaseLatestBlockHeight()
			time.Sleep(time.Millisecond * time.Duration(requestPeriod))
		}
	}
}

// Initailize the gateway.
// At first, query the recorded latest block height to the anchor contract
// in order to request next block to the private chain.
// If the that block height is zero, the gateway request the genesis block info of the private chain.
func initGW(a *types.App, contractAddr string) {
	queryMsg := xtypes.QueryMsg{
		ContractAddress: contractAddr,
		QueryMsg:        types.QueryLatestBlockMsg,
	}

	// Check the recorded latest block height in the contract.
	res, err := a.PubClient.QueryContract(queryMsg).Query()
	if err != nil {
		util.LogErr(types.ErrGw, err)
		panic(err)
	}

	var latestBlock types.QueryLatestBlockResponse
	responseData := util.JsonUnmarshalData(&latestBlock, []byte(res))
	mapstructure.Decode(responseData, &latestBlock)

	util.LogInfo(util.BB("recorded latest block height=") + latestBlock.Data.LatestHeight)

	if latestBlock.Data.LatestHeight == "0" {
		BlockListMng().NewLatestBlockHeight(types.GenesisBlockNum)
	} else {
		BlockListMng().NewLatestBlockHeight(latestBlock.Data.LatestHeight)
		BlockListMng().IncreaseLatestBlockHeight()
	}

	user, err := key.Bech32AddrString(a.PubClient.GetPrivateKey())
	if err != nil {
		util.LogErr(types.ErrGw, err)
		panic(err)
	}

	// Query the sequence number of the account in order to run the gateway.
	seqRes, err := querySequence(a.PubClient, user)
	if err != nil {
		util.LogErr(types.ErrGw, err)
		panic(err)
	}

	seq := util.ParsingQueryAccount(seqRes)
	SequenceMng().NewSequence(seq)

	a.Channels.HttpClientStartSignal <- true
}

// Request block info.
func DoRequest(a *types.App, blockApi string, blockHeight string, isGateway bool) ([]byte, error) {
	if isGateway {
		blockHeight = BlockListMng().NowLatestBlockHeight()
	}

	privLcdUrl := app.AppFile().Get().Config.PrivateChain.LCD + blockApi + "/" + blockHeight
	util.LogInfo(util.BB("URL=") + privLcdUrl)

	request, err := http.NewRequest("GET", privLcdUrl, nil)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	httpClient.Timeout = time.Second * 30

	defer func() {
		if err := recover(); err != nil {
			util.LogInfo(err)
		}
	}()

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if isGateway {
		go aggregate(a, responseBody)
	}

	return responseBody, nil
}

// check the sequence number of the anchor account.
func querySequence(xplac *client.XplaClient, addr string) (string, error) {
	queryAccAddressMsg := xtypes.QueryAccAddressMsg{
		Address: addr,
	}
	res, err := xplac.AccAddress(queryAccAddressMsg).Query()
	if err != nil {
		return "", err
	}

	return res, nil
}
