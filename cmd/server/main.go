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
connect database")
}

// Migrate Schema
if err := model.Migrate(db); err != nil {
migrate database")
}

// Initialize Router
r := gin.Default()

// CORS Configuration
r.Use(cors.New(cors.Config{
s:     []string{"*"}, // Allow all for simplicity in trial, or specific domains
   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
   []string{"Origin", "Content-Type", "Accept"},
  []string{"Content-Length"},
tials: true,
         12 * time.Hour,
}))

// Setup Handlers
expenseHandler := handler.NewExpenseHandler(db)
projectHandler := handler.NewProjectHandler(db)

api := r.Group("/api")
{
ses", expenseHandler.CreateExpense)
ses", expenseHandler.GetExpenses)
expenseHandler.GetSummary)
seHandler.SetBudget)
dler.CreateProject)
dler.GetProjects)
}

port := os.Getenv("PORT")
if port == "" {
"8002"
}

log.Println("Server starting on :" + port)
r.Run(":" + port)
}
