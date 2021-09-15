CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE SCHEMA IF NOT EXISTS store;
DROP TABLE IF EXISTS store.users cascade;
 
CREATE TABLE IF NOT EXISTS store.users
(
    id       SERIAL CONSTRAINT users_pkey PRIMARY KEY,
    name     VARCHAR(255) NOT NULL CONSTRAINT idx_user_name UNIQUE,
    user_uid UUID         NOT NULL CONSTRAINT idx_user_user_uid UNIQUE
);

INSERT INTO store.users(name, user_uid) VALUES ('Ivan', '6d2cb5a0-943c-4b96-9aa6-89eac7bdfd2b');
