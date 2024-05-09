-- drop foreign key constraints
ALTER TABLE "orders" DROP CONSTRAINT "orders_product_id_fkey";

-- drop unnecessary columns
ALTER TABLE "orders"
  DROP COLUMN "product_id",
  DROP COLUMN "quantity",
  DROP COLUMN "price";

CREATE TABLE IF NOT EXISTS "order_product" (
  "order_id" uuid NOT NULL,
  "product_id" uuid NOT NULL,
  "quantity" int NOT NULL,
  "price" numeric(10, 2) NOT NULL,
  "total_price" numeric(10, 2) NOT NULL GENERATED ALWAYS AS ("quantity" * "price") STORED,
  PRIMARY KEY ("order_id", "product_id"),
  FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE
);
