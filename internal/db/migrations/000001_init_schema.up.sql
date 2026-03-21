CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "email" varchar(254) UNIQUE NOT NULL,
  "password" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "auth_sessions" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "refresh_jti" uuid UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "expires_at" timestamptz NOT NULL,
  "revoked_at" timestamptz,
  "rotated_at" timestamptz,
  "last_seen_at" timestamptz
);

CREATE TABLE "accounts" (
  "id" uuid PRIMARY KEY,
  "name" varchar(50) NOT NULL,
  "balance" decimal(12,2) NOT NULL DEFAULT 0,
  "user_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "movement_types" (
  "id" uuid PRIMARY KEY,
  "name" varchar(25) UNIQUE NOT NULL,
  "key" varchar(20) UNIQUE NOT NULL,
  "description" text,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "categories" (
  "id" uuid PRIMARY KEY,
  "name" varchar(50) UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "movements" (
  "id" uuid PRIMARY KEY,
  "amount" decimal(12,2) NOT NULL,
  "description" text NOT NULL,
  "movement_type_id" uuid NOT NULL,
  "category_id" uuid NOT NULL,
  "account_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

-- Users
CREATE INDEX idx_users_email_active 
  ON users(email) 
  WHERE deleted_at IS NULL;
CREATE INDEX idx_users_id_active 
  ON users(id) 
  WHERE deleted_at IS NULL;

-- Auth Sessions
CREATE INDEX idx_auth_sessions_active_by_id 
  ON auth_sessions(id) 
  WHERE revoked_at IS NULL;
CREATE INDEX idx_auth_sessions_refresh_active 
  ON auth_sessions(id, refresh_jti) 
  WHERE revoked_at IS NULL;
CREATE INDEX idx_auth_sessions_active_by_user 
  ON auth_sessions(user_id, created_at DESC, id DESC) 
  WHERE revoked_at IS NULL;

-- Accounts
CREATE INDEX idx_accounts_user_cursor 
  ON accounts(user_id, created_at DESC, id DESC) 
  WHERE deleted_at IS NULL;

-- Movement Types
CREATE INDEX idx_movement_types_active 
  ON movement_types(id) 
  WHERE deleted_at IS NULL;
CREATE INDEX idx_movement_types_key_active 
  ON movement_types(key) 
  WHERE deleted_at IS NULL;
CREATE INDEX idx_movement_types_name_trgm 
  ON movement_types USING GIN (name gin_trgm_ops) 
  WHERE deleted_at IS NULL;

-- Categories
CREATE INDEX idx_categories_cursor 
  ON categories(created_at DESC, id DESC) 
  WHERE deleted_at IS NULL;
CREATE INDEX idx_categories_name_trgm 
  ON categories USING GIN (name gin_trgm_ops) 
  WHERE deleted_at IS NULL;

-- Movements
CREATE INDEX idx_movements_user_cursor 
  ON movements(user_id, created_at DESC, id DESC) 
  WHERE deleted_at IS NULL;
CREATE INDEX idx_movements_user_category 
  ON movements(user_id, category_id, created_at DESC, id DESC) 
  WHERE deleted_at IS NULL;
CREATE INDEX idx_movements_user_type 
  ON movements(user_id, movement_type_id, created_at DESC, id DESC) 
  WHERE deleted_at IS NULL;
CREATE INDEX idx_movements_user_amount 
  ON movements(user_id, amount, created_at DESC, id DESC) 
  WHERE deleted_at IS NULL;
CREATE INDEX idx_movements_account_cursor 
  ON movements(account_id, created_at DESC, id DESC) 
  WHERE deleted_at IS NULL;
CREATE INDEX idx_movements_id_active 
  ON movements(id) 
  WHERE deleted_at IS NULL;

ALTER TABLE "auth_sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "movements" ADD FOREIGN KEY ("movement_type_id") REFERENCES "movement_types" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "movements" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "movements" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "movements" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION DEFERRABLE INITIALLY IMMEDIATE;
