package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		return
	}

	logger, err := setupLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	connStr, ok := os.LookupEnv("DB_CONN_STRING")
	if !ok {
		log.Fatal("missing DB_CONN_STRING environment variable")
	}

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal("Unable to connect to the database")
	}
	defer db.Close()

	server := NewServer(logger, db)

	// Lookup port and start server
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	server.Run(port)
}

func setupLogger() (*zap.Logger, error) {
	if env, ok := os.LookupEnv("APP_ENV"); ok && env == "development" {
		return zap.NewDevelopment()
	} else {
		return zap.NewProduction()
	}
}
