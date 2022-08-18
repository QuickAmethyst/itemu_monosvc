CREATE TABLE IF NOT EXISTS user_profiles
(
    id         uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    user_id    uuid NOT NULL UNIQUE,
    full_name  text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id)
)
