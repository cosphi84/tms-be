CREATE TABLE IF NOT EXISTS stocks (
    id BIGSERIAL PRIMARY KEY,
    tools_id BIGINT NOT NULL REFERENCES master_tools(id) ON DELETE CASCADE,
    slock_ic BIGINT NOT NULL REFERENCES storage_locations(id) ON DELETE CASCADE,
    qty INT NOT NULL CHECK (qty >= 0),
    start_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    expired_date TIMESTAMPTZ NULL,
    remarks TEXT NULL,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by BIGINT,
    updated_at TIMESTAMPTZ,
    updated_by BIGINT,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT    
);

CREATE INDEX IF NOT EXISTS idx_stocks_tools_id ON stocks(tools_id);
CREATE INDEX IF NOT EXISTS idx_stocks_slock_ic ON stocks(slock_ic);
CREATE UNIQUE INDEX IF NOT EXISTS idx_stocks_tools_id_slock_ic ON stocks(tools_id, slock_ic);