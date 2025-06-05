-- name: ListCandidatos :many
SELECT * FROM candidatos
  LIMIT $1
  OFFSET $2;

-- name: GetCandidato :one
SELECT *
  FROM candidatos
  WHERE id = $1 LIMIT 1;

-- name: CreateCandidato :one
INSERT INTO candidatos (
	name ,
	email,
	cpf  ,
	phone
) VALUES ( $1, $2, $3, $4 )
  RETURNING *;

-- name: UpdateCandidato :one
UPDATE candidatos
  SET name  = $2,
      email = $3,
      cpf   = $4,
      phone = $5
  WHERE id = $1
  RETURNING *;

-- name: DeleteCandidato :exec
DELETE FROM candidatos
  WHERE id = $1;
