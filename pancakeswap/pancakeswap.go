// Package PancakeSwap provides a very simple interface to PancakeSwap's V2 API.
package pancakeswap

import (
	"io/ioutil"
	"net/http"
)

// apiRequest makes an API request to PancakeSwap V2.
func apiRequest(endpoint string) (body []byte, err error) {
	apiURL := "https://api.pancakeswap.info/api/v2/"
	res, err := http.Get(apiURL + endpoint)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	return
}
