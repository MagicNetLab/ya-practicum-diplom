CREATE TABLE files (
    id uuid PRIMARY KEY,
    uid uuid NOT NULL,
    name varchar NOT NULL,
    path varchar NOT NULL UNIQUE,
    meta text,
    size bigint NOT NULL,
    created_at timestamp NOT NULL
);

CREATE INDEX files_uid_idx ON files using hash(uid);