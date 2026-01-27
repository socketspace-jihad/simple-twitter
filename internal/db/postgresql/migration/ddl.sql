-- Jalankan ini di database default (biasanya postgres)
SELECT 'CREATE DATABASE "simple-twitter"'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'simple-twitter')\gexec

-- Pastikan extension uuid-ossp tersedia untuk generate UUID otomatis jika perlu
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    display_name VARCHAR(255) NOT NULL,
    username     VARCHAR(100) NOT NULL UNIQUE,
    born_date    DATE,
    address      TEXT
);

CREATE TABLE posts (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content    TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    user_id    UUID NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_users_username ON users(username);
