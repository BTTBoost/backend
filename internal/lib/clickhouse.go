package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func QueryClickhouse(query string) (string, error) {
	url := fmt.Sprintf("%v/?query=%v&user=%v&password=%v",
		os.Getenv("CLICKHOUSE_URL"),
		url.QueryEscape(query),
		os.Getenv("CLICKHOUSE_USER"),
		os.Getenv("CLICKHOUSE_PASSWORD"),
	)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	result := string(body)

	return result, nil
}
