CREATE TABLE IF NOT EXISTS tool_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by BIGINT,
    updated_at TIMESTAMPTZ,
    updated_by BIGINT,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tool_categories_name ON tool_categories(name);