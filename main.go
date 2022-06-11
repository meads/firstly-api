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

	q := api.New(db)

	router.GET("/app/image/", func(c *gin.Context) {

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
		apiImage, err := q.CreateImage(context.Background(), image.Data)
		if err != nil {
			log.Fatalf("Error calling CreateImage: %s", err)
			return
		}

		c.IndentedJSON(http.StatusCreated, apiImage)
	})

	router.Run(":" + os.Getenv("PORT"))
}
