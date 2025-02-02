-- query.sql
-- name: CreateUser :one
INSERT INTO users (name, email, password)
VALUES ($1, $2, $3::bytea)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateMachine :one
INSERT INTO machines (
    name,
    ram,
    cpu,
    memory,
    key,
    owner_id,
    buyer_id,
    host,
    ssh_user
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    NULL,
    $7,
    $8
) RETURNING *;

-- name: GetMachineByID :one
SELECT * FROM machines WHERE id = $1;

-- name: ListMachinesByOwnerID :many
SELECT * FROM machines WHERE owner_id = $1;

-- name: ListMachinesByBuyerID :many
SELECT ram, cpu, memory, name FROM machines WHERE buyer_id = $1;

-- name: ListAvailableMachines :many
SELECT * FROM machines WHERE buyer_id IS NULL;

-- name: UpdateMachineBuyer :one
UPDATE machines 
SET buyer_id = $1, key = $2
WHERE id = $3 AND buyer_id IS NULL
RETURNING *;

-- name: GetMachineByNameAndOwner :one
SELECT m.* 
FROM machines m
JOIN users u ON m.owner_id = u.id
WHERE m.name = $1 AND u.name = $2;