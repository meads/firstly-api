package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

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
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatalf("errror reading request body: %s", err)
			return
		}
		log.Println(string(jsonData))
		jsonHeaders, err := json.Marshal(c.Request.Header)
		if err != nil {
			log.Fatalf("error marshalling the request headers: %s", err)
			return
		}
		log.Println(string(jsonHeaders))
		c.JSON(http.StatusOK, strings.Join([]string{string(jsonHeaders), string(jsonData)}, ""))
		// return

		// var image Image

		// // bind the json to the struct
		// err := c.BindJSON(&image)
		// if err != nil {
		// 	log.Fatalf("error binding json to image struct: %s", err)
		// 	return
		// }

		// // insert the new image record
		// q := api.New(db)
		// apiImage, err := q.CreateImage(context.Background(), image.Data)
		// if err != nil {
		// 	log.Fatalf("Error calling CreateImage: %s", err)
		// 	return
		// }

		// c.IndentedJSON(http.StatusCreated, apiImage)
	})

	router.Run(":" + os.Getenv("PORT"))
}
