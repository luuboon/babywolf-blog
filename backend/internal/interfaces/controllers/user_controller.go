package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"babywolf-blog-backend/internal/application/usecases"
	"babywolf-blog-backend/internal/domain/entities"
)

var errUnauthorized = errors.New("unauthorized")

type UserController struct {
	manageUsersUseCase *usecases.ManageUsersUseCase
}

func NewUserController(useCase *usecases.ManageUsersUseCase) *UserController {
	return &UserController{manageUsersUseCase: useCase}
}

func (c *UserController) GetProfile(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := c.manageUsersUseCase.GetUserBySupabaseUID(ctx.Request.Context(), userId.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// executorRole resuelve el rol real (desde la BD, fuente de verdad) del usuario autenticado.
func (c *UserController) executorRole(ctx *gin.Context) (string, error) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		return "", errUnauthorized
	}
	user, err := c.manageUsersUseCase.GetUserBySupabaseUID(ctx.Request.Context(), userId.(string))
	if err != nil {
		return "", err
	}
	return user.Role, nil
}

func (c *UserController) ListUsers(ctx *gin.Context) {
	role, err := c.executorRole(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	users, err := c.manageUsersUseCase.ListUsers(ctx.Request.Context(), role)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Username string `json:"username" binding:"required"`
	Role     string `json:"role"`
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	role, err := c.executorRole(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.manageUsersUseCase.CreateUserByAdmin(ctx.Request.Context(), role, req.Email, req.Password, req.Username, req.Role, ctx.ClientIP()); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

type setActiveRequest struct {
	Active bool `json:"active"`
}

func (c *UserController) SetActive(ctx *gin.Context) {
	role, err := c.executorRole(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req setActiveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.manageUsersUseCase.SetActive(ctx.Request.Context(), ctx.Param("id"), req.Active, role, ctx.ClientIP()); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User status updated"})
}

type resetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

func (c *UserController) ResetPassword(ctx *gin.Context) {
	role, err := c.executorRole(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req resetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.manageUsersUseCase.AdminResetPassword(ctx.Request.Context(), ctx.Param("id"), req.Password, role, ctx.ClientIP()); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

type updateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user admin"`
}

func (c *UserController) UpdateRole(ctx *gin.Context) {
	role, err := c.executorRole(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req updateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.manageUsersUseCase.UpdateRole(ctx.Request.Context(), ctx.Param("id"), req.Role, role, ctx.ClientIP()); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Role updated"})
}

type updateUserRequest struct {
	Username string `json:"username" binding:"required"`
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	role, err := c.executorRole(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.manageUsersUseCase.UpdateUsername(ctx.Request.Context(), ctx.Param("id"), req.Username, role, ctx.ClientIP()); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func (c *UserController) Register(ctx *gin.Context) {
	var user entities.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.manageUsersUseCase.RegisterUser(ctx.Request.Context(), &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
