-- name: UpsertAnswer :one
INSERT INTO answers (
  attempt_id,
  question_id,
  selected_option_id,
  essay_answer,
  updated_at
) VALUES (
  $1, $2, $3, $4, now()
)
ON CONFLICT (attempt_id, question_id)
DO UPDATE SET
  selected_option_id = EXCLUDED.selected_option_id,
  essay_answer = EXCLUDED.essay_answer,
  updated_at = now()
RETURNING *;

-- name: AutoGradeMultipleChoice :exec
UPDATE answers a
SET is_correct = o.is_correct,
    score = CASE WHEN o.is_correct THEN q.points ELSE 0 END
FROM options o
JOIN questions q ON q.id = a.question_id
WHERE a.selected_option_id = o.id
  AND q.type = 'multiple_choice'
  AND a.id = $1;

-- name: HasUngradedEssay :one
SELECT EXISTS (
  SELECT 1
  FROM answers a
  JOIN questions q ON q.id = a.question_id
  WHERE a.attempt_id = $1
    AND q.type = 'essay'
    AND a.score IS NULL
);
