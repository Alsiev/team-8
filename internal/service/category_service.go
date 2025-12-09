package service

import (
	"errors"
	"healthy_body/internal/models"
	"healthy_body/internal/repository"
	"log/slog"
)

type CategoryServices interface {
	CreateCategory(req models.CreateCategoryRequest) (*models.Category, error)
	GetCategoryList()([]models.Category,error)
	GetCategoryByID(id uint) (*models.Category,error)
	UpdateCategory(id uint, req models.UpdateCategoryRequest) (*models.Category, error)
	DeleteCategory(id uint) error
}

type categoryServices struct {
	category repository.CategoryRepo
	log *slog.Logger
}


func NewCategoryServices(category repository.CategoryRepo, log *slog.Logger) CategoryServices{
	return  &categoryServices{
		category: category,
		log: log,
	}
}

func (c *categoryServices) CreateCategory(req models.CreateCategoryRequest) (*models.Category, error) {
	if req.Name == ""{
		return  nil , errors.New("empty name by your category")
	}

	if req.Price == 0{
		return  nil , errors.New("empty valid price your category")
	}

	if req.Description == ""{
		return  nil , errors.New("empty description by your category")
	}

	 category := &models.Category{
		Name: req.Name,
		Description: req.Description,
		Price: req.Price,
	 }

	  if err:= c.category.Create(category); err != nil {
		c.log.Error("error Create in category_service.go")
		return nil, err
	  }

	 return  category, nil
}

func (c *categoryServices) GetWithPlans(id uint) (*models.Category, error) {
    cat , err := c.category.GetWithPlans(id)
	if err != nil {
		c.log.Error("error preloads or id")
		return nil ,err
	}
   
	return cat, nil
}


func (c *categoryServices) GetCategoryList()([]models.Category,error){
	list , err := c.category.List()
	if err != nil {
		c.log.Error("error GetList in category_service.go")
		return nil, err
	}


	return  list , nil
}

func (c *categoryServices) GetCategoryByID(id uint) (*models.Category,error) {
	category, err := c.category.GetByID(id)
	if err != nil {
		c.log.Error("error GetCategoryByID in category_service.go")
		return  nil ,err
	}

	return  category, nil
}

func (c *categoryServices) UpdateCategory(id uint, req models.UpdateCategoryRequest) (*models.Category, error){
	category ,err :=  c.category.GetByID(id)
	if err != nil {
		c.log.Error("error UpdateCategory function in category_service.go")
		return &models.Category{} , err
	}

	c.Up(category, req)

	if err := c.category.Update(category); err != nil {
		c.log.Error("error UpdateCategory in category_service.go")
		return &models.Category{}, nil
	}

	return  category, nil
}


func (c *categoryServices) DeleteCategory(id uint) error{
	if err:= c.category.Delete(id); err != nil {
		c.log.Error("error UpdateCategory in category_service.go")
		return err
	}

	return  nil
}

func(c *categoryServices) Up(cat *models.Category, req models.UpdateCategoryRequest){
	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.Description != nil {
		cat.Description = *req.Description
	}
	if req.Price != nil {
		cat.Price = *req.Price
	}
}