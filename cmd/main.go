package main

import (
	"log"
	"os"
	"time"

	"crm-backend/internal/config"
	"crm-backend/internal/database"
	"crm-backend/internal/handlers"
	"crm-backend/internal/middleware"
	"crm-backend/internal/repositories"
	"crm-backend/internal/services"
	"crm-backend/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Inicializar logger
	logger.Init()
	logger.Info("Iniciando aplicação CRM Backend")

	// Carregar configurações
	cfg := config.Load()
	logger.Infof("Configurações carregadas - Environment: %s", cfg.Environment)

	// Conectar ao banco de dados
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Falha ao conectar com o banco de dados:", err)
	}
	logger.Info("Conexão com banco de dados estabelecida")

	// Executar migrações
	if err := database.Migrate(db); err != nil {
		logger.Fatal("Falha ao executar migrações:", err)
	}
	logger.Info("Migrações executadas com sucesso")

	// Inicializar repositórios
	userRepo := repositories.NewUserRepository(db)
	contactRepo := repositories.NewContactRepository(db)
	interactionRepo := repositories.NewInteractionRepository(db)
	taskRepo := repositories.NewTaskRepository(db)
	projectRepo := repositories.NewProjectRepository(db)

	// Inicializar serviços
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	userService := services.NewUserService(userRepo, contactRepo, taskRepo, projectRepo, interactionRepo)
	contactService := services.NewContactService(contactRepo, interactionRepo, taskRepo, projectRepo)
	interactionService := services.NewInteractionService(interactionRepo, contactRepo)
	taskService := services.NewTaskService(taskRepo, contactRepo, projectRepo)
	projectService := services.NewProjectService(projectRepo, contactRepo, taskRepo)

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	contactHandler := handlers.NewContactHandler(contactService)
	interactionHandler := handlers.NewInteractionHandler(interactionService)
	taskHandler := handlers.NewTaskHandler(taskService)
	projectHandler := handlers.NewProjectHandler(projectService)

	// Configurar Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))

	// Middleware global
	router.Use(middleware.CustomLogger()) // Usar o logger personalizado
	router.Use(middleware.ErrorHandler())

	logger.Info("Middlewares configurados")

	// Agrupar todas as rotas sob /api
	api := router.Group("/api")
	{
		// Rotas públicas
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/validate", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.ValidateToken)
			auth.POST("/logout", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.Logout)
		}

		// Rotas protegidas (agora como subgrupo de /api)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Rotas de usuários
			users := protected.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
				users.PUT("/change-password", userHandler.ChangePassword)
				users.DELETE("/delete-account", userHandler.DeleteAccount)
				users.GET("/stats", userHandler.GetStats)
				users.GET("/activities", userHandler.GetRecentActivities)
			}

			// Rotas de contatos
			contacts := protected.Group("/contacts")
			{
				contacts.POST("/create", contactHandler.Create)
				contacts.GET("/list", contactHandler.List)
				contacts.GET("/:id", contactHandler.GetByID)
				contacts.PUT("/:id", contactHandler.Update)
				contacts.DELETE("/:id", contactHandler.Delete)

				contacts.POST("/:id/interactions", interactionHandler.Create)
				contacts.GET("/:id/interactions", interactionHandler.ListByContact)
			}

			// Rotas de tarefas
			tasks := protected.Group("/tasks")
			{
				tasks.POST("/create", taskHandler.Create)
				tasks.GET("/list", taskHandler.List)
				tasks.GET("/:id", taskHandler.GetByID)
				tasks.PUT("/:id", taskHandler.Update)
				tasks.DELETE("/:id", taskHandler.Delete)
				tasks.PUT("/:id/complete", taskHandler.MarkTaskAsCompleted)
				tasks.PUT("/:id/uncomplete", taskHandler.MarkTaskAsPending)
			}

			// Rotas de projetos
			projects := protected.Group("/projects")
			{
				projects.POST("/create", projectHandler.Create)
				projects.GET("/list", projectHandler.List)
				projects.GET("/:id", projectHandler.GetByID)
				projects.PUT("/:id", projectHandler.Update)
				projects.DELETE("/:id", projectHandler.Delete)
			}

			// Rotas de interações (globais)
			interactions := protected.Group("/interactions")
			{
				interactions.GET("/list", interactionHandler.List)
				interactions.GET("/:id", interactionHandler.GetByID)
				interactions.PUT("/:id", interactionHandler.Update)
				interactions.DELETE("/:id", interactionHandler.Delete)
			}
		}

		// Iniciar servidor
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		logger.Infof("Servidor iniciando na porta %s", port)
		logger.WithFields("INFO", "Server Starting", map[string]interface{}{
			"port":        port,
			"environment": cfg.Environment,
			"address":     "0.0.0.0:" + port,
		})

		if err := router.Run("0.0.0.0:" + port); err != nil {
			logger.Fatal("Falha ao iniciar servidor:", err)
		}
	}
}
