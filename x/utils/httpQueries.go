package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//queryFromAPI is generic function that can query every api and endpoint
func queryFromAPI(endpoint string, ptr interface{}) error {
	resp, err := http.Get(endpoint)
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
