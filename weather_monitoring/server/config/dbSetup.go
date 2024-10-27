package config

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// type dbConfig struct {
// 	db *sql.DB
// }

func ConnectDB(connStr string) (*sql.DB, error) {
	// dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable"
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Println("Failed to connect to PostgreSQL database:", err)
	// 	return nil, err
	// }

	db, err := sql.Open("pgx", connStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}

type Config struct {
	DatabaseURL       string
	OpenWeatherAPIKey string
}

func Load() *Config {
	return &Config{
		DatabaseURL:       "postgres://tsdbadmin:egcdui1z6x8kxlkw@nuo1krqx3h.ck7w6yea7m.tsdb.cloud.timescale.com:31394/tsdb?sslmode=require",
		OpenWeatherAPIKey: "eabed5afa938251da25ac3a124c19c56",
	}
}
