package main

import (
	"awake/internal/lib"
	"log"
	"math/big"

	"github.com/joho/godotenv"
)

// Fetches and saves to db token holders at each day from Covalent API.
func main() {
	// load .env into ENV
	err := godotenv.Overload()
	if err != nil {
		log.Fatal(err)
	}

	// t0 := lib.TokenHoldersQuery(1, "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", 1636156800, big.NewInt(1))
	// t1 := lib.TokenHoldersQuery(1, "0x026224a2940bfe258d0dbe947919b62fe321f042", 1636156800, big.NewInt(1))
	// t2 := lib.TokenHoldersQuery(1, "0x5a98fcbea516cf06857215779fd812ca3bef1b32", 1636156800, big.NewInt(1))
	// t0xt1 := lib.JoinHolderQueries(t0, t1)
	// t0xt1xt2 := lib.JoinHolderQueries(t0xt1, t2)
	// groupt0 := lib.GroupHolderQuery(t0)
	// groupt0xt1xt2 := lib.GroupHolderQuery(t0xt1xt2)

	// log.Printf("T0: %v", t0)
	// log.Printf("T1: %v", t1)
	// log.Printf("T2: %v", t2)
	// log.Printf("T0xT1: %v", t0xt1)
	// log.Printf("(T0xT1)xT2: %v", t0xt1xt2)
	// log.Printf("GROUP(T0): %v", groupt0)
	// log.Printf("GROUP((T0xT1)xT2): %v", groupt0xt1xt2)

	// // params
	// var from int64 = 1636156800
	// tokens := []string{"0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", "0x026224a2940bfe258d0dbe947919b62fe321f042"}
	// amounts := []*big.Int{big.NewInt(1), big.NewInt(1)}

	// // generate sql query
	// query := lib.TokenHoldersQuery(1, tokens[0], from, amounts[0])
	// for i := 1; i < len(tokens); i++ {
	// 	qi := lib.TokenHoldersQuery(1, tokens[i], from, amounts[i])
	// 	query = lib.JoinHolderQueries(query, qi)
	// }
	// query = lib.GroupHolderQuery(query)
	// log.Printf("query: %v", query)

	var from int64 = 1636156800
	var days int64 = 90
	holdersQuery := lib.TokenHoldersQuery(1, "0x5a98fcbea516cf06857215779fd812ca3bef1b32", from, big.NewInt(1))
	eventsQuery := lib.EventQuery("aave_deposits", from, days)
	holderDailyEventsQuery := lib.HolderDailyEventsQuery(holdersQuery, eventsQuery, from, days)

	log.Printf("holders query: %v", holdersQuery)
	log.Printf("events query: %v", eventsQuery)
	log.Printf("holder daily events query: %v", holderDailyEventsQuery)
}
