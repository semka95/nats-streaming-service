-- name: CreateOrder :exec
INSERT INTO orders(
    data
) VALUES (
    $1
);

-- name: GetOrderByID :one
SELECT data 
FROM orders
WHERE data ->> 'order_uid' = $1;