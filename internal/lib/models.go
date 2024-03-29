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

type ProtocolUsage struct {
	Name           string `json:"name"`
	Logo           string `json:"logo,omitempty"`
	Url            string `json:"url,omitempty"`
	UsersLastMonth int64  `json:"users_last_month"`
	UsersInTotal   int64  `json:"users_in_total"`
}

type NetworkUsage struct {
	Name           string `json:"name"`
	Logo           string `json:"logo,omitempty"`
	Url            string `json:"url,omitempty"`
	UsersLastMonth int64  `json:"users_last_month"`
	UsersInTotal   int64  `json:"users_in_total"`
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
	FirstTransfer   int64   `json:"first_transfer"`
	TotalBalanceUsd float64 `json:"total_balance_usd"`
}

type NFTStats struct {
	Token     Token `json:"token"`
	Active1d  int64 `json:"active_1d"`
	Active7d  int64 `json:"active_7d"`
	Active30d int64 `json:"active_30d"`
	Total     int64 `json:"total"`
}
