CREATE TABLE IF NOT EXISTS account_groups
(
    id        SERIAL PRIMARY KEY,
    parent_id int  NOT NULL,
    class_id  int  NOT NULL,
    name      TEXT NOT NULL UNIQUE,
    inactive  BOOLEAN DEFAULT FALSE,

    CONSTRAINT parent FOREIGN KEY (parent_id) REFERENCES account_groups (id),
    CONSTRAINT fk_class_id FOREIGN KEY (class_id) REFERENCES account_classes (id)
)
