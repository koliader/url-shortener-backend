CREATE TABLE "users" (
  "username" varchar UNIQUE NOT NULL,
  "password" varchar ,
  "color" varchar NOT NULL
);

CREATE TABLE "urls" (
  "url" varchar NOT NULL,
  "code" varchar UNIQUE NOT NULL,
  "owner" varchar,
  "clicks" integer NOT NULL DEFAULT (0)
);

ALTER TABLE "urls" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
