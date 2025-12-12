package transport

import (
	"healthy_body/internal/models"
	"healthy_body/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Для Swagger лучше использовать отдельную структуру ответа
type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

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
	group := r.Group("/category")
	{
		group.POST("/", h.CreateCategory)
		group.GET("/", h.GetList)
		group.GET("/:id", h.GetByID)
		group.PATCH("/:id", h.UpdateCategory)
		group.DELETE("/:id", h.DeleteCategory)
	}
}

// CreateCategory godoc
// @Summary Создать категорию
// @Description Создает новую категорию
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body models.CreateCategoryRequest true "Данные категории"
// @Success 201 {object} CategoryResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /category/ [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Warn("invalid input", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := h.category.CreateCategory(input)
	if err != nil {
		h.log.Error("failed to create category", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

// GetByID godoc
// @Summary Получить категорию по ID
// @Description Возвращает категорию по ID
// @Tags Categories
// @Produce json
// @Param id path int true "ID категории"
// @Success 200 {object} CategoryResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /category/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Warn("invalid id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	cat, err := h.category.GetCategoryByID(uint(id))
	if err != nil {
		h.log.Error("category not found", "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, cat)
}

// GetList godoc
// @Summary Получить список категорий
// @Description Возвращает все категории
// @Tags Categories
// @Produce json
// @Success 200 {array} CategoryResponse
// @Failure 500 {object} map[string]string
// @Router /category/ [get]
func (h *CategoryHandler) GetList(c *gin.Context) {
	list, err := h.category.GetCategoryList()
	if err != nil {
		h.log.Error("failed to get category list", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get categories"})
		return
	}

	c.JSON(http.StatusOK, list)
}

// UpdateCategory godoc
// @Summary Обновить категорию
// @Description Обновляет категорию по ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "ID категории"
// @Param category body models.UpdateCategoryRequest true "Данные обновления"
// @Success 200 {object} CategoryResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /category/{id} [patch]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Warn("invalid id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Warn("invalid update data", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := h.category.UpdateCategory(uint(id), input)
	if err != nil {
		h.log.Error("failed to update category", "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found or update failed"})
		return
	}

	c.JSON(http.StatusOK, cat)
}

// DeleteCategory godoc
// @Summary Удалить категорию
// @Description Удаляет категорию по ID
// @Tags Categories
// @Produce json
// @Param id path int true "ID категории"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /category/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Warn("invalid id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.category.DeleteCategory(uint(id)); err != nil {
		h.log.Error("failed to delete category", "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found or delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
