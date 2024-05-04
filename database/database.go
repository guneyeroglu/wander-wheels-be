package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDb() *sql.DB {
	// err := godotenv.Load()

	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// dbUsername := os.Getenv("DB_USERNAME")
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbConnectionUrl := os.Getenv("DB_CONNECTION_URL")
	// dbPort := os.Getenv("DB_PORT")
	// dbName := os.Getenv("DB_NAME")
	// dbSslMode := os.Getenv("DB_SSL_MODE")
	dbConnection := os.Getenv("DB_CONNECTION")

	// connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", dbUsername, dbPassword, dbConnectionUrl, dbPort, dbName, dbSslMode)
	db, err := sql.Open("postgres", dbConnection)

	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	return db

}
