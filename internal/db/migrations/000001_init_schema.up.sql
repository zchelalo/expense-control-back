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
  "balance" decimal(12,2) NOT NULL DEFAULT 0,
  "user_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "movement_types" (
  "id" uuid PRIMARY KEY,
  "key" varchar(20) UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "movements" (
  "id" uuid PRIMARY KEY,
  "amount" decimal(12,2) NOT NULL,
  "description" text NOT NULL,
  "movement_type_id" uuid NOT NULL,
  "account_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE INDEX idx_users_active
  ON users(email)
  WHERE deleted_at IS NULL;

CREATE INDEX idx_auth_sessions_active_by_user
  ON auth_sessions(user_id, created_at, id)
  WHERE revoked_at IS NULL;

CREATE INDEX idx_auth_sessions_active_by_id
  ON auth_sessions(id)
  WHERE revoked_at IS NULL;

CREATE INDEX idx_accounts_active_by_user_created_id
  ON accounts(user_id, created_at, id)
  WHERE deleted_at IS NULL;

CREATE INDEX idx_movements_active_by_account_created_id
  ON movements(account_id, created_at, id)
  WHERE deleted_at IS NULL;

CREATE INDEX idx_movements_active_by_account_type_created_id
  ON movements(account_id, movement_type_id, created_at, id)
  WHERE deleted_at IS NULL;

ALTER TABLE "auth_sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION;

ALTER TABLE "accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION;

ALTER TABLE "movements" ADD FOREIGN KEY ("movement_type_id") REFERENCES "movement_types" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION;

ALTER TABLE "movements" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION;
