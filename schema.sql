CREATE TABLE orders(
    id BIGSERIAL PRIMARY KEY,
    order_uid VARCHAR UNIQUE NOT NULL,
    data JSON NOT NULL
);

CREATE INDEX idx_order_uid ON orders USING HASH (order_uid);