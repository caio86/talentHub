-- name: CreateVacancy :one
INSERT INTO vacancies (
  title, description, is_active, area_id, type_id, location, posted_date
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetVacancyByID :one
SELECT * FROM vacancies
WHERE id = $1;

-- name: GetFullVacancyByID :one
SELECT v.*, ea.name as area, et.name as type from vacancies v
LEFT JOIN employment_areas ea ON v.area_id = ea.id
LEFT JOIN employment_types et ON v.type_id = et.id
WHERE v.id = $1;

-- name: UpdateVacancy :one
UPDATE vacancies
SET
  title = $2,
  description = $3,
  location = $4,
  area_id = $5,
  type_id = $6
WHERE id = $1
RETURNING *;

-- name: DeleteVacancy :exec
DELETE FROM vacancies
WHERE id = $1;

-- name: ListVacancies :many
SELECT v.*, ea.name as area, et.name as type from vacancies v
LEFT JOIN employment_areas ea ON v.area_id = ea.id
LEFT JOIN employment_types et ON v.type_id = et.id
ORDER BY title
LIMIT $1 OFFSET $2;

-- name: ListAllVacancies :many
SELECT v.*, ea.name as area, et.name as type from vacancies v
LEFT JOIN employment_areas ea ON v.area_id = ea.id
LEFT JOIN employment_types et ON v.type_id = et.id
ORDER BY title;

-- name: CountVacancies :one
SELECT count(*) FROM vacancies;

-- name: OpenVacancy :exec
UPDATE vacancies
  SET is_active = true
  WHERE id = $1;

-- name: CloseVacancy :exec
UPDATE vacancies
  SET is_active = false
  WHERE id = $1;

-- Area
-- name: GetAreaByID :one
SELECT * FROM employment_areas
WHERE id = $1;

-- name: GetAreaByName :one
SELECT * FROM employment_areas
WHERE name = $1;

-- Type
-- name: GetTypeByID :one
SELECT * FROM employment_types
WHERE id = $1;

-- name: GetTypeByName :one
SELECT * FROM employment_types
WHERE name = $1;


-- Requirements

-- name: AddRequirement :one
INSERT INTO requirements (name)
  VALUES ($1)
  ON CONFLICT DO NOTHING
RETURNING *;

-- name: GetRequirementByName :one
SELECT * FROM requirements
WHERE name = $1;

-- name: GetRequirementsByVacancyID :many
SELECT name FROM requirements r
JOIN vacancy_requirements vr ON r.id = vr.requirement_id
JOIN vacancies v ON v.id = vr.vacancy_id
WHERE v.id = $1;


-- name: AddVacancyRequirement :exec
INSERT INTO vacancy_requirements (vacancy_id, requirement_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: RemoveVacancyRequirement :exec
DELETE FROM vacancy_requirements
WHERE vacancy_id = $1 AND requirement_id = $2;
