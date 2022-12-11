CREATE TABLE IF NOT EXISTS bank_account_types
(
    id   int PRIMARY KEY,
    name text
);

INSERT INTO bank_account_types (id, name)
VALUES (1, 'Cash Account'),
       (2, 'Chequing Account'),
       (3, 'Saving Account');

CREATE TABLE IF NOT EXISTS bank_accounts
(
    id          SERIAL PRIMARY KEY,
    account_id  int UNIQUE NOT NULL,
    type_id     int NOT NULL,
    bank_number text,
    inactive    BOOLEAN DEFAULT FALSE,

    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES accounts (id),
    CONSTRAINT fk_type_id FOREIGN KEY (type_id) REFERENCES bank_account_types (id)
);