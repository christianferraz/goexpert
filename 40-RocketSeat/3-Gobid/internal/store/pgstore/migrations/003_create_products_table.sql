-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS products (
  id UUID PRIMARY KEY,
  seller_id UUID NOT NULL REFERENCES users(id),
  product_name TEXT NOT NULL,
  description TEXT NOT NULL,
  baseprice FLOAT NOT NULL,
  auction_end TIMESTAMPTZ NOT NULL,
  is_sold BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
