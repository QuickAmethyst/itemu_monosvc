CREATE TABLE IF NOT EXISTS account_classes
(
    id       SERIAL PRIMARY KEY,
    name     TEXT          NOT NULL UNIQUE,
    type     NUMERIC(1, 0) NOT NULL,
    inactive BOOLEAN DEFAULT FALSE
)
