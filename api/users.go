package handler

import (
	"awake/internal/lib"
	"bufio"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	// tokens arg
	tokens := strings.Split(r.URL.Query().Get("tokens"), ",")
	if len(tokens) == 0 || !lib.IsValidAddressSlice(tokens) {
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

	// days arg
	daysQuery := r.URL.Query().Get("days")
	days, err := strconv.ParseInt(daysQuery, 10, 32)
	if err != nil || days <= 0 || days > 365 {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid days param")
		return
	}

	var from int64 = 1636156800

	// generate sql query
	query := lib.TokenHoldersQuery(1, tokens[0], from, amounts[0])
	for i := 1; i < len(tokens); i++ {
		qi := lib.TokenHoldersQuery(1, tokens[i], from, amounts[i])
		query = lib.JoinHolderQueries(query, qi)
	}
	query = lib.GroupHolderQuery(query)

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
		entries = append(entries, Entry{Time: ss[0], Amount: ss[1]})
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)

	lib.WriteErrorResponse(w, http.StatusBadRequest, "not implemented")
}
