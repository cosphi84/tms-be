CREATE TABLE IF NOT EXISTS offices (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT NULL REFERENCES offices(id)
                     ON DELETE RESTRICT,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL
        CHECK (
            type IN (
                'cabang',
                'sdss',
                'ssr',
                'sass',
                'tc',
                'hq'
                )
            ),

    lft INTEGER NOT NULL,
    rgt INTEGER NOT NULL,
    depth INTEGER NOT NULL DEFAULT 0,
    children_count INT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by INTEGER NULL,
    updated_at TIMESTAMPTZ NULL,
    updated_by BIGINT NULL,
    deleted_at TIMESTAMPTZ NULL,
    deleted_by INTEGER NULL
);

CREATE INDEX IF NOT EXISTS idx_offices_code ON offices(code);
CREATE INDEX IF NOT EXISTS idx_offices_parent_id ON offices(parent_id);