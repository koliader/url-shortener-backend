CREATE TABLE "users" (
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL
);

CREATE TABLE "urls" (
  "url" varchar NOT NULL,
  "code" varchar UNIQUE NOT NULL,
  "owner" varchar
);

ALTER TABLE "urls" ADD FOREIGN KEY ("owner") REFERENCES "users" ("email");
