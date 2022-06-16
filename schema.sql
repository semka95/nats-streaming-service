CREATE TABLE orders(
    id BIGSERIAL PRIMARY KEY,
    data JSONB NOT NULL
);

CREATE INDEX idx_order_uid ON orders USING HASH ((data->'order_uid'));  