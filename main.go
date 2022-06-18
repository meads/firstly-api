package main

import (
	"context"
	"database/sql"
	"io/ioutil"

	_ "github.com/lib/pq"

	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/heroku/x/hmetrics/onload"

	api "github.com/heroku/firstly-api/db/api"
)

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
	headersLogger := func(ctx *gin.Context) {
		log.Println(ctx.Request.Header)
		log.Println("-------------------------------------------------------")
	}
	router := gin.New()
	router.Use(headersLogger, gin.Logger())

	dbURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening postgres driver using url '%s', '%s'", dbURL, err)
		return
	}
	defer db.Close()

	m, err := migrate.New("file://./db/sql/migrations", dbURL)

	if err != nil {
		log.Fatalf("error calling New with sql-migration tool: %s", err)
		return
	}
	m.Up()

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
			body, _ := ioutil.ReadAll(c.Request.Body)
			log.Println(string(body))
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
