package main

import (
	"database/sql"
	"log"
	"merch-store/internal/app/di"
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
	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(database)

	// Создаем DI-контейнер
	container := di.BuildDependencies(database)

	// Запускаем сервер
	srv := server.NewServer(cfg, container)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
