CREATE TABLE IF NOT EXISTS employment_areas (
  id int PRIMARY KEY GENERATED ALWAYS as IDENTITY,
  name varchar(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS employment_types (
  id int PRIMARY KEY GENERATED ALWAYS as IDENTITY,
  name varchar(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS vacancies (
  id int PRIMARY KEY GENERATED ALWAYS as IDENTITY,
  title varchar(255) NOT NULL,
  description text,
  is_active boolean NOT NULL DEFAULT false,
  area_id int REFERENCES employment_areas(id),
  type_id int REFERENCES employment_types(id),
  location varchar(255),
  posted_date DATE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS requirements (
  id int PRIMARY KEY GENERATED ALWAYS as IDENTITY,
  name varchar(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS vacancy_requirements (
  vacancy_id int NOT NULL REFERENCES vacancies(id) ON DELETE CASCADE,
  requirement_id int NOT NULL REFERENCES requirements(id),
  PRIMARY KEY (vacancy_id, requirement_id)
);

--- INSERTS

-- Areas
INSERT INTO employment_areas (name) VALUES
('Tecnologia'),
('Marketing'),
('Recursos Humanos');

-- Employment Types
INSERT INTO employment_types (name) VALUES
('Tempo Integral'),
('Estágio');

-- Vancancies
INSERT INTO vacancies (title, description, area_id, type_id, location)
SELECT
  v.title,
  v.description,
  ea.id AS area_id,
  et.id AS type_id,
  v.location
FROM (VALUES
  ('Banco de Talentos', 'Reserva', 'Tecnologia', 'Tempo Integral', 'Remoto'),
  ('Desenvolvedor Frontend Angular Pleno', 'Procuramos um desenvolvedor Frontend com experiência em Angular para se juntar à nossa equipa inovadora.', 'Tecnologia', 'Tempo Integral', 'Lisboa, Portugal'),
  ('Analista de Marketing Digital', 'Vaga para profissional de marketing digital com foco em SEO e campanhas pagas.', 'Marketing', 'Tempo Integral', 'Porto, Portugal'),
  ('Estágio em Recursos Humanos', 'Oportunidade de estágio para estudantes de RH ou Psicologia Organizacional.', 'Recursos Humanos', 'Estágio', 'Remoto'),
  ('Estágio em Recursos Humanos', 'Oportunidade de estágio para estudantes de RH ou Psicologia Organizacional.', 'Recursos Humanos', 'Estágio', 'Remoto')
) AS v(title, description, area, type, location)
JOIN employment_areas ea ON ea.name = v.area
JOIN employment_types et ON et.name = v.type;

-- requirements
INSERT INTO requirements (name)
  SELECT DISTINCT
  UNNEST(v.reqs) FROM
  (VALUES
    ( ARRAY['Angular 12+','TypeScript','RxJS','HTML5','CSS3/SCSS','Git'] ),
    ( ARRAY['SEO','Google Ads','Facebook Ads','Google Analytics','Marketing de Conteúdo'] ),
    ( ARRAY['Cursando RH ou Psicologia','Boa comunicação','Proatividade'] )
  ) as v(reqs);

INSERT INTO vacancy_requirements (vacancy_id, requirement_id)
  SELECT v.id, r.id FROM vacancies v
  JOIN (VALUES
    (3, 'Angular 12+'),
    (3, 'TypeScript'),
    (3, 'RxJS'),
    (3, 'HTML5'),
    (3, 'CSS3/SCSS'),
    (3, 'Git'),
    (1, 'SEO'),
    (1, 'Google Ads'),
    (1, 'Facebook Ads'),
    (1, 'Google Analytics'),
    (1, 'Marketing de Conteúdo'),
    (4, 'Cursando RH ou Psicologia'),
    (4, 'Boa comunicação'),
    (4, 'Proatividade'),
    (5, 'Cursando RH ou Psicologia'),
    (5, 'Boa comunicação'),
    (5, 'Proatividade')
  ) as vr(vacancy_id, requirement_name)
  ON v.id = vr.vacancy_id
  JOIN requirements r ON r.name = vr.requirement_name;
