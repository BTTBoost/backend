package main

import (
	"awake/internal/lib"
	"bufio"
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

// TODO: move to lib
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

// TODO: support pagination for tokens with > 100k holders
// Fetches and saves to db holders of all tokens at some time from Covalent API.
func main() {
	network := 1

	// load .env into ENV
	err := godotenv.Overload()
	if err != nil {
		log.Fatal(err)
	}

	// init params
	if len(os.Args) != 3 {
		log.Fatal("wrong arguments. pass [tokenlist.txt] [time]")
	}

	// token list arg
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	tokens := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		token := scanner.Text()
		if !lib.IsValidAddress(token) {
			log.Fatalf("invalid token '%v'", token)
		}
		tokens = append(tokens, token)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading token list: %v", err)
	}

	// time
	timeParam := os.Args[2]
	timestamp, err := strconv.ParseInt(timeParam, 10, 64)
	if err != nil || timestamp <= 0 {
		log.Fatalf("invalid time param '%v': %v", timeParam, err)
	}
	t := time.Unix(timestamp, 0).Truncate(time.Hour * 24)

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
	block := timeBlocks[t.Unix()]

	log.Printf("fetching holders for %v tokens at block %v", len(tokens), block)

	for i, token := range tokens {
		// fetch
		var holders []lib.CovalentHolder
		for {
			holders, err = fetchHolders(network, token, block)
			if err == nil {
				break
			}
			log.Printf("[%v|%v] failed to fetch holders for token: %v", i, token, err)
			time.Sleep(5 * time.Second)
		}

		// save to db
		err = db.SaveLastTokenHolders(network, token, holders)
		if err != nil {
			log.Fatalf("[%v|%v] failed to save token holders: %v", i, token, err)
		}
		log.Printf("[%v|%v] fetched and saved %v holders", i, token, len(holders))
	}

	os.Exit(0)
}
