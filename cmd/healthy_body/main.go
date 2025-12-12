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
	"time"

	_ "healthy_body/internal/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	db := config.SetUpDatabaseConnection()
	server := gin.Default()

	// 🚀 ВКЛЮЧАЕМ CORS — ЭТО ГЛАВНОЕ
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	if err := db.AutoMigrate(
		&models.Categories{},
		&models.Subscription{},
		&models.User{},
		&models.UserPlan{},
		&models.UserSubscription{},
		&models.ExercisePlan{},
		&models.ExercisePlanItem{},
		&models.MealPlan{},
		&models.MealPlanItem{},
		&models.Reviews{},
	); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	categoryRepo := repository.NewCategoryRepo(db, logger)
	planRepo := repository.NewExercisePlanRepo(db, logger)
	mealPlanRepo := repository.NewMealPlanRepository(db, logger)
	mealPlanItemRepo := repository.NewMealPlanItemRepository(db, logger)
	subRepo := repository.NewSubscriptionRepo(db, logger)
	reviewsRepo := repository.NewReviewsRepository(db, logger)

	categoryServices := service.NewCategoryServices(categoryRepo, logger)
	planServices := service.NewExercisePlanServices(planRepo, logger, categoryServices)
	mealPlanService := service.NewMealPlanService(mealPlanRepo, logger, categoryServices)
	mealPlanItemService := service.NewMealPlanItemsService(mealPlanItemRepo, logger)
	userRepo := repository.NewUserRepository(db, logger)
	subService := service.NewSubscriptionService(subRepo, logger, categoryServices)
	notificationService := service.NewEmailNotificationService(
		os.Getenv("EMAIL_USER"),
		os.Getenv("EMAIL_PASS"),
		os.Getenv("EMAIL_HOST"),
		587,
		logger)
	userService := service.NewUserService(userRepo, logger, db, subService, categoryRepo, notificationService)
	reviewsService := service.NewReviewsService(reviewsRepo, logger)

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
		subService,
		reviewsService,
	)

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := server.Run(":8888"); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
