-- Run against your Postgres database once: psql $DATABASE_URL -f schema.sql

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS reviews (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    language TEXT NOT NULL,
    original_code TEXT NOT NULL,
    bug_explanation TEXT NOT NULL,
    fixed_code TEXT NOT NULL,
    suggestions TEXT NOT NULL,
    is_helpful BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- //create table users(
--     id text primary,
--     email text unique not Null
--     password_hash text not null
-- )

-- //
-- create table reviews(
--     id text primary,
    
-- )