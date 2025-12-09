package transport

import (
	"healthy_body/internal/models"
	"healthy_body/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MealPlanItemHandler struct {
	mealPlanItems service.MealPlanItemsService
	logger        *slog.Logger
}

func NewMealPlanItemHandler(mealPlanItems service.MealPlanItemsService, logger *slog.Logger) *MealPlanItemHandler {
	return &MealPlanItemHandler{
		mealPlanItems: mealPlanItems,
		logger:        logger,
	}
}

func (h *MealPlanItemHandler) RegisterRoutes(r *gin.Engine) {
	mealPlanItems := r.Group("/mealPlanItems")
	{
		mealPlanItems.POST("/", h.Create)
		mealPlanItems.GET("/", h.ListMealPlanItems)
		mealPlanItems.PATCH("/:id", h.Update)
		mealPlanItems.GET("/:id", h.GetMealPlanItemById)
		mealPlanItems.DELETE("/:id", h.DeleteMealPlanItem)
	}
}

func (h *MealPlanItemHandler) Create(c *gin.Context) {
	var req models.CreateMealPlanItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler: failed to bind JSON", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mealPlanItem, err := h.mealPlanItems.CreateMealPlanItem(req)
	if err != nil {
		h.logger.Error("handler: failed to create meal plan item", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("handler: meal plan item created successfully")
	c.JSON(http.StatusOK, mealPlanItem)
}

func (h *MealPlanItemHandler) ListMealPlanItems(c *gin.Context) {
	mealPlanItems, err := h.mealPlanItems.GetAllMealPlanItems()
	if err != nil {
		h.logger.Error("failed to fetch meal plan items")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("fetch to meal plan items successfully", "count", len(mealPlanItems))
	c.JSON(http.StatusOK, mealPlanItems)
}

func (h *MealPlanItemHandler) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error("handler: invalid meal plan item id", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}

	var req models.UpdateMealPlanItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler: failed to bind update request", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mealPlanItem, err := h.mealPlanItems.UpdateMealPlanItem(uint(id), &req)
	if err != nil {
		h.logger.Error("handler: failed to update meal plan item")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("handler: meal plan item updated successfully", "id", id)
	c.JSON(http.StatusOK, mealPlanItem)
}

func (h *MealPlanItemHandler) GetMealPlanItemById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error("handler: invalid meal plan item id", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}

	mealPlanItem, err := h.mealPlanItems.GetMealPlanItemById(uint(id))
	if err != nil {
		h.logger.Error("handler: failed to fetch meal plan item", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("handler: meal plan item fetch to successfully", "id", id)
	c.JSON(http.StatusOK, mealPlanItem)
}

func (h *MealPlanItemHandler) DeleteMealPlanItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error("handler: invalid meal plan item id", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}

	if err := h.mealPlanItems.DeleteMealPlanItem(uint(id)); err != nil {
		h.logger.Error("handler: failed to delete meal plan item", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка при удалении"})
		return
	}
	h.logger.Info("handler: meal plan item deleted successfully", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "удаление прошло успешно"})
}
