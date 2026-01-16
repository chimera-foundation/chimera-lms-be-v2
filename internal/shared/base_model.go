package shared

import (
	"context"
	"time"
	"github.com/uptrace/bun"
    "github.com/google/uuid"
)

type AuditModel struct {
	IsActive bool `bun:"is_active,notnull,default:true"`
	
	CreatedAt time.Time `bun:"created_at,type:timestamp,nullzero,notnull,default:current_timestamp"`
	CreatedBy uuid.UUID `bun:"created_by,type:uuid,notnull"`
	
	UpdatedAt time.Time `bun:"updated_at,type:timestamp,nullzero,notnull,default:current_timestamp"`
	UpdatedBy uuid.UUID `bun:"updated_by,type:uuid,notnull"`
	
	DeletedAt time.Time `bun:"deleted_at,type:timestamp,soft_delete,nullzero"`
	DeletedBy *uuid.UUID `bun:"deleted_by,type:uuid"` 
}

func (m *AuditModel) BeforeInsert(ctx context.Context, _ *bun.InsertQuery) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	
	if uid, ok := ctx.Value("id").(uuid.UUID); ok {
		m.CreatedBy = uid
		m.UpdatedBy = uid
	}
	return nil
}

func (m *AuditModel) BeforeUpdate(ctx context.Context, _ *bun.UpdateQuery) error {
	m.UpdatedAt = time.Now()
	
	if uid, ok := ctx.Value("id").(uuid.UUID); ok {
		m.UpdatedBy = uid
	}
	return nil
}