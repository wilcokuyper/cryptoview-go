package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		return
	}

	if env, ok := os.LookupEnv("APP_ENV"); ok && env == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("unable to create zap logger")
	}
	defer logger.Sync()

	db, err := sql.Open("mysql", "")
	defer db.Close()
	if err != nil {
		logger.Fatal("unable to connect to database")
	}

	server := NewServer(logger, db)

	// Lookup port and start server
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	server.Run(port)
}
