-- name: CreateNote :one
INSERT INTO notes (
    title, content
) VALUES (
    $1, $2
)
RETURNING *;
