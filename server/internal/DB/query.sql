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
    gpu,
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
    $7,
    NULL,
    $8,
    $9
) RETURNING *;

-- name: GetMachineByID :one
SELECT * FROM machines WHERE id = $1;

-- name: ListMachinesByOwnerID :many
SELECT * FROM machines WHERE owner_id = $1;

-- name: ListMachinesByBuyerID :many
SELECT ram, cpu, gpu, memory, name, owner_id FROM machines WHERE buyer_id = $1;

-- name: ListAvailableMachines :many
SELECT * FROM machines WHERE buyer_id IS NULL;

-- name: FilterAvailableMachines :many
SELECT * FROM machines
WHERE buyer_id IS NULL
  AND ($1::int IS NULL OR cpu >= $1)
  AND ($2::int IS NULL OR ram >= $2)
  AND ($3::int IS NULL OR gpu >= $3)
ORDER BY id DESC;

-- name: UpdateMachineBuyer :one
UPDATE machines 
SET buyer_id = $1, key = $2
WHERE id = $3 AND buyer_id IS NULL
RETURNING *;

-- name: GetMachineByNameAndOwner :one
SELECT * FROM machines 
WHERE name = $1 AND owner_id = $2;