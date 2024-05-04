package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func ConnectDb() *sql.DB {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.SetDefault("type", "dev")
	typeName := viper.Get("type").(string)

	fmt.Println(typeName)

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Error", err)

	}

	dbUsername := viper.Get("DB_USERNAME").(string)
	dbPassword := viper.Get("DB_PASSWORD").(string)
	dbConnectionUrl := viper.Get("DB_CONNECTION_URL").(string)
	dbPort := viper.Get("DB_PORT").(string)
	dbName := viper.Get("DB_NAME").(string)

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUsername, dbPassword, dbConnectionUrl, dbPort, dbName)
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	return db

}
