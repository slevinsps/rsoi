CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
DROP TABLE IF EXISTS Post cascade;
DROP TABLE IF EXISTS Users cascade;
DROP TABLE IF EXISTS Forum cascade;
DROP TABLE IF EXISTS Thread cascade;
DROP TABLE IF EXISTS Vote cascade;

CREATE TABLE IF NOT EXISTS Persons (
  id       bigserial    PRIMARY KEY,
  name     varchar(100) NOT NULL,
  age      int          NOT NULL,
  address  CITEXT       NOT NULL,
  work     CITEXT
);
