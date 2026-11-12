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

// Create inserta la noticia. Práctica 12 (defecto P12-D01): si el ID viene
// vacío se omite la columna para que la base genere el UUID por defecto;
// antes se enviaba id = '' y PostgreSQL lo rechazaba (SQLSTATE 22P02), de modo
// que POST /api/admin/posts devolvía 500 en el 100 % de los casos.
func (r *GormPostRepository) Create(ctx context.Context, post *entities.Post) error {
	tx := r.db.WithContext(ctx)
	if post.ID == "" {
		tx = tx.Omit("id")
	}
	return tx.Create(post).Error
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

// ListPublishedPage: Práctica 12. Una página de publicados con LIMIT/OFFSET
// (usa idx_posts_publicados_recientes, 009) y el total para X-Total-Count.
// created_at puede repetirse; id desempata para que las páginas no se solapen.
func (r *GormPostRepository) ListPublishedPage(ctx context.Context, limit, offset int) ([]*entities.Post, int64, error) {
	var total int64
	base := r.db.WithContext(ctx).Model(&entities.Post{}).Where("published = ?", true)
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var posts []*entities.Post
	err := r.db.WithContext(ctx).
		Where("published = ?", true).
		Order("created_at desc").Order("id").
		Limit(limit).Offset(offset).
		Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}
