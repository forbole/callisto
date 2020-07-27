package utils


import(
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func QueryCoinGecko(endpoint string,ptr []interface{})error{
	url:="https://api.coingecko.com/api/v3/"
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
		return err
	}

	return nil

}