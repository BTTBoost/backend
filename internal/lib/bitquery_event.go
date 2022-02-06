package lib

// event models
type BitqueryEvent struct {
	BlockData BitqueryEventBlock      `json:"block"`
	Tx        BitqueryEventTx         `json:"transaction"`
	Arguments []BitqueryEventArgument `json:"arguments"`
}

type BitqueryEventBlock struct {
	Block    int64                  `json:"height"`
	TimeData BitqueryEventBlockTime `json:"timestamp"`
}

type BitqueryEventBlockTime struct {
	Timestamp int64 `json:"unixtime"`
}

type BitqueryEventTx struct {
	Hash string `json:"hash"`
}

type BitqueryEventArgument struct {
	Name  string `json:"argument"`
	Value string `json:"value"`
}

// response models
type BitqueryResponse struct {
	Data   BitqueryResponseData    `json:"data"`
	Errors []BitqueryResponseError `json:"errors"`
}

type BitqueryResponseData struct {
	BitqueryResponseDataEthereum `json:"ethereum"`
}

type BitqueryResponseDataEthereum struct {
	Events []BitqueryEvent `json:"smartContractEvents"`
}

type BitqueryResponseError struct {
	Message string `json:"message"`
}
