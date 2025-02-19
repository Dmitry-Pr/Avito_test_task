package integration_test

import (
	"fmt"
	"log"
	"merch-shop/internal/pkg/db"
	"os"
	"testing"

	"gorm.io/gorm"
)

var TestDB *gorm.DB // Глобальная переменная для соединения с БД

func TestMain(m *testing.M) {
	err := os.Setenv("JWT_SECRET_KEY", "testsecret")
	if err != nil {
		return
	}
	setUpDB()
	tearDownDB()
	m.Run()
	tearDownDB()
}

func setUpDB() *gorm.DB {
	return db.InitDB()
}

func tearDownDB() {
	if TestDB != nil {
		var tableNames []string

		// Получаем список таблиц, исключая таблицу migrations
		err := TestDB.Raw(`
                        SELECT table_name
                        FROM information_schema.tables
                        WHERE table_schema = 'public'
                        AND table_name != 'migrations'
                `).Scan(&tableNames).Error

		if err != nil {
			log.Println("Ошибка получения списка таблиц:", err)
			return
		}

		for _, tableName := range tableNames {
			err := TestDB.Exec("TRUNCATE TABLE " + tableName + " CASCADE;").Error
			if err != nil {
				log.Println("Ошибка очистки таблицы", tableName, ":", err)
			}
		}

		sqlDB, err := TestDB.DB()
		if err != nil {
			log.Println("Ошибка получения *sql.DB при закрытии:", err)
			return
		}
		err = sqlDB.Close()
		if err != nil {
			log.Println("Ошибка закрытия соединения с БД:", err)
		}
		fmt.Println("Соединение с тестовой БД закрыто")
	}
}
