package repository

import (
	"errors"
	"healthy_body/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type CategoryRepo interface {
	 Create(category *models.Category) error
	 List() ([]models.Category, error)
	 GetByID(id uint) (*models.Category,error)
	 Update(category *models.Category) error
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


func (c *categoryRepo) Create(category *models.Category) error {
	 if category == nil {
		c.log.Error("error in Create function category_repository.go")
		return  errors.New("error create category in db")
	 }

	 return  c.db.Create(category).Error
}


func (c *categoryRepo) List() ([]models.Category, error){
	var list []models.Category
	if err:= c.db.Find(&list).Error; err != nil {
		c.log.Error("error in List function category_repository.go")
		return nil, err
	}

	return  list, nil
}


func (c *categoryRepo) GetByID(id uint) (*models.Category,error) {
	var category models.Category
	if err := c.db.First(&category,id).Error; err != nil {
		c.log.Error("error in GetByID function category_repository.go")
		return nil, err
	}

	return  &category, nil
}


func (c *categoryRepo) Update(category *models.Category) error {
	if category == nil {
		c.log.Error("error in Update function category_repository.go")
		return errors.New("error update in db") 
	}

	return  c.db.Save(category).Error
}


func (c *categoryRepo) Delete(id uint) error {
	if err := c.db.Delete(&models.Category{}, id).Error; err != nil {
		c.log.Error("error in Delete function category_repository.go")
		return errors.New("error delete in db") 
	}

	return  nil 
}