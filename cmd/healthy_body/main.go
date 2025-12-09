package main

import (
	"fmt"
	"healthy_body/internal/config"
	"healthy_body/internal/models"
	"healthy_body/internal/repository"
	"healthy_body/internal/service"
	"healthy_body/internal/transport"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.SetUpDatabaseConnection()
	server := gin.Default()

	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.ExercisePlan{},
		&models.ExercisePlanItem{},
		&models.MealPlan{},
		&models.MealPlanItem{}); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	categoryRepo := repository.NewCategoryRepo(db, logger)
	planRepo := repository.NewExercisePlanRepo(db, logger)

	mealPlanRepo := repository.NewMealPlanRepository(db, logger)
	mealPlanItemRepo := repository.NewMealPlanItemRepository(db, logger)

	categoryServices := service.NewCategoryServices(categoryRepo, logger)
	planServices := service.NewExercisePlanServices(planRepo, logger, categoryServices)
	mealPlanService := service.NewMealPlanService(mealPlanRepo, logger, categoryServices)
	mealPlanItemService := service.NewMealPlanItemsService(mealPlanItemRepo, logger)
	userRepo := repository.NewUserRepository(db, logger)
	userService := service.NewUserService(userRepo, logger, db)


	if tableList, err := db.Migrator().GetTables(); err == nil {
		fmt.Println("tables:", tableList)
	}

	transport.RegisterRoutes(
		server,
		logger,
		categoryServices,
		planServices,
		mealPlanService,
		mealPlanItemService,
		userService,
	)

	if err := server.Run(); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
