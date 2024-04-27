package util

import (
	"fmt"

	"github.com/google/uuid"
)

func ValidateUUID(id string) (uuid.UUID, error) {
	newID, err := uuid.Parse(id)
	if err != nil || newID == uuid.Nil || newID.String() == "" {
		return uuid.Nil, fmt.Errorf("invalid UUID: %w", err)
	}
	return newID, nil
}
