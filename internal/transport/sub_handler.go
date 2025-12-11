package transport

import (
	"healthy_body/internal/models"
	"healthy_body/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SubscriptionResponse используется в Swagger как безопасный ответ без gorm.Model
type SubscriptionResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Price        int    `json:"price"`
	DurationDays int    `json:"duration_days"`
	CategoriesID uint   `json:"categories_id"`
}

type SubscriptionHandler struct {
	sub service.SubscriptionService
	log *slog.Logger
}

func NewSubscriptionHandler(sub service.SubscriptionService, log *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		sub: sub,
		log: log,
	}
}

func (h *SubscriptionHandler) RegisterRoutes(r *gin.Engine) {
	subGroup := r.Group("/sub")
	{
		subGroup.POST("/", h.CreateSub)
		subGroup.GET("/", h.GetListSub)
		subGroup.GET("/:id", h.GetByID)
		subGroup.PATCH("/:id", h.Update)
		subGroup.DELETE("/:id", h.Delete)
	}
}

// CreateSub godoc
// @Summary Создать подписку
// @Description Создает новую подписку
// @Tags Subscription
// @Accept json
// @Produce json
// @Param subscription body models.CreateSubscriptionRequest true "Данные для создания подписки"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sub/ [post]
func (h *SubscriptionHandler) CreateSub(r *gin.Context) {
	var inputSub models.CreateSubscriptionRequest
	if err := r.ShouldBindJSON(&inputSub); err != nil {
		h.log.Warn("error create values invalid")
		r.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	sub, err := h.sub.CreateSub(&inputSub)
	if err != nil {
		h.log.Warn("error create in db")
		r.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	h.log.Info("created sub succes")
	r.IndentedJSON(http.StatusOK, sub)
}

// GetByID godoc
// @Summary Получить подписку по ID
// @Description Возвращает подписку по указанному ID
// @Tags Subscription
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sub/{id} [get]
func (h *SubscriptionHandler) GetByID(r *gin.Context) {
	idStr := r.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("error parse id")
		r.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sub, err := h.sub.GetSubByID(uint(id))
	if err != nil {
		h.log.Warn("error create in db")
		r.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	h.log.Info("sub finded", "sub_id", sub.ID)
	r.IndentedJSON(http.StatusOK, sub)
}

// GetListSub godoc
// @Summary Получить список подписок
// @Description Возвращает все подписки
// @Tags Subscription
// @Produce json
// @Success 200 {array} SubscriptionResponse
// @Failure 500 {object} map[string]string
// @Router /sub/ [get]
func (h *SubscriptionHandler) GetListSub(r *gin.Context) {
	list, err := h.sub.GetListSub()
	if err != nil {
		h.log.Error("error list not found in db")
		r.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	h.log.Info("sub list finded")
	r.IndentedJSON(http.StatusOK, list)
}

// Update godoc
// @Summary Обновить подписку
// @Description Обновляет подписку по ID
// @Tags Subscription
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Param subscription body models.UpdateSubscriptionRequest true "Данные для обновления подписки"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Router /sub/{id} [patch]
func (h *SubscriptionHandler) Update(r *gin.Context) {
	idStr := r.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("error parse id")
		r.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var upSub models.UpdateSubscriptionRequest
	if err := r.ShouldBindJSON(&upSub); err != nil {
		h.log.Warn("error type update values")
		r.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	sub, err := h.sub.UpdateSub(uint(id), upSub)
	if err != nil {
		h.log.Error("error type update values")
		r.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	h.log.Info("sub updated succes", "sub_id", sub.ID)
	r.IndentedJSON(http.StatusOK, sub)
}

// Delete godoc
// @Summary Удалить подписку
// @Description Удаляет подписку по ID
// @Tags subscription
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Router /sub/{id} [delete]
func (h *SubscriptionHandler) Delete(r *gin.Context) {
	idStr := r.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("error parse id")
		r.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.sub.Delete(uint(id)); err != nil {
		h.log.Error("error delete sub")
		r.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	h.log.Info("sub deleted succes", "sub_id", id)
	r.IndentedJSON(http.StatusOK, gin.H{"deleted": true})
}
