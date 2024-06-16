-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS estates (
  id                    UUID             DEFAULT uuid_generate_v4(),
  width               	INTEGER        	NOT NULL,
  length               	INTEGER         NOT NULL,
  created_at            TIMESTAMP        NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP        NOT NULL DEFAULT NOW(),
  deleted_at            TIMESTAMP        DEFAULT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS trees (
  id                    UUID             DEFAULT uuid_generate_v4(),
  estate_id				UUID			NOT NULL REFERENCES estates (id), 
  x       				INTEGER        	NOT NULL,
  y         			INTEGER         NOT NULL,
  height               	INTEGER         NOT NULL,
  created_at            TIMESTAMP        NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP        NOT NULL DEFAULT NOW(),
  deleted_at            TIMESTAMP        DEFAULT NULL,
  PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS index_tree ON trees(estate_id, x, y);