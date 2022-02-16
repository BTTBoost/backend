package lib

type ChartEntry struct {
	Time  int64 `json:"time"`
	Value int64 `json:"value"`
}

type HoldingRow struct {
	Token   Token   `json:"token"`
	Holders int64   `json:"holders"`
	Share   float64 `json:"share"`
}

type Token struct {
	Network int64  `json:"network"`
	Address string `json:"address"`
}

type CovalentHolder struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}
