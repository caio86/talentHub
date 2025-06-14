// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: candidato.sql

package repository

import (
	"context"
)

const countCandidatos = `-- name: CountCandidatos :one
SELECT count(*) FROM candidatos
`

func (q *Queries) CountCandidatos(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countCandidatos)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createCandidato = `-- name: CreateCandidato :one
INSERT INTO candidatos (
	name ,
	email,
	cpf  ,
	phone
) VALUES ( $1, $2, $3, $4 )
  RETURNING id, name, email, cpf, phone
`

type CreateCandidatoParams struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Cpf   string `json:"cpf"`
	Phone string `json:"phone"`
}

func (q *Queries) CreateCandidato(ctx context.Context, arg CreateCandidatoParams) (Candidato, error) {
	row := q.db.QueryRow(ctx, createCandidato,
		arg.Name,
		arg.Email,
		arg.Cpf,
		arg.Phone,
	)
	var i Candidato
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Cpf,
		&i.Phone,
	)
	return i, err
}

const deleteCandidato = `-- name: DeleteCandidato :exec
DELETE FROM candidatos
  WHERE id = $1
`

func (q *Queries) DeleteCandidato(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteCandidato, id)
	return err
}

const getCandidato = `-- name: GetCandidato :one
SELECT id, name, email, cpf, phone
  FROM candidatos
  WHERE id = $1 LIMIT 1
`

func (q *Queries) GetCandidato(ctx context.Context, id int64) (Candidato, error) {
	row := q.db.QueryRow(ctx, getCandidato, id)
	var i Candidato
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Cpf,
		&i.Phone,
	)
	return i, err
}

const listAllCandidatos = `-- name: ListAllCandidatos :many
SELECT id, name, email, cpf, phone FROM candidatos
`

func (q *Queries) ListAllCandidatos(ctx context.Context) ([]Candidato, error) {
	rows, err := q.db.Query(ctx, listAllCandidatos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Candidato
	for rows.Next() {
		var i Candidato
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Cpf,
			&i.Phone,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listCandidatos = `-- name: ListCandidatos :many
SELECT id, name, email, cpf, phone FROM candidatos
  LIMIT $1
  OFFSET $2
`

type ListCandidatosParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListCandidatos(ctx context.Context, arg ListCandidatosParams) ([]Candidato, error) {
	rows, err := q.db.Query(ctx, listCandidatos, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Candidato
	for rows.Next() {
		var i Candidato
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Cpf,
			&i.Phone,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCandidato = `-- name: UpdateCandidato :one
UPDATE candidatos
  SET name  = $2,
      email = $3,
      cpf   = $4,
      phone = $5
  WHERE id = $1
  RETURNING id, name, email, cpf, phone
`

type UpdateCandidatoParams struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Cpf   string `json:"cpf"`
	Phone string `json:"phone"`
}

func (q *Queries) UpdateCandidato(ctx context.Context, arg UpdateCandidatoParams) (Candidato, error) {
	row := q.db.QueryRow(ctx, updateCandidato,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Cpf,
		arg.Phone,
	)
	var i Candidato
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Cpf,
		&i.Phone,
	)
	return i, err
}
