-- name: CreateProducts :one
INSERT INTO products(
  "id", 
  "seller_id",
  "product_name",
  "description",
  "baseprice",
  "auction_end",
  "is_sold") VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;

-- name: GetProductsById :one
SELECT * FROM products WHERE id = $1;