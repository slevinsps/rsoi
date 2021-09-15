CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE SCHEMA IF NOT EXISTS session;
DROP TABLE IF EXISTS session.users cascade;
 
CREATE TABLE IF NOT EXISTS session.users
(
    id       SERIAL CONSTRAINT users_pkey PRIMARY KEY,
    login     VARCHAR(255) NOT NULL CONSTRAINT idx_user_name UNIQUE,
    password     VARCHAR(255) NOT NULL,
    user_uuid UUID         NOT NULL CONSTRAINT idx_user_user_uuid UNIQUE,
    is_admin boolean NOT NULL
);



CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE SCHEMA IF NOT EXISTS equipment;
DROP TABLE IF EXISTS equipment.equipment cascade;
DROP TABLE IF EXISTS equipment.equipment_model cascade;
CREATE TABLE IF NOT EXISTS equipment.equipment_model
(
    id       SERIAL CONSTRAINT equipment_models_pkey PRIMARY KEY,
    equipment_model_uuid UUID  NOT NULL CONSTRAINT idx_equipment_model_uuid UNIQUE,
    name     VARCHAR(255) NOT NULL
);

INSERT INTO equipment.equipment_model(equipment_model_uuid, name) VALUES ('ec3c0ea5-94eb-4253-b02f-83e465640440', 'IntelI7');
INSERT INTO equipment.equipment_model(equipment_model_uuid, name) VALUES ('93af6778-f1e5-4d1f-84c1-68bcb85904b3', 'AMD4');
INSERT INTO equipment.equipment_model(equipment_model_uuid, name) VALUES ('2c95038c-bb76-4a7d-abfe-bf8db50448ff', 'IntelI3');


CREATE TABLE IF NOT EXISTS equipment.equipment
(
    id       SERIAL CONSTRAINT equipments_pkey PRIMARY KEY,
    name     VARCHAR(255) NOT NULL,
    equipment_uuid UUID         NOT NULL CONSTRAINT idx_equipment_uuid UNIQUE,
    equipment_model_uuid  UUID  NOT NULL   REFERENCES equipment.equipment_model(equipment_model_uuid) ON DELETE CASCADE,
    status VARCHAR(255)
);

INSERT INTO equipment.equipment(name, equipment_uuid, equipment_model_uuid, status) VALUES ('intel_r3f5', '26065f2f-a0d1-403c-9a9a-c35d11fc8535', 'ec3c0ea5-94eb-4253-b02f-83e465640440', 'ACTIVE');
INSERT INTO equipment.equipment(name, equipment_uuid, equipment_model_uuid, status) VALUES ('intel_v5f8v4', '99d81285-142c-4578-bfbe-e39f7aeb2c9f', 'ec3c0ea5-94eb-4253-b02f-83e465640440', 'ACTIVE');

INSERT INTO equipment.equipment(name, equipment_uuid, equipment_model_uuid, status) VALUES ('amd_54', '76e4d634-e6f0-4558-b82c-7cd4948a36db', '93af6778-f1e5-4d1f-84c1-68bcb85904b3', 'ACTIVE');
INSERT INTO equipment.equipment(name, equipment_uuid, equipment_model_uuid, status) VALUES ('amd_78', '33ad32f1-0ee5-4271-925a-032013c362d3', '93af6778-f1e5-4d1f-84c1-68bcb85904b3', 'ACTIVE');

INSERT INTO equipment.equipment(name, equipment_uuid, equipment_model_uuid, status) VALUES ('intel_7f7f7f', 'd70260db-1cd1-4e93-a3da-ce1615cdef70', '2c95038c-bb76-4a7d-abfe-bf8db50448ff', 'ACTIVE');
INSERT INTO equipment.equipment(name, equipment_uuid, equipment_model_uuid, status) VALUES ('intel_65f65fg', '8bebea91-6fc8-47d0-9c43-1e4789047d57', '2c95038c-bb76-4a7d-abfe-bf8db50448ff', 'ACTIVE');



CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE SCHEMA IF NOT EXISTS monitor;
DROP TABLE IF EXISTS monitor.monitor cascade;
DROP TABLE IF EXISTS monitor.monitor_equipment_xref cascade;

CREATE TABLE IF NOT EXISTS monitor.monitor
(
    id       SERIAL CONSTRAINT users_pkey PRIMARY KEY,
    name     VARCHAR(255) NOT NULL,
    monitor_uuid UUID         NOT NULL CONSTRAINT idx_user_user_uuid UNIQUE,
    user_uuid    UUID       NOT NULL   REFERENCES session.users(user_uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS monitor.monitor_equipment_xref
(
    id       SERIAL CONSTRAINT equipment_models_pkey PRIMARY KEY,
    monitor_uuid UUID  NOT NULL   REFERENCES monitor.monitor(monitor_uuid) ON DELETE CASCADE,
    equipment_uuid UUID  NOT NULL   REFERENCES equipment.equipment(equipment_uuid) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_monitor_uuid_equipment_uuid on monitor.monitor_equipment_xref (monitor_uuid, equipment_uuid);






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



CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE SCHEMA IF NOT EXISTS gateway;
DROP TABLE IF EXISTS gateway.services_secrets cascade;
 
CREATE TABLE IF NOT EXISTS gateway.services_secrets
(
    id       SERIAL CONSTRAINT gateway_pkey PRIMARY KEY,
    login     VARCHAR(255) NOT NULL CONSTRAINT idx_login UNIQUE,
    password     VARCHAR(255) NOT NULL
);






CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE SCHEMA IF NOT EXISTS documentation;
DROP TABLE IF EXISTS documentation.documentation cascade;


CREATE TABLE IF NOT EXISTS documentation.documentation
(
    id       SERIAL CONSTRAINT documentation_pkey PRIMARY KEY,
    documentation_uuid UUID  NOT NULL CONSTRAINT idx_documentation_uuid UNIQUE,
    equipment_model_uuid UUID  NOT NULL   REFERENCES equipment.equipment_model(equipment_model_uuid) ON DELETE CASCADE,
    name     VARCHAR(255) NOT NULL,
    path    VARCHAR(255) NOT NULL
);

INSERT INTO documentation.documentation(documentation_uuid, equipment_model_uuid, name, path) VALUES ('9c95438c-bb76-4b7d-abfe-bf8db50448ff', '2c95038c-bb76-4a7d-abfe-bf8db50448ff', 'readme.pdf', './uploads/readme.pdf');


-- INSERT INTO session.users(login, password, user_uuid) VALUES ('admin', 'ivanivan', '6d2cb5a0-943c-4b96-9aa6-89eac7bdfd2b');
