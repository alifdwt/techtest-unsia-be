-- name: CountUserAttempts :one
SELECT COUNT(*)
FROM quiz_attempts
WHERE quiz_id = $1
  AND user_id = $2;

-- name: CreateQuizAttempt :one
INSERT INTO quiz_attempts (
  quiz_id,
  user_id,
  attempt_number,
  started_at,
  status
) VALUES (
  $1, $2, $3, now(), 'in_progress'
)
RETURNING *;

-- name: UpdateAttemptStatus :exec
UPDATE quiz_attempts
SET status = $2,
    finished_at = CASE
      WHEN $2 IN ('submitted', 'waiting_assessment', 'graded')
      THEN now()
      ELSE finished_at
    END
WHERE id = $1;

-- name: GetAttemptByID :one
SELECT *
FROM quiz_attempts
WHERE id = $1;

-- name: GetQuizDurationByAttemptID :one
SELECT q.duration_seconds
FROM quizzes q
JOIN quiz_attempts qa ON qa.quiz_id = q.id
WHERE qa.id = $1;

-- name: GetActiveAttempt :one
SELECT *
FROM quiz_attempts
WHERE quiz_id = $1
  AND user_id = $2
  AND status = 'in_progress'
ORDER BY started_at DESC
LIMIT 1;
