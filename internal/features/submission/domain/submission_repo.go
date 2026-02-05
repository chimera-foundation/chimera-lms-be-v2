package domain

import (
	"context"
)

type SubmissionRepository interface {
	Create(ctx context.Context, submission *Submission) error
}
