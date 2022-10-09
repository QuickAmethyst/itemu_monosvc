CREATE TABLE IF NOT EXISTS accounts
(
    id       SERIAL PRIMARY KEY,
    name     TEXT UNIQUE NOT NULL,
    group_id int NOT NULL,
    inactive BOOLEAN DEFAULT FALSE,

    CONSTRAINT fk_group_id FOREIGN KEY (group_id) REFERENCES account_groups (id)
)
