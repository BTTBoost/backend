package lib

import (
	"fmt"
	"math/big"
)

func TokenHoldersDaysQuery(network int, token string, from int64, days int64, amount *big.Int) string {
	to := from + days*86400
	return fmt.Sprintf("SELECT time, holder "+
		"FROM analytics.token_holders "+
		"WHERE network = %v AND token = '%v' AND time >= %v AND time < %v AND amount >= %v",
		network, token, from, to, amount,
	)
}

func JoinHolderDaysQueries(query0 string, query1 string) string {
	return fmt.Sprintf("SELECT T1.time, T1.holder "+
		"FROM (%v) T1 "+
		"INNER JOIN (%v) T2 "+
		"ON (T1.time = T2.time AND T1.holder = T2.holder)",
		query0, query1,
	)
}

func GroupHolderDaysQuery(query string, from int64, days int64) string {
	to := from + days*86400
	return fmt.Sprintf(
		"SELECT time, COUNT(*) "+
			"FROM (%v) "+
			"GROUP BY time "+
			"ORDER BY time ASC WITH FILL FROM %v TO %v STEP 86400",
		query, from, to,
	)
}

func EventQuery(event string, network int64, from int64, days int64) string {
	to := from + days*86400
	return fmt.Sprintf(
		"SELECT time_day as time, user, amount "+
			"FROM analytics.%v "+
			"WHERE network = %v AND time >= %v AND time < %v "+
			"ORDER BY toUnixTimestamp(toStartOfDay(toDate(time))) as time_day",
		event, network, from, to,
	)
}

func LooksRareTradesQuery(network int64, from int64, days int64) string {
	to := from + days*86400
	return fmt.Sprintf(
		"SELECT time_day as time, user "+
			"FROM ("+
			"SELECT time, maker as user "+
			"FROM analytics.looksrare_trades "+
			"WHERE network = %[1]v AND time >= %[2]v AND time < %[3]v "+
			"UNION ALL "+
			"SELECT time, taker as user "+
			"FROM analytics.looksrare_trades "+
			"WHERE network = %[1]v AND time >= %[2]v AND time < %[3]v "+
			") "+
			"ORDER BY toUnixTimestamp(toStartOfDay(toDate(time))) as time_day",
		network, from, to,
	)
}

func OpenseaTradesEthereumQuery(from int64, days int64) string {
	to := from + days*86400
	return fmt.Sprintf(
		"SELECT time_day as time, user "+
			"FROM ("+
			"SELECT time, taker as user "+
			"FROM analytics.opensea_trades "+
			"WHERE network = 1 AND time >= %[1]v AND time < %[2]v "+
			"UNION ALL "+
			"SELECT time, maker as user "+
			"FROM analytics.opensea_trades "+
			"WHERE network = 1 AND time >= %[1]v AND time < %[2]v "+
			")"+
			"ORDER BY toUnixTimestamp(toStartOfDay(toDate(time))) as time_day",
		from, to,
	)
}

func HolderDailyEventsQuery(holdersQuery string, eventQuery string, from int64, days int64) string {
	to := from + days*86400
	return fmt.Sprintf(
		"SELECT H.time, COUNT(*) "+
			"FROM (%v) H "+
			"INNER JOIN (%v) E "+
			"ON H.time = E.time AND H.holder = E.user "+
			"GROUP BY H.time "+
			"ORDER BY H.time WITH FILL FROM %v TO %v STEP 86400",
		holdersQuery, eventQuery, from, to,
	)
}
