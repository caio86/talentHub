CREATE TABLE IF NOT EXISTS application_status (
  id int PRIMARY KEY GENERATED ALWAYS as IDENTITY,
  status VARCHAR(20) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS applications (
  id int PRIMARY KEY GENERATED ALWAYS as IDENTITY,
  candidate_id int NOT NULL REFERENCES candidates(id),
  vacancy_id int NOT NULL REFERENCES vacancies(id),
  application_date TIMESTAMP NOT NULL DEFAULT NOW(),
  score int NOT NULL DEFAULT 0 CHECK (score BETWEEN 0 and 100),
  status_id int REFERENCES application_status(id),
  UNIQUE (candidate_id, vacancy_id)
);

--- INSERTS

INSERT INTO application_status (status)
VALUES
  ('Recebido'),
  ('Em análise'),
  ('Rejeitado'),
  ('Aprovado');

INSERT INTO applications (
    candidate_id, vacancy_id, status_id
  ) SELECT
    c.id,
    a.vacancy_id,
    s.id AS status_id
  FROM (VALUES
    ('paulo.silva@email.com', 2, NULL),
    ('ana.silva@email.com', 2, NULL),
    ('carlos.pereira@email.com', 2, NULL),
    ('paulo.silva@email.com', 3, 'Recebido'),
    ('ana.silva@email.com', 3, 'Recebido'),
    ('carlos.pereira@email.com', 4, NULL)
  ) AS a(email, vacancy_id, status)
  JOIN candidates c ON c.email = a.email
  LEFT JOIN application_status s ON s.status = a.status;
