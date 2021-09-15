CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE SCHEMA IF NOT EXISTS warranty;
DROP TABLE IF EXISTS warranty.warranty cascade;
 
CREATE TABLE warranty.warranty
(
    id            SERIAL CONSTRAINT warranty_pkey PRIMARY KEY,
    comment_      VARCHAR(1024),
    item_uid      UUID         NOT NULL CONSTRAINT idx_warranty_item_uid UNIQUE,
    status        VARCHAR(255) NOT NULL,
    warranty_date TIMESTAMP    NOT NULL
);