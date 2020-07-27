package utils


import(
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func QueryCoinGecko(endpoint string,ptr interface{})error{
	url:="https://api.coingecko.com/api/v3"
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