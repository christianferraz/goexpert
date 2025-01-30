-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS products (
  id UUID PRIMARY KEY,
  seller_id UUID NOT NULL REFERENCES users(id),
  product_name TEXT NOT NULL,
  description TEXT NOT NULL,
  baseprice FLOAT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
