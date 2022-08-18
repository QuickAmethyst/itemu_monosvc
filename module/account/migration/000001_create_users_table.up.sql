CREATE TABLE IF NOT EXISTS users
(
    id         uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    email      text NOT NULL,
    password   text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE
)
