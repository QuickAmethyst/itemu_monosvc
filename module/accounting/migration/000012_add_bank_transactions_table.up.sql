CREATE TABLE IF NOT EXISTS bank_transactions
(
    id              SERIAL PRIMARY KEY,
    journal_id      uuid                     NOT NULL,
    bank_account_id int                      NOT NULL,
    amount          numeric(18, 8)           NOT NULL DEFAULT 0,
    balance         numeric(18, 8)           NOT NULL DEFAULT 0,
    memo            text,
    created_by      uuid                     NOT NULL,
    trans_date      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_journal_id FOREIGN KEY (journal_id) REFERENCES journals (id),
    CONSTRAINT fk_bank_account_id FOREIGN KEY (bank_account_id) REFERENCES bank_accounts (id)
);
