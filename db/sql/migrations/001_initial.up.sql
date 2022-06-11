
-- +migrate Up
CREATE TABLE "images" (
  "id"   BIGSERIAL PRIMARY KEY,
  "name" TEXT NULL,
  "data" TEXT      NOT NULL,
  "created" VARCHAR(32) NOT NULL,
  "deleted" BOOLEAN DEFAULT FALSE
);
