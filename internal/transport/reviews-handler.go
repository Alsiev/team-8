package transport

import (
	"healthy_body/internal/models"
	"healthy_body/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewsHandler struct {
	review service.ReviewsService
	log    *slog.Logger
}

func NewReviewsHandler(review service.ReviewsService, log *slog.Logger) *ReviewsHandler {
	return &ReviewsHandler{review: review, log: log}
}

func (h *ReviewsHandler) CreateReview(c *gin.Context) {
	var req models.CreateReviewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Введены неверные данные",
			"err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Неверный формат данных",
			"error":   err.Error()})
		return
	}

	if req.UserID == 0 {
		h.log.Warn("UserID не указан")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "UserID обязателен для создания отзыва",
		})
		return
	}

	reviewID, err := h.review.CreateReview(req, req.UserID)
	if err != nil {
		h.log.Error("Ошибка при создании отзыва",
			"error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "ошибка при создании отзыва",
		})
		return
	}

	h.log.Info("Отзыв создан",
		"review_id", reviewID,
		"user_id", req.UserID)

	c.JSON(http.StatusCreated, gin.H{
		"message":   "отзыв создан",
		"review_id": reviewID,
	})
}

func (h *ReviewsHandler) GetReview(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("Некорректный ID отзыва",
			"id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "некорректный ID отзыва",
		})
		return
	}

	review, err := h.review.GetReview(uint(id))
	if err != nil {
		h.log.Error("Ошибка при получении отзыва",
			"id", id,
			"error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "ошибка при получении отзыва",
		})
		return
	}

	h.log.Info("Отзыв получен",
		"id", id)

	c.JSON(http.StatusOK, review)
}

func (h *ReviewsHandler) GetReviewsByUser(c *gin.Context) {
	userIDStr := c.Param("userID")

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		h.log.Warn("Некорректный ID пользователя",
			"user_id", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "некорректный ID пользователя",
		})
		return
	}

	reviews, err := h.review.GetReviewsByUser(uint(userID))
	if err != nil {
		h.log.Error("Ошибка при получении отзывов пользователя",
			"user_id", userID,
			"error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "ошибка при получении отзывов пользователя",
		})
		return
	}

	h.log.Info("Отзывы пользователя получены",
		"user_id", userID,
		"count", len(reviews))

	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
		"total":   len(reviews),
	})
}

func (h *ReviewsHandler) GetReviewsByCategory(c *gin.Context) {
	categoryIDStr := c.Param("categoryID")

	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		h.log.Warn("Некорректный ID категории",
			"category_id", categoryIDStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "некорректный ID категории",
		})
		return
	}

	reviews, err := h.review.GetReviewsByCategory(uint(categoryID))
	if err != nil {
		h.log.Error("Ошибка при получении отзывов по категории",
			"category_id", categoryID,
			"error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "ошибка при получении отзывов по категории",
		})
		return
	}

	h.log.Info("Отзывы по категории получены",
		"category_id", categoryID,
		"count", len(reviews))

	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
		"total":   len(reviews),
	})
}

func (h *ReviewsHandler) UpdateReview(c *gin.Context) {
	var req models.UpdateReviewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Введены неверные данные",
			"err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Неверный формат данных",
			"error":   err.Error()})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("Некорректный ID отзыва",
			"id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "некорректный ID отзыва",
		})
		return
	}

	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		var body struct {
			UserID uint `json:"user_id"`
		}
		if err := c.ShouldBindJSON(&body); err == nil {
			userIDStr = strconv.FormatUint(uint64(body.UserID), 10)
		}
	}

	if userIDStr == "" {
		h.log.Warn("UserID не указан")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "UserID обязателен для обновления отзыва",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		h.log.Warn("Некорректный UserID",
			"user_id", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "некорректный UserID",
		})
		return
	}

	err = h.review.UpdateReview(uint(id), req, uint(userID))
	if err != nil {
		h.log.Error("Ошибка при обновлении отзыва",
			"id", id,
			"user_id", userID,
			"error", err.Error())

		status := http.StatusInternalServerError
		if err.Error() == "нельзя обновлять чужой отзыв" {
			status = http.StatusForbidden
		}

		c.JSON(status, gin.H{
			"error":   err.Error(),
			"message": "ошибка при обновлении отзыва",
		})
		return
	}

	h.log.Info("Отзыв обновлен",
		"id", id,
		"user_id", userID)

	c.JSON(http.StatusOK, gin.H{
		"message": "отзыв успешно обновлен",
	})
}

func (h *ReviewsHandler) DeleteReview(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("Некорректный ID отзыва",
			"id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "некорректный ID отзыва",
		})
		return
	}

	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		h.log.Warn("UserID не указан")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "UserID обязателен для удаления отзыва",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		h.log.Warn("Некорректный UserID",
			"user_id", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "некорректный UserID",
		})
		return
	}

	err = h.review.DeleteReview(uint(id), uint(userID))
	if err != nil {
		h.log.Error("Ошибка при удалении отзыва",
			"id", id,
			"user_id", userID,
			"error", err.Error())

		status := http.StatusInternalServerError
		if err.Error() == "нельзя удалять чужой отзыв" {
			status = http.StatusForbidden
		}

		c.JSON(status, gin.H{
			"error":   err.Error(),
			"message": "ошибка при удалении отзыва",
		})
		return
	}

	h.log.Info("Отзыв удален",
		"id", id,
		"user_id", userID)

	c.JSON(http.StatusOK, gin.H{
		"message": "отзыв успешно удален",
	})
}

func (h *ReviewsHandler) RegisterRoutes(r *gin.Engine) {

	reviews := r.Group("/reviews")
	{
		reviews.POST("", h.CreateReview)
		reviews.GET("/:id", h.GetReview)
		reviews.GET("/user/:userID", h.GetReviewsByUser)
		reviews.GET("/category/:categoryID", h.GetReviewsByCategory)
		reviews.PUT("/:id", h.UpdateReview)
		reviews.DELETE("/:id", h.DeleteReview)
	}
}
