package main

import (
	"fmt"
	"log"

	"alhikmah-attendance-api/config"
	"alhikmah-attendance-api/internal/handler"
	"alhikmah-attendance-api/internal/middleware"
	"alhikmah-attendance-api/internal/repository"
	"alhikmah-attendance-api/internal/service"
	"alhikmah-attendance-api/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// 2. Connect to Database
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	// 3. Run Migrations
	runDBMigration(cfg)

	// 4. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	classRepo := repository.NewClassRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)

	// 5. Initialize Services
	userService := service.NewUserService(userRepo)
	classService := service.NewClassService(classRepo)
	studentService := service.NewStudentService(studentRepo)
	attendanceService := service.NewAttendanceService(attendanceRepo, studentRepo)

	// 6. Initialize Handlers
	authHandler := &handler.AuthHandler{
		DB:     db,
		Config: cfg,
	}
	userHandler := handler.NewUserHandler(userService)
	classHandler := handler.NewClassHandler(classService)
	studentHandler := handler.NewStudentHandler(studentService)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)

	// 7. Setup Gin Router
	r := gin.Default()

	// Global Middlewares (CORS could go here)

	// API Routes Group
	api := r.Group("/api/v1")
	{
		// Health Check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "UP"})
		})

		// Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", middleware.AuthMiddleware(cfg), authHandler.Logout)
		}

		// Protected Routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Users
			protected.GET("/users/me", userHandler.GetMe)
			protected.GET("/users", middleware.RoleMiddleware("admin"), userHandler.GetAll)
			protected.POST("/users", middleware.RoleMiddleware("admin"), userHandler.Create)
			protected.PUT("/users/:user_id", userHandler.Update)

			// Classes
			protected.GET("/classes", classHandler.GetAll)
			protected.GET("/classes/:class_id", classHandler.GetByID)

			// Students
			protected.POST("/students", middleware.RoleMiddleware("admin"), studentHandler.Create)
			protected.GET("/classes/:class_id/students", studentHandler.GetByClass)

			// Attendance
			protected.POST("/attendances/scan", attendanceHandler.ScanQR)
			protected.POST("/attendances/manual", attendanceHandler.ManualInput)
			protected.GET("/classes/:class_id/attendances/today", attendanceHandler.GetClassAttendanceForToday)
		}
	}

	// 8. Start Server

	log.Printf("Starting server on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runDBMigration(cfg config.Config) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations applied successfully")
}
