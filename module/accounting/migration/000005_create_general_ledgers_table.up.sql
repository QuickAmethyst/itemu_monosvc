CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS general_ledgers
(
    id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    journal_id uuid           NOT NULL,
    account_id int            NOT NULL,
    amount     numeric(18, 8) NOT NULL,
    created_by uuid           NOT NULL,

    CONSTRAINT fk_journal_id FOREIGN KEY (journal_id) REFERENCES journals (id),
    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES accounts (id)
);

CREATE INDEX idx_journal_id ON general_ledgers (journal_id);
CREATE INDEX idx_account_id ON general_ledgers (account_id);
