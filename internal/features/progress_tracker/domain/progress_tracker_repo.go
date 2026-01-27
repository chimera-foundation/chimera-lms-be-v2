package domain

import (
	"context"
)

type ProgressTrackerRepository interface {
	Create(ctx context.Context, tracker *ProgressTracker) error
	Update(ctx context.Context, tracker *ProgressTracker) error
}
