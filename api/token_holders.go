package handler

import (
	"awake/internal/lib"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Entry struct {
	Time  int64 `json:"time"`
	Value int64 `json:"value"`
}

func TokenHoldersHandler(w http.ResponseWriter, r *http.Request) {
	// network arg
	networkQuery := r.URL.Query().Get("network")
	network, err := strconv.ParseInt(networkQuery, 10, 32)
	if err != nil || network != 1 {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid network")
		return
	}

	// token arg
	token := r.URL.Query().Get("token")
	if !lib.IsValidAddress(token) {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid token")
		return
	}

	// amount arg
	minAmountQuery := r.URL.Query().Get("min_amount")
	minAmount, ok := lib.ParseBig256(minAmountQuery)
	if !ok || minAmount.Sign() != 1 {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid min_amount")
		return
	}

	// from
	from := 1636156800

	// make a query
	query := fmt.Sprintf("SELECT time, COUNT(*) "+
		"FROM analytics.token_holders "+
		"WHERE network = %v "+
		"AND token = '%v' "+
		"AND time >= %v "+
		"AND amount >= %v "+
		"GROUP BY time "+
		"ORDER BY time;",
		network, token, from, minAmount,
	)

	// query Clickhouse over HTTP
	result, err := lib.QueryClickhouse(query)
	if err != nil {
		log.Printf("failed to query clickhouse: %v", err)
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

	// write empty response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}
