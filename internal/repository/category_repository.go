package repository

import (
	"errors"
	"healthy_body/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type CategoryRepo interface {
	 Create(category *models.Categories) error
	 List() ([]models.Categories, error)
	 GetByID(id uint) (*models.Categories,error)
	 GetWithPlans(id uint) (*models.Categories, error)
	 Update(category *models.Categories) error
	 Delete(id uint) error
}

type categoryRepo struct {
	db *gorm.DB
	log *slog.Logger
}

func NewCategoryRepo(db *gorm.DB, log *slog.Logger) CategoryRepo{
	return  &categoryRepo{
		db: db,
		log: log,
	}
}


func (c *categoryRepo) Create(category *models.Categories) error {
	 if category == nil {
		c.log.Error("error in Create function category_repository.go")
		return  errors.New("error create category in db")
	 }

	 return  c.db.Create(category).Error
}


func (c *categoryRepo) List() ([]models.Categories, error){
	var list []models.Categories
	if err:= c.db.Find(&list).Error; err != nil {
		c.log.Error("error in List function category_repository.go")
		return nil, err
	}

	return  list, nil
}


func (c *categoryRepo) GetByID(id uint) (*models.Categories,error) {
	var category models.Categories
	if err := c.db.First(&category,id).Error; err != nil {
		c.log.Error("error in GetByID function category_repository.go")
		return nil, err
	}

	return  &category, nil
}


func (c *categoryRepo) GetWithPlans(id uint) (*models.Categories, error) {
    var category models.Categories

    err := c.db.Preload("ExercisePlans.Exercises").Preload("MealPlans.Meals").First(&category, id).Error

    if err != nil {
        c.log.Error("error in GetWithPlans function category_repository.go", "err", err)
        return nil, err
    }

    return &category, nil
}


func (c *categoryRepo) Update(category *models.Categories) error {
	if category == nil {
		c.log.Error("error in Update function category_repository.go")
		return errors.New("error update in db") 
	}

	return  c.db.Save(category).Error
}


func (c *categoryRepo) Delete(id uint) error {
	if err := c.db.Delete(&models.Categories{}, id).Error; err != nil {
		c.log.Error("error in Delete function category_repository.go")
		return errors.New("error delete in db") 
	}

	return  nil 
}