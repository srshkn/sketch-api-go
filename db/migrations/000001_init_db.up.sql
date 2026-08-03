CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(25) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);
