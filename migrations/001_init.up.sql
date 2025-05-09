CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    ip TEXT NOT NULL,
    country TEXT,
    password_hash TEXT NOT NULL
);
