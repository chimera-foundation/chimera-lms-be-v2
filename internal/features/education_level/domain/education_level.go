package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/subject/domain"
)

type EducationLevel struct {
	shared.Base
	
	Name string
	Code string	
	Subjects []domain.Subject
}