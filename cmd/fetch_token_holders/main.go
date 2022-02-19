package main

import (
	"awake/internal/lib"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const PAGE_SIZE = 1000000

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
// TODO: pagination is probably broken on server
func fetchHolders(network int, token string, block int64, pageSize int, page int) ([]lib.CovalentHolder, error) {
	apiKey := os.Getenv("COVALENT_API_KEY")

	url := fmt.Sprintf(
		"https://api.covalenthq.com/v1/%v/tokens/%v/token_holders/"+
			"?block-height=%v&page-number=%v&page-size=%v&key=%v",
		network, token, block, page, pageSize, apiKey,
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
	if result.Error {
		return nil, errors.New(result.ErrorMessage)
	}

	// TODO: uncomment once server pagination is fixed
	if len(result.Data.Items) == pageSize {
		return nil, fmt.Errorf("too many holders, pagination required")
	}

	return result.Data.Items, nil
}

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
		page := 0
		for {
			// fetch page with repeat on error
			var pageHolders []lib.CovalentHolder
			for {
				pageHolders, err = fetchHolders(network, token, block, PAGE_SIZE, page)
				if err == nil {
					log.Printf("[%v|%v|%v] fetched %v holders", i, token, page, len(pageHolders))
					break
				}
				log.Printf("[%v|%v|%v] failed to fetch holders for token: %v", i, token, page, err)
				time.Sleep(5 * time.Second)
			}

			holders = append(holders, pageHolders...)
			if len(pageHolders) < PAGE_SIZE {
				break
			}
			page++
		}

		// save to db
		err = db.SaveLastTokenHolders(network, token, holders)
		if err != nil {
			log.Fatalf("[%v|%v] failed to save token holders: %v", i, token, err)
		}
		log.Printf("[%v|%v] saved %v holders", i, token, len(holders))
	}

	os.Exit(0)
}
