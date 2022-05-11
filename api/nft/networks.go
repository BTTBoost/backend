package handler

import (
	"awake/internal/lib"
	"encoding/json"
	"net/http"
	"strconv"
)

func NFTNetworksHandler(w http.ResponseWriter, r *http.Request) {
	// token arg
	token := r.URL.Query().Get("token")
	if !lib.IsValidAddress(token) {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid token")
		return
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
	list, err := db.GetNFTNetworks(1, token, limit)
	if err != nil {
		lib.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
