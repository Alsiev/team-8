package transport

import (
	"healthy_body/internal/models"
	"healthy_body/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func (h *ExercisePlanHandler) RegisterRoutes(r *gin.Engine){
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

	h.log.Info("succes create plan",
		"plan", plan,
	)
	c.IndentedJSON(http.StatusOK, plan)
}

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
		h.log.Error("error found category in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.log.Info("succes plan found",
		"plan_id", plan.ID,
	)
	c.IndentedJSON(http.StatusOK, plan)
}

func (h *ExercisePlanHandler) GetAllPlan(c *gin.Context) {
	list, err := h.exer.GetListPlans()
	if err != nil {
		h.log.Error("error found plan list in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid list"})
		return
	}

	h.log.Info("list found succes")
	c.IndentedJSON(http.StatusOK, list)
}

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
		h.log.Error("error update category in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id or update"})
		return
	}

	h.log.Info("succes plan updated",
		"plan_id", plan.ID,
	)
	c.IndentedJSON(http.StatusOK, plan)
}

func (h *ExercisePlanHandler) DeletePlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("error parse id")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.exer.DeletePlan(uint(id)); err != nil {
		h.log.Error("error delete category in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id or update"})
		return
	}

	h.log.Info("succes plan deleted")
	c.IndentedJSON(http.StatusOK, gin.H{"deleted": true})
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

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

	h.log.Info("succes create plan item",
		"plan item", plan,
	)
	c.IndentedJSON(http.StatusOK, plan)
}

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

	h.log.Info("succes planItem found",
		"planItem_id", plan.ID,
	)
	c.IndentedJSON(http.StatusOK, plan)
}

func (h *ExercisePlanHandler) GetListPlanItem(c *gin.Context) {
	list, err := h.exer.GetAllPlanItem()
	if err != nil {
		h.log.Error("error found planItem list in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid list"})
		return
	}

	h.log.Info("list found succes")
	c.IndentedJSON(http.StatusOK, list)
}

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

	h.log.Info("succes planItem updated",
		"planItem_id", plan.ID,
	)
	c.IndentedJSON(http.StatusOK, plan)
}

func (h *ExercisePlanHandler) DeletePlanItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("error parse id")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.exer.DeletePlanItem(uint(id)); err != nil {
		h.log.Error("error delete category in db")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id or update"})
		return
	}

	h.log.Info("succes category deleted")
	c.IndentedJSON(http.StatusOK, gin.H{"deleted": true})
}
