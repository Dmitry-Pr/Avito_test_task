package main

import (
	"log"
	"merch-store/internal/config"
	"merch-store/internal/pkg/db"
	"merch-store/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Загружаем конфиг
	cfg := config.LoadConfig()

	// Инициализируем базу данных
	database := db.InitDB()
	defer database.Close()

	// Создаем сервер и пробрасываем базу
	srv := server.NewServer(cfg, database)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
