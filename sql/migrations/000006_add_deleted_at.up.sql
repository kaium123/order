-- Add 'deleted_at' column to the 'users' table
ALTER TABLE users
    ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL;

-- Add 'deleted_at' column to the 'access_tokens' table
ALTER TABLE access_tokens
    ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL;

-- Add 'deleted_at' column to the 'refresh_tokens' table
ALTER TABLE refresh_tokens
    ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL;
