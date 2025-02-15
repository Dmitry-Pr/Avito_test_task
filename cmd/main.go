package main

import (
	"log"
	"merch-shop/internal/app/di"
	"merch-shop/internal/config"
	"merch-shop/internal/pkg/db"
	"merch-shop/internal/server"

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
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("Ошибка получения *sql.DB из Gorm")
	}
	defer func() {
		log.Println("Закрываем соединение с базой данных")
		if err := sqlDB.Close(); err != nil {
			log.Println("Ошибка закрытия соединения с базой данных:", err)
		}
	}()

	// Создаем DI-контейнер
	dependencies := di.BuildDependencies(database)

	// Запускаем сервер
	srv := server.NewServer(cfg, dependencies)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
