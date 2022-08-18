CREATE TABLE IF NOT EXISTS uoms
(
    id          SERIAL PRIMARY KEY,
    name        TEXT UNIQUE NOT NULL,
    description TEXT,
    decimal     INTEGER     NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE,
    deleted_at  TIMESTAMP WITH TIME ZONE
);
