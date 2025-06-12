package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"

	_ "healthcare-portal/docs"
	"healthcare-portal/internal/config"
	"healthcare-portal/internal/database"
	"healthcare-portal/internal/handlers"
	"healthcare-portal/internal/middleware"
	"healthcare-portal/internal/repository"
	"healthcare-portal/internal/services"
)

// @title Healthcare Portal API
// @version 1.0
// @description API for Healthcare Portal with Receptionist and Doctor portals
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database
	database.Initialize()
	db := database.GetDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	patientRepo := repository.NewPatientRepository(db)
	appointmentRepo := repository.NewAppointmentRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	patientService := services.NewPatientService(patientRepo)
	appointmentService := services.NewAppointmentService(appointmentRepo, patientRepo, userRepo)

	// Initialize handlers - Now using services instead of repositories
	authHandler := handlers.NewAuthHandler(authService)
	patientHandler := handlers.NewPatientHandler(patientService)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)

	// Setup router
	router := setupRouter(authHandler, patientHandler, appointmentHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port // fallback to config/default (e.g., 8080) for local dev
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	router.Run(":" + port) // Start the server on the specified port
}

func setupRouter(authHandler *handlers.AuthHandler, patientHandler *handlers.PatientHandler, appointmentHandler *handlers.AppointmentHandler) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggerMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Healthcare Portal API is running",
		})
	})

	router.GET("/swagger/*any", ginSwagger.CustomWrapHandler(
		&ginSwagger.Config{
			URL: "/swagger/doc.json", // generated path
		},
		swaggerFiles.Handler,
	))

	// API routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
		}

		// Patient routes
		patients := api.Group("/patients")
		patients.Use(middleware.AuthMiddleware())
		{
			patients.GET("", patientHandler.GetAllPatients)
			patients.GET("/search", patientHandler.SearchPatients)
			patients.GET("/:id", patientHandler.GetPatientByID)

			// Receptionist only routes
			patients.POST("", middleware.RoleMiddleware("receptionist"), patientHandler.CreatePatient)
			patients.DELETE("/:id", middleware.RoleMiddleware("receptionist"), patientHandler.DeletePatient)

			// Both receptionist and doctor can update
			patients.PUT("/:id", middleware.RoleMiddleware("receptionist", "doctor"), patientHandler.UpdatePatient)
		}

		// Appointment routes
		appointments := api.Group("/appointments")
		appointments.Use(middleware.AuthMiddleware())
		{
			appointments.GET("", appointmentHandler.GetAllAppointments)
			appointments.GET("/date", appointmentHandler.GetAppointmentsByDate)
			appointments.GET("/:id", appointmentHandler.GetAppointmentByID)
			appointments.GET("/patient/:patientId", appointmentHandler.GetPatientAppointments)
			appointments.GET("/doctor/:doctorId", appointmentHandler.GetDoctorAppointments)

			// Receptionist only routes
			appointments.POST("", middleware.RoleMiddleware("receptionist"), appointmentHandler.CreateAppointment)
			appointments.DELETE("/:id", middleware.RoleMiddleware("receptionist"), appointmentHandler.DeleteAppointment)

			// Both can update status
			appointments.PATCH("/:id/status", middleware.RoleMiddleware("receptionist", "doctor"), appointmentHandler.UpdateAppointmentStatus)
		}
	}

	return router
}
