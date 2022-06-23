
-- +migrate Up
ALTER TABLE "images" ALTER COLUMN "created" TYPE VARCHAR(32) NOT NULL;

