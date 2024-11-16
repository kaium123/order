ALTER TABLE orders
    ADD COLUMN order_consignment_id VARCHAR(50) NOT NULL,  -- For the consignment ID
    ADD COLUMN order_type_id INT,          -- For the order type ID
    ADD COLUMN cod_fee INT,                -- For the COD fee
    ADD COLUMN promo_discount INT,        -- For the promotional discount
    ADD COLUMN discount INT,              -- For any additional discount
    ADD COLUMN delivery_fee INT,          -- For the delivery fee
    ADD COLUMN order_status VARCHAR(50) NOT NULL,         -- For the order status (Pending)
    ADD COLUMN order_type VARCHAR(50) NOT NULL,           -- For the order type (Delivery)
    ADD COLUMN order_amount INT,          -- For the order amount
    ADD COLUMN total_fee INT;             -- For the total fee
