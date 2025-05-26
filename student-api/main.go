package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/one2n/student-api/model"
	"github.com/one2n/student-api/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupDatabase() (*gorm.DB, error) {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "student_db")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Printf("Successfully connected to database at %s:%s", dbHost, dbPort)
	return db, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func setupRouter(studentService *service.StudentService) *gin.Engine {
	gin.DisableConsoleColor()
	r := gin.Default()

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] | %s | %d | %s | %s | %s | %s | %s | %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.ClientIP,
			param.StatusCode,
			param.Method,
			param.Path,
			param.Request.UserAgent(),
			param.Latency,
			param.ErrorMessage,
			param.Request.Host,
		)
	}))

	r.GET("/health", func(c *gin.Context) {
		log.Printf("Health check requested from %s", c.ClientIP())
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/v1/api")
	{
		v1.POST("/students", func(c *gin.Context) {
			log.Printf("Creating new student - Request from %s", c.ClientIP())
			var student model.Student
			if err := c.ShouldBindJSON(&student); err != nil {
				log.Printf("Invalid student data: %v", err)
				c.JSON(http.StatusBadRequest, model.StudentResponse{
					Success: false,
					Message: err.Error(),
				})
				return
			}

			createdStudent, err := studentService.CreateStudent(&student)
			if err != nil {
				log.Printf("Failed to create student: %v", err)
				c.JSON(http.StatusBadRequest, model.StudentResponse{
					Success: false,
					Message: err.Error(),
				})
				return
			}

			log.Printf("Successfully created student with ID: %s", createdStudent.ID)
			c.JSON(http.StatusCreated, model.StudentResponse{
				Success: true,
				Data:    createdStudent,
			})
		})

		v1.GET("/students", func(c *gin.Context) {
			log.Printf("Fetching all students - Request from %s", c.ClientIP())
			students, err := studentService.GetAllStudents()
			if err != nil {
				log.Printf("Failed to fetch students: %v", err)
				c.JSON(http.StatusInternalServerError, model.StudentResponse{
					Success: false,
					Message: err.Error(),
				})
				return
			}

			log.Printf("Successfully fetched %d students", len(students))
			c.JSON(http.StatusOK, model.StudentResponse{
				Success: true,
				Data:    students,
			})
		})

		v1.GET("/students/:id", func(c *gin.Context) {
			id := c.Param("id")
			log.Printf("Fetching student with ID: %s - Request from %s", id, c.ClientIP())
			student, err := studentService.GetStudentByID(id)
			if err != nil {
				log.Printf("Failed to fetch student %s: %v", id, err)
				c.JSON(http.StatusNotFound, model.StudentResponse{
					Success: false,
					Message: err.Error(),
				})
				return
			}

			log.Printf("Successfully fetched student with ID: %s", id)
			c.JSON(http.StatusOK, model.StudentResponse{
				Success: true,
				Data:    student,
			})
		})

		v1.PUT("/students/:id", func(c *gin.Context) {
			id := c.Param("id")
			log.Printf("Updating student with ID: %s - Request from %s", id, c.ClientIP())
			var student model.Student
			if err := c.ShouldBindJSON(&student); err != nil {
				log.Printf("Invalid student data for update: %v", err)
				c.JSON(http.StatusBadRequest, model.StudentResponse{
					Success: false,
					Message: err.Error(),
				})
				return
			}

			updatedStudent, err := studentService.UpdateStudent(id, &student)
			if err != nil {
				log.Printf("Failed to update student %s: %v", id, err)
				c.JSON(http.StatusNotFound, model.StudentResponse{
					Success: false,
					Message: err.Error(),
				})
				return
			}

			log.Printf("Successfully updated student with ID: %s", id)
			c.JSON(http.StatusOK, model.StudentResponse{
				Success: true,
				Data:    updatedStudent,
			})
		})

		v1.DELETE("/students/:id", func(c *gin.Context) {
			id := c.Param("id")
			log.Printf("Deleting student with ID: %s - Request from %s", id, c.ClientIP())
			err := studentService.DeleteStudent(id)
			if err != nil {
				log.Printf("Failed to delete student %s: %v", id, err)
				c.JSON(http.StatusNotFound, model.StudentResponse{
					Success: false,
					Message: err.Error(),
				})
				return
			}

			log.Printf("Successfully deleted student with ID: %s", id)
			c.JSON(http.StatusOK, model.StudentResponse{
				Success: true,
				Message: "Student deleted successfully",
			})
		})
	}

	return r
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	log.Println("Starting Student API server...")

	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	studentService := service.NewStudentService(db)
	r := setupRouter(studentService)

	port := getEnv("SERVER_PORT", "8080")
	log.Printf("Server is starting on port %s...", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
