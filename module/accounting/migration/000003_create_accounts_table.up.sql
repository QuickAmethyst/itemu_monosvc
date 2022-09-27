CREATE TABLE IF NOT EXISTS accounts
(
    id       SERIAL PRIMARY KEY,
    name     TEXT UNIQUE,
    group_id int,
    inactive BOOLEAN DEFAULT FALSE,

    CONSTRAINT fk_group_id FOREIGN KEY (group_id) REFERENCES account_groups (id)
)
