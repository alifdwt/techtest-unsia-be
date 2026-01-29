-- name: ListQuestionsByQuizID :many
SELECT
  id,
  quiz_id,
  type,
  question_text,
  points,
  order_index
FROM questions
WHERE quiz_id = $1
ORDER BY order_index ASC;
