package handler

import (
	"awake/internal/lib"
	"bufio"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func HoldingsHandler(w http.ResponseWriter, r *http.Request) {
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

	var time int64 = 1643932800

	// generate sql query
	query := lib.TokenHoldersQuery(int(networks[0]), tokens[0], time, amounts[0])
	for i := 1; i < len(tokens); i++ {
		qi := lib.TokenHoldersQuery(int(networks[i]), tokens[i], time, amounts[i])
		query = lib.JoinHolderQueries(query, qi)
	}
	query = lib.TopHoldingTokensQuery(query, time)

	// query Clickhouse over HTTP
	result, err := lib.QueryClickhouse(query)
	if err != nil {
		// log.Printf("failed to query clickhouse: %v", err)
		lib.WriteErrorResponse(w, http.StatusBadRequest, "internal error")
		return
	}

	// parse response
	rows := []lib.HoldingRow{}
	scanner := bufio.NewScanner(strings.NewReader(result))
	for scanner.Scan() {
		ss := strings.Split(scanner.Text(), "	")
		network, err := strconv.ParseInt(ss[0], 10, 64)
		if err != nil {
			lib.WriteErrorResponse(w, http.StatusBadRequest, "internal error")
			return
		}
		address := ss[1]
		token := lib.Token{Network: network, Address: address}
		holders, err := strconv.ParseInt(ss[2], 10, 64)
		if err != nil {
			lib.WriteErrorResponse(w, http.StatusBadRequest, "internal error")
			return
		}
		share, err := strconv.ParseFloat(ss[3], 64)
		if err != nil {
			lib.WriteErrorResponse(w, http.StatusBadRequest, "internal error")
			return
		}
		rows = append(rows, lib.HoldingRow{Token: token, Holders: holders, Share: share})
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}
