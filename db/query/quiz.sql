-- name: GetQuizByID :one
SELECT *
FROM quizzes
WHERE id = $1;
