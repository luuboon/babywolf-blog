package repositories

import (
	"context"
	"babywolf-blog-backend/internal/domain/entities"
)

// CommentRepository defines the interface for comment data operations
type CommentRepository interface {
	GetByID(ctx context.Context, id string) (*entities.Comment, error)
	GetByPostID(ctx context.Context, postID string) ([]*entities.Comment, error)
	Create(ctx context.Context, comment *entities.Comment) error
	Update(ctx context.Context, comment *entities.Comment) error
	Delete(ctx context.Context, id string) error
	DeleteAllByPostID(ctx context.Context, postID string) error
}
