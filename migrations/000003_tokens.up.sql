CREATE TABLE IF NOT EXISTS tokens (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    hashed_token VARCHAR(255) NOT NULL UNIQUE,
    user_id      UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    revoked      BOOLEAN      NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);
