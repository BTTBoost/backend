package lib

type OpenSeaTrade struct {
	block    int64
	time     int64
	tx       string
	buyHash  string
	sellHash string
	maker    string
	taker    string
	price    string
}

func ParseOpenseaTradeBitqueryEvent(event BitqueryEvent) (*OpenSeaTrade, error) {
	var buyHash string
	var sellHash string
	var maker string
	var taker string
	var price string
	for _, a := range event.Arguments {
		switch a.Name {
		case "buyHash":
			buyHash = a.Value
		case "sellHash":
			sellHash = a.Value
		case "maker":
			maker = a.Value
		case "taker":
			taker = a.Value
		case "price":
			price = a.Value
		}
	}

	return &OpenSeaTrade{
		event.BlockData.Block,
		event.BlockData.TimeData.Timestamp,
		event.Tx.Hash,
		buyHash,
		sellHash,
		maker,
		taker,
		price,
	}, nil
}
