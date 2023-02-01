-- name: CreateEntry :one
INSERT INTO entries (
  account_id, amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetEntriesOfAccount :many
SELECT * from entries WHERE account_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: GetEntryById :one
SELECT * from entries WHERE id = $1;

