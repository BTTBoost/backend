package handler

import (
	"awake/internal/lib"
	"encoding/json"
	"net/http"
)

func NFTAllHandler(w http.ResponseWriter, r *http.Request) {
	// connect db
	db, err := lib.CreateDB()
	if err != nil {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "internal error")
		return
	}
	defer db.Close()

	nfts, err := db.GetAllNFTs()
	if err != nil {
		lib.WriteErrorResponse(w, http.StatusBadRequest, "internal error")
		return
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nfts)
}
