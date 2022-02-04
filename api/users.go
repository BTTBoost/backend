package api

import (
	"awake/internal/lib"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
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

	// connect db
	db, err := lib.CreateDB()
	if err != nil {
		log.Printf("failed to connect to db: %v", err)
		lib.WriteErrorResponse(w, http.StatusBadRequest, "db error")
		return
	}
	defer db.Close()

	// // read holder count from db
	// from := time.Now().Truncate(time.Hour * 24).Add(-24 * time.Hour * time.Duration(days-1)).Unix()
	// log.Printf("query args: days = %v, from = %v", days, from)
	// entries, err := db.GetDailyGroupTokenHolders(tokens, from)
	// if err != nil {
	// 	log.Printf("failed to count token holders: %v", err)
	// 	lib.WriteErrorResponse(w, http.StatusBadRequest, "db error")
	// 	return
	// }

	// // write response
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(entries)

	// write empty response
	w.Header().Set("Content-Type", "application/json")
	resp := lib.SuccessResponse{Success: true}
	json.NewEncoder(w).Encode(resp)
}
