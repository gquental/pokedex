package importer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type listType struct {
	Count   int
	Results []struct {
		Url  string
		Name string
	}
	Next string
}

func requestAPI(endpoint string) ([]byte, error) {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error in the request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
