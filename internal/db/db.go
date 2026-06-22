package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // driver PostgreSQL
)

type DBConnection struct {
	HOST     string
	PORT     string
	USER     string
	PASSWORD string
	DB_NAME  string
	SSL_MODE string
}

func ContextDBURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSL_MODE"),
	)
}

func ContextDB() string {
	dbContext := DBConnection{
		HOST:     os.Getenv("DB_HOST"),
		PORT:     os.Getenv("DB_PORT"),
		USER:     os.Getenv("DB_USER"),
		PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_NAME:  os.Getenv("DB_NAME"),
		SSL_MODE: os.Getenv("SSL_MODE"),
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbContext.HOST, dbContext.PORT, dbContext.USER, dbContext.PASSWORD, dbContext.DB_NAME, dbContext.SSL_MODE,
	)
}

func Connect() *sql.DB {
	strConnection := ContextDB()

	db, err := sql.Open("postgres", strConnection)
	if err != nil {
		log.Fatal(err)
	}

	// ping connection
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
