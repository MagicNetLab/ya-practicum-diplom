CREATE TABLE accounts (
    id uuid PRIMARY KEY,
    uid uuid NOT NULL,
    login VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    url VARCHAR NOT NULL,
    description VARCHAR,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)