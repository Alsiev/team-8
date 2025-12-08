package transport

import (
	"healthy_body/internal/models"
	"healthy_body/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	category service.CategoryServices
	log      *slog.Logger
}

func NewCategoryHandler(category service.CategoryServices, log *slog.Logger) *CategoryHandler {
	return &CategoryHandler{
		category: category,
		log:      log,
	}
}

func (h *CategoryHandler) RegisterRoutes(r *gin.Engine) {
	categoryGroup := r.Group("/category")
	{
		categoryGroup.POST("/", h.CreateCategory)
		categoryGroup.GET("/:id", h.GetByID)
		categoryGroup.GET("/", h.GetList)
		categoryGroup.PATCH("/:id", h.UpdateCategory)
		categoryGroup.DELETE("/:id", h.DeleteCategory)

	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var inputCategory models.CreateCategoryRequest

	if err := c.ShouldBindJSON(&inputCategory); err != nil {
		h.log.Warn("error invalid input type information")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.category.CreateCategory(inputCategory)
	if err != nil {
		h.log.Error("error in db")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("succes create category",
		"category", category,
	)
	c.IndentedJSON(http.StatusCreated, category)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("error parse id")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	category, err := h.category.GetCategoryByID(uint(id))
	if err != nil {
		h.log.Error("error found category in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.log.Info("succes category found",
		"category_id", category.ID,
	)
	c.IndentedJSON(http.StatusOK, category)
}

func (h *CategoryHandler) GetList(c *gin.Context) {
	list, err := h.category.GetCategoryList()
	if err != nil {
		h.log.Error("error found category list in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid list"})
		return
	}

	h.log.Info("list found succes")
	c.IndentedJSON(http.StatusOK, list)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("error parse id")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var updateCategory models.UpdateCategoryRequest

	if err := c.ShouldBindJSON(&updateCategory); err != nil {
		h.log.Warn("error type update values")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	category, err := h.category.UpdateCategory(uint(id), updateCategory)
	if err != nil {
		h.log.Error("error update category in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id or update"})
		return
	}

	h.log.Info("succes category updated",
		"category_id", category.ID,
	)
	c.IndentedJSON(http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("error parse id")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.category.DeleteCategory(uint(id)); err != nil {
		h.log.Error("error delete category in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id or update"})
		return
	}

	h.log.Info("succes category deleted")
	c.IndentedJSON(http.StatusOK, gin.H{"deleted": true})
}
