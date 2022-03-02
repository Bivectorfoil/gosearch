package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	router.Run() // listen and serve on 0.0.0.0:8080
}
