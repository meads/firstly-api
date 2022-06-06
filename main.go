package main

import (
	"log"
	"net/http"
	"os"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	type Photo struct {
		Data   string  `json:"data"`
	}

	router := gin.New()
	router.Use(gin.Logger())
	photos := []Photo{}


	// list route for photos
	router.GET("/app/photo/", func(c *gin.Context) {
		b, err := json.Marshal(photos)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return;
		}
		c.String(http.StatusOK, string(b))
	})

	router.POST("/app/photo/", func(c *gin.Context) {
		var photo Photo

		if err := c.BindJSON(&photo); err != nil {
			return
		}

		// Add the new album to the slice.
		photos = append(photos, photo)
		c.IndentedJSON(http.StatusCreated, photo)
	})

	router.Run(":" + port)
}
