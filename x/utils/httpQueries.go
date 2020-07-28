package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//queryFromAPI is generic function that can query every api and endpoint
func queryFromAPI(url string, endpoint string, ptr interface{}) error {
	resp, err := http.Get(fmt.Sprintf("%s/%s", url, endpoint))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bz, &ptr); err != nil {
		log.Print(string(bz))
		return err
	}

	return nil
}

// QueryCoinGecko can query endpoint from coingecko
func QueryCoinGecko(endpoint string, ptr interface{}) error {
	url := "https://api.coingecko.com/api/v3"
	if err := queryFromAPI(url, endpoint, ptr); err != nil {
		return err
	}
	return nil
}
