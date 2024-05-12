ALTER TABLE "products"
  DROP COLUMN IF EXISTS "updated_by",
  DROP COLUMN IF EXISTS "deleted_by";
