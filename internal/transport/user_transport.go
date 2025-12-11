package transport

import (
	"healthy_body/internal/models"
	"healthy_body/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	user service.UserService
	log  *slog.Logger
}

func NewUserHandler(user service.UserService, log *slog.Logger) *UserHandler {
	return &UserHandler{user: user, log: log}
}

// Create godoc
// @Summary Создать пользователя
// @Description Создает нового пользователя
// @Tags User
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "Данные пользователя"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/ [post]
func (h *UserHandler) Create(c *gin.Context) {
	var user models.CreateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		h.log.Warn("Введены неверные данные", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат данных", "error": err.Error()})
		return
	}

	result, err := h.user.CreateUser(user)
	if err != nil {
		h.log.Error("Ошибка при создании пользователя", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "ошибка при создании пользователя"})
		return
	}

	h.log.Info("Пользователь создан", "имя", user.Name)
	c.JSON(http.StatusCreated, gin.H{"message": "пользователь создан", "users": result})
}

// GetAllUser godoc
// @Summary Получить всех пользователей
// @Description Возвращает список всех пользователей
// @Tags User
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /user/ [get]
func (h *UserHandler) GetAllUser(c *gin.Context) {
	result, err := h.user.GetAllUsers()
	if err != nil {
		h.log.Error("Ошибка при выводе всех пользователей", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	h.log.Info("Пользователи получены", "всего пользователей", len(result))
	c.JSON(http.StatusOK, gin.H{"users": result, "total": len(result)})
}

// GetUserByID godoc
// @Summary Получить пользователя по ID
// @Description Возвращает пользователя по указанному ID
// @Tags User
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("Некорректный ID")
		c.JSON(http.StatusBadRequest, gin.H{"message": "некорректный ID"})
		return
	}

	result, err := h.user.GetUserByID(uint(id))
	if err != nil {
		h.log.Error("Ошибка при поиске пользователя по ID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Ошибка при поиске пользователя по ID"})
		return
	}

	h.log.Info("Пользователь найден", "пользователь", result)
	c.JSON(http.StatusOK, result)
}

// Update godoc
// @Summary Обновить пользователя
// @Description Обновляет данные пользователя по ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param user body models.UpdateUserRequest true "Данные для обновления"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /user/{id} [patch]
func (h *UserHandler) Update(c *gin.Context) {
	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Введены неверные данные", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат данных", "error": err.Error()})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("Некорректный ID")
		c.JSON(http.StatusBadRequest, gin.H{"message": "некорректный ID"})
		return
	}

	result, err := h.user.UpdateUser(uint(id), req)
	if err != nil {
		h.log.Error("ошибка при обновлении пользователя", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Пользователь обновлен")
	c.JSON(http.StatusOK, result)
}

// Delete godoc
// @Summary Удалить пользователя
// @Description Удаляет пользователя по ID
// @Tags User
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /user/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Warn("Некорректный ID")
		c.JSON(http.StatusBadRequest, gin.H{"message": "некорректный ID"})
		return
	}

	if err := h.user.Delete(uint(id)); err != nil {
		h.log.Error("Ошибка при удалении пользователя", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Пользователь удален")
	c.JSON(http.StatusOK, gin.H{"message": "пользователь удален"})
}

// Payment godoc
// @Summary Оплата пользователем
// @Description Пользователь оплачивает выбранную категорию
// @Tags User
// @Produce json
// @Param userID path int true "ID пользователя"
// @Param categoryID path int true "ID категории"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /user/payment/{userID}/{categoryID} [post]
func (h *UserHandler) Payment(c *gin.Context) {
	userIDstr := c.Param("userID")
	categoryIDstr := c.Param("categoryID")

	userID, err := strconv.ParseUint(userIDstr, 10, 64)

	if err != nil {
		h.log.Error("Ошибка при получении ID пользователя")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	categoryID, err := strconv.ParseUint(categoryIDstr, 10, 64)

	if err != nil {
		h.log.Error("Ошибка при получении ID категории")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	if err := h.user.Payment(uint(userID), uint(categoryID)); err != nil {
		h.log.Error("Ошибка при оплате",
			"error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	h.log.Info("Оплата прошла успешно")
	c.JSON(http.StatusOK, gin.H{
		"message": "оплата прошла успешно",
	})
}

// PaymentToAnother godoc
// @Summary Оплата другому пользователю
// @Description Пользователь оплачивает категорию другому пользователю
// @Tags User
// @Produce json
// @Param userID path int true "ID пользователя"
// @Param categoryID path int true "ID категории"
// @Param secondUserID path int true "ID второго пользователя"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /user/present/{userID}/{categoryID}/{secondUserID} [post]
func (h *UserHandler) PaymentToAnother(c *gin.Context) {
	userIDstr := c.Param("userID")
	categoryIDstr := c.Param("categoryID")
	secondUserIDstr := c.Param("secondUserID")

	userID, err := strconv.ParseUint(userIDstr, 10, 64)

	if err != nil {
		h.log.Error("Ошибка при получении ID второго пользователя")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	categoryID, err := strconv.ParseUint(categoryIDstr, 10, 64)

	if err != nil {
		h.log.Error("Ошибка при получении ID категории")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	secondUserID, err := strconv.ParseUint(secondUserIDstr, 10, 64)

	if err != nil {
		h.log.Error("Ошибка при получении ID пользователя")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	if err := h.user.PaymentToAnother(uint(userID), uint(categoryID), uint(secondUserID)); err != nil {
		h.log.Error("Ошибка при оплате",
			"error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	h.log.Info("Оплата прошла успешно")
	c.JSON(http.StatusOK, gin.H{
		"message": "оплата прошла успешно",
	})
}

// GetUserWithPlan godoc
// @Summary Получить пользователя с планом питания
// @Description Возвращает пользователя вместе с его планом питания
// @Tags User
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /user/plan/{id} [get]
func (h *UserHandler) GetUserWithPlan(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		h.log.Warn("Некорректный ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "некорректный ID",
		})
		return
	}

	user, err := h.user.GetUserPlan(uint(id))
	if err != nil {
		h.log.Error("Ошибка при удалении пользователя",
			"error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

// GetUserCategory godoc
// @Summary Получить категорию/планы пользователя
// @Description Возвращает пользователя с предзагруженными планами и категориями
// @Tags User
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /user/userplans/{id} [get]
func (h *UserHandler) GetUserCategory(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		h.log.Warn("Ошибка при вводе ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	user, err := h.user.GetUserCategory(uint(userID))
	if err != nil {
		h.log.Error("Ошибка при получении ID пользователя")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

// GetUserSubs godoc
// @Summary Получить подписки пользователя
// @Description Возвращает все подписки пользователя
// @Tags User
// @Produce json
// @Param userID path int true "ID пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /user/usersub/{userID} [get]
func (h *UserHandler) GetUserSubs(c *gin.Context) {
	userIDstr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDstr, 10, 64)

	if err != nil {
		h.log.Warn("Ошибка при вводе ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	user, err := h.user.GetUserSub(uint(userID))
	if err != nil {
		h.log.Error("Ошибка при получении ID пользователя")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

// SubPayment godoc
// @Summary Оплата подписки пользователем
// @Description Пользователь оплачивает выбранную подписку
// @Tags User
// @Produce json
// @Param userID path int true "ID пользователя"
// @Param subID path int true "ID подписки"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /user/sub/{userID}/{subID} [post]
func (h *UserHandler) SubPayment(c *gin.Context) {
	userIDstr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDstr, 10, 64)
	if err != nil {
		h.log.Warn("Ошибка при вводе ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	subIDstr := c.Param("subID")
	subID, err := strconv.ParseUint(subIDstr, 10, 64)
	if err != nil {
		h.log.Warn("Ошибка при вводе ID подписки")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}
	if err := h.user.SubPayment(uint(userID), uint(subID)); err != nil {
		h.log.Error("Ошибка при оплате подписки")
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "подписка успешно оформлена",
	})
}

func (h *UserHandler) UserRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/", h.Create)
		userGroup.POST("/payment/:userID/:categoryID", h.Payment)
		userGroup.POST("/present/:userID/:categoryID/:secondUserID", h.PaymentToAnother)
		userGroup.POST("sub/:userID/:subID", h.SubPayment)
		userGroup.GET("/", h.GetAllUser)
		userGroup.GET("/:id", h.GetUserByID)
		userGroup.GET("/plan/:id", h.GetUserWithPlan)
		userGroup.GET("/userplans/:id", h.GetUserCategory)
		userGroup.GET("/usersub/:userID", h.GetUserSubs)
		userGroup.PATCH("/:id", h.Update)
		userGroup.DELETE("/:id", h.Delete)
	}
}
