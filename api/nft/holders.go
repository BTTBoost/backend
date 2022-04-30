package handler

import (
	"awake/internal/lib"
	"encoding/json"
	"net/http"
	"sort"
)

func NFTHoldersHandler(w http.ResponseWriter, r *http.Request) {
	// token arg
	token := r.URL.Query().Get("token")
	if !lib.IsValidAddress(token) {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid token")
		return
	}

	// connect db
	db, err := lib.CreateDB()
	if err != nil {
		lib.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer db.Close()

	// get holders from db
	holders, err := db.GetLastTokenHolders(1, token)
	if err != nil {
		lib.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// sort
	sort.SliceStable(holders, func(i, j int) bool {
		return holders[i].Amount > holders[j].Amount
	})

	// write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(holders)
}
