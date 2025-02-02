-- schema.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password BYTEA NOT NULL
);

CREATE TABLE machines (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    buyer_id INTEGER REFERENCES users(id),
    owner_id INTEGER NOT NULL REFERENCES users(id),
    ram INTEGER NOT NULL,
    cpu INTEGER NOT NULL,
    memory INTEGER NOT NULL,
    key TEXT,
    host VARCHAR(255) NOT NULL,
    ssh_user VARCHAR(255) NOT NULL
);