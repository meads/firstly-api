package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"

	migrate "github.com/golang-migrate/migrate/v4"
	postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"

	api "github.com/heroku/firstly-api/db/api"
)

type Image struct {
	Data string `json:"data"`
}

func main() {

	router := gin.New()
	router.Use(gin.Logger())

	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening postgres driver using url '%s', '%s'", dbURL, err)
	}
	defer db.Close()

	// Run migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("error initializing postgres migration tool with db handle: %s", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///db/sql/migrations",
		"postgres",
		driver)
	if err != nil {
		log.Fatalf("error running migration: %s", err)
		return
	}

	m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run

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
