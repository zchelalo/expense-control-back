ALTER TABLE categories
  ADD COLUMN is_system boolean NOT NULL DEFAULT FALSE,
  ADD COLUMN system_key varchar(50)
;

ALTER TABLE categories
  DROP CONSTRAINT categories_name_key
;

ALTER TABLE categories
  ADD CONSTRAINT chk_categories_system_fields
  CHECK (
    (is_system = TRUE AND system_key IS NOT NULL)
    OR (is_system = FALSE AND system_key IS NULL)
  )
;

CREATE UNIQUE INDEX idx_categories_name_unique_user
  ON categories(name)
  WHERE is_system = FALSE
;

CREATE UNIQUE INDEX idx_categories_system_key_unique
  ON categories(system_key)
  WHERE system_key IS NOT NULL
;
