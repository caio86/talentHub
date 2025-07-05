-- name: CreateCandidate :one
INSERT INTO candidates (
  name, email, password, phone, address, linkedin, resume_url
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetCandidateByID :one
SELECT * FROM candidates
WHERE id = $1;

-- name: GetCandidateByEmail :one
SELECT * FROM candidates
WHERE email = $1;

-- name: UpdateCandidate :one
UPDATE candidates
SET
  name = $2,
  phone = $3,
  address = $4,
  linkedin = $5,
  resume_url = $6
WHERE id = $1
RETURNING *;

-- name: DeleteCandidate :exec
DELETE FROM candidates
WHERE id = $1;

-- name: ListCandidates :many
SELECT * FROM candidates
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: ListAllCandidates :many
SELECT * FROM candidates
ORDER BY name;

-- name: CountCandidates :one
SELECT count(*) FROM candidates;


-- Skills


-- name: SearchCandidatesBySkill :many
SELECT c.*
FROM candidates c
JOIN candidate_skills cs ON c.id = cs.candidate_id
WHERE cs.skill = $1;

-- name: GetCandidateSkills :many
SELECT *
FROM candidate_skills
WHERE candidate_id = $1;

-- name: AddCandidateSkill :exec
INSERT INTO candidate_skills (candidate_id, skill)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: RemoveCandidateSkill :exec
DELETE FROM candidate_skills
WHERE candidate_id = $1 AND skill = $2;


-- Interests

-- name: SearchCandidatesByInterest :many
SELECT c.*
FROM candidates c
JOIN candidate_interests ci ON c.id = ci.candidate_id
WHERE ci.interest = $1;

-- name: GetCandidateInterests :many
SELECT *
FROM candidate_interests
WHERE candidate_id = $1;

-- name: AddCandidateInterest :exec
INSERT INTO candidate_interests (candidate_id, interest)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: RemoveCandidateInterest :exec
DELETE FROM candidate_interests
WHERE candidate_id = $1 AND interest = $2;

-- Education

-- name: GetCandidateEducations :many
SELECT *
FROM education
WHERE candidate_id = $1;

-- name: AddCandidateEducation :one
INSERT INTO education (candidate_id, institution, course, level)
VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: RemoveCandidateEducation :exec
DELETE FROM education
WHERE id = $1;

-- Experience

-- name: GetCandidateExperiences :many
SELECT *
FROM experiences
WHERE candidate_id = $1;

-- name: AddCandidateExperience :one
INSERT INTO experiences (candidate_id, company, role, years)
VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: RemoveCandidateExperience :exec
DELETE FROM experiences
WHERE id = $1;
