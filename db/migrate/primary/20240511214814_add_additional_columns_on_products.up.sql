ALTER TABLE "products" 
  ADD COLUMN IF NOT EXISTS "updated_by" uuid NULL,
  ADD CONSTRAINT "fk_products_updated_by" FOREIGN KEY ("updated_by") REFERENCES "users"("id"),
  ADD COLUMN IF NOT EXISTS "deleted_by" uuid NULL,
  ADD CONSTRAINT "fk_products_deleted_by" FOREIGN KEY ("deleted_by") REFERENCES "users"("id");

