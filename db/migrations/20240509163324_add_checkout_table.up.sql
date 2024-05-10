CREATE TABLE IF NOT EXISTS "checkout" (
    "id" uuid UNIQUE NOT NULL DEFAULT (gen_random_uuid()) PRIMARY KEY,
    "customerId" uuid NOT NULL,
    "total" int NOT NULL,
    "paid" int NOT NULL,
    "change" int NOT NULL,
    "createdAt" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "checkout" ADD FOREIGN KEY ("customerId") REFERENCES "customer" ("id");