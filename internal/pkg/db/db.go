package db

import (
	"fmt"
	"log"
	"merch-shop/internal/app/models"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB инициализирует соединение с БД
func InitDB() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Ошибка получения *sql.DB: %v", err)
	}

	// Настройка параметров пула соединений
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Проверяем соединение
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	log.Println("Подключение к БД установлено")

	err = db.AutoMigrate(&models.User{}, &models.Merch{}, models.Transaction{})
	if err != nil {
		log.Fatal("Ошибка миграции таблиц: ", err)
	}
	err = AddMerch(db)
	if err != nil {
		log.Fatal("Ошибка добавления мерча: ", err)
	}

	return db
}

func AddMerch(db *gorm.DB) error {
	merches := []models.Merch{
		{Name: "t-shirt", Price: 80},
		{Name: "cup", Price: 20},
		{Name: "book", Price: 50},
		{Name: "pen", Price: 10},
		{Name: "powerbank", Price: 200},
		{Name: "hoody", Price: 300},
		{Name: "umbrella", Price: 200},
		{Name: "socks", Price: 10},
		{Name: "wallet", Price: 50},
		{Name: "pink-hoody", Price: 500},
	}

	result := db.Create(&merches)
	return result.Error
}
