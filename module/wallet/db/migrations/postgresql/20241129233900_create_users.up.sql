CREATE TABLE IF NOT EXISTS "users" (
    "id"              serial        PRIMARY KEY,
    "name"            varchar(255)  NOT NULL,          
    "balance"         decimal(16,2) DEFAULT 0,
    "pending_balance" decimal(16,2) DEFAULT 0,
    "created_at"      timestamp     NOT NULL,
    "updated_at"      timestamp     NOT NULL
);