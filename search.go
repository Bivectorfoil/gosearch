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

func search(useProxy bool) []byte {
	// test go request
	// resp, err := requests.Get("http://www.zhanluejia.net.cn")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(resp.Text())

	// read CSE_ID and CSE_KEY from .env file with godotenv package
	godotenv.Load()
	CSEID := os.Getenv("CSE_ID")
	CSEKEY := os.Getenv("CSE_KEY")
	PROXY_HOST := ""
	if useProxy {
		// read proxy url from .env file with godotenv package
		PROXY_HOST = os.Getenv("PROXY_HOST")
	}

	fmt.Printf("cse_id is %s\n, cse_key is %s\n, proxy is %s\n", CSEID, CSEKEY, PROXY_HOST)

	params := map[string]interface{}{
		"cx":    CSEID,
		"q":     "golang",
		"key":   CSEKEY,
		"num":   10,
		"start": 1,
	}

	// creating proxy string
	proxyURL, err := url.Parse(PROXY_HOST)
	if err != nil {
		log.Fatal(err)
	}
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{Transport: transport}

	URL := "https://www.googleapis.com/customsearch/v1?"
	url, err := url.Parse(URL + paramsToQuery(params))
	if err != nil {
		log.Fatal(err)
	}
	// generating the HTTP GET request through http proxy
	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return data
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
