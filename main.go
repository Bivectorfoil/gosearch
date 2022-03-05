package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type result struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"Snippet"`
}

type resultItems struct {
	Items []result `json:"items"`
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
		results := &resultItems{}
		json.Unmarshal(resp, results)
		// POST redirect to avoid resubmit form
		c.Redirect(http.StatusMovedPermanently, "/")
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":      "query finished at",
			"query_item": query,
			"items":      results.Items,
		})

	})

	// basic search route
	router.GET("/search", func(c *gin.Context) {
		resp := search("golang")
		results := &resultItems{}

		json.Unmarshal(resp, results)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"search_item": "Golang",
			"items":       results.Items,
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
