CREATE TABLE "image" (
  "id"   BIGSERIAL  PRIMARY KEY,
  "data" TEXT       NOT NULL,
  "created" VARCHAR NOT NULL,
  "deleted" BOOLEAN NOT NULL DEFAULT FALSE
);
