package main

import (
	"encoding/json"
	"net/http"

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
		StartIndex int `json:"startIndex"`
	}
	NextPage []struct {
		StartIndex int `json:"startIndex"`
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
		resp := search(query)
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
	router.GET("/search", func(c *gin.Context) {
		resp := search("golang")
		results := &CSERespnse{}

		json.Unmarshal(resp, results)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"search_item": "Golang",
			"items":       results.Items,
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
