package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"

	"github.com/joho/godotenv"
)

const BASE_URL = "https://www.googleapis.com/customsearch/v1?"

func search(queryItem string, startIndex int) ([]byte, error) {
	// read CSE_ID and CSE_KEY from .env file with godotenv package
	godotenv.Load()
	CSEID := os.Getenv("CSE_ID")
	CSEKEY := os.Getenv("CSE_KEY")
	PROXY_HOST := os.Getenv("PROXY_HOST")

	params := map[string]interface{}{
		"cx":    CSEID,
		"q":     url.QueryEscape(queryItem),
		"key":   CSEKEY,
		"num":   10,
		"start": startIndex,
	}

	// creating proxy string
	proxyURL, err := url.Parse(PROXY_HOST)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{Transport: transport}

	rawUrl, err := url.Parse(BASE_URL + paramsToQuery(params))
	if err != nil {
		// log err and return
		log.Print(err)
		return nil, err
	}
	// generating the HTTP GET request through http proxy
	request, err := http.NewRequest("GET", rawUrl.String(), nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return data, nil
}

func paramsToQuery(params map[string]interface{}) string {
	var query string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		query += fmt.Sprintf("%s=%v&", k, params[k])
	}
	return query[:len(query)-1]
}
