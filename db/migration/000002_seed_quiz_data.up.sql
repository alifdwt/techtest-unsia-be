-- ==================================================
-- QUIZ
-- ==================================================
INSERT INTO quizzes (id, title, duration_seconds, max_attempts)
VALUES (
  '11111111-1111-1111-1111-111111111111',
  'Golang & PostgreSQL Fundamentals',
  2700,
  2
);

-- ==================================================
-- QUESTIONS
-- ==================================================

INSERT INTO questions (id, quiz_id, type, question_text, points, order_index)
VALUES (
  '22222222-2222-2222-2222-222222222201',
  '11111111-1111-1111-1111-111111111111',
  'multiple_choice',
  'Apa kegunaan utama goroutine dalam Golang?',
  10,
  1
);

INSERT INTO questions (id, quiz_id, type, question_text, points, order_index)
VALUES (
  '22222222-2222-2222-2222-222222222202',
  '11111111-1111-1111-1111-111111111111',
  'multiple_choice',
  'Package apa yang direkomendasikan untuk mengelola koneksi PostgreSQL di Golang?',
  10,
  2
);

INSERT INTO questions (id, quiz_id, type, question_text, points, order_index)
VALUES (
  '22222222-2222-2222-2222-222222222203',
  '11111111-1111-1111-1111-111111111111',
  'multiple_choice',
  'Apa perbedaan utama antara PRIMARY KEY dan UNIQUE di PostgreSQL?',
  10,
  3
);

INSERT INTO questions (id, quiz_id, type, question_text, points, order_index)
VALUES (
  '22222222-2222-2222-2222-222222222204',
  '11111111-1111-1111-1111-111111111111',
  'essay',
  'Jelaskan perbedaan antara goroutine dan thread, serta bagaimana scheduler Golang mengelolanya.',
  20,
  4
);

INSERT INTO questions (id, quiz_id, type, question_text, points, order_index)
VALUES (
  '22222222-2222-2222-2222-222222222205',
  '11111111-1111-1111-1111-111111111111',
  'essay',
  'Jelaskan fungsi index di PostgreSQL dan kapan penggunaannya dapat berdampak negatif.',
  20,
  5
);

-- ==================================================
-- OPTIONS
-- ==================================================

INSERT INTO options (id, question_id, option_text, is_correct)
VALUES (
  '33333333-3333-3333-3333-333333333201',
  '22222222-2222-2222-2222-222222222201',
  'Menjalankan fungsi secara concurrent dengan biaya ringan',
  TRUE
);

INSERT INTO options (id, question_id, option_text, is_correct)
VALUES (
  '33333333-3333-3333-3333-333333333202',
  '22222222-2222-2222-2222-222222222201',
  'Menggantikan fungsi main dalam program Go',
  FALSE
);

INSERT INTO options (id, question_id, option_text, is_correct)
VALUES (
  '33333333-3333-3333-3333-333333333203',
  '22222222-2222-2222-2222-222222222201',
  'Membuat proses OS baru',
  FALSE
);

INSERT INTO options (id, question_id, option_text, is_correct)
VALUES (
  '33333333-3333-3333-3333-333333333204',
  '22222222-2222-2222-2222-222222222202',
  'database/sql dengan driver pgx',
  TRUE
);

INSERT INTO options (id, question_id, option_text, is_correct)
VALUES (
  '33333333-3333-3333-3333-333333333205',
  '22222222-2222-2222-2222-222222222202',
  'gorm secara langsung tanpa driver',
  FALSE
);

INSERT INTO options (id, question_id, option_text, is_correct)
VALUES (
  '33333333-3333-3333-3333-333333333206',
  '22222222-2222-2222-2222-222222222202',
  'net/http',
  FALSE
);

INSERT INTO options (id, question_id, option_text, is_correct)
VALUES (
  '33333333-3333-3333-3333-333333333207',
  '22222222-2222-2222-2222-222222222203',
  'PRIMARY KEY tidak boleh NULL dan hanya satu per tabel',
  TRUE
);

INSERT INTO options (id, question_id, option_text, is_correct)
VALUES (
  '33333333-3333-3333-3333-333333333208',
  '22222222-2222-2222-2222-222222222203',
  'UNIQUE otomatis menjadi clustered index',
  FALSE
);

INSERT INTO options (id, question_id, option_text, is_correct)
VALUES (
  '33333333-3333-3333-3333-333333333209',
  '22222222-2222-2222-2222-222222222203',
  'UNIQUE dapat memiliki nilai duplikat',
  FALSE
);
