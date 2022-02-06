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
	// networks arg
	networksQuery := strings.Split(r.URL.Query().Get("networks"), ",")
	networks, err := lib.ParseInt64Slice(networksQuery)
	if len(networks) == 0 || err != nil {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid networks param")
		return
	}

	// tokens arg
	tokens := strings.Split(r.URL.Query().Get("tokens"), ",")
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
	eventTable := ""
	var eventNetwork int64
	switch event {
	case "aave_deposit_ethereum":
		eventTable = "aave_deposits"
		eventNetwork = 1
	case "aave_deposit_polygon":
		eventTable = "aave_deposits"
		eventNetwork = 137
	default:
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid event")
		return
	}

	// fixed time range
	var days int64 = 90
	var from int64 = 1636156800

	// generate sql query
	holdersQuery := lib.TokenHoldersQuery(int(networks[0]), tokens[0], from, days, amounts[0])
	for i := 1; i < len(tokens); i++ {
		qi := lib.TokenHoldersQuery(int(networks[i]), tokens[i], from, days, amounts[i])
		holdersQuery = lib.JoinHolderQueries(holdersQuery, qi)
	}
	eventQuery := lib.EventQuery(eventTable, eventNetwork, from, days)
	query := lib.HolderDailyEventsQuery(holdersQuery, eventQuery, from, days)

	// query Clickhouse over HTTP
	result, err := lib.QueryClickhouse(query)
	if err != nil {
		// log.Printf("failed to query clickhouse: %v", err)
		lib.WriteErrorResponse(w, http.StatusBadRequest, "internal error")
		return
	}

	// parse response
	entries := []Entry{}
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
		entries = append(entries, Entry{Time: time, Value: value})
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}
