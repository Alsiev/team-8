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

func (h *UserHandler) Create(c *gin.Context) {
	var user models.CreateUserRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		h.log.Warn("Введены неверные данные",
			"err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Неверный формат данных",
			"error":   err.Error()})
		return
	}

	result, err := h.user.CreateUser(user)

	if err != nil {
		h.log.Error("Ошибка при создании пользователя",
			"error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "ошибка при создании пользователя",
		})
		return
	}

	h.log.Info("Пользователь создан",
		"имя", user.Name)

	c.JSON(http.StatusCreated, gin.H{
		"message:": "пользователь создан",
		"users":    result,
	})
}

func (h *UserHandler) GetAllUser(c *gin.Context) {

	result, err := h.user.GetAllUsers()

	if err != nil {
		h.log.Error("Ошибка при выводе всех пользователей",
			"error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	h.log.Info("Пользователи получены",
		"пользователи", result,
		"всего пользователей", len(result))

	c.JSON(http.StatusOK, gin.H{
		"users": result,
		"total": len(result),
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		h.log.Warn("Некорректный ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "некорректный ID",
		})
		return
	}

	result, err := h.user.GetUserByID(uint(id))

	if err != nil {
		h.log.Error("Ошибка при поиске пользователя по ID")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Ошибка при поиске пользователя по ID",
		})
		return
	}

	h.log.Info("Пользователь найден",
		"пользователь", result)

	c.JSON(http.StatusOK, result)

}

func (h *UserHandler) Update(c *gin.Context) {
	var req models.UpdateUserRequest

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
		h.log.Warn("Некорректный ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "некорректный ID",
		})
		return
	}

	result, err := h.user.UpdateUser(uint(id), req)

	if err != nil {
		h.log.Error("ошибка при обновлении пользователя",
			"error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.log.Info("Пользователь обновлен")
	c.JSON(http.StatusOK, result)

}

func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		h.log.Warn("Некорректный ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "некорректный ID",
		})
		return
	}

	if err := h.user.Delete(uint(id)); err != nil {
		h.log.Error("Ошибка при удалении пользователя",
			"error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.log.Info("Пользователь удален")

	c.JSON(http.StatusOK, gin.H{
		"message": "пользователь удален",
	})

}

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

	h.log.Info("Оплата проошла успешно")
	c.JSON(http.StatusOK, gin.H{
		"message": "оплата прощла успешно",
	})
}

func (h *UserHandler) GetUserWithiPlan(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		h.log.Warn("Некорректный ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "некорректный ID",
		})
		return
	}

	user , err := h.user.GetUserPlan(uint(id))
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

func (h *UserHandler) UserRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/", h.Create)
		userGroup.POST("/buycategory/:userID/:categoryID", h.Payment)
		userGroup.GET("/", h.GetAllUser)
		userGroup.GET("/:id", h.GetUserByID)
		userGroup.GET("/plan/:id", h.GetUserWithiPlan)
		userGroup.PATCH("/:id", h.Update)
		userGroup.DELETE("/:id", h.Delete)
	}
}
