package main

import (
	"awake/internal/lib"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const NETWORK = "matic"
const CONTRACT_ADDRESS = "0x8dff5e27ea6b7ac08ebfdf9eb090f32ee9a30fcf"
const EVENT = "Deposit"
const FROM_BLOCK = 20836867
const TO_BLOCK = 24559039
const PAGE_SIZE = 10000

func fetchEvents(page int64) ([]lib.BitqueryEvent, error) {
	apiKey := os.Getenv("BITQUERY_API_KEY")
	url := "https://graphql.bitquery.io/"
	body := fmt.Sprintf(`{
		"query":"query ($network: EthereumNetwork!,$contract: String!,$event: String!, $fromBlock: Int!, $toBlock: Int!, $limit: Int!, $offset: Int!) {\n  ethereum(network: $network) {\n    smartContractEvents(\n      options: {asc: \"block.height\", limit: $limit, offset: $offset}\n      smartContractEvent: {is: $event }\n      smartContractAddress: {is: $contract}\n      height: {gteq: $fromBlock, lt: $toBlock }\n    ) {\n      block {\n        height\n        timestamp {\n          iso8601\n          unixtime\n        }\n      }\n      transaction {\n        hash\n      }\n      arguments {\n        value\n        argument\n      }\n    }\n  }\n}\n",
		"variables": {
			"network": "%v",
			"contract": "%v",
			"event": "%v",
			"fromBlock": %v,
			"toBlock": %v,
			"limit": %v,
			"offset": %v
		}
	}`, NETWORK, CONTRACT_ADDRESS, EVENT, FROM_BLOCK, TO_BLOCK, PAGE_SIZE, page*PAGE_SIZE,
	)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Fatalf("failed to create req")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	result := &lib.BitqueryResponse{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, errors.New(result.Errors[0].Message)
	}

	return result.Data.Events, nil
}

func main() {
	// load .env into ENV
	err := godotenv.Overload()
	if err != nil {
		log.Fatal(err)
	}

	// init params
	if len(os.Args) != 2 {
		log.Fatal("wrong arguments. pass [page]")
	}

	pageParam := os.Args[1]
	page, err := strconv.ParseInt(pageParam, 10, 64)
	if err != nil {
		log.Fatalf("invalid page param '%v': %v", pageParam, err)
	}
	log.Printf(`fetching events from bitquery:
- network: %v
- contract: %v
- event: %v
- fromBlock: %v
- toBlock: %v
- limit: %v
- offset: %v`,
		NETWORK, CONTRACT_ADDRESS, EVENT, FROM_BLOCK, TO_BLOCK, PAGE_SIZE, page*PAGE_SIZE,
	)

	// connect db
	db, err := lib.CreateDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// network
	var network int64
	switch NETWORK {
	case "ethereum":
		network = 1
	case "matic":
		network = 137
	}

	for {
		// fetch
		events, err := fetchEvents(page)
		if err != nil {
			log.Fatalf("failed to fetch events for page %v: %v", page, err)
		}
		if len(events) == 0 {
			log.Printf("fetch success with no events")
			break
		}

		// save to db
		err = db.SaveAaveDeposits(network, events)
		if err != nil {
			log.Fatalf("failed to save events to db: %v", err)
		}

		log.Printf("[%v - %v] fetched and saved %v events: [%v -> %v]",
			page*PAGE_SIZE+1, (page+1)*PAGE_SIZE, len(events),
			events[0].BlockData.Block, events[len(events)-1].BlockData.Block,
		)
		page++
	}

	os.Exit(0)
}

func logParsedEvents(events []lib.BitqueryEvent) {
	for i, e := range events {
		deposit, err := lib.ParseAaveDepositBitqueryEvent(e)
		if err != nil {
			log.Fatalf("failed to parse deposit: %v", err)
		}
		log.Printf("[%v] deposit = %v", i, deposit)
	}
}
