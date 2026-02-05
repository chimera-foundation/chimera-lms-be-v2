package domain

import (
	"context"
)

type AssessmentRepo interface {
	Create(ctx context.Context, assessment *Assessment) error
}