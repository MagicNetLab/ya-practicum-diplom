CREATE TABLE notes (
    id uuid PRIMARY KEY NOt NULL,
    uid uuid NOT NULL,
    title varchar NOT NULL,
    content text NOT NULL,
    meta varchar,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
)