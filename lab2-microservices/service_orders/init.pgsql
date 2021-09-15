CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE SCHEMA IF NOT EXISTS orders;
DROP TABLE IF EXISTS orders.orders cascade;
 
CREATE TABLE IF NOT EXISTS orders.orders (
  id         SERIAL CONSTRAINT orders_pkey PRIMARY KEY,
  item_uid   UUID         NOT NULL,
  order_date TIMESTAMP    NOT NULL,
  order_uid  UUID         NOT NULL CONSTRAINT idx_orders_order_uid UNIQUE,
  status     VARCHAR(255) NOT NULL,
  user_uid   UUID         NOT NULL
);
