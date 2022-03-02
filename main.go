package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()
	// basic route
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome")
	})

	// more result
	router.GET("/more", func(c *gin.Context) {
		c.String(200, "More page")
	})

	// query result page
	router.GET("/result", func(c *gin.Context) {
		c.String(200, "Result page")
	})

	// POST query and return json string with query item
	router.POST("/", func(c *gin.Context) {
		query := c.PostForm("query")
		fmt.Printf("Your query item is: %s", query)
		c.JSON(200, gin.H{"result": query, "status": http.StatusOK})
	})

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
	fmt.Printf("cse_id is %s, cse_key is %s", CSEID, CSEKEY)

	params := map[string]interface{}{
		"cx":    CSEID,
		"q":     "golang",
		"key":   CSEKEY,
		"num":   10,
		"start": 1,
	}

	// todo: use with proxy at test mode
	URL := "https://www.googleapis.com/customsearch/v1?"
	// read resp from request GET through http proxy
	resp, err := http.Get(URL + paramsToQuery(params))
	if err != nil {
		log.Fatal("request error:", err)
	}

	// read resp body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read body error:", err)
	}
	fmt.Printf("resp body %s", body)

	router.Run() // listen and serve on 0.0.0.0:8080
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
