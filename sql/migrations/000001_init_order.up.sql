CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    store_id INT NOT NULL,
    merchant_order_id VARCHAR(50) NOT NULL,
    recipient_name VARCHAR(50) NOT NULL,
    recipient_phone VARCHAR(50) NOT NULL,
    recipient_address TEXT NOT NULL,
    recipient_city INT NOT NULL,
    recipient_zone INT NOT NULL,
    recipient_area INT NOT NULL,
    delivery_type INT NOT NULL,
    item_type INT NOT NULL,
    special_instruction TEXT,
    item_quantity INT NOT NULL,
    item_weight DOUBLE PRECISION NOT NULL,
    amount_to_collect DOUBLE PRECISION NOT NULL,
    item_description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL
);
