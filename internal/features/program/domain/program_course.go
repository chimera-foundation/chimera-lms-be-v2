package domain

import (
	"github.com/google/uuid"
)

type ProgramCourse struct {
	ProgramID  uuid.UUID
	CourseID   uuid.UUID
	OrderIndex int
}
