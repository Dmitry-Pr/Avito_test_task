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
	defer func() {
		log.Println("Закрываем соединение с базой данных")
		if err := database.Close(); err != nil {
			log.Println("Ошибка закрытия соединения с базой данных:", err)
		}
	}()

	// Создаем DI-контейнер
	container := di.BuildDependencies(database)

	// Запускаем сервер
	srv := server.NewServer(cfg, container)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
