package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"

	_ "github.com/lib/pq"

	api "github.com/heroku/firstly-api/db/api"
	migrate "github.com/rubenv/sql-migrate"
)

// TODO:
// Create go types for nullable columns and create overrides for json encoding to handle db values in dtos
// Reconcile the differences between environments local/production. Favor working locally.
// Setup unit tests and integration tests.
// Fix authorization issues with android app POST request.

type Image struct {
	ID      int          `json:"id"`
	Created string       `json:"created"`
	Data    string       `json:"data"`
	Deleted sql.NullBool `json:"deleted"`
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

	// OR: Read migrations from a folder:
	migrations := &migrate.FileMigrationSource{
		Dir: "db/sql/migrations",
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("migration execution failed: %s", err)
		return
	}
	fmt.Printf("Applied %d migrations!\n", n)

	router.GET("/app/image/", func(c *gin.Context) {
		q := api.New(db)
		images, err := q.ListImages(context.Background())
		if err != nil {
			log.Fatalf("Error calling ListImages: %s", err)
			return
		}

		c.JSON(http.StatusOK, images)
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

		c.IndentedJSON(http.StatusCreated, apiImage)
	})

	router.Run(":" + os.Getenv("PORT"))
}
