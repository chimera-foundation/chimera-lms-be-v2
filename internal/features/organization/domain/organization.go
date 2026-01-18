package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
)

type Organization struct {
	shared.Base
	
	Name string `validate:"required,min=3,max=100"`
	Slug string `validate:"required,alphanum"`
	IsActive bool
}