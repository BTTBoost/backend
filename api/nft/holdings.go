package handler

import (
	"awake/internal/lib"
	"encoding/json"
	"net/http"
	"strconv"
)

func NFTHoldingsHandler(w http.ResponseWriter, r *http.Request) {
	// token arg
	token := r.URL.Query().Get("token")
	if !lib.IsValidAddress(token) {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid token")
		return
	}

	// type arg
	tokenType := r.URL.Query().Get("type")
	if tokenType == "" {
		tokenType = "erc20"
	}
	if tokenType != "erc20" && tokenType != "nft" {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid token type: must be 'erc20' or 'nft'")
		return
	}
	nft := false
	if tokenType == "nft" {
		nft = true
	}

	// limit arg
	limit := 20
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid limit")
	}

	// connect db
	db, err := lib.CreateDB()
	if err != nil {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "internal error")
		return
	}
	defer db.Close()

	// query db
	holdings, err := db.GetNFTTokenHoldings(1, token, nft, limit)
	if err != nil {
		lib.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(holdings)
}
