package repositories

import (
	"context"
	"babywolf-blog-backend/internal/domain/entities"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetBySupabaseUID(ctx context.Context, supabaseUID string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error
	ListAll(ctx context.Context) ([]*entities.User, error)
}
