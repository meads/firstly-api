package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heroku/firstly-api/db/api"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

func main() {
	type Image struct {
		Data string `json:"data"`
	}

	router := gin.New()
	router.Use(gin.Logger())

	// ####### DB Connection related----------------------------------
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	time.Sleep(5 * time.Second)
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// ######## Routes------------------------------------------------
	router.GET("/app/image/", func(c *gin.Context) {
		q := api.Queries{}
		images, err := q.ListImages(context.Background())
		if err != nil {
			log.Fatalf("Error calling ListImages: %s", err)
			return
		}

		c.JSON(http.StatusOK, images)
	})

	router.POST("/app/image/", func(c *gin.Context) {
		var image Image

		// bind the json to the struct
		err := c.BindJSON(&image)
		if err != nil {
			log.Fatalf("error binding json to image struct: %s", err)
			return
		}

		// create the db operation params
		params := api.CreateImageParams{}
		params.Data = image.Data
		params.Name = "test-name"

		// insert the new image record
		q := api.Queries{}
		apiImage, err := q.CreateImage(context.Background(), params)
		if err != nil {
			log.Fatalf("Error calling CreateImage: %s", err)
			return
		}

		c.IndentedJSON(http.StatusCreated, apiImage)
	})

	router.Run(":" + os.Getenv("PORT"))
}
