ALTER TABLE users
ADD COLUMN email VARCHAR(254) NOT NULL;

CREATE UNIQUE INDEX users_email_unique_idx
ON users (LOWER(email));
