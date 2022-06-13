-- +migrate Up
ALTER TABLE "images" DROP COLUMN "name";
ALTER TABLE "images" ALTER COLUMN deleted TYPE INTEGER USING deleted::INTEGER;
