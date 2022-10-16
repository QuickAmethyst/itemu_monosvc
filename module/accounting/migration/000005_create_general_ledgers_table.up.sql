CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS general_ledgers
(
    id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    journal_id uuid           NOT NULL,
    account_id int            NOT NULL,
    amount     numeric(18, 8) NOT NULL,
    created_by uuid           NOT NULL,

    UNIQUE (journal_id, account_id),
    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES accounts (id)
)
