FROM golang:1.23.1-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Устанавливаем необходимые инструменты для тестирования (например, mockgen)
RUN apk add --no-cache make gcc musl-dev # добавляем gcc и musl-dev
RUN go install github.com/golang/mock/mockgen@v1.6.0

# Запускаем тесты
CMD ["go", "test", "-cover", "-race", "-v", "-coverpkg=./...", "./tests/..."]