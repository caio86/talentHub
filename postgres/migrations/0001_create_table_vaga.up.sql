CREATE TABLE IF NOT EXISTS vagas (
  id bigint PRIMARY KEY GENERATED ALWAYS as IDENTITY,
  name varchar(80) NOT NULL,
  description text NOT NULL,
  open boolean NOT NULL DEFAULT false,
  created_at timestamp NOT NULL DEFAULT NOW(),
  expires_at timestamp NOT NULL DEFAULT NOW() + '1 month'::interval
);

INSERT INTO vagas (
  name,
  description
) VALUES
  ('Banco de Talentos', 'reserva'),
  ('T1', 'descT1'),
  ('T2', 'descT2'),
  ('T3', 'descT3'),
  ('T4', 'descT4'),
  ('T5', 'descT5'),
  ('T6', 'descT6'),
  ('T7', 'descT7'),
  ('T8', 'descT8');
