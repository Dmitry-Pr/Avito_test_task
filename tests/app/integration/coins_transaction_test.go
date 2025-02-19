package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"merch-shop/internal/app/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"merch-shop/internal/app/di"
	"merch-shop/internal/config"
	"merch-shop/internal/server"
)

// TestSendCoinIntegration проверяет сценарий передачи монет между пользователями.
func TestSendCoinIntegration(t *testing.T) {
	testCases := []struct {
		name                  string
		senderUsername        string
		senderPassword        string
		senderInitialCoins    int
		receiverUsername      string
		receiverPassword      string
		receiverInitialCoins  int
		transferAmount        int
		expectedStatusCode    int
		expectedSenderCoins   int
		expectedReceiverCoins int
		errorMessage          string // Ожидаемое сообщение об ошибке (если есть)
	}{
		{
			name:                  "Успешная передача монет",
			senderUsername:        "alice",
			senderPassword:        "testpass",
			senderInitialCoins:    1000,
			receiverUsername:      "bob",
			receiverPassword:      "testpass",
			receiverInitialCoins:  1000,
			transferAmount:        200,
			expectedStatusCode:    http.StatusOK,
			expectedSenderCoins:   800,  // 1000 - 200
			expectedReceiverCoins: 1200, // 1000 + 200
			errorMessage:          "",
		},
		{
			name:                  "Недостаточно средств",
			senderUsername:        "alice",
			senderPassword:        "testpass",
			senderInitialCoins:    100, // недостаточно средств
			receiverUsername:      "bob",
			receiverPassword:      "testpass",
			receiverInitialCoins:  1000,
			transferAmount:        200,
			expectedStatusCode:    http.StatusBadRequest,
			expectedSenderCoins:   100,  // баланс не изменился
			expectedReceiverCoins: 1000, // баланс не изменился
			errorMessage:          "недостаточно монет",
		},
		{
			name:                  "Передача монет себе",
			senderUsername:        "alice",
			senderPassword:        "testpass",
			senderInitialCoins:    1000,
			receiverUsername:      "alice", // попытка перевода самому себе
			receiverPassword:      "testpass",
			receiverInitialCoins:  1000,
			transferAmount:        100,
			expectedStatusCode:    http.StatusBadRequest,
			expectedSenderCoins:   1000, // баланс не изменился
			expectedReceiverCoins: 1000,
			errorMessage:          "нельзя отправить монеты себе",
		},
		{
			name:                  "Пользователь получатель не найден",
			senderUsername:        "alice",
			senderPassword:        "testpass",
			senderInitialCoins:    1000,
			receiverUsername:      "charlie", // получатель не создан
			receiverPassword:      "testpass",
			receiverInitialCoins:  0,
			transferAmount:        100,
			expectedStatusCode:    http.StatusBadRequest,
			expectedSenderCoins:   1000,
			expectedReceiverCoins: 0,
			errorMessage:          "пользователь получатель не найден",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Сброс тестовой базы
			tearDownDB()
			TestDB = setUpDB()
			if TestDB == nil {
				t.Fatalf("Ошибка установки тестовой БД для теста %s", tc.name)
			}

			// Предварительное создание пользователя-отправителя с нужным балансом
			sender := models.User{
				Username: tc.senderUsername,
				// Пароль не имеет значения, т.к. при аутентификации, если пользователя нет,
				// он создаётся автоматически, но для контроля баланса мы создаём его вручную.
				Password: "$2a$10$gpqFoNQodGlJV.mx04Jfj.Y88eXm8tXWy5ZgRU6uepYhZHc4HUuqW",
				Coins:    tc.senderInitialCoins,
			}
			TestDB.Create(&sender)

			// Если тестовый случай подразумевает существование получателя (и не передача самому себе)
			if tc.receiverUsername != tc.senderUsername && tc.expectedReceiverCoins != 0 {
				receiver := models.User{
					Username: tc.receiverUsername,
					Password: "$2a$10$gpqFoNQodGlJV.mx04Jfj.Y88eXm8tXWy5ZgRU6uepYhZHc4HUuqW",
					Coins:    tc.receiverInitialCoins,
				}
				TestDB.Create(&receiver)
			}

			// Загружаем конфигурацию и собираем зависимости
			cfg := config.LoadConfig()
			dependencies := di.BuildDependencies(TestDB)
			srv := server.NewServer(cfg, dependencies)
			handler := srv.HTTPHandler()
			testServer := httptest.NewServer(handler)
			defer testServer.Close()

			// 1. Аутентификация отправителя
			authURL := testServer.URL + "/api/auth"
			authPayload := map[string]string{
				"username": tc.senderUsername,
				"password": tc.senderPassword,
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

			// 2. Передача монет: отправляем запрос на /api/sendCoin
			sendCoinURL := testServer.URL + "/api/sendCoin"
			sendCoinPayload := map[string]interface{}{
				"toUser": tc.receiverUsername,
				"amount": tc.transferAmount,
			}
			sendCoinPayloadBytes, err := json.Marshal(sendCoinPayload)
			if err != nil {
				t.Fatalf("Ошибка маршалинга запроса передачи монет: %v", err)
			}
			req, err := http.NewRequest("POST", sendCoinURL, bytes.NewBuffer(sendCoinPayloadBytes))
			if err != nil {
				t.Fatalf("Ошибка создания запроса на передачу монет: %v", err)
			}
			req.Header.Set("Authorization", "Bearer "+authRes.Token)
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			sendResp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Ошибка выполнения запроса передачи монет: %v", err)
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					return
				}
			}(sendResp.Body)
			if sendResp.StatusCode != tc.expectedStatusCode {
				bodyBytes, _ := io.ReadAll(sendResp.Body)
				t.Fatalf("Ожидался статус %d, получен %d, тело ответа: %s",
					tc.expectedStatusCode, sendResp.StatusCode, string(bodyBytes))
			}

			// Если ожидается сообщение об ошибке, проверяем его
			if tc.errorMessage != "" {
				bodyBytes, _ := io.ReadAll(sendResp.Body)
				if !bytes.Contains(bodyBytes, []byte(tc.errorMessage)) {
					t.Errorf("Ожидалось сообщение об ошибке '%s', получено '%s'",
						tc.errorMessage, string(bodyBytes))
				}
			}

			// 3. Получение информации для отправителя через /api/info
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
			var senderInfo InfoResponse
			if err := json.NewDecoder(infoResp.Body).Decode(&senderInfo); err != nil {
				t.Fatalf("Ошибка декодирования ответа информации: %v", err)
			}
			if senderInfo.Coins != tc.expectedSenderCoins {
				t.Errorf("У отправителя ожидалось %d монет, получено %d", tc.expectedSenderCoins, senderInfo.Coins)
			}

			// 4. Если получатель существует, аутентифицируем его и проверяем баланс
			if tc.receiverUsername != tc.senderUsername && tc.expectedReceiverCoins != 0 {
				authPayload = map[string]string{
					"username": tc.receiverUsername,
					"password": tc.receiverPassword,
				}
				authPayloadBytes, err = json.Marshal(authPayload)
				if err != nil {
					t.Fatalf("Ошибка маршалинга запроса аутентификации получателя: %v", err)
				}
				resp, err = http.Post(authURL, "application/json", bytes.NewBuffer(authPayloadBytes))
				if err != nil {
					t.Fatalf("Ошибка выполнения запроса аутентификации получателя: %v", err)
				}
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						return
					}
				}(resp.Body)
				if resp.StatusCode != http.StatusOK {
					t.Fatalf("Ожидался статус 200 при аутентификации получателя, получен %d", resp.StatusCode)
				}
				var receiverAuth AuthResponse
				if err := json.NewDecoder(resp.Body).Decode(&receiverAuth); err != nil {
					t.Fatalf("Ошибка декодирования ответа аутентификации получателя: %v", err)
				}
				if receiverAuth.Token == "" {
					t.Fatal("Получен пустой токен для получателя")
				}
				// Получаем информацию о получателе
				infoReq, err = http.NewRequest("GET", infoURL, nil)
				if err != nil {
					t.Fatalf("Ошибка создания запроса на получение информации для получателя: %v", err)
				}
				infoReq.Header.Set("Authorization", "Bearer "+receiverAuth.Token)
				infoResp, err = client.Do(infoReq)
				if err != nil {
					t.Fatalf("Ошибка выполнения запроса на получение информации для получателя: %v", err)
				}
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						return
					}
				}(infoResp.Body)
				if infoResp.StatusCode != http.StatusOK {
					bodyBytes, _ := io.ReadAll(infoResp.Body)
					t.Fatalf("Ожидался статус 200 при получении информации для получателя, получен %d, тело ответа: %s",
						infoResp.StatusCode, string(bodyBytes))
				}
				var receiverInfo InfoResponse
				if err := json.NewDecoder(infoResp.Body).Decode(&receiverInfo); err != nil {
					t.Fatalf("Ошибка декодирования ответа информации для получателя: %v", err)
				}
				if receiverInfo.Coins != tc.expectedReceiverCoins {
					t.Errorf("У получателя ожидалось %d монет, получено %d",
						tc.expectedReceiverCoins, receiverInfo.Coins)
				}
			}
		})
	}
}
