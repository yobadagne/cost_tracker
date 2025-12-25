package main

import (
	"construction-tracker/internal/handler"
	"construction-tracker/internal/model"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialize Database
	// Note: For production with Render Free Tier, SQLite is ephemeral (data lost on restart).
	// For a real app, use PostgreSQL. Keeping SQLite as per plan for simple trial.
	db, err := gorm.Open(sqlite.Open("construction_tracker.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}

	// Migrate Schema
	if err := model.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database")
	}

	// Initialize Router
	r := gin.Default()

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all for simplicity in trial, or specific domains
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Content-Length"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup Handlers
	expenseHandler := handler.NewExpenseHandler(db)
	projectHandler := handler.NewProjectHandler(db)

	api := r.Group("/api")
	{
		api.POST("/expenses", expenseHandler.CreateExpense)
		api.GET("/expenses", expenseHandler.GetExpenses)
		api.GET("/summary", expenseHandler.GetSummary)
		api.POST("/budget", expenseHandler.SetBudget)
		api.POST("/projects", projectHandler.CreateProject)
		api.GET("/projects", projectHandler.GetProjects)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	log.Println("Server starting on :" + port)
	r.Run(":" + port)
}
