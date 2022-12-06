CREATE TABLE IF NOT EXISTS bank_accounts
(
    id          SERIAL PRIMARY KEY,
    account_id  int           NOT NULL,
    type        NUMERIC(1, 0) NOT NULL,
    bank_number text,
    inactive    BOOLEAN DEFAULT FALSE,

    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES accounts (id)
);
