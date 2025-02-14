package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

var DB *sql.DB // Глобальная переменная для БД

// InitDB инициализирует соединение с БД
func InitDB() *sql.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	// Настройка параметров пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	log.Println("Подключение к БД установлено")
	DB = db
	return db
}
