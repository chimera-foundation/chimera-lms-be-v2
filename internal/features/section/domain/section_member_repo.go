package domain

import (
	"context"
)

type SectionMemberRepository interface {
	Create(ctx context.Context, member *SectionMember) error
}
