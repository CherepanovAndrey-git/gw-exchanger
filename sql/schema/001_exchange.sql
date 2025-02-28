-- +goose Up
CREATE TABLE exchange_rates (
        id SERIAL PRIMARY KEY,
        from_currency VARCHAR(3) NOT NULL,
        to_currency VARCHAR(3) NOT NULL,
        rate NUMERIC(18, 6) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_exchange_rates_pair ON exchange_rates (from_currency, to_currency);


INSERT INTO exchange_rates (from_currency, to_currency, rate) VALUES
      ('USD', 'RUB', 93.000000),
      ('RUB', 'USD', 0.010800),
      ('USD', 'EUR', 0.920000),
      ('EUR', 'USD', 1.087000),
      ('EUR', 'RUB', 101.500000),
      ('RUB', 'EUR', 0.009850);

-- +goose Down
DROP TABLE IF EXISTS exchange_rates;