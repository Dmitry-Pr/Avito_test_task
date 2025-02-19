package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"merch-shop/internal/app/models"
	"merch-shop/internal/pkg/errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"merch-shop/internal/app/di"
	"merch-shop/internal/config"
	"merch-shop/internal/server"
)

// AuthResponse описывает структуру ответа от /api/auth.
type AuthResponse struct {
	Token string `json:"token"`
}

// InfoResponse описывает структуру ответа от /api/info.
type InfoResponse struct {
	Coins     int `json:"coins"`
	Inventory []struct {
		Type     string `json:"type"`
		Quantity int    `json:"quantity"`
	} `json:"inventory"`
	// Поле coinHistory опущено для краткости
}

func TestBuyMerchIntegration(t *testing.T) {
	testCases := []struct {
		name               string
		username           string
		password           string
		merchType          string
		expectedStatusCode int
		expectedCoins      int
		users              []models.User
		expectedInventory  []struct {
			Type     string `json:"type"`
			Quantity int    `json:"quantity"`
		}
		errorMessage errors.ErrorResponse // Сообщение об ошибке (если ожидается)
	}{
		{
			name:               "Успешная покупка t-shirt",
			username:           "testuser",
			password:           "testpass",
			merchType:          "t-shirt",
			expectedStatusCode: http.StatusOK,
			expectedCoins:      1000 - 80, // Цена t-shirt
			expectedInventory: []struct {
				Type     string `json:"type"`
				Quantity int    `json:"quantity"`
			}{
				{Type: "t-shirt", Quantity: 1},
			},
		},
		{
			name:               "Успешная покупка cup",
			username:           "testuser",
			password:           "testpass",
			merchType:          "cup",
			expectedStatusCode: http.StatusOK,
			expectedCoins:      1000 - 20, // Цена cup
			expectedInventory: []struct {
				Type     string `json:"type"`
				Quantity int    `json:"quantity"`
			}{
				{Type: "cup", Quantity: 1},
			},
		},
		{
			name:               "Недостаточно средств",
			username:           "testuser",
			password:           "testpass",
			merchType:          "powerbank",           // Цена 200
			expectedStatusCode: http.StatusBadRequest, // Или другой код, если обрабатываете эту ситуацию иначе
			expectedCoins:      100,
			users: []models.User{{Username: "testuser",
				Password: "$2a$10$gpqFoNQodGlJV.mx04Jfj.Y88eXm8tXWy5ZgRU6uepYhZHc4HUuqW", Coins: 100}},
			expectedInventory: []struct { // Инвентарь не должен измениться
				Type     string `json:"type"`
				Quantity int    `json:"quantity"`
			}{},
			errorMessage: errors.ErrorResponse{Errors: "недостаточно монет"}, // Ожидаемое сообщение об ошибке
		},
		{
			name:               "Неверный тип мерча",
			username:           "testuser",
			password:           "testpass",
			merchType:          "unknown",
			expectedStatusCode: http.StatusBadRequest, // Или другой код, если обрабатываете эту ситуацию иначе
			expectedCoins:      1000,
			expectedInventory: []struct { // Инвентарь не должен измениться
				Type     string `json:"type"`
				Quantity int    `json:"quantity"`
			}{},
			errorMessage: errors.ErrorResponse{Errors: "мерч не найден: record not found"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup для *текущего* тестового случая
			tearDownDB()
			TestDB = setUpDB()
			if TestDB == nil {
				t.Fatal("Ошибка установки тестовой БД для теста", tc.name)
			}
			TestDB.Create(&tc.users)

			// Загружаем конфигурацию
			cfg := config.LoadConfig()

			// Собираем зависимости
			dependencies := di.BuildDependencies(TestDB)

			// Создаем экземпляр сервера
			srv := server.NewServer(cfg, dependencies)

			// Получаем HTTP-обработчик сервера.
			handler := srv.HTTPHandler()
			testServer := httptest.NewServer(handler)
			defer testServer.Close()

			// 1. Аутентификация: получаем JWT-токен
			authURL := testServer.URL + "/api/auth"
			authPayload := map[string]string{
				"username": tc.username,
				"password": tc.password,
			}
			authPayloadBytes, err := json.Marshal(authPayload)
			if err != nil {
				t.Fatalf("Ошибка маршалинга запроса аутентификации: %v", err)
			}
			resp, err := http.Post(authURL, "application/json", bytes.NewBuffer(authPayloadBytes))
			if err != nil {
				t.Fatalf("Ошибка выполнения запроса аутентификации: %v", err)
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					return
				}
			}(resp.Body)
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("Ожидался статус 200 при аутентификации, получен %d", resp.StatusCode)
			}
			var authRes AuthResponse
			if err := json.NewDecoder(resp.Body).Decode(&authRes); err != nil {
				t.Fatalf("Ошибка декодирования ответа аутентификации: %v", err)
			}
			if authRes.Token == "" {
				t.Fatal("Получен пустой токен")
			}

			// 2. Покупка мерча: отправляем запрос на покупку
			buyURL := testServer.URL + "/api/buy/" + tc.merchType // Используем tc.merchType
			req, err := http.NewRequest("POST", buyURL, nil)
			if err != nil {
				t.Fatalf("Ошибка создания запроса на покупку мерча: %v", err)
			}
			req.Header.Set("Authorization", "Bearer "+authRes.Token)
			client := &http.Client{}
			buyResp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Ошибка выполнения запроса на покупку мерча: %v", err)
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					return
				}
			}(buyResp.Body)
			if buyResp.StatusCode != tc.expectedStatusCode { // Проверяем ожидаемый статус код
				bodyBytes, _ := io.ReadAll(buyResp.Body)
				t.Fatalf("Ожидался статус %d при покупке мерча, получен %d, тело ответа: %s",
					tc.expectedStatusCode, buyResp.StatusCode, string(bodyBytes))
			}

			// Проверяем сообщение об ошибке, если оно ожидается
			if tc.errorMessage.Errors != "" {
				bodyBytes, _ := io.ReadAll(buyResp.Body)
				responseString := string(bodyBytes)
				if !bytes.Contains(bodyBytes, []byte(tc.errorMessage.Errors)) {
					t.Errorf("Ожидалось сообщение об ошибке '%s', получено '%s'", tc.errorMessage, responseString)
				}
			}

			// 3. Получаем информацию о пользователе для проверки корректного списания монет и обновления инвентаря
			infoURL := testServer.URL + "/api/info"
			infoReq, err := http.NewRequest("GET", infoURL, nil)
			if err != nil {
				t.Fatalf("Ошибка создания запроса на получение информации: %v", err)
			}
			infoReq.Header.Set("Authorization", "Bearer "+authRes.Token)
			infoResp, err := client.Do(infoReq)
			if err != nil {
				t.Fatalf("Ошибка выполнения запроса на получение информации: %v", err)
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					return
				}
			}(infoResp.Body)
			if infoResp.StatusCode != http.StatusOK {
				bodyBytes, _ := io.ReadAll(infoResp.Body)
				t.Fatalf("Ожидался статус 200 при получении информации, получен %d, тело ответа: %s",
					infoResp.StatusCode, string(bodyBytes))
			}
			var infoRes InfoResponse
			if err := json.NewDecoder(infoResp.Body).Decode(&infoRes); err != nil {
				t.Fatalf("Ошибка декодирования ответа информации: %v", err)
			}

			// Проверяем количество монет
			if infoRes.Coins != tc.expectedCoins {
				t.Errorf("Ожидалось, что количество монет будет %d, но получено %d",
					tc.expectedCoins, infoRes.Coins)
			}

			// Проверяем инвентарь
			for _, expectedItem := range tc.expectedInventory {
				found := false
				for _, item := range infoRes.Inventory {
					if item.Type == expectedItem.Type && item.Quantity == expectedItem.Quantity {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Ожидаемый предмет '%s' в количестве %d не найден в инвентаре пользователя",
						expectedItem.Type, expectedItem.Quantity)
				}
			}
		})
	}
}
