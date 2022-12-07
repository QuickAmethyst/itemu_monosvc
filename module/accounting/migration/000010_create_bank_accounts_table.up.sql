CREATE TYPE bank_account_type AS ENUM (1, 2, 3);
CREATE TABLE IF NOT EXISTS bank_accounts
(
    id          SERIAL PRIMARY KEY,
    account_id  int           NOT NULL,
    type        bank_account_type,
    bank_number text,
    inactive    BOOLEAN DEFAULT FALSE,

    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES accounts (id)
);
