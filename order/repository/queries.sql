-- name: CreateOrder :exec
INSERT INTO orders(order_uid, data) 
VALUES ($1, $2);

-- name: GetOrderByID :one
SELECT order_uid, data
FROM orders
WHERE order_uid = $1;

-- name: ListOrders :many
SELECT order_uid, data
FROM orders
WHERE id > $1
LIMIT $2;