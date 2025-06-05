-- name: ListVagas :many
SELECT * FROM vagas
  LIMIT $1
  OFFSET $2;

-- name: GetVaga :one
SELECT *
  FROM vagas
  WHERE id = $1 LIMIT 1;

-- name: CreateVaga :one
INSERT INTO vagas (
  name       ,
  description,
  open       ,
  created_at ,
  expires_at
) VALUES ( $1, $2, $3, $4, $5 )
  RETURNING *;

-- name: UpdateVaga :one
UPDATE vagas
  SET name        = $2,
      description = $3,
      open        = $4,
      expires_at  = $5
  WHERE id = $1
  RETURNING *;

-- name: DeleteVaga :exec
DELETE FROM vagas
  WHERE id = $1;

