-- name: RegisterApplication :one
INSERT INTO applications (
  candidate_id, vacancy_id, status_id
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetApplicationById :one
SELECT * FROM applications
WHERE id = $1;

-- name: GetFullApplicationById :one
SELECT a.*, s.status as status FROM applications a
LEFT JOIN application_status s ON a.status_id = s.id
WHERE a.id = $1;

-- name: ListApplications :many
SELECT * FROM applications
LIMIT $1 OFFSET $2;

-- name: ListAllApplications :many
SELECT * FROM applications;

-- name: ListFullApplications :many
SELECT a.*, s.status as status FROM applications a
LEFT JOIN application_status s ON a.status_id = s.id
LIMIT $1 OFFSET $2;

-- name: ListAllFullApplications :many
SELECT a.*, s.status as status FROM applications a
LEFT JOIN application_status s ON a.status_id = s.id;

-- name: CountApplications :one
SELECT count(*) FROM applications;

-- name: SearchApplicationsByCandidateId :many
SELECT a.* FROM applications a
JOIN candidates c ON a.candidate_id = c.id
WHERE c.id = $1;

-- name: SearchApplicationsByVacancyId :many
SELECT a.* FROM applications a
JOIN vacancies v ON a.vacancy_id = v.id
WHERE v.id = $1;

-- name: UnregisterApplication :exec
DELETE FROM applications
  WHERE id = $1;

-- name: UpdateApplication :one
UPDATE applications
  SET
    score = $2,
    status_id = $3
  WHERE id = $1
  RETURNING *;


-- application status

-- name: GetApplicationStatusByName :one
SELECT * FROM application_status
WHERE status = $1;

-- name: GetApplicationStatusById :one
SELECT * FROM application_status
WHERE id = $1;

-- name: ListApplicationStatuses :many
SELECT * FROM application_status;

