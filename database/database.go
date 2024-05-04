package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func ConnectDb() *sql.DB {
	viper.AutomaticEnv()

	dbUsername := viper.GetString("DB_USERNAME")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbConnectionUrl := viper.GetString("DB_CONNECTION_URL")
	dbPort := viper.GetString("DB_PORT")
	dbName := viper.GetString("DB_NAME")

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUsername, dbPassword, dbConnectionUrl, dbPort, dbName)
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	return db

}
