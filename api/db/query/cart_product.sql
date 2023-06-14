-- name: CreateCartProduct :one
INSERT INTO cart_products (
  user_id,
  product_id,
  quantity
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: UpdateCartProduct :one
UPDATE cart_products
SET
  quantity = $3
WHERE user_id = $1 AND product_id = $2
RETURNING *;

-- name: DeleteCartProduct :exec
DELETE FROM cart_products
WHERE user_id = $1 AND product_id = $2;

-- name: GetCartProductByUserIDAndProductID :one
SELECT * FROM cart_products
WHERE user_id = $1 AND product_id = $2
ORDER BY created_at;

-- name: GetCartProductsByUserID :many
SELECT * FROM cart_products
WHERE user_id = $1;

-- name: TruncateCartProductsTable :exec
TRUNCATE TABLE cart_products CASCADE;
