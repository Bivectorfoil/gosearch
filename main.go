package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

type SearchInformation struct {
	SearchTime   int    `json:"searchTime"`
	TotalResults string `json:"totalResults"`
}

type Queries struct {
	PreviousPage []struct {
		StartIndex  int    `json:"startIndex"`
		SearchTerms string `json:"searchTerms"`
	}
	NextPage []struct {
		StartIndex  int    `json:"startIndex"`
		SearchTerms string `json:"searchTerms"`
	}
	Request []struct {
		SearchTerms string `json:"searchTerms"`
	}
}

type CSERespnse struct {
	Items      []Result          `json:"items"`
	SearchInfo SearchInformation `json:"searchInformation"`
	Queries    Queries           `json:"queries"`
}

func main() {
	initProxy()
	router := gin.Default()
	// load template file
	router.LoadHTMLGlob("templates/*")

	// basic route
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})
	// post route
	router.POST("/", func(c *gin.Context) {
		query := c.PostForm("search")
		resp := search(query, 1)
		results := &CSERespnse{}
		json.Unmarshal(resp, results)
		// POST redirect to avoid resubmit form
		c.Redirect(http.StatusMovedPermanently, "/")
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"query_item": query,
			"items":      results.Items,
			"searchInfo": results.SearchInfo,
			"results":    results,
		})

	})

	// basic search route
	router.GET("/result", func(c *gin.Context) {
		q := c.Query("q")
		startIndex, _ := strconv.Atoi(c.Query("startIndex"))
		resp := search(q, startIndex)
		results := &CSERespnse{}

		json.Unmarshal(resp, results)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"query_item": q,
			"items":      results.Items,
			"searchInfo": results.SearchInfo,
			"results":    results,
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}

func initProxy() {
	fmt.Println("init proxy")
	// init proxy with script/setProxy.sh file
	cmd := exec.Command("/bin/bash", "-c", "./script/setProxy.sh")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s\nErr msg: %s\n", err, stderr.String())
	} else {
		fmt.Printf("%s\n", out.String())
	}
}
