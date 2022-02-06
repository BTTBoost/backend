package lib

import (
	"fmt"
	"math/big"
)

func TokenHoldersQuery(network int, token string, from int64, amount *big.Int) string {
	return fmt.Sprintf("SELECT time, holder "+
		"FROM analytics.token_holders "+
		"WHERE network = %v AND token = '%v' AND time >= %v AND amount >= %v",
		network, token, from, amount,
	)
}

func JoinHolderQueries(query0 string, query1 string) string {
	return fmt.Sprintf("SELECT T1.time, T1.holder "+
		"FROM (%v) T1 "+
		"INNER JOIN (%v) T2 "+
		"ON (T1.time = T2.time AND T1.holder = T2.holder)",
		query0, query1,
	)
}

func GroupHolderQuery(query string) string {
	return fmt.Sprintf(
		"SELECT time, COUNT(*) "+
			"FROM (%v) "+
			"GROUP BY time "+
			"ORDER BY time ASC",
		query,
	)
}

func EventQuery(event string, from int64, days int64) string {
	to := 1636156800 + days*86400
	return fmt.Sprintf(
		"SELECT time_day as time, user, amount "+
			"FROM analytics.%v "+
			"WHERE time >= %v AND time < %v "+
			"ORDER BY toUnixTimestamp(toStartOfDay(toDate(time))) as time_day",
		event, from, to,
	)
}

func HolderDailyEventsQuery(holdersQuery string, eventQuery string, from int64, days int64) string {
	to := 1636156800 + days*86400
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
