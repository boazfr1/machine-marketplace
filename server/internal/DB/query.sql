-- query.sql
-- name: CreateUser :one
INSERT INTO users (name, email, password)
VALUES ($1, $2, $3::bytea)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateCreditCard :one
INSERT INTO credit_cards (owner_id, number, expiration_date, secret)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCreditCardsByOwnerID :many
SELECT * FROM credit_cards WHERE owner_id = $1;

-- name: CreateMachine :one
INSERT INTO machines (name, owner_id, ram, cpu, memory)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMachineByID :one
SELECT * FROM machines WHERE id = $1;

-- name: ListMachinesByOwnerID :many
SELECT * FROM machines WHERE owner_id = $1;

-- name: ListMachinesByBuyerID :many
SELECT * FROM machines WHERE buyer_id = $1;

-- name: ListAvailableMachines :many
SELECT * FROM machines WHERE buyer_id IS NULL;

-- name: UpdateMachineBuyer :one
UPDATE machines 
SET buyer_id = $1, key = $2
WHERE id = $3 AND buyer_id IS NULL
RETURNING *;

-- name: GetUserCreditCards :many
SELECT * FROM credit_cards WHERE owner_id = $1;