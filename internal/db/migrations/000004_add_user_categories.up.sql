CREATE TABLE "user_categories" (
  "user_id" uuid NOT NULL,
  "category_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz,
  PRIMARY KEY ("user_id", "category_id")
);

CREATE INDEX idx_user_categories_user_cursor
  ON user_categories(user_id, created_at DESC, category_id DESC)
  WHERE deleted_at IS NULL;

ALTER TABLE "user_categories"
  ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id")
  ON DELETE RESTRICT ON UPDATE NO ACTION DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "user_categories"
  ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id")
  ON DELETE RESTRICT ON UPDATE NO ACTION DEFERRABLE INITIALLY IMMEDIATE;

INSERT INTO user_categories (
  user_id,
  category_id,
  created_at,
  updated_at
)
SELECT
  m.user_id,
  m.category_id,
  MIN(m.created_at) AS created_at,
  MAX(m.updated_at) AS updated_at
FROM movements m
INNER JOIN categories c
  ON c.id = m.category_id
 AND c.deleted_at IS NULL
GROUP BY m.user_id, m.category_id
ON CONFLICT (user_id, category_id) DO UPDATE
SET
  updated_at = EXCLUDED.updated_at,
  deleted_at = NULL
;
