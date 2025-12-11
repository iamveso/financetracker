-- name: GetUser :one
-- Retrieves a single user by the email
SELECT * FROM users WHERE email = @email LIMIT 1;

-- name: CreateUser :exec
-- Creates a new user by email
INSERT INTO users (email, password) VALUES (@email, @password);
