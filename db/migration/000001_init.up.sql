CREATE TABLE "users" (
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL
);

CREATE TABLE "urls" (
  "url" varchar NOT NULL,
  "code" varchar UNIQUE NOT NULL,
  "owner" varchar NOT NULL
);

ALTER TABLE "urls" ADD FOREIGN KEY ("owner") REFERENCES "users" ("email");
