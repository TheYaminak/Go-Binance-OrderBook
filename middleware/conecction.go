package middleware

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetConnection() *sql.DB {
	connStr := os.Getenv("POSTGRES_CONNECTION_STRING")
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func CreateTables(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS ask (
		symbol VARCHAR(20) NOT NULL,
		price DOUBLE PRECISION DEFAULT NULL,
		s_price VARCHAR(33) NOT NULL,
		quantity DOUBLE PRECISION DEFAULT NULL,
		PRIMARY KEY (symbol, s_price)
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bid (
		symbol VARCHAR(20) NOT NULL,
		price DOUBLE PRECISION DEFAULT NULL,
		s_price VARCHAR(33) NOT NULL,
		quantity DOUBLE PRECISION DEFAULT NULL,
		PRIMARY KEY (symbol, s_price)
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS order_book (
		id SERIAL PRIMARY KEY,
		symbol TEXT NOT NULL,
		bids JSONB NOT NULL,
		asks JSONB NOT NULL
	)`)
	if err != nil {
		return err
	}

	return nil
}
