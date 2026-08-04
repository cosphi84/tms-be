CREATE TABLE IF NOT EXISTS stocks (
    id BIGSERIAL PRIMARY KEY,
    tools_id BIGINT NOT NULL REFERENCES master_tools(id) ON DELETE CASCADE,
    slock_id BIGINT NOT NULL REFERENCES storage_locations(id) ON DELETE CASCADE,
    qty INT NOT NULL CHECK (qty >= 0),
    start_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    expired_date TIMESTAMPTZ NULL,
    cont_stock INT NOT NULL DEFAULT 1 CHECK (cont_stock >= 1),
    remarks TEXT NULL,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by BIGINT,
    updated_at TIMESTAMPTZ,
    updated_by BIGINT,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT    
);

CREATE INDEX IF NOT EXISTS idx_stocks_tools_id ON stocks(tools_id);
CREATE INDEX IF NOT EXISTS idx_stocks_slock_id ON stocks(slock_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_stocks_tools_id_slock_id ON stocks(tools_id, slock_id);