package usecases

import (
	"context"
	"errors"
	"babywolf-blog-backend/internal/domain/entities"
	"babywolf-blog-backend/internal/domain/repositories"
	"babywolf-blog-backend/internal/infrastructure/database"
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

func (u *ManageUsersUseCase) UpdateRole(ctx context.Context, id string, newRole string, executorRole string, ip string) error {
	if executorRole != "admin" {
		return errors.New("unauthorized: only admins can update roles")
	}

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.Role = newRole
	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}
	database.LogAuditAction(user.ID, user.Email, "role_changed", "nuevo rol: "+newRole, ip)
	return nil
}

// UpdateUsername permite al administrador editar los datos básicos de otro usuario
// (Practica 11-12, Parte 1: edición de usuarios).
func (u *ManageUsersUseCase) UpdateUsername(ctx context.Context, id, newUsername, executorRole, ip string) error {
	if executorRole != "admin" {
		return errors.New("unauthorized: only admins can edit users")
	}
	if newUsername == "" {
		return errors.New("username is required")
	}

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	oldUsername := user.Username
	user.Username = newUsername
	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}
	database.LogAuditAction(user.ID, user.Email, "user_updated", "username: "+oldUsername+" -> "+newUsername, ip)
	return nil
}

func (u *ManageUsersUseCase) ListUsers(ctx context.Context, executorRole string) ([]*entities.User, error) {
	if executorRole != "admin" {
		return nil, errors.New("unauthorized: only admins can list users")
	}
	return u.userRepo.ListAll(ctx)
}

// CreateUserByAdmin da de alta un usuario directamente en Supabase Auth
// (Practica 11-12, Parte 1: alta de usuarios).
func (u *ManageUsersUseCase) CreateUserByAdmin(ctx context.Context, executorRole, email, password, username, role, ip string) error {
	if executorRole != "admin" {
		return errors.New("unauthorized: only admins can create users")
	}
	if email == "" || password == "" {
		return errors.New("email and password are required")
	}

	newID, err := database.AdminCreateAuthUser(email, password, username)
	if err != nil {
		return err
	}
	database.LogAuditAction(newID, email, "user_created", "usuario "+username+" creado por administrador", ip)
	return nil
}

// SetActive activa o desactiva (eliminación lógica) una cuenta.
func (u *ManageUsersUseCase) SetActive(ctx context.Context, id string, active bool, executorRole, ip string) error {
	if executorRole != "admin" {
		return errors.New("unauthorized: only admins can activate/deactivate users")
	}

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.Active = active
	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}
	action := "user_deactivated"
	if active {
		action = "user_activated"
	}
	database.LogAuditAction(user.ID, user.Email, action, "", ip)
	return nil
}

// AdminResetPassword permite al administrador restablecer la contraseña de un usuario.
func (u *ManageUsersUseCase) AdminResetPassword(ctx context.Context, id, newPassword, executorRole, ip string) error {
	if executorRole != "admin" {
		return errors.New("unauthorized: only admins can reset passwords")
	}

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := database.AdminSetUserPassword(user.SupabaseUID, newPassword); err != nil {
		return err
	}
	database.LogAuditAction(user.ID, user.Email, "password_change", "restablecida por administrador", ip)
	return nil
}
