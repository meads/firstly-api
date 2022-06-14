package main

import (
	"context"

	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"

	api "github.com/heroku/firstly-api/db/api"
)

// TODO:
// - Create go types for nullable columns and create overrides for json encoding to handle db values in dtos
// Reconcile the differences between environments local/production. Favor working locally.
// Setup unit tests and integration tests.
// Fix authorization issues with android app POST request.

type Image struct {
	ID      int64  `json:"id"`
	Created string `json:"created"`
	Data    string `json:"data"`
	Deleted bool   `json:"deleted"`
}

func (to Image) fromDbAPIType(from *api.Image) *Image {
	if from == nil {
		return &Image{}
	}

	deleted := false
	if from.Deleted.Valid && from.Deleted.Int32 == 1 {
		deleted = true
	}

	return &Image{
		ID:      from.ID,
		Created: from.Created,
		Data:    from.Data,
		Deleted: deleted,
	}
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

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("error calling WithInstance: %s", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/sql/migrations", "postgres", driver)

	if err != nil {
		log.Fatalf("error calling NewWithDatabaseInstance: %s", err)
		return
	}

	m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run

	router.GET("/app/image/", func(c *gin.Context) {
		q := api.New(db)
		images, err := q.ListImages(context.Background())
		if err != nil {
			log.Fatalf("Error calling ListImages: %s", err)
			return
		}
		dtoImages := []Image{}
		for _, img := range images {
			dtoImages = append(dtoImages, *Image{}.fromDbAPIType(&img))
		}
		c.JSON(http.StatusOK, dtoImages)
	})

	router.POST("/app/image/", func(c *gin.Context) {
		var image Image
		if err := c.BindJSON(&image); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		q := api.New(db)
		apiImage, err := q.CreateImage(context.Background(), image.Data)
		if err != nil {
			log.Fatalf("Error calling CreateImage: %s", err)
			return
		}

		c.JSON(http.StatusCreated, Image{}.fromDbAPIType(&apiImage))
	})

	router.Run(":" + os.Getenv("PORT"))
}
