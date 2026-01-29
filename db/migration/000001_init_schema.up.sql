CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE quizzes (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  title TEXT NOT NULL,
  duration_seconds INT NOT NULL,
  max_attempts INT NOT NULL,
  auto_grading BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE questions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
  type TEXT NOT NULL CHECK (type IN ('multiple_choice', 'essay')),
  question_text TEXT NOT NULL,
  points INT NOT NULL,
  order_index INT NOT NULL
);

CREATE TABLE options (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  option_text TEXT NOT NULL,
  is_correct BOOLEAN DEFAULT FALSE
);

CREATE TABLE quiz_attempts (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  quiz_id UUID NOT NULL REFERENCES quizzes(id),
  user_id UUID NOT NULL,
  attempt_number INT NOT NULL,
  started_at TIMESTAMP NOT NULL,
  finished_at TIMESTAMP,
  status TEXT NOT NULL CHECK (
    status IN ('in_progress', 'submitted', 'waiting_assessment', 'graded')
  ),
  UNIQUE (quiz_id, user_id, attempt_number)
);

CREATE TABLE answers (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  attempt_id UUID NOT NULL REFERENCES quiz_attempts(id) ON DELETE CASCADE,
  question_id UUID NOT NULL REFERENCES questions(id),
  selected_option_id UUID,
  essay_answer TEXT,
  is_correct BOOLEAN,
  score INT,
  updated_at TIMESTAMP DEFAULT now(),
  UNIQUE (attempt_id, question_id)
);

CREATE TABLE manual_grades (
  answer_id UUID PRIMARY KEY REFERENCES answers(id) ON DELETE CASCADE,
  grader_id UUID NOT NULL,
  score INT NOT NULL,
  feedback TEXT,
  graded_at TIMESTAMP NOT NULL DEFAULT now()
);
