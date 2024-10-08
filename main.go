package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/heroku/x/hmetrics/onload"

	db "github.com/meads/firstly-api/db"
	http_api "github.com/meads/firstly-api/http"
	"github.com/meads/firstly-api/security"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening postgres database '%s'", err)
		return
	}
	defer conn.Close()

	m, err := migrate.New("file://./db/migration", dbURL)

	if err != nil {
		log.Fatalf("error calling New with sql-migration tool: %s", err)
		return
	}

	m.Up()

	fmt.Print("\nmigrations were a success. 🎉\n")

	claimer := security.NewClaimsValidator()
	hasher := security.NewHasher()
	router := gin.Default()
	store := db.NewStore(conn)

	server := http_api.NewFirstlyServer(claimer, hasher, router, store)

	err = server.Start(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
