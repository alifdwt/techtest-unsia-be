-- name: ListOptionsByQuestionID :many
SELECT
  id,
  question_id,
  option_text
FROM options
WHERE question_id = $1
ORDER BY id;
