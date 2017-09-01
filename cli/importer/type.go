package importer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gquental/pokedex/config"
	"github.com/gquental/pokedex/data"
)

func ImportTypes() {
	fetchTypesFromAPI()
}

func fetchTypesFromAPI() {
	typeCountCh := make(chan int)
	typesCh := make(chan data.DamageType)

	var totalTypes int
	var countTypes int

	go fetchTypeList("", typeCountCh, typesCh)
	totalTypes = <-typeCountCh
	close(typeCountCh)

	for {
		select {
		case typeDetail := <-typesCh:
			fmt.Println(typeDetail)
			countTypes++
		}

		if countTypes == totalTypes {
			close(typesCh)
			break
		}
	}

}

type listType struct {
	Count   int
	Results []struct {
		Url  string
		Name string
	}
	Next string
}

func fetchTypeList(url string, typeCountCh chan int, typesCh chan data.DamageType) {
	firstTime := false
	if url == "" {
		firstTime = true
		url = fmt.Sprintf("%s%s", config.Config.APIEndpoint, "type")
	}

	body, err := requestAPI(url)
	if err != nil {

	}

	types := listType{}
	json.Unmarshal(body, &types)

	if firstTime {
		typeCountCh <- types.Count
	}

	for _, item := range types.Results {
		fmt.Println(item)
	}
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
