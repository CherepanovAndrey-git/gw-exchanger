-- name: GetRate :one
SELECT rate FROM exchange_rates
WHERE from_currency = $1 AND to_currency = $2;

-- name: GetAllRates :many
SELECT from_currency, to_currency, rate FROM exchange_rates;

-- name: UpdateRate :exec
UPDATE exchange_rates
SET rate = $3, updated_at = NOW()
WHERE from_currency = $1 AND to_currency = $2;