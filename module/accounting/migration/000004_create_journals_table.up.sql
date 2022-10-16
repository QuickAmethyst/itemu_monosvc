CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS journals
(
    id         uuid PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    amount     numeric(18, 8)           NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
)
