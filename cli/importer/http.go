package importer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func requestAPI(endpoint string) ([]byte, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
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

type ListDefinition struct {
	Count   int
	Results []struct {
		Url  string
		Name string
	}
	Next string
}

type ItemEntry struct {
	Url  string
	Type ImportableDetail
}

type ExpandEntry struct {
	Url    string
	Index  int
	Assign func(interface{}, []byte, int) ImportableDetail
}

type Importable interface {
	GetCount() int
	GetNext() string
	EraseNext()
	List() []ItemEntry
	GetEndpoint() string
}

type ImportableDetail interface {
	Store()
	ExpandDetails() []ExpandEntry
}

func fetchFromAPI(importEntry Importable) {
	countCh := make(chan int)
	entryCh := make(chan ImportableDetail)

	var total int
	var count int

	go fetchList(0, importEntry, importEntry.GetEndpoint(), countCh, entryCh)
	total = <-countCh
	close(countCh)

	for {
		select {
		case entryDetail := <-entryCh:
			entryDetail.Store()
			count++
		}

		if count == total {
			close(entryCh)
			break
		}
	}

}

func fetchList(callNumber int, importEntry Importable, url string, countCh chan int, entryCh chan ImportableDetail) {
	body, err := requestAPI(url)
	if err != nil {
		logrus.Error(fmt.Errorf("Problem in list request: %v", err))
		close(entryCh)
		if callNumber == 0 {
			countCh <- 0
		}
		return
	}

	importEntry.EraseNext()

	json.Unmarshal(body, importEntry)
	if callNumber == 0 {
		countCh <- importEntry.GetCount()
	}

	if importEntry.GetNext() != "" {
		callNumber++
		go fetchList(callNumber, importEntry, importEntry.GetNext(), countCh, entryCh)
	}

	for _, item := range importEntry.List() {
		fetchDetail(item, entryCh)
	}
}

func fetchDetail(entry ItemEntry, entryCh chan ImportableDetail) {
	body, err := requestAPI(entry.Url)
	if err != nil {
		logrus.Error(fmt.Errorf("Problem in detail request: %v", err))
		close(entryCh)
		return
	}

	json.Unmarshal(body, entry.Type)

	for _, item := range entry.Type.ExpandDetails() {
		body, err := requestAPI(item.Url)
		if err != nil {
			logrus.Error(fmt.Errorf("Problem in expanded detail request: %v", err))
			continue
		}

		entry.Type = item.Assign(entry.Type, body, item.Index)
	}

	entryCh <- entry.Type
}
