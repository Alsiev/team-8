package transport

import (
	"healthy_body/internal/models"
	"healthy_body/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ExercisePlanResponse используется в Swagger как безопасный ответ без gorm.Model
type ExercisePlanResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	DurationWeeks int    `json:"duration_weeks"`
	CategoriesID  uint   `json:"categories_id"`
}

type ExercisePlanHandler struct {
	exer service.ExercisePlanServices
	log  *slog.Logger
}

func NewExercisePlanHandler(exer service.ExercisePlanServices, log *slog.Logger) *ExercisePlanHandler {
	return &ExercisePlanHandler{
		exer: exer,
		log:  log,
	}
}

func (h *ExercisePlanHandler) RegisterRoutes(r *gin.Engine) {
	planGroup := r.Group("/plan")
	{
		planGroup.POST("/", h.CreatePlan)
		planGroup.GET("/:id", h.GetByID)
		planGroup.GET("/", h.GetAllPlan)
		planGroup.PATCH("/:id", h.UpdatePlan)
		planGroup.DELETE("/:id", h.DeletePlan)

		planGroup.POST("/planItem", h.CreatePlanItem)
		planGroup.GET("/planItem/:id", h.GetPlanItemByID)
		planGroup.GET("/planItem/", h.GetListPlanItem)
		planGroup.PATCH("/planItem/:id", h.UpdatePlanItem)
		planGroup.DELETE("/planItem/:id", h.DeletePlanItem)
	}
}

// CreatePlan godoc
// @Summary Создание тренировочного плана
// @Description Создаёт новый тренировочный план
// @Tags ExercisePlan
// @Accept json
// @Produce json
// @Param plan body models.CreateExercesicePlanRequest true "Данные тренировочного плана"
// @Success 200 {object} ExercisePlanResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /plan/ [post]
func (h *ExercisePlanHandler) CreatePlan(c *gin.Context) {
	var inputPlan models.CreateExercesicePlanRequest

	if err := c.ShouldBindJSON(&inputPlan); err != nil {
		h.log.Warn("error invalid input type information")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.exer.CreatePlan(inputPlan)
	if err != nil {
		h.log.Error("error in db")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("succes create plan", "plan", plan)
	c.IndentedJSON(http.StatusOK, plan)
}

// GetByID godoc
// @Summary Получить тренировочный план по ID
// @Tags ExercisePlan
// @Produce json
// @Param id path int true "ID плана"
// @Success 200 {object} ExercisePlanResponse
// @Failure 400 {object} map[string]string
// @Router /plan/{id} [get]
func (h *ExercisePlanHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("error parse id")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	plan, err := h.exer.GetPlanByID(uint(id))
	if err != nil {
		h.log.Error("error found plan in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.log.Info("success plan found", "plan_id", plan.ID)
	c.IndentedJSON(http.StatusOK, plan)
}

// GetAllPlan godoc
// @Summary Получить список тренировочных планов
// @Tags ExercisePlan
// @Produce json
// @Success 200 {array} ExercisePlanResponse
// @Failure 400 {object} map[string]string
// @Router /plan/ [get]
func (h *ExercisePlanHandler) GetAllPlan(c *gin.Context) {
	list, err := h.exer.GetListPlans()
	if err != nil {
		h.log.Error("error found plan list in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid list"})
		return
	}

	h.log.Info("list found success")
	c.IndentedJSON(http.StatusOK, list)
}

// UpdatePlan godoc
// @Summary Обновить тренировочный план
// @Tags ExercisePlan
// @Accept json
// @Produce json
// @Param id path int true "ID плана"
// @Param plan body models.UpdateExercesicePlanRequest true "Обновлённые данные"
// @Success 200 {object} ExercisePlanResponse
// @Failure 400 {object} map[string]string
// @Router /plan/{id} [patch]
func (h *ExercisePlanHandler) UpdatePlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("error parse id")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var updatePlan models.UpdateExercesicePlanRequest

	if err := c.ShouldBindJSON(&updatePlan); err != nil {
		h.log.Warn("error type update values")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	plan, err := h.exer.UpdatePlan(uint(id), updatePlan)
	if err != nil {
		h.log.Error("error update plan in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id or update"})
		return
	}

	h.log.Info("success plan updated", "plan_id", plan.ID)
	c.IndentedJSON(http.StatusOK, plan)
}

// DeletePlan godoc
// @Summary Удалить тренировочный план
// @Tags ExercisePlan
// @Produce json
// @Param id path int true "ID плана"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Router /plan/{id} [delete]
func (h *ExercisePlanHandler) DeletePlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("error parse id")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.exer.DeletePlan(uint(id)); err != nil {
		h.log.Error("error delete plan in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id or update"})
		return
	}

	h.log.Info("success plan deleted")
	c.IndentedJSON(http.StatusOK, gin.H{"deleted": true})
}

// CreatePlanItem godoc
// @Summary Создать упражнение в плане
// @Tags ExercisePlanItem
// @Accept json
// @Produce json
// @Param item body models.CreateExercisePlanItemRequest true "Данные элемента плана"
// @Success 200 {object} models.ExercisePlanItem
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /plan/planItem [post]
func (h *ExercisePlanHandler) CreatePlanItem(c *gin.Context) {
	var inputPlanItem models.CreateExercisePlanItemRequest

	if err := c.ShouldBindJSON(&inputPlanItem); err != nil {
		h.log.Warn("error invalid input type information")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.exer.CreatePlanItem(inputPlanItem)
	if err != nil {
		h.log.Error("error in db")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("success create plan item", "planItem", plan)
	c.IndentedJSON(http.StatusOK, plan)
}

// GetPlanItemByID godoc
// @Summary Получить элемент плана по ID
// @Tags ExercisePlanItem
// @Produce json
// @Param id path int true "ID элемента плана"
// @Success 200 {object} models.ExercisePlanItem
// @Failure 400 {object} map[string]string
// @Router /plan/planItem/{id} [get]
func (h *ExercisePlanHandler) GetPlanItemByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("error parse id")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	plan, err := h.exer.GetByIDPlanItem(uint(id))
	if err != nil {
		h.log.Error("error found planItem in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.log.Info("success planItem found", "planItem_id", plan.ID)
	c.IndentedJSON(http.StatusOK, plan)
}

// GetListPlanItem godoc
// @Summary Получить список элементов плана
// @Tags ExercisePlanItem
// @Produce json
// @Success 200 {array} models.ExercisePlanItem
// @Failure 400 {object} map[string]string
// @Router /plan/planItem/ [get]
func (h *ExercisePlanHandler) GetListPlanItem(c *gin.Context) {
	list, err := h.exer.GetAllPlanItem()
	if err != nil {
		h.log.Error("error found planItem list in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid list"})
		return
	}

	h.log.Info("success list found")
	c.IndentedJSON(http.StatusOK, list)
}

// UpdatePlanItem godoc
// @Summary Обновить элемент плана
// @Tags ExercisePlanItem
// @Accept json
// @Produce json
// @Param id path int true "ID элемента"
// @Param item body models.UpdateExercisePlanItemRequest true "Обновление"
// @Success 200 {object} models.ExercisePlanItem
// @Failure 400 {object} map[string]string
// @Router /plan/planItem/{id} [patch]
func (h *ExercisePlanHandler) UpdatePlanItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("error parse id")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var updatePlan models.UpdateExercisePlanItemRequest

	if err := c.ShouldBindJSON(&updatePlan); err != nil {
		h.log.Warn("error type update values")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	plan, err := h.exer.UpdatePlanItem(uint(id), updatePlan)
	if err != nil {
		h.log.Error("error update planItem in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id or update"})
		return
	}

	h.log.Info("success planItem updated", "planItem_id", plan.ID)
	c.IndentedJSON(http.StatusOK, plan)
}

// DeletePlanItem godoc
// @Summary Удалить элемент плана
// @Tags ExercisePlanItem
// @Produce json
// @Param id path int true "ID элемента"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Router /plan/planItem/{id} [delete]
func (h *ExercisePlanHandler) DeletePlanItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("error parse id")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.exer.DeletePlanItem(uint(id)); err != nil {
		h.log.Error("error delete planItem in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id or update"})
		return
	}

	h.log.Info("success planItem deleted")
	c.IndentedJSON(http.StatusOK, gin.H{"deleted": true})
}
