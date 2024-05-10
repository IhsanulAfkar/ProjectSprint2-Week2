CREATE TABLE IF NOT EXISTS "customer" (
    "id" uuid UNIQUE NOT NULL DEFAULT (gen_random_uuid()) PRIMARY KEY,
    "name" varchar(50) NOT NULL,
    "phone" varchar(20) UNIQUE NOT NULL,
    "createdAt" timestamp NOT NULL DEFAULT (now()),
    "updatedAt" timestamp NOT NULL DEFAULT (now()),
    "deletedAt" timestamp DEFAULT NULL
);