CREATE TABLE IF NOT EXISTS "checkoutItem" (
    "id" uuid UNIQUE NOT NULL DEFAULT (gen_random_uuid()) PRIMARY KEY,
    "checkoutId" uuid NOT NULL,
    "productId" uuid NOT NULL,
    "price" int NOT NULL,
    "quantity" int NOT NULL,
    "createdAt" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "checkoutItem" ADD FOREIGN KEY ("checkoutId") REFERENCES "checkout" ("id");