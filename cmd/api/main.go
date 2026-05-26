package main

import (
	"fmt"
	"log"
	"time"

	"alhikmah-attendance-api/config"
	"alhikmah-attendance-api/internal/handler"
	"alhikmah-attendance-api/internal/middleware"
	"alhikmah-attendance-api/internal/repository"
	"alhikmah-attendance-api/internal/service"
	"alhikmah-attendance-api/pkg/database"

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
	reportRepo := repository.NewReportRepository(db)
	reportCacheRepo := repository.NewReportCacheRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)

	// 5. Initialize Services
	userService := service.NewUserService(userRepo)
	classService := service.NewClassService(classRepo)
	studentService := service.NewStudentService(studentRepo)
	attendanceService := service.NewAttendanceService(attendanceRepo, studentRepo)
	reportService := service.NewReportService(reportRepo, studentRepo, reportCacheRepo)
	dashboardService := service.NewDashboardService(dashboardRepo)

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
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	// 7. Setup Gin Router
	r := gin.Default()

	// Global Middlewares
	frontendURL := cfg.FrontendURL
	if frontendURL == "" {
		frontendURL = "http://localhost:3000" // Default for local development
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
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
			protected.PUT("/users/:user_id",middleware.RoleMiddleware("admin"), userHandler.Update)
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

			// Dashboard
			protected.GET("/dashboard/recent-activity", dashboardHandler.GetRecentActivity)
			protected.GET("/dashboard/attendance-trend", dashboardHandler.GetAttendanceTrend)

			// Reports
			protected.GET("/reports/daily", reportHandler.GetDailyReport)
			protected.GET("/reports/monthly", reportHandler.GetMonthlyReport)
			protected.GET("/reports/semester", reportHandler.GetSemesterReport)
			protected.GET("/reports/student/:student_id", reportHandler.GetStudentReport)
			protected.POST("/reports/export", reportHandler.Export)
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
