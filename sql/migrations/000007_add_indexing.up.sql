
-- Add indexes for email and user_name in users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_user_name ON users(user_name);
CREATE INDEX idx_access_tokens_user_id ON access_tokens(user_id);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

ALTER TABLE orders
    ADD COLUMN transfer_status INT DEFAULT 1 NOT NULL, -- Status of transfer, default 1
    ADD COLUMN archive INT DEFAULT 0 NOT NULL;        -- Archive flag, default 0

-- Add indexes for transfer_status and archive in orders
CREATE INDEX idx_orders_transfer_status ON orders(transfer_status);
CREATE INDEX idx_orders_archive ON orders(archive);



