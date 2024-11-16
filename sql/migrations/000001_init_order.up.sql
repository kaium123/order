CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY, -- Auto-incrementing primary key
    store_id INT NOT NULL, -- Store ID, required
    merchant_order_id VARCHAR(50) NOT NULL, -- Merchant Order ID, optional
    recipient_name VARCHAR(50) NOT NULL, -- Recipient name, required
    recipient_phone VARCHAR(50) NOT NULL, -- Recipient phone, required
    recipient_address TEXT NOT NULL, -- Recipient address, required
    recipient_city INT NOT NULL, -- Recipient city, required
    recipient_zone INT NOT NULL, -- Recipient zone, required
    recipient_area INT NOT NULL, -- Recipient area, required
    delivery_type INT NOT NULL, -- Delivery type, required
    item_type INT NOT NULL, -- Item type, required
    special_instruction TEXT, -- Special instruction, optional
    item_quantity INT NOT NULL, -- Item quantity, required
    item_weight DOUBLE PRECISION NOT NULL, -- Item weight, required
    amount_to_collect DOUBLE PRECISION NOT NULL, -- Amount to collect, required
    item_description TEXT, -- Item description, optional
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL, -- Created timestamp
    updated_at TIMESTAMP DEFAULT NULL, -- Updated timestamp
    deleted_at TIMESTAMP DEFAULT NULL -- Deleted timestamp for soft delete
);
