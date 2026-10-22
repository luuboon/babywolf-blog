package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"babywolf-blog-backend/internal/application/connect"
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
	soporteRepo := repositories.NewGormSoporteRepository(database.DB)
	// commentRepo would go here...

	// 3. Setup UseCases
	manageUsersUseCase := usecases.NewManageUsersUseCase(userRepo)
	managePostsUseCase := usecases.NewManagePostsUseCase(postRepo)
	storageUseCase := usecases.NewStorageUseCase()

	// 4. Setup Controllers
	userController := controllers.NewUserController(manageUsersUseCase)
	postController := controllers.NewPostController(managePostsUseCase)
	storageController := controllers.NewStorageController(storageUseCase)

	// Desarrollo móvil (app BabyWolf): prácticas 3, 6 y 7.
	consultaController := controllers.NewConsultaController(managePostsUseCase)
	connectController := controllers.NewConnectController(connect.NewStore())
	soporteController := controllers.NewSoporteController(soporteRepo, manageUsersUseCase)

	// 5. Setup Gin Router
	r := gin.Default()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	// Public Routes
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Babywolf Blog API is running"})
		})
		
		api.GET("/posts", postController.ListPublicPosts)
		api.GET("/posts/:slug", postController.GetPostBySlug)

		// Práctica 3: URL preasignada (firma HMAC, 5 min).
		api.GET("/consulta/url", consultaController.URL)
		api.GET("/consulta/recurso", consultaController.Recurso)

		// Práctica 6: MobileConnect, canal efímero entre dos dispositivos.
		api.GET("/connect/salas/:sala/mensajes", connectController.Listar)
		api.POST("/connect/salas/:sala/mensajes", connectController.Enviar)
		api.DELETE("/connect/salas/:sala/mensajes", connectController.Limpiar)

		// Práctica 7: unirse a la mesa de ayuda con la URL preasignada de rol.
		api.GET("/soporte/unirse", soporteController.InfoInvitacion)
		api.POST("/soporte/unirse", soporteController.Unirse)
	}

	// Práctica 7: mesa de ayuda. Requiere el token que entrega /unirse o
	// /admin/soporte/sesion.
	sop := api.Group("/soporte")
	sop.Use(soporteController.Auth())
	{
		sop.POST("/latido", soporteController.Latido)
		sop.POST("/salir", soporteController.Salir)
		sop.GET("/monitor", soporteController.Monitor)
		sop.POST("/invitaciones", soporteController.Invitar)
		sop.DELETE("/participantes/:id", soporteController.Expulsar)

		sop.GET("/reportes", soporteController.ListarReportes)
		sop.POST("/reportes", soporteController.CrearReporte)
		sop.POST("/reportes/asumir-desconectados", soporteController.AsumirDesconectados)
		sop.GET("/reportes/:id", soporteController.VerReporte)
		sop.PATCH("/reportes/:id", soporteController.ActualizarReporte)
		sop.DELETE("/reportes/:id", soporteController.EliminarReporte)
		sop.POST("/reportes/:id/asignar", soporteController.Asignar)

		sop.GET("/mensajes", soporteController.ListarMensajes)
		sop.POST("/mensajes", soporteController.EnviarMensaje)
		sop.DELETE("/mensajes", soporteController.LimpiarMensajes)
		sop.POST("/archivos/subida", soporteController.FirmarSubida)
		sop.POST("/archivos/descarga", soporteController.FirmarDescarga)
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

		// Práctica 7: el admin del blog abre su sesión en la mesa de ayuda.
		protected.POST("/soporte/sesion", soporteController.SesionAdmin)
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
