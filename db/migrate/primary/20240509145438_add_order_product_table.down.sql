ALTER TABLE "orders"
  ADD COLUMN "product_id",
  ADD COLUMN "quantity",
  ADD COLUMN "price";

ALTER TABLE "orders"
  ADD CONSTRAINT FOREIGN KEY ("product_id") REFERENCES "products" ("id");

DROP TABLE IF EXISTS "order_product";
