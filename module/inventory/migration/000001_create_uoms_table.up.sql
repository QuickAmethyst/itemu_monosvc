CREATE TABLE IF NOT EXISTS uoms
(
    id          SERIAL PRIMARY KEY,
    name        TEXT UNIQUE NOT NULL,
    description TEXT,
    decimal     INTEGER     NOT NULL,
    created_at  timestamp with time zone DEFAULT now(),
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone
);
