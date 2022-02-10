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

const OPENSEA_ADDRESS = "0x7be8076f4ea4a4ad08075c2508e481d6c946d12b"
const TRADE_TOPIC = "0xc4109843e0b7d514e4c093114b863f8e7d8d9a458c372cd51bfe526b588006c9"
const PAGE_SIZE = 25000

type Response struct {
	Data         ResponseData `json:"data"`
	Error        bool         `json:"error"`
	ErrorMessage string       `json:"error_message"`
}

type ResponseData struct {
	Items []lib.OpenseaTradeEvent `json:"items"`
}

func fetchTrades(fromBlock int64) ([]lib.OpenseaTradeEvent, error) {
	apiKey := os.Getenv("COVALENT_API_KEY")
	toBlock := fromBlock + 1000000 // max allowed

	url := fmt.Sprintf(
		"https://api.covalenthq.com/v1/1/events/topics/%v/"+
			"?sender-address=%v&starting-block=%v&ending-block=%v&page-number=0&page-size=%v&key=%v",
		TRADE_TOPIC, OPENSEA_ADDRESS, fromBlock, toBlock,
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

// Fetches and saves to db OpenSea trade events from Covalent API.
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
	log.Printf("fetching trades from [%v]...", fromBlock)

	// connect db
	db, err := lib.CreateDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// fetch and save trades in a loop
	for {
		trades, err := fetchTrades(fromBlock)
		if err != nil {
			log.Fatalf("failed to fetch trades from block %v: %v", fromBlock, err)
		}
		if len(trades) == 0 {
			log.Printf("fetch success with no trades")
			break
		}

		err = db.SaveOpenseaTrades(trades)
		if err != nil {
			log.Fatalf("failed to save trades to db: %v", err)
		}

		log.Printf("[%v -> %v] fetched and saved %v trades", trades[0].BlockHeight, trades[len(trades)-1].BlockHeight, len(trades))
		fromBlock = trades[len(trades)-1].BlockHeight
	}

	os.Exit(0)
}
