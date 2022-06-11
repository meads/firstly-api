CREATE TABLE images (
  "id"   BIGSERIAL PRIMARY KEY,
  "name" TEXT NULL,
  "data" TEXT      NOT NULL,
  "created" NVARCHAR(8) NOT NULL,
  "deleted" BOOLEAN DEFAULT FALSE
);
