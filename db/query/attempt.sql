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
