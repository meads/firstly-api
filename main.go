package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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
	photos := []Photo{}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
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
	router.GET("/app/photo/", func(c *gin.Context) {
		b, err := json.Marshal(photos)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		c.String(http.StatusOK, string(b))
	})

	router.POST("/app/photo/", func(c *gin.Context) {
		var photo Photo

		if err := c.BindJSON(&photo); err != nil {
			return
		}

		tx, err := db.Begin()
		if err != nil {
			fmt.Errorf("error starting tx %s", err)
			return
		}
		tx.Exec("")

		// Add the new album to the slice.
		photos = append(photos, photo)
		c.IndentedJSON(http.StatusCreated, photo)
	})

	router.Run(":" + port)
}
