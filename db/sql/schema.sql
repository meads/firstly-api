CREATE TABLE images (
  "id"   BIGSERIAL  PRIMARY KEY,
  "data" TEXT       NOT NULL,
  "created" VARCHAR NOT NULL,
  "deleted" BOOLEAN DEFAULT FALSE
);
