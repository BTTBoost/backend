package lib

const EVENT_TIME_LAYOUT = "2006-01-02T15:04:05Z"

type CovalentEvent struct {
	BlockSignedAt string               `json:"block_signed_at"`
	BlockHeight   int64                `json:"block_height"`
	TxHash        string               `json:"tx_hash"`
	Decoded       CovalentEventDecoded `json:"decoded"`
}

type CovalentEventDecoded struct {
	Params []CovalentEventParam `json:"params"`
}

type CovalentEventParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
