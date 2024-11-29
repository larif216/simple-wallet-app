CREATE TABLE IF NOT EXISTS "disbursements" (
    "id"         serial        PRIMARY KEY,
    "user_id"    int           NOT NULL,          
    "amount"     decimal(16,2) NOT NULL,
    "status"     smallint      NOT NULL,
    "created_at" timestamp     NOT NULL,
    "updated_at" timestamp     NOT NULL
);