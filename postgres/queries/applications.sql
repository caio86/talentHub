-- name: GetApplicationByID :one
SELECT 
    a.id,
    a.candidate_id,
    a.vacancy_id,
    a.application_date,
    a.score,
    s.status
FROM applications a
LEFT JOIN application_status s ON a.status_id = s.id
WHERE a.id = $1;

-- name: GetApplicationsByCandidateID :many
SELECT 
    a.id,
    a.candidate_id,
    a.vacancy_id,
    a.application_date,
    a.score,
    s.status
FROM applications a
LEFT JOIN application_status s ON a.status_id = s.id
WHERE a.candidate_id = $1
ORDER BY a.application_date DESC;

-- name: GetApplicationsByVacancyID :many
SELECT 
    a.id,
    a.candidate_id,
    a.vacancy_id,
    a.application_date,
    a.score,
    s.status
FROM applications a
LEFT JOIN application_status s ON a.status_id = s.id
WHERE a.vacancy_id = $1
ORDER BY a.application_date DESC;

-- name: GetAllApplications :many
SELECT 
    a.id,
    a.candidate_id,
    a.vacancy_id,
    a.application_date,
    a.score,
    s.status
FROM applications a
LEFT JOIN application_status s ON a.status_id = s.id
ORDER BY a.application_date DESC;

-- name: CreateApplication :one
INSERT INTO applications (candidate_id, vacancy_id, score, status_id)
VALUES ($1, $2, $3, (SELECT id FROM application_status WHERE status = $4))
RETURNING id, candidate_id, vacancy_id, application_date, score, 
    (SELECT status FROM application_status WHERE id = status_id) as status;

-- name: UpdateApplication :one
UPDATE applications 
SET score = COALESCE($2, score),
    status_id = COALESCE((SELECT id FROM application_status WHERE status = $3), status_id)
WHERE id = $1
RETURNING id, candidate_id, vacancy_id, application_date, score,
    (SELECT status FROM application_status WHERE id = status_id) as status;

-- name: DeleteApplication :exec
DELETE FROM applications WHERE id = $1;

-- name: GetApplicationStatusList :many
SELECT id, status FROM application_status ORDER BY id;

