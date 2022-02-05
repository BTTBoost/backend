package main

import (
	"awake/internal/lib"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const NETWORK = 1
const LENDING_POOL_V2_ADDRESS = "0x7d2768dE32b0b80b7a3454c06BdAc94A69DDc7A9"
const DEPOSIT_TOPIC = "0xde6857219544bb5b7746f48ed30be6386fefc61b2f864cacf559893bf50fd951"
const TO_BLOCK = 14142694
const PAGE_SIZE = 50000

type Response struct {
	Data         ResponseData `json:"data"`
	Error        bool         `json:"error"`
	ErrorMessage string       `json:"error_message"`
}

type ResponseData struct {
	Items []lib.CovalentEvent `json:"items"`
}

func fetchDeposits(fromBlock int64) ([]lib.CovalentEvent, error) {
	apiKey := os.Getenv("COVALENT_API_KEY")

	url := fmt.Sprintf(
		"https://api.covalenthq.com/v1/%v/events/topics/%v/"+
			"?sender-address=%v&starting-block=%v&ending-block=%v&page-number=0&page-size=%v&key=%v",
		NETWORK, DEPOSIT_TOPIC, LENDING_POOL_V2_ADDRESS, fromBlock, TO_BLOCK,
		PAGE_SIZE, apiKey,
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

	return result.Data.Items, nil
}

// Fetches and saves to db Aave v2 deposit events from Covalent API.
func main() {
	// load .env into ENV
	err := godotenv.Overload()
	if err != nil {
		log.Fatal(err)
	}

	// init params
	if len(os.Args) != 2 {
		log.Fatal("wrong arguments. pass [from_block]")
	}

	fromBlockParam := os.Args[1]
	fromBlock, err := strconv.ParseInt(fromBlockParam, 10, 64)
	if err != nil {
		log.Fatalf("invalid from_block param '%v': %v", fromBlockParam, err)
	}
	log.Printf("fetching events from [%v]...", fromBlock)

	// fetch
	events, err := fetchDeposits(fromBlock)
	if err != nil {
		log.Fatalf("failed to fetch events from block %v: %v", fromBlock, err)
	}
	if len(events) == 0 {
		log.Printf("fetch success with no events")
		return
	}

	// connect db
	db, err := lib.CreateDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// save to db
	err = db.SaveAaveDeposits(events)
	if err != nil {
		log.Fatalf("failed to save events to db: %v", err)
	}

	log.Printf("[%v -> %v] fetched and saved %v events", events[0].BlockHeight, events[len(events)-1].BlockHeight, len(events))
	os.Exit(0)
}
