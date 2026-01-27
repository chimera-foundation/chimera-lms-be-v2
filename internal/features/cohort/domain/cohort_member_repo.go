package domain

import (
	"context"
)

type CohortMemberRepository interface {
	Create(ctx context.Context, member *CohortMember) error
}
