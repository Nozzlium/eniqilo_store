
CREATE TABLE IF NOT EXISTS "users" (
  "id" uuid NOT NULL,
  "name" varchar(50),
  -- "username" varchar(255) NOT NULL,
  "password" varchar(100) NOT NULL,
  -- "email" varchar(255) NULL DEFAULT NULL,
  "phone_number" varchar(20) NOT NULL,
  -- "email_verified_at" timestamp NULL DEFAULT NULL,
  "phone_verified_at" timestamp NULL DEFAULT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp NULL DEFAULT NULL,
  PRIMARY KEY ("id"),
  UNIQUE ("phone_number")
);

-- CREATE TABLE IF NOT EXISTS "password_resets" (
--   "id" bigserial NOT NULL,
--   "email" varchar(255) NOT NULL,
--   "token" varchar(255) NOT NULL,
--   "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
--   PRIMARY KEY ("id")
-- );
-- 
-- CREATE TABLE IF NOT EXISTS "access_tokens" (
--   "id" bigserial NOT NULL,
--   "tokenable_type" varchar(255) NOT NULL,
--   "tokenable_id" uuid NOT NULL,
--   "name" varchar(255) NOT NULL,
--   "token" varchar(64) NOT NULL,
--   "abilities" text DEFAULT NULL,
--   "last_used_at" timestamp NULL DEFAULT NULL,
--   "expires_at" timestamp NULL DEFAULT NULL,
--   "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
--   "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
--   PRIMARY KEY ("id"),
--   UNIQUE ("token")
-- );
-- 
-- CREATE INDEX "access_tokens_tokenable_type_tokenable_id_index" ON "access_tokens" ("tokenable_type","tokenable_id");

CREATE TYPE "category" AS ENUM ('clothing', 'accessories', 'footwear', 'beverages');

CREATE TABLE IF NOT EXISTS "products" (
  id uuid NOT NULL,
  name varchar(30) NOT NULL,
  sku varchar(30) NOT NULL,
  stock int NOT NULL,
  price numeric(10,2) NOT NULL,
  category CATEGORY NOT NULL,
  notes text NOT NULL,
  location varchar(200) NOT NULL,
  is_available boolean NOT NULL DEFAULT TRUE,
  image_url varchar(255) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp NULL DEFAULT NULL,
  created_by uuid NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("created_by") REFERENCES "users" ("id") ON DELETE CASCADE,
  UNIQUE ("sku")
);

CREATE TABLE IF NOT EXISTS "customers" (
  id uuid NOT NULL,
  name varchar(50) NOT NULL,
  -- email varchar(255) NULL DEFAULT NULL,
  phone_number varchar(20) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY ("id"),
  UNIQUE ("phone_number")
);

CREATE TABLE IF NOT EXISTS "orders" (
  id uuid NOT NULL,
  customer_id uuid NOT NULL,
  product_id uuid NOT NULL,
  quantity int NOT NULL,
  price numeric(10,2) NOT NULL,
  total_price numeric(10,2) NOT NULL,
  payment_amount numeric(10,2) NOT NULL,
  change numeric(10,2) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp NULL DEFAULT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("customer_id") REFERENCES "customers" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE
);
