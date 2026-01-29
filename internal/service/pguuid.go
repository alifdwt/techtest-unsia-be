package service

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func toPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}

func toPgText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}
