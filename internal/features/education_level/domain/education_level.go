package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
)

type EducationLevel struct {
	shared.Base
	
	Name string
	Code string	
}