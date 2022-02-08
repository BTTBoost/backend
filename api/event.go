package handler

import (
	"awake/internal/lib"
	"bufio"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func EventHandler(w http.ResponseWriter, r *http.Request) {
	// fixed time range
	var days int64 = 90
	var from int64 = 1636156800

	// networks arg
	networksQuery := strings.Split(r.URL.Query().Get("networks"), ",")
	networks, err := lib.ParseInt64Slice(networksQuery)
	if len(networks) == 0 || err != nil {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid networks param")
		return
	}

	// tokens arg
	tokensParam := strings.ToLower(r.URL.Query().Get("tokens"))
	tokens := strings.Split(tokensParam, ",")
	if !lib.IsValidAddressSlice(tokens) || len(tokens) != len(networks) {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid tokens param")
		return
	}

	// amounts arg
	amountsQuery := strings.Split(r.URL.Query().Get("amounts"), ",")
	amounts, ok := lib.ParseTokenAmountSlice(amountsQuery)
	if !ok || len(tokens) != len(amounts) {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid amounts param")
		return
	}

	// event arg
	event := r.URL.Query().Get("event")

	// generate sql query
	holdersQuery := lib.TokenHoldersQuery(int(networks[0]), tokens[0], from, days, amounts[0])
	for i := 1; i < len(tokens); i++ {
		qi := lib.TokenHoldersQuery(int(networks[i]), tokens[i], from, days, amounts[i])
		holdersQuery = lib.JoinHolderQueries(holdersQuery, qi)
	}

	// event query
	eventQuery := ""
	switch event {
	case "aave_deposit_ethereum":
		eventQuery = lib.EventQuery("aave_deposits", 1, from, days)
	case "aave_deposit_polygon":
		eventQuery = lib.EventQuery("aave_deposits", 137, from, days)
	case "looksrare_trade_ethereum":
		// TODO: currently ignores takers (!)
		eventQuery = lib.LooksRareTradesQuery(1, from, days)
	default:
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid event")
		return
	}

	query := lib.HolderDailyEventsQuery(holdersQuery, eventQuery, from, days)

	// query Clickhouse over HTTP
	result, err := lib.QueryClickhouse(query)
	if err != nil {
		// log.Printf("failed to query clickhouse: %v", err)
		lib.WriteErrorResponse(w, http.StatusBadRequest, "internal error")
		return
	}

	// parse response
	entries := []lib.ChartEntry{}
	scanner := bufio.NewScanner(strings.NewReader(result))
	for scanner.Scan() {
		ss := strings.Split(scanner.Text(), "	")
		time, err := strconv.ParseInt(ss[0], 10, 64)
		if err != nil {
			lib.WriteErrorResponse(w, http.StatusBadRequest, "internal error")
			return
		}
		value, err := strconv.ParseInt(ss[1], 10, 64)
		if err != nil {
			lib.WriteErrorResponse(w, http.StatusBadRequest, "internal error")
			return
		}
		entries = append(entries, lib.ChartEntry{Time: time, Value: value})
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}
