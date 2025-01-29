-- name: CreateUser :one
INSERT INTO USERS("id", "user_name", "email", "password_hash", "bio") VALUES ($1, $2, $3, $4, $5) RETURNING id;

-- name: GetUserById :one
SELECT id, user_name, email, password_hash, bio, created_at, updated_at FROM USERS WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, user_name, email, password_hash, bio, created_at, updated_at FROM USERS WHERE email = $1;

