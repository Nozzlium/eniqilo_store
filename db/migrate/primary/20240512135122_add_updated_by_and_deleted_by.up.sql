ALTER TABLE "products"
  ADD COLUMN IF NOT EXISTS "updated_by" uuid REFERENCES "users" ("id") ON DELETE SET NULL,
  ADD COLUMN IF NOT EXISTS "deleted_by" uuid REFERENCES "users" ("id") ON DELETE SET NULL;
