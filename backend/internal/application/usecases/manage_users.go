package usecases

import (
	"context"
	"errors"
	"babywolf-blog-backend/internal/domain/entities"
	"babywolf-blog-backend/internal/domain/repositories"
)

type ManageUsersUseCase struct {
	userRepo repositories.UserRepository
}

func NewManageUsersUseCase(repo repositories.UserRepository) *ManageUsersUseCase {
	return &ManageUsersUseCase{userRepo: repo}
}

func (u *ManageUsersUseCase) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u *ManageUsersUseCase) GetUserBySupabaseUID(ctx context.Context, supabaseUID string) (*entities.User, error) {
	return u.userRepo.GetBySupabaseUID(ctx, supabaseUID)
}

func (u *ManageUsersUseCase) RegisterUser(ctx context.Context, user *entities.User) error {
	if user.Email == "" || user.SupabaseUID == "" {
		return errors.New("email and supabase_uid are required")
	}
	if user.Role == "" {
		user.Role = "user" // default fallback
	}
	return u.userRepo.Create(ctx, user)
}

func (u *ManageUsersUseCase) UpdateRole(ctx context.Context, id string, newRole string, executorRole string) error {
	if executorRole != "admin" {
		return errors.New("unauthorized: only admins can update roles")
	}

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.Role = newRole
	return u.userRepo.Update(ctx, user)
}
