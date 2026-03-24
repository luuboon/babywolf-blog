package repositories

import (
	"context"

	"gorm.io/gorm"
	"babywolf-blog-backend/internal/domain/entities"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetBySupabaseUID(ctx context.Context, supabaseUID string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).Where("supabase_uid = ?", supabaseUID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *GormUserRepository) Update(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *GormUserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.User{}).Error
}

func (r *GormUserRepository) ListAll(ctx context.Context) ([]*entities.User, error) {
	var users []*entities.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
