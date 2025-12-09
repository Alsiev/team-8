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
	mealPlan service.MealPlanService,
	mealPlanItem service.MealPlanItemsService,
	user service.UserService,
) {

	categoryHandler := NewCategoryHandler(category, log)
	planHandler := NewExercisePlanHandler(plan, log)
	bmiHand := NewBmiHandler(log)
	userHandler := NewUserHandler(user, log)
	mealPlanHandler := NewMealPlanHandler(mealPlan, log)
	mealPlanItemHandler := NewMealPlanItemHandler(mealPlanItem, log)
	mealPlanHandler.RegisterRoutes(router)
	mealPlanItemHandler.RegisterRoutes(router)
	categoryHandler.RegisterRoutes(router)
	planHandler.RegisterRoutes(router)
	bmiHand.RegisterRoutes(router)
	userHandler.UserRoutes(router)

}
