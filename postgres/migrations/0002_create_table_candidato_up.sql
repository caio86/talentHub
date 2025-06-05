CREATE TABLE IF NOT EXISTS candidatos (
  id    bigint PRIMARY KEY GENERATED ALWAYS as IDENTITY,
	name  varchar(80) NOT NULL,
	email varchar(80) NOT NULL,
	cpf   varchar(80) NOT NULL,
	phone varchar(80) NOT NULL
);

INSERT INTO candidatos (
  name,
  email,
  cpf,
  phone
) VALUES
  ("c1", "c1@teste.com", "00000000100", "083900000001"),
  ("c2", "c2@teste.com", "00000000200", "083900000002"),
  ("c3", "c3@teste.com", "00000000300", "083900000003"),
  ("c4", "c4@teste.com", "00000000400", "083900000004"),
  ("c5", "c5@teste.com", "00000000500", "083900000005"),
  ("c6", "c6@teste.com", "00000000600", "083900000006"),
  ("c7", "c7@teste.com", "00000000700", "083900000007"),
  ("c8", "c8@teste.com", "00000000800", "083900000008");
