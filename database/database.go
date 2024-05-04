package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDb() *sql.DB {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbConnectionUrl := os.Getenv("DB_CONNECTION_URL")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUsername, dbPassword, dbConnectionUrl, dbPort, dbName)
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	return db

}
