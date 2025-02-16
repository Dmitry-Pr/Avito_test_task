# Тестовое задание стажёра Backend-направления (зимняя волна 2025)

## Оглавление

- [Магазин мерча](#магазин-мерча)
- [Описание задачи](#описание-задачи)
- [Запуск](#запуск)
- [Стек технологий](#стек-технологий)
- [Архитектура](#архитектура)
- [Заполнение `.env` файла](#заполнение-env-файла)
- [API Specification (OpenAPI)](#api-specification-openapi)

## Магазин мерча

В Авито существует внутренний магазин мерча, где сотрудники могут приобретать товары за монеты (coin). Каждому новому сотруднику выделяется 1000 монет, которые можно использовать для покупки товаров. Кроме того, монеты можно передавать другим сотрудникам в знак благодарности или как подарок.

## Описание задачи

Необходимо реализовать сервис, который позволит сотрудникам обмениваться монетками и приобретать на них мерч. Каждый сотрудник должен иметь возможность видеть:

- Список купленных им мерчовых товаров
- Сгруппированную информацию о перемещении монеток в его кошельке, включая:
- Кто ему передавал монетки и в каком количестве
- Кому сотрудник передавал монетки и в каком количестве

Количество монеток не может быть отрицательным, запрещено уходить в минус при операциях с монетками.

## Запуск

Чтобы запустить приложение, необходимо:
- иметь открытый порт 8080 
- запустить docker на устройстве
- [создать .env файл в корне проекта](#заполнение-env-файла) (__ВНИМАНИЕ__, в целях упрощения запуска проверяющим, я добавил .env файл в проект, в проде так делать нельзя!)
- выполнить команду: ```docker compose up -d --build```
- приложение будет доступно по адресу: `localhost:8080/api/auth`

## Стек технологий

- Go
- net/http
- gorm
- PostgreSQL

## Архитектура

Выбрана __слоистая архитектура__, есть слои:
- models
- repositories
- services
- handlers

Слои связаны с помощью __Dependency Injection__

Дополнительно используется __middleware__, чтобы оборачивать в них все запросы

## Заполнение `.env` файла

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=shop
SERVER_ADDRESS=:8080
JWT_SECRET_KEY=df3ce9cd-6084-47ac-bd89-a4d89bdad6f3
```

## API Specification (OpenAPI)

```yaml
openapi: 3.0.0
info:
  title: API Avito shop
  version: 1.0.0

servers:
  - url: http://localhost:8080

security:
  - BearerAuth: []

paths:
  /api/info:
    get:
      summary: Получить информацию о монетах, инвентаре и истории транзакций.
      security:
        - BearerAuth: []
      responses:
        '200':
          description: Успешный ответ.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/InfoResponse'
        '400':
          description: Неверный запрос.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '401':
          description: Неавторизован.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Внутренняя ошибка сервера.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /api/sendCoin:
    post:
      summary: Отправить монеты другому пользователю.
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/SendCoinRequest'
      responses:
        '200':
          description: Успешный ответ.
        '400':
          description: Неверный запрос.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '401':
          description: Неавторизован.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Внутренняя ошибка сервера.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /api/buy/{item}:
    get:
      summary: Купить предмет за монеты.
      security:
        - BearerAuth: []
      parameters:
        - name: item
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Успешный ответ.
        '400':
          description: Неверный запрос.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '401':
          description: Неавторизован.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Внутренняя ошибка сервера.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /api/auth:
    post:
      summary: Аутентификация и получение JWT-токена. При первой аутентификации пользователь создается автоматически. 
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/AuthRequest'
      responses:
        '200':
          description: Успешная аутентификация.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AuthResponse'
        '400':
          description: Неверный запрос.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '401':
          description: Неавторизован.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: В