package services

import (
	"database/sql"
	"fmt"
	domain "main/domain/model"
	"main/domain/repositories"
	"main/utils"
)

type categoryService struct {
	db          *sql.DB
	CategoryRepo domain.ICategoryRepository
}

func NewCategoryService(db *sql.DB) domain.ICategoryService {
	return &categoryService{db: db, CategoryRepo: repositories.NewCategoryRepo(db)}
}

func (c categoryService) FindCategory() ([]*domain.Category, utils.MessageErr) {
	articles,err := c.CategoryRepo.Find()
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (c categoryService) CreateCategory(category *domain.Category) (*domain.Category, utils.MessageErr) {
	category, err := c.CategoryRepo.Create(category)
	fmt.Println("Service :", category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (c categoryService) FindCategoryById(id int) (*domain.Category, utils.MessageErr) {
	category, err := c.CategoryRepo.FindById(id)
	if err != nil{
		return nil, err
	}

	return category, nil
}

func (c categoryService) UpdateCategory(category *domain.Category, id int) (*domain.Category, utils.MessageErr) {
	data,errfound := c.FindCategoryById(id)
	if errfound != nil {
		return nil, errfound
	}

	data.NameCategory = category.NameCategory

	NewCategory, err := c.CategoryRepo.Update(data)

	if err != nil{
		return nil, err
	}

	return NewCategory, nil

}

func (c categoryService) DeleteCategory(id int) (int64, utils.MessageErr) {
	category, err := c.CategoryRepo.FindById(id)
	if err != nil{
		return 0, err
	}

	result, err := c.CategoryRepo.Delete(category.Id)
	if err != nil {
		return result, err
	}
	return result, nil
}
