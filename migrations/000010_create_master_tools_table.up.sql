-- Sesuaikan nomor urut file ini. WAJIB jalan SETELAH tool_categories ada.

CREATE TABLE IF NOT EXISTS master_tools (
    id           BIGSERIAL PRIMARY KEY,
    code         VARCHAR(20)  NOT NULL,
    name         VARCHAR(150) NOT NULL,
    brand         VARCHAR(100) NOT NULL,
    type         VARCHAR(100) NOT NULL,
    serial_num   VARCHAR(100),
    category_id  BIGINT       NOT NULL REFERENCES tool_categories(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    group_id   BIGINT NOT NULL  REFERENCES tool_groups(id) ON UPDATE CASCADE ON DELETE SET NULL,
    photo_id   BIGINT REFERENCES asset_files(id) ON UPDATE CASCADE ON DELETE SET NULL,
    price        NUMERIC(15, 2) NOT NULL CHECK (price >= 0),
    usage_period VARCHAR(10)  NOT NULL, -- format: "2y", "6m", "30d"
    is_active    BOOLEAN      NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by BIGINT,
    updated_at TIMESTAMPTZ,
    updated_by BIGINT,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT
);

CREATE INDEX IF NOT EXISTS idx_master_tools_category_id ON master_tools (category_id);
CREATE INDEX IF NOT EXISTS idx_master_tools_is_active ON master_tools (is_active);

-- Partial unique -- code unik di antara row yang masih aktif (belum di-soft-delete)
CREATE UNIQUE INDEX IF NOT EXISTS idx_master_tools_code_active
    ON master_tools (code, name)
    WHERE deleted_at IS NULL;