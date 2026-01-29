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

func (s *QuizService) StartQuiz(ctx context.Context, quizID, userID uuid.UUID) (*db.QuizAttempt, *db.Quiz, time.Time, error) {
	quiz, err := s.q.GetQuizByID(ctx, toPgUUID(quizID))
	if err != nil {
		return nil, nil, time.Time{}, err
	}

	count, err := s.q.CountUserAttempts(ctx, db.CountUserAttemptsParams{
		QuizID: toPgUUID(quizID),
		UserID: toPgUUID(userID),
	})
	if err != nil {
		return nil, nil, time.Time{}, err
	}

	if count >= int64(quiz.MaxAttempts) {
		return nil, nil, time.Time{}, errors.New("max attempts exceeded")
	}

	attempt, err := s.q.CreateQuizAttempt(ctx, db.CreateQuizAttemptParams{
		QuizID:        toPgUUID(quizID),
		UserID:        toPgUUID(userID),
		AttemptNumber: int32(count + 1),
	})
	if err != nil {
		return nil, nil, time.Time{}, err
	}

	deadline := attempt.StartedAt.Time.Add(time.Second * time.Duration(quiz.DurationSeconds))
	return &attempt, &quiz, deadline, nil
}
