package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/heroku/x/hmetrics/onload"

	db "github.com/meads/firstly-api/db/sqlc"
	http_api "github.com/meads/firstly-api/http"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	fmt.Printf("Here is the DATABASE_URL: %s", dbURL)

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening postgres driver using url '%s', '%s'", dbURL, err)
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

	store := db.NewStore(conn)
	server := http_api.NewServer(store)
	// config.ServerAddress
	err = server.Start(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
