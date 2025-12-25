package model

import (
"time"

"gorm.io/gorm"
)

type ExpenseType string

const (
Labor     ExpenseType = "Labor"
Material  ExpenseType = "Material"
Equipment ExpenseType = "Equipment"
)

type Project struct {
ID        uint      `gorm:"primaryKey" json:"id"`
Name      string    `json:"name"`
CreatedAt time.Time `json:"created_at"`
}

type Expense struct {
ID          uint        `gorm:"primaryKey" json:"id"`
Date        time.Time   `json:"date"`
Description string      `json:"description"`
Type        ExpenseType `json:"type"`
Amount      float64     `json:"amount"`
ProjectID   uint        `json:"project_id"` 
CreatedAt   time.Time   `json:"created_at"`
UpdatedAt   time.Time   `json:"updated_at"`
}

type Budget struct {
ID        uint    `gorm:"primaryKey" json:"id"`
ProjectID uint    `json:"project_id" gorm:"uniqueIndex"`
Amount    float64 `json:"amount"`
}

// Migrate initializes the database schema
func Migrate(db *gorm.DB) error {
return db.AutoMigrate(&Project{}, &Expense{}, &Budget{})
}
