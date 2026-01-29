package service

import "github.com/google/uuid"

type QuestionDTO struct {
	ID      uuid.UUID   `json:"id"`
	Type    string      `json:"type"`
	Text    string      `json:"text"`
	Points  int32       `json:"points"`
	Options []OptionDTO `json:"options"`
}

type OptionDTO struct {
	ID   uuid.UUID `json:"id"`
	Text string    `json:"text"`
}

type AnswerResultDTO struct {
	QuestionID    uuid.UUID `json:"question_id"`
	Question      string    `json:"question"`
	Type          string    `json:"type"`
	UserAnswer    string    `json:"user_answer"`
	CorrectAnswer *string   `json:"correct_answer"`
	Score         *int32    `json:"score"`
	MaxScore      int32     `json:"max_score"`
}

type ResultDTO struct {
	Status string `json:"status"`
	Score  struct {
		Auto   int32  `json:"auto"`
		Manual *int32 `json:"manual"`
		Final  *int32 `json:"final"`
	} `json:"score"`
	Answers []AnswerResultDTO `json:"answers"`
}
