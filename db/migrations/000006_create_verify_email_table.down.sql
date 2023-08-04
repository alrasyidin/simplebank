DROP TABLE verify_emails CASCADE;

ALTER TABLE users DROP CONSTRAINT fk_username_user_verify_email;