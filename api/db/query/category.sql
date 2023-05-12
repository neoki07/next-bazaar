-- name: CreateCategory :one
INSERT INTO categories (
  name
) VALUES (
  $1
) RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: GetCategoriesByIDs :many
SELECT * FROM categories
WHERE id = ANY((sqlc.arg('ids'))::uuid[])
ORDER BY id;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;

-- name: TruncateCategoriesTable :exec
TRUNCATE TABLE categories CASCADE;
