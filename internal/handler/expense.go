package handler

import (
"construction-tracker/internal/model"
"net/http"
"strconv"

"github.com/gin-gonic/gin"
"gorm.io/gorm"
)

type ExpenseHandler struct {
DB *gorm.DB
}

func NewExpenseHandler(db *gorm.DB) *ExpenseHandler {
return &ExpenseHandler{DB: db}
}

// CreateExpense records a new daily expense
func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
var expense model.Expense
if err := c.ShouldBindJSON(&expense); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

// Validate ProjectID
if expense.ProjectID == 0 {
c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
return
}

if result := h.DB.Create(&expense); result.Error != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
return
}

c.JSON(http.StatusCreated, expense)
}

// GetExpenses returns list of expenses
func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
projectID := c.Query("project_id")
if projectID == "" {
c.JSON(http.StatusBadRequest, gin.H{"error": "project_id query param is required"})
return
}

var expenses []model.Expense
if result := h.DB.Where("project_id = ?", projectID).Order("date desc").Find(&expenses); result.Error != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
return
}

c.JSON(http.StatusOK, expenses)
}

// SetBudget sets the budget for a project
func (h *ExpenseHandler) SetBudget(c *gin.Context) {
var budget model.Budget
if err := c.ShouldBindJSON(&budget); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
    
    if budget.ProjectID == 0 {
c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
return
    }

// Upsert based on ProjectID
var existing model.Budget
result := h.DB.Where("project_id = ?", budget.ProjectID).First(&existing)
if result.Error == nil {
existing.Amount = budget.Amount
h.DB.Save(&existing)
} else {
h.DB.Create(&budget)
}

c.JSON(http.StatusOK, budget)
}

// GetSummary returns total cost, breakdown, and comparison to budget
func (h *ExpenseHandler) GetSummary(c *gin.Context) {
projectIDStr := c.Query("project_id")
    if projectIDStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "project_id query param is required"})
        return
    }
    
    // safe ignore error as we check empty above, but good practice to check
    projectID, _ := strconv.Atoi(projectIDStr)

var total float64
var breakdown []struct {
Type  model.ExpenseType `json:"type"`
Total float64           `json:"total"`
}

// Calculate overall total
h.DB.Model(&model.Expense{}).Where("project_id = ?", projectID).Select("coalesce(sum(amount), 0)").Scan(&total)

// Calculate breakdown by type
h.DB.Model(&model.Expense{}).Where("project_id = ?", projectID).Select("type, coalesce(sum(amount), 0) as total").Group("type").Scan(&breakdown)

// Get Budget
var budget model.Budget
h.DB.Where("project_id = ?", projectID).First(&budget)

alert := false
if budget.Amount > 0 && total >= budget.Amount {
alert = true
}

c.JSON(http.StatusOK, gin.H{
"total_cost": total,
"breakdown":  breakdown,
"budget":     budget.Amount,
"alert":      alert,
})
}
