package lib

import (
	"fmt"
	"time"
)

type OpenSeaTrade struct {
	block int64
	time  int64
	tx    string
	maker string
	taker string
	price string
}

func OpenseaEventToTrade(event CovalentEvent) (*OpenSeaTrade, error) {
	time, err := time.Parse(EVENT_TIME_LAYOUT, event.BlockSignedAt)
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
