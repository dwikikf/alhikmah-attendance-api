package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"alhikmah-attendance-api/config"
	"alhikmah-attendance-api/core/handler"
	"alhikmah-attendance-api/core/middleware"
	"alhikmah-attendance-api/core/repository"
	"alhikmah-attendance-api/core/service"
	"alhikmah-attendance-api/pkg/cache"
	"alhikmah-attendance-api/pkg/database"
	"alhikmah-attendance-api/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		slog.Error("Could not load config", slog.Any("error", err))
		os.Exit(1)
	}

	// 1.5 Initialize Logger
	logger.InitLogger(cfg.AppEnv)
	slog.Info("Logger initialized", "env", cfg.AppEnv)

	// 2. Connect to Database
	db, err := database.ConnectDB(cfg)
	if err != nil {
		slog.Error("Could not connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	// 3. Run Migrations
	runDBMigration(cfg)

	// 4. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	classRepo := repository.NewClassRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)
	reportRepo := repository.NewReportRepository(db)
	reportCacheRepo := repository.NewReportCacheRepository(db)

	// 5. Initialize Services
	userService := service.NewUserService(userRepo)
	classService := service.NewClassService(classRepo)
	studentService := service.NewStudentService(studentRepo)

	qrCache := cache.NewCache()
	attendanceService := service.NewAttendanceService(attendanceRepo, studentRepo, classRepo, qrCache)
	reportService := service.NewReportService(reportRepo, studentRepo, reportCacheRepo)

	// 6. Initialize Handlers
	authHandler := &handler.AuthHandler{
		DB:     db,
		Config: cfg,
	}
	userHandler := handler.NewUserHandler(userService)
	classHandler := handler.NewClassHandler(classService, studentService)
	studentHandler := handler.NewStudentHandler(studentService, classService)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)
	reportHandler := handler.NewReportHandler(reportService)

	// 7. Setup Gin Router
	r := gin.Default()

	// Global Middlewares
	r.Use(middleware.LoggerMiddleware())

	frontendURL := cfg.FrontendURL
	if frontendURL == "" {
		frontendURL = "http://localhost:3000" // Default for local development
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "sentry-trace", "baggage"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
			auth.POST("/refresh-token", authHandler.Refresh)
			auth.GET("/me", middleware.AuthMiddleware(cfg), userHandler.GetMe)

			// Password Reset Stubs
			auth.POST("/reset-password", authHandler.ResetPasswordRequest)
			auth.POST("/reset-password/confirm", authHandler.ResetPasswordConfirm)
			auth.POST("/reset-password/change", middleware.AuthMiddleware(cfg), authHandler.ResetPasswordChange)
		}

		// Protected Routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Users
			protected.GET("/users/me", userHandler.GetMe)
			protected.GET("/users", middleware.RoleMiddleware("admin"), userHandler.GetAll)
			protected.POST("/users", middleware.RoleMiddleware("admin"), userHandler.Create)
			protected.PUT("/users/:user_id", middleware.RoleMiddleware("admin"), userHandler.Update)
			protected.DELETE("/users/:user_id", middleware.RoleMiddleware("admin"), userHandler.Delete)

			// Classes
			protected.GET("/classes", classHandler.GetAll)
			protected.GET("/classes/:class_id", classHandler.GetByID)
			protected.GET("/classes/:class_id/export-excel", classHandler.ExportExcel)
			protected.GET("/classes/:class_id/export-qrcode", classHandler.ExportQRCode)
			protected.POST("/classes", middleware.RoleMiddleware("admin", "teacher"), classHandler.Create)
			protected.PUT("/classes/:class_id", middleware.RoleMiddleware("admin", "teacher"), classHandler.Update)
			protected.DELETE("/classes/:class_id", middleware.RoleMiddleware("admin", "teacher"), classHandler.Delete)

			// Students
			protected.GET("/students", studentHandler.GetAll)
			protected.GET("/students/:student_id", studentHandler.GetByID)
			protected.POST("/students", studentHandler.Create)
			protected.POST("/students/import", studentHandler.ImportCSV)
			protected.PUT("/students/:student_id", studentHandler.Update)
			protected.DELETE("/students/:student_id", studentHandler.Delete)
			protected.GET("/classes/:class_id/students", studentHandler.GetByClass)
			protected.GET("/students/:student_id/qrcode", studentHandler.GetQRCode)

			// Attendance
			protected.POST("/attendances/qr-scan", attendanceHandler.ScanQR)
			protected.POST("/attendances/manual", attendanceHandler.ManualInput)
			protected.PUT("/attendances/:attendance_id", attendanceHandler.Update)
			protected.GET("/classes/:class_id/attendances/today", attendanceHandler.GetClassAttendanceForToday)
			protected.GET("/attendances/:class_id/:date", attendanceHandler.GetByClassAndDate)

			// Reports
			protected.GET("/reports/daily", reportHandler.GetDailyReport)
			protected.GET("/reports/monthly", reportHandler.GetMonthlyReport)
			protected.GET("/reports/semester", reportHandler.GetSemesterReport)
			protected.GET("/reports/student/:student_id", reportHandler.GetStudentReport)
			protected.POST("/reports/export", reportHandler.Export)
		}
	}

	// 8. Start Server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	slog.Info("Starting server", "port", port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}

func runDBMigration(cfg config.Config) {
	var dsn string
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		dsn = dbURL
	} else {
		sslmode := cfg.DBSSLMode
		if sslmode == "" {
			if cfg.AppEnv == "production" {
				sslmode = "require"
			} else {
				sslmode = "disable"
			}
		}
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, sslmode)
	}

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		slog.Error("Failed to initialize migrations", slog.Any("error", err))
		os.Exit(1)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("Failed to run migrations", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Database migrations applied successfully")
}
