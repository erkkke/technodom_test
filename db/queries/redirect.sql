-- name: CreateRedirect :one
INSERT INTO redirects (
    active_link, history_link
) VALUES (
             $1, $2
         )
RETURNING *;

-- name: GetRedirect :one
SELECT * FROM redirects
WHERE id = $1 LIMIT 1;

-- name: ListRedirects :many
SELECT * FROM redirects
ORDER BY id
LIMIT $1
    OFFSET $2;

-- name: UpdateRedirect :one
UPDATE redirects
set active_link = $2, history_link = $3
WHERE id = $1
RETURNING *;

-- name: DeleteRedirect :one
DELETE FROM redirects
WHERE id = $1
RETURNING *;

-- name: GetRedirectByHistoryLink :one
SELECT active_link FROM redirects WHERE history_link = $1;

-- name: GetRedirectByActiveLink :one
SELECT active_link FROM redirects WHERE active_link = $1;