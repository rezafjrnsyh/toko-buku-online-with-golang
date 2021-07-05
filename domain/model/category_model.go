package domain

import "main/utils"

type Category struct {
	Id		int	`json:"id"`
	NameCategory string	`json:"nameCategory"`
}

type ICategoryRepository interface {
	Find() ([]*Category, utils.MessageErr)
	Create(category *Category) (*Category, utils.MessageErr)
	FindById(id int) (*Category, utils.MessageErr)
	Update(category *Category) (*Category, utils.MessageErr)
	Delete(id int) (int64,utils.MessageErr)
}

type ICategoryService interface {
	FindCategory() ([]*Category, utils.MessageErr)
	CreateCategory(category *Category) (*Category, utils.MessageErr)
	FindCategoryById(id int) (*Category, utils.MessageErr)
	UpdateCategory(category *Category, id int) (*Category, utils.MessageErr)
	DeleteCategory(id int) (int64,utils.MessageErr)
}
