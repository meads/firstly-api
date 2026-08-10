package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	db "github.com/meads/firstly-api/internal/db/sqlc"
	"github.com/meads/firstly-api/internal/handler"
	"github.com/meads/firstly-api/internal/repository"
	"github.com/meads/firstly-api/internal/security"
	"github.com/meads/firstly-api/internal/service"
)

func dbConnect(retries int, dbUrl string) *sql.DB {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	retryCount := 0
	for {
		err := db.Ping()
		if err != nil {
			retryCount += 1
			time.Sleep(time.Second * 2)
			if retryCount == retries {
				log.Fatal(err)
			}
			fmt.Printf("Database connect retry attempt %d...\n", retryCount)
		} else {
			break
		}
	}

	return db
}

func main() {
	dbConnectionString := os.Getenv("DATABASE_URL")
	secretKey := os.Getenv("SECRET_KEY")
	// conn := dbConnect(10, dbURL)
	// defer conn.Close()
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbConnectionString)
	if err != nil {
		log.Fatalf("Unable to connect to the database %v\n", err)
	}
	defer pool.Close()

	m, err := migrate.New("file://internal/db/migration", dbConnectionString)
	if err != nil {
		log.Fatalf("error calling New with sql-migration tool: %s", err)
		return
	}
	m.Up()

	fmt.Print("\nmigrations were a success. 🎉\n")

	tokener := security.NewTokenManager(secretKey)
	hasher := security.NewHasher()
	// store := db.New(conn)
	queries := db.New(pool)

	sessionRepo := repository.NewSessionRepository(queries)
	userRepo := repository.NewUserRepository(queries)

	authService := service.NewAuthService(userRepo, sessionRepo, tokener, hasher)
	authHandler := handler.NewAuthHandler(authService)

	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	noteRepo := repository.NewNoteRepository(queries)
	noteService := service.NewNoteService(noteRepo, userRepo)
	noteHandler := handler.NewNoteHandler(noteService)

	router := handler.SetupRouter(authHandler, userHandler, noteHandler, tokener)

	err = router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
