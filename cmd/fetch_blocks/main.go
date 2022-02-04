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

type GetBlockResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func fetchBlock(network int, timestamp int64) (int64, error) {
	var apiKey string
	var baseUrl string
	switch network {
	case 1:
		apiKey = os.Getenv("ETHERSCAN_API_KEY")
		baseUrl = "https://api.etherscan.io"
	case 137:
		apiKey = os.Getenv("POLYGONSCAN_API_KEY")
		baseUrl = "https://api.polygonscan.com"
	default:
		log.Fatalf("unsupported network %v", network)
	}
	url := fmt.Sprintf("%v/api?apikey=%v&module=block&action=getblocknobytime&closest=before&timestamp=%v",
		baseUrl, apiKey, timestamp,
	)

	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	result := &GetBlockResponse{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	block, err := strconv.ParseInt(result.Result, 10, 64)
	if err != nil {
		return 0, err
	}
	return block, nil
}

func main() {
	// load .env into ENV
	err := godotenv.Overload()
	if err != nil {
		log.Fatal(err)
	}

	// read args
	if len(os.Args) != 3 {
		log.Fatal("wrong arguments, pass: [network] [from_timestamp]")
	}

	var network int
	networkParam := os.Args[1]
	switch networkParam {
	case "ethereum":
		network = 1
	case "polygon":
		network = 137
	default:
		log.Fatalf("unsupported network %v", networkParam)
	}

	fromParam := os.Args[2]
	fromTimestamp, err := strconv.ParseInt(fromParam, 10, 64)
	if err != nil || fromTimestamp < 0 {
		log.Fatalf("invalid from_timestamp param '%v'", fromParam)
	}

	to := time.Now().Truncate(time.Hour * 24)
	from := time.Unix(fromTimestamp, 0).Truncate(time.Hour * 24)
	if from.After(to) {
		log.Fatal("from_timestamp must be earlier than now")
	}
	log.Printf("updating blocks from [%v] to [%v]...", from.Format("02-01-2006"), to.Format("02-01-2006"))

	// connect db
	db, err := lib.CreateDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// TODO: save in batches
	// fetch and save blocks
	for t := from; !t.After(to); t = t.Add(time.Hour * 24) {
		block, err := fetchBlock(network, t.Unix())
		if err != nil {
			log.Fatalf("failed to fetch block for %v: %v", t.String(), err)
		}

		log.Printf("[%v|%v] block = %v", network, t.Format("02-01-2006"), block)

		err = db.SaveBlockTime(network, t.Unix(), block)
		if err != nil {
			log.Fatalf("failed to save block to db for %v: %v", t.String(), err)
		}
	}

	os.Exit(0)
}
