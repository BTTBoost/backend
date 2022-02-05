package lib

import (
	"fmt"
	"time"
)

type AaveDeposit struct {
	block      int64
	time       int64
	tx         string
	user       string
	onBehalfOf string
	reserve    string
	amount     string
}

func ParseAaveDepositEvent(event CovalentEvent) (*AaveDeposit, error) {
	time, err := time.Parse(EVENT_TIME_LAYOUT, event.BlockSignedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse block time '%v': %v", event.BlockSignedAt, err)
	}

	var user string
	var onBehalfOf string
	var reserve string
	var amount string
	for _, p := range event.Decoded.Params {
		switch p.Name {
		case "user":
			user = p.Value
		case "onBehalfOf":
			onBehalfOf = p.Value
		case "reserve":
			reserve = p.Value
		case "amount":
			amount = p.Value
		}
	}

	return &AaveDeposit{
		event.BlockHeight,
		time.Unix(),
		event.TxHash,
		user,
		onBehalfOf,
		reserve,
		amount,
	}, nil
}
