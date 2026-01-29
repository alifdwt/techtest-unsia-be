-- name: ListAnswersWithQuestions :many
SELECT
  q.id            AS question_id,
  q.question_text,
  q.type,
  q.points,
  a.selected_option_id,
  o.option_text   AS selected_option_text,
  o.is_correct,
  a.essay_answer,
  a.score
FROM answers a
JOIN questions q ON q.id = a.question_id
LEFT JOIN options o ON o.id = a.selected_option_id
WHERE a.attempt_id = $1
ORDER BY q.order_index;
