-- name: CreateProduct :one
INSERT INTO products (
  name,
  description,
  price,
  stock_quantity,
  category_id,
  seller_id,
  image_url
) VALUES (
  sqlc.arg('name'),
  sqlc.narg('description'),
  sqlc.arg('price'),
  sqlc.arg('stock_quantity'),
  sqlc.arg('category_id'),
  sqlc.arg('seller_id'),
  sqlc.narg('image_url')
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
WHERE category_id = sqlc.narg('category_id') OR sqlc.narg('category_id') IS NULL
ORDER BY created_at
LIMIT $1
OFFSET $2;

-- name: CountProducts :one
SELECT count(*) FROM products;

-- name: ListProductsBySeller :many
SELECT * FROM products
WHERE seller_id = sqlc.arg('seller_id')
ORDER BY created_at
LIMIT $1
OFFSET $2;

-- name: CountProductsBySeller :one
SELECT count(*) FROM products
WHERE seller_id = sqlc.arg('seller_id');

-- name: AddProduct :one
INSERT INTO products (
  name,
  description,
  price,
  stock_quantity,
  category_id,
  seller_id,
  image_url
) VALUES (
  sqlc.arg('name'),
  sqlc.narg('description'),
  sqlc.arg('price'),
  sqlc.arg('stock_quantity'),
  sqlc.arg('category_id'),
  sqlc.arg('seller_id'),
  sqlc.narg('image_url')
) RETURNING *;

-- name: TruncateProductsTable :exec
TRUNCATE TABLE products CASCADE;
