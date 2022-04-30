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
	Network  int64  `json:"network"`
	Address  string `json:"address"`
	Name     string `json:"name,omitempty"`
	Symbol   string `json:"symbol,omitempty"`
	Logo     string `json:"logo,omitempty"`
	Decimals int64  `json:"decimals,omitempty"`
}

type CovalentHolder struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

type NFTCollection struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Logo    string `json:"logo"`
}

type NFTHolder struct {
	Address         string  `json:"address"`
	Amount          int64   `json:"amount"`
	TotalBalanceUsd float64 `json:"total_balance_usd"`
}
