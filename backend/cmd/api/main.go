package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"babywolf-blog-backend/internal/application/usecases"
	"babywolf-blog-backend/internal/infrastructure/database"
	"babywolf-blog-backend/internal/infrastructure/http/middleware"
	"babywolf-blog-backend/internal/infrastructure/repositories"
	"babywolf-blog-backend/internal/interfaces/controllers"
)

func main() {
	// 1. Setup Database
	database.InitSupabase() // Used for remote storage if needed
	database.InitGORM()

	// 2. Setup Repositories
	userRepo := repositories.NewGormUserRepository(database.DB)
	postRepo := repositories.NewGormPostRepository(database.DB)
	// commentRepo would go here...

	// 3. Setup UseCases
	manageUsersUseCase := usecases.NewManageUsersUseCase(userRepo)
	managePostsUseCase := usecases.NewManagePostsUseCase(postRepo)
	storageUseCase := usecases.NewStorageUseCase()

	// 4. Setup Controllers
	userController := controllers.NewUserController(manageUsersUseCase)
	postController := controllers.NewPostController(managePostsUseCase)
	storageController := controllers.NewStorageController(storageUseCase)

	// 5. Setup Gin Router
	r := gin.Default()
	r.Use(gin.Logger(), gin.Recovery())

	// Public Routes
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Babywolf Blog API is running"})
		})
		
		api.GET("/posts", postController.ListPublicPosts)
		api.GET("/posts/:slug", postController.GetPostBySlug)
	}

	// Protected Routes
	protected := api.Group("/admin")
	protected.Use(middleware.SupabaseAuthMiddleware())
	{
		protected.GET("/profile", userController.GetProfile)
		protected.POST("/posts", postController.CreatePost)
		protected.DELETE("/posts/:id", postController.DeletePost)
        protected.POST("/upload", storageController.UploadFile)

		// Practica 11-12: administracion de usuarios (el rol real se valida
		// contra la BD dentro de cada handler, no solo con el JWT).
		protected.GET("/users", userController.ListUsers)
		protected.POST("/users", userController.CreateUser)
		protected.PATCH("/users/:id", userController.UpdateUser)
		protected.PATCH("/users/:id/active", userController.SetActive)
		protected.PATCH("/users/:id/password", userController.ResetPassword)
		protected.PATCH("/users/:id/role", userController.UpdateRole)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
