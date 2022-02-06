package lib

type LooksRareTrade struct {
	block    int64
	time     int64
	tx       string
	maker    string
	taker    string
	currency string
	price    string
}

func ParseLooksRareTradeBitqueryEvent(event BitqueryEvent) (*LooksRareTrade, error) {
	var maker string
	var taker string
	var currency string
	var price string
	for _, a := range event.Arguments {
		switch a.Name {
		case "maker":
			maker = a.Value
		case "taker":
			taker = a.Value
		case "currency":
			currency = a.Value
		case "price":
			price = a.Value
		}
	}

	return &LooksRareTrade{
		event.BlockData.Block,
		event.BlockData.TimeData.Timestamp,
		event.Tx.Hash,
		maker,
		taker,
		currency,
		price,
	}, nil
}
