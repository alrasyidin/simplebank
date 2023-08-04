DROP TABLE "users" CASCADE;

ALTER TABLE "accounts" DROP FOREIGN KEY "owner";

ALTER TABLE "accounts" DROP CONSTRAINT "onwer_currency_key";
