package types

const (
	// Latest block message.
	QueryLatestBlockMsg = `{"latest_block":{}}`
)

// The type of the sending transaction for anchring.
type Anchoring struct {
	Data   []Data `json:"data"`
	Latest string `json:"latest"`
}

func NewAncoring(data []Data, latest string) Anchoring {
	var anchoring Anchoring

	anchoring.Data = data
	anchoring.Latest = latest

	return anchoring
}

type Data struct {
	Height     string `json:"height"`
	BlockHash  string `json:"block_hash"`
	DataMerkle string `json:"data_merkle"`
	Timestamp  string `json:"timestamp"`
}

func NewData(height, blockHash, dataMerkle, timestamp string) Data {
	var data Data

	data.Height = height
	data.BlockHash = blockHash
	data.DataMerkle = dataMerkle
	data.Timestamp = timestamp

	return data
}

type QueryLatestBlockResponse struct {
	Data struct {
		LatestHeight string `json:"latest_height"`
	} `json:"data"`
}

type QueryBlockInfoResponse struct {
	Data struct {
		Height     string `json:"height"`
		BlockHash  string `json:"block_hash"`
		DataMerkle string `json:"data_merkle"`
		Timestamp  string `json:"timestamp"`
	} `json:"data"`
}
