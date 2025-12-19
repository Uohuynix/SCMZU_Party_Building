package main

import (
	"SCMZU_Party_Building/api"
	"SCMZU_Party_Building/config"
	"SCMZU_Party_Building/dao"
	"SCMZU_Party_Building/db"
	"SCMZU_Party_Building/middleware"
	"SCMZU_Party_Building/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration
	cfg := config.LoadConfigFromEnv()
	log.Printf("Starting Party Building Management System...")
	log.Printf("Database config: %s:%s@%s:%s/%s", cfg.DBUser, "****", cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Initialize database
	gormDB, sqlDB, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer sqlDB.Close()

	log.Println("? Database connected successfully")

	// Check if database has data, if not create tables and sample data
	var tableCount int
	err = sqlDB.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ?", cfg.DBName).Scan(&tableCount)
	if err != nil {
		log.Printf("Warning: Unable to check database tables: %v", err)
	} else if tableCount == 0 {
		log.Println("Empty database detected, creating tables and sample data...")
		if err := db.CreateTables(sqlDB); err != nil {
			log.Printf("Warning: Failed to create tables: %v", err)
		} else {
			log.Println("? Database tables created successfully")
		}

		if err := db.InsertSampleData(sqlDB); err != nil {
			log.Printf("Warning: Failed to insert sample data: %v", err)
		} else {
			log.Println("? Sample data inserted successfully")
		}
	}

	// Initialize DAO
	userDAO := dao.NewUserDAO(gormDB)
	studentDAO := dao.NewStudentDAO(gormDB)
	trainingDAO := dao.NewTrainingDAO(gormDB)
	developmentDAO := dao.NewDevelopmentDAO(gormDB)
	rewardDAO := dao.NewRewardDAO(gormDB)

	// Initialize Service
	authService := service.NewAuthService(userDAO)
	studentService := service.NewStudentService(studentDAO, developmentDAO, rewardDAO)
	trainingService := service.NewTrainingService(trainingDAO)
	statsService := service.NewStatsService(trainingDAO)

	// Initialize Handler
	authHandler := api.NewAuthHandler(authService)
	// 学员相关接口需要根据用户角色（学生/教师/管理员）做权限控制，因此注入 authService
	studentHandler := api.NewStudentHandler(studentService, authService)
	trainingHandler := api.NewTrainingHandler(trainingService)
	statsHandler := api.NewStatsHandler(statsService)

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(gin.Recovery())

	// Static files for uploaded photos
	router.Static("/uploads", "./uploads")

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Party Building Management System is running",
		})
	})

	// API v1 route group
	v1 := router.Group("/api/v1")
	{
		// Authentication routes (no auth required)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}

		// Routes requiring authentication
		authenticated := v1.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			// User profile
			authenticated.GET("/profile", authHandler.Profile)

			// Student management
			students := authenticated.Group("/students")
			{
				students.GET("", studentHandler.ListStudents)
				students.GET("/:id", studentHandler.GetStudentDetail)
				students.POST("", studentHandler.CreateStudent)
				students.PUT("/:id", studentHandler.UpdateStudent)
				students.DELETE("/:id", studentHandler.DeleteStudent)
				students.POST("/:id/photo", studentHandler.UploadPhoto)

				// Student development history
				students.GET("/:id/development", studentHandler.GetDevelopment)
				students.PUT("/:id/development", studentHandler.UpdateDevelopment)

				// Student reward records
				students.GET("/:id/rewards", studentHandler.GetRewards)
				students.POST("/:id/rewards", studentHandler.CreateReward)
			}

			// Training management
			trainings := authenticated.Group("/trainings")
			{
				// 所有角色都可以查看
				trainings.GET("", trainingHandler.ListTrainings)
				trainings.GET("/:id", trainingHandler.GetTraining)

				// 只有教师和管理员可以编辑（创建、更新、删除）
				trainings.POST("", middleware.TeacherOrAdminMiddleware(), trainingHandler.CreateTraining)
				trainings.PUT("/:id", middleware.TeacherOrAdminMiddleware(), trainingHandler.UpdateTraining)
				trainings.DELETE("/:id", middleware.TeacherOrAdminMiddleware(), trainingHandler.DeleteTraining)
			}

			// Statistics
			stats := authenticated.Group("/statistics")
			{
				stats.GET("/trainings", statsHandler.TrainingStats)
			}
		}
	}

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server started successfully")
	log.Printf("Access URL: http://localhost%s", serverAddr)
	log.Printf("API docs: http://localhost%s/api/v1", serverAddr)
	log.Printf("Health check: http://localhost%s/health", serverAddr)
	log.Println("\nPress Ctrl+C to stop the server")

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
