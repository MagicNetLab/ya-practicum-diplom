CREATE TABLE cards (
   id uuid PRIMARY KEY NOT NULL,
   uid uuid NOT NULL,
   name varchar NULL,
   number varchar NOT NULL UNIQUE,
   mask varchar NOT NULL,
   month integer NULL,
   year integer NULL,
   cvc varchar NULL,
   pin varchar NULL,
   created_at timestamp NOT NULL
);

CREATE INDEX cards_uid_idx ON cards USING HASH(uid);