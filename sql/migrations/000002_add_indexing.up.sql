-- Create an index for the `store_id` column
CREATE INDEX idx_orders_store_id ON orders(store_id);

-- Create an index for the `merchant_order_id` column
CREATE INDEX idx_orders_merchant_order_id ON orders(merchant_order_id);