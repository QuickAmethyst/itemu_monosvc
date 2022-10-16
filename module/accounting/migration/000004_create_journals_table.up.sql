CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS journals
(
    id         uuid PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    account_id int                      NOT NULL,
    memo       text,
    amount     numeric(18, 8),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES accounts (id)
)
