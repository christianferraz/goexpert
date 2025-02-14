-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS tasks (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  priority INT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
)
---- create above / drop below ----
DROP TABLE IF EXISTS taskS;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
