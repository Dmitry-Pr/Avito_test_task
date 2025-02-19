FROM golang:1.23.1-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN apk add --no-cache make gcc musl-dev # добавляем gcc и musl-dev
RUN go install github.com/golang/mock/mockgen@v1.6.0

COPY . .

# Запускаем тесты
CMD ["make", "test"]
