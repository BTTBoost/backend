package handler

import (
	"awake/internal/lib"
	"encoding/json"
	"net/http"
)

func NFTStatsHandler(w http.ResponseWriter, r *http.Request) {
	// token arg
	token := r.URL.Query().Get("token")
	if !lib.IsValidAddress(token) {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "invalid token")
		return
	}

	// connect db
	db, err := lib.CreateDB()
	if err != nil {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "internal error")
		return
	}
	defer db.Close()

	stats, err := db.GetNFTStats(token)
	if err != nil {
		lib.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
