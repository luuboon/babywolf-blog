package repositories

import (
	"context"

	"gorm.io/gorm"
	"babywolf-blog-backend/internal/domain/entities"
)

type GormPostRepository struct {
	db *gorm.DB
}

func NewGormPostRepository(db *gorm.DB) *GormPostRepository {
	return &GormPostRepository{db: db}
}

func (r *GormPostRepository) GetByID(ctx context.Context, id string) (*entities.Post, error) {
	var post entities.Post
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *GormPostRepository) GetBySlug(ctx context.Context, slug string) (*entities.Post, error) {
	var post entities.Post
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *GormPostRepository) Create(ctx context.Context, post *entities.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *GormPostRepository) Update(ctx context.Context, post *entities.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *GormPostRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Post{}).Error
}

func (r *GormPostRepository) ListAll(ctx context.Context, includeUnpublished bool) ([]*entities.Post, error) {
	var posts []*entities.Post
	query := r.db.WithContext(ctx)
	
	if !includeUnpublished {
		query = query.Where("published = ?", true)
	}

	if err := query.Order("created_at desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
