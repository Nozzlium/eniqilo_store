ALTER TABLE "orders"
ADD COLUMN product_id uuid NULL,
ADD COLUMN quantity int NULL,
ADD COLUMN price numeric(10,2) NULL;

ALTER TABLE "orders"
  ADD CONSTRAINT orders_product_id_fk FOREIGN KEY ("product_id") REFERENCES "products" ("id");

DROP TABLE IF EXISTS "order_product";
