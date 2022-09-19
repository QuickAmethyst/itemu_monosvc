CREATE TABLE IF NOT EXISTS account_classes
(
    id       SERIAL PRIMARY KEY,
    name     TEXT          NOT NULL UNIQUE,
    type_id  NUMERIC(1, 0) NOT NULL,
    inactive BOOLEAN DEFAULT FALSE
)
