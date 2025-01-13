CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE credit_cards (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES users(id),
    number INTEGER NOT NULL,
    expiration_date TEXT NOT NULL,
    secret INTEGER NOT NULL
);

CREATE TABLE machines (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    buyer_id INTEGER REFERENCES users(id),
    owner_id INTEGER NOT NULL REFERENCES users(id),
    ram INTEGER NOT NULL,
    cpu INTEGER NOT NULL,
    memory INTEGER NOT NULL,
    key TEXT
);