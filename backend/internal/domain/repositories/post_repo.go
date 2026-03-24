package repositories

import (
	"context"
	"babywolf-blog-backend/internal/domain/entities"
)

// PostRepository defines the interface for post data operations
type PostRepository interface {
	GetByID(ctx context.Context, id string) (*entities.Post, error)
	GetBySlug(ctx context.Context, slug string) (*entities.Post, error)
	Create(ctx context.Context, post *entities.Post) error
	Update(ctx context.Context, post *entities.Post) error
	Delete(ctx context.Context, id string) error
	ListAll(ctx context.Context, includeUnpublished bool) ([]*entities.Post, error)
}
