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
	type Photo struct {
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
		photos, err := q.ListImages(context.Background())
		if err != nil {
			log.Fatalf("Error calling ListImages: %s", err)
			return
		}

		c.JSON(http.StatusOK, photos)
	})

	router.POST("/app/image/", func(c *gin.Context) {
		var photo Photo

		// bind the json to the struct
		err := c.BindJSON(&photo)
		if err != nil {
			log.Fatalf("error binding json to photo struct: %s", err)
			return
		}

		// create the db operation params
		params := api.CreateImageParams{}
		params.Data = photo.Data
		params.Name = "test-name"

		// insert the new image record
		q := api.Queries{}
		image, err := q.CreateImage(context.Background(), params)
		if err != nil {
			log.Fatalf("Error calling CreateImage: %s", err)
			return
		}

		c.IndentedJSON(http.StatusCreated, image)
	})

	router.Run(":" + os.Getenv("PORT"))
}
