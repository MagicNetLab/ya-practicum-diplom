CREATE TABLE users (
    uid uuid PRIMARY KEY,
    login VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);