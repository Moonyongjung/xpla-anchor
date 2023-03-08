package types

type Channels struct {
	HttpClientStartSignal chan bool
	AnchringTx            chan Anchoring
}
