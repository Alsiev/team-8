package transport

import (
	"healthy_body/internal/models"
	"healthy_body/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MealPlanHandler struct {
	mealPlans service.MealPlanService
	logger    *slog.Logger
}

func NewMealPlanHandler(mealPlans service.MealPlanService, logger *slog.Logger) *MealPlanHandler {
	return &MealPlanHandler{
		mealPlans: mealPlans,
		logger:    logger,
	}
}

func (h *MealPlanHandler) RegisterRoutes(r *gin.Engine) {
	mealPlans := r.Group("/mealPlans")
	{
		mealPlans.POST("", h.Create)
	}
}

func (h *MealPlanHandler) Create(c *gin.Context) {
	var req models.CreateMealPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler: failed to bind JSON", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mealPlan, err := h.mealPlans.CreateMealPlan(req)
	if err != nil {
		h.logger.Error("handler: failed to create meal plan", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("handler: meal plan created successfully")
	c.JSON(http.StatusOK, mealPlan)
}

func (h *MealPlanHandler) GetAllMealPlans(c *gin.Context) {
	mealPlans, err := h.mealPlans.ListMealPlan()
	if err != nil {
		h.logger.Error("failed to fetch meal plans")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("fetch to meal plans successfully", "count", len(mealPlans))
	c.JSON(http.StatusOK, mealPlans)
}

func (h *MealPlanHandler) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error("handler: invalid meal plan id", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}

	var req models.UpdateMealPlanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("handler: failed to bind update request", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mealPlan, err := h.mealPlans.UpdateMealPlan(uint(id), &req)
	if err != nil {
		h.logger.Error("handler: failed to update meal plan")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("hadnler: meal plan updated successfully", "id", id)
	c.JSON(http.StatusOK, mealPlan)
}

func (h *MealPlanHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error("handler: invalid meal plan id", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id"})
		return
	}

	if err := h.mealPlans.DeleteMealPlan(uint(id)); err != nil {
		h.logger.Error("handler: failed to delete meal plan", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка при удалении плана"})
		return
	}
	h.logger.Info("handler: meal plan deleted successfully", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "план успешно удалён"})
}

func (h *MealPlanHandler) GetMealPlanByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error("handler: invalid meal plan id", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mealPlan, err := h.mealPlans.GetMealPlanByID(uint(id))
	if err != nil {
		h.logger.Error("handler: failed to fetch meal plan", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("hadnler: fetch to meal plan successfully")
	c.JSON(http.StatusOK, mealPlan)
}
