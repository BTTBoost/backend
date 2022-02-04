package main

import (
	"awake/internal/lib"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Response struct {
	Data         ResponseData `json:"data"`
	Error        bool         `json:"error"`
	ErrorMessage string       `json:"error_message"`
}

type ResponseData struct {
	Items []lib.CovalentHolder `json:"items"`
	// Pagination ResponseDataPagination `json:"pagination"` returns null values
}

// TODO: support pagination
func fetchHolders(network int, token string, block int64) ([]lib.CovalentHolder, error) {
	apiKey := os.Getenv("COVALENT_API_KEY")
	pageSize := 100000

	url := fmt.Sprintf(
		"https://api.covalenthq.com/v1/%v/tokens/%v/token_holders/"+
			"?block-height=%v&page-number=0&page-size=%v&key=%v",
		network, token, block, pageSize, apiKey,
	)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	result := &Response{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(result.Data.Items) == pageSize {
		return nil, fmt.Errorf("too many holders, pagination required")
	}

	return result.Data.Items, nil
}

// Fetches and saves to db token holders at each day from Covalent API.
func main() {
	// load .env into ENV
	err := godotenv.Overload()
	if err != nil {
		log.Fatal(err)
	}

	// init params
	if len(os.Args) != 4 {
		log.Fatal("wrong arguments. pass [network] [token] [days]")
	}

	var network int
	networkParam := os.Args[1]
	switch networkParam {
	case "ethereum":
		network = 1
	case "polygon":
		network = 137
	default:
		log.Fatalf("unsupported network '%v'. supported values [ethereum, polygon]", networkParam)
	}

	token := os.Args[2]
	if !lib.IsValidAddress(token) {
		log.Fatalf("invalid token param '%v'", token)
	}

	daysParam := os.Args[3]
	days, err := strconv.ParseInt(daysParam, 10, 64)
	if err != nil || days < 1 || days > 10000 {
		log.Fatalf("invalid days param '%v'. supported range is [1, 10000]", err)
	}

	var to = time.Now().Truncate(time.Hour * 24)
	var from = to.Add(-time.Hour * 24 * time.Duration(days))
	log.Printf("fetching holders from [%v] to [%v]...", from.Format("02-01-2006"), to.Format("02-01-2006"))

	// connect db
	db, err := lib.CreateDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// read all time-block pairs from db
	timeBlocks, err := db.GetDailyBlocks(network)
	if err != nil {
		log.Fatalf("failed to read blocks from db: %v", err)
	}

	// fetch and save holders
	for t := from; !t.After(to); t = t.Add(time.Hour * 24) {
		block := timeBlocks[t.Unix()]
		holders, err := fetchHolders(network, token, block)
		if err != nil {
			log.Fatalf("failed to fetch holders at block %v: %v", block, err)
		}

		// save holders to db
		err = db.SaveCovalentTokenHolders(network, token, t.Unix(), holders)
		if err != nil {
			log.Fatalf("failed to save holders at [%v|%v|%v]: %v", t.Unix(), t.Format("02-01-2006"), block, err)
		}

		log.Printf("[%v|%v|%v]: %v holders", t.Unix(), t.Format("02-01-2006"), block, len(holders))
	}
	os.Exit(0)
}
