package main

import (
	"fmt"
	"healthy_body/internal/config"
	"healthy_body/internal/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	db := config.SetUpDatabaseConnection()
	server := gin.Default()
	
	if err:= db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.ExercisePlan{},
		&models.ExercisePlanItem{},
		&models.MealPlan{},
		&models.MealPlanItem{},); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}

	if tableList, err := db.Migrator().GetTables(); err == nil {
		fmt.Println("tables:", tableList)
	}


	if err := server.Run(); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
