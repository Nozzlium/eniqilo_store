ALTER TABLE "orders"
  ADD COLUMN "product_id" uuid not null,
  ADD COLUMN "quantity" int not null,
  ADD COLUMN "price" numeric(10,2) NOT NULL;

ALTER TABLE "orders"
  ADD CONSTRAINT fk_order_product_product FOREIGN KEY ("product_id") REFERENCES "products" ("id");

DROP TABLE IF EXISTS "order_product";
