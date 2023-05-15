-- name: CreateCartProduct :one
INSERT INTO cart_products (
  product_id,
  user_id,
  quantity
) VALUES (
  sqlc.arg('product_id'),
  sqlc.arg('user_id'),
  sqlc.arg('quantity')
) RETURNING *;

-- name: GetCartProductsByUserId :many
SELECT * FROM cart_products
WHERE user_id = $1;

-- name: TruncateCartProductsTable :exec
TRUNCATE TABLE cart_products CASCADE;
