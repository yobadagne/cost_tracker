package handler

import (
"construction-tracker/internal/model"
"net/http"

"github.com/gin-gonic/gin"
"gorm.io/gorm"
)

type ProjectHandler struct {
DB *gorm.DB
}

func NewProjectHandler(db *gorm.DB) *ProjectHandler {
return &ProjectHandler{DB: db}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
var project model.Project
if err := c.ShouldBindJSON(&project); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

if result := h.DB.Create(&project); result.Error != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
return
}

c.JSON(http.StatusCreated, project)
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {
var projects []model.Project
if result := h.DB.Find(&projects); result.Error != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
return
}
c.JSON(http.StatusOK, projects)
}
