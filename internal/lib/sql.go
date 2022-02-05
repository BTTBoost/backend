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
