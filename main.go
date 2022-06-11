package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
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

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	time.Sleep(5 * time.Second)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file:///db/sql/migrations",
		"postgres", driver)
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
