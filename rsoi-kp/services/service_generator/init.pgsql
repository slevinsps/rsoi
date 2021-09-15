CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE SCHEMA IF NOT EXISTS generator;
DROP TABLE IF EXISTS generator.generator cascade;
 
CREATE TABLE IF NOT EXISTS generator.generator
(
    id       SERIAL CONSTRAINT users_pkey PRIMARY KEY,
    data_uuid UUID         NOT NULL CONSTRAINT idx_data_uuid UNIQUE,
    equipment_uuid    UUID       NOT NULL   REFERENCES equipment.equipment(equipment_uuid) ON DELETE CASCADE,
    temperature     float  DEFAULT 0,
    voltage float  DEFAULT 0,
    frequency float  DEFAULT 0,
    load_level float  DEFAULT 0,
    timestamp_ TIMESTAMPTZ
);

