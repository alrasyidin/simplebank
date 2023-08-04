CREATE TABLE verify_emails(
  "id" BIGSERIAL PRIMARY KEY,
  "username" VARCHAR NOT NULL,
  "email" VARCHAR NOT NULL,
  "secret_code" varchar NOT NULL,
  "is_used" BOOL NOT NULL DEFAULT FALSE, 
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "expired_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);

ALTER TABLE users ADD CONSTRAINT fk_username_user_verify_email FOREIGN KEY("username") REFERENCES "users"("username");