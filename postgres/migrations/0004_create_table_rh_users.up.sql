CREATE TABLE IF NOT EXISTS rh_users (
  id int PRIMARY KEY GENERATED ALWAYS as IDENTITY,
  name varchar(80) NOT NULL,
  email varchar(80) NOT NULL UNIQUE,
  password varchar(255) NOT NULL,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);

-- Insert default RH user
INSERT INTO rh_users (name, email, password) VALUES 
  ('Gestor RH', 'rh@empresa.com', 'rhpassword');

