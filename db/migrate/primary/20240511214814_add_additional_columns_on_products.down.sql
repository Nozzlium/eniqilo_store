
ALTER TABLE "products" 
  DROP CONSTRAINT "fk_products_updated_by",
  DROP COLUMN IF EXISTS "updated_by",
  DROP CONSTRAINT "fk_products_deleted_by",
  DROP COLUMN IF EXISTS "deleted_by";

