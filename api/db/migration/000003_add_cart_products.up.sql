CREATE TABLE "cart_products" (
  "user_id" uuid NOT NULL,
  "product_id" uuid NOT NULL,
  "quantity" int NOT NULL,
  PRIMARY KEY ("user_id", "product_id")
);

CREATE INDEX ON "cart_products" ("user_id");

ALTER TABLE "cart_products" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "cart_products" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");
