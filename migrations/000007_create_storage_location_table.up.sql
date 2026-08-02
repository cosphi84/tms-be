-- Sesuaikan nomor urut file ini sama konvensi migration tool kamu.

CREATE TABLE IF NOT EXISTS storage_locations (
     id        BIGSERIAL PRIMARY KEY,
     code      VARCHAR(20)  NOT NULL,
     name      VARCHAR(100) NOT NULL,
     office_id BIGINT       NOT NULL REFERENCES offices(id) ON UPDATE CASCADE ON DELETE RESTRICT,
     member    VARCHAR(10)  NOT NULL CHECK (member IN ('SEID', 'OTS', 'SASS')),

     created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
     created_by BIGINT,
     updated_at TIMESTAMPTZ,
     updated_by BIGINT,
     deleted_at TIMESTAMPTZ,
     deleted_by BIGINT
);

CREATE INDEX IF NOT EXISTS idx_storage_locations_office_id ON storage_locations (office_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_storage_locations_code_active
    ON storage_locations (code)
    WHERE deleted_at IS NULL;