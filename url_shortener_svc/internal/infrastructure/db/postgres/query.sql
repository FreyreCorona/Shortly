-- name: SaveURL :one
INSERT INTO urls(original_url,short_code) VALUES($1,$2) RETURNING id,created_at;

-- name: GetByShortCode :one
SELECT id,original_url,short_code,created_at FROM urls WHERE short_code = $1;
