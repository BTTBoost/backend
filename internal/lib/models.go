package lib

type ChartEntry struct {
	Time  int64 `json:"time"`
	Value int64 `json:"value"`
}

type CovalentHolder struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}
