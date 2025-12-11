package transport

import (
	"healthy_body/internal/models"
	"healthy_body/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MealPlanResponse используется в Swagger как безопасный ответ без gorm.Model
type MealPlanResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	CategoriesID *uint  `json:"categories_id"`
	TotalDays    int    `json:"total_days"`
}

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
		mealPlans.POST("/", h.Create)
		mealPlans.GET("/", h.GetAllMealPlans)
		mealPlans.GET("/:id", h.GetMealPlanByID)
		mealPlans.PATCH("/:id", h.Update)
		mealPlans.DELETE("/:id", h.Delete)
	}
}

// @Summary Create Meal Plan
// @Tags MealPlans
// @Accept json
// @Produce json
// @Param mealPlan body models.CreateMealPlanRequest true "Meal Plan Data"
// @Success 200 {object} MealPlanResponse
// @Failure 400 {object} map[string]string
// @Router /mealPlans/ [post]
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

// @Summary Get All Meal Plans
// @Tags MealPlans
// @Produce json
// @Success 200 {array} MealPlanResponse
// @Failure 400 {object} map[string]string
// @Router /mealPlans/ [get]
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

// @Summary Update Meal Plan
// @Tags MealPlans
// @Accept json
// @Produce json
// @Param id path int true "Meal Plan ID"
// @Param mealPlan body models.UpdateMealPlanRequest true "Update data"
// @Success 200 {object} MealPlanResponse
// @Failure 400 {object} map[string]string
// @Router /mealPlans/{id} [patch]
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

	h.logger.Info("handler: meal plan updated successfully", "id", id)
	c.JSON(http.StatusOK, mealPlan)
}

// @Summary Delete Meal Plan
// @Tags MealPlans
// @Param id path int true "Meal Plan ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /mealPlans/{id} [delete]
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

// @Summary Get Meal Plan By ID
// @Tags MealPlans
// @Produce json
// @Param id path int true "Meal Plan ID"
// @Success 200 {object} MealPlanResponse
// @Failure 400 {object} map[string]string
// @Router /mealPlans/{id} [get]
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
	h.logger.Info("handler: fetch to meal plan successfully")
	c.JSON(http.StatusOK, mealPlan)
}
