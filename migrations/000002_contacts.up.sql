CREATE TABLE IF NOT EXISTS contacts (
    user_id    UUID        NOT NULL REFERENCES users(id),
    contact_id UUID        NOT NULL REFERENCES users(id),
    active     BOOLEAN     NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, contact_id),
    CONSTRAINT chk_no_self_contact CHECK (user_id <> contact_id)
);
