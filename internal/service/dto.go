package service

import "github.com/google/uuid"

type QuestionDTO struct {
	ID      uuid.UUID
	Type    string
	Text    string
	Points  int32
	Options []OptionDTO
}

type OptionDTO struct {
	ID   uuid.UUID
	Text string
}
