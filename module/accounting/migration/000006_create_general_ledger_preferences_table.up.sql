CREATE TABLE IF NOT EXISTS general_ledger_preferences
(
    id         SERIAL PRIMARY KEY,
    account_id int,

    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES accounts (id)
);

INSERT INTO general_ledger_preferences (id)
VALUES (1),
       (2);
