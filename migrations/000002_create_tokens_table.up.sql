CREATE TABLE tokens (
    id BIGSERIAL PRIMARY KEY,
    token TEXT NOT NULL UNIQUE,
    uid uuid NOT NULL,
    is_refresh BOOLEAN DEFAULT false,
    expired TIMESTAMP NOT NULL
);