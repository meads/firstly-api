-- +migrate Up
ALTER TABLE "images" DROP COLUMN "name";
ALTER TABLE "images" DROP COLUMN deleted;
ALTER TABLE "images" ADD  COLUMN deleted INT;
