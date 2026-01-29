package service

import (
	"context"
	"errors"
	"time"

	"github.com/alifdwt/techtest-unsia-be/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type QuizService struct {
	q *db.Queries
}

func NewQuizService(q *db.Queries) *QuizService {
	return &QuizService{
		q: q,
	}
}

func (s *QuizService) StartQuiz(ctx context.Context, quizID, userID uuid.UUID) (*db.QuizAttempt, *db.Quiz, time.Time, []QuestionDTO, error) {
	pgQuizID := toPgUUID(quizID)
	pgUserID := toPgUUID(userID)

	quiz, err := s.q.GetQuizByID(ctx, pgQuizID)
	if err != nil {
		return nil, nil, time.Time{}, nil, err
	}

	activeAttempt, err := s.q.GetActiveAttempt(ctx, db.GetActiveAttemptParams{
		QuizID: pgQuizID,
		UserID: pgUserID,
	})

	if err == nil {
		deadline := activeAttempt.StartedAt.Time.Add(
			time.Duration(quiz.DurationSeconds) * time.Second,
		)

		if time.Now().Before(deadline) {
			questions, err := s.loadQuestions(ctx, pgQuizID)
			if err != nil {
				return nil, nil, time.Time{}, nil, err
			}
			return &activeAttempt, &quiz, deadline, questions, nil
		}
	}

	count, err := s.q.CountUserAttempts(ctx, db.CountUserAttemptsParams{
		QuizID: pgQuizID,
		UserID: pgUserID,
	})
	if err != nil {
		return nil, nil, time.Time{}, nil, err
	}

	if int(count) >= int(quiz.MaxAttempts) {
		return nil, nil, time.Time{}, nil, errors.New("max attempts reached")
	}

	attempt, err := s.q.CreateQuizAttempt(ctx, db.CreateQuizAttemptParams{
		QuizID:        pgQuizID,
		UserID:        pgUserID,
		AttemptNumber: int32(count + 1),
	})

	deadline := attempt.StartedAt.Time.Add(
		time.Duration(quiz.DurationSeconds) * time.Second,
	)

	questions, err := s.loadQuestions(ctx, pgQuizID)
	if err != nil {
		return nil, nil, time.Time{}, nil, err
	}

	return &attempt, &quiz, deadline, questions, nil
}

func (s *QuizService) SubmitAnswer(
	ctx context.Context,
	attemptID, questionID uuid.UUID,
	selectedOptionID *uuid.UUID,
	essayAnswer *string,
	final bool,
) error {
	pgAttemptID := toPgUUID(attemptID)
	pgQuestionID := toPgUUID(questionID)

	attempt, err := s.q.GetAttemptByID(ctx, pgAttemptID)
	if err != nil {
		return errors.New("attempt not found")
	}

	if attempt.Status != "in_progress" {
		return errors.New("attempt already submitted")
	}

	duration, err := s.q.GetQuizDurationByAttemptID(ctx, pgAttemptID)
	if err != nil {
		return err
	}

	deadline := attempt.StartedAt.Time.Add(
		time.Duration(duration) * time.Second,
	)
	if time.Now().After(deadline) {
		return errors.New("attempt expired")
	}

	var pgOptionID pgtype.UUID
	if selectedOptionID != nil {
		pgOptionID = toPgUUID(*selectedOptionID)
	}

	answer, err := s.q.UpsertAnswer(ctx, db.UpsertAnswerParams{
		AttemptID:        pgAttemptID,
		QuestionID:       pgQuestionID,
		SelectedOptionID: pgOptionID,
		EssayAnswer:      toPgText(*essayAnswer),
	})
	if err != nil {
		return err
	}

	if selectedOptionID != nil {
		_ = s.q.AutoGradeMultipleChoice(ctx, answer.ID)
	}

	if final {
		hasEssay, err := s.q.HasUngradedEssay(ctx, pgAttemptID)
		if err != nil {
			return err
		}

		status := "graded"
		if hasEssay {
			status = "waiting_assessment"
		}

		return s.q.UpdateAttemptStatus(ctx, db.UpdateAttemptStatusParams{
			ID:     pgAttemptID,
			Status: status,
		})
	}

	return nil
}

func (s *QuizService) GetResult(
	ctx context.Context,
	attemptID uuid.UUID,
) (*ResultDTO, error) {

	pgAttemptID := toPgUUID(attemptID)

	attempt, err := s.q.GetAttemptByID(ctx, pgAttemptID)
	if err != nil {
		return nil, errors.New("attempt not found")
	}

	// if attempt.Status == "in_progress" {
	// 	return nil, errors.New("attempt not submitted")
	// }

	rows, err := s.q.ListAnswersWithQuestions(ctx, pgAttemptID)
	if err != nil {
		return nil, err
	}

	result := &ResultDTO{
		Status: attempt.Status,
	}

	var autoScore int32
	var manualScore int32

	for _, r := range rows {
		var score *int32
		if r.Score.Valid {
			v := r.Score.Int32
			score = &v
		}

		dto := AnswerResultDTO{
			QuestionID: uuid.UUID(r.QuestionID.Bytes),
			Question:   r.QuestionText,
			Type:       r.Type,
			MaxScore:   r.Points,
			Score:      score,
		}

		if r.Type == "multiple_choice" {
			if r.SelectedOptionText.Valid {
				dto.UserAnswer = r.SelectedOptionText.String
			}
			if r.IsCorrect.Bool {
				autoScore += r.Points
				if r.SelectedOptionText.Valid {
					c := r.SelectedOptionText.String
					dto.CorrectAnswer = &c
				}

			}
		} else {
			dto.UserAnswer = r.EssayAnswer.String
			if r.Score.Valid {
				manualScore += r.Score.Int32
			}

		}

		result.Answers = append(result.Answers, dto)
	}

	result.Score.Auto = autoScore

	if attempt.Status == "graded" {
		result.Score.Manual = &manualScore
		final := autoScore + manualScore
		result.Score.Final = &final
	} else {
		result.Score.Manual = nil
		result.Score.Final = nil
	}

	return result, nil
}

func (s *QuizService) loadQuestions(
	ctx context.Context,
	pgQuizID pgtype.UUID,
) ([]QuestionDTO, error) {
	questionsDB, err := s.q.ListQuestionsByQuizID(ctx, pgQuizID)
	if err != nil {
		return nil, err
	}

	var result []QuestionDTO

	for _, q := range questionsDB {
		dto := QuestionDTO{
			ID:     uuid.UUID(q.ID.Bytes),
			Type:   q.Type,
			Text:   q.QuestionText,
			Points: q.Points,
		}

		if q.Type == "multiple_choice" {
			optsDB, err := s.q.ListOptionsByQuestionID(ctx, q.ID)
			if err != nil {
				return nil, err
			}

			var opts []OptionDTO
			for _, o := range optsDB {
				opts = append(opts, OptionDTO{
					ID:   uuid.UUID(o.ID.Bytes),
					Text: o.OptionText,
				})
			}

			dto.Options = opts
		} else {
			dto.Options = nil
		}

		result = append(result, dto)
	}

	return result, nil
}
