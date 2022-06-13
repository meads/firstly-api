-- +migrate Up
ALTER TABLE "images" DROP COLUMN "name";
ALTER TABLE "images" ALTER COLUMN "deleted" TYPE INT DEFAULT 0;
