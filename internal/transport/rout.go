package transport

import (
	"healthy_body/internal/service"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	log *slog.Logger,
	category service.CategoryServices,
	plan service.ExercisePlanServices,
	user service.UserService,
) {

	categoryHandler := NewCategoryHandler(category, log)
	planHandler := NewExercisePlanHandler(plan,log)
	bmiHand:= NewBmiHandler(log)
	categoryHandler.RegisterRoutes(router)
	planHandler.RegisterRoutes(router)
	bmiHand.RegisterRoutes(router)
}

