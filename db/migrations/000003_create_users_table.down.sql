ALTER TABLE "accounts" DROP FOREIGN KEY "owner";

ALTER TABLE "accounts" DROP CONSTRAINT "onwer_balance_key";

DROP TABLE "users";