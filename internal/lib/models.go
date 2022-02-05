package lib

import (
	"fmt"
	"time"
)

const OPENSEA_TRADE_TIME_LAYOUT = "2006-01-02T15:04:05Z"

type CovalentHolder struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

// OpenseaTradeEvent is a trade event of OpenSea from Covalent API
type OpenseaTradeEvent struct {
	BlockSignedAt string                   `json:"block_signed_at"`
	BlockHeight   int64                    `json:"block_height"`
	TxHash        string                   `json:"tx_hash"`
	Decoded       OpenseaTradeEventDecoded `json:"decoded"`
}

type OpenseaTradeEventDecoded struct {
	Params []OpenseaTradeEventDecodedParam `json:"params"`
}

type OpenseaTradeEventDecodedParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type OpenSeaTrade struct {
	block int64
	time  int64
	tx    string
	maker string
	taker string
	price string
}

func OpenseaEventToTrade(event OpenseaTradeEvent) (*OpenSeaTrade, error) {
	time, err := time.Parse(OPENSEA_TRADE_TIME_LAYOUT, event.BlockSignedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse block time '%v': %v", event.BlockSignedAt, err)
	}

	var maker string
	var taker string
	var price string
	for _, p := range event.Decoded.Params {
		switch p.Name {
		case "maker":
			maker = p.Value
		case "taker":
			taker = p.Value
		case "price":
			price = p.Value
		}
	}

	return &OpenSeaTrade{
		event.BlockHeight,
		time.Unix(),
		event.TxHash,
		maker,
		taker,
		price,
	}, nil
}
