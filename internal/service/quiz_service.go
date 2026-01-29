package service

import (
	"context"
	"errors"
	"time"

	"github.com/alifdwt/techtest-unsia-be/internal/db"
	"github.com/google/uuid"
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

	count, err := s.q.CountUserAttempts(ctx, db.CountUserAttemptsParams{
		QuizID: pgQuizID,
		UserID: pgUserID,
	})
	if err != nil {
		return nil, nil, time.Time{}, nil, err
	}

	if count >= int64(quiz.MaxAttempts) {
		return nil, nil, time.Time{}, nil, errors.New("max attempts exceeded")
	}

	attempt, err := s.q.CreateQuizAttempt(ctx, db.CreateQuizAttemptParams{
		QuizID:        pgQuizID,
		UserID:        pgUserID,
		AttemptNumber: int32(count + 1),
	})
	if err != nil {
		return nil, nil, time.Time{}, nil, err
	}

	deadline := attempt.StartedAt.Time.Add(time.Second * time.Duration(quiz.DurationSeconds))

	questionsDB, err := s.q.ListQuestionsByQuizID(ctx, pgQuizID)
	if err != nil {
		return nil, nil, time.Time{}, nil, err
	}

	var questions []QuestionDTO

	for _, q := range questionsDB {
		optsDB, err := s.q.ListOptionsByQuestionID(ctx, q.ID)
		if err != nil {
			return nil, nil, time.Time{}, nil, err
		}

		var opts []OptionDTO
		for _, o := range optsDB {
			opts = append(opts, OptionDTO{
				ID:   uuid.UUID(o.ID.Bytes),
				Text: o.OptionText,
			})
		}

		questions = append(questions, QuestionDTO{
			ID:      uuid.UUID(q.ID.Bytes),
			Type:    q.Type,
			Text:    q.QuestionText,
			Points:  q.Points,
			Options: opts,
		})
	}

	return &attempt, &quiz, deadline, questions, nil
}
