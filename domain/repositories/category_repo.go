package repositories

import (
	"database/sql"
	"fmt"
	"log"
	domain "main/domain/model"
	"main/utils"
	"strings"
)

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) domain.ICategoryRepository {
	return &categoryRepo{db: db}
}

func (c *categoryRepo) Find() ([]*domain.Category, utils.MessageErr) {
	// Membuat object slice category
	articles := make([]*domain.Category, 0)
	//defer a.db.Close()


	// Untuk format query
	query := fmt.Sprintf(`SELECT id, name_category FROM category`)

	// Eksekusi query
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, utils.ParseError(err)
	}
	defer rows.Close()

	for rows.Next() {
		category := &domain.Category{}
		getError := rows.Scan(&category.Id,&category.NameCategory)
		if err != nil {
			return nil, utils.NewInternalServerError(fmt.Sprintf("Error when trying to get message: %s", getError.Error()))
		}
		articles = append(articles, category)
	}
	if len(articles) == 0 {
		return nil, utils.NewNotFoundError("no records found")
	}
	return articles, nil
}

func (c *categoryRepo) Create(category *domain.Category) (*domain.Category, utils.MessageErr) {
	query := fmt.Sprintf(`INSERT INTO category(name_category) VALUES(?)`)
	result, err := c.db.Exec(query, category.NameCategory)
	if err != nil {
		s := strings.Split(err.Error(), ":")
		log.Println(s[1])
		return nil, utils.NewInternalServerError(fmt.Sprintf("error when trying to prepare user to save: %s", err.Error()))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, utils.NewInternalServerError(fmt.Sprintf("error when trying to save message: %s", err.Error()))
	}

	category.Id = int(id)

	return category, nil
}

func (c *categoryRepo) FindById(id int) (*domain.Category, utils.MessageErr) {
	category := new(domain.Category)

	query := fmt.Sprintf(`SELECT id, name_category FROM category WHERE id=?`)
	if getError := c.db.QueryRow(query, id).
		Scan(&category.Id,&category.NameCategory); getError != nil {
		fmt.Println("this is the error man: ", getError)
		return nil,  utils.ParseError(getError)
	}
	return category, nil
}

func (c *categoryRepo) Update(category *domain.Category) (*domain.Category, utils.MessageErr) {
	query := fmt.Sprintf("UPDATE category SET name_category = ? WHERE id = ?")
	_, updateErr := c.db.Exec(query,&category.NameCategory, &category.Id, )
	if updateErr != nil {
		s := strings.Split(updateErr.Error(), ":")
		log.Println(s[1])
		if updateErr != nil {
			return nil,  utils.ParseError(updateErr)
		}
	}

	return category, nil
}

func (c *categoryRepo) Delete(id int) (int64, utils.MessageErr) {
	query := fmt.Sprintf("DELETE FROM category WHERE id = ?")
	result, err := c.db.Exec(query, id)
	if err != nil {
		return 0, utils.NewInternalServerError(fmt.Sprintf("error when trying to delete message %s", err.Error()))
	}
	RowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, utils.NewInternalServerError(fmt.Sprintf("error when trying to get rows affected %s", err.Error()))
	}
	return RowsAffected,nil
}
