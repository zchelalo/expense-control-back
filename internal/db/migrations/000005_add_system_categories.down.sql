UPDATE categories
SET name = 'sys_' || REPLACE(id::text, '-', '')
WHERE is_system = TRUE
;

DROP INDEX IF EXISTS idx_categories_system_key_unique;
DROP INDEX IF EXISTS idx_categories_name_unique_user;

ALTER TABLE categories
  DROP CONSTRAINT IF EXISTS chk_categories_system_fields
;

ALTER TABLE categories
  DROP COLUMN IF EXISTS system_key,
  DROP COLUMN IF EXISTS is_system
;

ALTER TABLE categories
  ADD CONSTRAINT categories_name_key UNIQUE (name)
;
