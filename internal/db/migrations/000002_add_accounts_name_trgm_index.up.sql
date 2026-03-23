CREATE INDEX idx_accounts_name_trgm 
  ON accounts USING GIN (name gin_trgm_ops) 
  WHERE deleted_at IS NULL;
