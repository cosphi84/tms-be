CREATE TABLE IF NOT EXISTS office_leaders (
          id         BIGSERIAL PRIMARY KEY,
          office_id  BIGINT NOT NULL REFERENCES offices(id) ON UPDATE CASCADE ON DELETE RESTRICT,
          user_id    BIGINT NOT NULL REFERENCES users(id)   ON UPDATE CASCADE ON DELETE RESTRICT,
          start_date DATE   NOT NULL,
          end_date   DATE,                      -- NULL = masih aktif/menjabat

          created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
          created_by BIGINT,
          updated_at TIMESTAMPTZ,
          updated_by BIGINT,
          deleted_at TIMESTAMPTZ,
          deleted_by BIGINT,

          CONSTRAINT chk_office_leaders_date_range CHECK (end_date IS NULL OR end_date >= start_date)
);

CREATE INDEX IF NOT EXISTS idx_office_leaders_office_id ON office_leaders (office_id);
CREATE INDEX IF NOT EXISTS idx_office_leaders_user_id   ON office_leaders (user_id);

-- ATURAN #1: cuma boleh ada 1 leader AKTIF (end_date IS NULL) per office.
-- Ini yang jamin "satu office hanya punya 1 leader aktif" di level DB,
-- gak cuma di kode Go -- kalau 2 request assign leader ke office yang sama
-- kebetulan lolos row-lock (harusnya gak mungkin, tapi ini backstop terakhir),
-- index ini yang bakal nolak insert kedua dengan error unique violation.
CREATE UNIQUE INDEX IF NOT EXISTS idx_office_leaders_active_per_office
    ON office_leaders (office_id)
    WHERE end_date IS NULL AND deleted_at IS N  ULL;

-- ATURAN #2: cuma boleh ada 1 assignment AKTIF per user (satu user gak bisa
-- jadi leader di 2 office bersamaan).
CREATE UNIQUE INDEX IF NOT EXISTS idx_office_leaders_active_per_user
    ON office_leaders (user_id)
    WHERE end_date IS NULL AND deleted_at IS NULL;