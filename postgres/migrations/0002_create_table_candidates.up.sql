CREATE TABLE IF NOT EXISTS candidates (
  id int PRIMARY KEY GENERATED ALWAYS as IDENTITY,
	name varchar(80) NOT NULL,
	email varchar(80) NOT NULL UNIQUE,
  password varchar(255) NOT NULL,
	phone varchar(80),
  address text,
  linkedin text,
  resume_url text
);

CREATE TABLE IF NOT EXISTS candidate_skills (
  candidate_id int NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
  skill text NOT NULL,
  PRIMARY KEY(candidate_id, skill)
);

CREATE TABLE IF NOT EXISTS candidate_interests (
  candidate_id int NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
  interest text NOT NULL,
  PRIMARY KEY(candidate_id, interest)
);

CREATE TABLE IF NOT EXISTS experiences (
  id int PRIMARY KEY GENERATED ALWAYS as IDENTITY,
  candidate_id int NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
  company text NOT NULL,
  role text NOT NULL,
  years int NOT NULL CHECK (years > 0)
);

CREATE TABLE IF NOT EXISTS education (
  id int PRIMARY KEY GENERATED ALWAYS as IDENTITY,
  candidate_id int NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
  institution text NOT NULL,
  course text NOT NULL,
  level text NOT NULL
);

--- INSERTS

INSERT INTO candidates (
    name, email, password, phone, address, linkedin, resume_url
  ) VALUES
    ('Ana Silva', 'ana.silva@email.com', '123', '912345678', 'Rua Exemplo, 123, Lisboa', 'linkedin.com/in/anasilva', NULL),
    ('Paulo Silva', 'paulo.silva@email.com', '123', '912345678', 'Rua Exemplo, 123, Lisboa', 'linkedin.com/in/anasilva', NULL),
    ('Carlos Pereira', 'carlos.pereira@email.com', '456', '923456789', 'Avenida Teste, 456, Porto', NULL, '/path/to/carlos_cv.pdf'),
    ('Mariana Costa', 'mariana.costa@email.com', '789', '934567890', 'Travessa Modelo, 789, Coimbra', 'linkedin.com/in/marianacosta', NULL);

INSERT INTO candidate_skills (
    candidate_id, skill
  ) SELECT c.id, v.skill FROM (VALUES
    ('ana.silva@email.com', 'Angular'),
    ('ana.silva@email.com', 'TypeScript'),
    ('ana.silva@email.com', 'JavaScript'),
    ('ana.silva@email.com', 'HTML'),
    ('ana.silva@email.com', 'CSS'),
    ('paulo.silva@email.com', 'Angular'),
    ('paulo.silva@email.com', 'TypeScript'),
    ('paulo.silva@email.com', 'JavaScript'),
    ('paulo.silva@email.com', 'HTML'),
    ('paulo.silva@email.com', 'CSS'),
    ('carlos.pereira@email.com', 'SEO'),
    ('carlos.pereira@email.com', 'Google Ads'),
    ('carlos.pereira@email.com', 'Social Media')
  ) as v(email, skill)
  JOIN candidates c ON c.email = v.email;

INSERT INTO candidate_interests (
    candidate_id, interest
  ) SELECT c.id, v.interest FROM (VALUES
      ('ana.silva@email.com', 'Tecnologia'),
      ('ana.silva@email.com', 'Desenvolvimento Web'),
      ('paulo.silva@email.com', 'Tecnologia'),
      ('paulo.silva@email.com', 'Desenvolvimento Web'),
      ('carlos.pereira@email.com', 'Marketing'),
      ('carlos.pereira@email.com', 'Publicidade'),
      ('mariana.costa@email.com', 'Recursos Humanos'),
      ('mariana.costa@email.com', 'Gestão de Pessoas')
  ) as v(email, interest)
  JOIN candidates c ON c.email = v.email;

-- Experiences
INSERT INTO experiences (
    candidate_id, company, role, years
  ) SELECT c.id, v.company, v.role, v.years FROM ( VALUES
    ('ana.silva@email.com', 'Empresa X', 'Desenvolvedor Frontend Jr', 2),
    ('paulo.silva@email.com', 'Empresa X', 'Desenvolvedor Frontend Jr', 2)
  ) as v(email, company, role, years)
  JOIN candidates c ON c.email = v.email;

-- Education
INSERT INTO education (
    candidate_id, institution, course, level
  ) SELECT c.id, v.institution, v.course, v.level FROM ( VALUES
    ('ana.silva@email.com', 'Universidade Y', 'Engenharia Informática', 'Licenciatura'),
    ('paulo.silva@email.com', 'Universidade Y', 'Engenharia Informática', 'Licenciatura'),
    ('carlos.pereira@email.com', 'Escola Z', 'Marketing Digital', 'Curso Técnico')
  ) as v(email, institution, course, level)
  JOIN candidates c ON c.email = v.email;
