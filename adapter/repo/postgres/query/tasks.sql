-- name: GetTaskById :one
SELECT *
FROM tasks AS t
WHERE t.id = sqlc.arg(id);

-- name: GetTasks :many
SELECT *
FROM tasks;

-- name: SaveTask :one
INSERT INTO tasks (id,
                   title,
                   description,
                   status,
                   due_date,
                   created_at)
VALUES (@id,
        @title,
        @description,
        @status,
        @due_date,
        now())
RETURNING *;
