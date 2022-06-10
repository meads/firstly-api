package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heroku/firstly-api/db/api"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

const (
	host     = "ec2-52-204-195-41.compute-1.amazonaws.com"
	port     = 5432
	user     = "irdanpwkdvbzxg"
	password = "c29b40a6619957f7b572795b81f7805414be54ccd3675c07761f4cc61894d83e"
	dbname   = "d89dudhb3lei05"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	type Photo struct {
		Data string `json:"data"`
	}

	router := gin.New()
	router.Use(gin.Logger())

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	time.Sleep(5 * time.Second)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	// list route for photos
	router.GET("/app/image/", func(c *gin.Context) {
		q := api.Queries{}
		photos, err := q.ListImages(c.Request.Context())
		if err != nil {
			fmt.Printf("Error calling ListImages: %s", err)
			return
		}

		c.JSON(http.StatusOK, photos)
	})

	router.POST("/app/image/", func(c *gin.Context) {
		var photo Photo

		// bind the json to the struct
		err := c.BindJSON(&photo)
		if err != nil {
			fmt.Printf("error binding json to photo struct: %s", err)
			return
		}

		// create the db operation params
		params := api.CreateImageParams{}
		params.Data = photo.Data
		params.Name = "test-name"

		// insert the new image record
		q := api.Queries{}
		image, err := q.CreateImage(c.Request.Context(), params)
		if err != nil {
			fmt.Printf("Error calling CreateImage: %s", err)
			return
		}

		c.IndentedJSON(http.StatusCreated, image)
	})

	router.Run(":" + port)
}
