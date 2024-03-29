package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

type CSEResponse struct {
	Items      []Result          `json:"items"`
	SearchInfo SearchInformation `json:"searchInformation"`
	Queries    Queries           `json:"queries"`
}

func main() {
	// set run mode
	initEnv()
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
		resp, err := search(query, 1)
		if err != nil {
			fmt.Println(err)
			c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{
				"error": err,
			})
			return
		}
		results := &CSEResponse{}
		err = json.Unmarshal(resp, results)
		if err != nil {
			render500ErrorResponse(c, err)
			return
		}
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
		resp, err := search(q, startIndex)
		if err != nil {
			fmt.Println(err)
			c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{
				"error": err,
			})
			return
		}
		results := &CSEResponse{}

		err = json.Unmarshal(resp, results)
		if err != nil {
			render500ErrorResponse(c, err)
			return
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"query_item": q,
			"items":      results.Items,
			"searchInfo": results.SearchInfo,
			"results":    results,
		})
	})

	router.POST("/result", func(c *gin.Context) {
		query := c.PostForm("search")
		resp, err := search(query, 1)
		if err != nil {
			fmt.Println(err)
			c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{
				"error": err,
			})
			return
		}
		results := &CSEResponse{}
		err = json.Unmarshal(resp, results)
		if err != nil {
			render500ErrorResponse(c, err)
			return
		}
		// POST redirect to avoid resubmit form
		c.Redirect(http.StatusMovedPermanently, "/")
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"query_item": query,
			"items":      results.Items,
			"searchInfo": results.SearchInfo,
			"results":    results,
		})
	})

	err := router.Run()
	if err != nil {
		render500ErrorResponse(nil, err)
		return
	} // listen and serve on 0.0.0.0:8080
}

func render500ErrorResponse(c *gin.Context, err error) {
	if c == nil {
		c = &gin.Context{}
	}
	fmt.Println(err)
	c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{
		"error": err,
	})
}

func initProxy() {
	if gin.Mode() != gin.DebugMode {
		return
	}
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

func initEnv() {
	RunMode := os.Getenv("RUN_MODE")
	switch RunMode {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}
