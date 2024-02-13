-- name: CreateNote :one
INSERT INTO notes (
    title, content
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetNote :one
SELECT * FROM notes WHERE id = $1;

-- name: ListNotes :many
SELECT * FROM notes;
