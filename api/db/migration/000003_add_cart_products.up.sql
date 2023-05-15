CREATE TABLE "cart_products" (
  "id" uuid PRIMARY KEY,
  "product_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "quantity" int NOT NULL
);

ALTER TABLE "cart_products" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "cart_products" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");
