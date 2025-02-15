-- Таблица пользователей
CREATE TABLE users
(
    id           BIGSERIAL PRIMARY KEY,
    username     VARCHAR(255) UNIQUE NOT NULL,
    coin_balance BIGINT DEFAULT 1000 CHECK (coin_balance >= 0)
);

-- Таблица транзакций
CREATE TABLE transactions
(
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT REFERENCES users (id) NOT NULL,
    type         VARCHAR(255)                 NOT NULL, -- 'transfer', 'purchase'
    amount       BIGINT                       NOT NULL,
    from_user_id BIGINT REFERENCES users (id),          -- Может быть NULL для 'purchase'
    to_user_id   BIGINT REFERENCES users (id),          -- Может быть NULL для 'purchase'
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица товаров
CREATE TABLE merch
(
    id    BIGSERIAL PRIMARY KEY,
    name  VARCHAR(255) UNIQUE NOT NULL,
    price BIGINT              NOT NULL
);
