package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"babywolf-blog-backend/internal/application/usecases"
	"babywolf-blog-backend/internal/domain/entities"
)

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
