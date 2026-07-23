CREATE TABLE IF NOT EXISTS office_leaders (
    id BIGSERIAL PRIMARY KEY,
    office_id BIGINT NOT NULL REFERENCES offices(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date DATE NULL,
    end_date DATE NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_office_leaders_office_id ON office_leaders(office_id);
CREATE INDEX IF NOT EXISTS idx_office_leaders_user_id ON office_leaders(user_id);