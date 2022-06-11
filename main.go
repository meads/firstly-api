package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/heroku/firstly-api/db/api"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

func main() {
	// Run migrations

	m, err := migrate.New("file:///db/sql/migrations", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("error initializing postgres migration tool: %s", err)
		return
	}

	err = m.Up()
	if err != nil {
		log.Fatalf("error running migration: %s", err)
		return
	}

	m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run

	type Image struct {
		Data string `json:"data"`
	}

	router := gin.New()
	router.Use(gin.Logger())

	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(dbURL)
		log.Fatal(err)
	}
	defer db.Close()

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

		// insert the new image record
		q := api.Queries{}
		apiImage, err := q.CreateImage(context.Background(), image.Data)
		if err != nil {
			log.Fatalf("Error calling CreateImage: %s", err)
			return
		}

		c.IndentedJSON(http.StatusCreated, apiImage)
	})

	router.Run(":" + os.Getenv("PORT"))
}
