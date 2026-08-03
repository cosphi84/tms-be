CREATE TABLE IF NOT EXISTS storage_locations (
     id        BIGSERIAL PRIMARY KEY,
     code      VARCHAR(20)  NOT NULL,
     name      VARCHAR(100) NOT NULL,
     office_id BIGINT       NOT NULL REFERENCES offices(id) ON UPDATE CASCADE ON DELETE RESTRICT,
     member    VARCHAR(10)  NOT NULL CHECK (member IN ('SEID', 'OTS', 'SASS')),
     group_id   BIGINT NOT NULL  REFERENCES tool_groups(id) ON UPDATE CASCADE ON DELETE SET NULL,

     created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
     created_by BIGINT,
     updated_at TIMESTAMPTZ,
     updated_by BIGINT,
     deleted_at TIMESTAMPTZ,
     deleted_by BIGINT
);

CREATE INDEX IF NOT EXISTS idx_storage_locations_office_id ON storage_locations (office_id);
CREATE INDEX IF NOT EXISTS idx_storage_locations_member ON storage_locations (member);
CREATE INDEX IF NOT EXISTS idx_storage_locations_group_id ON storage_locations (group_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_storage_locations_code_active
    ON storage_locations (code)
    WHERE deleted_at IS NULL;