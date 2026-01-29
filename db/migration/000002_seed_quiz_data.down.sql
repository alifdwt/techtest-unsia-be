-- =========================
-- CLEAN RUNTIME DATA FIRST
-- =========================
DELETE FROM manual_grades;
DELETE FROM answers;
DELETE FROM quiz_attempts;

-- =========================
-- THEN SEED DATA
-- =========================
DELETE FROM options;
DELETE FROM questions;
DELETE FROM quizzes;
