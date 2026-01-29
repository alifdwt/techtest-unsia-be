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
